package meta

// Metadata represents a metadata.json file
type Metadata struct {
	ID                string   `json:"@id"`
	Typ               VarStr   `json:"@type"`
	Title             string   `json:"title"`
	Created           *W3CDate `json:"created,omitempty"`
	Creator           Agent    `json:"creator,omitempty"`
	Source            VarStr   `json:"source,omitempty"`
	Series            string   `json:"series,omitempty"`
	Consignment       string   `json:"consignment,omitempty"`
	DisposalRule      Disposal `json:"disposalRule,omitempty"`
	Duration          string   `json:"duration,omitempty"`
	Language          VarStr   `json:"language,omitempty"`
	Subtitles         VarStr   `json:"subtitles,omitempty"`
	Director          VarStr   `json:"director,omitempty"`
	Actor             VarStr   `json:"actor,omitempty"`
	ProductionCompany VarStr   `json:"productionCompany,omitempty"`
	About             Thing    `json:"about,omitempty"`
	Context           Context  `json:"@context"`
}

// VarStr represents a single string or a slice of strings
type VarStr interface{}

func copyVarStr(v VarStr) VarStr {
	if v == nil {
		return nil
	}
	if str, ok := v.(string); ok {
		return str
	}
	strs := v.([]string)
	ret := make([]string, len(strs))
	copy(ret, strs)
	return ret
}

// Disposal can be a single DisposalRule{} or a slice of []DisposalRule{}
type Disposal interface{}

type DisposalRule struct {
	Authority string `json:"authority"`
	Class     string `json:"class"`
}

// Agent can be a string e.g. "Richard Lehane" or an object with @id/@type, or a slice of Agents
type Agent interface{}

// Thing can be anything that a metadata is "about"
type Thing interface{}

// Business is a type of Thing. It is used for the BRS project
type Business struct {
	Typ                string   `json:"@type,omitempty"`
	LegalName          string   `json:"legalName,omitempty"`
	CommencedTrading   *W3CDate `json:"commencedTrading,omitempty"`
	CeasedTrading      *W3CDate `json:"ceasedTrading,omitempty"`
	RenewalDueDate     *W3CDate `json:"renewalDueDate,omitempty"`
	RegistrationNumber string   `json:"registrationNumber,omitempty"`
	ABN                string   `json:"abn,omitempty"`
	Proprietor         Agent    `json:"proprietor,omitempty"`
}

// NewMetadata returns a Metadata with the supplied title. It also sets the @type.
func NewMetadata(id int, title string) *Metadata {
	return &Metadata{
		ID:    ReferenceObject(id),
		Typ:   "http://records.nsw.gov.au/terms/DigitalArchive",
		Title: title,
	}
}

// Metadata can have multiple types e.g. both an DigitalArchive and a Movie
func (m *Metadata) AddType(typ string) {
	if str, ok := m.Typ.(string); ok {
		m.Typ = []string{str, typ}
		return
	}
	strs := m.Typ.([]string)
	m.Typ = append(strs, typ)
	return
}

// MakeAgent returns an Agent with the given name, @id and @type.
// If @id and @type aren't given, the Agent is a simple string e.g. "Richard Lehane"
func MakeAgent(name, id, typ string) Agent {
	if id == "" && typ == "" {
		return name
	}
	return Obj{
		ID:   id,
		Typ:  typ,
		Name: name,
	}
}

// MakeSDOPerson creates an Agent that is of @type schema.org/Person. Does not set an @id.
func MakeSDOPerson(name string) Agent {
	return MakeAgent(name, "", "http://schema.org/Person")
}

// MakePerson creates an Agent that is of @type terms/Person. Sets the @id to the supplied value.
func MakePerson(name string, id int) Agent {
	return MakeAgent(name, ToID(id, "http://records.nsw.gov.au/persons/"), "http://records.nsw.gov.au/terms/Person")
}

// MakePerson creates an Agent that is of @type terms/Agency. Sets the @id to the supplied value.
func MakeAgency(name string, id int) Agent {
	return MakeAgent(name, ToID(id, "http://records.nsw.gov.au/agencies/"), "http://records.nsw.gov.au/terms/Agency")
}

// MakePerson creates an Agent that is of @type schema.org/Organization.
func MakeOrganization(name string) Agent {
	return MakeAgent(name, "", "http://schema.org/Organization")
}

// ToSeries turns a series number into an IRI @id
func ToSeries(i int) string {
	return ToID(i, "http://records.nsw.gov.au/series/")
}

// ToConsignment turns a consignment number into an IRI @id
func ToConsignment(i int) string {
	return ToID(i, "http://records.nsw.gov.au/consignment/")
}

// MakeBusiness returns a Thing of @type schema.org/Organization and sets the supplied fields.
// A variable number of proprietors can be supplied and these will be set as a slice of Organizations.
func MakeBusiness(legalName, commencedTrading, ceasedTrading, renewalDueDate, registrationNumber, abn string, proprietors ...string) Thing {
	var props Agent
	if len(proprietors) > 0 {
		if len(proprietors) > 1 {
			props := make([]Agent, len(proprietors))
			for i, v := range proprietors {
				props[i] = MakeOrganization(v)
			}
		} else {
			props = MakeOrganization(proprietors[0])
		}
	}
	return Business{
		Typ:                "http://schema.org/Organization",
		LegalName:          legalName,
		CommencedTrading:   NewDate(commencedTrading),
		CeasedTrading:      NewDate(ceasedTrading),
		RenewalDueDate:     NewDate(renewalDueDate),
		RegistrationNumber: registrationNumber,
		ABN:                abn,
		Proprietor:         props,
	}

}

// ReferenceObject makes a temporary reference to another object in the consignment.
// This reference is swapped for a UUID by the migrate tool.
func ReferenceObject(i int) string {
	return ToRef(i, "obj")
}

// ReferenceMigration makes a temporary reference to another migration (that doesn't yet have a UUID).
// This reference is swapped for a UUID by the migrate tool.
func ReferenceMigration(i int) string {
	return ToRef(i, "mig")
}

var metadataContext = Context{
	"abn":              "https://www.wikidata.org/wiki/Q4823913",
	"about":            "http://schema.org/about",
	"actor":            "http://schema.org/actor",
	"authority":        "http://records.nsw.gov.au/terms/disposalAuthority",
	"ceasedTrading":    "http://schema.org/dissolutionDate",
	"class":            "http://records.nsw.gov.au/terms/disposalClass",
	"commencedTrading": "http://schema.org/foundingDate",
	"consignment": Obj{
		ID:  "http://records.nsw.gov.au/terms/consignment",
		Typ: "@id",
	},
	"created": Obj{
		ID:  "http://purl.org/dc/terms/created",
		Typ: "http://www.w3.org/2001/XMLSchema#date",
	},
	"creator":  "http://purl.org/dc/terms/creator",
	"director": "http://schema.org/director",
	"disposalRule": Obj{
		ID:  "http://records.nsw.gov.au/terms/disposalRule",
		Typ: "http://records.nsw.gov.au/terms/DisposalRule",
	},
	"duration":           "http://schema.org/duration",
	"language":           "http://schema.org/inLanguage",
	"legalName":          "http://schema.org/legalName",
	"name":               "http://schema.org/name",
	"productionCompany":  "http://schema.org/productionCompany",
	"proprietor":         "http://records.nsw.gov.au/terms/proprietor",
	"registrationNumber": "http://records.nsw.gov.au/terms/registrationNumber",
	"renewalDueDate": Obj{
		ID:  "http://records.nsw.gov.au/terms/renewalDueDate",
		Typ: "http://www.w3.org/2001/XMLSchema#date",
	},
	"series": Obj{
		ID:  "http://records.nsw.gov.au/terms/series",
		Typ: "@id",
	},
	"source":    "http://purl.org/dc/terms/source",
	"subtitles": "http://schema.org/subtitleLanguage",
	"title":     "http://purl.org/dc/terms/title",
}

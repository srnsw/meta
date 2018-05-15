package meta

// Metadata represents a metadata.json file
type Metadata struct {
	Typ               []string       `json:"@type"`
	Title             string         `json:"title"`
	Created           *W3CDate       `json:"created,omitempty"`
	Creators          []Agent        `json:"creators,omitempty"`
	Source            string         `json:"source,omitempty"`
	Series            string         `json:"series,omitempty"`
	Consignment       string         `json:"consignment,omitempty"`
	DisposalRules     []DisposalRule `json:"disposalRules,omitempty"`
	Duration          string         `json:"duration,omitempty"`
	Language          string         `json:"language,omitempty"`
	Subtitles         string         `json:"subtitles,omitempty"`
	Director          string         `json:"director,omitempty"`
	Actor             string         `json:"actor,omitempty"`
	ProductionCompany string         `json:"productionCompany,omitempty"`
	About             Thing          `json:"about,omitempty"`
	Context           Context        `json:"@context"`
}

// DisposalRule represents disposalRules
type DisposalRule struct {
	Authority string `json:"authority"`
	Class     string `json:"class"`
}

// Agent can be a string e.g. "Richard Lehane" or an object with @id/@type
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
	Proprietors        []Agent  `json:"proprietors,omitempty"`
}

// NewMetadata returns a Metadata with the supplied title. It also sets the @type.
func NewMetadata(title string) *Metadata {
	return &Metadata{
		Typ:   []string{"http://www.records.nsw.gov.au/terms/DigitalArchive"},
		Title: title,
	}
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

// MakePerson creates an Agent that is of @type terms/Person. Sets the @id to the supplied value.
func MakePerson(name string, id int) Agent {
	return MakeAgent(name, ToID(id, "http://www.records.nsw.gov.au/persons/"), "http://www.records.nsw.gov.au/terms/Person")
}

// MakePerson creates an Agent that is of @type terms/Agency. Sets the @id to the supplied value.
func MakeAgency(name string, id int) Agent {
	return MakeAgent(name, ToID(id, "http://records.nsw.gov.au/agencies/"), "http://www.records.nsw.gov.au/terms/Agency")
}

// MakePerson creates an Agent that is of @type schema.org/Organization.
func MakeOrganization(name string) Agent {
	return MakeAgent(name, "", "http://schema.org/Organization")
}

// ToSeries turns a series number into an IRI @id
func ToSeries(i int) string {
	return ToID(i, "http://www.records.nsw.gov.au/series/")
}

// ToConsignment turns a consignment number into an IRI @id
func ToConsignment(i int) string {
	return ToID(i, "http://www.records.nsw.gov.au/consignment/")
}

// MakeBusiness returns a Thing of @type schema.org/Organization and sets the supplied fields.
// A variable number of proprietors can be supplied and these will be set as a slice of Organizations.
func MakeBusiness(legalName, commencedTrading, ceasedTrading, renewalDueDate, registrationNumber, abn string, proprietors ...string) Thing {
	props := make([]Agent, len(proprietors))
	for i, v := range proprietors {
		props[i] = MakeOrganization(v)
	}
	return Business{
		Typ:                "http://schema.org/Organization",
		LegalName:          legalName,
		CommencedTrading:   NewDate(commencedTrading),
		CeasedTrading:      NewDate(ceasedTrading),
		RenewalDueDate:     NewDate(renewalDueDate),
		RegistrationNumber: registrationNumber,
		ABN:                abn,
		Proprietors:        props,
	}

}

var metadataContext = Context{
	"abn":              "https://www.wikidata.org/wiki/Q4823913",
	"about":            "http://schema.org/about",
	"actor":            "http://schema.org/actor",
	"authority":        "http://www.records.nsw.gov.au/terms/disposalAuthority",
	"ceasedTrading":    "http://schema.org/dissolutionDate",
	"class":            "http://www.records.nsw.gov.au/terms/disposalClass",
	"commencedTrading": "http://schema.org/foundingDate",
	"consignment": Obj{
		ID:  "http://records.nsw.gov.au/terms/consignment",
		Typ: "@id",
	},
	"created": Obj{
		ID:  "http://purl.org/dc/terms/created",
		Typ: "http://www.w3.org/2001/XMLSchema#date",
	},
	"creators": "http://purl.org/dc/terms/creator",
	"director": "http://schema.org/director",
	"disposalRules": Obj{
		ID:  "http://www.records.nsw.gov.au/terms/disposalRules",
		Typ: "http://www.records.nsw.gov.au/terms/DisposalRule",
	},
	"duration":           "http://schema.org/duration",
	"language":           "http://schema.org/inLanguage",
	"legalName":          "http://schema.org/legalName",
	"name":               "http://schema.org/name",
	"productionCompany":  "http://schema.org/productionCompany",
	"proprietors":        "http://www.records.nsw.gov.au/terms/proprietors",
	"registrationNumber": "http://www.records.nsw.gov.au/terms/registrationNumber",
	"renewalDueDate": Obj{
		ID:  "http://www.records.nsw.gov.au/terms/renewalDueDate",
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

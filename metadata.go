package meta

// represents a metadata.json file
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

type DisposalRule struct {
	Authority string `json:"authority"`
	Class     string `json:"class"`
}

type Agent interface{}

func MakeAgent(name, id, typ string) Agent {
	if id == "" && typ == "" {
		return name
	}
	return ObjField{
		ID:   id,
		Typ:  typ,
		Name: name,
	}
}

func MakePerson(name string, id int) Agent {
	return MakeAgent(name, ToID(id, "http://www.records.nsw.gov.au/persons/"), "http://www.records.nsw.gov.au/terms/Person")
}

func MakeAgency(name string, id int) Agent {
	return MakeAgent(name, ToID(id, "http://records.nsw.gov.au/agencies/"), "http://www.records.nsw.gov.au/terms/Agency")
}

func MakeOrganization(name string) Agent {
	return MakeAgent(name, "", "http://schema.org/Organization")
}

type Thing interface{}

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

func NewMetadata(title, created string, series int) *Metadata {
	return &Metadata{
		Typ:     []string{"http://www.records.nsw.gov.au/terms/DigitalArchive"},
		Title:   title,
		Created: NewDate(created),
		Series:  ToID(series, "http://www.records.nsw.gov.au/series/"),
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
	"consignment": ObjField{
		ID:  "http://records.nsw.gov.au/terms/consignment",
		Typ: "@id",
	},
	"created": ObjField{
		ID:  "http://purl.org/dc/terms/created",
		Typ: "http://www.w3.org/2001/XMLSchema#date",
	},
	"creators": "http://purl.org/dc/terms/creator",
	"director": "http://schema.org/director",
	"disposalRules": ObjField{
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
	"renewalDueDate": ObjField{
		ID:  "http://www.records.nsw.gov.au/terms/renewalDueDate",
		Typ: "http://www.w3.org/2001/XMLSchema#date",
	},
	"series": ObjField{
		ID:  "http://records.nsw.gov.au/terms/series",
		Typ: "@id",
	},
	"source":    "http://purl.org/dc/terms/source",
	"subtitles": "http://schema.org/subtitleLanguage",
	"title":     "http://purl.org/dc/terms/title",
}

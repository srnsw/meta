// Copyright 2018 State of New South Wales through the State Archives and Records Authority of NSW
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package meta

// Metadata represents a metadata.json file
type Metadata struct {
	ID                string   `json:"@id"`
	Migration         string   `json:"migration"`
	Typ               VarStr   `json:"@type"`
	Title             string   `json:"title"`
	Description       string   `json:"description,omitempty"`
	Created           *W3CDate `json:"created,omitempty"`
	Modified          *W3CDate `json:"modified,omitempty"`
	Creator           Agent    `json:"creator,omitempty"`
	Source            VarStr   `json:"source,omitempty"`
	IsPartOf          VarStr   `json:"isPartOf,omitempty"`
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
		ID:        ReferenceObject(id),
		Migration: "mig:0",
		Typ:       "http://records.nsw.gov.au/terms/DigitalArchive",
		Title:     title,
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

var metadataContext = Context{
	"abn":              "http://www.wikidata.org/wiki/Q4823913",
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
	"creator":     "http://purl.org/dc/terms/creator",
	"description": "http://purl.org/dc/terms/description",
	"director":    "http://schema.org/director",
	"disposalRule": Obj{
		ID:  "http://records.nsw.gov.au/terms/disposalRule",
		Typ: "http://records.nsw.gov.au/terms/DisposalRule",
	},
	"duration":  "http://schema.org/duration",
	"isPartOf":  "http://purl.org/dc/terms/isPartOf",
	"language":  "http://schema.org/inLanguage",
	"legalName": "http://schema.org/legalName",
	"migration": Obj{
		ID:  "http://records.nsw.gov.au/terms/migration",
		Typ: "@id",
	},
	"modified": Obj{
		ID:  "http://purl.org/dc/terms/modified",
		Typ: "http://www.w3.org/2001/XMLSchema#date",
	},
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

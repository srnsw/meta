package meta

import "time"

// represents a manifest.json file
type Manifest struct {
	Context     Context      `json:"@context"`
	AccessRules []AccessRule `json:"accessRules,omitempty"`
	Versions    []Version    `json:"versions,omitempty"`
}

type AccessRule struct {
	ID           string    `json:"@id"`
	ExecuteDate  time.Time `json:"executeDate"` // time.Time marshals in RFC3339 format
	Scope        string    `json:"scope"`
	Publish      bool      `json:"publish"`
	Basis        *Basis    `json:"basis,omitempty"`
	Patch        *int      `json:"metadataPatch,omitempty"`
	FullManifest *bool     `json:"fullManifest,omitempty"`
	Display      []string  `json:"displayTarget,omitempty"`
	Preview      []string  `json:"previewTarget,omitempty"`
	Text         []string  `json:"textTarget,omitempty"`
}

type Basis struct {
	AccessDirection   int    `json:"accessDirection"`
	AccessDescription string `json:"accessDescription,omitempty"`
}

type Version struct {
	ID             string   `json:"@id"`
	Base           string   `json:"base,omitempty"`
	HasAccessRules []string `json:"hasAccessRules,omitempty"`
	Files          []File   `json:"files"`
}

type File struct {
	ID             string     `json:"@id"`
	Name           string     `json:"name"`
	OriginalName   string     `json:"originalName,omitempty"`
	Size           int64      `json:"size"`
	Created        *time.Time `json:"created,omitempty"`
	Modified       *time.Time `json:"modified,omitempty"`
	MIME           string     `json:"mime,omitempty"`
	PUID           string     `json:"puid,omitempty"`
	Hash           *Hash      `json:"hasHash,omitempty"`
	HasAccessRules []string   `json:"hasAccessRules,omitempty"`
}

type Hash struct {
	Algorithm string `json:"hashAlgorithm,omitempty"`
	Value     string `json:"hashValue,omitempty"`
}

var manifestContext = Context{
	"accessRules": "http://www.records.nsw.gov.au/repo/accessRules",
	"versions": ObjField{
		ID:        "http://www.records.nsw.gov.au/repo/versions",
		Container: "@list",
	},
	"executeDate": ObjField{
		ID:  "http://www.records.nsw.gov.au/repo/executeDate",
		Typ: "http://www.w3.org/2001/XMLSchema#dateTime",
	},
	"scope": "http://www.records.nsw.gov.au/repo/scope",
	"publish": ObjField{
		ID:  "http://www.records.nsw.gov.au/repo/publish",
		Typ: "http://www.w3.org/2001/XMLSchema#boolean",
	},
	"basis": "http://www.records.nsw.gov.au/repo/basis",
	"accessDirection": ObjField{
		ID:  "http://www.records.nsw.gov.au/repo/accessDirection",
		Typ: "http://www.w3.org/2001/XMLSchema#integer",
	},
	"fullManifest": ObjField{
		ID:  "http://www.records.nsw.gov.au/repo/fullManifest",
		Typ: "http://www.w3.org/2001/XMLSchema#boolean",
	},
	"displayTarget": ObjField{
		ID:        "http://www.records.nsw.gov.au/repo/displayTarget",
		Typ:       "@id",
		Container: "@list",
	},
	"previewTarget": ObjField{
		ID:        "http://www.records.nsw.gov.au/repo/previewTarget",
		Typ:       "@id",
		Container: "@list",
	},
	"textTarget": ObjField{
		ID:        "http://www.records.nsw.gov.au/repo/textTarget",
		Typ:       "@id",
		Container: "@list",
	},
	"hasAccessRules": ObjField{
		ID:        "http://www.records.nsw.gov.au/repo/hasAccessRules",
		Typ:       "@id",
		Container: "@list",
	},
	"base": "http://www.records.nsw.gov.au/repo/base",
	"files": ObjField{
		ID:        "http://www.openarchives.org/ore/0.9/jsonld/aggregates",
		Container: "@list"},
	"name":         "http://www.semanticdesktop.org/ontologies/2007/03/22/nfo#fileName",
	"originalName": "http://id.loc.gov/vocabulary/preservation/hasOriginalName",
	"size": ObjField{
		ID:  "http://www.semanticdesktop.org/ontologies/2007/03/22/nfo#fileSize",
		Typ: "http://www.w3.org/2001/XMLSchema#integer",
	},
	"created": ObjField{
		ID:  "http://www.semanticdesktop.org/ontologies/2007/03/22/nfo#fileCreated",
		Typ: "http://www.w3.org/2001/XMLSchema#dateTime",
	},
	"modified": ObjField{
		ID:  "http://www.semanticdesktop.org/ontologies/2007/03/22/nfo#fileLastModified",
		Typ: "http://www.w3.org/2001/XMLSchema#dateTime",
	},
	"mime": ObjField{
		ID:  "http://purl.org/dc/terms/format",
		Typ: "http://purl.org/dc/terms/FileFormat",
	},
	"puid": ObjField{
		ID:  "http://purl.org/dc/terms/format",
		Typ: "http://purl.org/dc/terms/FileFormat",
	},
	"hash": ObjField{
		ID:  "http://www.semanticdesktop.org/ontologies/2007/03/22/nfo#hasHash",
		Typ: "http://www.semanticdesktop.org/ontologies/2007/03/22/nfo#FileHash",
	},
	"hashAlgorithm": "http://www.semanticdesktop.org/ontologies/2007/03/22/nfo#hashAlgorithm",
	"hashValue":     "http://www.semanticdesktop.org/ontologies/2007/03/22/nfo#hashValue",
}

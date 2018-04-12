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
	Text         []string  `json:"textTarget,omitempty"`
	Display      []string  `json:"displayTarget,omitempty"`
	Preview      []string  `json:"previewTarget,omitempty"`
}

type Basis struct {
	AccessDirection   int    `json:"accessDirection"`
	AccessDescription string `json:"accessDescription,omitempty"`
}

type Version struct {
	Version        int      `json:"versionNumber"`
	Files          []File   `json:"ore:aggregates"`
	HasAccessRules []string `json:"hasAccessRules,omitempty"`
}

type File struct {
	Name        string   `json:"fileName"`
	OldPath     string   `json:"originalName,omitempty"`
	Size        int64    `json:"fileSize"`
	Format      *Format  `json:"format,omitempty"`
	Hash        *Hash    `json:"hasHash,omitempty"`
	Mime        string   `json:"format,omitempty"`
	AccessRules []string `json:"accessRules,omitempty"`
}

type Format struct {
	Registry string `json:"formatRegistry,omitempty"`
	Puid     string `json:"formatDesignation,omitempty"`
}

type Hash struct {
	Algorithm string `json:"hashAlgorithm,omitempty"`
	Value     string `json:"hashValue,omitempty"`
}

var manifestContext = Context{
	"accessRules": "http://www.records.nsw.gov.au/repo/accessRules",
	"executeDate": ObjField{
		ID:  "http://www.records.nsw.gov.au/repo/executeDate",
		Typ: "http://www.w3.org/2001/XMLSchema#dateTime",
	},
	"scope":             "http://www.records.nsw.gov.au/repo/scope",
	"publish":           "http://www.records.nsw.gov.au/repo/publish",
	"basis":             "http://www.records.nsw.gov.au/repo/basis",
	"accessDirection":   "http://www.records.nsw.gov.au/repo/accessDirection",
	"accessDescription": "http://www.records.nsw.gov.au/repo/accessDescription",
	"fullManifest":      "http://www.records.nsw.gov.au/repo/fullManifest",
	"displayTarget": ObjField{
		ID:  "http://www.records.nsw.gov.au/repo/displayTarget",
		Typ: "@id",
	},
	"previewTarget": ObjField{
		ID:  "http://www.records.nsw.gov.au/repo/previewTarget",
		Typ: "@id",
	},
	"textTarget": ObjField{
		ID:  "http://www.records.nsw.gov.au/repo/textTarget",
		Typ: "@id",
	},
	"hasAccessRules": ObjField{
		ID:  "http://www.records.nsw.gov.au/repo/hasAccessRules",
		Typ: "@id",
	},
	"versions": "http://www.records.nsw.gov.au/repo/versions",
	"files":    "http://www.openarchives.org/ore/0.9/jsonld/aggregates",
	"name":     "http://www.semanticdesktop.org/ontologies/2007/03/22/nfo#fileName",
	"size":     "http://www.semanticdesktop.org/ontologies/2007/03/22/nfo#fileSize",
	"created":  "http://www.semanticdesktop.org/ontologies/2007/03/22/nfo#fileCreated",
	"modified": "http://www.semanticdesktop.org/ontologies/2007/03/22/nfo#fileLastModified",
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

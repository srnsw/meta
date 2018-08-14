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

import (
	"errors"
	"strconv"
	"time"
)

// Manifest represents a manifest.json file
type Manifest struct {
	Type        string       `json:"@type"`
	AccessRules []AccessRule `json:"accessRules,omitempty"`
	Versions    []Version    `json:"versions,omitempty"`
	Context     Context      `json:"@context"`
}

// AccessRule represents accessRules
type AccessRule struct {
	ID           string  `json:"@id"`
	ExecuteDate  W3CDate `json:"executeDate"`
	Scope        string  `json:"scope"`
	Publish      bool    `json:"publish"`
	Basis        *Basis  `json:"basis,omitempty"`
	Patch        *int    `json:"metadataPatch,omitempty"`
	FullManifest *bool   `json:"fullManifest,omitempty"`
	Display      VarStr  `json:"displayTarget,omitempty"`
	Preview      VarStr  `json:"previewTarget,omitempty"`
	Text         VarStr  `json:"textTarget,omitempty"`
}

// Basis is a json basis
type Basis struct {
	AccessDirection   string `json:"accessDirection"`
	AccessDescription string `json:"accessDescription,omitempty"`
}

// Version represents versions
type Version struct {
	ID             string   `json:"@id"`
	Base           string   `json:"base,omitempty"`
	DerivedFrom    string   `json:"derivedFrom,omitempty"`
	GeneratedBy    string   `json:"generatedBy,omitempty"`
	HasAccessRules []string `json:"hasAccessRules,omitempty"`
	Files          []File   `json:"files"`
}

// File represents files
type File struct {
	ID             string     `json:"@id"`
	Name           string     `json:"name"`
	OriginalName   string     `json:"originalName,omitempty"`
	Size           int64      `json:"size"`
	Created        *time.Time `json:"fileCreated,omitempty"`
	Modified       *time.Time `json:"modified,omitempty"`
	MIME           string     `json:"mime,omitempty"`
	PUID           string     `json:"puid,omitempty"`
	Hash           *Hash      `json:"hash,omitempty"`
	HasAccessRules []string   `json:"hasAccessRules,omitempty"`
}

// Hash is a json hash
type Hash struct {
	Algorithm string `json:"hashAlgorithm,omitempty"`
	Value     string `json:"hashValue,omitempty"`
}

// NewManifest returns a reference to a manifest. It also sets the @type field.
func NewManifest() *Manifest {
	return &Manifest{
		Type: "http://records.nsw.gov.au/terms/Manifest",
	}
}

// AddAR adds a new access rule to a Manifest
// Provide access rule fields in the supplied arguments. Because of their rarity, Patch and FullManifest fields aren't supplied as
// arguments but should be manipulated directly.
// Returns the new access rule's @id and an error.
func (m *Manifest) AddAR(
	executeDate, scope string,
	publish bool,
	accessDirection int,
	accessDescription string,
	display, preview, text []FileTarget,
) (arid string, err error) {
	m.AccessRules, arid, err = AppendAR(m.AccessRules, executeDate, scope, publish, accessDirection, accessDescription, display, preview, text)
	return arid, err
}

// AppendAR adds a new access rule to a set of rules (or create a new set of rules by appending to nil)
// Provide access rule fields in the supplied arguments. Because of their rarity, Patch and FullManifest fields aren't supplied as
// arguments but should be manipulated directly.
// Returns the list of access rules, the new access rule's @id, and an error.
func AppendAR(rules []AccessRule,
	executeDate, scope string,
	publish bool,
	accessDirection int,
	accessDescription string,
	display, preview, text []FileTarget,
) ([]AccessRule, string, error) {
	t, err := ParseDate(executeDate)
	if err != nil {
		return nil, "", err
	}
	switch scope {
	case "root", "global", "local":
	case "":
		scope = "global"
	default:
		return nil, "", errors.New("scope must be either root, global or local")
	}
	var basis *Basis
	if accessDirection > 0 {
		basis = &Basis{
			AccessDirection:   ToID(accessDirection, "http://records.nsw.gov.au/accessDirection/"),
			AccessDescription: accessDescription,
		}
	}
	arid := ReferenceN(Ref{"ar", len(rules)})
	return append(rules, AccessRule{
		ID:          arid,
		ExecuteDate: t,
		Scope:       scope,
		Publish:     publish,
		Basis:       basis,
		Display:     ReferenceFiles(display),
		Preview:     ReferenceFiles(preview),
		Text:        ReferenceFiles(text),
	}), arid, nil
}

// AddVersion adds a new version to a Manifest.
// It takes a slice of access rules, and a slice of files.
// It returns the version @id.
func (m *Manifest) AddVersion(files []File) string {
	vid := ReferenceVersion(len(m.Versions))
	for i, _ := range files {
		files[i].ID = ReferenceN(Ref{"v", len(m.Versions)}, Ref{"f", i})
	}
	base := "versions/" + strconv.Itoa(len(m.Versions))
	m.Versions = append(m.Versions, Version{
		ID:    vid,
		Base:  base,
		Files: files,
	})
	return vid
}

// CopyARs copies a slice of AccessRules. This allows you to create a template access rule and copy it e.g. to change file references.
func CopyARs(ar []AccessRule) []AccessRule {
	ret := make([]AccessRule, len(ar))
	for i, v := range ar {
		ret[i] = copyAR(v)
	}
	return ret
}

// copy by reference is fine for most fields of an AccessRule,
// except we don't want display/preview/text to be copied by reference
func copyAR(ar AccessRule) AccessRule {
	return AccessRule{
		ID:           ar.ID,
		ExecuteDate:  ar.ExecuteDate,
		Scope:        ar.Scope,
		Publish:      ar.Publish,
		Basis:        ar.Basis,
		Patch:        ar.Patch,
		FullManifest: ar.FullManifest,
		Display:      copyVarStr(ar.Display),
		Preview:      copyVarStr(ar.Preview),
		Text:         copyVarStr(ar.Text),
	}
}

// Ref supports the construction of internal reference strings (blank node IDs) e.g. _:ar1
type Ref struct {
	Prefix string
	N      int
}

// ReferenceN constructs an internal reference string from the supplied Ref(s)
func ReferenceN(refs ...Ref) string {
	base := "_:"
	for _, ref := range refs {
		base += ref.Prefix + strconv.Itoa(ref.N)
	}
	return base
}

// FileTarget is an internal reference for a file target
type FileTarget [2]int

// String is a string representation of a FileTarget
func (ft FileTarget) String() string {
	return ReferenceN(
		Ref{"v", ft[0]},
		Ref{"f", ft[1]},
	)
}

func ReferenceVersion(v int) string {
	return ReferenceN(Ref{"v", v})
}

func ReferenceFiles(ts []FileTarget) VarStr {
	if len(ts) == 0 {
		return nil
	}
	if len(ts) == 1 {
		return ts[0].String()
	}
	ret := make([]string, len(ts))
	for i, v := range ts {
		ret[i] = v.String()
	}
	return ret
}

var manifestContext = Context{
	"accessDescription": "http://records.nsw.gov.au/terms/accessDescription",
	"accessDirection": Obj{
		ID:  "http://records.nsw.gov.au/terms/accessDirection",
		Typ: "@id",
	},
	"accessRules": Obj{
		ID:  "http://records.nsw.gov.au/terms/accessRules",
		Typ: "http://records.nsw.gov.au/terms/AccessRule",
	},
	"base": "http://records.nsw.gov.au/terms/base",
	"basis": Obj{
		ID:  "http://records.nsw.gov.au/terms/basis",
		Typ: "http://records.nsw.gov.au/terms/Basis",
	},
	"derivedFrom": Obj{
		ID:  "http://www.w3.org/ns/prov#wasDerivedFrom",
		Typ: "@id",
	},
	"displayTarget": Obj{
		ID:  "http://records.nsw.gov.au/terms/displayTarget",
		Typ: "@id",
	},
	"executeDate": Obj{
		ID:  "http://records.nsw.gov.au/terms/executeDate",
		Typ: "http://www.w3.org/2001/XMLSchema#date",
	},
	"fileCreated": Obj{
		ID:  "http://www.semanticdesktop.org/ontologies/2007/03/22/nfo#fileCreated",
		Typ: "http://www.w3.org/2001/XMLSchema#dateTime",
	},
	"files": Obj{
		ID:  "http://www.openarchives.org/ore/0.9/jsonld#aggregates",
		Typ: "http://www.semanticdesktop.org/ontologies/2007/03/22/nfo#FileDataObject",
	},
	"fullManifest": Obj{
		ID:  "http://records.nsw.gov.au/terms/fullManifest",
		Typ: "http://www.w3.org/2001/XMLSchema#boolean",
	},
	"generatedBy": Obj{
		ID:  "http://www.w3.org/ns/prov#wasGeneratedBy",
		Typ: "@id",
	},
	"hasAccessRules": Obj{
		ID:  "http://records.nsw.gov.au/terms/hasAccessRules",
		Typ: "@id",
	},
	"hash":          "http://www.semanticdesktop.org/ontologies/2007/03/22/nfo#hasHash",
	"hashAlgorithm": "http://www.semanticdesktop.org/ontologies/2007/03/22/nfo#hashAlgorithm",
	"hashValue":     "http://www.semanticdesktop.org/ontologies/2007/03/22/nfo#hashValue",
	"mime": Obj{
		ID:  "http://purl.org/dc/terms/format",
		Typ: "http://purl.org/dc/terms/MediaType",
	},
	"modified": Obj{
		ID:  "http://www.semanticdesktop.org/ontologies/2007/03/22/nfo#fileLastModified",
		Typ: "http://www.w3.org/2001/XMLSchema#dateTime",
	},
	"name":         "http://www.semanticdesktop.org/ontologies/2007/03/22/nfo#fileName",
	"originalName": "http://id.loc.gov/vocabulary/preservation/hasOriginalName",
	"previewTarget": Obj{
		ID:  "http://records.nsw.gov.au/terms/previewTarget",
		Typ: "@id",
	},
	"publish": Obj{
		ID:  "http://records.nsw.gov.au/terms/publish",
		Typ: "http://www.w3.org/2001/XMLSchema#boolean",
	},
	"puid": Obj{
		ID:  "http://purl.org/dc/terms/format",
		Typ: "@id",
	},
	"scope": "http://records.nsw.gov.au/terms/scope",
	"size": Obj{
		ID:  "http://www.semanticdesktop.org/ontologies/2007/03/22/nfo#fileSize",
		Typ: "http://www.w3.org/2001/XMLSchema#integer",
	},
	"textTarget": Obj{
		ID:  "http://records.nsw.gov.au/terms/textTarget",
		Typ: "@id",
	},
	"versions": Obj{
		ID:  "http://records.nsw.gov.au/terms/versions",
		Typ: "http://records.nsw.gov.au/terms/Version",
	},
}

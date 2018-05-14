package meta

import (
	"errors"
	"strconv"
	"time"
)

// represents a manifest.json file
type Manifest struct {
	Type        string       `json:"@type"`
	AccessRules []AccessRule `json:"accessRules,omitempty"`
	Versions    []Version    `json:"versions,omitempty"`
	Context     Context      `json:"@context"`
}

type AccessRule struct {
	ID           string   `json:"@id"`
	ExecuteDate  W3CDate  `json:"executeDate"`
	Scope        string   `json:"scope"`
	Publish      bool     `json:"publish"`
	Basis        *Basis   `json:"basis,omitempty"`
	Patch        *int     `json:"metadataPatch,omitempty"`
	FullManifest *bool    `json:"fullManifest,omitempty"`
	Display      []string `json:"displayTarget,omitempty"`
	Preview      []string `json:"previewTarget,omitempty"`
	Text         []string `json:"textTarget,omitempty"`
}

type Basis struct {
	AccessDirection   string `json:"accessDirection"`
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
	Created        *time.Time `json:"fileCreated,omitempty"`
	Modified       *time.Time `json:"modified,omitempty"`
	MIME           string     `json:"mime,omitempty"`
	PUID           string     `json:"puid,omitempty"`
	Hash           *Hash      `json:"hash,omitempty"`
	HasAccessRules []string   `json:"hasAccessRules,omitempty"`
}

type Hash struct {
	Algorithm string `json:"hashAlgorithm,omitempty"`
	Value     string `json:"hashValue,omitempty"`
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

func (ft FileTarget) String() string {
	return ReferenceN(
		Ref{"v", ft[0]},
		Ref{"f", ft[1]},
	)
}

func referenceFiles(ts []FileTarget) []string {
	ret := make([]string, len(ts))
	for i, v := range ts {
		ret[i] = v.String()
	}
	return ret
}

// copy by reference is fine for most fields of an AccessRule,
// except we don't want display/preview/text to be copied by reference
func copyAR(ar AccessRule) AccessRule {
	dis := make([]string, len(ar.Display))
	copy(dis, ar.Display)
	prev := make([]string, len(ar.Preview))
	copy(dis, ar.Preview)
	txt := make([]string, len(ar.Text))
	copy(dis, ar.Text)
	return AccessRule{
		ID:           ar.ID,
		ExecuteDate:  ar.ExecuteDate,
		Scope:        ar.Scope,
		Publish:      ar.Publish,
		Basis:        ar.Basis,
		Patch:        ar.Patch,
		FullManifest: ar.FullManifest,
		Display:      dis,
		Preview:      prev,
		Text:         txt,
	}
}

func CopyARs(ar []AccessRule) []AccessRule {
	ret := make([]AccessRule, len(ar))
	for i, v := range ar {
		ret[i] = copyAR(v)
	}
	return ret
}

func NewManifest() *Manifest {
	return &Manifest{
		Type: "http://www.records.nsw.gov.au/terms/Manifest",
	}
}

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

// AppendAccessRule is a helper function to add a new access rule to a set of rules (or create a new set of rules by appending to nil)
// Provide access rule fields in the supplied arguments. Because of their rarity, Patch and FullManifest fields aren't supplied as
// arguments but should be manipulated directly. Returns the list of access rules, the new access rule's ID and an error.
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
			AccessDirection:   ToID(accessDirection, "http://www.records.nsw.gov.au/accessDirection/"),
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
		Display:     referenceFiles(display),
		Preview:     referenceFiles(preview),
		Text:        referenceFiles(text),
	}), arid, nil
}

// add version applies IDs to both the version and the files
func (m *Manifest) AddVersion(base string, rules []string, files []File) string {
	vid := ReferenceN(Ref{"v", len(m.Versions)})
	for i, _ := range files {
		files[i].ID = ReferenceN(Ref{"v", len(m.Versions)}, Ref{"f", i})
	}
	m.Versions = append(m.Versions, Version{
		ID:             vid,
		Base:           base,
		HasAccessRules: rules,
		Files:          files,
	})

	return vid
}

var manifestContext = Context{
	"accessDescription": "http://www.records.nsw.gov.au/terms/accessDescription",
	"accessDirection": ObjField{
		ID:  "http://www.records.nsw.gov.au/terms/accessDirection",
		Typ: "@id",
	},
	"accessRules": ObjField{
		ID:  "http://www.records.nsw.gov.au/terms/accessRules",
		Typ: "http://www.records.nsw.gov.au/terms/AccessRule",
	},
	"base": "http://www.records.nsw.gov.au/terms/base",
	"basis": ObjField{
		ID:  "http://www.records.nsw.gov.au/terms/basis",
		Typ: "http://www.records.nsw.gov.au/terms/Basis",
	},
	"displayTarget": ObjField{
		ID:  "http://www.records.nsw.gov.au/terms/displayTarget",
		Typ: "@id",
	},
	"executeDate": ObjField{
		ID:  "http://www.records.nsw.gov.au/terms/executeDate",
		Typ: "http://www.w3.org/2001/XMLSchema#date",
	},
	"fileCreated": ObjField{
		ID:  "http://www.semanticdesktop.org/ontologies/2007/03/22/nfo#fileCreated",
		Typ: "http://www.w3.org/2001/XMLSchema#dateTime",
	},
	"files": ObjField{
		ID:  "http://www.openarchives.org/ore/0.9/jsonld#aggregates",
		Typ: "http://www.semanticdesktop.org/ontologies/2007/03/22/nfo#FileDataObject",
	},
	"fullManifest": ObjField{
		ID:  "http://www.records.nsw.gov.au/terms/fullManifest",
		Typ: "http://www.w3.org/2001/XMLSchema#boolean",
	},
	"hasAccessRules": ObjField{
		ID:  "http://www.records.nsw.gov.au/terms/hasAccessRules",
		Typ: "@id",
	},
	"hash":          "http://www.semanticdesktop.org/ontologies/2007/03/22/nfo#hasHash",
	"hashAlgorithm": "http://www.semanticdesktop.org/ontologies/2007/03/22/nfo#hashAlgorithm",
	"hashValue":     "http://www.semanticdesktop.org/ontologies/2007/03/22/nfo#hashValue",
	"mime": ObjField{
		ID:  "http://purl.org/dc/terms/format",
		Typ: "http://purl.org/dc/terms/FileFormat",
	},
	"modified": ObjField{
		ID:  "http://www.semanticdesktop.org/ontologies/2007/03/22/nfo#fileLastModified",
		Typ: "http://www.w3.org/2001/XMLSchema#dateTime",
	},
	"name":         "http://www.semanticdesktop.org/ontologies/2007/03/22/nfo#fileName",
	"originalName": "http://id.loc.gov/vocabulary/preservation/hasOriginalName",
	"previewTarget": ObjField{
		ID:  "http://www.records.nsw.gov.au/terms/previewTarget",
		Typ: "@id",
	},
	"publish": ObjField{
		ID:  "http://www.records.nsw.gov.au/terms/publish",
		Typ: "http://www.w3.org/2001/XMLSchema#boolean",
	},
	"puid": ObjField{
		ID:  "http://purl.org/dc/terms/format",
		Typ: "@id",
	},
	"scope": "http://www.records.nsw.gov.au/terms/scope",
	"size": ObjField{
		ID:  "http://www.semanticdesktop.org/ontologies/2007/03/22/nfo#fileSize",
		Typ: "http://www.w3.org/2001/XMLSchema#integer",
	},
	"textTarget": ObjField{
		ID:  "http://www.records.nsw.gov.au/terms/textTarget",
		Typ: "@id",
	},
	"versions": ObjField{
		ID:  "http://www.records.nsw.gov.au/terms/versions",
		Typ: "http://www.records.nsw.gov.au/terms/Version",
	},
}

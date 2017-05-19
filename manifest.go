package meta

// represents a manifest.json file
type Manifest struct {
	AccessRules []AccessRule `json:"repo:accessRules,omitempty"`
	Versions    []Version    `json:"repo:versions,omitempty"`
}

type AccessRule struct {
	ID            string `json:"@id"`
	ExecuteDate   string `json:"repo:executeDate"`
	Scope         string `json:"repo:scope,omitempty"`
	Publish       bool   `json:"repo:publish,omitempty"`
	Basis         Basis  `json:"repo:basis,omitempty"`
	Patch         int    `json:"repo:metadataPatch,omitempty"`
	FullManifest  bool   `json:"repo:fullManifest,omitempty"`
	Text          string `json:"repo:textTarget,omitempty"`
	DisplayTarget string `json:"repo:displayTarget,omitempty"`
	DisplayType   string `json:"repo:displayType,omitempty"`
}

type Basis struct {
	AccessDirection   int    `json:"repo:accessDirection,omitempty"`
	AccessDescription string `json:"repo:accessDescription,omitempty"`
}

type Version struct {
	Version        int      `json:"repo:versionNumber"`
	Files          []File   `json:"ore:aggregates"`
	HasAccessRules []string `json:"repo:hasAccessRules,omitempty"`
}

type File struct {
	Name        string   `json:"nfo:fileName"`
	OldPath     string   `json:"premis:originalName,omitempty"`
	Size        int64    `json:"nfo:fileSize"`
	Format      *Format  `json:"premis:format,omitempty"`
	Hash        *Hash    `json:"nfo:hasHash,omitempty"`
	Mime        string   `json:"dc:format,omitempty"`
	AccessRules []string `json:"repo:accessRules,omitempty"`
}

type Format struct {
	Registry string `json:"premis:formatRegistry,omitempty"`
	Puid     string `json:"premis:formatDesignation,omitempty"`
}

type Hash struct {
	Algorithm string `json:"nfo:hashAlgorithm,omitempty"`
	Value     string `json:"nfo:hashValue,omitempty"`
}

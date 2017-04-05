package meta

// represents a manifest.json file
type Manifest struct {
	AccessRule AccessRule `json:"accessRule"`
	Preview    []string   `json:"preview"`
	Versions   []Version  `json:"versions,omitempty"`
}

type Version struct {
	Version int    `json:"version"`
	Files   []File `json:"files"`
}

type File struct {
	Name    string   `json:"name"`
	Display bool     `json:"display,omitempty"`
	Path    []string `json:"path,omitempty"`
	Size    int64    `json:"fileSize"`
	Puid    string   `json:"format,omitempty"`
	Mime    string   `json:"mimeType,omitempty"`
}

type AccessRule struct {
	Direction     int    `json:"accessDirection"`
	Effect        string `json:"effect"`
	CalculateFrom string `json:"calculateFrom,omitempty"`
}

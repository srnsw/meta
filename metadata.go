package meta

// represents a metadata.json file
type Metadata struct {
	Title        string       `json:"title"`
	Created      string       `json:"created,omitempty"`
	Creators     []Agent      `json:"creators,omitempty"`
	Series       int          `json:"series"`
	DisposalRule DisposalRule `json:"disposalRule"`
}

type DisposalRule struct {
	Authority string `json:"authority"`
	Class     string `json:"class"`
}

type Agent interface{}

type Organisation struct {
	Organisation OrganisationFields `json:"organisation,omitempty"`
}

type OrganisationFields struct {
	Name     string `json:"name,omitempty"`
	AgencyID int    `json:"agencyID,omitempty"`
}

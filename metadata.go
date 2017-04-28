package meta

// represents a metadata.json file
type Metadata struct {
	Title         string         `json:"dct:title"`
	Created       string         `json:"created,omitempty"`
	Creators      []Agent        `json:"creators,omitempty"`
	Series        int            `json:"series"`
	DisposalRules []DisposalRule `json:"disposalRules"`
	Audio         *MediaFields   `json:"audio,omitempty"`
	Video         *MediaFields   `json:"video,omitempty"`
}

type MediaFields struct {
	Language          string `json:"lang,omitempty"`
	Duration          string `json:"duration,omitempty"`
	Subtitles         string `json:"subtitles,omitempty"`
	Director          string `json:"director,omitempty"`
	Actor             string `json:"actor,omitempty"`
	ProductionCompany string `json:"productionCompany,omitempty"`
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

package meta

import "time"

// Log represents a preservation event e.g. format migration.
// The PROV and PREMIS ontologies are primarily used for this metadata.
type Log struct {
	ID      string     `json:"@id"`
	Typ     string     `json:"@type"` // from http://id.loc.gov/vocabulary/preservation/eventType.html e.g. http://id.loc.gov/vocabulary/preservation/eventType/mig
	Start   *time.Time `json:"startTime,omitempty"`
	End     *time.Time `json:"endTime"`
	Detail  string     `json:"detail"`
	Agent   Agent      `json:"agent"`
	Context Context    `json:"@context"`
}

const (
	ModificationEvent = "http://id.loc.gov/vocabulary/preservation/eventType/mod"
	MigrationEvent    = "http://id.loc.gov/vocabulary/preservation/eventType/mig"
)

// NewLog creates a *Log
func NewLog(id int, typ string) *Log {
	return &Log{
		ID:  ReferenceLog(id),
		Typ: typ,
	}
}

// ReferenceLog makes a temporary reference to a log event.
// This reference is swapped for a UUID by the migrate tool.
func ReferenceLog(i int) string {
	return ToRef(i, "log")
}

var logContext = Context{
	"agent": Obj{
		ID:  "https://www.w3.org/ns/prov#wasAssociatedWith",
		Typ: "https://www.w3.org/ns/prov#Agent",
	},
	"detail": "http://id.loc.gov/vocabulary/preservation/hasNote",
	"endTime": Obj{
		ID:  "https://www.w3.org/ns/prov#endedAtTime",
		Typ: "http://www.w3.org/2001/XMLSchema#dateTime",
	},
	"name": "http://schema.org/name",
	"startTime": Obj{
		ID:  "https://www.w3.org/ns/prov#startedAtTime",
		Typ: "http://www.w3.org/2001/XMLSchema#dateTime",
	},
}

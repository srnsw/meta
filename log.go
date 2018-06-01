package meta

import (
	"fmt"
	"time"
)

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
	MigrationEvent = "http://id.loc.gov/vocabulary/preservation/eventType/mig"
)

func NewLog(id int, typ string) *Log {
	return &Log{
		ID:  ReferenceLog(id),
		Typ: typ,
	}
}

func ReferenceLog(i int) string {
	return fmt.Sprintf("$log{%d}", i)
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

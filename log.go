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
	"name":            "http://schema.org/name",
	"softwareVersion": "https://schema.org/softwareVersion",
	"startTime": Obj{
		ID:  "https://www.w3.org/ns/prov#startedAtTime",
		Typ: "http://www.w3.org/2001/XMLSchema#dateTime",
	},
}

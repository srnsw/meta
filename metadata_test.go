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
	"path/filepath"
	"testing"
	"time"
)

func TestMetadata(t *testing.T) {
	m := NewMetadata(0, "Business Name Registration - Duntryleague Country Club")
	m.AddType("http://schema.org/Movie")
	m.Description = "A very nice record"
	m.Creator = []Agent{MakeAgency("Office of Fair Trading", 0), MakePerson("Michael Bruce Baird", 288)}
	m.Created = NewDate("1902-01-01")
	m.Modified = NewDate("1903-02-02")
	m.AgencyID = "TRAN.001.890"
	m.Provenance = "https://twitter.com/MelGibson"
	m.Source = "https://twitter.com/"
	m.IsPartOf = "Exhibits 50. The trials of Mel Gibson"
	m.DeliveryMethod = "AUSTRALIA POST"
	m.DocumentType = "Report"
	m.Series = ToSeries(21404)
	m.Consignment = ToConsignment(189087)
	m.DisposalRule = DisposalRule{Authority: "DA48", Class: "1.1.1.2"}
	m.Duration = "13:47:30"
	m.Language = "en"
	m.Subtitles = "zh"
	m.Director = "Mel Gibson"
	m.Actor = "Mel Gibson"
	m.ProductionCompany = "Icon Films"
	commencedTrading, _ := time.Parse(w3cymd, "1990-01-01")
	ceasedTrading, _ := time.Parse(w3cymd, "1999-03-29")
	renewalDueDate, _ := time.Parse(w3cymd, "1998-11-30")
	m.About = MakeBusiness("Duntryleague Country Club", "A0369711", "Unknown", commencedTrading, ceasedTrading, renewalDueDate, "The Orange Golf Club Ltd")
	ctx, err := populate(metadataContext, m)
	if err != nil {
		t.Fatal(err)
	}
	m.Context = ctx
	if err = compare(m, "metadata.json"); err != nil {
		t.Fatal(err)
	}
}

func TestMetadataConst(t *testing.T) {
	m := NewMetadata(0, "State Records Act 1998 No 17")
	m.AddType("http://schema.org/Legislation")
	m.Description = "An Act to make provision for the creation, management and protection of the records of public offices of the State and to provide for public access to those records, to establish the State Archives and Records Authority; and for other purposes."
	m.Created = NewDate("1998")
	m.Creator = []Agent{MakeAgency("Department of Premier and Cabinet", 10)}
	m.Series = ToSeries(21404)
	m.Consignment = ToConsignment(189087)
	m.DisposalRule = DisposalRule{Authority: "DA48", Class: "1.1.1.2"}
	m.Language = "en"
	ctx, err := populate(metadataContext, m)
	if err != nil {
		t.Fatal(err)
	}
	m.Context = ctx
	if err = compare(m, filepath.Join("project-0", "0", "metadata.json")); err != nil {
		t.Fatal(err)
	}
}

func TestAppendAgent(t *testing.T) {
	agents := AppendAgent(nil, MakeSDOPerson("Richard Lehane"))
	agents = AppendAgent(agents, MakeOrganization("The ANZ Bank"))
	agents = AppendAgent(agents, MakeSDOPerson("Prince Richard"))
	if slc, ok := agents.([]Agent); !ok {
		t.Fatalf("Expecting a slice of agents, got %t", agents)
	} else {
		if len(slc) != 3 {
			t.Fatalf("Expecting 3 agents, got %d", len(slc))
		}
	}
}

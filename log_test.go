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
)

func TestLog(t *testing.T) {
	l := NewLog(0, MigrationEvent)
	l.Start, l.End = NewDateTime("2015-04-20T17:41:48+10:00"), NewDateTime("2015-04-20T17:42:00+10:00")
	l.Detail = "Manually created CSV using MS Excel 2013"
	l.Agent = MakeSDOPerson("Richard Lehane")
	ctx, err := populate(logContext, l)
	if err != nil {
		t.Fatal(err)
	}
	l.Context = ctx
	if err = compare(l, "log.json"); err != nil {
		t.Fatal(err)
	}
}

func TestLogConst(t *testing.T) {
	l := NewLog(0, MigrationEvent)
	l.Start, l.End = NewDateTime("2015-04-20T17:41:48+10:00"), NewDateTime("2015-04-20T17:42:00+10:00")
	l.Detail = "Manually created CSV using MS Excel 2013"
	l.Agent = MakeSDOPerson("Richard Lehane")
	ctx, err := populate(logContext, l)
	if err != nil {
		t.Fatal(err)
	}
	l.Context = ctx
	if err = compare(l, filepath.Join("project-0", "0", "logs", "0.json")); err != nil {
		t.Fatal(err)
	}
}

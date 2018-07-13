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

import "testing"

const testJson = `
	[{
		"apple": "orange",
		"banana": {
			"carrot": 16,
			"tee pee": false
		},
		"watermelons": ["one", "two", {"tricky": 1}, "four"],
		"mango": 5
	},
	"velvet",
	{"calypso": "crazy"}]

`

func TestKeys(t *testing.T) {
	ks, _ := keys([]byte(testJson))
	if len(ks) != 8 {
		t.Errorf("Expecting 8 keys, got %d with %v", len(ks), ks)
	}

}

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

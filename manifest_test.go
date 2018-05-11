package meta

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"testing"
	"time"
)

func compareJSON(b1, b2 []byte) (bool, error) {
	d1, d2 := json.NewDecoder(bytes.NewBuffer(b1)), json.NewDecoder(bytes.NewBuffer(b2))
	var i int
	for {
		t1, err1 := d1.Token()
		t2, err2 := d2.Token()
		if t1 != t2 {
			return false, fmt.Errorf("json tokens aren't equal: %v and %v; successfully matched %d tokens", t1, t2, i)
		}
		if err1 != err2 {
			return false, fmt.Errorf("json errors aren't equal: %v and %v; successfully matched %d tokens", err1, err2, i)
		}
		if err1 == io.EOF {
			return true, nil
		}
		i++
	}
}

func compare(v interface{}, fn string) error {
	byts, err := ioutil.ReadFile(filepath.Join("examples", fn))
	if err != nil {
		return err
	}
	byts2, err := json.Marshal(v)
	if err != nil {
		return err
	}
	if ok, err := compareJSON(byts, byts2); !ok {
		return fmt.Errorf("Marshalled JSON does not match expected JSON:\nError: %v\n------\nGot:\n%s\n------\nExpect:\n%s", err, string(byts2), string(byts))
	}
	return nil
}

func newTime(s string) *time.Time {
	t, _ := ParseDateTime(s)
	return &t
}

func TestManifest(t *testing.T) {
	m := NewManifest()
	arid, err := m.AddAR("2016-05-23", "global", true, "http://www.records.nsw.gov.au/accessDirection/1296", "Early", []FileTarget{{0, 0}, {0, 2}, {0, 3}}, []FileTarget{{1, 0}}, []FileTarget{{1, 1}})
	if err != nil {
		t.Fatal(err)
	}
	m.AccessRules[0].FullManifest = new(bool) // fullManifest needs to be manually set
	m.AddVersion("versions/0", []string{arid}, []File{
		{
			Name:         "index.html",
			OriginalName: "Teddies photos\\Materials_teddies_44250.jpg",
			Size:         4026,
			Created:      newTime("2015-04-20T17:41:48+10:00"),
			Modified:     newTime("2015-04-20T17:41:48+10:00"),
			MIME:         "text/html",
			PUID:         "http://www.nationalarchives.gov.uk/pronom/fmt/24",
			Hash: &Hash{
				Algorithm: "md5",
				Value:     "hfuehwoiuhfoeihjwpoih0197626",
			},
		}})
	// create manifest.json
	ctx, err := populate(manifestContext, m)
	if err != nil {
		t.Fatal(err)
	}
	m.Context = ctx
	if err = compare(m, "manifest.json"); err != nil {
		t.Fatal(err)
	}
}

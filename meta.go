package meta

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

// Cap defines the capacity of the index slice. Edit for large jobs to an approximate number of objects
var Cap int = 1000

// Loader is anything that can load data into a Meta
type Loader interface {
	Load(m *Meta) error
}

// Meta is a set of metadata (metadata.json and manifest.json)
// The Index field provides ordering.
// The Store field can be used to store arbitrary data needed for particular projects.
type Meta struct {
	Index    []string
	Metadata map[string]*Metadata
	Manifest map[string]*Manifest
	Store    map[string]interface{}
}

// Action is a function (such as file copy) that is called when generating output
type Action func(target, index string, metadata *Metadata, manifest *Manifest, store interface{}) error

// New creates a new meta from the supplied loaders
func New(loaders ...Loader) (*Meta, error) {
	m := &Meta{
		make([]string, 0, Cap),
		make(map[string]*Metadata),
		make(map[string]*Manifest),
		make(map[string]interface{})}
	for _, l := range loaders {
		if err := l.Load(m); err != nil {
			return nil, err
		}
	}
	return m, nil
}

// Output generates metadata.json and manifest.json files for all of a Meta's metadata.
// Arbitrary actions based on that data can also be called by this function.
// For testing provide a sample size e.g. 10. When running in production, give a negative int for sample.
// Target is the target output directory.
func (m *Meta) Output(sample int, target string, actions ...Action) error {
	for i, v := range m.Index {
		if sample == 0 {
			return nil
		}
		sample--
		meta, man, store := m.Metadata[v], m.Manifest[v], m.Store[v]
		// make the output directory, which is an incrementing integer
		out := filepath.Join(target, strconv.Itoa(i))
		err := os.MkdirAll(out, os.ModePerm)
		if err != nil {
			return err
		}
		// create metadata.json
		ctx, err := populate(metadataContext, meta)
		if err != nil {
			return err
		}
		meta.Context = ctx
		meta.Typ = metadataTyp
		j, err := json.MarshalIndent(meta, "", "  ")
		if err != nil {
			return err
		}
		if err = ioutil.WriteFile(filepath.Join(out, "metadata.json"), j, os.ModePerm); err != nil {
			return err
		}
		// create manifest.json
		ctx, err = populate(manifestContext, man)
		if err != nil {
			return err
		}
		man.Context = ctx
		j, err = json.MarshalIndent(man, "", "  ")
		if err != nil {
			return err
		}
		if err = ioutil.WriteFile(filepath.Join(out, "manifest.json"), j, os.ModePerm); err != nil {
			return err
		}
		// finally, execute the actions
		for _, a := range actions {
			if err := a(out, v, meta, man, store); err != nil {
				return err
			}
		}
	}
	return nil
}

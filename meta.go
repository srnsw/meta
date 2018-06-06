package meta

import (
	"encoding/csv"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

// Meta is a set of metadata (metadata.json and manifest.json)
// The Index field provides ordering.
// The Store field can be used to store arbitrary data needed for particular projects.
type Meta struct {
	Index    []string
	Metadata map[string]*Metadata
	Manifest map[string]*Manifest
	Logs     map[string][]*Log
	Store    map[string]interface{}
}

// Cap defines the capacity of the index slice. Edit for large jobs to an approximate number of objects
var Cap int = 1000

// Loader is anything that can load data into a Meta
type Loader interface {
	Load(m *Meta) error
}

// Action is a function (such as file copy) that is called when generating output
type Action func(meta *Meta, target, index string) error

// New creates a new meta from the supplied loaders
func New(loaders ...Loader) (*Meta, error) {
	m := &Meta{
		make([]string, 0, Cap),
		make(map[string]*Metadata),
		make(map[string]*Manifest),
		make(map[string][]*Log),
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
// Target is the target output directory.
func (m *Meta) Output(target string, actions ...Action) error {
	return m.Sample(0, -1, target, actions...)
}

// Sample generates metadata.json and manifest.json files for a sample of a Meta's metadata.
// Arbitrary actions based on that data can also be called by this function.
// For testing provide a sample size e.g. 10 and an index where you'd like the sample to start.
// If a negative index is provided then the index will be calculated from the end. I.e. -10 will return the final 10.
// Target is the target output directory.
func (m *Meta) Sample(index, sample int, target string, actions ...Action) error {
	if index < 0 && index > 0-len(m.Index) {
		index = len(m.Index) + index
	}
	for i, v := range m.Index {
		if i < index {
			continue
		}
		if sample == 0 {
			return nil
		}
		sample--
		// make the output directory, which is an incrementing integer
		out := filepath.Join(target, strconv.Itoa(i))
		err := os.MkdirAll(out, os.ModePerm)
		if err != nil {
			return err
		}
		// execute the actions
		for _, a := range actions {
			if err := a(m, out, v); err != nil {
				return err
			}
		}
		meta, man := m.Metadata[v], m.Manifest[v]
		// create metadata.json
		ctx, err := populate(metadataContext, meta)
		if err != nil {
			return err
		}
		meta.Context = ctx
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
		// create logs
		logs, ok := m.Logs[v]
		if !ok {
			continue
		}
		logdir := filepath.Join(out, "logs")
		err = os.MkdirAll(logdir, os.ModePerm)
		if err != nil {
			return err
		}
		for ii, log := range logs {
			ctx, err = populate(logContext, log)
			if err != nil {
				return err
			}
			log.Context = ctx
			j, err = json.MarshalIndent(log, "", "  ")
			if err != nil {
				return err
			}
			if err = ioutil.WriteFile(filepath.Join(logdir, strconv.Itoa(ii)+".json"), j, os.ModePerm); err != nil {
				return err
			}
		}
	}
	return nil
}

// Helpers

// ToID is a helper function that turns an integer identififier into a string @id
// It is useful for places where we use integer IDs e.g. series numbers but want an IRI for the @id
func ToID(i int, pat string) string {
	if i <= 0 {
		return ""
	}
	return pat + strconv.Itoa(i)
}

// ToRef is a helper function that turns an integer identififier into a string ref:I reference.
// It is used for placedholder IDs like log:1 that get swapped out by the Migrate tool.
func ToRef(i int, ref string) string {
	if i < 0 {
		return ""
	}
	return ref + ":" + strconv.Itoa(i)
}

// ToSz is a helper function that turns a string int into an int64
func ToSz(s string) int64 {
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}

// ReadAll is a helper function that opens a file at path and reads as a CSV.
// Returns a slice of string slices and an error.
func ReadAll(path string) ([][]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	reader := csv.NewReader(file)
	return reader.ReadAll()
}

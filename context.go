package meta

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"sort"
)

// Context generates the @context field in json output
type Context map[string]Field

// Fields are typically plain strings or objects with @id/@type
type Field interface{}

// Obj is a json object. Used with @id/@type in @context.
// Can also be used to generate generic objects e.g. Agents and containers in metadata are Objs
type Obj struct {
	ID        string `json:"@id,omitempty"`
	Typ       string `json:"@type,omitempty"`
	Container string `json:"@container,omitempty"`
	Name      string `json:"name,omitempty"`
	Title     string `json:"title,omitempty"`
}

// populate reads v as json and infers the @context necessary to describe this json file.
// This function is called internally by the Output() function.
func populate(templ Context, v interface{}) (Context, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	ks, err := keys(b)
	if err != nil {
		return nil, err
	}
	ret := make(Context)
	for _, k := range ks {
		fld, ok := templ[k]
		if ok {
			ret[k] = fld
		}
	}
	return ret, nil
}

// reads json and returns a list of unique keys
func keys(b []byte) ([]string, error) {
	dec := json.NewDecoder(bytes.NewBuffer(b))
	// store unique keys in kmap
	kmap := make(map[string]struct{})
	// is Token a key
	var key bool
	// keep track of both object and array parents with a slice of bools:
	//   - an object parent is true, an array parent is false
	parents := make([]bool, 0, 10)
	for {
		t, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		del, ok := t.(json.Delim)
		if ok {
			if del == '{' {
				// push an object parent
				parents = append(parents, true)
			}
			if del == '[' {
				// push an array parent
				parents = append(parents, false)
			}
			if del == '}' || del == ']' {
				if len(parents) == 0 {
					return nil, errors.New("bad json: unexpected } or ] delim")
				}
				// pop the last parent
				parents = parents[:len(parents)-1]
			}
			if len(parents) > 0 && parents[len(parents)-1] {
				// if we are within an object, the next token must be a key
				key = true
			} else {
				// otherwise we are in an array, and the next token is an array entry
				key = false
			}
			continue
		}
		if key {
			str, ok := t.(string)
			if !ok {
				return nil, errors.New("bad json: keys must be strings")
			}
			kmap[str] = struct{}{}
			// if this is a key, then the next token is the value
			key = false
		} else if len(parents) > 0 && parents[len(parents)-1] {
			// if this is a value, and we are within an object, then the next token is a new key
			key = true
		}
	}
	// now turn our map of keys into a sorted slice
	ret := make([]string, len(kmap))
	var i int
	for k := range kmap {
		ret[i] = k
		i++
	}
	sort.Strings(ret)
	return ret, nil
}

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
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/richardlehane/siegfried/pkg/reader"
)

// Siegfried loader. Reads a siegfried file (or droid or fido) to generate generic digital objects
type Siegfried struct {
	Blacklist []string
	reader.Reader
}

// NewSiegfried takes an io Reader for a siegfried/droid/fido results file and an optional blacklist.
// The blacklist is IDs you'd like to exclude e.g. to prune thumbs db files
func NewSiegfried(rdr io.Reader, blacklist ...string) (*Siegfried, error) {
	srdr, err := reader.New(rdr, "")
	return &Siegfried{blacklist, srdr}, err
}

func (s *Siegfried) Load(m *Meta) error {
	var (
		f   reader.File
		err error
	)
	// first check that we just have a single indentifier
	if len(s.Head().Identifiers) != 1 {
		return fmt.Errorf("meta: siegfried loader can only process single IDs, have %d IDs for file %s", len(f.IDs), f.Path)
	}
	mimeField := -1
	for i, v := range s.Head().Fields[0] {
		if v == "mime" || v == "MIME" {
			mimeField = i
			break
		}
	}
	if mimeField < 0 {
		return fmt.Errorf("meta: siegfried loader expects a single identifier that has a MIME field")
	}
	for f, err = s.Next(); err == nil; f, err = s.Next() {
		id := f.IDs[0]
		// check the blacklist
		var isBlackListed bool
		for _, black := range s.Blacklist {
			if black == id.String() {
				isBlackListed = true
				break
			}
		}
		if isBlackListed {
			continue
		}
		fname := filepath.Base(f.Path)
		met, man := NewMetadata(len(m.Index), strings.TrimSuffix(fname, filepath.Ext(fname))), NewManifest()
		modT := NewDateTime(f.Mod)
		met.Created = WrapDate(*modT)
		var hash *Hash
		if s.Head().HashHeader != "" {
			hash = &Hash{
				Algorithm: s.Head().HashHeader,
				Value:     string(f.Hash),
			}
		}
		man.AddVersion([]File{{
			Name:     fname,
			Size:     f.Size,
			Modified: modT,
			MIME:     id.Values()[mimeField],
			PUID:     "http://www.nationalarchives.gov.uk/pronom/" + id.String(),
			Hash:     hash,
		}})
		m.Index = append(m.Index, f.Path)
		m.Metadata[f.Path] = met
		m.Manifest[f.Path] = man
	}
	if err == io.EOF {
		err = nil
	}
	return err
}

// GlobalAccess loader. Applies a simple, global access rule to all digital objects
type GlobalAccess struct {
	AccessDir    int
	AccessEffect string
	Execute      string
}

func (g GlobalAccess) Load(m *Meta) error {
	for _, k := range m.Index {
		m.Manifest[k].AddAR(g.Execute,
			"global",
			true,
			g.AccessDir,
			g.AccessEffect,
			[]FileTarget{{}},
			nil,
			nil)
	}
	return nil
}

// DisposalRule loader. Applies a single disposal rule to all digital objects
func (d DisposalRule) Load(m *Meta) error {
	for _, k := range m.Index {
		m.Metadata[k].DisposalRule = d
	}
	return nil
}

// Series loader. Applies a single series to all digital objects
type Series int

func (s Series) Load(m *Meta) error {
	for _, k := range m.Index {
		m.Metadata[k].Series = ToSeries(int(s))
	}
	return nil
}

// Agency loader. Applies a single agency creator to all digital objects
type Agency struct {
	Name string
	ID   int
}

func (a Agency) Load(m *Meta) error {
	for _, k := range m.Index {
		m.Metadata[k].Creator = MakeAgency(a.Name, a.ID)
	}
	return nil
}

// Title func loader. Allows you to customise a title based on other metadata (e.g. manipulate the file name to make a title)
type TitleFn func(m *Meta, index string) string

func (t TitleFn) Load(m *Meta) error {
	for _, k := range m.Index {
		m.Metadata[k].Title = t(m, k)
	}
	return nil
}

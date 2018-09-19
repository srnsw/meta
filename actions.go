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
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/richardlehane/siegfried"
	"github.com/richardlehane/siegfried/pkg/config"
	"github.com/richardlehane/siegfried/pkg/decompress"
	"github.com/richardlehane/siegfried/pkg/pronom"
	"github.com/srnsw/wincommands"
)

// ManifestCopy copies files and versions as listed in the manifest
// Supply a pathfunc takes the Meta and index as parameters. The output of the pathfunc will be joined with the filename as listed in manifest.
func ManifestCopy(pathfunc func(m *Meta, index string) string) Action {
	return func(m *Meta, target, index string) error {
		man := m.Manifest[index]
		for vidx, v := range man.Versions {
			for _, f := range v.Files {
				if err := wincommands.FileCopy(
					filepath.Join(pathfunc(m, index), f.Name),
					filepath.Join(target, "versions", strconv.Itoa(vidx))+string(filepath.Separator),
					false); err != nil {
					return err
				}
			}
		}
		return nil
	}
}

// IndexPath is an example function that could be supplied to ManifestCopy
// IndexPath assumes that the index holds the full path to a file, so takes the dir of that path.
func IndexPath(m *Meta, index string) string {
	return filepath.Dir(index)
}

// SimpleManifest observes the files in the "versions" folder and builds a simple manifest based on that information.
// It takes a fmtmap argument and path to a siegfried file.
// The fmtmap links file extensions e.g. "pdf" to PUID + mimetype. It can be nil if you want siegfried identification only.
// The siegfried path can be an empty string if you don't want siegfried scanning.
func SimpleManifest(fmtmap map[string][2]string, sfpath string) Action {
	var s *siegfried.Siegfried
	var err error
	if sfpath != "" {
		s, err = siegfried.Load(sfpath)
		if err != nil {
			panic(err)
		}
	}
	if fmtmap == nil {
		fmtmap = make(map[string][2]string)
	}
	return func(m *Meta, target, index string) error {
		man, ok := m.Manifest[index]
		if !ok {
			man = NewManifest()
			m.Manifest[index] = man
		}
		for i := 0; ; i++ {
			_, err := os.Stat(filepath.Join(target, "versions", strconv.Itoa(i)))
			if err != nil {
				if os.IsNotExist(err) {
					return nil
				}
				return err
			}
			files := make([]File, 0, 10)
			err = filepath.Walk(filepath.Join(target, "versions", strconv.Itoa(i)), func(path string, info os.FileInfo, err error) error {
				if info.IsDir() || !info.Mode().IsRegular() {
					return nil
				}
				fnames := filepath.SplitList(info.Name())
				idx := -1
				for ii, vv := range fnames {
					if vv == "versions" {
						idx = ii
						break
					}
				}
				fname := info.Name()
				if idx >= 0 && idx+2 < len(fnames) {
					fname = filepath.Join(fnames[idx+2:]...)
				}
				fmt := [2]string{"UNKNOWN", ""}
				var ok bool
				fmt, ok = fmtmap[strings.TrimPrefix(filepath.Ext(fname), ".")]
				if !ok && s != nil {
					f, err := os.Open(path)
					if err == nil {
						ids, _ := s.Identify(f, path, "")
						if len(ids) == 1 {
							fmt[0] = ids[0].String()
							fmt[1] = ids[0].(pronom.Identification).MIME
						}
					}
					f.Close()
				}
				t := info.ModTime().Truncate(time.Second)
				files = append(files, File{
					Name:     fname,
					Size:     info.Size(),
					Modified: &t,
					MIME:     fmt[1],
					PUID:     ToPUID(fmt[0]),
				})
				return nil
			})
			if err != nil {
				return err
			}
			man.AddVersion(files)
		}
		return nil
	}
}

// Progress prints progress message every n'th item processed
func Progress(i int) Action {
	var n, j int
	return func(m *Meta, target, index string) error {
		n++
		j++
		if j == i {
			j = 0
			log.Printf("Processing number %d (%s)\n", n, index)
		}
		return nil
	}
}

// Decompress takes a path to a siegfried signature file and a pathfunc.
// The pathfunc returns the directory that will be joined to the filename out of the manifest.
// It returns an action that:
// - checks if a PUID (assuming a single version 0/ file 0) is a compressed type and recursively decompresses,
// - adding new files to manifest and copying them to output.
func Decompress(sfpath string) Action {
	var sf *siegfried.Siegfried
	var err error
	if sfpath != "" {
		sf, err = siegfried.Load(sfpath)
		if err != nil {
			panic(err)
		}
	}
	repl := strings.NewReplacer(".zip#", "_zip/")
	return func(m *Meta, target, index string) error {
		man := m.Manifest[index]
		if len(man.Versions) != 1 || len(man.Versions[0].Files) != 1 { // only operate on manifests with a single version/file
			return nil
		}
		if config.IsArchive(strings.TrimPrefix(man.Versions[0].Files[0].PUID, "http://www.nationalarchives.gov.uk/pronom/")) == 0 {
			return nil
		}
		basedir := filepath.Join(target, "versions", "1")
		files := make([]File, 0, 10)
		var idRdr func(rdr io.Reader, name, mime string, sz int64) error
		idRdr = func(rdr io.Reader, name, mime string, sz int64) error {
			buf, err := sf.Buffer(rdr)
			defer sf.Put(buf)
			if err != nil && err.Error() != "empty source" {
				return err
			}
			ids, _ := sf.IdentifyBuffer(buf, nil, name, mime)
			arc := decompress.IsArc(ids)
			if arc > 0 {
				dec, err := decompress.New(arc, buf, name, sz)
				if err != nil {
					return err
				}
				for err = dec.Next(); err == nil; err = dec.Next() {
					err = idRdr(dec.Reader(), dec.Path(), dec.MIME(), dec.Size()) // recurse on the contents of the archive
					if err != nil && err != io.EOF {
						return err
					}
				}
				return nil
			}
			fname := strings.TrimPrefix(name, "#")
			path := fname
			dir := basedir
			fname = repl.Replace(fname)
			var extradirs string
			extradirs, fname = filepath.Split(fname)
			if len(extradirs) > 0 {
				dir = filepath.Join(dir, extradirs)
			}
			os.MkdirAll(dir, 0666)
			f, err := os.Create(filepath.Join(dir, fname))
			if err != nil {
				return err
			}
			_, err = io.Copy(f, buf.Reader())
			if err != nil && err.Error() == "empty source" {
				err = nil
			}
			f.Close()
			if err != nil {
				return err
			}
			fi, err := os.Stat(filepath.Join(dir, fname))
			if err != nil {
				return err
			}
			t := fi.ModTime()
			fmt := [2]string{"UNKNOWN", ""}
			if len(ids) == 1 {
				fmt[0] = ids[0].String()
				fmt[1] = ids[0].(pronom.Identification).MIME
			}
			files = append(files, File{
				Name:     path,
				Size:     fi.Size(),
				Modified: &t,
				MIME:     fmt[1],
				PUID:     ToPUID(fmt[0]),
			})
			return nil
		}
		path := filepath.Join(target, "versions", "0", man.Versions[0].Files[0].Name)
		fi, err := os.Stat(path)
		if err != nil {
			return err
		}
		f, err := os.Open(path)
		defer f.Close()
		if err != nil {
			return err
		}
		err = idRdr(f, "", "", fi.Size())
		if err != nil && err != io.EOF {
			return err
		}
		man.AddVersion(files)
		// now update access rules that have v0f0 as display target
		fts := make([]FileTarget, len(files))
		for i := range fts {
			fts[i][0] = 1
			fts[i][1] = i
		}
		blank := FileTarget{}
		for i, ar := range man.AccessRules {
			str, ok := ar.Display.(string)
			if ok && str == blank.String() {
				man.AccessRules[i].Display = ReferenceFiles(fts)
			}
		}
		return nil
	}
}

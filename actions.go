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
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/richardlehane/siegfried"
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

				}
				t := info.ModTime()
				files = append(files, File{
					Name:     fname,
					Size:     info.Size(),
					Modified: &t,
					MIME:     fmt[1],
					PUID:     fmt[0],
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

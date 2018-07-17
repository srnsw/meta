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
	"path/filepath"
	"strconv"

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

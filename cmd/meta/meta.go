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

// Meta cmd provides a simple tool for creating SIPS from the command line.
// Designed for simple use cases where most of metadata is in a siegfried file.
// Also serves to showcase use of generic loaders and actions available from the meta package.
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"bitbucket.org/srnsw/meta"
)

var (
	blacklistf  = flag.String("blacklist", "", "comma-separated list of IDs to blacklist e.g. x-fmt/111,fmt/10")
	agencyf     = flag.Int("agency", 0, "agency ID e.g. 15")
	agencyNamef = flag.String("agencyName", "", "agency name e.g. State Archives and Records Authority of NSW")
	seriesf     = flag.Int("series", 0, "series ID e.g. 15")
	authorityf  = flag.String("authority", "", "disposal authority e.g. GA28")
	classf      = flag.String("class", "", "disposal class e.g. 1.1.1")
	accessf     = flag.Int("access", 0, "access direction e.g. 15")
	effectf     = flag.String("effect", "", "access direction effect e.g. Early")
	executef    = flag.String("execute", "", "access rule execution date e.g. 2015-01-31")
	outputf     = flag.String("output", "", "output directory e.g. c:/users/richardl/Desktop")
	contentf    = flag.String("content", "", "content directory e.g. c:/users/richardl/stuff")
)

func main() {
	flag.Parse()
	if len(flag.Args()) < 1 {
		fmt.Print("meta: expecting a siegfried results file as input e.g. `meta my_results.yaml`")
		os.Exit(1)
	}
	f, err := os.Open(flag.Arg(0))
	if err != nil {
		fmt.Printf("meta: error opening results file: %v", err)
		os.Exit(1)
	}
	defer f.Close()
	var blacklist []string
	if *blacklistf != "" {
		blacklist = strings.Split(*blacklistf, ",")
	}
	sl, err := meta.NewSiegfried(f, blacklist...)
	if err != nil {
		fmt.Printf("meta: error creating siegfried loader: %v", err)
		os.Exit(1)
	}
	loaders := []meta.Loader{sl}
	// now deal with flags
	if *agencyf > 0 {
		loaders = append(loaders, meta.Agency{*agencyNamef, *agencyf})
	}
	if *seriesf > 0 {
		loaders = append(loaders, meta.Series(*seriesf))
	}
	if *authorityf != "" {
		loaders = append(loaders, meta.DisposalRule{*authorityf, *classf})
	}
	if *accessf > 0 {
		loaders = append(loaders, meta.GlobalAccess{*accessf, *effectf, *executef})
	}
	output := "."
	if *outputf != "" {
		err = os.MkdirAll(*outputf, 0777)
		if err != nil && !os.IsExist(err) {
			fmt.Printf("meta: error creating output folder %s, got %v", *outputf, err)
			os.Exit(1)
		}
		output = *outputf
	}
	pathfunc := meta.IndexPath
	if *contentf != "" {
		pathfunc = func(m *meta.Meta, index string) string {
			return *contentf
		}
	}
	m, err := meta.New(loaders...)
	if err != nil {
		fmt.Printf("meta: error creating meta: %v", err)
		os.Exit(1)
	}
	fmt.Print(m.Output(output, meta.ManifestCopy(pathfunc), meta.Progress(1)))
}

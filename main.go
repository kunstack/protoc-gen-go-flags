// Copyright 2021 Aapeli <aapeli.nian@gmail.com> All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"os"

	"github.com/kunstack/protoc-gen-go-flags/module"
	"github.com/kunstack/protoc-gen-go-flags/version"
	pgs "github.com/lyft/protoc-gen-star/v2"
	pgsgo "github.com/lyft/protoc-gen-star/v2/lang/go"
	"google.golang.org/protobuf/types/pluginpb"
)

func main() {
	// Parse command line flags
	versionFlag := flag.Bool("version", false, "Print version information and exit")
	flag.Parse()

	// Handle --version flag
	if *versionFlag {
		version.Print()
		os.Exit(0)
	}

	var ver = uint64(
		pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL,
	)

	pgs.Init(
		pgs.DebugEnv("DEBUG"),
		pgs.SupportedFeatures(&ver),
	).RegisterModule(
		module.Flags(),
	).RegisterPostProcessor(
		pgsgo.GoFmt(),
	).Render()
}

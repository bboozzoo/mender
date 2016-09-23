// Copyright 2016 Mender Software AS
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

// +build local

package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/mendersoftware/log"
)

const (
	defaultPrefix = "prefix"
)

func mustMkdirAll(path string) {
	err := os.MkdirAll(path, os.FileMode(0755))
	if err != nil {
		panic(fmt.Sprintf("failed to create path %s: %s", path, err))
	}
}

func getPrefixPath() string {
	ep := os.Getenv("MENDER_PREFIX")
	if ep != "" {
		return ep
	}

	p := path.Join(getRunningBinaryPath(), defaultPrefix)
	log.Warnf("MENDER_PREFIX unset, using default '%s'", p)
	return p
}

func getRunningBinaryPath() string {
	return filepath.Dir(os.Args[0])
}

func getDataDirPath() string {
	p := path.Join(getPrefixPath(), "share", "mender")
	mustMkdirAll(p)
	return p
}

func getStateDirPath() string {
	p := path.Join(getPrefixPath(), "var", "lib", "mender")
	mustMkdirAll(p)
	return p
}

func getConfDirPath() string {
	p := path.Join(getPrefixPath(), "etc", "mender")
	mustMkdirAll(p)
	return p
}

func getDevDirPath() string {
	p := path.Join(getPrefixPath(), "dev")
	mustMkdirAll(p)
	return p
}

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
	"encoding/json"
	"io/ioutil"
	"os"
	"path"

	"github.com/mendersoftware/log"
	"github.com/pkg/errors"
)

var (
	fakeEnvPath = path.Join(getStateDirPath(), "fake-env")
)

type fakeEnv struct {
}

func NewEnvironment(cmd Commander) *fakeEnv {
	return &fakeEnv{}
}

func (e *fakeEnv) ReadEnv(names ...string) (BootVars, error) {
	log.Infof("reading environment from %v", fakeEnvPath)
	data, err := ioutil.ReadFile(fakeEnvPath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to load environment")
	}

	var env BootVars
	if len(data) > 0 {
		if err := json.Unmarshal(data, &env); err != nil {
			return nil, errors.Wrapf(err, "failed to decode environment")
		}
	} else {
		env = BootVars{}
	}
	log.Infof("environment: %v", env)

	return env, nil
}

func (e *fakeEnv) WriteEnv(vars BootVars) error {
	log.Infof("writing environment %v to %v", vars, fakeEnvPath)
	env, err := e.ReadEnv()
	if err != nil {
		return errors.Wrapf(err, "failed to load current environment")
	}

	for k, v := range vars {
		env[k] = v
	}

	data, err := json.Marshal(env)
	if err != nil {
		return errors.Wrapf(err, "failed to encode environment")
	}

	if err := ioutil.WriteFile(fakeEnvPath, data, os.FileMode(0644)); err != nil {
		return errors.Wrapf(err, "failed to save environment")
	}
	return nil
}

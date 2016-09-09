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
	"io"

	"github.com/mendersoftware/log"
	"github.com/pkg/errors"
)

type deviceConfig struct {
	rootfsPartA string
	rootfsPartB string
}

type device struct {
	BootEnvReadWriter
	Commander
	parts partitions
}

func NewDevice(env BootEnvReadWriter, sc StatCommander, config deviceConfig) *device {
	partitions := partitions{
		StatCommander:     sc,
		BootEnvReadWriter: env,
		rootfsPartA:       config.rootfsPartA,
		rootfsPartB:       config.rootfsPartB,
		active:            "",
		inactive:          "",
	}
	device := device{env, sc, partitions}
	return &device
}

func (d *device) Reboot() error {
	log.Infof("reboot")
	return nil
}

func (d *device) Rollback() error {
	log.Infof("rollback")
	err := d.WriteEnv(BootVars{
		"upgrade_available": "0",
	})
	if err != nil {
		return errors.Wrapf(err, "failed to update environment after install")
	}
	return nil
}

func (d *device) InstallUpdate(image io.ReadCloser, size int64) error {
	log.Infof("install update of size %v", size)

	err := d.WriteEnv(BootVars{
		"upgrade_available": "1",
	})
	if err != nil {
		return errors.Wrapf(err, "failed to update environment after install")
	}
	return nil
}

func (d *device) EnableUpdatedPartition() error {

	log.Infof("enable updated partition")
	return nil
}

func (d *device) CommitUpdate() error {
	log.Info("Commiting update")
	return nil
}

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
	"os"
	"path"

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
	log.Infof("new device using config: %+v", config)

	partitions := partitions{
		StatCommander:     sc,
		BootEnvReadWriter: env,
		rootfsPartA:       config.rootfsPartA,
		rootfsPartB:       config.rootfsPartB,
		active:            "",
		inactive:          "",
	}

	if partitions.rootfsPartA == "" {
		partitions.rootfsPartA = "mmcblk0p1"
	}
	if partitions.rootfsPartB == "" {
		partitions.rootfsPartB = "mmcblk0p2"
	}
	partitions.rootfsPartA = path.Join(getDevDirPath(),
		partitions.rootfsPartA)
	partitions.rootfsPartB = path.Join(getDevDirPath(),
		partitions.rootfsPartB)

	partitions.active = partitions.rootfsPartA
	partitions.inactive = partitions.rootfsPartB

	log.Infof("partitions mapped to %s %s",
		partitions.rootfsPartA, partitions.rootfsPartB)

	device := device{env, sc, partitions}

	device.doLinks()

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
	log.Infof("install update of size %v to partition %s",
		size, d.parts.inactive)

	f, err := os.Create(d.parts.inactive)
	if err != nil {
		return errors.Wrapf(err, "failed to open file: %s", d.parts.inactive)
	}

	_, err = io.Copy(f, image)
	if err != nil {
		return errors.Wrap(err, "failed to copy update data")
	}

	err = d.WriteEnv(BootVars{
		"upgrade_available": "1",
	})
	if err != nil {
		return errors.Wrapf(err, "failed to update environment after install")
	}
	return nil
}

func (d *device) EnableUpdatedPartition() error {

	log.Infof("enable updated partition")

	d.parts.active, d.parts.inactive = d.parts.inactive, d.parts.active
	d.doLinks()
	return nil
}

func (d *device) CommitUpdate() error {
	log.Info("Commiting update")

	return d.WriteEnv(BootVars{"upgrade_available": "0"})
}

func (d *device) doLinks() {
	ap := path.Join(getDevDirPath(), "active")
	ip := path.Join(getDevDirPath(), "inactive")
	os.Remove(ap)
	os.Remove(ip)
	os.Symlink(path.Base(d.parts.active), ap)
	os.Symlink(path.Base(d.parts.inactive), ip)
}

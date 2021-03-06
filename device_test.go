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
package main

import (
	"errors"
	"os"
	"testing"
)

func Test_commitUpdate(t *testing.T) {
	runner := newTestOSCalls("", 0)
	fakeEnv := uBootEnv{&runner}
	device := device{}
	device.BootEnvReadWriter = &fakeEnv

	if err := device.CommitUpdate(); err != nil {
		t.FailNow()
	}

	runner = newTestOSCalls("", 1)
	if err := device.CommitUpdate(); err == nil {
		t.FailNow()
	}
}

func Test_enableUpdatedPartition_wrongPartitinNumber_fails(t *testing.T) {
	runner := newTestOSCalls("", 0)
	fakeEnv := uBootEnv{&runner}

	testPart := partitions{}
	testPart.inactive = "inactive"

	testDevice := device{}
	testDevice.partitions = &testPart
	testDevice.BootEnvReadWriter = &fakeEnv

	if err := testDevice.EnableUpdatedPartition(); err == nil {
		t.FailNow()
	}
}

func Test_enableUpdatedPartition_correctPartitinNumber(t *testing.T) {
	runner := newTestOSCalls("", 0)
	fakeEnv := uBootEnv{&runner}

	testPart := partitions{}
	testPart.inactive = "inactive2"

	testDevice := device{}
	testDevice.partitions = &testPart
	testDevice.BootEnvReadWriter = &fakeEnv

	if err := testDevice.EnableUpdatedPartition(); err != nil {
		t.FailNow()
	}

	runner = newTestOSCalls("", 1)
	if err := testDevice.EnableUpdatedPartition(); err == nil {
		t.FailNow()
	}
}

func Test_installUpdate_existingAndNonInactivePartition(t *testing.T) {
	testDevice := device{}

	fakePartitions := partitions{}
	fakePartitions.inactive = "/non/existing"
	testDevice.partitions = &fakePartitions

	if err := testDevice.InstallUpdate(nil, 0); err == nil {
		t.FailNow()
	}

	os.Create("inactivePart")
	fakePartitions.inactive = "inactivePart"
	defer os.Remove("inactivePart")

	image, _ := os.Create("imageFile")
	defer os.Remove("imageFile")

	imageContent := "test content"
	image.WriteString(imageContent)
	// rewind to the beginning of file
	image.Seek(0, 0)

	fakePartitions.blockDevSizeGetFunc = func(file *os.File) (uint64, error) { return uint64(len(imageContent)), nil }

	if err := testDevice.InstallUpdate(image, int64(len(imageContent))); err != nil {
		t.FailNow()
	}

	fakePartitions.blockDevSizeGetFunc = func(file *os.File) (uint64, error) { return 0, errors.New("") }
	if err := testDevice.InstallUpdate(image, int64(len(imageContent))); err == nil {
		t.FailNow()
	}
}

func Test_writeImageToPartition_existingAndNonInactivePartition(t *testing.T) {
	os.Create("partitionFile")
	image, _ := os.Create("imageFile")
	imageContent := "test content"
	image.WriteString(imageContent)
	// rewind to the beginning of file
	image.Seek(0, 0)

	defer os.Remove("partitionFile")
	defer os.Remove("imageFile")

	// test writing to correct partition
	err := writeToPartition(image, int64(len(imageContent)), "partitionFile")
	if err != nil {
		t.Fatal(err)
	}

	// test write to non-existing partition
	err = writeToPartition(image, int64(len(imageContent)), "non-existing")
	if err == nil {
		t.Fatal(err)
	}
}

func Test_FetchUpdate_existingAndNonExistingUpdateFile(t *testing.T) {
	image, _ := os.Create("imageFile")
	imageContent := "test content"
	image.WriteString(imageContent)
	file, size, err := FetchUpdateFromFile("imageFile")
	if file == nil || size != int64(len(imageContent)) || err != nil {
		t.FailNow()
	}

	file, _, err = FetchUpdateFromFile("non-existing")
	if file != nil || err == nil {
		t.FailNow()
	}
}

func Test_Rollback_OK(t *testing.T) {
	runner := newTestOSCalls("", 0)
	fakeEnv := uBootEnv{&runner}

	testPart := partitions{}
	testPart.inactive = "part2"

	testDevice := device{}
	testDevice.partitions = &testPart
	testDevice.BootEnvReadWriter = &fakeEnv

	if err := testDevice.Rollback(); err != nil {
		t.FailNow()
	}
}

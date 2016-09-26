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
	"os"

	"github.com/mendersoftware/log"
	"github.com/mendersoftware/mender/utils"
)

var (
	BlockDeviceSize = getBlockDeviceSize
)

// Helper for obtaining the size of a block device.
type BlockDeviceGetSizeFunc func(file *os.File) (uint64, error)

// Wrapper for a block device, implements io.Writer and io.Closer interfaces.
type BlockDevice struct {
	Path string // device path, ex. /dev/mmcblk0p1
	out  *os.File
	size uint64
	w    *utils.LimitedWriter
}

// Write data `p` to underlying block device. Will automatically open device for
// write. Otherwise, behaves like io.Writer.
func (bd *BlockDevice) Write(p []byte) (int, error) {
	if bd.out == nil {
		log.Infof("opening partition %s for writing", bd.Path)
		out, err := os.OpenFile(bd.Path, os.O_WRONLY, 0)
		if err != nil {
			return 0, err
		}

		size, err := BlockDeviceSize(bd.out)
		if err != nil {
			return 0, nil
		}
		log.Infof("   partition %s size: %u", bd.Path, size)

		bd.out = out
		bd.size = size
		bd.w = &utils.LimitedWriter{out, size}
	}

	w, err := bd.w.Write(p)
	if err != nil {
		log.Errorf("written %u out of %u bytes to partition %s: %v",
			w, len(p), bd.Path, err)
	}
	return w, err
}

// Close underlying block device automatically syncing any unwritten data.
// Othewise, behaves like io.Closer.
func (bd *BlockDevice) Close() error {
	if bd.out != nil {
		if err := bd.out.Sync(); err != nil {
			log.Errorf("failed to fsync partition %s: %v", bd.Path, err)
			return err
		}
		if err := bd.out.Close(); err != nil {
			log.Errorf("failed to close partition %s: %v", bd.Path, err)
			return err
		}
	}

	return nil
}

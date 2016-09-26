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

package image

import (
	"fmt"
	"io"

	"github.com/mendersoftware/artifacts/parser"
	"github.com/mendersoftware/artifacts/reader"
	"github.com/pkg/errors"
)

// Install update to output stream `to` by reading & pasing `from` stream.
// Automatically closes the output stream.
func FromArtifact(from io.ReadCloser, to io.Writer, dt string) error {

	ar := areader.NewReader(from)
	defer ar.Close()

	info, err := ar.ReadInfo()
	if err != nil {
		return err
	}

	switch info.Version {
	case 1:
		//
	default:
		return errors.Errorf("unsupported version %v", info.Version)
	}

	hInfo, err := ar.ReadHeaderInfo()
	if err != nil {
		return err
	}

	if len(hInfo.Updates) != 1 {
		return errors.Errorf("unexpected update count %v", len(hInfo.Updates))
	}

	// we will have just one here
	for cnt, update := range hInfo.Updates {
		if update.Type != "rootfs-image" {
			return errors.Errorf("unexpected update type %v", update.Type)
		}
		rp := parser.NewRootfsParser(to, "")
		ar.PushWorker(rp, fmt.Sprintf("%04d", cnt))
	}

	hdr, err := ar.ReadHeader()
	if err != nil {
		return err
	}

	if len(hdr) != 1 {
		return errors.Errorf("expected one header")
	}

	if hdr["0000"].GetDeviceType() != dt {
		return errors.Errorf("unsupported device type %v", hdr["0000"].GetDeviceType())
	}
	updateFiles := hdr["0000"].GetUpdateFiles()
	if len(updateFiles) != 1 {
		return errors.Errorf("expected exacly one update file, got %v", len(updateFiles))
	}

	// check signatures later updateFiles[0].Signature
	// note: updateFiles[0].Size is not there yet
	_, err = ar.ReadData()
	if err != nil {
		return errors.Wrapf(err, "update read failed")
	}

	return nil
}

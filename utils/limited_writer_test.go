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

package utils

import (
	"bytes"
	"io/ioutil"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLimitedWriter(t *testing.T) {
	lw := LimitedWriter{ioutil.Discard, 5}
	assert.NotNil(t, lw)

	// limit to 5 bytes
	_, err := lw.Write([]byte("abcde"))
	assert.NoError(t, err)

	// ENOSPC
	_, err = lw.Write([]byte("foo"))
	assert.EqualError(t, err, syscall.ENOSPC.Error())

	b := &bytes.Buffer{}
	lw = LimitedWriter{b, 5}
	// try to write more than 5 bytes
	w, err := lw.Write([]byte("abcdefg"))
	assert.Equal(t, 5, w)
	assert.EqualError(t, err, syscall.ENOSPC.Error())
	assert.Equal(t, []byte("abcde"), b.Bytes())

	// success write
	b = &bytes.Buffer{}
	lw = LimitedWriter{b, 5}
	w, err = lw.Write([]byte("foo"))
	assert.NoError(t, err)
	assert.Equal(t, len([]byte("foo")), w)

	lw = LimitedWriter{nil, 100}
	_, err = lw.Write([]byte("foo"))
	assert.Error(t, err)
}

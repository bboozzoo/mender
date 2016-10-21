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
	"testing"

	"github.com/mendersoftware/mender/client"
	"github.com/stretchr/testify/assert"
)

func TestInventoryDataDecoder(t *testing.T) {

	idec := NewInventoryDataDecoder()
	assert.NotNil(t, idec)

	_, err := idec.Write([]byte(`foo=bar`))
	assert.NoError(t, err)

	assert.Contains(t, idec.GetInventoryData(), client.InventoryAttribute{"foo", "bar"})

	_, err = idec.Write([]byte(`foo=baz`))
	assert.NoError(t, err)
	assert.Contains(t, idec.data, "foo")
	assert.Contains(t, idec.GetInventoryData(),
		client.InventoryAttribute{"foo", []string{"bar", "baz"}})

	_, err = idec.Write([]byte(`bar=zen`))
	assert.NoError(t, err)
	assert.Contains(t, idec.GetInventoryData(),
		client.InventoryAttribute{"foo", []string{"bar", "baz"}})
	assert.Contains(t, idec.GetInventoryData(), client.InventoryAttribute{"bar", "zen"})

	_, err = idec.Write([]byte("adada"))
	assert.Error(t, err)

	idata := idec.GetInventoryData()
	// added only foo & bar keys
	assert.Len(t, idata, 2)
	assert.Contains(t, idata, client.InventoryAttribute{"foo", []string{"bar", "baz"}})
	assert.Contains(t, idata, client.InventoryAttribute{"bar", "zen"})
}

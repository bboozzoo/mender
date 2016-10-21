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
	"encoding/json"
	"testing"

	"github.com/mendersoftware/mender/client"
	"github.com/stretchr/testify/assert"
)

func TestAuthManager(t *testing.T) {
	ms := NewMemStore()

	cmdr := newTestOSCalls("", 0)
	am := NewAuthManager(ms, "devkey", IdentityDataRunner{
		cmdr: &cmdr,
	})
	assert.IsType(t, &MenderAuthManager{}, am)

	assert.False(t, am.HasKey())
	assert.NoError(t, am.GenerateKey())
	assert.True(t, am.HasKey())

	assert.False(t, am.IsAuthorized())

	code, err := am.AuthToken()
	assert.Equal(t, noAuthToken, code)
	assert.NoError(t, err)

	ms.WriteAll(authTokenName, []byte("footoken"))
	// disable store access
	ms.Disable(true)
	code, err = am.AuthToken()
	assert.Error(t, err)
	ms.Disable(false)

	code, err = am.AuthToken()
	assert.Equal(t, client.AuthToken("footoken"), code)
	assert.NoError(t, err)
}

func TestAuthManagerRequest(t *testing.T) {
	ms := NewMemStore()

	var err error

	badcmdr := newTestOSCalls("mac=foobar", -1)
	am := NewAuthManager(ms, "devkey", IdentityDataRunner{
		cmdr: &badcmdr,
	})
	_, err = am.MakeAuthRequest()
	assert.Error(t, err, "should fail, cannot obtain identity data")
	assert.Contains(t, err.Error(), "identity data")

	cmdr := newTestOSCalls("mac=foobar", 0)
	am = NewAuthManager(ms, "devkey", IdentityDataRunner{
		cmdr: &cmdr,
	})
	_, err = am.MakeAuthRequest()
	assert.Error(t, err, "should fail, no device keys are present")
	assert.Contains(t, err.Error(), "device public key")

	// generate key first
	assert.NoError(t, am.GenerateKey())

	_, err = am.MakeAuthRequest()
	assert.Error(t, err, "should fail, no tenant token present")
	assert.Contains(t, err.Error(), "tenant token")

	// setup tenant token
	ms.WriteAll(authTenantTokenName, []byte("tenant"))
	// setup sequence number
	ms.WriteAll(authSeqName, []byte("12"))

	req, err := am.MakeAuthRequest()
	assert.NoError(t, err)
	assert.NotEmpty(t, req.Data)
	assert.Equal(t, client.AuthToken("tenant"), req.Token)
	assert.NotEmpty(t, req.Signature)

	var ard client.AuthReqData
	err = json.Unmarshal(req.Data, &ard)
	assert.NoError(t, err)

	mam := am.(*MenderAuthManager)
	pempub, _ := mam.keyStore.PublicPEM()
	assert.Equal(t, client.AuthReqData{
		IdData:      "{\"mac\":\"foobar\"}",
		TenantToken: "tenant",
		Pubkey:      pempub,
		SeqNumber:   13,
	}, ard)

	sign, err := mam.keyStore.Sign(req.Data)
	assert.Equal(t, sign, req.Signature)
}

func TestAuthManagerResponse(t *testing.T) {
	ms := NewMemStore()

	cmdr := newTestOSCalls("mac=foobar", 0)
	am := NewAuthManager(ms, "devkey", IdentityDataRunner{
		cmdr: &cmdr,
	})

	var err error
	err = am.RecvAuthResponse([]byte{})
	// should fail with empty response
	assert.Error(t, err)

	// make storage RO
	ms.ReadOnly(true)
	err = am.RecvAuthResponse([]byte("fooresp"))
	assert.Error(t, err)

	ms.ReadOnly(false)
	err = am.RecvAuthResponse([]byte("fooresp"))
	tokdata, err := ms.ReadAll(authTokenName)
	assert.NoError(t, err)
	assert.Equal(t, []byte("fooresp"), tokdata)
	assert.True(t, am.IsAuthorized())
}

package rest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/uniris/uniris-core/pkg/crypto"
	"github.com/uniris/uniris-core/pkg/shared"
)

/*
Scenario: Get shared keys with no emitter public key in the query
	Given no emitter public key in the query
	When I request to get the last shared keys
	THen I get an error
*/
func TestGetSharedKeysWhenMissingPublicKey(t *testing.T) {
	r := gin.New()
	r.GET("/api/sharedkeys", GetSharedKeysHandler(&mockSharedKeyReader{}))

	path := fmt.Sprintf("http://localhost/api/sharedkeys")
	req, _ := http.NewRequest("GET", path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	resBytes, _ := ioutil.ReadAll(w.Body)
	var err httpError
	json.Unmarshal(resBytes, &err)
	assert.Equal(t, "emitter_public_key: public key is empty", err.Error)
	assert.Equal(t, http.StatusText(http.StatusBadRequest), err.Status)
	assert.Equal(t, time.Now().Unix(), err.Timestamp)
}

/*
Scenario: Get shared keys with an invalid emitter public key
	Given an invalid public key
	When I request to get the last shared keys
	THen I get an error
*/
func TestGetSharedKeysWithInvalidPublicKey(t *testing.T) {
	r := gin.New()
	r.GET("/api/sharedkeys", GetSharedKeysHandler(&mockSharedKeyReader{}))

	path := fmt.Sprintf("http://localhost/api/sharedkeys?emitter_public_key=abc")
	req, _ := http.NewRequest("GET", path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	resBytes, _ := ioutil.ReadAll(w.Body)
	var err httpError
	json.Unmarshal(resBytes, &err)
	assert.Equal(t, "emitter_public_key: public key is not in hexadecimal format", err.Error)
	assert.Equal(t, http.StatusText(http.StatusBadRequest), err.Status)
	assert.Equal(t, time.Now().Unix(), err.Timestamp)
}

/*
Scenario: Get shared keys
	Given a valid public key
	When I request to get the last shared keys
	THen I get the last shared keuys
*/
func TestGetSharedKeys(t *testing.T) {

	pub, pv := crypto.GenerateKeys()

	sharedKeyReader := &mockSharedKeyReader{}
	encPv, _ := crypto.Encrypt(pv, pub)
	emKP, _ := shared.NewEmitterCrossKeyPair(encPv, pub)
	nodeKey, _ := shared.NewNodeCrossKeyPair(pub, pv)
	sharedKeyReader.crossNodeKeys = append(sharedKeyReader.crossNodeKeys, nodeKey)
	sharedKeyReader.crossEmitterKeys = append(sharedKeyReader.crossEmitterKeys, emKP)

	r := gin.New()
	r.GET("/api/sharedkeys", GetSharedKeysHandler(sharedKeyReader))

	path := fmt.Sprintf("http://localhost/api/sharedkeys?emitter_public_key=%s", pub)
	req, _ := http.NewRequest("GET", path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	resBytes, _ := ioutil.ReadAll(w.Body)
	var res sharedKeysResponse
	json.Unmarshal(resBytes, &res)
	assert.NotEmpty(t, res.NodePublicKey)
	assert.NotEmpty(t, res.EmitterKeys)

	assert.Equal(t, pub, res.NodePublicKey)

	assert.Len(t, res.EmitterKeys, 1)
	assert.Equal(t, pub, res.EmitterKeys[0].PublicKey)

	emPvKey, _ := crypto.Decrypt(res.EmitterKeys[0].EncryptedPrivateKey, pv)
	assert.Equal(t, pv, emPvKey)
}

type mockSharedKeyReader struct {
	crossNodeKeys    []shared.NodeCrossKeyPair
	crossEmitterKeys []shared.EmitterCrossKeyPair
}

func (r mockSharedKeyReader) EmitterCrossKeypairs() ([]shared.EmitterCrossKeyPair, error) {
	return r.crossEmitterKeys, nil
}

func (r mockSharedKeyReader) FirstNodeCrossKeypair() (shared.NodeCrossKeyPair, error) {
	return r.crossNodeKeys[0], nil
}

func (r mockSharedKeyReader) LastNodeCrossKeypair() (shared.NodeCrossKeyPair, error) {
	return r.crossNodeKeys[len(r.crossNodeKeys)-1], nil
}

func (r mockSharedKeyReader) AuthorizedNodesPublicKeys() ([]string, error) {
	return []string{
		"pub1",
		"pub2",
		"pub3",
		"pub4",
		"pub5",
		"pub6",
		"pub7",
		"pub8",
	}, nil
}

func (r mockSharedKeyReader) CrossEmitterPublicKeys() (pubKeys []string, err error) {
	for _, kp := range r.crossEmitterKeys {
		pubKeys = append(pubKeys, kp.PublicKey())
	}
	return
}

func (r mockSharedKeyReader) FirstEmitterCrossKeypair() (shared.EmitterCrossKeyPair, error) {
	return r.crossEmitterKeys[0], nil
}

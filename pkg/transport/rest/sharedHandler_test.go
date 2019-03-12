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
	r.GET("/api/sharedkeys", GetSharedKeysHandler(&mockTechDB{}))

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
	r.GET("/api/sharedkeys", GetSharedKeysHandler(&mockTechDB{}))

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

	techDB := &mockTechDB{}
	encPv, _ := crypto.Encrypt(pv, pub)
	emKP, _ := shared.NewEmitterKeyPair(encPv, pub)
	nodeKey, _ := shared.NewKeyPair(pub, pv)
	techDB.nodeKeys = append(techDB.nodeKeys, nodeKey)
	techDB.emKeys = append(techDB.emKeys, emKP)

	r := gin.New()
	r.GET("/api/sharedkeys", GetSharedKeysHandler(techDB))

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

type mockTechDB struct {
	emKeys   shared.EmitterKeys
	nodeKeys []shared.KeyPair
}

func (db mockTechDB) EmitterKeys() (shared.EmitterKeys, error) {
	return db.emKeys, nil
}

func (db mockTechDB) NodeLastKeys() (shared.KeyPair, error) {
	return db.nodeKeys[len(db.nodeKeys)-1], nil
}

func (db mockTechDB) NodeFirstKeys() (shared.KeyPair, error) {
	return db.nodeKeys[0], nil
}

func (db mockTechDB) AuthorizedPublicKeys() ([]string, error) {
	return []string{
		"pub1",
		"pub2",
		"pub3",
	}, nil
}

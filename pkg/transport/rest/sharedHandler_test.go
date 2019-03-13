package rest

import (
	"crypto/rand"
	"encoding/hex"
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
	assert.Equal(t, "emitter public key is missing", err.Error)
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
	assert.Equal(t, "emitter public key is not in hexadecimal", err.Error)
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

	pv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

	techDB := &mockTechDB{}

	pvB, _ := pv.Marshal()

	encPv, _ := pub.Encrypt(pvB)
	emKP, _ := shared.NewEmitterKeyPair(encPv, pub)
	nodeKey, _ := shared.NewNodeKeyPair(pub, pv)
	techDB.nodeKeys = append(techDB.nodeKeys, nodeKey)
	techDB.emKeys = append(techDB.emKeys, emKP)

	pubB, _ := pub.Marshal()

	r := gin.New()
	r.GET("/api/sharedkeys", GetSharedKeysHandler(techDB))

	path := fmt.Sprintf("http://localhost/api/sharedkeys?emitter_public_key=%s", hex.EncodeToString(pubB))
	req, _ := http.NewRequest("GET", path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	resBytes, _ := ioutil.ReadAll(w.Body)
	var res sharedKeysResponse
	json.Unmarshal(resBytes, &res)
	assert.NotEmpty(t, res.NodePublicKey)
	assert.NotEmpty(t, res.EmitterKeys)

	assert.EqualValues(t, hex.EncodeToString(pubB), res.NodePublicKey)

	assert.Len(t, res.EmitterKeys, 1)
	assert.EqualValues(t, hex.EncodeToString(pubB), res.EmitterKeys[0].PublicKey)

	encPvBytes, _ := hex.DecodeString(res.EmitterKeys[0].EncryptedPrivateKey)
	emPvKey, _ := pv.Decrypt(encPvBytes)
	assert.EqualValues(t, pvB, emPvKey)
}

type mockTechDB struct {
	emKeys   shared.EmitterKeys
	nodeKeys []shared.NodeKeyPair
}

func (db mockTechDB) EmitterKeys() (shared.EmitterKeys, error) {
	return db.emKeys, nil
}

func (db mockTechDB) NodeLastKeys() (shared.NodeKeyPair, error) {
	return db.nodeKeys[len(db.nodeKeys)-1], nil
}

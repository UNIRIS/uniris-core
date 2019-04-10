package rest

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/uniris/uniris-core/pkg/crypto"
	"github.com/uniris/uniris-core/pkg/logging"
	"github.com/uniris/uniris-core/pkg/shared"
)

/*
Scenario: Get shared keys with no emitter public key in the query
	Given no emitter public key in the query
	When I request to get the last shared keys
	THen I get an error
*/
func TestGetSharedKeysWhenMissingPublicKey(t *testing.T) {
	l := logging.NewLogger("stdout", log.New(os.Stdout, "", 0), "test", net.ParseIP("127.0.0.1"), logging.ErrorLogLevel)
	r := gin.New()
	r.GET("/api/sharedkeys", GetSharedKeysHandler(&mockSharedKeyReader{}, l))

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
	l := logging.NewLogger("stdout", log.New(os.Stdout, "", 0), "test", net.ParseIP("127.0.0.1"), logging.ErrorLogLevel)
	r := gin.New()
	r.GET("/api/sharedkeys", GetSharedKeysHandler(&mockSharedKeyReader{}, l))

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

	sharedKeyReader := &mockSharedKeyReader{}

	pvB, _ := pv.Marshal()

	encPv, _ := pub.Encrypt(pvB)
	emKP, _ := shared.NewEmitterCrossKeyPair(encPv, pub)
	nodeKey, _ := shared.NewNodeCrossKeyPair(pub, pv)
	sharedKeyReader.crossNodeKeys = append(sharedKeyReader.crossNodeKeys, nodeKey)
	sharedKeyReader.crossEmitterKeys = append(sharedKeyReader.crossEmitterKeys, emKP)
	l := logging.NewLogger("stdout", log.New(os.Stdout, "", 0), "test", net.ParseIP("127.0.0.1"), logging.ErrorLogLevel)
	pubB, _ := pub.Marshal()

	r := gin.New()
	r.GET("/api/sharedkeys", GetSharedKeysHandler(sharedKeyReader, l))

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

type mockSharedKeyReader struct {
	crossNodeKeys    []shared.NodeCrossKeyPair
	crossEmitterKeys []shared.EmitterCrossKeyPair
	authKeys         []crypto.PublicKey
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

func (r mockSharedKeyReader) AuthorizedNodesPublicKeys() ([]crypto.PublicKey, error) {
	return r.authKeys, nil
}

func (r mockSharedKeyReader) CrossEmitterPublicKeys() (pubKeys []crypto.PublicKey, err error) {
	for _, kp := range r.crossEmitterKeys {
		pubKeys = append(pubKeys, kp.PublicKey())
	}
	return
}

func (r mockSharedKeyReader) FirstEmitterCrossKeypair() (shared.EmitterCrossKeyPair, error) {
	return r.crossEmitterKeys[0], nil
}

func (r mockSharedKeyReader) IsAuthorizedNode(pub crypto.PublicKey) bool {
	found := false
	for _, k := range r.authKeys {
		if k.Equals(pub) {
			found = true
			break
		}
	}
	return found
}

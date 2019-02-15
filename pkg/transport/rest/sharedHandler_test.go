package rest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	api "github.com/uniris/uniris-core/api/protobuf-spec"
	"github.com/uniris/uniris-core/pkg/crypto"
	"github.com/uniris/uniris-core/pkg/shared"
	"github.com/uniris/uniris-core/pkg/transport/rpc"
	"google.golang.org/grpc"
)

/*
Scenario: Get shared keys with no emitter public key in the query
	Given no emitter public key in the query
	When I request to get the last shared keys
	THen I get an error
*/
func TestGetSharedKeysWhenMissingPublicKey(t *testing.T) {
	r := gin.Default()
	apiGroup := r.Group("/api")
	NewSharedHandler(apiGroup, 1717)

	path := fmt.Sprintf("http://localhost:3000/api/sharedkeys")
	req, _ := http.NewRequest("GET", path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	resBytes, _ := ioutil.ReadAll(w.Body)
	var errorRes map[string]string
	json.Unmarshal(resBytes, &errorRes)
	assert.Equal(t, errorRes["error"], "emitter_public_key: public key is empty")
}

/*
Scenario: Get shared keys with an invalid emitter public key
	Given an invalid public key
	When I request to get the last shared keys
	THen I get an error
*/
func TestGetSharedKeysWithInvalidPublicKey(t *testing.T) {
	r := gin.Default()
	apiGroup := r.Group("/api")
	NewSharedHandler(apiGroup, 1717)

	path := fmt.Sprintf("http://localhost:3000/api/sharedkeys?emitter_public_key=abc")
	req, _ := http.NewRequest("GET", path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	resBytes, _ := ioutil.ReadAll(w.Body)
	var errorRes map[string]string
	json.Unmarshal(resBytes, &errorRes)
	assert.Equal(t, errorRes["error"], "emitter_public_key: public key is not in hexadecimal format")
}

/*
Scenario: Get shared keys with an invalid emitter public key
	Given a valid public key
	When I request to get the last shared keys
	THen I get the last shared keuys
*/
func TestGetSharedKeys(t *testing.T) {
	r := gin.Default()
	apiGroup := r.Group("/api")
	NewSharedHandler(apiGroup, 1717)

	pub, pv := crypto.GenerateKeys()

	techDB := &mockTechDB{}
	encPv, _ := crypto.Encrypt(pv, pub)
	emKP, _ := shared.NewEmitterKeyPair(encPv, pub)
	minerKey, _ := shared.NewMinerKeyPair(pub, pv)
	techDB.minerKeys = append(techDB.minerKeys, minerKey)
	techDB.emKeys = append(techDB.emKeys, emKP)

	pr := rpc.NewPoolRequester(techDB)

	intSrv := rpc.NewInternalServer(techDB, pr)
	lisInt, _ := net.Listen("tcp", ":1717")
	defer lisInt.Close()
	grpcServerInt := grpc.NewServer()
	api.RegisterInternalServiceServer(grpcServerInt, intSrv)
	go grpcServerInt.Serve(lisInt)

	path := fmt.Sprintf("http://localhost:3000/api/sharedkeys?emitter_public_key=%s", pub)
	req, _ := http.NewRequest("GET", path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	resBytes, _ := ioutil.ReadAll(w.Body)
	var res map[string]interface{}
	json.Unmarshal(resBytes, &res)
	assert.NotEmpty(t, res["shared_miner_public_key"])
	assert.NotEmpty(t, res["shared_emitter_keys"])

	assert.Equal(t, pub, res["shared_miner_public_key"])

	emKeys := res["shared_emitter_keys"].([]interface{})
	assert.Len(t, emKeys, 1)
	emKey0 := emKeys[0].(map[string]interface{})
	assert.Equal(t, pub, emKey0["public_key"])

	emPvKey, _ := crypto.Decrypt(emKey0["encrypted_private_key"].(string), pv)
	assert.Equal(t, pv, emPvKey)
}

type mockTechDB struct {
	emKeys    shared.EmitterKeys
	minerKeys []shared.MinerKeyPair
}

func (db mockTechDB) EmitterKeys() (shared.EmitterKeys, error) {
	return db.emKeys, nil
}

func (db mockTechDB) LastMinerKeys() (shared.MinerKeyPair, error) {
	return db.minerKeys[len(db.minerKeys)-1], nil
}

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
	"github.com/uniris/uniris-core/pkg/transaction"
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
	NewSharedHandler(apiGroup, 3545)

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
	NewSharedHandler(apiGroup, 3545)

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

	encPv, _ := crypto.Encrypt(pv, pub)
	minerKP, _ := shared.NewMinerKeyPair(pub, pv)
	emKP, _ := shared.NewEmitterKeyPair(encPv, pub)

	sharedRepo := &mockSharedRepo{}
	sharedRepo.emKeys = []shared.EmitterKeyPair{emKP}
	sharedRepo.minerKeys = minerKP
	poolR := &mockPoolRequester{}

	sharedSrv := shared.NewService(sharedRepo)
	poolingSrv := transaction.NewPoolFindingService(rpc.NewPoolRetriever(sharedSrv))
	miningSrv := transaction.NewMiningService(poolR, poolingSrv, sharedSrv, "127.0.0.1", pub, pv)

	intSrv := rpc.NewInternalServer(poolingSrv, miningSrv, sharedSrv)
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

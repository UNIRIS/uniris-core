package rest

import (
	"encoding/hex"
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
Scenario: Get transactions status with receipt non hexadecimal
	Given an invalid receipt no hexadecimal
	When I want to get the transaction status
	Then I get an error
*/
func TestGetTransactionStatusWithNoHexaReceipt(t *testing.T) {
	r := gin.Default()
	apiGroup := r.Group("/api")
	NewTransactionHandler(apiGroup, 1717)

	path := fmt.Sprintf("http://localhost:3000/api/transaction/abc/status")
	req, _ := http.NewRequest("GET", path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	resBytes, _ := ioutil.ReadAll(w.Body)
	var errorRes map[string]string
	json.Unmarshal(resBytes, &errorRes)
	assert.Equal(t, errorRes["error"], "tx receipt decoding: must be hexadecimal")
}

/*
Scenario: Get transactions status with receipt bad length
	Given an invalid receipt no hexadecimal
	When I want to get the transaction status
	Then I get an error
*/
func TestGetTransactionStatusWithBadReceiptLength(t *testing.T) {
	r := gin.Default()
	apiGroup := r.Group("/api")
	NewTransactionHandler(apiGroup, 1717)

	path := fmt.Sprintf("http://localhost:3000/api/transaction/%s/status", hex.EncodeToString([]byte("abc")))
	req, _ := http.NewRequest("GET", path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	resBytes, _ := ioutil.ReadAll(w.Body)
	var errorRes map[string]string
	json.Unmarshal(resBytes, &errorRes)
	assert.Equal(t, errorRes["error"], "tx receipt decoding: invalid length")
}

/*
Scenario: Get transactions status with unknown transaction hash
	Given an unknown transaction hash
	When I want to get the transaction status
	Then I get the status unknown
*/
func TestGetTransactionStatusUnknown(t *testing.T) {

	pub, pv := crypto.GenerateKeys()

	techDB := &mockTechDB{}
	nodeKey, _ := shared.NewKeyPair(pub, pv)
	techDB.nodeKeys = append(techDB.nodeKeys, nodeKey)

	chainDB := &mockChainDB{}
	locker := &mockLocker{}

	pr := rpc.NewPoolRequester(techDB)

	storageSrv := rpc.NewStorageServer(chainDB, locker, techDB, pr)
	intSrv := rpc.NewInternalServer(techDB, pr)

	//Start transaction server
	lisTx, _ := net.Listen("tcp", ":5000")
	defer lisTx.Close()
	grpcServer := grpc.NewServer()
	api.RegisterStorageServiceServer(grpcServer, storageSrv)
	go grpcServer.Serve(lisTx)

	//Start internal server
	lisInt, _ := net.Listen("tcp", ":1717")
	defer lisInt.Close()
	grpcServerInt := grpc.NewServer()
	api.RegisterInternalServiceServer(grpcServerInt, intSrv)
	go grpcServerInt.Serve(lisInt)

	//Start API
	r := gin.Default()
	apiGroup := r.Group("/api")
	NewTransactionHandler(apiGroup, 1717)

	txReceipt := fmt.Sprintf("%s%s", crypto.HashString("hash"), crypto.HashString("abc"))

	path := fmt.Sprintf("http://localhost:3000/api/transaction/%s/status", txReceipt)
	req, _ := http.NewRequest("GET", path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	resBytes, _ := ioutil.ReadAll(w.Body)
	var res map[string]interface{}
	json.Unmarshal(resBytes, &res)
	assert.Equal(t, res["status"], "UNKNOWN")
	assert.NotEmpty(t, res["timestamp"])
	assert.NotEmpty(t, res["signature"])
}

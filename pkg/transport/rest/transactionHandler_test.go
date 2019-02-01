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
Scenario: Get transactions status with invalid address
	Given an invalid address (not hexa or not hash)
	When I want to get the transaction status
	Then I get an error
*/
func TestGetTransactionStatusWithInvalidAddress(t *testing.T) {
	r := gin.Default()
	apiGroup := r.Group("/api")
	NewTransactionHandler(apiGroup, 3545)

	path := fmt.Sprintf("http://localhost:3000/api/transaction/abc/status/abc")
	req, _ := http.NewRequest("GET", path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	resBytes, _ := ioutil.ReadAll(w.Body)
	var errorRes map[string]string
	json.Unmarshal(resBytes, &errorRes)
	assert.Equal(t, errorRes["error"], "address: hash is not in hexadecimal format")

}

/*
Scenario: Get transactions status with invalid transaction hash
	Given an invalid transaction hash (not hexa or not hash)
	When I want to get the transaction status
	Then I get an error
*/
func TestGetTransactionStatusWithInvalidTxHash(t *testing.T) {
	r := gin.Default()
	apiGroup := r.Group("/api")
	NewTransactionHandler(apiGroup, 3545)

	addr := crypto.HashString("abc")
	path := fmt.Sprintf("http://localhost:3000/api/transaction/%s/status/abc", addr)
	req, _ := http.NewRequest("GET", path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	resBytes, _ := ioutil.ReadAll(w.Body)
	var errorRes map[string]string
	json.Unmarshal(resBytes, &errorRes)
	assert.Equal(t, errorRes["error"], "hash: hash is not in hexadecimal format")
}

/*
Scenario: Get transactions status with unknown transaction hash
	Given an unknown transaction hash
	When I want to get the transaction status
	Then I get the status unknown
*/
func TestGetTransactionStatusUnknown(t *testing.T) {

	pub, pv := crypto.GenerateKeys()

	txRepo := &mockTxRepository{}
	lockRepo := &mockLockRepository{}
	sharedRepo := &mockSharedRepo{}

	encPv, _ := crypto.Encrypt(pv, pub)
	minerKP, _ := shared.NewMinerKeyPair(pub, pv)
	emKP, _ := shared.NewEmitterKeyPair(encPv, pub)
	sharedRepo.emKeys = []shared.EmitterKeyPair{emKP}
	sharedRepo.minerKeys = minerKP

	poolR := &mockPoolRequester{
		repo: txRepo,
	}

	sharedSrv := shared.NewService(sharedRepo)
	poolingSrv := transaction.NewPoolFindingService(rpc.NewPoolRetriever(sharedSrv))
	miningSrv := transaction.NewMiningService(poolR, poolingSrv, sharedSrv, "127.0.0.1", pub, pv)

	storageSrv := transaction.NewStorageService(txRepo, miningSrv)
	lockSrv := transaction.NewLockService(lockRepo)

	txSrv := rpc.NewTransactionServer(storageSrv, lockSrv, miningSrv, sharedSrv)

	//Start transaction server
	lisTx, _ := net.Listen("tcp", ":3545")
	defer lisTx.Close()
	grpcServer := grpc.NewServer()
	api.RegisterTransactionServiceServer(grpcServer, txSrv)
	go grpcServer.Serve(lisTx)

	//Start internal server
	intSrv := rpc.NewInternalServer(poolingSrv, miningSrv, sharedSrv)
	lisInt, _ := net.Listen("tcp", ":1717")
	defer lisInt.Close()
	grpcServerInt := grpc.NewServer()
	api.RegisterInternalServiceServer(grpcServerInt, intSrv)
	go grpcServerInt.Serve(lisInt)

	//Start API
	r := gin.Default()
	apiGroup := r.Group("/api")
	NewTransactionHandler(apiGroup, 1717)

	addr := crypto.HashString("abc")
	txHash := crypto.HashString("hash")

	path := fmt.Sprintf("http://localhost:3000/api/transaction/%s/status/%s", addr, txHash)
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

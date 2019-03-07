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
	"time"

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
	r := gin.New()
	r.GET("/api/transaction/:txReceipt/status", GetTransactionStatusHandler(&mockTechDB{}))

	path := fmt.Sprintf("http://localhost/api/transaction/abc/status")
	req, _ := http.NewRequest("GET", path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	resBytes, _ := ioutil.ReadAll(w.Body)
	var err httpError
	json.Unmarshal(resBytes, &err)
	assert.Equal(t, "tx receipt decoding: must be hexadecimal", err.Error)
	assert.Equal(t, http.StatusText(http.StatusBadRequest), err.Status)
	assert.Equal(t, time.Now().Unix(), err.Timestamp)
}

/*
Scenario: Get transactions status with receipt bad length
	Given an invalid receipt no hexadecimal
	When I want to get the transaction status
	Then I get an error
*/
func TestGetTransactionStatusWithBadReceiptLength(t *testing.T) {
	r := gin.New()
	r.GET("/api/transaction/:txReceipt/status", GetTransactionStatusHandler(&mockTechDB{}))

	path := fmt.Sprintf("http://localhost/api/transaction/%s/status", hex.EncodeToString([]byte("abc")))
	req, _ := http.NewRequest("GET", path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	resBytes, _ := ioutil.ReadAll(w.Body)
	var err httpError
	json.Unmarshal(resBytes, &err)
	assert.Equal(t, "tx receipt decoding: invalid length", err.Error)
	assert.Equal(t, http.StatusText(http.StatusBadRequest), err.Status)
	assert.Equal(t, time.Now().Unix(), err.Timestamp)
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

	pr := rpc.NewPoolRequester(techDB)

	txSrv := rpc.NewTransactionService(chainDB, techDB, pr, pub, pv)

	//Start transaction server
	lisTx, _ := net.Listen("tcp", ":5000")
	defer lisTx.Close()
	grpcServer := grpc.NewServer()
	api.RegisterTransactionServiceServer(grpcServer, txSrv)
	go grpcServer.Serve(lisTx)

	//Start API
	r := gin.New()
	r.GET("/api/transaction/:txReceipt/status", GetTransactionStatusHandler(techDB))

	txReceipt := fmt.Sprintf("%s%s", crypto.HashString("hash"), crypto.HashString("abc"))

	path := fmt.Sprintf("http://localhost/api/transaction/%s/status", txReceipt)
	req, _ := http.NewRequest("GET", path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	resBytes, _ := ioutil.ReadAll(w.Body)
	var res transactionStatusResponse
	json.Unmarshal(resBytes, &res)
	assert.Equal(t, "UNKNOWN", res.Status)
	// assert.Equal(t, time.Now().Unix(), res.Timestamp)
	// assert.NotEmpty(t, res.Signature)
}

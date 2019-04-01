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

	"github.com/uniris/uniris-core/pkg/consensus"
	"github.com/uniris/uniris-core/pkg/logging"
	"github.com/uniris/uniris-core/pkg/shared"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	api "github.com/uniris/uniris-core/api/protobuf-spec"
	"github.com/uniris/uniris-core/pkg/crypto"
	"github.com/uniris/uniris-core/pkg/transport/rpc"
	"google.golang.org/grpc"

	stdcrypto "crypto"
)

/*
Scenario: Get transactions status with receipt non hexadecimal
	Given an invalid receipt no hexadecimal
	When I want to get the transaction status
	Then I get an error
*/
func TestGetTransactionStatusWithNoHexaReceipt(t *testing.T) {
	l := logging.NewLogger("stdout", log.New(os.Stdout, "", 0), "test", net.ParseIP("127.0.0.1"), logging.ErrorLogLevel)
	r := gin.New()
	r.GET("/api/transaction/:txReceipt/status", GetTransactionStatusHandler(&mockSharedKeyReader{}, &mockNodeDatabase{}, l))

	path := fmt.Sprintf("http://localhost/api/transaction/abc/status")
	req, _ := http.NewRequest("GET", path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	resBytes, _ := ioutil.ReadAll(w.Body)
	var err httpError
	json.Unmarshal(resBytes, &err)
	assert.Equal(t, "transaction receipt is not in hexadecimal", err.Error)
	assert.Equal(t, http.StatusText(http.StatusBadRequest), err.Status)
	assert.Equal(t, time.Now().Unix(), err.Timestamp)
}

/*
Scenario: Get transactions status with receipt bad format
	Given an invalid receipt with a bad hash algorithm identifier
	When I want to get the transaction status
	Then I get an error
*/
func TestGetTransactionStatusWithBadFormat(t *testing.T) {
	l := logging.NewLogger("stdout", log.New(os.Stdout, "", 0), "test", net.ParseIP("127.0.0.1"), logging.ErrorLogLevel)
	r := gin.New()
	r.GET("/api/transaction/:txReceipt/status", GetTransactionStatusHandler(&mockSharedKeyReader{}, &mockNodeDatabase{}, l))

	path := fmt.Sprintf("http://localhost/api/transaction/%s/status", hex.EncodeToString([]byte("abc")))
	req, _ := http.NewRequest("GET", path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	resBytes, _ := ioutil.ReadAll(w.Body)
	var err httpError
	json.Unmarshal(resBytes, &err)
	assert.Equal(t, "transaction receipt is invalid", err.Error)
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
	l := logging.NewLogger("stdout", log.New(os.Stdout, "", 0), "test", net.ParseIP("127.0.0.1"), logging.ErrorLogLevel)
	r := gin.New()
	r.GET("/api/transaction/:txReceipt/status", GetTransactionStatusHandler(&mockSharedKeyReader{}, &mockNodeDatabase{}, l))

	receipt := make([]byte, 10)
	receipt[0] = byte(int(stdcrypto.SHA256))
	copy(receipt[1:], []byte("bc6"))

	path := fmt.Sprintf("http://localhost/api/transaction/%s/status", hex.EncodeToString(receipt))
	req, _ := http.NewRequest("GET", path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	resBytes, _ := ioutil.ReadAll(w.Body)
	var err httpError
	json.Unmarshal(resBytes, &err)
	assert.Equal(t, "transaction receipt is invalid", err.Error)
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

	pv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

	sharedKeyReader := &mockSharedKeyReader{}
	nodeKey, _ := shared.NewNodeCrossKeyPair(pub, pv)
	sharedKeyReader.crossNodeKeys = append(sharedKeyReader.crossNodeKeys, nodeKey)

	chainDB := &mockChainDB{}

	l := logging.NewLogger("stdout", log.New(os.Stdout, "", 0), "test", net.ParseIP("127.0.0.1"), logging.ErrorLogLevel)
	pr := rpc.NewPoolRequester(sharedKeyReader, l)

	nodeReader := &mockNodeDatabase{
		nodes: []consensus.Node{
			consensus.NewNode(net.ParseIP("127.0.0.1"), 5000, pub, consensus.NodeOK, "", 300, "1.0", 0, 1, 30.0, -10.0, consensus.GeoPatch{}, true),
		},
	}

	txSrv := rpc.NewTransactionService(chainDB, sharedKeyReader, nodeReader, pr, pub, pv, l)

	//Start transaction server
	lisTx, _ := net.Listen("tcp", ":5000")
	defer lisTx.Close()
	grpcServer := grpc.NewServer()
	api.RegisterTransactionServiceServer(grpcServer, txSrv)
	go grpcServer.Serve(lisTx)

	//Start API
	r := gin.New()
	r.GET("/api/transaction/:txReceipt/status", GetTransactionStatusHandler(sharedKeyReader, nodeReader, l))

	addr := crypto.Hash([]byte("hash"))
	hash := crypto.Hash([]byte("abc"))

	path := fmt.Sprintf("http://localhost/api/transaction/%s/status", fmt.Sprintf("%x%x", addr, hash))
	req, _ := http.NewRequest("GET", path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	resBytes, _ := ioutil.ReadAll(w.Body)
	var res transactionStatusResponse
	json.Unmarshal(resBytes, &res)
	assert.Equal(t, "UNKNOWN", res.Status)
}

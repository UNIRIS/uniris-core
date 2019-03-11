package rest

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"
	"time"

	"github.com/uniris/uniris-core/pkg/chain"
	"github.com/uniris/uniris-core/pkg/consensus"

	api "github.com/uniris/uniris-core/api/protobuf-spec"
	"github.com/uniris/uniris-core/pkg/shared"
	"google.golang.org/grpc"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/uniris/uniris-core/pkg/crypto"
	"github.com/uniris/uniris-core/pkg/transport/rpc"
)

/*
Scenario: Get account request with an ID hash not a valid hash
	Given an invalid hash (not hexa)
	When I want to request to retrieve an account
	Then I got a 400 (Bad request) response status and an error message
*/
func TestGetAccountWhenInvalidHash(t *testing.T) {
	r := gin.New()
	r.GET("/api/account/:idHash", GetAccountHandler(&mockTechDB{}))

	path := fmt.Sprintf("http://localhost/api/account/abc")
	req, _ := http.NewRequest("GET", path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	resBytes, _ := ioutil.ReadAll(w.Body)
	var err httpError
	json.Unmarshal(resBytes, &err)
	assert.Equal(t, "id hash: must be hexadecimal", err.Error)
	assert.Equal(t, http.StatusText(http.StatusBadRequest), err.Status)
	assert.Equal(t, time.Now().Unix(), err.Timestamp)
}

/*
Scenario: Get account request with an invalid idSignature
	Given a hash and an invalid idSignature
	When I want to request to retrieve an account
	Then I got a 400 (Bad request) response status and an error message
*/
func TestGetAccountWhenInvalidSignature(t *testing.T) {
	r := gin.New()
	r.GET("/api/account/:idHash", GetAccountHandler(&mockTechDB{}))

	path1 := fmt.Sprintf("http://localhost/api/account/%s", crypto.HashString("abc"))
	req1, _ := http.NewRequest("GET", path1, nil)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusBadRequest, w1.Code)
	resBytes, _ := ioutil.ReadAll(w1.Body)
	var err httpError
	json.Unmarshal(resBytes, &err)
	assert.Equal(t, "signature request: signature is empty", err.Error)
	assert.Equal(t, http.StatusText(http.StatusBadRequest), err.Status)
	assert.Equal(t, time.Now().Unix(), err.Timestamp)

	path2 := fmt.Sprintf("http://localhost/api/account/%s?signature=%s", crypto.HashString("abc"), "idSig")
	req2, _ := http.NewRequest("GET", path2, nil)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusBadRequest, w2.Code)
	resBytes2, _ := ioutil.ReadAll(w2.Body)
	var err2 httpError
	json.Unmarshal(resBytes2, &err2)
	assert.Equal(t, "signature request: signature is not in hexadecimal format", err2.Error)
	assert.Equal(t, http.StatusText(http.StatusBadRequest), err2.Status)
	assert.Equal(t, time.Now().Unix(), err2.Timestamp)

	path3 := fmt.Sprintf("http://localhost/api/account/%s?signature=%s", crypto.HashString("abc"), hex.EncodeToString([]byte("idSig")))
	req3, _ := http.NewRequest("GET", path3, nil)
	w3 := httptest.NewRecorder()
	r.ServeHTTP(w3, req3)
	assert.Equal(t, http.StatusBadRequest, w3.Code)
	resBytes3, _ := ioutil.ReadAll(w3.Body)
	var err3 httpError
	json.Unmarshal(resBytes3, &err3)
	assert.Equal(t, "signature request: signature is not valid", err3.Error)
	assert.Equal(t, http.StatusText(http.StatusBadRequest), err3.Status)
	assert.Equal(t, time.Now().Unix(), err3.Timestamp)
}

/*
Scenario: Get account request with an ID not existing
	Given an ID hash and a valid idSignature related to no real ID transaction
	When I want to request to retrieve an account
	Then I got a 404 (Not found) response status and an error message
*/
func TestGetAccountWhenIDNotExist(t *testing.T) {

	pub, pv := crypto.GenerateKeys()

	techDB := &mockTechDB{}
	nodeKey, _ := shared.NewKeyPair(pub, pv)
	techDB.nodeKeys = append(techDB.nodeKeys, nodeKey)
	emKey, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("ov")), pub)
	techDB.emKeys = append(techDB.emKeys, emKey)

	chainDB := &mockChainDB{}
	pr := &mockPoolRequester{
		repo: chainDB,
	}

	lis, _ := net.Listen("tcp", ":5000")
	defer lis.Close()
	grpcServer := grpc.NewServer()
	api.RegisterTransactionServiceServer(grpcServer, rpc.NewTransactionService(chainDB, techDB, pr, pub, pv))
	go grpcServer.Serve(lis)

	r := gin.New()
	r.GET("/api/account/:idHash", GetAccountHandler(techDB))

	idHash := crypto.HashString("abc")
	encIDHash, _ := crypto.Encrypt(idHash, pub)
	sig, _ := crypto.Sign(encIDHash, pv)

	path := fmt.Sprintf("http://localhost/api/account/%s?signature=%s", encIDHash, sig)
	req, _ := http.NewRequest("GET", path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	resBytes, _ := ioutil.ReadAll(w.Body)
	var err httpError
	json.Unmarshal(resBytes, &err)
	assert.Equal(t, "ID: transaction does not exist", err.Error)
	assert.Equal(t, http.StatusText(http.StatusNotFound), err.Status)
	assert.Equal(t, time.Now().Unix(), err.Timestamp)
}

/*
Scenario: Get account request with a Keychain not existing
	Given an ID hash and a valid idSignature related to no real Keychain transaction
	When I want to request to retrieve an account
	Then I got a 404 (Not found) response status and an error message
*/
func TestGetAccountWhenKeychainNotExist(t *testing.T) {
	pub, pv := crypto.GenerateKeys()

	chainDB := &mockChainDB{}
	techDB := &mockTechDB{}
	nodeKey, _ := shared.NewKeyPair(pub, pv)
	techDB.nodeKeys = append(techDB.nodeKeys, nodeKey)
	emKey, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("ov")), pub)
	techDB.emKeys = append(techDB.emKeys, emKey)
	pr := &mockPoolRequester{
		repo: chainDB,
	}

	lis, _ := net.Listen("tcp", ":5000")
	defer lis.Close()
	grpcServer := grpc.NewServer()
	api.RegisterTransactionServiceServer(grpcServer, rpc.NewTransactionService(chainDB, techDB, pr, pub, pv))
	go grpcServer.Serve(lis)

	//Start API
	r := gin.New()
	r.GET("/api/account/:idHash", GetAccountHandler(techDB))

	//Create transactions
	encAddr, _ := crypto.Encrypt(hex.EncodeToString([]byte("addr")), pub)

	idData := map[string]string{
		"encrypted_address_by_node": encAddr,
		"encrypted_address_by_id":   encAddr,
		"encrypted_aes_key":         hex.EncodeToString([]byte("aes_key")),
	}
	prop, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("encPV")), pub)
	idHash := crypto.HashString("abc")
	idTxRaw := map[string]interface{}{
		"addr":                    idHash,
		"data":                    idData,
		"timestamp":               time.Now().Unix(),
		"type":                    chain.IDTransactionType,
		"public_key":              pub,
		"em_shared_keys_proposal": prop,
	}
	idtxBytes, _ := json.Marshal(idTxRaw)
	idSig, _ := crypto.Sign(string(idtxBytes), pv)
	idTxRaw["signature"] = idSig

	idtxbytesWithSig, _ := json.Marshal(idTxRaw)
	emSig, _ := crypto.Sign(string(idtxbytesWithSig), pv)
	idTxRaw["em_signature"] = emSig

	idtxBytes, _ = json.Marshal(idTxRaw)

	idTx, _ := chain.NewTransaction(idHash, chain.IDTransactionType, idData, time.Now(), pub, prop, idSig, emSig, crypto.HashBytes(idtxBytes))
	id, _ := chain.NewID(idTx)
	chainDB.WriteID(id)

	encIDHash, _ := crypto.Encrypt(idHash, pub)
	sig, _ := crypto.Sign(encIDHash, pv)

	path := fmt.Sprintf("http://localhost/api/account/%s?signature=%s", encIDHash, sig)
	req, _ := http.NewRequest("GET", path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	resBytes, _ := ioutil.ReadAll(w.Body)
	var err httpError
	json.Unmarshal(resBytes, &err)
	assert.Equal(t, "Keychain: transaction does not exist", err.Error)
	assert.Equal(t, http.StatusText(http.StatusNotFound), err.Status)
	assert.Equal(t, time.Now().Unix(), err.Timestamp)
}

/*
Scenario: Get an account after its creation
	Given an account created (keychain and ID transaction)
	When I want to retrieve it
	Then I can get encrypted wallet and encrypted aes key
*/
func TestGetAccount(t *testing.T) {
	pub, pv := crypto.GenerateKeys()

	chainDB := &mockChainDB{}
	techDB := &mockTechDB{}
	nodeKey, _ := shared.NewKeyPair(pub, pv)
	techDB.nodeKeys = append(techDB.nodeKeys, nodeKey)
	emKey, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("ov")), pub)
	techDB.emKeys = append(techDB.emKeys, emKey)
	pr := &mockPoolRequester{
		repo: chainDB,
	}

	lis, _ := net.Listen("tcp", ":5000")
	defer lis.Close()
	grpcServer := grpc.NewServer()
	api.RegisterTransactionServiceServer(grpcServer, rpc.NewTransactionService(chainDB, techDB, pr, pub, pv))
	go grpcServer.Serve(lis)

	//Create transactions
	addr := crypto.HashString("addr")
	encAddr, _ := crypto.Encrypt(addr, pub)

	idData := map[string]string{
		"encrypted_address_by_node": encAddr,
		"encrypted_address_by_id":   encAddr,
		"encrypted_aes_key":         hex.EncodeToString([]byte("aes_key")),
	}
	prop, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("encPV")), pub)

	idHash := crypto.HashString("abc")
	idTxRaw := map[string]interface{}{
		"addr":                    idHash,
		"data":                    idData,
		"timestamp":               time.Now().Unix(),
		"type":                    chain.IDTransactionType,
		"public_key":              pub,
		"em_shared_keys_proposal": prop,
	}
	idtxBytes, _ := json.Marshal(idTxRaw)
	idSig, _ := crypto.Sign(string(idtxBytes), pv)
	idTxRaw["signature"] = idSig

	idtxbytesWithSig, _ := json.Marshal(idTxRaw)
	emSig, _ := crypto.Sign(string(idtxbytesWithSig), pv)
	idTxRaw["em_signature"] = emSig

	idtxBytes, _ = json.Marshal(idTxRaw)

	idTx, _ := chain.NewTransaction(idHash, chain.IDTransactionType, idData, time.Now(), pub, prop, idSig, emSig, crypto.HashBytes(idtxBytes))
	id, _ := chain.NewID(idTx)
	chainDB.WriteID(id)

	keychainData := map[string]string{
		"encrypted_address_by_node": encAddr,
		"encrypted_wallet":          hex.EncodeToString([]byte("wallet")),
	}

	keychainTxRaw := map[string]interface{}{
		"addr":                    crypto.HashString("addr"),
		"data":                    keychainData,
		"timestamp":               time.Now().Unix(),
		"type":                    chain.KeychainTransactionType,
		"public_key":              pub,
		"em_shared_keys_proposal": prop,
	}
	txKeychainBytes, _ := json.Marshal(keychainTxRaw)
	keychainSig, _ := crypto.Sign(string(txKeychainBytes), pv)
	keychainTxRaw["signature"] = keychainSig

	keychaintxbytesWithSig, _ := json.Marshal(keychainTxRaw)
	keychainEmSig, _ := crypto.Sign(string(keychaintxbytesWithSig), pv)
	keychainTxRaw["em_signature"] = keychainEmSig

	keychainTxRaw["em_signature"] = keychainSig
	txKeychainBytes, _ = json.Marshal(keychainTxRaw)

	keychainTx, _ := chain.NewTransaction(addr, chain.KeychainTransactionType, keychainData, time.Now(), pub, prop, keychainSig, keychainEmSig, crypto.HashBytes(txKeychainBytes))
	keychain, _ := chain.NewKeychain(keychainTx)
	chainDB.WriteKeychain(keychain)

	encIDHash, _ := crypto.Encrypt(idHash, pub)
	sig, _ := crypto.Sign(encIDHash, pv)

	r := gin.New()
	r.GET("/api/account/:idHash", GetAccountHandler(techDB))
	path := fmt.Sprintf("http://localhost/api/account/%s?signature=%s", encIDHash, sig)

	req, _ := http.NewRequest("GET", path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	resBytes, _ := ioutil.ReadAll(w.Body)
	var res accountFindResponse
	json.Unmarshal(resBytes, &res)
	assert.Equal(t, hex.EncodeToString([]byte("aes_key")), res.EncryptedAESKey)
	assert.Equal(t, hex.EncodeToString([]byte("wallet")), res.EncryptedWallet)
	assert.Equal(t, time.Now().Unix(), res.Timestamp)
	assert.NotEmpty(t, res.Signature)

	resBytes, _ = json.Marshal(struct {
		EncryptedWallet string `json:"encrypted_wallet"`
		EncryptedAESKey string `json:"encrypted_aes_key"`
		Timestamp       int64  `json:"timestamp"`
	}{
		EncryptedAESKey: res.EncryptedAESKey,
		EncryptedWallet: res.EncryptedWallet,
		Timestamp:       res.Timestamp,
	})
	assert.Nil(t, crypto.VerifySignature(string(resBytes), techDB.nodeKeys[0].PublicKey(), res.Signature))
}

/*
Scenario: Create account request with an invalid encrypted ID
	Given an invalid encrypted id (not hexa)
	When I want to request to create an account
	Then I got a 400 (Bad request) response status and an error message
*/
func TestCreationAccountWhenInvalidID(t *testing.T) {
	r := gin.New()
	r.POST("/api/account", CreateAccountHandler(&mockTechDB{}))

	form, _ := json.Marshal(map[string]string{
		"encrypted_id":       "abc",
		"encrypted_keychain": "abc",
		"signature":          "abc",
	})

	path := "http://localhost/api/account"
	req, _ := http.NewRequest("POST", path, bytes.NewBuffer(form))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	resBytes, _ := ioutil.ReadAll(w.Body)
	var err httpError
	json.Unmarshal(resBytes, &err)
	assert.Equal(t, "encrypted id: must be hexadecimal", err.Error)
	assert.Equal(t, http.StatusText(http.StatusBadRequest), err.Status)
	assert.Equal(t, time.Now().Unix(), err.Timestamp)

}

/*
Scenario: Create account request with an invalid encrypted Keychain
	Given an invalid encrypted Keychain (not hexa)
	When I want to request to create an account
	Then I got a 400 (Bad request) response status and an error message
*/
func TestCreationAccountWhenKeychainInvalid(t *testing.T) {
	r := gin.New()
	r.POST("/api/account", CreateAccountHandler(&mockTechDB{}))

	form, _ := json.Marshal(map[string]string{
		"encrypted_id":       hex.EncodeToString([]byte("id")),
		"encrypted_keychain": "abc",
		"signature":          "abc",
	})

	path := "http://localhost/api/account"
	req, _ := http.NewRequest("POST", path, bytes.NewBuffer(form))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	resBytes, _ := ioutil.ReadAll(w.Body)
	var err httpError
	json.Unmarshal(resBytes, &err)
	assert.Equal(t, "encrypted keychain: must be hexadecimal", err.Error)
	assert.Equal(t, http.StatusText(http.StatusBadRequest), err.Status)
	assert.Equal(t, time.Now().Unix(), err.Timestamp)

}

/*
Scenario: Create account request with an invalid signature
	Given an invalid signature (not hexa and not valid)
	When I want to request to create an account
	Then I got a 400 (Bad request) response status and an error message
*/
func TestCreationAccountWhenSignatureInvalid(t *testing.T) {

	pub, pv := crypto.GenerateKeys()

	techDB := &mockTechDB{}
	nodeKey, _ := shared.NewKeyPair(pub, pv)
	emKey, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("pv")), pub)
	techDB.nodeKeys = append(techDB.nodeKeys, nodeKey)
	techDB.emKeys = append(techDB.emKeys, emKey)

	r := gin.New()
	r.POST("/api/account", CreateAccountHandler(techDB))

	form, _ := json.Marshal(map[string]string{
		"encrypted_id":       hex.EncodeToString([]byte("id")),
		"encrypted_keychain": hex.EncodeToString([]byte("keychain")),
		"signature":          "abc",
	})

	path := "http://localhost/api/account"
	req1, _ := http.NewRequest("POST", path, bytes.NewBuffer(form))
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)

	assert.Equal(t, http.StatusBadRequest, w1.Code)
	resBytes, _ := ioutil.ReadAll(w1.Body)

	var err httpError
	json.Unmarshal(resBytes, &err)
	assert.Equal(t, "signature request: signature is not in hexadecimal format", err.Error)
	assert.Equal(t, http.StatusText(http.StatusBadRequest), err.Status)
	assert.Equal(t, time.Now().Unix(), err.Timestamp)

	form2, _ := json.Marshal(map[string]string{
		"encrypted_id":       hex.EncodeToString([]byte("id")),
		"encrypted_keychain": hex.EncodeToString([]byte("keychain")),
		"signature":          hex.EncodeToString([]byte("abc")),
	})

	req2, _ := http.NewRequest("POST", path, bytes.NewBuffer(form2))
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusBadRequest, w2.Code)
	resBytes2, _ := ioutil.ReadAll(w2.Body)
	var err2 httpError
	json.Unmarshal(resBytes2, &err2)
	assert.Equal(t, "signature request: signature is not valid", err2.Error)
	assert.Equal(t, http.StatusText(http.StatusBadRequest), err2.Status)
	assert.Equal(t, time.Now().Unix(), err2.Timestamp)

	sig, _ := crypto.Sign(hex.EncodeToString([]byte("hello")), pv)

	form3, _ := json.Marshal(map[string]string{
		"encrypted_id":       hex.EncodeToString([]byte("id")),
		"encrypted_keychain": hex.EncodeToString([]byte("keychain")),
		"signature":          sig,
	})

	req3, _ := http.NewRequest("POST", path, bytes.NewBuffer(form3))
	w3 := httptest.NewRecorder()
	r.ServeHTTP(w3, req3)

	assert.Equal(t, http.StatusBadRequest, w3.Code)
	resBytes3, _ := ioutil.ReadAll(w3.Body)
	var err3 httpError
	json.Unmarshal(resBytes3, &err3)
	assert.Equal(t, "signature request: signature is not valid", err3.Error)
	assert.Equal(t, http.StatusText(http.StatusBadRequest), err3.Status)
	assert.Equal(t, time.Now().Unix(), err3.Timestamp)

}

/*
Scenario: Create an account including ID and keychain transaction
	Given a valid ID and keychain transaction
	When I want to create an account
	Then two transaction are created (ID/Keychain) and the data is stored
*/
func TestCreateAccount(t *testing.T) {
	pub, pv := crypto.GenerateKeys()

	chainDB := &mockChainDB{}
	techDB := &mockTechDB{}

	encKey, _ := crypto.Encrypt(pv, pub)
	emKey, _ := shared.NewEmitterKeyPair(encKey, pub)
	techDB.emKeys = append(techDB.emKeys, emKey)

	nodeKey, _ := shared.NewKeyPair(pub, pv)
	techDB.nodeKeys = append(techDB.nodeKeys, nodeKey)

	pr := &mockPoolRequester{
		repo: chainDB,
	}

	txSrv := rpc.NewTransactionService(chainDB, techDB, pr, pub, pv)

	//Start transaction server
	lis, _ := net.Listen("tcp", ":5000")
	defer lis.Close()
	grpcServer := grpc.NewServer()
	api.RegisterTransactionServiceServer(grpcServer, txSrv)
	go grpcServer.Serve(lis)

	//Start API
	r := gin.New()
	r.POST("/api/account", CreateAccountHandler(techDB))

	//Create transactions
	addr := crypto.HashString("addr")
	encAddr, _ := crypto.Encrypt(addr, pub)

	idTx := map[string]interface{}{
		"addr": crypto.HashString("abc"),
		"data": map[string]string{
			"encrypted_address_by_node": encAddr,
			"encrypted_address_by_id":   encAddr,
			"encrypted_aes_key":         hex.EncodeToString([]byte("aes_key")),
		},
		"timestamp":  time.Now().Unix(),
		"type":       int(chain.IDTransactionType),
		"public_key": pub,
		"em_shared_keys_proposal": map[string]string{
			"encrypted_private_key": hex.EncodeToString([]byte("encPV")),
			"public_key":            pub,
		},
	}

	idTxBytes, _ := json.Marshal(idTx)
	idSig, _ := crypto.Sign(string(idTxBytes), pv)
	idTx["signature"] = idSig

	idTxByteWithSig, _ := json.Marshal(idTx)
	emSig, _ := crypto.Sign(string(idTxByteWithSig), pv)
	idTx["em_signature"] = emSig

	keychainTx := map[string]interface{}{
		"addr": addr,
		"data": map[string]string{
			"encrypted_address_by_node": encAddr,
			"encrypted_wallet":          hex.EncodeToString([]byte("wallet")),
		},
		"timestamp":  time.Now().Unix(),
		"type":       int(chain.KeychainTransactionType),
		"public_key": pub,
		"em_shared_keys_proposal": map[string]string{
			"encrypted_private_key": hex.EncodeToString([]byte("encPV")),
			"public_key":            pub,
		},
	}

	keychainTxBytes, _ := json.Marshal(keychainTx)
	keychainSig, _ := crypto.Sign(string(keychainTxBytes), pv)
	keychainTx["signature"] = keychainSig

	keychainTxByteWithSig, _ := json.Marshal(keychainTx)
	keychainEmSig, _ := crypto.Sign(string(keychainTxByteWithSig), pv)
	keychainTx["em_signature"] = keychainEmSig

	idTxBytes, _ = json.Marshal(idTx)
	keychainTxBytes, _ = json.Marshal(keychainTx)

	encryptedID, _ := crypto.Encrypt(string(idTxBytes), pub)
	encryptedKeychain, _ := crypto.Encrypt(string(keychainTxBytes), pub)

	form := map[string]string{
		"encrypted_id":       encryptedID,
		"encrypted_keychain": encryptedKeychain,
	}
	formB, _ := json.Marshal(form)
	sig, _ := crypto.Sign(string(formB), pv)

	form["signature"] = sig

	formB, _ = json.Marshal(form)

	path := "http://localhost/api/account"
	req, _ := http.NewRequest("POST", path, bytes.NewBuffer(formB))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	resBytes, _ := ioutil.ReadAll(w.Body)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resTx accountCreationResponse
	json.Unmarshal(resBytes, &resTx)

	assert.NotEmpty(t, resTx.IDTransaction.TransactionReceipt)
	assert.NotEmpty(t, resTx.IDTransaction.Timestamp)
	assert.NotEmpty(t, resTx.IDTransaction.Signature)
	assert.Equal(t, time.Now().Unix(), resTx.IDTransaction.Timestamp)

	idTxHash := crypto.HashBytes(idTxBytes)
	assert.Equal(t, fmt.Sprintf("%s%s", crypto.HashString("abc"), idTxHash), resTx.IDTransaction.TransactionReceipt)

	idResBytes, _ := json.Marshal(transactionResponse{
		TransactionReceipt: resTx.IDTransaction.TransactionReceipt,
		Timestamp:          resTx.IDTransaction.Timestamp,
	})
	assert.Nil(t, crypto.VerifySignature(string(idResBytes), techDB.nodeKeys[0].PublicKey(), resTx.IDTransaction.Signature))

	assert.NotEmpty(t, resTx.KeychainTransaction.TransactionReceipt)
	assert.NotEmpty(t, resTx.KeychainTransaction.Timestamp)
	assert.NotEmpty(t, resTx.KeychainTransaction.Signature)
	assert.Equal(t, time.Now().Unix(), resTx.KeychainTransaction.Timestamp)

	keychainTxHash := crypto.HashBytes(keychainTxBytes)
	assert.Equal(t, fmt.Sprintf("%s%s", addr, keychainTxHash), resTx.KeychainTransaction.TransactionReceipt)

	keychainResBytes, _ := json.Marshal(transactionResponse{
		TransactionReceipt: resTx.KeychainTransaction.TransactionReceipt,
		Timestamp:          resTx.KeychainTransaction.Timestamp,
	})
	assert.Nil(t, crypto.VerifySignature(string(keychainResBytes), techDB.nodeKeys[0].PublicKey(), resTx.KeychainTransaction.Signature))

	time.Sleep(50 * time.Millisecond)

	assert.Len(t, chainDB.keychains, 1)
	assert.Equal(t, addr, chainDB.keychains[0].Address())
	assert.Len(t, chainDB.ids, 1)
	assert.Equal(t, crypto.HashString("abc"), chainDB.ids[0].Address())

}

type mockPoolRequester struct {
	stores []chain.Transaction
	repo   *mockChainDB
}

func (pr mockPoolRequester) RequestLastTransaction(pool consensus.Pool, txAddr string, txType chain.TransactionType) (*chain.Transaction, error) {
	switch txType {
	case chain.KeychainTransactionType:
		kc, _ := pr.repo.LastKeychain(txAddr)
		if kc == nil {
			return nil, nil
		}
		return &kc.Transaction, nil
	case chain.IDTransactionType:
		id, _ := pr.repo.ID(txAddr)
		if id == nil {
			return nil, nil
		}
		return &id.Transaction, nil
	}

	return nil, nil
}

func (pr mockPoolRequester) RequestTransactionTimeLock(pool consensus.Pool, txHash string, txAddr string, masterPublicKey string) error {
	return nil
}

func (pr mockPoolRequester) RequestTransactionUnlock(pool consensus.Pool, txHash string, txAddr string) error {
	return nil
}

func (pr mockPoolRequester) RequestTransactionValidations(pool consensus.Pool, tx chain.Transaction, minValids int, masterValid chain.MasterValidation) ([]chain.Validation, error) {
	pub, pv := crypto.GenerateKeys()

	vRaw := map[string]interface{}{
		"status":     chain.ValidationOK,
		"public_key": pub,
		"timestamp":  time.Now().Unix(),
	}
	vBytes, _ := json.Marshal(vRaw)
	sig, _ := crypto.Sign(string(vBytes), pv)
	v, _ := chain.NewValidation(chain.ValidationOK, time.Now(), pub, sig)

	return []chain.Validation{v}, nil
}

func (pr *mockPoolRequester) RequestTransactionStorage(pool consensus.Pool, minReplicas int, tx chain.Transaction) error {
	pr.stores = append(pr.stores, tx)
	if tx.TransactionType() == chain.KeychainTransactionType {
		k, _ := chain.NewKeychain(tx)
		pr.repo.keychains = append(pr.repo.keychains, k)
	}
	if tx.TransactionType() == chain.IDTransactionType {
		id, _ := chain.NewID(tx)
		pr.repo.ids = append(pr.repo.ids, id)
	}
	return nil
}

type mockChainDB struct {
	kos       []chain.Transaction
	keychains []chain.Keychain
	ids       []chain.ID
}

func (r mockChainDB) LastKeychain(txAddr string) (*chain.Keychain, error) {
	sort.Slice(r.keychains, func(i, j int) bool {
		return r.keychains[i].Timestamp().Unix() > r.keychains[j].Timestamp().Unix()
	})

	if len(r.keychains) > 0 {
		return &r.keychains[0], nil
	}
	return nil, nil
}

func (r mockChainDB) FullKeychain(txAddr string) (*chain.Keychain, error) {
	sort.Slice(r.keychains, func(i, j int) bool {
		return r.keychains[i].Timestamp().Unix() > r.keychains[j].Timestamp().Unix()
	})

	if len(r.keychains) > 0 {
		return &r.keychains[0], nil
	}
	return nil, nil
}

func (r mockChainDB) KeychainByHash(txHash string) (*chain.Keychain, error) {
	for _, tx := range r.keychains {
		if tx.TransactionHash() == txHash {
			return &tx, nil
		}
	}
	return nil, nil
}

func (r mockChainDB) IDByHash(txHash string) (*chain.ID, error) {
	for _, tx := range r.ids {
		if tx.TransactionHash() == txHash {
			return &tx, nil
		}
	}
	return nil, nil
}

func (r mockChainDB) ID(addr string) (*chain.ID, error) {
	for _, tx := range r.ids {
		if tx.Address() == addr {
			return &tx, nil
		}
	}
	return nil, nil
}

func (r mockChainDB) KOByHash(txHash string) (*chain.Transaction, error) {
	for _, tx := range r.kos {
		if tx.TransactionHash() == txHash {
			return &tx, nil
		}
	}
	return nil, nil
}

func (r *mockChainDB) WriteKeychain(kc chain.Keychain) error {
	r.keychains = append(r.keychains, kc)
	return nil
}

func (r *mockChainDB) WriteID(id chain.ID) error {
	r.ids = append(r.ids, id)
	return nil
}

func (r *mockChainDB) WriteKO(tx chain.Transaction) error {
	r.kos = append(r.kos, tx)
	return nil
}

package rest

import (
	"bytes"
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
	"github.com/uniris/uniris-core/pkg/logging"
	"github.com/uniris/uniris-core/pkg/transport/rpc"
)

/*
Scenario: Get account request with an ID hash not a valid hash
	Given an invalid hash (not hexadecimal)
	When I want to request to retrieve an account
	Then I got a 400 (Bad request) response status and an error message
*/
func TestGetAccountWhenInvalidHash(t *testing.T) {
	r := gin.New()
	l := logging.NewLogger(log.New(os.Stdout, "", 0), "test", net.ParseIP("127.0.0.1"), "debug")
	r.GET("/api/account/:idHash", GetAccountHandler(&mockSharedKeyReader{}, &mockNodeDatabase{}, l))

	path := fmt.Sprintf("http://localhost/api/account/abc")
	req, _ := http.NewRequest("GET", path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	resBytes, _ := ioutil.ReadAll(w.Body)
	var err httpError
	json.Unmarshal(resBytes, &err)
	assert.Equal(t, "id hash is not in hexadecimal", err.Error)
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

	_, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
	emK, _ := shared.NewEmitterCrossKeyPair([]byte("enc"), pub)
	l := logging.NewLogger(log.New(os.Stdout, "", 0), "test", net.ParseIP("127.0.0.1"), "debug")

	r := gin.New()
	r.GET("/api/account/:idHash", GetAccountHandler(&mockSharedKeyReader{
		crossEmitterKeys: []shared.EmitterCrossKeyPair{emK},
	}, &mockNodeDatabase{}, l))

	path1 := fmt.Sprintf("http://localhost/api/account/%s", hex.EncodeToString(crypto.Hash(([]byte("abc")))))
	req1, _ := http.NewRequest("GET", path1, nil)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusBadRequest, w1.Code)
	resBytes, _ := ioutil.ReadAll(w1.Body)
	var err httpError
	json.Unmarshal(resBytes, &err)
	assert.Equal(t, "signature is missing", err.Error)
	assert.Equal(t, http.StatusText(http.StatusBadRequest), err.Status)
	assert.Equal(t, time.Now().Unix(), err.Timestamp)

	path2 := fmt.Sprintf("http://localhost/api/account/%s?signature=%s", hex.EncodeToString(crypto.Hash([]byte("abc"))), "idSig")
	req2, _ := http.NewRequest("GET", path2, nil)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusBadRequest, w2.Code)
	resBytes2, _ := ioutil.ReadAll(w2.Body)
	var err2 httpError
	json.Unmarshal(resBytes2, &err2)
	assert.Equal(t, "signature is not in hexadecimal", err2.Error)
	assert.Equal(t, http.StatusText(http.StatusBadRequest), err2.Status)
	assert.Equal(t, time.Now().Unix(), err2.Timestamp)

	path3 := fmt.Sprintf("http://localhost/api/account/%s?signature=%s", hex.EncodeToString(crypto.Hash([]byte("abc"))), hex.EncodeToString([]byte("idSig")))
	req3, _ := http.NewRequest("GET", path3, nil)
	w3 := httptest.NewRecorder()
	r.ServeHTTP(w3, req3)
	assert.Equal(t, http.StatusBadRequest, w3.Code)
	resBytes3, _ := ioutil.ReadAll(w3.Body)
	var err3 httpError
	json.Unmarshal(resBytes3, &err3)
	assert.Equal(t, "signature is invalid", err3.Error)
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

	pv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

	sharedKeyReader := &mockSharedKeyReader{}
	nodeKey, _ := shared.NewNodeCrossKeyPair(pub, pv)
	sharedKeyReader.crossNodeKeys = append(sharedKeyReader.crossNodeKeys, nodeKey)
	emKey, _ := shared.NewEmitterCrossKeyPair([]byte("encPV"), pub)
	sharedKeyReader.crossEmitterKeys = append(sharedKeyReader.crossEmitterKeys, emKey)

	chainDB := &mockChainDB{}
	pr := &mockPoolRequester{
		repo: chainDB,
	}

	nodeReader := &mockNodeDatabase{
		nodes: []consensus.Node{
			consensus.NewNode(net.ParseIP("127.0.0.1"), 5000, pub, consensus.NodeOK, "", 300, "1.0", 0, 1, 30.0, -10.0, consensus.GeoPatch{}, true),
		},
	}

	lis, _ := net.Listen("tcp", ":5000")
	defer lis.Close()
	grpcServer := grpc.NewServer()
	l := logging.NewLogger(log.New(os.Stdout, "", 0), "test", net.ParseIP("127.0.0.1"), "debug")
	api.RegisterTransactionServiceServer(grpcServer, rpc.NewTransactionService(chainDB, sharedKeyReader, nodeReader, pr, pub, pv, l))
	go grpcServer.Serve(lis)

	r := gin.New()
	r.GET("/api/account/:idHash", GetAccountHandler(sharedKeyReader, nodeReader, l))

	idHash := crypto.Hash([]byte("abc"))
	encIDHash, _ := pub.Encrypt(idHash)
	sig, _ := pv.Sign(encIDHash)

	path := fmt.Sprintf("http://localhost/api/account/%s?signature=%s", hex.EncodeToString(encIDHash), hex.EncodeToString(sig))
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
	pv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

	chainDB := &mockChainDB{}
	sharedKeyReader := &mockSharedKeyReader{}
	nodeKey, _ := shared.NewNodeCrossKeyPair(pub, pv)
	sharedKeyReader.crossNodeKeys = append(sharedKeyReader.crossNodeKeys, nodeKey)
	emKey, _ := shared.NewEmitterCrossKeyPair([]byte("encPV"), pub)
	sharedKeyReader.crossEmitterKeys = append(sharedKeyReader.crossEmitterKeys, emKey)
	pr := &mockPoolRequester{
		repo: chainDB,
	}

	nodeReader := &mockNodeDatabase{
		nodes: []consensus.Node{
			consensus.NewNode(net.ParseIP("127.0.0.1"), 5000, pub, consensus.NodeOK, "", 300, "1.0", 0, 1, 30.0, -10.0, consensus.GeoPatch{}, true),
		},
	}

	lis, _ := net.Listen("tcp", ":5000")
	defer lis.Close()
	l := logging.NewLogger(log.New(os.Stdout, "", 0), "test", net.ParseIP("127.0.0.1"), "debug")
	grpcServer := grpc.NewServer()
	api.RegisterTransactionServiceServer(grpcServer, rpc.NewTransactionService(chainDB, sharedKeyReader, nodeReader, pr, pub, pv, l))
	go grpcServer.Serve(lis)

	//Start API
	r := gin.New()
	r.GET("/api/account/:idHash", GetAccountHandler(sharedKeyReader, nodeReader, l))

	//Create transactions

	prop, _ := shared.NewEmitterCrossKeyPair([]byte("encpv"), pub)
	idHash := crypto.Hash([]byte("abc"))

	encAddr, _ := pub.Encrypt([]byte("addr"))

	idData := map[string][]byte{
		"encrypted_address_by_node": encAddr,
		"encrypted_address_by_id":   encAddr,
		"encrypted_aes_key":         []byte("aes_key"),
	}

	pubB, _ := pub.Marshal()

	idTxRaw := map[string]interface{}{
		"addr": hex.EncodeToString(idHash),
		"data": map[string]string{
			"encrypted_address_by_node": hex.EncodeToString(encAddr),
			"encrypted_address_by_id":   hex.EncodeToString(encAddr),
			"encrypted_aes_key":         hex.EncodeToString([]byte("aes_key")),
		},
		"timestamp":  time.Now().Unix(),
		"type":       chain.IDTransactionType,
		"public_key": hex.EncodeToString(pubB),
		"em_shared_keys_proposal": map[string]string{
			"encrypted_private_key": hex.EncodeToString(prop.EncryptedPrivateKey()),
			"public_key":            hex.EncodeToString(pubB),
		},
	}
	idtxBytes, _ := json.Marshal(idTxRaw)
	idSig, _ := pv.Sign(idtxBytes)
	idTxRaw["signature"] = hex.EncodeToString(idSig)

	idtxbytesWithSig, _ := json.Marshal(idTxRaw)
	emSig, _ := pv.Sign(idtxbytesWithSig)
	idTxRaw["em_signature"] = hex.EncodeToString(emSig)

	idtxBytes, _ = json.Marshal(idTxRaw)

	idTx, _ := chain.NewTransaction(idHash, chain.IDTransactionType, idData, time.Now(), pub, prop, idSig, emSig, crypto.Hash(idtxBytes))
	id, _ := chain.NewID(idTx)
	chainDB.WriteID(id)

	encIDHash, _ := pub.Encrypt(idHash)
	sig, _ := pv.Sign(encIDHash)

	path := fmt.Sprintf("http://localhost/api/account/%s?signature=%s", hex.EncodeToString(encIDHash), hex.EncodeToString(sig))

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
	pv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
	pubB, _ := pub.Marshal()

	chainDB := &mockChainDB{}
	sharedKeyReader := &mockSharedKeyReader{}
	nodeKey, _ := shared.NewNodeCrossKeyPair(pub, pv)
	sharedKeyReader.crossNodeKeys = append(sharedKeyReader.crossNodeKeys, nodeKey)
	emKey, _ := shared.NewEmitterCrossKeyPair([]byte("encPV"), pub)
	sharedKeyReader.crossEmitterKeys = append(sharedKeyReader.crossEmitterKeys, emKey)
	pr := &mockPoolRequester{
		repo: chainDB,
	}

	nodeReader := &mockNodeDatabase{
		nodes: []consensus.Node{
			consensus.NewNode(net.ParseIP("127.0.0.1"), 5000, pub, consensus.NodeOK, "", 300, "1.0", 0, 1, 30.0, -10.0, consensus.GeoPatch{}, true),
		},
	}

	lis, _ := net.Listen("tcp", ":5000")
	defer lis.Close()
	l := logging.NewLogger(log.New(os.Stdout, "", 0), "test", net.ParseIP("127.0.0.1"), "debug")
	grpcServer := grpc.NewServer()
	api.RegisterTransactionServiceServer(grpcServer, rpc.NewTransactionService(chainDB, sharedKeyReader, nodeReader, pr, pub, pv, l))
	go grpcServer.Serve(lis)

	//Create transactions
	addr := crypto.Hash([]byte("addr"))
	encAddr, _ := pub.Encrypt(addr)

	idData := map[string][]byte{
		"encrypted_address_by_node": encAddr,
		"encrypted_address_by_id":   encAddr,
		"encrypted_aes_key":         []byte("aes_key"),
	}
	prop, _ := shared.NewEmitterCrossKeyPair([]byte("encpv"), pub)

	idHash := crypto.Hash([]byte("abc"))
	idTxRaw := map[string]interface{}{
		"addr": hex.EncodeToString(idHash),
		"data": map[string]string{
			"encrypted_address_by_node": hex.EncodeToString(encAddr),
			"encrypted_address_by_id":   hex.EncodeToString(encAddr),
			"encrypted_aes_key":         hex.EncodeToString([]byte("aes_key")),
		},
		"timestamp":  time.Now().Unix(),
		"type":       chain.IDTransactionType,
		"public_key": hex.EncodeToString(pubB),
		"em_shared_keys_proposal": map[string]string{
			"encrypted_private_key": hex.EncodeToString(prop.EncryptedPrivateKey()),
			"public_key":            hex.EncodeToString(pubB),
		},
	}
	idtxBytes, _ := json.Marshal(idTxRaw)
	idSig, _ := pv.Sign(idtxBytes)
	idTxRaw["signature"] = hex.EncodeToString(idSig)

	idtxbytesWithSig, _ := json.Marshal(idTxRaw)
	emSig, _ := pv.Sign(idtxbytesWithSig)
	idTxRaw["em_signature"] = hex.EncodeToString(emSig)

	idtxBytes, _ = json.Marshal(idTxRaw)

	idTx, _ := chain.NewTransaction(idHash, chain.IDTransactionType, idData, time.Now(), pub, prop, idSig, emSig, crypto.Hash(idtxBytes))
	id, _ := chain.NewID(idTx)
	chainDB.WriteID(id)

	keychainData := map[string][]byte{
		"encrypted_address_by_node": encAddr,
		"encrypted_wallet":          []byte("wallet"),
	}

	keychainTxRaw := map[string]interface{}{
		"addr": hex.EncodeToString(crypto.Hash([]byte("addr"))),
		"data": map[string]string{
			"encrypted_address_by_node": hex.EncodeToString(encAddr),
			"encrypted_wallet":          hex.EncodeToString([]byte("wallet")),
		},
		"timestamp":  time.Now().Unix(),
		"type":       chain.KeychainTransactionType,
		"public_key": hex.EncodeToString(pubB),
		"em_shared_keys_proposal": map[string]string{
			"encrypted_private_key": hex.EncodeToString(prop.EncryptedPrivateKey()),
			"public_key":            hex.EncodeToString(pubB),
		},
	}
	txKeychainBytes, _ := json.Marshal(keychainTxRaw)
	keychainSig, _ := pv.Sign(txKeychainBytes)
	keychainTxRaw["signature"] = hex.EncodeToString(keychainSig)

	keychaintxbytesWithSig, _ := json.Marshal(keychainTxRaw)
	keychainEmSig, _ := pv.Sign(keychaintxbytesWithSig)
	keychainTxRaw["em_signature"] = hex.EncodeToString(keychainEmSig)

	keychainTxRaw["em_signature"] = keychainSig
	txKeychainBytes, _ = json.Marshal(keychainTxRaw)

	keychainTx, _ := chain.NewTransaction(addr, chain.KeychainTransactionType, keychainData, time.Now(), pub, prop, keychainSig, keychainEmSig, crypto.Hash(txKeychainBytes))
	keychain, _ := chain.NewKeychain(keychainTx)
	chainDB.WriteKeychain(keychain)

	encIDHash, _ := pub.Encrypt(idHash)
	sig, _ := pv.Sign(encIDHash)

	r := gin.New()
	r.GET("/api/account/:idHash", GetAccountHandler(sharedKeyReader, nodeReader, l))
	path := fmt.Sprintf("http://localhost/api/account/%s?signature=%s", hex.EncodeToString(encIDHash), hex.EncodeToString(sig))

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
	sigBytes, _ := hex.DecodeString(res.Signature)
	firstNodeKey, _ := sharedKeyReader.FirstNodeCrossKeypair()
	assert.True(t, firstNodeKey.PublicKey().Verify(resBytes, sigBytes))
}

/*
Scenario: Create account request with an invalid signature
	Given an invalid signature (not hexadecimal and not valid)
	When I want to request to create an account
	Then I got a 400 (Bad request) response status and an error message
*/
func TestCreationAccountWhenInvalidID(t *testing.T) {
	r := gin.New()

	pv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
	l := logging.NewLogger(log.New(os.Stdout, "", 0), "test", net.ParseIP("127.0.0.1"), "debug")

	sharedKeyReader := &mockSharedKeyReader{}
	nodeKey, _ := shared.NewNodeCrossKeyPair(pub, pv)
	emKey, _ := shared.NewEmitterCrossKeyPair([]byte("pv"), pub)
	sharedKeyReader.crossNodeKeys = append(sharedKeyReader.crossNodeKeys, nodeKey)
	sharedKeyReader.crossEmitterKeys = append(sharedKeyReader.crossEmitterKeys, emKey)

	r.POST("/api/account", CreateAccountHandler(sharedKeyReader, mockNodeDatabase{}, pub, pv, l))

	form, _ := json.Marshal(map[string]string{
		"encrypted_id":       hex.EncodeToString([]byte("id")),
		"encrypted_keychain": hex.EncodeToString([]byte("keychain")),
		"signature":          hex.EncodeToString([]byte("abc")),
	})

	path := "http://localhost/api/account"
	req1, _ := http.NewRequest("POST", path, bytes.NewBuffer(form))
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)

	assert.Equal(t, http.StatusBadRequest, w1.Code)
	resBytes, _ := ioutil.ReadAll(w1.Body)

	var err httpError
	json.Unmarshal(resBytes, &err)
	assert.Equal(t, "signature is invalid", err.Error)
	assert.Equal(t, http.StatusText(http.StatusBadRequest), err.Status)
	assert.Equal(t, time.Now().Unix(), err.Timestamp)
}

/*
Scenario: Create account request with an invalid encrypted transaction raw
	Given an invalid transaction raw (not encrypted, not JSON or missing fields)
	When I want to request to create an account
	Then I got a 400 (Bad request) response status and an error message
*/
func TestCreationAccountWhenInvalidTransactionRaw(t *testing.T) {

	pv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

	sharedKeyReader := &mockSharedKeyReader{}
	nodeKey, _ := shared.NewNodeCrossKeyPair(pub, pv)
	emKey, _ := shared.NewEmitterCrossKeyPair([]byte("pv"), pub)
	sharedKeyReader.crossNodeKeys = append(sharedKeyReader.crossNodeKeys, nodeKey)
	sharedKeyReader.crossEmitterKeys = append(sharedKeyReader.crossEmitterKeys, emKey)
	l := logging.NewLogger(log.New(os.Stdout, "", 0), "test", net.ParseIP("127.0.0.1"), "debug")

	r := gin.New()
	r.POST("/api/account", CreateAccountHandler(sharedKeyReader, mockNodeDatabase{}, pub, pv, l))

	form := accountCreationRequest{
		EncryptedID:       hex.EncodeToString([]byte("abc")),
		EncryptedKeychain: hex.EncodeToString([]byte("abc")),
	}

	formBytes, _ := json.Marshal(form)
	sig, _ := pv.Sign(formBytes)
	form.Signature = hex.EncodeToString(sig)
	formBytes, _ = json.Marshal(form)

	path := "http://localhost/api/account"
	req, _ := http.NewRequest("POST", path, bytes.NewBuffer(formBytes))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	resBytes, _ := ioutil.ReadAll(w.Body)

	var err httpError
	json.Unmarshal(resBytes, &err)
	assert.Equal(t, "invalid message", err.Error)
	assert.Equal(t, http.StatusText(http.StatusBadRequest), err.Status)
	assert.Equal(t, time.Now().Unix(), err.Timestamp)

	encID, _ := pub.Encrypt([]byte("abc"))
	form = accountCreationRequest{
		EncryptedID:       hex.EncodeToString(encID),
		EncryptedKeychain: hex.EncodeToString([]byte("abc")),
	}
	formBytes, _ = json.Marshal(form)
	sig, _ = pv.Sign(formBytes)
	form.Signature = hex.EncodeToString(sig)
	formBytes, _ = json.Marshal(form)

	req2, _ := http.NewRequest("POST", path, bytes.NewBuffer(formBytes))
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	resBytes2, _ := ioutil.ReadAll(w2.Body)
	assert.Equal(t, http.StatusBadRequest, w2.Code)

	var err2 httpError
	json.Unmarshal(resBytes2, &err2)
	assert.Equal(t, "invalid JSON", err2.Error)
	assert.Equal(t, http.StatusText(http.StatusBadRequest), err2.Status)
	assert.Equal(t, time.Now().Unix(), err2.Timestamp)

	fakeJSON, _ := json.Marshal(map[string]string{
		"hello": "text",
	})
	encID, _ = pub.Encrypt(fakeJSON)
	form = accountCreationRequest{
		EncryptedID:       hex.EncodeToString(encID),
		EncryptedKeychain: hex.EncodeToString([]byte("abc")),
	}
	formBytes, _ = json.Marshal(form)
	sig, _ = pv.Sign(formBytes)
	form.Signature = hex.EncodeToString(sig)
	formBytes, _ = json.Marshal(form)

	req3, _ := http.NewRequest("POST", path, bytes.NewBuffer(formBytes))
	w3 := httptest.NewRecorder()
	r.ServeHTTP(w3, req3)
	resBytes3, _ := ioutil.ReadAll(w3.Body)
	assert.Equal(t, http.StatusBadRequest, w3.Code)

	var err3 httpError
	json.Unmarshal(resBytes3, &err3)
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
	pv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

	chainDB := &mockChainDB{}
	sharedKeyReader := &mockSharedKeyReader{}

	pvB, _ := pv.Marshal()

	encKey, _ := pub.Encrypt(pvB)
	emKey, _ := shared.NewEmitterCrossKeyPair(encKey, pub)
	sharedKeyReader.crossEmitterKeys = append(sharedKeyReader.crossEmitterKeys, emKey)

	nodeKey, _ := shared.NewNodeCrossKeyPair(pub, pv)
	sharedKeyReader.crossNodeKeys = append(sharedKeyReader.crossNodeKeys, nodeKey)

	pr := &mockPoolRequester{
		repo: chainDB,
	}

	_, pub1, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
	_, pub2, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
	_, pub3, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
	_, pub4, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
	_, pub5, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
	_, pub6, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
	_, pub7, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
	_, pub8, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

	sharedKeyReader.authKeys = []crypto.PublicKey{pub1, pub2, pub3, pub4, pub5, pub6, pub7, pub8}

	nodeReader := &mockNodeDatabase{
		nodes: []consensus.Node{
			consensus.NewNode(net.ParseIP("127.0.0.1"), 5000, pub1, consensus.NodeOK, "", 300, "1.0", 0, 1, 30.0, -10.0, consensus.GeoPatch{}, true),
			consensus.NewNode(net.ParseIP("127.0.0.1"), 5000, pub2, consensus.NodeOK, "", 300, "1.0", 0, 1, 30.0, -10.0, consensus.GeoPatch{}, true),
			consensus.NewNode(net.ParseIP("127.0.0.1"), 5000, pub3, consensus.NodeOK, "", 300, "1.0", 0, 1, 30.0, -10.0, consensus.GeoPatch{}, true),
			consensus.NewNode(net.ParseIP("127.0.0.1"), 5000, pub4, consensus.NodeOK, "", 300, "1.0", 0, 1, 30.0, -10.0, consensus.GeoPatch{}, true),
			consensus.NewNode(net.ParseIP("127.0.0.1"), 5000, pub5, consensus.NodeOK, "", 300, "1.0", 0, 1, 30.0, -10.0, consensus.GeoPatch{}, true),
			consensus.NewNode(net.ParseIP("127.0.0.1"), 5000, pub6, consensus.NodeOK, "", 300, "1.0", 0, 1, 30.0, -10.0, consensus.GeoPatch{}, true),
			consensus.NewNode(net.ParseIP("127.0.0.1"), 5000, pub7, consensus.NodeOK, "", 300, "1.0", 0, 1, 30.0, -10.0, consensus.GeoPatch{}, true),
			consensus.NewNode(net.ParseIP("127.0.0.1"), 5000, pub8, consensus.NodeOK, "", 300, "1.0", 0, 1, 30.0, -10.0, consensus.GeoPatch{}, true),
		},
	}
	l := logging.NewLogger(log.New(os.Stdout, "", 0), "test", net.ParseIP("127.0.0.1"), "debug")
	txSrv := rpc.NewTransactionService(chainDB, sharedKeyReader, nodeReader, pr, pub, pv, l)

	//Start transaction server
	lis, _ := net.Listen("tcp", ":5000")
	defer lis.Close()
	grpcServer := grpc.NewServer()
	api.RegisterTransactionServiceServer(grpcServer, txSrv)
	go grpcServer.Serve(lis)

	//Start API
	r := gin.New()
	r.POST("/api/account", CreateAccountHandler(sharedKeyReader, mockNodeDatabase{
		nodes: []consensus.Node{
			consensus.NewNode(net.ParseIP("127.0.0.1"), 5000, pub1, consensus.NodeOK, "", 100, "1.0", 1, 1, 12.5, -100, consensus.GeoPatch{}, true),
			consensus.NewNode(net.ParseIP("127.0.0.1"), 5000, pub2, consensus.NodeOK, "", 500, "1.0", 1, 1, 50.5, -80, consensus.GeoPatch{}, true),
			consensus.NewNode(net.ParseIP("127.0.0.1"), 5000, pub3, consensus.NodeOK, "", 300, "1.0", 1, 1, -30, 20.5, consensus.GeoPatch{}, true),
			consensus.NewNode(net.ParseIP("127.0.0.1"), 5000, pub4, consensus.NodeOK, "", 100, "1.0", 1, 1, 12.5, -100, consensus.GeoPatch{}, true),
			consensus.NewNode(net.ParseIP("127.0.0.1"), 5000, pub5, consensus.NodeOK, "", 500, "1.0", 1, 1, 50.5, -80, consensus.GeoPatch{}, true),
			consensus.NewNode(net.ParseIP("127.0.0.1"), 5000, pub6, consensus.NodeOK, "", 300, "1.0", 1, 1, -30, 20.5, consensus.GeoPatch{}, true),
			consensus.NewNode(net.ParseIP("127.0.0.1"), 5000, pub7, consensus.NodeOK, "", 100, "1.0", 1, 1, 12.5, -100, consensus.GeoPatch{}, true),
			consensus.NewNode(net.ParseIP("127.0.0.1"), 5000, pub8, consensus.NodeOK, "", 500, "1.0", 1, 1, 50.5, -80, consensus.GeoPatch{}, true),
		},
	}, pub, pv, l))

	//Create transactions
	addr := crypto.Hash([]byte("addr"))
	encAddr, _ := pub.Encrypt(addr)
	pubB, _ := pub.Marshal()

	idTx := map[string]interface{}{
		"addr": hex.EncodeToString(crypto.Hash([]byte("abc"))),
		"data": map[string]string{
			"encrypted_address_by_node": hex.EncodeToString(encAddr),
			"encrypted_address_by_id":   hex.EncodeToString(encAddr),
			"encrypted_aes_key":         hex.EncodeToString([]byte("aes_key")),
		},
		"timestamp":  time.Now().Unix(),
		"type":       chain.IDTransactionType,
		"public_key": hex.EncodeToString(pubB),
		"em_shared_keys_proposal": map[string]string{
			"encrypted_private_key": hex.EncodeToString([]byte("encPV")),
			"public_key":            hex.EncodeToString(pubB),
		},
	}

	idTxBytes, _ := json.Marshal(idTx)
	idSig, _ := pv.Sign(idTxBytes)
	idTx["signature"] = hex.EncodeToString(idSig)

	idTxByteWithSig, _ := json.Marshal(idTx)
	emSig, _ := pv.Sign(idTxByteWithSig)
	idTx["em_signature"] = hex.EncodeToString(emSig)

	keychainTx := map[string]interface{}{
		"addr": hex.EncodeToString(crypto.Hash([]byte("abc"))),
		"data": map[string]string{
			"encrypted_address_by_node": hex.EncodeToString(encAddr),
			"encrypted_wallet":          hex.EncodeToString([]byte("wallet")),
		},
		"timestamp":  time.Now().Unix(),
		"type":       chain.KeychainTransactionType,
		"public_key": hex.EncodeToString(pubB),
		"em_shared_keys_proposal": map[string]string{
			"encrypted_private_key": hex.EncodeToString([]byte("encPV")),
			"public_key":            hex.EncodeToString(pubB),
		},
	}

	keychainTxBytes, _ := json.Marshal(keychainTx)
	keychainSig, _ := pv.Sign(keychainTxBytes)
	keychainTx["signature"] = hex.EncodeToString(keychainSig)

	keychainTxByteWithSig, _ := json.Marshal(keychainTx)
	keychainEmSig, _ := pv.Sign(keychainTxByteWithSig)
	keychainTx["em_signature"] = hex.EncodeToString(keychainEmSig)

	idTxBytes, _ = json.Marshal(idTx)

	keychainTxBytes, _ = json.Marshal(keychainTx)

	encryptedID, _ := pub.Encrypt(idTxBytes)
	encryptedKeychain, _ := pub.Encrypt(keychainTxBytes)

	form := accountCreationRequest{
		EncryptedID:       hex.EncodeToString(encryptedID),
		EncryptedKeychain: hex.EncodeToString(encryptedKeychain),
	}
	formB, _ := json.Marshal(form)
	sig, _ := pv.Sign(formB)

	form.Signature = hex.EncodeToString(sig)

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

	idTxHash := crypto.Hash(idTxBytes)
	assert.EqualValues(t, fmt.Sprintf("%x%x", crypto.Hash([]byte("abc")), idTxHash), resTx.IDTransaction.TransactionReceipt)

	idResBytes, _ := json.Marshal(transactionResponse{
		TransactionReceipt: resTx.IDTransaction.TransactionReceipt,
		Timestamp:          resTx.IDTransaction.Timestamp,
	})
	idSigBytes, _ := hex.DecodeString(resTx.IDTransaction.Signature)
	assert.True(t, sharedKeyReader.crossNodeKeys[0].PublicKey().Verify(idResBytes, idSigBytes))

	assert.NotEmpty(t, resTx.KeychainTransaction.TransactionReceipt)
	assert.NotEmpty(t, resTx.KeychainTransaction.Timestamp)
	assert.NotEmpty(t, resTx.KeychainTransaction.Signature)
	assert.Equal(t, time.Now().Unix(), resTx.KeychainTransaction.Timestamp)

	keychainTxHash := crypto.Hash(keychainTxBytes)
	assert.EqualValues(t, fmt.Sprintf("%x%x", crypto.Hash([]byte("abc")), keychainTxHash), resTx.KeychainTransaction.TransactionReceipt)

	keychainResBytes, _ := json.Marshal(transactionResponse{
		TransactionReceipt: resTx.KeychainTransaction.TransactionReceipt,
		Timestamp:          resTx.KeychainTransaction.Timestamp,
	})
	keychainSigBytes, _ := hex.DecodeString(resTx.KeychainTransaction.Signature)
	assert.True(t, sharedKeyReader.crossNodeKeys[0].PublicKey().Verify(keychainResBytes, keychainSigBytes))

	time.Sleep(50 * time.Millisecond)

	assert.Len(t, chainDB.keychains, 5) //because we send to 5 master and because we do not provide sharding yet
	assert.EqualValues(t, crypto.Hash([]byte("abc")), chainDB.keychains[0].Address())
	assert.Len(t, chainDB.ids, 5) //because we send to 5 master and because we do not provide sharding yet
	assert.EqualValues(t, crypto.Hash([]byte("abc")), chainDB.ids[0].Address())

}

type mockPoolRequester struct {
	stores []chain.Transaction
	repo   *mockChainDB
}

func (pr mockPoolRequester) RequestLastTransaction(pool consensus.Pool, txAddr crypto.VersionnedHash, txType chain.TransactionType) (*chain.Transaction, error) {
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

func (pr mockPoolRequester) RequestTransactionTimeLock(pool consensus.Pool, txHash crypto.VersionnedHash, txAddr crypto.VersionnedHash, masterPublicKey crypto.PublicKey) error {
	return nil
}

func (pr mockPoolRequester) RequestTransactionUnlock(pool consensus.Pool, txHash crypto.VersionnedHash, txAddr crypto.VersionnedHash) error {
	return nil
}

func (pr mockPoolRequester) RequestTransactionValidations(pool consensus.Pool, tx chain.Transaction, minValids int, masterValid chain.MasterValidation) ([]chain.Validation, error) {
	pv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

	pubB, _ := pub.Marshal()
	vRaw := map[string]interface{}{
		"status":     chain.ValidationOK,
		"public_key": pubB,
		"timestamp":  time.Now().Unix(),
	}
	vBytes, _ := json.Marshal(vRaw)
	sig, _ := pv.Sign(vBytes)
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

func (r mockChainDB) LastKeychain(txAddr crypto.VersionnedHash) (*chain.Keychain, error) {
	sort.Slice(r.keychains, func(i, j int) bool {
		return r.keychains[i].Timestamp().Unix() > r.keychains[j].Timestamp().Unix()
	})

	if len(r.keychains) > 0 {
		return &r.keychains[0], nil
	}
	return nil, nil
}

func (r mockChainDB) FullKeychain(txAddr crypto.VersionnedHash) (*chain.Keychain, error) {
	sort.Slice(r.keychains, func(i, j int) bool {
		return r.keychains[i].Timestamp().Unix() > r.keychains[j].Timestamp().Unix()
	})

	if len(r.keychains) > 0 {
		return &r.keychains[0], nil
	}
	return nil, nil
}

func (r mockChainDB) KeychainByHash(txHash crypto.VersionnedHash) (*chain.Keychain, error) {
	for _, tx := range r.keychains {
		if bytes.Equal(tx.TransactionHash(), txHash) {
			return &tx, nil
		}
	}
	return nil, nil
}

func (r mockChainDB) IDByHash(txHash crypto.VersionnedHash) (*chain.ID, error) {
	for _, tx := range r.ids {
		if bytes.Equal(tx.TransactionHash(), txHash) {
			return &tx, nil
		}
	}
	return nil, nil
}

func (r mockChainDB) ID(addr crypto.VersionnedHash) (*chain.ID, error) {
	for _, tx := range r.ids {
		if bytes.Equal(tx.Address(), addr) {
			return &tx, nil
		}
	}
	return nil, nil
}

func (r mockChainDB) KOByHash(txHash crypto.VersionnedHash) (*chain.Transaction, error) {
	for _, tx := range r.kos {
		if bytes.Equal(tx.TransactionHash(), txHash) {
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

type mockNodeDatabase struct {
	nodes []consensus.Node
}

func (db mockNodeDatabase) Reachables() (reachables []consensus.Node, err error) {
	for _, n := range db.nodes {
		if n.IsReachable() {
			reachables = append(reachables, n)
		}
	}
	return
}

func (db mockNodeDatabase) Unreachables() (unreachables []consensus.Node, err error) {
	for _, n := range db.nodes {
		if !n.IsReachable() {
			unreachables = append(unreachables, n)
		}
	}
	return
}

func (db mockNodeDatabase) CountReachables() (nb int, err error) {
	for _, n := range db.nodes {
		if n.IsReachable() {
			nb++
		}
	}
	return
}

func (db mockNodeDatabase) FindByPublicKey(publicKey crypto.PublicKey) (found consensus.Node, err error) {
	for _, n := range db.nodes {
		if n.PublicKey().Equals(publicKey) {
			return n, nil
		}
	}
	return
}

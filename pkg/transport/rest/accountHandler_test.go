package rest

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"
	"time"

	api "github.com/uniris/uniris-core/api/protobuf-spec"
	"github.com/uniris/uniris-core/pkg/shared"
	"google.golang.org/grpc"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/uniris/uniris-core/pkg/crypto"
	"github.com/uniris/uniris-core/pkg/transaction"
	"github.com/uniris/uniris-core/pkg/transport/rpc"
)

/*
Scenario: Get account request with an ID hash not a valid hash
	Given an invalid hash (not hexa)
	When I want to request to retrieve an account
	Then I got a 400 (Bad request) response status and an error message
*/
func TestGetAccountWhenInvalidHash(t *testing.T) {
	r := gin.Default()
	apiGroup := r.Group("/api")
	NewAccountHandler(apiGroup, 3545, shared.Service{})

	path := fmt.Sprintf("http://localhost:3000/api/account/abc")
	req, _ := http.NewRequest("GET", path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	resBytes, _ := ioutil.ReadAll(w.Body)
	var errorRes map[string]string
	json.Unmarshal(resBytes, &errorRes)
	assert.Equal(t, errorRes["error"], "id hash: must be hexadecimal")
}

/*
Scenario: Get account request with an invalid idSignature
	Given a hash and an invalid idSignature
	When I want to request to retrieve an account
	Then I got a 400 (Bad request) response status and an error message
*/
func TestGetAccountWhenInvalidSignature(t *testing.T) {
	r := gin.Default()
	apiGroup := r.Group("/api")
	NewAccountHandler(apiGroup, 3545, shared.Service{})

	path1 := fmt.Sprintf("http://localhost:3000/api/account/%s", crypto.HashString("abc"))
	req1, _ := http.NewRequest("GET", path1, nil)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusBadRequest, w1.Code)
	resBytes, _ := ioutil.ReadAll(w1.Body)
	var errorRes map[string]string
	json.Unmarshal(resBytes, &errorRes)
	assert.Equal(t, errorRes["error"], "request signature: signature is empty")

	path2 := fmt.Sprintf("http://localhost:3000/api/account/%s?signature=%s", crypto.HashString("abc"), "idSig")
	req2, _ := http.NewRequest("GET", path2, nil)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusBadRequest, w2.Code)
	resBytes2, _ := ioutil.ReadAll(w2.Body)
	var errorRes2 map[string]string
	json.Unmarshal(resBytes2, &errorRes2)
	assert.Equal(t, errorRes2["error"], "request signature: signature is not in hexadecimal format")

	path3 := fmt.Sprintf("http://localhost:3000/api/account/%s?signature=%s", crypto.HashString("abc"), hex.EncodeToString([]byte("idSig")))
	req3, _ := http.NewRequest("GET", path3, nil)
	w3 := httptest.NewRecorder()
	r.ServeHTTP(w3, req3)
	assert.Equal(t, http.StatusBadRequest, w3.Code)
	resBytes3, _ := ioutil.ReadAll(w3.Body)
	var errorRes3 map[string]string
	json.Unmarshal(resBytes3, &errorRes3)
	assert.Equal(t, errorRes3["error"], "request signature: signature is not valid")
}

/*
Scenario: Get account request with an ID not existing
	Given an ID hash and a valid idSignature related to no real ID transaction
	When I want to request to retrieve an account
	Then I got a 404 (Not found) response status and an error message
*/
func TestGetAccountWhenIDNotExist(t *testing.T) {

	pub, pv := crypto.GenerateKeys()

	encPv, _ := crypto.Encrypt(pv, pub)
	minerKP, _ := shared.NewMinerKeyPair(pub, pv)
	emKP, _ := shared.NewEmitterKeyPair(encPv, pub)

	txRepo := &mockTxRepository{}
	lockRepo := &mockLockRepository{}
	sharedRepo := &mockSharedRepo{}
	sharedRepo.emKeys = []shared.EmitterKeyPair{emKP}
	sharedRepo.minerKeys = minerKP
	poolR := &mockPoolRequester{}

	sharedSrv := shared.NewService(sharedRepo)
	poolingSrv := transaction.NewPoolFindingService(rpc.NewPoolRetriever(sharedSrv))
	miningSrv := transaction.NewMiningService(poolR, poolingSrv, sharedSrv, "127.0.0.1", pub, pv)

	storageSrv := transaction.NewStorageService(txRepo, miningSrv)
	lockSrv := transaction.NewLockService(lockRepo)

	txSrv := rpc.NewTransactionServer(storageSrv, lockSrv, miningSrv, sharedSrv)

	lisTx, _ := net.Listen("tcp", ":3545")
	defer lisTx.Close()
	grpcServer := grpc.NewServer()
	api.RegisterTransactionServiceServer(grpcServer, txSrv)
	go grpcServer.Serve(lisTx)

	intSrv := rpc.NewInternalServer(poolingSrv, miningSrv, sharedSrv)
	lisInt, _ := net.Listen("tcp", ":1717")
	defer lisInt.Close()
	grpcServerInt := grpc.NewServer()
	api.RegisterInternalServiceServer(grpcServerInt, intSrv)
	go grpcServerInt.Serve(lisInt)

	r := gin.Default()
	apiGroup := r.Group("/api")
	NewAccountHandler(apiGroup, 1717, sharedSrv)

	idHash := crypto.HashString("abc")
	encIDHash, _ := crypto.Encrypt(idHash, pub)
	sig, _ := crypto.Sign(encIDHash, pv)

	path := fmt.Sprintf("http://localhost:3000/api/account/%s?signature=%s", encIDHash, sig)
	req, _ := http.NewRequest("GET", path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	resBytes, _ := ioutil.ReadAll(w.Body)
	var errorRes map[string]string
	json.Unmarshal(resBytes, &errorRes)
	assert.Equal(t, errorRes["error"], "ID does not exist")

}

/*
Scenario: Get account request with a Keychain not existing
	Given an ID hash and a valid idSignature related to no real Keychain transaction
	When I want to request to retrieve an account
	Then I got a 404 (Not found) response status and an error message
*/
func TestGetAccountWhenKeychainNotExist(t *testing.T) {
	pub, pv := crypto.GenerateKeys()

	txRepo := &mockTxRepository{}
	lockRepo := &mockLockRepository{}
	poolR := &mockPoolRequester{}

	sharedRepo := &mockSharedRepo{}
	encPv, _ := crypto.Encrypt(pv, pub)
	minerKP, _ := shared.NewMinerKeyPair(pub, pv)
	emKP, _ := shared.NewEmitterKeyPair(encPv, pub)
	sharedRepo.emKeys = []shared.EmitterKeyPair{emKP}
	sharedRepo.minerKeys = minerKP

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
	NewAccountHandler(apiGroup, 1717, sharedSrv)

	//Create transactions
	encAddr, _ := crypto.Encrypt(hex.EncodeToString([]byte("addr")), pub)

	idData := map[string]string{
		"encrypted_address_by_robot": encAddr,
		"encrypted_address_by_id":    encAddr,
		"encrypted_aes_key":          hex.EncodeToString([]byte("aes_key")),
	}
	idHash := crypto.HashString("abc")
	txIDRaw, _ := json.Marshal(txRaw{
		Address: crypto.HashString("idHash"),
		Data:    idData,
		Proposal: txProp{
			SharedEmitterKeys: txSharedKeys{
				EncryptedPrivateKey: hex.EncodeToString([]byte("encPV")),
				PublicKey:           pub,
			},
		},
		Timestamp: time.Now().Unix(),
		Type:      int(transaction.IDType),
		PublicKey: pub,
	})
	idSig, _ := crypto.Sign(string(txIDRaw), pv)
	txIDSigned, _ := json.Marshal(txSigned{
		Address: crypto.HashString("idHash"),
		Data:    idData,
		Proposal: txProp{
			SharedEmitterKeys: txSharedKeys{
				EncryptedPrivateKey: hex.EncodeToString([]byte("encPV")),
				PublicKey:           pub,
			},
		},
		Timestamp:        time.Now().Unix(),
		Type:             int(transaction.IDType),
		PublicKey:        pub,
		EmitterSignature: idSig,
		Signature:        idSig,
	})

	kp, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("encPV")), pub)
	prop, _ := transaction.NewProposal(kp)

	idTx, _ := transaction.New(idHash, transaction.IDType, idData, time.Now(), pub, idSig, idSig, prop, crypto.HashBytes(txIDSigned))
	id, _ := transaction.NewID(idTx)
	txRepo.StoreID(id)

	encIDHash, _ := crypto.Encrypt(idHash, pub)
	sig, _ := crypto.Sign(encIDHash, pv)

	path := fmt.Sprintf("http://localhost:3000/api/account/%s?signature=%s", encIDHash, sig)
	req, _ := http.NewRequest("GET", path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)

	resBytes, _ := ioutil.ReadAll(w.Body)
	var errorRes map[string]string
	json.Unmarshal(resBytes, &errorRes)
	assert.Equal(t, errorRes["error"], "Keychain does not exist")

}

/*
Scenario: Get an account after its creation
	Given an account created (keychain and ID transaction)
	When I want to retrieve it
	Then I can get encrypted wallet and encrypted aes key
*/
func TestGetAccount(t *testing.T) {
	pub, pv := crypto.GenerateKeys()

	txRepo := &mockTxRepository{}
	lockRepo := &mockLockRepository{}
	poolR := &mockPoolRequester{}

	sharedRepo := &mockSharedRepo{}
	encPv, _ := crypto.Encrypt(pv, pub)
	minerKP, _ := shared.NewMinerKeyPair(pub, pv)
	emKP, _ := shared.NewEmitterKeyPair(encPv, pub)
	sharedRepo.emKeys = []shared.EmitterKeyPair{emKP}
	sharedRepo.minerKeys = minerKP

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
	NewAccountHandler(apiGroup, 1717, sharedSrv)

	//Create transactions
	addr := crypto.HashString("addr")
	encAddr, _ := crypto.Encrypt(addr, pub)

	idData := map[string]string{
		"encrypted_address_by_robot": encAddr,
		"encrypted_address_by_id":    encAddr,
		"encrypted_aes_key":          hex.EncodeToString([]byte("aes_key")),
	}
	idHash := crypto.HashString("abc")
	txIDRaw, _ := json.Marshal(txRaw{
		Address: idHash,
		Data:    idData,
		Proposal: txProp{
			SharedEmitterKeys: txSharedKeys{
				EncryptedPrivateKey: hex.EncodeToString([]byte("encPV")),
				PublicKey:           pub,
			},
		},
		Timestamp: time.Now().Unix(),
		Type:      int(transaction.IDType),
		PublicKey: pub,
	})
	idSig, _ := crypto.Sign(string(txIDRaw), pv)
	txIDSigned, _ := json.Marshal(txSigned{
		Address: idHash,
		Data:    idData,
		Proposal: txProp{
			SharedEmitterKeys: txSharedKeys{
				EncryptedPrivateKey: hex.EncodeToString([]byte("encPV")),
				PublicKey:           pub,
			},
		},
		Timestamp:        time.Now().Unix(),
		Type:             int(transaction.IDType),
		PublicKey:        pub,
		EmitterSignature: idSig,
		Signature:        idSig,
	})

	kp, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("encPV")), pub)
	prop, _ := transaction.NewProposal(kp)

	idTx, _ := transaction.New(idHash, transaction.IDType, idData, time.Now(), pub, idSig, idSig, prop, crypto.HashBytes(txIDSigned))
	id, _ := transaction.NewID(idTx)
	txRepo.StoreID(id)

	keychainData := map[string]string{
		"encrypted_address": encAddr,
		"encrypted_wallet":  hex.EncodeToString([]byte("wallet")),
	}

	txKeychainRaw, _ := json.Marshal(txRaw{
		Address: addr,
		Data:    keychainData,
		Proposal: txProp{
			SharedEmitterKeys: txSharedKeys{
				EncryptedPrivateKey: hex.EncodeToString([]byte("encPV")),
				PublicKey:           pub,
			},
		},
		Timestamp: time.Now().Unix(),
		Type:      int(transaction.KeychainType),
		PublicKey: pub,
	})
	keychainSig, _ := crypto.Sign(string(txKeychainRaw), pv)
	txKeychainSigned, _ := json.Marshal(txSigned{
		Address: addr,
		Data:    keychainData,
		Proposal: txProp{
			SharedEmitterKeys: txSharedKeys{
				EncryptedPrivateKey: hex.EncodeToString([]byte("encPV")),
				PublicKey:           pub,
			},
		},
		Timestamp:        time.Now().Unix(),
		Type:             int(transaction.KeychainType),
		PublicKey:        pub,
		EmitterSignature: keychainSig,
		Signature:        keychainSig,
	})

	keychainTx, _ := transaction.New(addr, transaction.KeychainType, keychainData, time.Now(), pub, keychainSig, keychainSig, prop, crypto.HashBytes(txKeychainSigned))
	keychain, _ := transaction.NewKeychain(keychainTx)
	txRepo.StoreKeychain(keychain)

	encIDHash, _ := crypto.Encrypt(idHash, pub)
	sig, _ := crypto.Sign(encIDHash, pv)

	path := fmt.Sprintf("http://localhost:3000/api/account/%s?signature=%s", encIDHash, sig)

	req, _ := http.NewRequest("GET", path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	resBytes, _ := ioutil.ReadAll(w.Body)
	var accountRes api.GetAccountResponse
	json.Unmarshal(resBytes, &accountRes)
	assert.Equal(t, hex.EncodeToString([]byte("aes_key")), accountRes.EncryptedAesKey)
	assert.Equal(t, hex.EncodeToString([]byte("wallet")), accountRes.EncryptedWallet)

}

/*
Scenario: Create account request with an invalid encrypted ID
	Given an invalid encrypted id (not hexa)
	When I want to request to create an account
	Then I got a 400 (Bad request) response status and an error message
*/
func TestCreationAccountWhenInvalidID(t *testing.T) {
	r := gin.Default()
	apiGroup := r.Group("/api")
	NewAccountHandler(apiGroup, 3545, shared.Service{})

	form, _ := json.Marshal(map[string]string{
		"encrypted_id":       "abc",
		"encrypted_keychain": "abc",
		"signature":          "abc",
	})

	path := "http://localhost:3000/api/account"
	req, _ := http.NewRequest("POST", path, bytes.NewBuffer(form))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	resBytes, _ := ioutil.ReadAll(w.Body)
	var errorRes map[string]string
	json.Unmarshal(resBytes, &errorRes)
	assert.Equal(t, errorRes["error"], "encrypted id: must be hexadecimal")

}

/*
Scenario: Create account request with an invalid encrypted Keychain
	Given an invalid encrypted Keychain (not hexa)
	When I want to request to create an account
	Then I got a 400 (Bad request) response status and an error message
*/
func TestCreationAccountWhenKeychainInvalid(t *testing.T) {
	r := gin.Default()
	apiGroup := r.Group("/api")
	NewAccountHandler(apiGroup, 3545, shared.Service{})

	form, _ := json.Marshal(map[string]string{
		"encrypted_id":       hex.EncodeToString([]byte("id")),
		"encrypted_keychain": "abc",
		"signature":          "abc",
	})
	log.Print(string(form))

	path := "http://localhost:3000/api/account"
	req, _ := http.NewRequest("POST", path, bytes.NewBuffer(form))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	resBytes, _ := ioutil.ReadAll(w.Body)
	var errorRes map[string]string
	json.Unmarshal(resBytes, &errorRes)
	assert.Equal(t, errorRes["error"], "encrypted keychain: must be hexadecimal")

}

/*
Scenario: Create account request with an invalid signature
	Given an invalid signature (not hexa and not valid)
	When I want to request to create an account
	Then I got a 400 (Bad request) response status and an error message
*/
func TestCreationAccountWhenSignatureInvalid(t *testing.T) {

	pub, pv := crypto.GenerateKeys()

	sharedRepo := &mockSharedRepo{}
	encPv, _ := crypto.Encrypt(pv, pub)
	minerKP, _ := shared.NewMinerKeyPair(pub, pv)
	emKP, _ := shared.NewEmitterKeyPair(encPv, pub)
	sharedRepo.emKeys = []shared.EmitterKeyPair{emKP}
	sharedRepo.minerKeys = minerKP

	sharedSrv := shared.NewService(sharedRepo)
	r := gin.Default()
	apiGroup := r.Group("/api")
	NewAccountHandler(apiGroup, 3545, sharedSrv)

	form, _ := json.Marshal(map[string]string{
		"encrypted_id":       hex.EncodeToString([]byte("id")),
		"encrypted_keychain": hex.EncodeToString([]byte("keychain")),
		"signature":          "abc",
	})

	path := "http://localhost:3000/api/account"
	req1, _ := http.NewRequest("POST", path, bytes.NewBuffer(form))
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)

	assert.Equal(t, http.StatusBadRequest, w1.Code)
	resBytes, _ := ioutil.ReadAll(w1.Body)

	var errorRes map[string]string
	json.Unmarshal(resBytes, &errorRes)
	assert.Equal(t, errorRes["error"], "signature request: signature is not in hexadecimal format")

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
	var errorRes2 map[string]string
	json.Unmarshal(resBytes2, &errorRes2)
	assert.Equal(t, errorRes2["error"], "signature request: signature is not valid")

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
	var errorRes3 map[string]string
	json.Unmarshal(resBytes3, &errorRes3)
	assert.Equal(t, errorRes3["error"], "signature request: signature is not valid")

}

/*
Scenario: Create an account including ID and keychain transaction
	Given a valid ID and keychain transaction
	When I want to create an account
	Then two transaction are created (ID/Keychain) and the data is stored
*/
func TestCreateAccount(t *testing.T) {
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
	NewAccountHandler(apiGroup, 1717, sharedSrv)

	//Create transactions
	addr := crypto.HashString("addr")
	encAddr, _ := crypto.Encrypt(addr, pub)

	idTx := map[string]interface{}{
		"address": crypto.HashString("abc"),
		"data": map[string]string{
			"encrypted_address_by_robot": encAddr,
			"encrypted_address_by_id":    encAddr,
			"encrypted_aes_key":          hex.EncodeToString([]byte("aes_key")),
		},
		"timestamp":  time.Now().Unix(),
		"type":       int(transaction.IDType),
		"public_key": pub,
		"proposal": map[string]interface{}{
			"shared_emitter_keys": map[string]string{
				"encrypted_private_key": hex.EncodeToString([]byte("encPV")),
				"public_key":            pub,
			},
		},
	}

	idTxBytes, _ := json.Marshal(idTx)
	idSig, _ := crypto.Sign(string(idTxBytes), pv)
	idTx["signature"] = idSig
	idTx["em_signature"] = idSig

	keychainTx := map[string]interface{}{
		"address": addr,
		"data": map[string]string{
			"encrypted_address": encAddr,
			"encrypted_wallet":  hex.EncodeToString([]byte("wallet")),
		},
		"timestamp":  time.Now().Unix(),
		"type":       int(transaction.KeychainType),
		"public_key": pub,
		"proposal": map[string]interface{}{
			"shared_emitter_keys": map[string]string{
				"encrypted_private_key": hex.EncodeToString([]byte("encPV")),
				"public_key":            pub,
			},
		},
	}

	keychainTxBytes, _ := json.Marshal(keychainTx)
	keychainSig, _ := crypto.Sign(string(keychainTxBytes), pv)
	keychainTx["signature"] = keychainSig
	keychainTx["em_signature"] = keychainSig

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

	path := "http://localhost:3000/api/account"
	req, _ := http.NewRequest("POST", path, bytes.NewBuffer(formB))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	resBytes, _ := ioutil.ReadAll(w.Body)

	var resTx map[string]interface{}
	json.Unmarshal(resBytes, &resTx)

	idTxRes := resTx["id_transaction"].(map[string]interface{})
	assert.NotEmpty(t, idTxRes["transaction_hash"])
	assert.NotEmpty(t, idTxRes["timestamp"])
	assert.NotEmpty(t, idTxRes["signature"])

	keychainTxRes := resTx["keychain_transaction"].(map[string]interface{})
	assert.NotEmpty(t, keychainTxRes["transaction_hash"])
	assert.NotEmpty(t, keychainTxRes["timestamp"])
	assert.NotEmpty(t, keychainTxRes["signature"])

	time.Sleep(1 * time.Second)

	assert.Len(t, txRepo.keychains, 1)
	assert.Equal(t, addr, txRepo.keychains[0].Address())
	assert.Len(t, txRepo.ids, 1)
	assert.Equal(t, crypto.HashString("abc"), txRepo.ids[0].Address())

}

type mockSharedRepo struct {
	emKeys    shared.EmitterKeys
	minerKeys shared.MinerKeyPair
}

func (r mockSharedRepo) ListSharedEmitterKeyPairs() (shared.EmitterKeys, error) {
	return r.emKeys, nil
}
func (r *mockSharedRepo) StoreSharedEmitterKeyPair(kp shared.EmitterKeyPair) error {
	r.emKeys = append(r.emKeys, kp)
	return nil
}

func (r *mockSharedRepo) GetLastSharedMinersKeyPair() (shared.MinerKeyPair, error) {
	return r.minerKeys, nil
}

type mockPoolRequester struct {
	stores []transaction.Transaction
	repo   *mockTxRepository
}

func (pr mockPoolRequester) RequestTransactionLock(pool transaction.Pool, txLock transaction.Lock) error {
	return nil
}

func (pr mockPoolRequester) RequestTransactionUnlock(pool transaction.Pool, txLock transaction.Lock) error {
	return nil
}

func (pr mockPoolRequester) RequestTransactionValidations(pool transaction.Pool, tx transaction.Transaction, masterValid transaction.MasterValidation, validChan chan<- transaction.MinerValidation) {
	pub, pv := crypto.GenerateKeys()

	v := transaction.MinerValidation{}
	vRaw := map[string]interface{}{
		"status":     transaction.ValidationOK,
		"public_key": pub,
		"timestamp":  time.Now().Unix(),
	}
	vBytes, _ := json.Marshal(vRaw)
	sig, _ := crypto.Sign(string(vBytes), pv)
	v, _ = transaction.NewMinerValidation(transaction.ValidationOK, time.Now(), pub, sig)

	validChan <- v
}

func (pr *mockPoolRequester) RequestTransactionStorage(pool transaction.Pool, tx transaction.Transaction, ackChan chan<- bool) {
	pr.stores = append(pr.stores, tx)
	if tx.Type() == transaction.KeychainType {
		k, _ := transaction.NewKeychain(tx)
		pr.repo.keychains = append(pr.repo.keychains, k)
	}
	if tx.Type() == transaction.IDType {
		id, _ := transaction.NewID(tx)
		pr.repo.ids = append(pr.repo.ids, id)
	}
	ackChan <- true
}

type mockTxRepository struct {
	pendings  []transaction.Transaction
	kos       []transaction.Transaction
	keychains []transaction.Keychain
	ids       []transaction.ID
}

func (r mockTxRepository) FindPendingTransaction(txHash string) (*transaction.Transaction, error) {
	for _, tx := range r.pendings {
		if tx.TransactionHash() == txHash {
			return &tx, nil
		}
	}
	return nil, nil
}

func (r mockTxRepository) GetKeychain(txAddr string) (*transaction.Keychain, error) {
	sort.Slice(r.keychains, func(i, j int) bool {
		return r.keychains[i].Timestamp().Unix() > r.keychains[j].Timestamp().Unix()
	})

	if len(r.keychains) > 0 {
		return &r.keychains[0], nil
	}
	return nil, nil
}

func (r mockTxRepository) FindKeychainByHash(txHash string) (*transaction.Keychain, error) {
	for _, tx := range r.keychains {
		if tx.TransactionHash() == txHash {
			return &tx, nil
		}
	}
	return nil, nil
}

func (r mockTxRepository) FindLastKeychain(addr string) (*transaction.Keychain, error) {

	sort.Slice(r.keychains, func(i, j int) bool {
		return r.keychains[i].Timestamp().Unix() > r.keychains[j].Timestamp().Unix()
	})

	for _, tx := range r.keychains {
		if tx.Address() == addr {
			return &tx, nil
		}
	}
	return nil, nil
}

func (r mockTxRepository) FindIDByHash(txHash string) (*transaction.ID, error) {
	for _, tx := range r.ids {
		if tx.TransactionHash() == txHash {
			return &tx, nil
		}
	}
	return nil, nil
}

func (r mockTxRepository) FindIDByAddress(addr string) (*transaction.ID, error) {
	for _, tx := range r.ids {
		if tx.Address() == addr {
			return &tx, nil
		}
	}
	return nil, nil
}

func (r mockTxRepository) FindKOTransaction(txHash string) (*transaction.Transaction, error) {
	for _, tx := range r.kos {
		if tx.TransactionHash() == txHash {
			return &tx, nil
		}
	}
	return nil, nil
}

func (r *mockTxRepository) StoreKeychain(kc transaction.Keychain) error {
	r.keychains = append(r.keychains, kc)
	return nil
}

func (r *mockTxRepository) StoreID(id transaction.ID) error {
	r.ids = append(r.ids, id)
	return nil
}

func (r *mockTxRepository) StoreKO(tx transaction.Transaction) error {
	r.kos = append(r.kos, tx)
	return nil
}

type mockLockRepository struct {
	locks []transaction.Lock
}

func (r *mockLockRepository) StoreLock(l transaction.Lock) error {
	r.locks = append(r.locks, l)
	return nil
}

func (r *mockLockRepository) RemoveLock(l transaction.Lock) error {
	pos := r.findLockPosition(l)
	if pos > -1 {
		r.locks = append(r.locks[:pos], r.locks[pos+1:]...)
	}
	return nil
}

func (r mockLockRepository) ContainsLock(l transaction.Lock) (bool, error) {
	return r.findLockPosition(l) > -1, nil
}

func (r mockLockRepository) findLockPosition(l transaction.Lock) int {
	for i, lock := range r.locks {
		if lock.TransactionHash() == l.TransactionHash() && l.MasterRobotKey() == lock.MasterRobotKey() {
			return i
		}
	}
	return -1
}

type txSigned struct {
	Address          string            `json:"address"`
	Data             map[string]string `json:"data"`
	Timestamp        int64             `json:"timestamp"`
	Type             int               `json:"type"`
	PublicKey        string            `json:"public_key"`
	Proposal         txProp            `json:"proposal"`
	Signature        string            `json:"idSignature"`
	EmitterSignature string            `json:"em_idSignature"`
}

type txRaw struct {
	Address   string            `json:"address"`
	Data      map[string]string `json:"data"`
	Timestamp int64             `json:"timestamp"`
	Type      int               `json:"type"`
	PublicKey string            `json:"public_key"`
	Proposal  txProp            `json:"proposal"`
}

type txProp struct {
	SharedEmitterKeys txSharedKeys `json:"shared_emitter_keys"`
}

type txSharedKeys struct {
	EncryptedPrivateKey string `json:"encrypted_private_key"`
	PublicKey           string `json:"public_key"`
}

func formatLeadMiningRequest(tx txSigned, txHash string, minValidations int) *api.LeadMiningRequest {
	return &api.LeadMiningRequest{
		MinimumValidations: int32(minValidations),
		Timestamp:          time.Now().Unix(),
		Transaction: &api.Transaction{
			Address:          tx.Address,
			Data:             tx.Data,
			Type:             api.TransactionType(tx.Type),
			Timestamp:        tx.Timestamp,
			PublicKey:        tx.PublicKey,
			Signature:        tx.Signature,
			EmitterSignature: tx.EmitterSignature,
			Proposal: &api.TransactionProposal{
				SharedEmitterKeys: &api.SharedKeyPair{
					EncryptedPrivateKey: tx.Proposal.SharedEmitterKeys.EncryptedPrivateKey,
					PublicKey:           tx.Proposal.SharedEmitterKeys.PublicKey,
				},
			},
			TransactionHash: txHash,
		},
	}
}

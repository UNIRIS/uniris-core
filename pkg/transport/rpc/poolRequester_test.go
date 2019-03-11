package rpc

import (
	"encoding/hex"
	"encoding/json"
	"net"
	"testing"
	"time"

	"github.com/uniris/uniris-core/pkg/chain"

	"github.com/uniris/uniris-core/pkg/consensus"

	"github.com/uniris/uniris-core/pkg/shared"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

	api "github.com/uniris/uniris-core/api/protobuf-spec"
	"github.com/uniris/uniris-core/pkg/crypto"
)

/*
Scenario: Request transction lock on a pool
	Given a transaction to lock
	When I request to lock it
	Then the lock is stored in the database
*/
func TestRequestTransactionLock(t *testing.T) {

	pub, pv := crypto.GenerateKeys()

	techDB := &mockTechDB{}
	nodeKey, _ := shared.NewKeyPair(pub, pv)
	techDB.nodeKeys = append(techDB.nodeKeys, nodeKey)

	chainDB := &mockChainDB{}
	pr := NewPoolRequester(techDB)
	txSrv := NewTransactionService(chainDB, techDB, pr, pub, pv)

	lis, _ := net.Listen("tcp", ":5000")
	defer lis.Close()
	grpcServer := grpc.NewServer()
	api.RegisterTransactionServiceServer(grpcServer, txSrv)
	go grpcServer.Serve(lis)

	pool, _ := consensus.FindStoragePool("addr")
	assert.Nil(t, pr.RequestTransactionTimeLock(pool, crypto.HashString("tx"), crypto.HashString("addr"), pub))

	assert.True(t, chain.ContainsTimeLock(crypto.HashString("tx"), crypto.HashString("addr")))
}

/*
Scenario: Request transaction validation confirmation
	Given a transaction to validate
	When I request to confirm the validation
	Then I get a validation
*/
func TestRequestConfirmValidation(t *testing.T) {

	pub, pv := crypto.GenerateKeys()

	kp, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("pvkey")), pub)

	techDB := &mockTechDB{}
	techDB.emKeys = append(techDB.emKeys, kp)
	nodeKey, _ := shared.NewKeyPair(pub, pv)
	techDB.nodeKeys = append(techDB.nodeKeys, nodeKey)

	pr := NewPoolRequester(techDB)

	miningSrv := NewTransactionService(nil, techDB, pr, pub, pv)

	lis, err := net.Listen("tcp", ":5000")
	defer lis.Close()
	grpcServer := grpc.NewServer()
	api.RegisterTransactionServiceServer(grpcServer, miningSrv)
	go grpcServer.Serve(lis)

	data := map[string]string{
		"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
		"encrypted_wallet":          hex.EncodeToString([]byte("wallet")),
	}
	prop := kp
	txRaw := map[string]interface{}{
		"addr":                    crypto.HashString("addr"),
		"data":                    data,
		"type":                    chain.KeychainTransactionType,
		"timestamp":               time.Now().Unix(),
		"public_key":              pub,
		"em_shared_keys_proposal": prop,
	}
	txBytes, _ := json.Marshal(txRaw)
	sig, _ := crypto.Sign(string(txBytes), pv)
	txRaw["signature"] = sig
	txByteWithSig, _ := json.Marshal(txRaw)
	emSig, _ := crypto.Sign(string(txByteWithSig), pv)
	txRaw["em_signature"] = emSig
	txBytes, _ = json.Marshal(txRaw)

	tx, _ := chain.NewTransaction(crypto.HashString("addr"), chain.KeychainTransactionType, data, time.Now(), pub, prop, sig, emSig, crypto.HashBytes(txBytes))
	vBytes, _ := json.Marshal(map[string]interface{}{
		"status":     chain.ValidationOK,
		"public_key": pub,
		"timestamp":  time.Now().Unix(),
	})
	vSig, _ := crypto.Sign(string(vBytes), pv)
	v, _ := chain.NewValidation(chain.ValidationOK, time.Now(), pub, vSig)
	wHeaders := []chain.NodeHeader{chain.NewNodeHeader("pub", false, false, 0, true)}
	vHeaders := []chain.NodeHeader{chain.NewNodeHeader("pub", false, false, 0, true)}
	sHeaders := []chain.NodeHeader{chain.NewNodeHeader("pub", false, false, 0, true)}
	mv, _ := chain.NewMasterValidation([]string{}, pub, v, wHeaders, vHeaders, sHeaders)

	pool, _ := consensus.FindValidationPool(tx)
	valids, err := pr.RequestTransactionValidations(pool, tx, 1, mv)
	assert.Nil(t, err)

	assert.Len(t, valids, 1)
	assert.Equal(t, pub, valids[0].PublicKey())
	assert.Equal(t, chain.ValidationOK, valids[0].Status())
	ok, err := valids[0].IsValid()
	assert.Nil(t, err)
	assert.True(t, ok)
}

/*
Scenario: Request transaction store
	Given a transaction to store
	When I request to store the validation
	Then the transaction is stored
*/
func TestRequestStorage(t *testing.T) {

	pub, pv := crypto.GenerateKeys()
	kp, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("pvkey")), pub)

	chainDB := &mockChainDB{}
	techDB := &mockTechDB{}
	nodeKey, _ := shared.NewKeyPair(pub, pv)
	techDB.nodeKeys = append(techDB.nodeKeys, nodeKey)

	techDB.emKeys = append(techDB.emKeys, kp)

	pr := NewPoolRequester(techDB)

	txSrv := NewTransactionService(chainDB, techDB, pr, pub, pv)

	lis, _ := net.Listen("tcp", ":5000")
	defer lis.Close()
	grpcServer := grpc.NewServer()
	api.RegisterTransactionServiceServer(grpcServer, txSrv)
	go grpcServer.Serve(lis)

	data := map[string]string{
		"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
		"encrypted_wallet":          hex.EncodeToString([]byte("wallet")),
	}
	prop := kp
	txRaw := map[string]interface{}{
		"addr":                    crypto.HashString("addr"),
		"data":                    data,
		"type":                    chain.KeychainTransactionType,
		"timestamp":               time.Now().Unix(),
		"public_key":              pub,
		"em_shared_keys_proposal": prop,
	}
	txBytes, _ := json.Marshal(txRaw)
	sig, _ := crypto.Sign(string(txBytes), pv)
	txRaw["signature"] = sig
	txByteWithSig, _ := json.Marshal(txRaw)
	emSig, _ := crypto.Sign(string(txByteWithSig), pv)
	txRaw["em_signature"] = emSig
	txBytes, _ = json.Marshal(txRaw)

	tx, _ := chain.NewTransaction(crypto.HashString("addr"), chain.KeychainTransactionType, data, time.Now(), pub, prop, sig, emSig, crypto.HashBytes(txBytes))
	vBytes, _ := json.Marshal(map[string]interface{}{
		"status":     chain.ValidationOK,
		"public_key": pub,
		"timestamp":  time.Now().Unix(),
	})
	vSig, _ := crypto.Sign(string(vBytes), pv)
	v, _ := chain.NewValidation(chain.ValidationOK, time.Now(), pub, vSig)
	wHeaders := []chain.NodeHeader{chain.NewNodeHeader("pub", false, false, 0, true)}
	vHeaders := []chain.NodeHeader{chain.NewNodeHeader("pub", false, false, 0, true)}
	sHeaders := []chain.NodeHeader{chain.NewNodeHeader("pub", false, false, 0, true)}
	mv, _ := chain.NewMasterValidation([]string{}, pub, v, wHeaders, vHeaders, sHeaders)

	pool, _ := consensus.FindStoragePool("addr")
	tx.Mined(mv, []chain.Validation{v})
	assert.Nil(t, pr.RequestTransactionStorage(pool, 1, tx))

	assert.Len(t, chainDB.keychains, 1)
	assert.Equal(t, crypto.HashBytes(txBytes), chainDB.keychains[0].TransactionHash())
}

/*
Scenario: Send request to get last transaction
	Given a keychain transaction stored
	When I want to request a node to get the last transaction from the address
	Then I get the last transaction
*/
func TestSendGetLastTransaction(t *testing.T) {

	pub, pv := crypto.GenerateKeys()

	chainDB := &mockChainDB{}
	techDB := &mockTechDB{}
	nodeKey, _ := shared.NewKeyPair(pub, pv)
	techDB.nodeKeys = append(techDB.nodeKeys, nodeKey)

	pr := NewPoolRequester(techDB)

	txSrv := NewTransactionService(chainDB, techDB, pr, pub, pv)

	lis, _ := net.Listen("tcp", ":5000")
	defer lis.Close()
	grpcServer := grpc.NewServer()
	api.RegisterTransactionServiceServer(grpcServer, txSrv)
	go grpcServer.Serve(lis)

	data := map[string]string{
		"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
		"encrypted_wallet":          hex.EncodeToString([]byte("wallet")),
	}

	prop, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("encPV")), pub)
	txRaw := map[string]interface{}{
		"address":                 crypto.HashString("addr"),
		"data":                    data,
		"timestamp":               time.Now().Unix(),
		"type":                    chain.KeychainTransactionType,
		"public_key":              pub,
		"em_shared_keys_proposal": prop,
	}
	txBytes, _ := json.Marshal(txRaw)
	sig, _ := crypto.Sign(string(txBytes), pv)
	txRaw["signature"] = sig
	txByteWithSig, _ := json.Marshal(txRaw)
	emSig, _ := crypto.Sign(string(txByteWithSig), pv)
	txRaw["em_signature"] = emSig
	txBytes, _ = json.Marshal(txRaw)

	tx, _ := chain.NewTransaction(crypto.HashString("addr"), chain.KeychainTransactionType, data, time.Now(), pub, prop, sig, sig, crypto.HashBytes(txBytes))
	keychain, _ := chain.NewKeychain(tx)
	chainDB.keychains = append(chainDB.keychains, keychain)

	pool, _ := consensus.FindStoragePool("address")

	txRes, err := pr.RequestLastTransaction(pool, crypto.HashString("addr"), chain.KeychainTransactionType)
	assert.Nil(t, err)
	assert.Equal(t, chain.KeychainTransactionType, txRes.TransactionType())
	assert.Equal(t, crypto.HashBytes(txBytes), txRes.TransactionHash())
}

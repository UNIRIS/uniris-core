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
	minerKey, _ := shared.NewMinerKeyPair(pub, pv)
	techDB.minerKeys = append(techDB.minerKeys, minerKey)

	lockDB := &mockLockDb{}
	lockSrv := NewLockServer(lockDB, techDB)
	pr := NewPoolRequester(techDB)

	lis, _ := net.Listen("tcp", ":5000")
	defer lis.Close()
	grpcServer := grpc.NewServer()
	api.RegisterLockServiceServer(grpcServer, lockSrv)
	go grpcServer.Serve(lis)

	pool, _ := consensus.FindStoragePool("addr")
	assert.Nil(t, pr.RequestTransactionLock(pool, crypto.HashString("tx"), crypto.HashString("addr"), pub))

	assert.Len(t, lockDB.locks, 1)
	assert.Equal(t, crypto.HashString("addr"), lockDB.locks[0]["transaction_address"])
}

/*
Scenario: Request transction lock on a pool
	Given a transaction already locked
	When I request to unlock it
	Then the lock is removed from the database
*/
func TestRequestTransactionUnlock(t *testing.T) {

	pub, pv := crypto.GenerateKeys()

	techDB := &mockTechDB{}
	minerKey, _ := shared.NewMinerKeyPair(pub, pv)
	techDB.minerKeys = append(techDB.minerKeys, minerKey)

	lockDB := &mockLockDb{}
	lockSrv := NewLockServer(lockDB, techDB)
	lockDB.locks = append(lockDB.locks, map[string]string{
		"transaction_hash":    crypto.HashString("tx"),
		"transaction_address": crypto.HashString("addr"),
	})

	pr := NewPoolRequester(techDB)

	lis, _ := net.Listen("tcp", ":5000")
	defer lis.Close()
	grpcServer := grpc.NewServer()
	api.RegisterLockServiceServer(grpcServer, lockSrv)
	go grpcServer.Serve(lis)

	pool, _ := consensus.FindStoragePool("addr")
	assert.Nil(t, pr.RequestTransactionUnlock(pool, crypto.HashString("tx"), crypto.HashString("addr")))
	assert.Len(t, lockDB.locks, 0)
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
	minerKey, _ := shared.NewMinerKeyPair(pub, pv)
	techDB.minerKeys = append(techDB.minerKeys, minerKey)

	pr := NewPoolRequester(techDB)

	miningSrv := NewMiningServer(techDB, pr, pub, pv)

	lis, err := net.Listen("tcp", ":5000")
	defer lis.Close()
	grpcServer := grpc.NewServer()
	api.RegisterMiningServiceServer(grpcServer, miningSrv)
	go grpcServer.Serve(lis)

	data := map[string]string{
		"encrypted_address_by_miner": hex.EncodeToString([]byte("addr")),
		"encrypted_wallet":           hex.EncodeToString([]byte("wallet")),
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
	v, _ := chain.NewMinerValidation(chain.ValidationOK, time.Now(), pub, vSig)
	mv, _ := chain.NewMasterValidation([]string{}, pub, v)

	pool, _ := consensus.FindValidationPool(tx)
	valids, err := pr.RequestTransactionValidations(pool, tx, 1, mv)
	assert.Nil(t, err)

	assert.Len(t, valids, 1)
	assert.Equal(t, pub, valids[0].MinerPublicKey())
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
	minerKey, _ := shared.NewMinerKeyPair(pub, pv)
	techDB.minerKeys = append(techDB.minerKeys, minerKey)

	techDB.emKeys = append(techDB.emKeys, kp)

	pr := NewPoolRequester(techDB)

	chainSrv := NewChainServer(chainDB, techDB, pr)

	lis, _ := net.Listen("tcp", ":5000")
	defer lis.Close()
	grpcServer := grpc.NewServer()
	api.RegisterChainServiceServer(grpcServer, chainSrv)
	go grpcServer.Serve(lis)

	data := map[string]string{
		"encrypted_address_by_miner": hex.EncodeToString([]byte("addr")),
		"encrypted_wallet":           hex.EncodeToString([]byte("wallet")),
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
	v, _ := chain.NewMinerValidation(chain.ValidationOK, time.Now(), pub, vSig)
	mv, _ := chain.NewMasterValidation([]string{}, pub, v)

	pool, _ := consensus.FindStoragePool("addr")
	tx.Mined(mv, []chain.MinerValidation{v})
	assert.Nil(t, pr.RequestTransactionStorage(pool, 1, tx))

	assert.Len(t, chainDB.keychains, 1)
	assert.Equal(t, crypto.HashBytes(txBytes), chainDB.keychains[0].TransactionHash())
}

/*
Scenario: Send request to get last transaction
	Given a keychain transaction stored
	When I want to request a miner to get the last transaction from the address
	Then I get the last transaction
*/
func TestSendGetLastTransaction(t *testing.T) {

	pub, pv := crypto.GenerateKeys()

	chainDB := &mockChainDB{}
	techDB := &mockTechDB{}
	minerKey, _ := shared.NewMinerKeyPair(pub, pv)
	techDB.minerKeys = append(techDB.minerKeys, minerKey)

	pr := NewPoolRequester(techDB)

	chainSrv := NewChainServer(chainDB, techDB, pr)

	lis, _ := net.Listen("tcp", ":5000")
	defer lis.Close()
	grpcServer := grpc.NewServer()
	api.RegisterChainServiceServer(grpcServer, chainSrv)
	go grpcServer.Serve(lis)

	data := map[string]string{
		"encrypted_address_by_miner": hex.EncodeToString([]byte("addr")),
		"encrypted_wallet":           hex.EncodeToString([]byte("wallet")),
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

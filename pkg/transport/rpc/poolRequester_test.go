package rpc

import (
	"encoding/hex"
	"encoding/json"
	"net"
	"testing"
	"time"

	"github.com/uniris/uniris-core/pkg/shared"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

	api "github.com/uniris/uniris-core/api/protobuf-spec"
	"github.com/uniris/uniris-core/pkg/crypto"

	"github.com/uniris/uniris-core/pkg/transaction"
)

/*
Scenario: Request transction lock on a pool
	Given a transaction to lock
	When I request to lock it
	Then the lock is stored in the database
*/
func TestRequestTransactionLock(t *testing.T) {

	pub, pv := crypto.GenerateKeys()

	lockRepo := &mockLockRepository{}
	txRepo := &mockTxRepository{}
	sharedRepo := &mockSharedRepo{}
	txSrv := newTransactionServer(txRepo, lockRepo, sharedRepo, pub, pv)

	sharedSrv := shared.NewService(sharedRepo)
	pr := NewPoolRequester(sharedSrv)

	lis, _ := net.Listen("tcp", ":3545")
	defer lis.Close()
	grpcServer := grpc.NewServer()
	api.RegisterTransactionServiceServer(grpcServer, txSrv)
	go grpcServer.Serve(lis)

	lock, _ := transaction.NewLock(crypto.HashString("tx"), crypto.HashString("addr"), pub)
	pool := transaction.Pool{
		transaction.NewPoolMember(net.ParseIP("127.0.0.1"), 3545),
	}
	assert.Nil(t, pr.RequestTransactionLock(pool, lock))

	assert.Len(t, lockRepo.locks, 1)
	assert.Equal(t, crypto.HashString("addr"), lockRepo.locks[0].Address())
}

/*
Scenario: Request transction lock on a pool
	Given a transaction already locked
	When I request to unlock it
	Then the lock is removed from the database
*/
func TestRequestTransactionUnlock(t *testing.T) {

	pub, pv := crypto.GenerateKeys()

	lockRepo := &mockLockRepository{}
	txRepo := &mockTxRepository{}
	sharedRepo := &mockSharedRepo{}
	txSrv := newTransactionServer(txRepo, lockRepo, sharedRepo, pub, pv)

	sharedSrv := shared.NewService(sharedRepo)
	pr := NewPoolRequester(sharedSrv)

	lis, _ := net.Listen("tcp", ":3545")
	defer lis.Close()
	grpcServer := grpc.NewServer()
	api.RegisterTransactionServiceServer(grpcServer, txSrv)
	go grpcServer.Serve(lis)

	lock, _ := transaction.NewLock(crypto.HashString("tx"), crypto.HashString("addr"), pub)
	pool := transaction.Pool{
		transaction.NewPoolMember(net.ParseIP("127.0.0.1"), 3545),
	}
	assert.Nil(t, pr.RequestTransactionLock(pool, lock))

	assert.Len(t, lockRepo.locks, 1)
	assert.Equal(t, crypto.HashString("addr"), lockRepo.locks[0].Address())

	assert.Nil(t, pr.RequestTransactionUnlock(pool, lock))
	assert.Len(t, lockRepo.locks, 0)
}

/*
Scenario: Request transaction validation confirmation
	Given a transaction to validate
	When I request to confirm the validation
	Then I get a validation
*/
func TestRequestConfirmValidation(t *testing.T) {

	pub, pv := crypto.GenerateKeys()

	lockRepo := &mockLockRepository{}
	txRepo := &mockTxRepository{}
	sharedRepo := &mockSharedRepo{}
	txSrv := newTransactionServer(txRepo, lockRepo, sharedRepo, pub, pv)

	lis, err := net.Listen("tcp", ":3545")
	defer lis.Close()
	grpcServer := grpc.NewServer()
	api.RegisterTransactionServiceServer(grpcServer, txSrv)
	go grpcServer.Serve(lis)

	sharedSrv := shared.NewService(sharedRepo)
	pr := NewPoolRequester(sharedSrv)

	sk, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("pvkey")), pub)
	sharedSrv.StoreSharedEmitterKeyPair(sk)

	data := map[string]string{
		"encrypted_address": hex.EncodeToString([]byte("addr")),
		"encrypted_wallet":  hex.EncodeToString([]byte("wallet")),
	}
	prop, _ := transaction.NewProposal(sk)
	txRaw := map[string]interface{}{
		"address":    crypto.HashString("addr"),
		"data":       data,
		"type":       transaction.KeychainType,
		"timestamp":  time.Now().Unix(),
		"public_key": pub,
		"proposal":   prop,
	}
	txBytes, _ := json.Marshal(txRaw)
	sig, _ := crypto.Sign(string(txBytes), pv)
	txRaw["signature"] = sig
	txRaw["em_signature"] = sig
	txBytes, _ = json.Marshal(txRaw)

	tx, _ := transaction.New(crypto.HashString("addr"), transaction.KeychainType, data, time.Now(), pub, sig, sig, prop, crypto.HashBytes(txBytes))
	vBytes, _ := json.Marshal(map[string]interface{}{
		"status":     transaction.ValidationOK,
		"public_key": pub,
		"timestamp":  time.Now().Unix(),
	})
	vSig, _ := crypto.Sign(string(vBytes), pv)
	v, _ := transaction.NewMinerValidation(transaction.ValidationOK, time.Now(), pub, vSig)
	mv, _ := transaction.NewMasterValidation(transaction.Pool{}, pub, v)

	vChan := make(chan transaction.MinerValidation)

	pool := transaction.Pool{
		transaction.NewPoolMember(net.ParseIP("127.0.0.1"), 3545),
	}
	go pr.RequestTransactionValidations(pool, tx, mv, vChan)

	valid := <-vChan
	assert.Equal(t, pub, valid.MinerPublicKey())
	assert.Equal(t, transaction.ValidationOK, valid.Status())
	ok, err := valid.IsValid()
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

	lockRepo := &mockLockRepository{}
	txRepo := &mockTxRepository{}
	sharedRepo := &mockSharedRepo{}
	txSrv := newTransactionServer(txRepo, lockRepo, sharedRepo, pub, pv)

	lis, _ := net.Listen("tcp", ":3545")
	defer lis.Close()
	grpcServer := grpc.NewServer()
	api.RegisterTransactionServiceServer(grpcServer, txSrv)
	go grpcServer.Serve(lis)

	sharedSrv := shared.NewService(sharedRepo)
	pr := NewPoolRequester(sharedSrv)

	sk, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("pvkey")), pub)
	sharedSrv.StoreSharedEmitterKeyPair(sk)

	data := map[string]string{
		"encrypted_address": hex.EncodeToString([]byte("addr")),
		"encrypted_wallet":  hex.EncodeToString([]byte("wallet")),
	}
	prop, _ := transaction.NewProposal(sk)
	txRaw := map[string]interface{}{
		"address":    crypto.HashString("addr"),
		"data":       data,
		"type":       transaction.KeychainType,
		"timestamp":  time.Now().Unix(),
		"public_key": pub,
		"proposal":   prop,
	}
	txBytes, _ := json.Marshal(txRaw)
	sig, _ := crypto.Sign(string(txBytes), pv)
	txRaw["signature"] = sig
	txRaw["em_signature"] = sig
	txBytes, _ = json.Marshal(txRaw)

	tx, _ := transaction.New(crypto.HashString("addr"), transaction.KeychainType, data, time.Now(), pub, sig, sig, prop, crypto.HashBytes(txBytes))
	vBytes, _ := json.Marshal(map[string]interface{}{
		"status":     transaction.ValidationOK,
		"public_key": pub,
		"timestamp":  time.Now().Unix(),
	})
	vSig, _ := crypto.Sign(string(vBytes), pv)
	v, _ := transaction.NewMinerValidation(transaction.ValidationOK, time.Now(), pub, vSig)
	mv, _ := transaction.NewMasterValidation(transaction.Pool{}, pub, v)

	ackChan := make(chan bool)

	pool := transaction.Pool{
		transaction.NewPoolMember(net.ParseIP("127.0.0.1"), 3545),
	}
	tx.AddMining(mv, []transaction.MinerValidation{v})
	go pr.RequestTransactionStorage(pool, tx, ackChan)

	<-ackChan

	assert.Len(t, txRepo.keychains, 1)
	assert.Equal(t, crypto.HashBytes(txBytes), txRepo.keychains[0].TransactionHash())
}

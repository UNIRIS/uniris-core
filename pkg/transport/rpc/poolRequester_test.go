package rpc

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
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

	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pub, _ := x509.MarshalPKIXPublicKey(key.Public())
	pv, _ := x509.MarshalECPrivateKey(key)

	pr := NewPoolRequester(hex.EncodeToString(pub), hex.EncodeToString(pv))

	lockRepo := &mockLockRepository{}
	lockSrv := transaction.NewLockService(lockRepo)
	txSrv := NewTransactionServer(transaction.StorageService{}, lockSrv, transaction.MiningService{}, hex.EncodeToString(pub), hex.EncodeToString(pv))

	lis, _ := net.Listen("tcp", ":3545")
	defer lis.Close()
	grpcServer := grpc.NewServer()
	api.RegisterTransactionServiceServer(grpcServer, txSrv)
	go grpcServer.Serve(lis)

	lock, _ := transaction.NewLock(crypto.HashString("tx"), crypto.HashString("addr"), hex.EncodeToString(pub))
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

	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pub, _ := x509.MarshalPKIXPublicKey(key.Public())
	pv, _ := x509.MarshalECPrivateKey(key)

	pr := NewPoolRequester(hex.EncodeToString(pub), hex.EncodeToString(pv))

	lockRepo := &mockLockRepository{}
	lockSrv := transaction.NewLockService(lockRepo)
	txSrv := NewTransactionServer(transaction.StorageService{}, lockSrv, transaction.MiningService{}, hex.EncodeToString(pub), hex.EncodeToString(pv))

	lis, _ := net.Listen("tcp", ":3545")
	defer lis.Close()
	grpcServer := grpc.NewServer()
	api.RegisterTransactionServiceServer(grpcServer, txSrv)
	go grpcServer.Serve(lis)

	lock, _ := transaction.NewLock(crypto.HashString("tx"), crypto.HashString("addr"), hex.EncodeToString(pub))
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

	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pub, _ := x509.MarshalPKIXPublicKey(key.Public())
	pv, _ := x509.MarshalECPrivateKey(key)
	miningSrv := transaction.NewMiningService(nil, transaction.PoolFindingService{}, shared.Service{}, "127.0.0.1", hex.EncodeToString(pub), hex.EncodeToString(pv))

	txSrv := NewTransactionServer(transaction.StorageService{}, transaction.LockService{}, miningSrv, hex.EncodeToString(pub), hex.EncodeToString(pv))
	lis, err := net.Listen("tcp", ":3545")
	defer lis.Close()
	grpcServer := grpc.NewServer()
	api.RegisterTransactionServiceServer(grpcServer, txSrv)
	go grpcServer.Serve(lis)

	pr := NewPoolRequester(hex.EncodeToString(pub), hex.EncodeToString(pv))

	sharedSrv := shared.NewService(&mockSharedRepo{})
	sk, _ := shared.NewKeyPair(hex.EncodeToString([]byte("pvkey")), hex.EncodeToString(pub))
	sharedSrv.StoreSharedEmitterKeyPair(sk)

	data := map[string]string{
		"encrypted_address": hex.EncodeToString([]byte("addr")),
		"encrypted_wallet":  hex.EncodeToString([]byte("wallet")),
	}
	prop, _ := transaction.NewProposal(sk)
	txRaw, _ := json.Marshal(struct {
		Address   string               `json:"address"`
		Data      map[string]string    `json:"data"`
		Timestamp int64                `json:"timestamp"`
		Type      transaction.Type     `json:"type"`
		PublicKey string               `json:"public_key"`
		Proposal  transaction.Proposal `json:"proposal"`
	}{
		Address:   crypto.HashString("addr"),
		Type:      transaction.KeychainType,
		Data:      data,
		Timestamp: time.Now().Unix(),
		PublicKey: hex.EncodeToString(pub),
		Proposal:  prop,
	})

	sig, _ := crypto.Sign(string(txRaw), hex.EncodeToString(pv))

	txSigned, _ := json.Marshal(struct {
		Address          string               `json:"address"`
		Data             map[string]string    `json:"data"`
		Timestamp        int64                `json:"timestamp"`
		Type             transaction.Type     `json:"type"`
		PublicKey        string               `json:"public_key"`
		Proposal         transaction.Proposal `json:"proposal"`
		Signature        string               `json:"signature"`
		EmitterSignature string               `json:"em_signature"`
	}{
		Address:          crypto.HashString("addr"),
		Type:             transaction.KeychainType,
		Data:             data,
		Timestamp:        time.Now().Unix(),
		PublicKey:        hex.EncodeToString(pub),
		Proposal:         prop,
		Signature:        sig,
		EmitterSignature: sig,
	})

	tx, _ := transaction.New(crypto.HashString("addr"), transaction.KeychainType, data, time.Now(), hex.EncodeToString(pub), sig, sig, prop, crypto.HashBytes(txSigned))
	vBytes, _ := json.Marshal(struct {
		Status         transaction.ValidationStatus `json:"status"`
		MinerPublicKey string                       `json:"public_key"`
		Timestamp      int64                        `json:"timestamp"`
	}{
		Status:         transaction.ValidationOK,
		MinerPublicKey: hex.EncodeToString(pub),
		Timestamp:      time.Now().Unix(),
	})
	vSig, _ := crypto.Sign(string(vBytes), hex.EncodeToString(pv))
	v, _ := transaction.NewMinerValidation(transaction.ValidationOK, time.Now(), hex.EncodeToString(pub), vSig)
	mv, _ := transaction.NewMasterValidation(transaction.Pool{}, hex.EncodeToString(pub), v)

	vChan := make(chan transaction.MinerValidation)

	pool := transaction.Pool{
		transaction.NewPoolMember(net.ParseIP("127.0.0.1"), 3545),
	}
	go pr.RequestTransactionValidations(pool, tx, mv, vChan)

	valid := <-vChan
	assert.Equal(t, hex.EncodeToString(pub), valid.MinerPublicKey())
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

	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pub, _ := x509.MarshalPKIXPublicKey(key.Public())
	pv, _ := x509.MarshalECPrivateKey(key)
	miningSrv := transaction.NewMiningService(nil, transaction.PoolFindingService{}, shared.Service{}, "127.0.0.1", hex.EncodeToString(pub), hex.EncodeToString(pv))
	txRepo := &mockTxRepository{}
	storeSrv := transaction.NewStorageService(txRepo, miningSrv)

	txSrv := NewTransactionServer(storeSrv, transaction.LockService{}, miningSrv, hex.EncodeToString(pub), hex.EncodeToString(pv))
	lis, _ := net.Listen("tcp", ":3545")
	defer lis.Close()
	grpcServer := grpc.NewServer()
	api.RegisterTransactionServiceServer(grpcServer, txSrv)
	go grpcServer.Serve(lis)

	pr := NewPoolRequester(hex.EncodeToString(pub), hex.EncodeToString(pv))

	sharedSrv := shared.NewService(&mockSharedRepo{})
	sk, _ := shared.NewKeyPair(hex.EncodeToString([]byte("pvkey")), hex.EncodeToString(pub))
	sharedSrv.StoreSharedEmitterKeyPair(sk)

	data := map[string]string{
		"encrypted_address": hex.EncodeToString([]byte("addr")),
		"encrypted_wallet":  hex.EncodeToString([]byte("wallet")),
	}
	prop, _ := transaction.NewProposal(sk)
	txRaw, _ := json.Marshal(struct {
		Address   string               `json:"address"`
		Data      map[string]string    `json:"data"`
		Timestamp int64                `json:"timestamp"`
		Type      transaction.Type     `json:"type"`
		PublicKey string               `json:"public_key"`
		Proposal  transaction.Proposal `json:"proposal"`
	}{
		Address:   crypto.HashString("addr"),
		Type:      transaction.KeychainType,
		Data:      data,
		Timestamp: time.Now().Unix(),
		PublicKey: hex.EncodeToString(pub),
		Proposal:  prop,
	})

	sig, _ := crypto.Sign(string(txRaw), hex.EncodeToString(pv))

	txSigned, _ := json.Marshal(struct {
		Address          string               `json:"address"`
		Data             map[string]string    `json:"data"`
		Timestamp        int64                `json:"timestamp"`
		Type             transaction.Type     `json:"type"`
		PublicKey        string               `json:"public_key"`
		Proposal         transaction.Proposal `json:"proposal"`
		Signature        string               `json:"signature"`
		EmitterSignature string               `json:"em_signature"`
	}{
		Address:          crypto.HashString("addr"),
		Type:             transaction.KeychainType,
		Data:             data,
		Timestamp:        time.Now().Unix(),
		PublicKey:        hex.EncodeToString(pub),
		Proposal:         prop,
		Signature:        sig,
		EmitterSignature: sig,
	})

	tx, _ := transaction.New(crypto.HashString("addr"), transaction.KeychainType, data, time.Now(), hex.EncodeToString(pub), sig, sig, prop, crypto.HashBytes(txSigned))
	vBytes, _ := json.Marshal(struct {
		Status         transaction.ValidationStatus `json:"status"`
		MinerPublicKey string                       `json:"public_key"`
		Timestamp      int64                        `json:"timestamp"`
	}{
		Status:         transaction.ValidationOK,
		MinerPublicKey: hex.EncodeToString(pub),
		Timestamp:      time.Now().Unix(),
	})
	vSig, _ := crypto.Sign(string(vBytes), hex.EncodeToString(pv))
	v, _ := transaction.NewMinerValidation(transaction.ValidationOK, time.Now(), hex.EncodeToString(pub), vSig)
	mv, _ := transaction.NewMasterValidation(transaction.Pool{}, hex.EncodeToString(pub), v)

	ackChan := make(chan bool)

	pool := transaction.Pool{
		transaction.NewPoolMember(net.ParseIP("127.0.0.1"), 3545),
	}
	tx.AddMining(mv, []transaction.MinerValidation{v})
	go pr.RequestTransactionStorage(pool, tx, ackChan)

	<-ackChan

	assert.Len(t, txRepo.keychains, 1)
	assert.Equal(t, crypto.HashBytes(txSigned), txRepo.keychains[0].TransactionHash())
}

package rpc

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

	api "github.com/uniris/uniris-core/api/protobuf-spec"
	"github.com/uniris/uniris-core/pkg/crypto"
	"github.com/uniris/uniris-core/pkg/shared"
	"github.com/uniris/uniris-core/pkg/transaction"
)

/*
Scenario: Receive a get transaction status from the API
	Given a transaction status request incoming from the API
	When I want to get the transaction status
	Then I get the right pool and send a GRPC request to the pool to get the status
*/
func TestHandleGetTransactionStatusInternal(t *testing.T) {

	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pub, _ := x509.MarshalPKIXPublicKey(key.Public())
	pv, _ := x509.MarshalECPrivateKey(key)

	txRepo := &mockTxRepository{}
	lockSrv := transaction.NewLockService(&mockLockRepository{})
	poolFindSrv := transaction.NewPoolFindingService(NewPoolRetriever(hex.EncodeToString(pub), hex.EncodeToString(pv)))
	sharedService := shared.NewService(&mockSharedRepo{})
	miningSrv := transaction.NewMiningService(&mockPoolRequester{}, poolFindSrv, sharedService, "127.0.0.1", hex.EncodeToString(pub), hex.EncodeToString(pv))
	storeSrv := transaction.NewStorageService(txRepo, miningSrv)

	txSrv := NewTransactionServer(storeSrv, lockSrv, miningSrv, hex.EncodeToString(pub), hex.EncodeToString(pv))
	intSrv := NewInternalServer(poolFindSrv, miningSrv, hex.EncodeToString(pub), hex.EncodeToString(pv))

	lis, _ := net.Listen("tcp", ":3545")
	defer lis.Close()
	grpcServer := grpc.NewServer()
	api.RegisterTransactionServiceServer(grpcServer, txSrv)
	go grpcServer.Serve(lis)

	req := &api.InternalTransactionStatusRequest{
		TransactionHash:    crypto.HashString("tx"),
		TransactionAddress: crypto.HashString("address"),
		Timestamp:          time.Now().Unix(),
	}
	res, err := intSrv.GetTransactionStatus(context.TODO(), req)
	assert.Nil(t, err)
	assert.Equal(t, api.TransactionStatusResponse_UNKNOWN, res.Status)

	resBytes, _ := json.Marshal(&api.TransactionStatusResponse{
		Status:    res.Status,
		Timestamp: res.Timestamp,
	})
	assert.Nil(t, crypto.VerifySignature(string(resBytes), hex.EncodeToString(pub), res.SignatureResponse))

}

/*
Scenario: Forward an transaction request from the API to the master miner
	Given a transaction come from the API
	When I want to process it
	Then I forward to the master miner and reply the transaction hash
*/
func TestHandleIncomingTransaction(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pub, _ := x509.MarshalPKIXPublicKey(key.Public())
	pv, _ := x509.MarshalECPrivateKey(key)

	txRepo := &mockTxRepository{}
	poolR := &mockPoolRequester{
		repo: txRepo,
	}
	lockSrv := transaction.NewLockService(&mockLockRepository{})
	poolFindSrv := transaction.NewPoolFindingService(NewPoolRetriever(hex.EncodeToString(pub), hex.EncodeToString(pv)))
	sharedService := shared.NewService(&mockSharedRepo{})
	miningSrv := transaction.NewMiningService(poolR, poolFindSrv, sharedService, "127.0.0.1", hex.EncodeToString(pub), hex.EncodeToString(pv))
	storeSrv := transaction.NewStorageService(txRepo, miningSrv)

	sk, _ := shared.NewKeyPair(hex.EncodeToString([]byte("pvkey")), hex.EncodeToString(pub))
	sharedService.StoreSharedEmitterKeyPair(sk)

	txSrv := NewTransactionServer(storeSrv, lockSrv, miningSrv, hex.EncodeToString(pub), hex.EncodeToString(pv))
	intSrv := NewInternalServer(poolFindSrv, miningSrv, hex.EncodeToString(pub), hex.EncodeToString(pv))

	lis, _ := net.Listen("tcp", ":3545")
	defer lis.Close()
	grpcServer := grpc.NewServer()
	api.RegisterTransactionServiceServer(grpcServer, txSrv)
	go grpcServer.Serve(lis)

	tx := map[string]interface{}{
		"address": crypto.HashString("addr"),
		"data": map[string]string{
			"encrypted_address": hex.EncodeToString([]byte("addr")),
			"encrypted_wallet":  hex.EncodeToString([]byte("wallet")),
		},
		"timestamp":  time.Now().Unix(),
		"type":       int(transaction.KeychainType),
		"public_key": hex.EncodeToString(pub),
		"proposal": map[string]interface{}{
			"shared_emitter_keys": map[string]string{
				"encrypted_private_key": hex.EncodeToString([]byte("encPV")),
				"public_key":            hex.EncodeToString(pub),
			},
		},
	}
	txBytes, _ := json.Marshal(tx)
	sig, _ := crypto.Sign(string(txBytes), hex.EncodeToString(pv))
	tx["signature"] = sig
	tx["em_signature"] = sig
	txBytes, _ = json.Marshal(tx)

	cipherTx, _ := crypto.Encrypt(string(txBytes), hex.EncodeToString(pub))
	res, err := intSrv.HandleTransaction(context.TODO(), &api.IncomingTransaction{
		EncryptedTransaction: cipherTx,
		Timestamp:            time.Now().Unix(),
	})
	assert.Nil(t, err)
	assert.NotEmpty(t, res.TransactionHash)

	resBytes, _ := json.Marshal(&api.TransactionResult{
		Timestamp:       res.Timestamp,
		TransactionHash: res.TransactionHash,
	})
	assert.Nil(t, crypto.VerifySignature(string(resBytes), hex.EncodeToString(pub), res.Signature))
	assert.Equal(t, crypto.HashBytes(txBytes), res.TransactionHash)

	time.Sleep(2 * time.Second)

	assert.Len(t, txRepo.keychains, 1)
	assert.Equal(t, res.TransactionHash, txRepo.keychains[0].TransactionHash())

}

/*
Scenario: Get account created
	Given a account created (id and keychain transaction proceed)
	When I want to retreive the account
	Then I provide the ID hash and I get retrieve the encrypted wallet and the encrypted aes key
*/
func TestHandleGetAccount(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pub, _ := x509.MarshalPKIXPublicKey(key.Public())
	pv, _ := x509.MarshalECPrivateKey(key)

	txRepo := &mockTxRepository{}
	poolR := &mockPoolRequester{
		repo: txRepo,
	}
	lockSrv := transaction.NewLockService(&mockLockRepository{})
	poolFindSrv := transaction.NewPoolFindingService(NewPoolRetriever(hex.EncodeToString(pub), hex.EncodeToString(pv)))
	sharedService := shared.NewService(&mockSharedRepo{})
	miningSrv := transaction.NewMiningService(poolR, poolFindSrv, sharedService, "127.0.0.1", hex.EncodeToString(pub), hex.EncodeToString(pv))
	storeSrv := transaction.NewStorageService(txRepo, miningSrv)

	sk, _ := shared.NewKeyPair(hex.EncodeToString([]byte("pvkey")), hex.EncodeToString(pub))
	sharedService.StoreSharedEmitterKeyPair(sk)

	txSrv := NewTransactionServer(storeSrv, lockSrv, miningSrv, hex.EncodeToString(pub), hex.EncodeToString(pv))
	intSrv := NewInternalServer(poolFindSrv, miningSrv, hex.EncodeToString(pub), hex.EncodeToString(pv))

	lis, _ := net.Listen("tcp", ":3545")
	defer lis.Close()
	grpcServer := grpc.NewServer()
	api.RegisterTransactionServiceServer(grpcServer, txSrv)
	go grpcServer.Serve(lis)

	encAddr, _ := crypto.Encrypt(crypto.HashString("addr"), hex.EncodeToString(pub))

	//First send the ID transaction
	txID := map[string]interface{}{
		"address": crypto.HashString("idHash"),
		"data": map[string]string{
			"encrypted_address_by_id":    encAddr,
			"encrypted_address_by_robot": encAddr,
			"encrypted_aes_key":          hex.EncodeToString([]byte("aesKey")),
		},
		"timestamp":  time.Now().Unix(),
		"type":       int(transaction.IDType),
		"public_key": hex.EncodeToString(pub),
		"proposal": map[string]interface{}{
			"shared_emitter_keys": map[string]string{
				"encrypted_private_key": hex.EncodeToString([]byte("encPV")),
				"public_key":            hex.EncodeToString(pub),
			},
		},
	}
	txIDBytes, _ := json.Marshal(txID)
	sigID, _ := crypto.Sign(string(txIDBytes), hex.EncodeToString(pv))
	txID["signature"] = sigID
	txID["em_signature"] = sigID
	txIDBytes, _ = json.Marshal(txID)

	cipherTx, _ := crypto.Encrypt(string(txIDBytes), hex.EncodeToString(pub))
	res, err := intSrv.HandleTransaction(context.TODO(), &api.IncomingTransaction{
		EncryptedTransaction: cipherTx,
		Timestamp:            time.Now().Unix(),
	})
	assert.Nil(t, err)
	assert.NotEmpty(t, res.TransactionHash)

	time.Sleep(1 * time.Second)
	assert.Equal(t, crypto.HashString("idHash"), txRepo.ids[0].Address())

	//Then send the keychain transaction
	txKeychain := map[string]interface{}{
		"address": crypto.HashString("addr"),
		"data": map[string]string{
			"encrypted_address": hex.EncodeToString([]byte("addr")),
			"encrypted_wallet":  hex.EncodeToString([]byte("wallet")),
		},
		"timestamp":  time.Now().Unix(),
		"type":       int(transaction.KeychainType),
		"public_key": hex.EncodeToString(pub),
		"proposal": map[string]interface{}{
			"shared_emitter_keys": map[string]string{
				"encrypted_private_key": hex.EncodeToString([]byte("encPV")),
				"public_key":            hex.EncodeToString(pub),
			},
		},
	}
	txKeychainBytes, _ := json.Marshal(txKeychain)
	sigKeychain, _ := crypto.Sign(string(txKeychainBytes), hex.EncodeToString(pv))
	txKeychain["signature"] = sigKeychain
	txKeychain["em_signature"] = sigKeychain
	txKeychainBytes, _ = json.Marshal(txKeychain)

	cipherTx2, _ := crypto.Encrypt(string(txKeychainBytes), hex.EncodeToString(pub))
	res2, err := intSrv.HandleTransaction(context.TODO(), &api.IncomingTransaction{
		EncryptedTransaction: cipherTx2,
		Timestamp:            time.Now().Unix(),
	})
	assert.Nil(t, err)
	assert.NotEmpty(t, res2.TransactionHash)

	time.Sleep(1 * time.Second)
	assert.Equal(t, crypto.HashString("addr"), txRepo.keychains[0].Address())

	encIDHash, _ := crypto.Encrypt(crypto.HashString("idHash"), hex.EncodeToString(pub))
	req := &api.GetAccountRequest{
		EncryptedIdAddress: encIDHash,
		Timestamp:          time.Now().Unix(),
	}

	resGet, err := intSrv.GetAccount(context.TODO(), req)
	assert.Nil(t, err)
	assert.NotEmpty(t, resGet.EncryptedAesKey)
	assert.NotEmpty(t, resGet.EncryptedWallet)
}

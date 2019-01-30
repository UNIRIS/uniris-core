package rpc

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"log"
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

	req := &api.TransactionStatusRequest{
		TransactionHash: crypto.HashString("tx"),
		Timestamp:       time.Now().Unix(),
	}
	reqBytes, _ := json.Marshal(req)
	reqSig, _ := crypto.Sign(string(reqBytes), hex.EncodeToString(pv))
	req.SignatureRequest = reqSig
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

	data := map[string]string{
		"encrypted_address": hex.EncodeToString([]byte("addr")),
		"encrypted_wallet":  hex.EncodeToString([]byte("wallet")),
	}

	tx, _ := json.Marshal(txRaw{
		Address: crypto.HashString("addr"),
		Data:    data,
		Proposal: txProp{
			SharedEmitterKeys: txSharedKeys{
				EncryptedPrivateKey: hex.EncodeToString([]byte("encPV")),
				PublicKey:           hex.EncodeToString(pub),
			},
		},
		Timestamp: time.Now().Unix(),
		Type:      int(transaction.KeychainType),
		PublicKey: hex.EncodeToString(pub),
	})
	sig, _ := crypto.Sign(string(tx), hex.EncodeToString(pv))
	txSigned, _ := json.Marshal(txSigned{
		Address: crypto.HashString("addr"),
		Data:    data,
		Proposal: txProp{
			SharedEmitterKeys: txSharedKeys{
				EncryptedPrivateKey: hex.EncodeToString([]byte("encPV")),
				PublicKey:           hex.EncodeToString(pub),
			},
		},
		Timestamp:        time.Now().Unix(),
		Type:             int(transaction.KeychainType),
		PublicKey:        hex.EncodeToString(pub),
		EmitterSignature: sig,
		Signature:        sig,
	})

	cipherTx, _ := crypto.Encrypt(string(txSigned), hex.EncodeToString(pub))
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
	assert.Equal(t, crypto.HashBytes(txSigned), res.TransactionHash)

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
	iddata := map[string]string{
		"encrypted_address_by_id":    encAddr,
		"encrypted_address_by_robot": encAddr,
		"encrypted_aes_key":          hex.EncodeToString([]byte("aesKey")),
	}

	txRaw1, _ := json.Marshal(txRaw{
		Address: crypto.HashString("idHash"),
		Data:    iddata,
		Proposal: txProp{
			SharedEmitterKeys: txSharedKeys{
				EncryptedPrivateKey: hex.EncodeToString([]byte("encPV")),
				PublicKey:           hex.EncodeToString(pub),
			},
		},
		Timestamp: time.Now().Unix(),
		Type:      int(transaction.IDType),
		PublicKey: hex.EncodeToString(pub),
	})
	sig, _ := crypto.Sign(string(txRaw1), hex.EncodeToString(pv))
	txSigned1, _ := json.Marshal(txSigned{
		Address: crypto.HashString("idHash"),
		Data:    iddata,
		Proposal: txProp{
			SharedEmitterKeys: txSharedKeys{
				EncryptedPrivateKey: hex.EncodeToString([]byte("encPV")),
				PublicKey:           hex.EncodeToString(pub),
			},
		},
		Timestamp:        time.Now().Unix(),
		Type:             int(transaction.IDType),
		PublicKey:        hex.EncodeToString(pub),
		EmitterSignature: sig,
		Signature:        sig,
	})

	cipherTx, _ := crypto.Encrypt(string(txSigned1), hex.EncodeToString(pub))
	res, err := intSrv.HandleTransaction(context.TODO(), &api.IncomingTransaction{
		EncryptedTransaction: cipherTx,
		Timestamp:            time.Now().Unix(),
	})
	assert.Nil(t, err)
	assert.NotEmpty(t, res.TransactionHash)

	time.Sleep(1 * time.Second)
	assert.Equal(t, crypto.HashString("idHash"), txRepo.ids[0].Address())
	log.Print(crypto.HashString("idHash"))

	//Then send the keychain transaction
	keychainData := map[string]string{
		"encrypted_address": encAddr,
		"encrypted_wallet":  hex.EncodeToString([]byte("wallet")),
	}

	txRaw2, _ := json.Marshal(txRaw{
		Address: crypto.HashString("addr"),
		Data:    keychainData,
		Proposal: txProp{
			SharedEmitterKeys: txSharedKeys{
				EncryptedPrivateKey: hex.EncodeToString([]byte("encPV")),
				PublicKey:           hex.EncodeToString(pub),
			},
		},
		Timestamp: time.Now().Unix(),
		Type:      int(transaction.KeychainType),
		PublicKey: hex.EncodeToString(pub),
	})
	sig2, _ := crypto.Sign(string(txRaw2), hex.EncodeToString(pv))
	txSigned2, _ := json.Marshal(txSigned{
		Address: crypto.HashString("addr"),
		Data:    keychainData,
		Proposal: txProp{
			SharedEmitterKeys: txSharedKeys{
				EncryptedPrivateKey: hex.EncodeToString([]byte("encPV")),
				PublicKey:           hex.EncodeToString(pub),
			},
		},
		Timestamp:        time.Now().Unix(),
		Type:             int(transaction.KeychainType),
		PublicKey:        hex.EncodeToString(pub),
		EmitterSignature: sig2,
		Signature:        sig2,
	})

	cipherTx2, _ := crypto.Encrypt(string(txSigned2), hex.EncodeToString(pub))
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
	reqBytes, _ := json.Marshal(req)
	sigReq, _ := crypto.Sign(string(reqBytes), hex.EncodeToString(pv))
	req.SignatureRequest = sigReq

	resGet, err := intSrv.GetAccount(context.TODO(), req)
	assert.Nil(t, err)
	assert.NotEmpty(t, resGet.EncryptedAesKey)
	assert.NotEmpty(t, resGet.EncryptedWallet)
}

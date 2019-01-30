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

	"github.com/stretchr/testify/assert"

	api "github.com/uniris/uniris-core/api/protobuf-spec"
	"github.com/uniris/uniris-core/pkg/crypto"
	"github.com/uniris/uniris-core/pkg/shared"
	"github.com/uniris/uniris-core/pkg/transaction"
	"google.golang.org/grpc"
)

/*
Scenario: Send request to get last transaction
	Given a keychain transaction stored
	When I want to request a miner to get the last transaction from the address
	Then I get the last transaction
*/
func TestSendGetLastTransaction(t *testing.T) {

	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pub, _ := x509.MarshalPKIXPublicKey(key.Public())
	pv, _ := x509.MarshalECPrivateKey(key)

	txRepo := &mockTxRepository{}
	lockSrv := transaction.NewLockService(&mockLockRepository{})
	poolFindSrv := transaction.NewPoolFindingService(NewPoolRetriever(hex.EncodeToString(pub), hex.EncodeToString(pv)))
	sharedService := shared.NewService(&mockSharedRepo{})

	sk, _ := shared.NewKeyPair(hex.EncodeToString([]byte("pvkey")), hex.EncodeToString(pub))
	sharedService.StoreSharedEmitterKeyPair(sk)

	miningSrv := transaction.NewMiningService(&mockPoolRequester{
		repo: txRepo,
	}, poolFindSrv, sharedService, "127.0.0.1", hex.EncodeToString(pub), hex.EncodeToString(pv))
	storeSrv := transaction.NewStorageService(txRepo, miningSrv)

	txSrv := NewTransactionServer(storeSrv, lockSrv, miningSrv, hex.EncodeToString(pub), hex.EncodeToString(pv))

	lis, _ := net.Listen("tcp", ":3545")
	defer lis.Close()
	grpcServer := grpc.NewServer()
	api.RegisterTransactionServiceServer(grpcServer, txSrv)
	go grpcServer.Serve(lis)

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
	keychain, _ := transaction.NewKeychain(tx)
	txRepo.StoreKeychain(keychain)

	ret := NewPoolRetriever(hex.EncodeToString(pub), hex.EncodeToString(pv))
	txRes, err := ret.RequestLastTransaction(transaction.Pool{
		transaction.NewPoolMember(net.ParseIP("127.0.0.1"), 3545),
	}, crypto.HashString("addr"), transaction.KeychainType)
	assert.Nil(t, err)
	assert.Equal(t, transaction.KeychainType, txRes.Type())
	assert.Equal(t, crypto.HashBytes(txSigned), txRes.TransactionHash())
}

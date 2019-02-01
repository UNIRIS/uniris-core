package rpc

import (
	"encoding/hex"
	"encoding/json"
	"net"
	"testing"
	"time"

	"github.com/uniris/uniris-core/pkg/shared"

	"github.com/stretchr/testify/assert"

	api "github.com/uniris/uniris-core/api/protobuf-spec"
	"github.com/uniris/uniris-core/pkg/crypto"
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

	data := map[string]string{
		"encrypted_address": hex.EncodeToString([]byte("addr")),
		"encrypted_wallet":  hex.EncodeToString([]byte("wallet")),
	}

	propKP, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("encPV")), pub)
	prop, _ := transaction.NewProposal(propKP)
	txRaw := map[string]interface{}{
		"address":    crypto.HashString("addr"),
		"data":       data,
		"timestamp":  time.Now().Unix(),
		"type":       transaction.KeychainType,
		"public_key": pub,
		"proposal":   prop,
	}
	txBytes, _ := json.Marshal(txRaw)
	sig, _ := crypto.Sign(string(txBytes), pv)
	txRaw["signature"] = sig
	txRaw["em_signature"] = sig
	txBytes, _ = json.Marshal(txRaw)

	tx, _ := transaction.New(crypto.HashString("addr"), transaction.KeychainType, data, time.Now(), pub, sig, sig, prop, crypto.HashBytes(txBytes))
	keychain, _ := transaction.NewKeychain(tx)
	txRepo.StoreKeychain(keychain)

	sharedSrv := shared.NewService(sharedRepo)
	ret := NewPoolRetriever(sharedSrv)
	txRes, err := ret.RequestLastTransaction(transaction.Pool{
		transaction.NewPoolMember(net.ParseIP("127.0.0.1"), 3545),
	}, crypto.HashString("addr"), transaction.KeychainType)
	assert.Nil(t, err)
	assert.Equal(t, transaction.KeychainType, txRes.Type())
	assert.Equal(t, crypto.HashBytes(txBytes), txRes.TransactionHash())
}

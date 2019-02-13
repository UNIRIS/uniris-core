package rpc

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"net"
	"testing"
	"time"

	"github.com/uniris/uniris-core/pkg/shared"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

	api "github.com/uniris/uniris-core/api/protobuf-spec"
	"github.com/uniris/uniris-core/pkg/chain"
	"github.com/uniris/uniris-core/pkg/crypto"
)

/*
Scenario: Receive a get transaction status from the API
	Given a transaction status request incoming from the API
	When I want to get the transaction status
	Then I get the right pool and send a GRPC request to the pool to get the status
*/
func TestHandleGetTransactionStatusInternal(t *testing.T) {

	pub, pv := crypto.GenerateKeys()

	chainDB := &mockChainDB{}
	techDB := &mockTechDB{}
	minerKey, _ := shared.NewMinerKeyPair(pub, pv)
	techDB.minerKeys = append(techDB.minerKeys, minerKey)

	pr := NewPoolRequester(techDB)

	chainSrv := NewChainServer(chainDB, techDB, pr)
	intSrv := NewInternalServer(techDB, pr)

	lis, _ := net.Listen("tcp", ":3545")
	defer lis.Close()
	grpcServer := grpc.NewServer()
	api.RegisterChainServiceServer(grpcServer, chainSrv)
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
	assert.Nil(t, crypto.VerifySignature(string(resBytes), pub, res.SignatureResponse))

}

/*
Scenario: Forward an transaction request from the API to the master miner
	Given a transaction come from the API
	When I want to process it
	Then I forward to the master miner and reply the transaction hash
*/
func TestHandleIncomingTransaction(t *testing.T) {
	pub, pv := crypto.GenerateKeys()

	chainDB := &mockChainDB{}
	techDB := &mockTechDB{}

	encKey, _ := crypto.Encrypt(pv, pub)
	emKey, _ := shared.NewEmitterKeyPair(encKey, pub)
	techDB.emKeys = append(techDB.emKeys, emKey)

	minerKey, _ := shared.NewMinerKeyPair(pub, pv)
	techDB.minerKeys = append(techDB.minerKeys, minerKey)

	pr := NewPoolRequester(techDB)

	chainSrv := NewChainServer(chainDB, techDB, pr)
	intSrv := NewInternalServer(techDB, pr)
	miningSrv := NewMiningServer(techDB, pr, pub, pv)

	lis, _ := net.Listen("tcp", ":3545")
	defer lis.Close()
	grpcServer := grpc.NewServer()
	api.RegisterChainServiceServer(grpcServer, chainSrv)
	api.RegisterMiningServiceServer(grpcServer, miningSrv)
	go grpcServer.Serve(lis)

	tx := map[string]interface{}{
		"addr": crypto.HashString("addr"),
		"data": map[string]string{
			"encrypted_address": hex.EncodeToString([]byte("addr")),
			"encrypted_wallet":  hex.EncodeToString([]byte("wallet")),
		},
		"timestamp":  time.Now().Unix(),
		"type":       int(chain.KeychainTransactionType),
		"public_key": pub,
		"em_shared_keys_proposal": map[string]string{
			"encrypted_private_key": hex.EncodeToString([]byte("encPV")),
			"public_key":            pub,
		},
	}
	txBytes, _ := json.Marshal(tx)
	sig, _ := crypto.Sign(string(txBytes), pv)
	tx["signature"] = sig
	tx["em_signature"] = sig
	txBytes, _ = json.Marshal(tx)

	cipherTx, _ := crypto.Encrypt(string(txBytes), pub)
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
	assert.Nil(t, crypto.VerifySignature(string(resBytes), pub, res.Signature))
	assert.Equal(t, crypto.HashBytes(txBytes), res.TransactionHash)

	time.Sleep(2 * time.Second)

	assert.Len(t, chainDB.keychains, 1)
	assert.Equal(t, res.TransactionHash, chainDB.keychains[0].TransactionHash())

}

/*
Scenario: Get account created
	Given a account created (id and keychain transaction proceed)
	When I want to retreive the account
	Then I provide the ID hash and I get retrieve the encrypted wallet and the encrypted aes key
*/
func TestHandleGetAccount(t *testing.T) {
	pub, pv := crypto.GenerateKeys()

	chainDB := &mockChainDB{}
	techDB := &mockTechDB{}
	minerKey, _ := shared.NewMinerKeyPair(pub, pv)
	techDB.minerKeys = append(techDB.minerKeys, minerKey)
	emKey, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("pv")), pub)
	techDB.emKeys = append(techDB.emKeys, emKey)

	pr := NewPoolRequester(techDB)

	chainSrv := NewChainServer(chainDB, techDB, pr)
	miningSrv := NewMiningServer(techDB, pr, pub, pv)
	intSrv := NewInternalServer(techDB, pr)

	lis, _ := net.Listen("tcp", ":3545")
	defer lis.Close()
	grpcServer := grpc.NewServer()
	api.RegisterChainServiceServer(grpcServer, chainSrv)
	api.RegisterMiningServiceServer(grpcServer, miningSrv)
	go grpcServer.Serve(lis)

	encAddr, _ := crypto.Encrypt(crypto.HashString("addr"), pub)

	//First send the ID transaction
	txID := map[string]interface{}{
		"addr": crypto.HashString("idHash"),
		"data": map[string]string{
			"encrypted_address_by_id":    encAddr,
			"encrypted_address_by_robot": encAddr,
			"encrypted_aes_key":          hex.EncodeToString([]byte("aesKey")),
		},
		"timestamp":  time.Now().Unix(),
		"type":       int(chain.IDTransactionType),
		"public_key": pub,
		"em_shared_keys_proposal": map[string]string{
			"encrypted_private_key": hex.EncodeToString([]byte("encPV")),
			"public_key":            pub,
		},
	}
	txIDBytes, _ := json.Marshal(txID)
	sigID, _ := crypto.Sign(string(txIDBytes), pv)
	txID["signature"] = sigID
	txID["em_signature"] = sigID
	txIDBytes, _ = json.Marshal(txID)

	cipherTx, _ := crypto.Encrypt(string(txIDBytes), pub)
	res, err := intSrv.HandleTransaction(context.TODO(), &api.IncomingTransaction{
		EncryptedTransaction: cipherTx,
		Timestamp:            time.Now().Unix(),
	})
	assert.Nil(t, err)
	assert.NotEmpty(t, res.TransactionHash)

	time.Sleep(1 * time.Second)
	assert.Equal(t, crypto.HashString("idHash"), chainDB.ids[0].Address())

	//Then send the keychain transaction
	txKeychain := map[string]interface{}{
		"addr": crypto.HashString("addr"),
		"data": map[string]string{
			"encrypted_address": hex.EncodeToString([]byte("addr")),
			"encrypted_wallet":  hex.EncodeToString([]byte("wallet")),
		},
		"timestamp":  time.Now().Unix(),
		"type":       int(chain.KeychainTransactionType),
		"public_key": pub,
		"em_shared_keys_proposal": map[string]string{
			"encrypted_private_key": hex.EncodeToString([]byte("encPV")),
			"public_key":            pub,
		},
	}
	txKeychainBytes, _ := json.Marshal(txKeychain)
	sigKeychain, _ := crypto.Sign(string(txKeychainBytes), pv)
	txKeychain["signature"] = sigKeychain
	txKeychain["em_signature"] = sigKeychain
	txKeychainBytes, _ = json.Marshal(txKeychain)

	cipherTx2, _ := crypto.Encrypt(string(txKeychainBytes), pub)
	res2, err := intSrv.HandleTransaction(context.TODO(), &api.IncomingTransaction{
		EncryptedTransaction: cipherTx2,
		Timestamp:            time.Now().Unix(),
	})
	assert.Nil(t, err)
	assert.NotEmpty(t, res2.TransactionHash)

	time.Sleep(1 * time.Second)
	assert.Equal(t, crypto.HashString("addr"), chainDB.keychains[0].Address())

	encIDHash, _ := crypto.Encrypt(crypto.HashString("idHash"), pub)
	req := &api.GetAccountRequest{
		EncryptedIdAddress: encIDHash,
		Timestamp:          time.Now().Unix(),
	}

	resGet, err := intSrv.GetAccount(context.TODO(), req)
	assert.Nil(t, err)
	assert.NotEmpty(t, resGet.EncryptedAesKey)
	assert.NotEmpty(t, resGet.EncryptedWallet)
}

/*
Scenario: Get last shared keys
	Given a emitter authorized and stored shared keys
	When I want to get the lastshared keys
	Then I get the last shared keys with encryption
*/
func TestGetLastSharedKeys(t *testing.T) {
	pub, pv := crypto.GenerateKeys()

	techDB := &mockTechDB{}
	minerKey, _ := shared.NewMinerKeyPair(pub, pv)
	techDB.minerKeys = append(techDB.minerKeys, minerKey)

	encKey, _ := crypto.Encrypt(pv, pub)
	emKey, _ := shared.NewEmitterKeyPair(encKey, pub)

	techDB.minerKeys = append(techDB.minerKeys, minerKey)
	techDB.emKeys = append(techDB.emKeys, emKey)

	pr := NewPoolRequester(techDB)

	intSrv := NewInternalServer(techDB, pr)

	res, err := intSrv.GetLastSharedKeys(context.TODO(), &api.LastSharedKeysRequest{
		EmitterPublicKey: pub,
		Timestamp:        time.Now().Unix(),
	})
	assert.Nil(t, err)
	assert.Equal(t, pub, res.MinerPublicKey)
	assert.Len(t, res.EmitterKeys, 1)
	assert.Equal(t, pub, res.EmitterKeys[0].PublicKey)

	emPv, _ := crypto.Decrypt(res.EmitterKeys[0].EncryptedPrivateKey, pv)
	assert.Equal(t, pv, emPv)
}

type mockTechDB struct {
	emKeys    shared.EmitterKeys
	minerKeys []shared.MinerKeyPair
}

func (db mockTechDB) EmitterKeys() (shared.EmitterKeys, error) {
	return db.emKeys, nil
}

func (db mockTechDB) LastMinerKeys() (shared.MinerKeyPair, error) {
	return db.minerKeys[len(db.minerKeys)-1], nil
}

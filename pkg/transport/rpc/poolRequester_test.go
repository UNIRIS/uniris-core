package rpc

import (
	"crypto/rand"
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

	pv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

	sharedKeyReader := &mockSharedKeyReader{}
	nodeKey, _ := shared.NewNodeCrossKeyPair(pub, pv)
	sharedKeyReader.crossNodeKeys = append(sharedKeyReader.crossNodeKeys, nodeKey)

	nodeReader := &mockNodeReader{
		nodes: []consensus.Node{
			consensus.NewNode(net.ParseIP("127.0.0.1"), 5000, pub, consensus.NodeOK, "", 300, "1.0", 0, 1, 30.0, -10.0, consensus.GeoPatch{}, true),
		},
	}

	chainDB := &mockChainDB{}
	pr := NewPoolRequester(sharedKeyReader)
	txSrv := NewTransactionService(chainDB, sharedKeyReader, nodeReader, pr, pub, pv)

	lis, _ := net.Listen("tcp", ":5000")
	defer lis.Close()
	grpcServer := grpc.NewServer()
	api.RegisterTransactionServiceServer(grpcServer, txSrv)
	go grpcServer.Serve(lis)

	pool, _ := consensus.FindStoragePool([]byte("addr"), nodeReader)
	assert.Nil(t, pr.RequestTransactionTimeLock(pool, crypto.Hash([]byte("tx")), crypto.Hash([]byte("addr")), pub))
	assert.True(t, chain.ContainsTimeLock(crypto.Hash([]byte("tx")), crypto.Hash([]byte("addr"))))
}

/*
Scenario: Request transaction validation confirmation
	Given a transaction to validate
	When I request to confirm the validation
	Then I get a validation
*/
func TestRequestConfirmValidation(t *testing.T) {

	pv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
	pubB, _ := pub.Marshal()

	kp, _ := shared.NewEmitterCrossKeyPair(([]byte("pvkey")), pub)

	sharedKeyReader := &mockSharedKeyReader{}
	sharedKeyReader.crossEmitterKeys = append(sharedKeyReader.crossEmitterKeys, kp)
	nodeKey, _ := shared.NewNodeCrossKeyPair(pub, pv)
	sharedKeyReader.crossNodeKeys = append(sharedKeyReader.crossNodeKeys, nodeKey)
	sharedKeyReader.authKeys = append(sharedKeyReader.authKeys, pub)

	nodeReader := &mockNodeReader{
		nodes: []consensus.Node{
			consensus.NewNode(net.ParseIP("127.0.0.1"), 5000, pub, consensus.NodeOK, "", 300, "1.0", 0, 1, 30.0, -10.0, consensus.GeoPatch{}, true),
		},
	}

	pr := NewPoolRequester(sharedKeyReader)

	miningSrv := NewTransactionService(nil, sharedKeyReader, nodeReader, pr, pub, pv)

	lis, err := net.Listen("tcp", ":5000")
	defer lis.Close()
	grpcServer := grpc.NewServer()
	api.RegisterTransactionServiceServer(grpcServer, miningSrv)
	go grpcServer.Serve(lis)

	data := map[string][]byte{
		"encrypted_address_by_node": []byte("addr"),
		"encrypted_wallet":          []byte("wallet"),
	}
	prop := kp
	txRaw := map[string]interface{}{
		"addr": hex.EncodeToString(crypto.Hash([]byte("addr"))),
		"data": map[string]string{
			"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
			"encrypted_wallet":          hex.EncodeToString([]byte("wallet")),
		},
		"timestamp":  time.Now().Unix(),
		"type":       chain.KeychainTransactionType,
		"public_key": hex.EncodeToString(pubB),
		"em_shared_keys_proposal": map[string]string{
			"encrypted_private_key": hex.EncodeToString([]byte("pvkey")),
			"public_key":            hex.EncodeToString(pubB),
		},
	}
	txBytes, _ := json.Marshal(txRaw)
	sig, _ := pv.Sign(txBytes)
	txRaw["signature"] = hex.EncodeToString(sig)
	txByteWithSig, _ := json.Marshal(txRaw)
	emSig, _ := pv.Sign(txByteWithSig)
	txRaw["em_signature"] = hex.EncodeToString(emSig)
	txBytes, _ = json.Marshal(txRaw)

	tx, _ := chain.NewTransaction(crypto.Hash([]byte("addr")), chain.KeychainTransactionType, data, time.Now(), pub, prop, sig, emSig, crypto.Hash(txBytes))

	vBytes, _ := json.Marshal(map[string]interface{}{
		"status":     chain.ValidationOK,
		"public_key": pubB,
		"timestamp":  time.Now().Unix(),
	})
	vSig, _ := pv.Sign(vBytes)
	v, err := chain.NewValidation(chain.ValidationOK, time.Now(), pub, vSig)

	wHeaders := chain.NewWelcomeNodeHeader(pub, []chain.NodeHeader{chain.NewNodeHeader(pub, false, false, 0, true)}, []byte("sig"))
	vHeaders := []chain.NodeHeader{chain.NewNodeHeader(pub, false, false, 0, true)}
	sHeaders := []chain.NodeHeader{chain.NewNodeHeader(pub, false, false, 0, true)}
	mv, _ := chain.NewMasterValidation([]crypto.PublicKey{}, pub, v, wHeaders, vHeaders, sHeaders)

	pool, _ := consensus.FindValidationPool(tx.Address(), 1, pub, nodeReader, sharedKeyReader)
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

	pv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
	kp, _ := shared.NewEmitterCrossKeyPair([]byte("pvkey"), pub)
	pubB, _ := pub.Marshal()

	chainDB := &mockChainDB{}
	sharedKeyReader := &mockSharedKeyReader{}
	nodeKey, _ := shared.NewNodeCrossKeyPair(pub, pv)
	sharedKeyReader.crossNodeKeys = append(sharedKeyReader.crossNodeKeys, nodeKey)

	nodeReader := &mockNodeReader{
		nodes: []consensus.Node{
			consensus.NewNode(net.ParseIP("127.0.0.1"), 5000, pub, consensus.NodeOK, "", 300, "1.0", 0, 1, 30.0, -10.0, consensus.GeoPatch{}, true),
		},
	}

	sharedKeyReader.crossEmitterKeys = append(sharedKeyReader.crossEmitterKeys, kp)

	pr := NewPoolRequester(sharedKeyReader)

	txSrv := NewTransactionService(chainDB, sharedKeyReader, nodeReader, pr, pub, pv)

	lis, _ := net.Listen("tcp", ":5000")
	defer lis.Close()
	grpcServer := grpc.NewServer()
	api.RegisterTransactionServiceServer(grpcServer, txSrv)
	go grpcServer.Serve(lis)

	data := map[string][]byte{
		"encrypted_address_by_node": []byte("addr"),
		"encrypted_wallet":          []byte("wallet"),
	}
	prop := kp
	txRaw := map[string]interface{}{
		"addr": hex.EncodeToString(crypto.Hash([]byte("addr"))),
		"data": map[string]string{
			"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
			"encrypted_wallet":          hex.EncodeToString([]byte("wallet")),
		},
		"timestamp":  time.Now().Unix(),
		"type":       chain.KeychainTransactionType,
		"public_key": hex.EncodeToString(pubB),
		"em_shared_keys_proposal": map[string]string{
			"encrypted_private_key": hex.EncodeToString([]byte("pvkey")),
			"public_key":            hex.EncodeToString(pubB),
		},
	}
	txBytes, _ := json.Marshal(txRaw)
	sig, _ := pv.Sign(txBytes)
	txRaw["signature"] = hex.EncodeToString(sig)
	txByteWithSig, _ := json.Marshal(txRaw)
	emSig, _ := pv.Sign(txByteWithSig)
	txRaw["em_signature"] = hex.EncodeToString(emSig)
	txBytes, _ = json.Marshal(txRaw)

	tx, _ := chain.NewTransaction(crypto.Hash([]byte("addr")), chain.KeychainTransactionType, data, time.Now(), pub, prop, sig, emSig, crypto.Hash(txBytes))
	vBytes, _ := json.Marshal(map[string]interface{}{
		"status":     chain.ValidationOK,
		"public_key": pubB,
		"timestamp":  time.Now().Unix(),
	})
	vSig, _ := pv.Sign(vBytes)
	v, _ := chain.NewValidation(chain.ValidationOK, time.Now(), pub, vSig)
	wHeaders := chain.NewWelcomeNodeHeader(pub, []chain.NodeHeader{chain.NewNodeHeader(pub, false, false, 0, true)}, []byte("sig"))
	vHeaders := []chain.NodeHeader{chain.NewNodeHeader(pub, false, false, 0, true)}
	sHeaders := []chain.NodeHeader{chain.NewNodeHeader(pub, false, false, 0, true)}
	mv, _ := chain.NewMasterValidation([]crypto.PublicKey{}, pub, v, wHeaders, vHeaders, sHeaders)
	pool, _ := consensus.FindStoragePool([]byte("addr"), nodeReader)
	tx.Mined(mv, []chain.Validation{v})
	assert.Nil(t, pr.RequestTransactionStorage(pool, 1, tx))

	assert.Len(t, chainDB.keychains, 1)
	assert.EqualValues(t, crypto.Hash(txBytes), chainDB.keychains[0].TransactionHash())
}

/*
Scenario: Send request to get last transaction
	Given a keychain transaction stored
	When I want to request a node to get the last transaction from the address
	Then I get the last transaction
*/
func TestSendGetLastTransaction(t *testing.T) {

	pv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
	pubB, _ := pub.Marshal()

	chainDB := &mockChainDB{}
	sharedKeyReader := &mockSharedKeyReader{}
	nodeKey, _ := shared.NewNodeCrossKeyPair(pub, pv)
	sharedKeyReader.crossNodeKeys = append(sharedKeyReader.crossNodeKeys, nodeKey)

	nodeReader := &mockNodeReader{
		nodes: []consensus.Node{
			consensus.NewNode(net.ParseIP("127.0.0.1"), 5000, pub, consensus.NodeOK, "", 300, "1.0", 0, 1, 30.0, -10.0, consensus.GeoPatch{}, true),
		},
	}

	pr := NewPoolRequester(sharedKeyReader)

	txSrv := NewTransactionService(chainDB, sharedKeyReader, nodeReader, pr, pub, pv)

	lis, _ := net.Listen("tcp", ":5000")
	defer lis.Close()
	grpcServer := grpc.NewServer()
	api.RegisterTransactionServiceServer(grpcServer, txSrv)
	go grpcServer.Serve(lis)

	data := map[string][]byte{
		"encrypted_address_by_node": []byte("addr"),
		"encrypted_wallet":          []byte("wallet"),
	}

	prop, _ := shared.NewEmitterCrossKeyPair([]byte("encPV"), pub)
	txRaw := map[string]interface{}{
		"addr": hex.EncodeToString(crypto.Hash([]byte("addr"))),
		"data": map[string]string{
			"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
			"encrypted_wallet":          hex.EncodeToString([]byte("wallet")),
		},
		"timestamp":  time.Now().Unix(),
		"type":       chain.KeychainTransactionType,
		"public_key": hex.EncodeToString(pubB),
		"em_shared_keys_proposal": map[string]string{
			"encrypted_private_key": hex.EncodeToString([]byte("pvkey")),
			"public_key":            hex.EncodeToString(pubB),
		},
	}
	txBytes, _ := json.Marshal(txRaw)
	sig, _ := pv.Sign(txBytes)
	txRaw["signature"] = hex.EncodeToString(sig)
	txByteWithSig, _ := json.Marshal(txRaw)
	emSig, _ := pv.Sign(txByteWithSig)
	txRaw["em_signature"] = hex.EncodeToString(emSig)
	txBytes, _ = json.Marshal(txRaw)

	tx, _ := chain.NewTransaction(crypto.Hash([]byte("addr")), chain.KeychainTransactionType, data, time.Now(), pub, prop, sig, sig, crypto.Hash(txBytes))
	keychain, _ := chain.NewKeychain(tx)
	chainDB.keychains = append(chainDB.keychains, keychain)

	pool, _ := consensus.FindStoragePool([]byte("address"), nodeReader)

	txRes, err := pr.RequestLastTransaction(pool, crypto.Hash([]byte("addr")), chain.KeychainTransactionType)
	assert.Nil(t, err)
	assert.Equal(t, chain.KeychainTransactionType, txRes.TransactionType())
	assert.Equal(t, crypto.Hash(txBytes), txRes.TransactionHash())
}

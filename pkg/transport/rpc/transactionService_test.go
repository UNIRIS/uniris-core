package rpc

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log"
	"net"
	"os"
	"sort"
	"testing"
	"time"

	"github.com/uniris/uniris-core/pkg/consensus"

	"github.com/uniris/uniris-core/pkg/chain"
	"google.golang.org/grpc/codes"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/status"

	api "github.com/uniris/uniris-core/api/protobuf-spec"
	"github.com/uniris/uniris-core/pkg/logging"
	"github.com/uniris/uniris-core/pkg/shared"

	"github.com/uniris/uniris-core/pkg/crypto"
)

/*
Scenario: Receive  get last transction about an unknown transaction
	Given no transaction store for an address
	When I want to request to retrieve the last transaction keychain of this unknown address
	Then I get an error
*/
func TestHandleGetLastTransactionWhenNotExist(t *testing.T) {
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

	poolR := &mockPoolRequester{
		repo: chainDB,
	}

	l := logging.NewLogger("stdout", log.New(os.Stdout, "", 0), "test", net.ParseIP("127.0.0.1"), logging.ErrorLogLevel)
	txSrv := NewTransactionService(chainDB, sharedKeyReader, nodeReader, poolR, pub, pv, l)

	req := &api.GetLastTransactionRequest{
		Timestamp:          time.Now().Unix(),
		TransactionAddress: crypto.Hash([]byte("address")),
		Type:               api.TransactionType_KEYCHAIN,
	}
	reqBytes, _ := json.Marshal(req)
	sig, _ := pv.Sign(reqBytes)
	req.SignatureRequest = sig

	_, err := txSrv.GetLastTransaction(context.TODO(), req)
	assert.NotNil(t, err)
	statusCode, _ := status.FromError(err)
	assert.Equal(t, codes.NotFound, statusCode.Code())
	assert.Equal(t, statusCode.Message(), "transaction does not exist")
}

/*
Scenario: Receive  get last transaction request
	Given a keychain transaction stored
	When I want to request to retrieve the last transaction keychain of this address
	Then I get an error
*/
func TestHandleGetLastTransaction(t *testing.T) {

	pv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
	pubB, _ := pub.Marshal()

	sharedKeyReader := &mockSharedKeyReader{}
	nodeKey, _ := shared.NewNodeCrossKeyPair(pub, pv)
	sharedKeyReader.crossNodeKeys = append(sharedKeyReader.crossNodeKeys, nodeKey)

	nodeReader := &mockNodeReader{
		nodes: []consensus.Node{
			consensus.NewNode(net.ParseIP("127.0.0.1"), 5000, pub, consensus.NodeOK, "", 300, "1.0", 0, 1, 30.0, -10.0, consensus.GeoPatch{}, true),
		},
	}

	chainDB := &mockChainDB{}

	poolR := &mockPoolRequester{
		repo: chainDB,
	}

	l := logging.NewLogger("stdout", log.New(os.Stdout, "", 0), "test", net.ParseIP("127.0.0.1"), logging.ErrorLogLevel)
	txSrv := NewTransactionService(chainDB, sharedKeyReader, nodeReader, poolR, pub, pv, l)

	data := map[string][]byte{
		"encrypted_address_by_node": []byte("addr"),
		"encrypted_wallet":          []byte("wallet"),
	}

	prop, _ := shared.NewEmitterCrossKeyPair([]byte("pvkey"), pub)

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

	txByteWithSign, _ := json.Marshal(txRaw)
	emSig, _ := pv.Sign(txByteWithSign)
	txRaw["em_signature"] = hex.EncodeToString(emSig)
	txBytes, _ = json.Marshal(txRaw)

	tx, _ := chain.NewTransaction(crypto.Hash([]byte("addr")), chain.KeychainTransactionType, data, time.Now(), pub, prop, sig, emSig, crypto.Hash([]byte(txBytes)))
	keychain, _ := chain.NewKeychain(tx)
	chainDB.keychains = append(chainDB.keychains, keychain)

	req := &api.GetLastTransactionRequest{
		Timestamp:          time.Now().Unix(),
		TransactionAddress: crypto.Hash([]byte("addr")),
		Type:               api.TransactionType_KEYCHAIN,
	}
	reqBytes, _ := json.Marshal(req)
	sigReq, _ := pv.Sign(reqBytes)
	req.SignatureRequest = sigReq

	res, err := txSrv.GetLastTransaction(context.TODO(), req)
	assert.Nil(t, err)
	assert.NotEmpty(t, res.SignatureResponse)
	assert.NotNil(t, res.Transaction)
	assert.EqualValues(t, crypto.Hash(txBytes), res.Transaction.TransactionHash)

	resBytes, _ := json.Marshal(&api.GetLastTransactionResponse{
		Timestamp:   res.Timestamp,
		Transaction: res.Transaction,
	})
	assert.True(t, pub.Verify(resBytes, res.SignatureResponse))
}

/*
Scenario: Receive get transaction status request
	Given no transaction stored
	When I want to request the transactions status for this transaction hash
	Then I get a status unknown
*/
func TestHandleGetTransactionStatus(t *testing.T) {

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

	poolR := &mockPoolRequester{
		repo: chainDB,
	}

	l := logging.NewLogger("stdout", log.New(os.Stdout, "", 0), "test", net.ParseIP("127.0.0.1"), logging.ErrorLogLevel)
	txSrv := NewTransactionService(chainDB, sharedKeyReader, nodeReader, poolR, pub, pv, l)

	req := &api.GetTransactionStatusRequest{
		Timestamp:       time.Now().Unix(),
		TransactionHash: crypto.Hash([]byte("tx1")),
	}
	reqBytes, _ := json.Marshal(req)
	sig, _ := pv.Sign(reqBytes)
	req.SignatureRequest = sig

	res, err := txSrv.GetTransactionStatus(context.TODO(), req)
	assert.Nil(t, err)
	assert.Equal(t, api.TransactionStatus_UNKNOWN, res.Status)
	resBytes, _ := json.Marshal(&api.GetTransactionStatusResponse{
		Timestamp: res.Timestamp,
		Status:    res.Status,
	})
	assert.True(t, pub.Verify(resBytes, res.SignatureResponse))
}

/*
Scenario: Receive storage  transaction request
	Given a transaction
	When I want to request to store of the transaction
	Then the transaction is stored
*/
func TestHandleStoreTransaction(t *testing.T) {

	pv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
	pubB, _ := pub.Marshal()

	sharedKeyReader := &mockSharedKeyReader{}
	nodeKey, _ := shared.NewNodeCrossKeyPair(pub, pv)
	sharedKeyReader.crossNodeKeys = append(sharedKeyReader.crossNodeKeys, nodeKey)

	nodeReader := &mockNodeReader{
		nodes: []consensus.Node{
			consensus.NewNode(net.ParseIP("127.0.0.1"), 5000, pub, consensus.NodeOK, "", 300, "1.0", 0, 1, 30.0, -10.0, consensus.GeoPatch{}, true),
		},
	}

	chainDB := &mockChainDB{}

	poolR := &mockPoolRequester{
		repo: chainDB,
	}

	l := logging.NewLogger("stdout", log.New(os.Stdout, "", 0), "test", net.ParseIP("127.0.0.1"), logging.ErrorLogLevel)
	txSrv := NewTransactionService(chainDB, sharedKeyReader, nodeReader, poolR, pub, pv, l)

	prop, _ := shared.NewEmitterCrossKeyPair([]byte("pvkey"), pub)

	data := map[string][]byte{
		"encrypted_address_by_node": []byte("addr"),
		"encrypted_wallet":          []byte("wallet"),
	}
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

	txByteWithSign, _ := json.Marshal(txRaw)
	emSig, _ := pv.Sign(txByteWithSign)
	txRaw["em_signature"] = hex.EncodeToString(emSig)
	txBytes, _ = json.Marshal(txRaw)

	txBytes, _ = json.Marshal(txRaw)

	tx, _ := chain.NewTransaction(crypto.Hash([]byte("addr")), chain.KeychainTransactionType, data, time.Now(), pub, prop, sig, emSig, crypto.Hash(txBytes))

	vRaw := map[string]interface{}{
		"status":     chain.ValidationOK,
		"public_key": pubB,
		"timestamp":  time.Now().Unix(),
	}
	vBytes, _ := json.Marshal(vRaw)
	vSig, _ := pv.Sign(vBytes)
	v, _ := chain.NewValidation(chain.ValidationOK, time.Now(), pub, vSig)
	wHeaders := chain.NewWelcomeNodeHeader(pub, []chain.NodeHeader{chain.NewNodeHeader(pub, false, false, 0, true)}, []byte("sig"))
	vHeaders := []chain.NodeHeader{chain.NewNodeHeader(pub, false, false, 0, true)}
	sHeaders := []chain.NodeHeader{chain.NewNodeHeader(pub, false, false, 0, true)}
	mv, _ := chain.NewMasterValidation([]crypto.PublicKey{}, pub, v, wHeaders, vHeaders, sHeaders)

	txf, _ := formatAPITransaction(tx)
	mvf, _ := formatAPIMasterValidation(mv)
	vf, _ := formatAPIValidation(v)

	req := &api.StoreTransactionRequest{
		Timestamp: time.Now().Unix(),
		MinedTransaction: &api.MinedTransaction{
			Transaction:        txf,
			MasterValidation:   mvf,
			ConfirmValidations: []*api.Validation{vf},
		},
	}

	reqBytes, _ := json.Marshal(req)
	sigReq, _ := pv.Sign(reqBytes)
	req.SignatureRequest = sigReq

	res, err := txSrv.StoreTransaction(context.TODO(), req)
	assert.Nil(t, err)

	resBytes, _ := json.Marshal(&api.StoreTransactionResponse{
		Timestamp: res.Timestamp,
	})
	assert.True(t, pub.Verify(resBytes, res.SignatureResponse))

	assert.Len(t, chainDB.keychains, 1)
	assert.EqualValues(t, crypto.Hash(txBytes), chainDB.keychains[0].TransactionHash())

}

/*
Scenario: Receive lock transaction request
	Given a transaction to lock
	When I want to request to lock it
	Then I get not error and the lock is stored
*/
func TestHandleLockTransaction(t *testing.T) {

	pv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

	chainDB := &mockChainDB{}
	sharedKeyReader := &mockSharedKeyReader{}
	nodeKey, _ := shared.NewNodeCrossKeyPair(pub, pv)
	sharedKeyReader.crossNodeKeys = append(sharedKeyReader.crossNodeKeys, nodeKey)

	nodeReader := &mockNodeReader{
		nodes: []consensus.Node{
			consensus.NewNode(net.ParseIP("127.0.0.1"), 5000, pub, consensus.NodeOK, "", 300, "1.0", 0, 1, 30.0, -10.0, consensus.GeoPatch{}, true),
		},
	}

	poolR := &mockPoolRequester{}
	l := logging.NewLogger("stdout", log.New(os.Stdout, "", 0), "test", net.ParseIP("127.0.0.1"), logging.ErrorLogLevel)
	txSrv := NewTransactionService(chainDB, sharedKeyReader, nodeReader, poolR, pub, pv, l)

	pubB, _ := pub.Marshal()

	req := &api.TimeLockTransactionRequest{
		Timestamp:           time.Now().Unix(),
		TransactionHash:     crypto.Hash([]byte("tx1")),
		MasterNodePublicKey: pubB,
		Address:             crypto.Hash([]byte("addr1")),
	}
	reqBytes, _ := json.Marshal(req)
	sig, _ := pv.Sign(reqBytes)
	req.SignatureRequest = sig

	res, err := txSrv.TimeLockTransaction(context.TODO(), req)
	assert.Nil(t, err)
	resBytes, _ := json.Marshal(&api.TimeLockTransactionResponse{
		Timestamp: res.Timestamp,
	})
	assert.True(t, pub.Verify(resBytes, res.SignatureResponse))
	assert.True(t, chain.ContainsTimeLock(crypto.Hash([]byte("tx1")), crypto.Hash([]byte("addr1"))))
}

/*
Scenario: Receive lead mining transaction request
	Given a transaction to validate
	When I want to request to lead mining of the transaction
	Then I get not error
*/
func TestHandleLeadTransactionMining(t *testing.T) {

	pv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
	pubB, _ := pub.Marshal()

	sharedKeyReader := &mockSharedKeyReader{}
	nodeKey, _ := shared.NewNodeCrossKeyPair(pub, pv)
	sharedKeyReader.crossNodeKeys = append(sharedKeyReader.crossNodeKeys, nodeKey)
	emKey, _ := shared.NewEmitterCrossKeyPair([]byte("encpv"), pub)
	sharedKeyReader.crossEmitterKeys = append(sharedKeyReader.crossEmitterKeys, emKey)
	sharedKeyReader.authKeys = append(sharedKeyReader.authKeys, pub)

	nodeReader := &mockNodeReader{
		nodes: []consensus.Node{
			consensus.NewNode(net.ParseIP("127.0.0.1"), 5000, pub, consensus.NodeOK, "", 300, "1.0", 0, 1, 30.0, -10.0, consensus.GeoPatch{}, true),
		},
	}

	chainDB := &mockChainDB{}

	poolR := &mockPoolRequester{
		repo: chainDB,
	}
	l := logging.NewLogger("stdout", log.New(os.Stdout, "", 0), "test", net.ParseIP("127.0.0.1"), logging.ErrorLogLevel)
	txSrv := NewTransactionService(chainDB, sharedKeyReader, nodeReader, poolR, pub, pv, l)
	data := map[string][]byte{
		"encrypted_address_by_node": []byte("addr"),
		"encrypted_wallet":          []byte("wallet"),
	}

	prop, _ := shared.NewEmitterCrossKeyPair([]byte("pvkey"), pub)
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
	txf, _ := formatAPITransaction(tx)
	ml := []*api.NodeHeader{
		&api.NodeHeader{
			IsMaster:      true,
			IsUnreachable: false,
			PatchNumber:   1,
			PublicKey:     pubB,
		}}
	req := &api.LeadTransactionMiningRequest{
		Timestamp:          time.Now().Unix(),
		MinimumValidations: 1,
		WelcomeHeaders: &api.WelcomeNodeHeader{
			PublicKey:   pubB,
			MastersList: ml,
			Signature:   []byte("sig"),
		},
		Transaction: txf,
	}

	reqBytes, _ := json.Marshal(req)
	sigReq, _ := pv.Sign(reqBytes)
	req.SignatureRequest = sigReq

	res, err := txSrv.LeadTransactionMining(context.TODO(), req)
	assert.Nil(t, err)

	time.Sleep(1 * time.Second)

	resBytes, _ := json.Marshal(&api.LeadTransactionMiningResponse{
		Timestamp: res.Timestamp,
	})
	assert.True(t, pub.Verify(resBytes, res.SignatureResponse))

	assert.Len(t, chainDB.keychains, 1)
	assert.EqualValues(t, crypto.Hash([]byte("addr")), chainDB.keychains[0].Address())
}

/*
Scenario: Receive confirmation of validations transaction request
	Given a transaction to validate
	When I want to request to validation of the transaction
	Then I get the node validation
*/
func TestHandleConfirmValiation(t *testing.T) {

	pv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
	pubB, _ := pub.Marshal()

	sharedKeyReader := &mockSharedKeyReader{}
	nodeKey, _ := shared.NewNodeCrossKeyPair(pub, pv)
	sharedKeyReader.crossNodeKeys = append(sharedKeyReader.crossNodeKeys, nodeKey)
	emKey, _ := shared.NewEmitterCrossKeyPair([]byte("encpv"), pub)
	sharedKeyReader.crossEmitterKeys = append(sharedKeyReader.crossEmitterKeys, emKey)

	nodeReader := &mockNodeReader{
		nodes: []consensus.Node{
			consensus.NewNode(net.ParseIP("127.0.0.1"), 5000, pub, consensus.NodeOK, "", 300, "1.0", 0, 1, 30.0, -10.0, consensus.GeoPatch{}, true),
		},
	}

	chainDB := &mockChainDB{}

	poolR := &mockPoolRequester{
		repo: chainDB,
	}
	l := logging.NewLogger("stdout", log.New(os.Stdout, "", 0), "test", net.ParseIP("127.0.0.1"), logging.ErrorLogLevel)
	txSrv := NewTransactionService(chainDB, sharedKeyReader, nodeReader, poolR, pub, pv, l)

	prop, _ := shared.NewEmitterCrossKeyPair([]byte("pvkey"), pub)
	data := map[string][]byte{
		"encrypted_address_by_node": []byte("addr"),
		"encrypted_wallet":          []byte("wallet"),
	}

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

	vRaw := map[string]interface{}{
		"status":     chain.ValidationOK,
		"public_key": pubB,
		"timestamp":  time.Now().Unix(),
	}

	vBytes, _ := json.Marshal(vRaw)
	vSig, _ := pv.Sign(vBytes)
	v, _ := chain.NewValidation(chain.ValidationOK, time.Now(), pub, vSig)
	wHeaders := chain.NewWelcomeNodeHeader(pub, []chain.NodeHeader{chain.NewNodeHeader(pub, false, false, 0, true)}, []byte("sig"))
	vHeaders := []chain.NodeHeader{chain.NewNodeHeader(pub, false, false, 0, true)}
	sHeaders := []chain.NodeHeader{chain.NewNodeHeader(pub, false, false, 0, true)}
	mv, _ := chain.NewMasterValidation([]crypto.PublicKey{}, pub, v, wHeaders, vHeaders, sHeaders)

	txv, _ := formatAPITransaction(tx)
	mvf, _ := formatAPIMasterValidation(mv)

	req := &api.ConfirmTransactionValidationRequest{
		Transaction:      txv,
		Timestamp:        time.Now().Unix(),
		MasterValidation: mvf,
	}

	reqBytes, _ := json.Marshal(req)
	sigReq, _ := pv.Sign(reqBytes)
	req.SignatureRequest = sigReq

	res, err := txSrv.ConfirmTransactionValidation(context.TODO(), req)
	assert.Nil(t, err)

	resBytes, _ := json.Marshal(&api.ConfirmTransactionValidationResponse{
		Timestamp:  res.Timestamp,
		Validation: res.Validation,
	})
	assert.True(t, pub.Verify(resBytes, res.SignatureResponse))

	assert.NotNil(t, res.Validation)
	assert.Equal(t, api.Validation_OK, res.Validation.Status)
	assert.EqualValues(t, pubB, res.Validation.PublicKey)
}

type mockPoolRequester struct {
	stores []chain.Transaction
	repo   *mockChainDB
}

func (pr mockPoolRequester) RequestLastTransaction(pool consensus.Pool, txAddr crypto.VersionnedHash, txType chain.TransactionType) (*chain.Transaction, error) {
	return nil, nil
}

func (pr mockPoolRequester) RequestTransactionTimeLock(pool consensus.Pool, txHash crypto.VersionnedHash, txAddr crypto.VersionnedHash, masterPublicKey crypto.PublicKey) error {
	return nil
}

func (pr mockPoolRequester) RequestTransactionUnlock(pool consensus.Pool, txHash crypto.VersionnedHash, txAddr crypto.VersionnedHash) error {
	return nil
}

func (pr mockPoolRequester) RequestTransactionValidations(pool consensus.Pool, tx chain.Transaction, minValids int, masterValid chain.MasterValidation) ([]chain.Validation, error) {
	pv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

	pubB, _ := pub.Marshal()
	vRaw := map[string]interface{}{
		"status":     chain.ValidationOK,
		"public_key": pubB,
		"timestamp":  time.Now().Unix(),
	}
	vBytes, _ := json.Marshal(vRaw)
	sig, _ := pv.Sign(vBytes)
	v, _ := chain.NewValidation(chain.ValidationOK, time.Now(), pub, sig)

	return []chain.Validation{v}, nil
}

func (pr *mockPoolRequester) RequestTransactionStorage(pool consensus.Pool, minReplicas int, tx chain.Transaction) error {
	pr.stores = append(pr.stores, tx)
	if tx.TransactionType() == chain.KeychainTransactionType {
		k, _ := chain.NewKeychain(tx)
		pr.repo.keychains = append(pr.repo.keychains, k)
	}
	if tx.TransactionType() == chain.IDTransactionType {
		id, _ := chain.NewID(tx)
		pr.repo.ids = append(pr.repo.ids, id)
	}
	return nil
}

type mockChainDB struct {
	kos       []chain.Transaction
	keychains []chain.Keychain
	ids       []chain.ID
}

func (r mockChainDB) FullKeychain(txAddr crypto.VersionnedHash) (*chain.Keychain, error) {
	sort.Slice(r.keychains, func(i, j int) bool {
		return r.keychains[i].Timestamp().Unix() > r.keychains[j].Timestamp().Unix()
	})

	if len(r.keychains) > 0 {
		return &r.keychains[0], nil
	}
	return nil, nil
}

func (r mockChainDB) LastKeychain(txAddr crypto.VersionnedHash) (*chain.Keychain, error) {
	sort.Slice(r.keychains, func(i, j int) bool {
		return r.keychains[i].Timestamp().Unix() > r.keychains[j].Timestamp().Unix()
	})

	if len(r.keychains) > 0 {
		return &r.keychains[0], nil
	}
	return nil, nil
}

func (r mockChainDB) KeychainByHash(txHash crypto.VersionnedHash) (*chain.Keychain, error) {
	for _, tx := range r.keychains {
		if bytes.Equal(tx.TransactionHash(), txHash) {
			return &tx, nil
		}
	}
	return nil, nil
}

func (r mockChainDB) IDByHash(txHash crypto.VersionnedHash) (*chain.ID, error) {
	for _, tx := range r.ids {
		if bytes.Equal(tx.TransactionHash(), txHash) {
			return &tx, nil
		}
	}
	return nil, nil
}

func (r mockChainDB) ID(addr crypto.VersionnedHash) (*chain.ID, error) {
	for _, tx := range r.ids {
		if bytes.Equal(tx.Address(), addr) {
			return &tx, nil
		}
	}
	return nil, nil
}

func (r mockChainDB) KOByHash(txHash crypto.VersionnedHash) (*chain.Transaction, error) {
	for _, tx := range r.kos {
		if bytes.Equal(tx.TransactionHash(), txHash) {
			return &tx, nil
		}
	}
	return nil, nil
}

func (r *mockChainDB) WriteKeychain(kc chain.Keychain) error {
	r.keychains = append(r.keychains, kc)
	return nil
}

func (r *mockChainDB) WriteID(id chain.ID) error {
	r.ids = append(r.ids, id)
	return nil
}

func (r *mockChainDB) WriteKO(tx chain.Transaction) error {
	r.kos = append(r.kos, tx)
	return nil
}

type mockSharedKeyReader struct {
	crossNodeKeys    []shared.NodeCrossKeyPair
	crossEmitterKeys []shared.EmitterCrossKeyPair
	authKeys         []crypto.PublicKey
}

func (r mockSharedKeyReader) EmitterCrossKeypairs() ([]shared.EmitterCrossKeyPair, error) {
	return r.crossEmitterKeys, nil
}

func (r mockSharedKeyReader) FirstNodeCrossKeypair() (shared.NodeCrossKeyPair, error) {
	return r.crossNodeKeys[0], nil
}

func (r mockSharedKeyReader) LastNodeCrossKeypair() (shared.NodeCrossKeyPair, error) {
	return r.crossNodeKeys[len(r.crossNodeKeys)-1], nil
}

func (r mockSharedKeyReader) AuthorizedNodesPublicKeys() ([]crypto.PublicKey, error) {
	return r.authKeys, nil
}

func (r mockSharedKeyReader) CrossEmitterPublicKeys() (pubKeys []crypto.PublicKey, err error) {
	for _, kp := range r.crossEmitterKeys {
		pubKeys = append(pubKeys, kp.PublicKey())
	}
	return
}

func (r mockSharedKeyReader) FirstEmitterCrossKeypair() (shared.EmitterCrossKeyPair, error) {
	return r.crossEmitterKeys[0], nil
}

type mockNodeReader struct {
	nodes []consensus.Node
}

func (db mockNodeReader) Reachables() (reachables []consensus.Node, err error) {
	for _, n := range db.nodes {
		if n.IsReachable() {
			reachables = append(reachables, n)
		}
	}
	return
}

func (db mockNodeReader) Unreachables() (unreachables []consensus.Node, err error) {
	for _, n := range db.nodes {
		if !n.IsReachable() {
			unreachables = append(unreachables, n)
		}
	}
	return
}

func (db mockNodeReader) CountReachables() (nb int, err error) {
	for _, n := range db.nodes {
		if n.IsReachable() {
			nb++
		}
	}
	return
}

func (db *mockNodeReader) FindByPublicKey(publicKey crypto.PublicKey) (found consensus.Node, err error) {
	for _, n := range db.nodes {
		if n.PublicKey().Equals(publicKey) {
			return n, nil
		}
	}
	return
}

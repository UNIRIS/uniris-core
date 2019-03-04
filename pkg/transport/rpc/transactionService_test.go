package rpc

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"sort"
	"testing"
	"time"

	"github.com/uniris/uniris-core/pkg/consensus"

	"github.com/uniris/uniris-core/pkg/chain"
	"google.golang.org/grpc/codes"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/status"

	api "github.com/uniris/uniris-core/api/protobuf-spec"
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

	techR := &mockTechDB{}
	nodeKey, _ := shared.NewNodeKeyPair(pub, pv)
	techR.nodeKeys = append(techR.nodeKeys, nodeKey)

	chainDB := &mockChainDB{}

	poolR := &mockPoolRequester{
		repo: chainDB,
	}
	txSrv := NewTransactionService(chainDB, techR, poolR, pub, pv)

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

	techR := &mockTechDB{}
	nodeKey, _ := shared.NewNodeKeyPair(pub, pv)
	techR.nodeKeys = append(techR.nodeKeys, nodeKey)

	chainDB := &mockChainDB{}

	poolR := &mockPoolRequester{
		repo: chainDB,
	}
	txSrv := NewTransactionService(chainDB, techR, poolR, pub, pv)

	data := map[string][]byte{
		"encrypted_address_by_node": []byte("addr"),
		"encrypted_wallet":          []byte("wallet"),
	}

	prop, _ := shared.NewEmitterKeyPair([]byte("pvkey"), pub)

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

	techR := &mockTechDB{}
	nodeKey, _ := shared.NewNodeKeyPair(pub, pv)
	techR.nodeKeys = append(techR.nodeKeys, nodeKey)

	chainDB := &mockChainDB{}

	poolR := &mockPoolRequester{
		repo: chainDB,
	}
	txSrv := NewTransactionService(chainDB, techR, poolR, pub, pv)

	req := &api.GetTransactionStatusRequest{
		Timestamp:       time.Now().Unix(),
		TransactionHash: crypto.Hash([]byte("tx")),
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

	techR := &mockTechDB{}
	nodeKey, _ := shared.NewNodeKeyPair(pub, pv)
	techR.nodeKeys = append(techR.nodeKeys, nodeKey)

	chainDB := &mockChainDB{}

	poolR := &mockPoolRequester{
		repo: chainDB,
	}
	txSrv := NewTransactionService(chainDB, techR, poolR, pub, pv)

	prop, _ := shared.NewEmitterKeyPair([]byte("pvkey"), pub)

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
	mv, _ := chain.NewMasterValidation([]crypto.PublicKey{}, pub, v)

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
	techDB := &mockTechDB{}
	nodeKey, _ := shared.NewNodeKeyPair(pub, pv)
	techDB.nodeKeys = append(techDB.nodeKeys, nodeKey)

	poolR := &mockPoolRequester{}
	txSrv := NewTransactionService(chainDB, techDB, poolR, pub, pv)

	pubB, _ := pub.Marshal()

	pubB, _ := pub.Marshal()

	req := &api.LockTransactionRequest{
		Timestamp:           time.Now().Unix(),
		TransactionHash:     crypto.Hash([]byte("tx")),
		MasterNodePublicKey: pubB,
		Address:             crypto.Hash([]byte("addr")),
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
<<<<<<< HEAD:pkg/transport/rpc/transactionService_test.go
	assert.True(t, chain.ContainsTimeLock(crypto.HashString("tx1"), crypto.HashString("addr1")))
=======

	assert.Len(t, locker.locks, 1)
	assert.EqualValues(t, crypto.Hash([]byte("addr")), locker.locks[0]["transaction_address"])
>>>>>>> Enable ed25519 curve, adaptative signature/encryption based on multi-crypto algo key and multi-support of hash:pkg/transport/rpc/tranasctionService_test.go
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

	techDB := &mockTechDB{}
	nodeKey, _ := shared.NewNodeKeyPair(pub, pv)
	techDB.nodeKeys = append(techDB.nodeKeys, nodeKey)
	emKey, _ := shared.NewEmitterKeyPair([]byte("encpv"), pub)
	techDB.emKeys = append(techDB.emKeys, emKey)

	chainDB := &mockChainDB{}

	poolR := &mockPoolRequester{
		repo: chainDB,
	}
	txSrv := NewTransactionService(chainDB, locker, techDB, poolR, pub, pv)
	data := map[string][]byte{
		"encrypted_address_by_node": []byte("addr"),
		"encrypted_wallet":          []byte("wallet"),
	}

	prop, _ := shared.NewEmitterKeyPair([]byte("pvkey"), pub)
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
	req := &api.LeadTransactionMiningRequest{
		Timestamp:          time.Now().Unix(),
		MinimumValidations: 1,
		Transaction:        txf,
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

	techDB := &mockTechDB{}
	nodeKey, _ := shared.NewNodeKeyPair(pub, pv)
	techDB.nodeKeys = append(techDB.nodeKeys, nodeKey)
	emKey, _ := shared.NewEmitterKeyPair([]byte("encpv"), pub)
	techDB.emKeys = append(techDB.emKeys, emKey)

	chainDB := &mockChainDB{}

	poolR := &mockPoolRequester{
		repo: chainDB,
	}
	txSrv := NewTransactionService(chainDB, techDB, poolR, pub, pv)

	prop, _ := shared.NewEmitterKeyPair([]byte("pvkey"), pub)
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
	mv, _ := chain.NewMasterValidation([]crypto.PublicKey{}, pub, v)

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

<<<<<<< HEAD:pkg/transport/rpc/transactionService_test.go
func (pr mockPoolRequester) RequestTransactionTimeLock(pool consensus.Pool, txHash crypto.VersionnedHash, txAddr crypto.VersionnedHash, masterPublicKey crypto.PublicKey) error {
=======
func (pr mockPoolRequester) RequestTransactionLock(pool consensus.Pool, txHash crypto.VersionnedHash, txAddr crypto.VersionnedHash, masterPublicKey crypto.PublicKey) error {
>>>>>>> Enable ed25519 curve, adaptative signature/encryption based on multi-crypto algo key and multi-support of hash:pkg/transport/rpc/tranasctionService_test.go
	return nil
}

func (pr mockPoolRequester) RequestTransactionUnlock(pool consensus.Pool, txHash crypto.VersionnedHash, txAddr []byte) error {
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
<<<<<<< HEAD:pkg/transport/rpc/transactionService_test.go
	kos       []chain.Transaction
	keychains []chain.Keychain
	ids       []chain.ID
=======
	inprogress []chain.Transaction
	kos        []chain.Transaction
	keychains  []chain.Keychain
	ids        []chain.ID
}

func (r mockChainDB) InProgressByHash(txHash crypto.VersionnedHash) (*chain.Transaction, error) {
	for _, tx := range r.inprogress {
		if bytes.Equal(tx.TransactionHash(), txHash) {
			return &tx, nil
		}
	}
	return nil, nil
>>>>>>> Enable ed25519 curve, adaptative signature/encryption based on multi-crypto algo key and multi-support of hash:pkg/transport/rpc/tranasctionService_test.go
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

<<<<<<< HEAD:pkg/transport/rpc/transactionService_test.go
=======
func (r *mockChainDB) WriteInProgress(tx chain.Transaction) error {
	r.inprogress = append(r.inprogress, tx)
	return nil
}

type mockLocker struct {
	locks []map[string][]byte
}

func (l *mockLocker) WriteLock(txHash crypto.VersionnedHash, txAddr crypto.VersionnedHash, masterPubKey crypto.PublicKey) error {
	masterPubk, _ := masterPubKey.Marshal()
	l.locks = append(l.locks, map[string][]byte{
		"transaction_address": txAddr,
		"transaction_hash":    txHash,
		"master_public_key":   masterPubk,
	})
	return nil
}
func (l *mockLocker) RemoveLock(txHash crypto.VersionnedHash, txAddr crypto.VersionnedHash) error {
	pos := l.findLockPosition(txHash, txAddr)
	if pos > -1 {
		l.locks = append(l.locks[:pos], l.locks[pos+1:]...)
	}
	return nil
}
func (l mockLocker) ContainsLock(txHash crypto.VersionnedHash, txAddr crypto.VersionnedHash) (bool, error) {
	return l.findLockPosition(txHash, txAddr) > -1, nil
}

func (l mockLocker) findLockPosition(txHash crypto.VersionnedHash, txAddr crypto.VersionnedHash) int {
	for i, lock := range l.locks {
		if bytes.Equal(lock["transaction_hash"], txHash) && bytes.Equal(lock["transaction_address"], txAddr) {
			return i
		}
	}
	return -1
}

>>>>>>> Enable ed25519 curve, adaptative signature/encryption based on multi-crypto algo key and multi-support of hash:pkg/transport/rpc/tranasctionService_test.go
type mockTechDB struct {
	emKeys   shared.EmitterKeys
	nodeKeys []shared.NodeKeyPair
}

func (db mockTechDB) EmitterKeys() (shared.EmitterKeys, error) {
	return db.emKeys, nil
}

func (db mockTechDB) NodeLastKeys() (shared.NodeKeyPair, error) {
	return db.nodeKeys[len(db.nodeKeys)-1], nil
}

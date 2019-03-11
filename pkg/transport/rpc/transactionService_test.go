package rpc

import (
	"context"
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
	pub, pv := crypto.GenerateKeys()

	techR := &mockTechDB{}
	nodeKey, _ := shared.NewKeyPair(pub, pv)
	techR.nodeKeys = append(techR.nodeKeys, nodeKey)

	chainDB := &mockChainDB{}

	poolR := &mockPoolRequester{
		repo: chainDB,
	}
	txSrv := NewTransactionService(chainDB, techR, poolR, pub, pv)

	req := &api.GetLastTransactionRequest{
		Timestamp:          time.Now().Unix(),
		TransactionAddress: crypto.HashString("address"),
		Type:               api.TransactionType_KEYCHAIN,
	}
	reqBytes, _ := json.Marshal(req)
	sig, _ := crypto.Sign(string(reqBytes), pv)
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

	pub, pv := crypto.GenerateKeys()

	techR := &mockTechDB{}
	nodeKey, _ := shared.NewKeyPair(pub, pv)
	techR.nodeKeys = append(techR.nodeKeys, nodeKey)

	chainDB := &mockChainDB{}

	poolR := &mockPoolRequester{
		repo: chainDB,
	}
	txSrv := NewTransactionService(chainDB, techR, poolR, pub, pv)

	data := map[string]string{
		"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
		"encrypted_wallet":          hex.EncodeToString([]byte("wallet")),
	}

	prop, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("pvkey")), pub)

	txRaw := map[string]interface{}{
		"addr":                    crypto.HashString("addr"),
		"data":                    data,
		"timestamp":               time.Now().Unix(),
		"type":                    chain.KeychainTransactionType,
		"public_key":              pub,
		"em_shared_keys_proposal": prop,
	}
	txBytes, _ := json.Marshal(txRaw)
	sig, _ := crypto.Sign(string(txBytes), pv)
	txRaw["signature"] = sig

	txByteWithSign, _ := json.Marshal(txRaw)
	emSig, _ := crypto.Sign(string(txByteWithSign), pv)
	txRaw["em_signature"] = emSig
	txBytes, _ = json.Marshal(txRaw)

	tx, _ := chain.NewTransaction(crypto.HashString("addr"), chain.KeychainTransactionType, data, time.Now(), pub, prop, sig, emSig, crypto.HashBytes(txBytes))
	keychain, _ := chain.NewKeychain(tx)
	chainDB.keychains = append(chainDB.keychains, keychain)

	req := &api.GetLastTransactionRequest{
		Timestamp:          time.Now().Unix(),
		TransactionAddress: crypto.HashString("addr"),
		Type:               api.TransactionType_KEYCHAIN,
	}
	reqBytes, _ := json.Marshal(req)
	sigReq, _ := crypto.Sign(string(reqBytes), pv)
	req.SignatureRequest = sigReq

	res, err := txSrv.GetLastTransaction(context.TODO(), req)
	assert.Nil(t, err)
	assert.NotEmpty(t, res.SignatureResponse)
	assert.NotNil(t, res.Transaction)
	assert.Equal(t, crypto.HashBytes(txBytes), res.Transaction.TransactionHash)

	resBytes, _ := json.Marshal(&api.GetLastTransactionResponse{
		Timestamp:   res.Timestamp,
		Transaction: res.Transaction,
	})
	assert.Nil(t, crypto.VerifySignature(string(resBytes), pub, res.SignatureResponse))
}

/*
Scenario: Receive get transaction status request
	Given no transaction stored
	When I want to request the transactions status for this transaction hash
	Then I get a status unknown
*/
func TestHandleGetTransactionStatus(t *testing.T) {

	pub, pv := crypto.GenerateKeys()

	techR := &mockTechDB{}
	nodeKey, _ := shared.NewKeyPair(pub, pv)
	techR.nodeKeys = append(techR.nodeKeys, nodeKey)

	chainDB := &mockChainDB{}

	poolR := &mockPoolRequester{
		repo: chainDB,
	}
	txSrv := NewTransactionService(chainDB, techR, poolR, pub, pv)

	req := &api.GetTransactionStatusRequest{
		Timestamp:       time.Now().Unix(),
		TransactionHash: crypto.HashString("tx1"),
	}
	reqBytes, _ := json.Marshal(req)
	sig, _ := crypto.Sign(string(reqBytes), pv)
	req.SignatureRequest = sig

	res, err := txSrv.GetTransactionStatus(context.TODO(), req)
	assert.Nil(t, err)
	assert.Equal(t, api.TransactionStatus_UNKNOWN, res.Status)
	resBytes, _ := json.Marshal(&api.GetTransactionStatusResponse{
		Timestamp: res.Timestamp,
		Status:    res.Status,
	})
	assert.Nil(t, crypto.VerifySignature(string(resBytes), pub, res.SignatureResponse))
}

/*
Scenario: Receive storage  transaction request
	Given a transaction
	When I want to request to store of the transaction
	Then the transaction is stored
*/
func TestHandleStoreTransaction(t *testing.T) {

	pub, pv := crypto.GenerateKeys()

	techR := &mockTechDB{}
	nodeKey, _ := shared.NewKeyPair(pub, pv)
	techR.nodeKeys = append(techR.nodeKeys, nodeKey)

	chainDB := &mockChainDB{}

	poolR := &mockPoolRequester{
		repo: chainDB,
	}
	txSrv := NewTransactionService(chainDB, techR, poolR, pub, pv)

	prop, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("pvkey")), pub)

	data := map[string]string{
		"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
		"encrypted_wallet":          hex.EncodeToString([]byte("wallet")),
	}
	txRaw := map[string]interface{}{
		"addr":                    crypto.HashString("addr"),
		"data":                    data,
		"timestamp":               time.Now().Unix(),
		"type":                    chain.KeychainTransactionType,
		"public_key":              pub,
		"em_shared_keys_proposal": prop,
	}
	txBytes, _ := json.Marshal(txRaw)
	sig, _ := crypto.Sign(string(txBytes), pv)
	txRaw["signature"] = sig

	txByteWithSign, _ := json.Marshal(txRaw)
	emSig, _ := crypto.Sign(string(txByteWithSign), pv)
	txRaw["em_signature"] = emSig
	txBytes, _ = json.Marshal(txRaw)

	txBytes, _ = json.Marshal(txRaw)

	tx, _ := chain.NewTransaction(crypto.HashString("addr"), chain.KeychainTransactionType, data, time.Now(), pub, prop, sig, emSig, crypto.HashBytes(txBytes))

	vRaw := map[string]interface{}{
		"status":     chain.ValidationOK,
		"public_key": pub,
		"timestamp":  time.Now().Unix(),
	}
	vBytes, _ := json.Marshal(vRaw)
	vSig, _ := crypto.Sign(string(vBytes), pv)
	v, _ := chain.NewValidation(chain.ValidationOK, time.Now(), pub, vSig)
	wHeaders := []chain.NodeHeader{chain.NewNodeHeader("pub", false, false, 0)}
	vHeaders := []chain.NodeHeader{chain.NewNodeHeader("pub", false, false, 0)}
	sHeaders := []chain.NodeHeader{chain.NewNodeHeader("pub", false, false, 0)}
	mv, _ := chain.NewMasterValidation([]string{}, pub, v, wHeaders, vHeaders, sHeaders)

	req := &api.StoreTransactionRequest{
		Timestamp: time.Now().Unix(),
		MinedTransaction: &api.MinedTransaction{
			Transaction:        formatAPITransaction(tx),
			MasterValidation:   formatAPIMasterValidation(mv),
			ConfirmValidations: []*api.Validation{formatAPIValidation(v)},
		},
	}

	reqBytes, _ := json.Marshal(req)
	sigReq, _ := crypto.Sign(string(reqBytes), pv)
	req.SignatureRequest = sigReq

	res, err := txSrv.StoreTransaction(context.TODO(), req)
	assert.Nil(t, err)

	resBytes, _ := json.Marshal(&api.StoreTransactionResponse{
		Timestamp: res.Timestamp,
	})
	assert.Nil(t, crypto.VerifySignature(string(resBytes), pub, res.SignatureResponse))

	assert.Len(t, chainDB.keychains, 1)
	assert.Equal(t, crypto.HashBytes(txBytes), chainDB.keychains[0].TransactionHash())

}

/*
Scenario: Receive lock transaction request
	Given a transaction to lock
	When I want to request to lock it
	Then I get not error and the lock is stored
*/
func TestHandleLockTransaction(t *testing.T) {

	pub, pv := crypto.GenerateKeys()

	chainDB := &mockChainDB{}
	techDB := &mockTechDB{}
	nodeKey, _ := shared.NewKeyPair(pub, pv)
	techDB.nodeKeys = append(techDB.nodeKeys, nodeKey)

	poolR := &mockPoolRequester{}
	txSrv := NewTransactionService(chainDB, techDB, poolR, pub, pv)

	req := &api.TimeLockTransactionRequest{
		Timestamp:           time.Now().Unix(),
		TransactionHash:     crypto.HashString("tx1"),
		MasterNodePublicKey: pub,
		Address:             crypto.HashString("addr1"),
	}
	reqBytes, _ := json.Marshal(req)
	sig, _ := crypto.Sign(string(reqBytes), pv)
	req.SignatureRequest = sig

	res, err := txSrv.TimeLockTransaction(context.TODO(), req)
	assert.Nil(t, err)
	resBytes, _ := json.Marshal(&api.TimeLockTransactionResponse{
		Timestamp: res.Timestamp,
	})
	assert.Nil(t, crypto.VerifySignature(string(resBytes), pub, res.SignatureResponse))
	assert.True(t, chain.ContainsTimeLock(crypto.HashString("tx1"), crypto.HashString("addr1")))
}

/*
Scenario: Receive lead mining transaction request
	Given a transaction to validate
	When I want to request to lead mining of the transaction
	Then I get not error
*/
func TestHandleLeadTransactionMining(t *testing.T) {

	pub, pv := crypto.GenerateKeys()

	techDB := &mockTechDB{}
	nodeKey, _ := shared.NewKeyPair(pub, pv)
	techDB.nodeKeys = append(techDB.nodeKeys, nodeKey)
	emKey, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("encpv")), pub)
	techDB.emKeys = append(techDB.emKeys, emKey)

	chainDB := &mockChainDB{}

	poolR := &mockPoolRequester{
		repo: chainDB,
	}
	txSrv := NewTransactionService(chainDB, techDB, poolR, pub, pv)
	data := map[string]string{
		"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
		"encrypted_wallet":          hex.EncodeToString([]byte("wallet")),
	}

	prop, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("pvkey")), pub)
	txRaw := map[string]interface{}{
		"addr":                    crypto.HashString("addr"),
		"data":                    data,
		"timestamp":               time.Now().Unix(),
		"type":                    chain.KeychainTransactionType,
		"public_key":              pub,
		"em_shared_keys_proposal": prop,
	}
	txBytes, _ := json.Marshal(txRaw)
	sig, _ := crypto.Sign(string(txBytes), pv)
	txRaw["signature"] = sig

	txByteWithSig, _ := json.Marshal(txRaw)
	emSig, _ := crypto.Sign(string(txByteWithSig), pv)
	txRaw["em_signature"] = emSig

	txBytes, _ = json.Marshal(txRaw)

	tx, _ := chain.NewTransaction(crypto.HashString("addr"), chain.KeychainTransactionType, data, time.Now(), pub, prop, sig, emSig, crypto.HashBytes(txBytes))
	req := &api.LeadTransactionMiningRequest{
		Timestamp:          time.Now().Unix(),
		MinimumValidations: 1,
		WelcomeHeaders: []*api.NodeHeader{
			&api.NodeHeader{
				IsMaster:      true,
				IsUnreachable: false,
				PatchNumber:   1,
				PublicKey:     pub,
			},
		},
		Transaction: formatAPITransaction(tx),
	}

	reqBytes, _ := json.Marshal(req)
	sigReq, _ := crypto.Sign(string(reqBytes), pv)
	req.SignatureRequest = sigReq

	res, err := txSrv.LeadTransactionMining(context.TODO(), req)
	assert.Nil(t, err)

	time.Sleep(1 * time.Second)

	resBytes, _ := json.Marshal(&api.LeadTransactionMiningResponse{
		Timestamp: res.Timestamp,
	})
	assert.Nil(t, crypto.VerifySignature(string(resBytes), pub, res.SignatureResponse))

	assert.Len(t, chainDB.keychains, 1)
	assert.Equal(t, crypto.HashString("addr"), chainDB.keychains[0].Address())
}

/*
Scenario: Receive confirmation of validations transaction request
	Given a transaction to validate
	When I want to request to validation of the transaction
	Then I get the node validation
*/
func TestHandleConfirmValiation(t *testing.T) {

	pub, pv := crypto.GenerateKeys()

	techDB := &mockTechDB{}
	nodeKey, _ := shared.NewKeyPair(pub, pv)
	techDB.nodeKeys = append(techDB.nodeKeys, nodeKey)
	emKey, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("encpv")), pub)
	techDB.emKeys = append(techDB.emKeys, emKey)

	chainDB := &mockChainDB{}

	poolR := &mockPoolRequester{
		repo: chainDB,
	}
	txSrv := NewTransactionService(chainDB, techDB, poolR, pub, pv)

	prop, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("pvkey")), pub)
	data := map[string]string{
		"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
		"encrypted_wallet":          hex.EncodeToString([]byte("wallet")),
	}

	txRaw := map[string]interface{}{
		"addr":                    crypto.HashString("addr"),
		"data":                    data,
		"timestamp":               time.Now().Unix(),
		"type":                    chain.KeychainTransactionType,
		"public_key":              pub,
		"em_shared_keys_proposal": prop,
	}
	txBytes, _ := json.Marshal(txRaw)
	sig, _ := crypto.Sign(string(txBytes), pv)
	txRaw["signature"] = sig
	txByteWithSig, _ := json.Marshal(txRaw)
	emSig, _ := crypto.Sign(string(txByteWithSig), pv)
	txRaw["em_signature"] = emSig
	txBytes, _ = json.Marshal(txRaw)
	tx, _ := chain.NewTransaction(crypto.HashString("addr"), chain.KeychainTransactionType, data, time.Now(), pub, prop, sig, emSig, crypto.HashBytes(txBytes))

	vRaw := map[string]interface{}{
		"status":     chain.ValidationOK,
		"public_key": pub,
		"timestamp":  time.Now().Unix(),
	}

	vBytes, _ := json.Marshal(vRaw)
	vSig, _ := crypto.Sign(string(vBytes), pv)
	v, _ := chain.NewValidation(chain.ValidationOK, time.Now(), pub, vSig)
	wHeaders := []chain.NodeHeader{chain.NewNodeHeader("pub", false, false, 0)}
	vHeaders := []chain.NodeHeader{chain.NewNodeHeader("pub", false, false, 0)}
	sHeaders := []chain.NodeHeader{chain.NewNodeHeader("pub", false, false, 0)}
	mv, _ := chain.NewMasterValidation([]string{}, pub, v, wHeaders, vHeaders, sHeaders)
	req := &api.ConfirmTransactionValidationRequest{
		Transaction:      formatAPITransaction(tx),
		Timestamp:        time.Now().Unix(),
		MasterValidation: formatAPIMasterValidation(mv),
	}

	reqBytes, _ := json.Marshal(req)
	sigReq, _ := crypto.Sign(string(reqBytes), pv)
	req.SignatureRequest = sigReq

	res, err := txSrv.ConfirmTransactionValidation(context.TODO(), req)
	assert.Nil(t, err)

	resBytes, _ := json.Marshal(&api.ConfirmTransactionValidationResponse{
		Timestamp:  res.Timestamp,
		Validation: res.Validation,
	})
	assert.Nil(t, crypto.VerifySignature(string(resBytes), pub, res.SignatureResponse))

	assert.NotNil(t, res.Validation)
	assert.Equal(t, api.Validation_OK, res.Validation.Status)
	assert.Equal(t, pub, res.Validation.PublicKey)
}

type mockPoolRequester struct {
	stores []chain.Transaction
	repo   *mockChainDB
}

func (pr mockPoolRequester) RequestLastTransaction(pool consensus.Pool, txAddr string, txType chain.TransactionType) (*chain.Transaction, error) {
	return nil, nil
}

func (pr mockPoolRequester) RequestTransactionTimeLock(pool consensus.Pool, txHash string, txAddr string, masterPublicKey string) error {
	return nil
}

func (pr mockPoolRequester) RequestTransactionUnlock(pool consensus.Pool, txHash string, txAddr string) error {
	return nil
}

func (pr mockPoolRequester) RequestTransactionValidations(pool consensus.Pool, tx chain.Transaction, minValids int, masterValid chain.MasterValidation) ([]chain.Validation, error) {
	pub, pv := crypto.GenerateKeys()

	vRaw := map[string]interface{}{
		"status":     chain.ValidationOK,
		"public_key": pub,
		"timestamp":  time.Now().Unix(),
	}
	vBytes, _ := json.Marshal(vRaw)
	sig, _ := crypto.Sign(string(vBytes), pv)
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

func (r mockChainDB) FullKeychain(txAddr string) (*chain.Keychain, error) {
	sort.Slice(r.keychains, func(i, j int) bool {
		return r.keychains[i].Timestamp().Unix() > r.keychains[j].Timestamp().Unix()
	})

	if len(r.keychains) > 0 {
		return &r.keychains[0], nil
	}
	return nil, nil
}

func (r mockChainDB) LastKeychain(txAddr string) (*chain.Keychain, error) {
	sort.Slice(r.keychains, func(i, j int) bool {
		return r.keychains[i].Timestamp().Unix() > r.keychains[j].Timestamp().Unix()
	})

	if len(r.keychains) > 0 {
		return &r.keychains[0], nil
	}
	return nil, nil
}

func (r mockChainDB) KeychainByHash(txHash string) (*chain.Keychain, error) {
	for _, tx := range r.keychains {
		if tx.TransactionHash() == txHash {
			return &tx, nil
		}
	}
	return nil, nil
}

func (r mockChainDB) IDByHash(txHash string) (*chain.ID, error) {
	for _, tx := range r.ids {
		if tx.TransactionHash() == txHash {
			return &tx, nil
		}
	}
	return nil, nil
}

func (r mockChainDB) ID(addr string) (*chain.ID, error) {
	for _, tx := range r.ids {
		if tx.Address() == addr {
			return &tx, nil
		}
	}
	return nil, nil
}

func (r mockChainDB) KOByHash(txHash string) (*chain.Transaction, error) {
	for _, tx := range r.kos {
		if tx.TransactionHash() == txHash {
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

type mockTechDB struct {
	emKeys   shared.EmitterKeys
	nodeKeys []shared.KeyPair
}

func (db mockTechDB) EmitterKeys() (shared.EmitterKeys, error) {
	return db.emKeys, nil
}

func (db mockTechDB) NodeLastKeys() (shared.KeyPair, error) {
	return db.nodeKeys[len(db.nodeKeys)-1], nil
}

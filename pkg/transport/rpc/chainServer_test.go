package rpc

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"log"
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
	minerKey, _ := shared.NewMinerKeyPair(pub, pv)
	techR.minerKeys = append(techR.minerKeys, minerKey)

	chainDB := &mockChainDB{}

	poolR := &mockPoolRequester{
		repo: chainDB,
	}
	chainSrv := NewChainServer(chainDB, techR, poolR)

	req := &api.LastTransactionRequest{
		Timestamp:          time.Now().Unix(),
		TransactionAddress: crypto.HashString("address"),
		Type:               api.TransactionType_KEYCHAIN,
	}
	reqBytes, _ := json.Marshal(req)
	sig, _ := crypto.Sign(string(reqBytes), pv)
	req.SignatureRequest = sig

	_, err := chainSrv.GetLastTransaction(context.TODO(), req)
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
	minerKey, _ := shared.NewMinerKeyPair(pub, pv)
	techR.minerKeys = append(techR.minerKeys, minerKey)

	chainDB := &mockChainDB{}

	poolR := &mockPoolRequester{
		repo: chainDB,
	}
	chainSrv := NewChainServer(chainDB, techR, poolR)

	data := map[string]string{
		"encrypted_address": hex.EncodeToString([]byte("addr")),
		"encrypted_wallet":  hex.EncodeToString([]byte("wallet")),
	}

	prop, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("pvkey")), pub)

	txRaw := map[string]interface{}{
		"address":                 crypto.HashString("addr"),
		"data":                    data,
		"timestamp":               time.Now().Unix(),
		"type":                    chain.KeychainTransactionType,
		"public_key":              pub,
		"em_shared_keys_proposal": prop,
	}
	txBytes, _ := json.Marshal(txRaw)
	sig, _ := crypto.Sign(string(txBytes), pv)
	txRaw["signature"] = sig
	txRaw["em_signature"] = sig
	txBytes, _ = json.Marshal(txRaw)

	tx, _ := chain.NewTransaction(crypto.HashString("addr"), chain.KeychainTransactionType, data, time.Now(), pub, prop, sig, sig, crypto.HashBytes(txBytes))
	keychain, _ := chain.NewKeychain(tx)
	chainDB.keychains = append(chainDB.keychains, keychain)

	req := &api.LastTransactionRequest{
		Timestamp:          time.Now().Unix(),
		TransactionAddress: crypto.HashString("addr"),
		Type:               api.TransactionType_KEYCHAIN,
	}
	reqBytes, _ := json.Marshal(req)
	sigReq, _ := crypto.Sign(string(reqBytes), pv)
	req.SignatureRequest = sigReq

	res, err := chainSrv.GetLastTransaction(context.TODO(), req)
	assert.Nil(t, err)
	assert.NotEmpty(t, res.SignatureResponse)
	assert.NotNil(t, res.Transaction)
	assert.Equal(t, crypto.HashBytes(txBytes), res.Transaction.TransactionHash)

	resBytes, _ := json.Marshal(&api.LastTransactionResponse{
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
	minerKey, _ := shared.NewMinerKeyPair(pub, pv)
	techR.minerKeys = append(techR.minerKeys, minerKey)

	chainDB := &mockChainDB{}

	poolR := &mockPoolRequester{
		repo: chainDB,
	}
	chainSrv := NewChainServer(chainDB, techR, poolR)

	req := &api.TransactionStatusRequest{
		Timestamp:       time.Now().Unix(),
		TransactionHash: crypto.HashString("tx"),
	}
	reqBytes, _ := json.Marshal(req)
	sig, _ := crypto.Sign(string(reqBytes), pv)
	req.SignatureRequest = sig

	res, err := chainSrv.GetTransactionStatus(context.TODO(), req)
	assert.Nil(t, err)
	assert.Equal(t, api.TransactionStatusResponse_UNKNOWN, res.Status)
	resBytes, _ := json.Marshal(&api.TransactionStatusResponse{
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
	minerKey, _ := shared.NewMinerKeyPair(pub, pv)
	techR.minerKeys = append(techR.minerKeys, minerKey)

	chainDB := &mockChainDB{}

	poolR := &mockPoolRequester{
		repo: chainDB,
	}
	chainSrv := NewChainServer(chainDB, techR, poolR)

	prop, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("pvkey")), pub)

	data := map[string]string{
		"encrypted_address": hex.EncodeToString([]byte("addr")),
		"encrypted_wallet":  hex.EncodeToString([]byte("wallet")),
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
	txRaw["em_signature"] = sig
	txBytes, _ = json.Marshal(txRaw)

	tx, _ := chain.NewTransaction(crypto.HashString("addr"), chain.KeychainTransactionType, data, time.Now(), pub, prop, sig, sig, crypto.HashBytes(txBytes))

	vRaw := map[string]interface{}{
		"status":     chain.ValidationOK,
		"public_key": pub,
		"timestamp":  time.Now().Unix(),
	}
	vBytes, _ := json.Marshal(vRaw)
	vSig, _ := crypto.Sign(string(vBytes), pv)
	v, _ := chain.NewMinerValidation(chain.ValidationOK, time.Now(), pub, vSig)
	mv, _ := chain.NewMasterValidation([]string{}, pub, v)

	req := &api.StoreRequest{
		Timestamp: time.Now().Unix(),
		MinedTransaction: &api.MinedTransaction{
			Transaction:        formatAPITransaction(tx),
			MasterValidation:   formatAPIMasterValidation(mv),
			ConfirmValidations: []*api.MinerValidation{formatAPIValidation(v)},
		},
	}

	reqBytes, _ := json.Marshal(req)
	sigReq, _ := crypto.Sign(string(reqBytes), pv)
	req.SignatureRequest = sigReq

	res, err := chainSrv.StoreTransaction(context.TODO(), req)
	assert.Nil(t, err)

	resBytes, _ := json.Marshal(&api.StoreResponse{
		Timestamp: res.Timestamp,
	})
	assert.Nil(t, crypto.VerifySignature(string(resBytes), pub, res.SignatureResponse))

	assert.Len(t, chainDB.keychains, 1)
	assert.Equal(t, crypto.HashBytes(txBytes), chainDB.keychains[0].TransactionHash())

}

type mockPoolRequester struct {
	stores []chain.Transaction
	repo   *mockChainDB
}

func (pr mockPoolRequester) RequestLastTransaction(pool consensus.Pool, txAddr string, txType chain.TransactionType) (*chain.Transaction, error) {
	return nil, nil
}

func (pr mockPoolRequester) RequestTransactionLock(pool consensus.Pool, txHash string, txAddr string, masterPublicKey string) error {
	return nil
}

func (pr mockPoolRequester) RequestTransactionUnlock(pool consensus.Pool, txHash string, txAddr string) error {
	return nil
}

func (pr mockPoolRequester) RequestTransactionValidations(pool consensus.Pool, tx chain.Transaction, minValids int, masterValid chain.MasterValidation) ([]chain.MinerValidation, error) {
	pub, pv := crypto.GenerateKeys()

	vRaw := map[string]interface{}{
		"status":     chain.ValidationOK,
		"public_key": pub,
		"timestamp":  time.Now().Unix(),
	}
	vBytes, _ := json.Marshal(vRaw)
	sig, _ := crypto.Sign(string(vBytes), pv)
	v, _ := chain.NewMinerValidation(chain.ValidationOK, time.Now(), pub, sig)

	return []chain.MinerValidation{v}, nil
}

func (pr *mockPoolRequester) RequestTransactionStorage(pool consensus.Pool, minReplicas int, tx chain.Transaction) error {
	pr.stores = append(pr.stores, tx)
	if tx.TransactionType() == chain.KeychainTransactionType {
		k, _ := chain.NewKeychain(tx)
		pr.repo.keychains = append(pr.repo.keychains, k)
	}
	if tx.TransactionType() == chain.IDTransactionType {
		id, err := chain.NewID(tx)
		log.Print(err)
		pr.repo.ids = append(pr.repo.ids, id)
	}
	return nil
}

type mockChainDB struct {
	inprogress []chain.Transaction
	kos        []chain.Transaction
	keychains  []chain.Keychain
	ids        []chain.ID
}

func (r mockChainDB) InProgressByHash(txHash string) (*chain.Transaction, error) {
	for _, tx := range r.inprogress {
		if tx.TransactionHash() == txHash {
			return &tx, nil
		}
	}
	return nil, nil
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

func (r *mockChainDB) WriteInProgress(tx chain.Transaction) error {
	r.inprogress = append(r.inprogress, tx)
	return nil
}

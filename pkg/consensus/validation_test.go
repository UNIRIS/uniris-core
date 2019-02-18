package consensus

import (
	"testing"

	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/uniris/uniris-core/pkg/chain"
	"github.com/uniris/uniris-core/pkg/crypto"
	"github.com/uniris/uniris-core/pkg/shared"
)

/*
Scenario: request transaction validations
	Given a transaction to validate
	When I aprop validations to a pool
	Then I get validations from them
*/
func TestRequestValidations(t *testing.T) {
	poolR := &mockPoolRequester{}
	pub, pv := crypto.GenerateKeys()

	v, _ := buildValidation(chain.ValidationOK, pub, pv)
	mv, _ := chain.NewMasterValidation([]string{}, pub, v)

	prop, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("pvKey")), pub)

	data := map[string]string{
		"encrypted_aes_key":         hex.EncodeToString([]byte("aesKey")),
		"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
		"encrypted_address_by_id":   hex.EncodeToString([]byte("addr")),
	}

	txRaw, _ := json.Marshal(map[string]interface{}{
		"addr":                    crypto.HashString("addr"),
		"data":                    data,
		"timestamp":               time.Now().Unix(),
		"type":                    chain.KeychainTransactionType,
		"public_key":              pub,
		"em_shared_keys_proposal": prop,
	})
	sig, _ := crypto.Sign(string(txRaw), pv)
	txSigned, _ := json.Marshal(map[string]interface{}{
		"addr":                    crypto.HashString("addr"),
		"data":                    data,
		"timestamp":               time.Now().Unix(),
		"type":                    chain.KeychainTransactionType,
		"public_key":              pub,
		"em_shared_keys_proposal": prop,
		"signature":               sig,
	})
	emSig, _ := crypto.Sign(string(txSigned), pv)
	txEmSigned, _ := json.Marshal(map[string]interface{}{
		"addr":                    crypto.HashString("addr"),
		"data":                    data,
		"timestamp":               time.Now().Unix(),
		"type":                    chain.KeychainTransactionType,
		"public_key":              pub,
		"em_shared_keys_proposal": prop,
		"signature":               sig,
		"em_signature":            emSig,
	})
	tx, err := chain.NewTransaction(crypto.HashString("addr"), chain.KeychainTransactionType, data, time.Now(), pub, prop, sig, emSig, crypto.HashBytes(txEmSigned))

	valids, err := requestValidations(tx, mv, Pool{}, 1, poolR)
	assert.Nil(t, err)
	assert.NotEmpty(t, valids)
	assert.Equal(t, chain.ValidationOK, valids[0].Status())
}

/*
Scenario: Create a node validation
	Given a validation status
	When I want to create node validation
	Then I get a validation signed
*/
func TestBuildValidation(t *testing.T) {
	pub, pv := crypto.GenerateKeys()

	v, err := buildValidation(chain.ValidationOK, pub, pv)
	assert.Nil(t, err)
	assert.Equal(t, pub, v.PublicKey())
	assert.Nil(t, err)
	assert.Equal(t, time.Now().Unix(), v.Timestamp().Unix())
	assert.Equal(t, chain.ValidationOK, v.Status())
	ok, err := v.IsValid()
	assert.True(t, ok)
}

/*
Scenario: Validate an incoming transaction
	Given a valid transaction
	When I want to valid the transaction
	Then I get a validation with status OK
*/
func TestValidateTransaction(t *testing.T) {
	pub, pv := crypto.GenerateKeys()

	prop, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("pvKey")), pub)

	data := map[string]string{
		"encrypted_aes_key":         hex.EncodeToString([]byte("aesKey")),
		"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
		"encrypted_address_by_id":   hex.EncodeToString([]byte("addr")),
	}

	txRaw, _ := json.Marshal(map[string]interface{}{
		"addr":                    crypto.HashString("addr"),
		"data":                    data,
		"timestamp":               time.Now().Unix(),
		"type":                    chain.KeychainTransactionType,
		"public_key":              pub,
		"em_shared_keys_proposal": prop,
	})
	sig, _ := crypto.Sign(string(txRaw), pv)
	txSigned, _ := json.Marshal(map[string]interface{}{
		"addr":                    crypto.HashString("addr"),
		"data":                    data,
		"timestamp":               time.Now().Unix(),
		"type":                    chain.KeychainTransactionType,
		"public_key":              pub,
		"em_shared_keys_proposal": prop,
		"signature":               sig,
	})
	emSig, _ := crypto.Sign(string(txSigned), pv)
	txEmSigned, _ := json.Marshal(map[string]interface{}{
		"addr":                    crypto.HashString("addr"),
		"data":                    data,
		"timestamp":               time.Now().Unix(),
		"type":                    chain.KeychainTransactionType,
		"public_key":              pub,
		"em_shared_keys_proposal": prop,
		"signature":               sig,
		"em_signature":            emSig,
	})
	tx, err := chain.NewTransaction(crypto.HashString("addr"), chain.KeychainTransactionType, data, time.Now(), pub, prop, sig, emSig, crypto.HashBytes(txEmSigned))
	assert.Nil(t, err)

	v, _ := buildValidation(chain.ValidationOK, pub, pv)
	mv, _ := chain.NewMasterValidation([]string{}, pub, v)

	valid, err := ConfirmTransactionValidation(tx, mv, pub, pv)
	assert.Nil(t, err)
	assert.Equal(t, chain.ValidationOK, valid.Status())
}

/*
Scenario: Validate an incoming transaction with invalid integrity
	Given a transaction with invalid transaction hash or signature
	When I want to valid the transaction
	Then I get a validation with status KO
*/
func TestValidateTransactionWithBadIntegrity(t *testing.T) {
	pub, pv := crypto.GenerateKeys()

	prop, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("pvKey")), pub)

	data := map[string]string{
		"encrypted_aes_key":         hex.EncodeToString([]byte("aesKey")),
		"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
		"encrypted_address_by_id":   hex.EncodeToString([]byte("addr")),
	}

	sig, _ := crypto.Sign("hello", pv)
	tx, _ := chain.NewTransaction(crypto.HashString("addr"), chain.IDTransactionType, data, time.Now(), pub, prop, sig, sig, crypto.HashString("hash"))

	v, _ := buildValidation(chain.ValidationOK, pub, pv)
	mv, _ := chain.NewMasterValidation([]string{}, pub, v)
	valid, err := ConfirmTransactionValidation(tx, mv, pub, pv)
	assert.Nil(t, err)
	assert.Equal(t, chain.ValidationKO, valid.Status())
}

/*
Scenario: Perform Proof of work
	Given a transaction and em chain keypair stored
	When I want to perform the proof of work of this transaction
	Then I get the valid public key
*/
func TestPerformPOW(t *testing.T) {

	pub, pv := crypto.GenerateKeys()

	emReader := &mockEmitterReader{}
	emKP, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("pvKey")), pub)
	emReader.emKeys = append(emReader.emKeys, emKP)

	prop, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("pvKey")), pub)

	data := map[string]string{
		"encrypted_aes_key":         hex.EncodeToString([]byte("aesKey")),
		"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
		"encrypted_address_by_id":   hex.EncodeToString([]byte("addr")),
	}

	txRaw, _ := json.Marshal(map[string]interface{}{
		"addr":                    crypto.HashString("addr"),
		"data":                    data,
		"timestamp":               time.Now().Unix(),
		"type":                    chain.KeychainTransactionType,
		"public_key":              pub,
		"em_shared_keys_proposal": prop,
	})
	sig, _ := crypto.Sign(string(txRaw), pv)
	txSigned, _ := json.Marshal(map[string]interface{}{
		"addr":                    crypto.HashString("addr"),
		"data":                    data,
		"timestamp":               time.Now().Unix(),
		"type":                    chain.KeychainTransactionType,
		"public_key":              pub,
		"em_shared_keys_proposal": prop,
		"signature":               sig,
	})
	emSig, _ := crypto.Sign(string(txSigned), pv)
	txEmSigned, _ := json.Marshal(map[string]interface{}{
		"addr":                    crypto.HashString("addr"),
		"data":                    data,
		"timestamp":               time.Now().Unix(),
		"type":                    chain.KeychainTransactionType,
		"public_key":              pub,
		"em_shared_keys_proposal": prop,
		"signature":               sig,
		"em_signature":            emSig,
	})
	tx, err := chain.NewTransaction(crypto.HashString("addr"), chain.KeychainTransactionType, data, time.Now(), pub, prop, sig, emSig, crypto.HashBytes(txEmSigned))
	assert.Nil(t, err)

	pow, err := proofOfWork(tx, emReader)
	assert.Nil(t, err)
	assert.Equal(t, pub, pow)
}

/*
Scenario: Pre-validate a transaction
	Given a transaction
	When I want to prevalidate this transaction
	Then I get the node validation and the proof of work
*/
func TestPreValidateTransaction(t *testing.T) {

	pub, pv := crypto.GenerateKeys()
	emReader := &mockEmitterReader{}
	emKP, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("pvKey")), pub)
	emReader.emKeys = append(emReader.emKeys, emKP)

	prop, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("pvKey")), pub)

	data := map[string]string{
		"encrypted_aes_key":         hex.EncodeToString([]byte("aesKey")),
		"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
		"encrypted_address_by_id":   hex.EncodeToString([]byte("addr")),
	}

	txRaw, _ := json.Marshal(map[string]interface{}{
		"addr":                    crypto.HashString("addr"),
		"data":                    data,
		"timestamp":               time.Now().Unix(),
		"type":                    chain.KeychainTransactionType,
		"public_key":              pub,
		"em_shared_keys_proposal": prop,
	})
	sig, _ := crypto.Sign(string(txRaw), pv)
	txSigned, _ := json.Marshal(map[string]interface{}{
		"addr":                    crypto.HashString("addr"),
		"data":                    data,
		"timestamp":               time.Now().Unix(),
		"type":                    chain.KeychainTransactionType,
		"public_key":              pub,
		"em_shared_keys_proposal": prop,
		"signature":               sig,
	})
	emSig, _ := crypto.Sign(string(txSigned), pv)
	txEmSigned, _ := json.Marshal(map[string]interface{}{
		"addr":                    crypto.HashString("addr"),
		"data":                    data,
		"timestamp":               time.Now().Unix(),
		"type":                    chain.KeychainTransactionType,
		"public_key":              pub,
		"em_shared_keys_proposal": prop,
		"signature":               sig,
		"em_signature":            emSig,
	})
	tx, err := chain.NewTransaction(crypto.HashString("addr"), chain.KeychainTransactionType, data, time.Now(), pub, prop, sig, emSig, crypto.HashBytes(txEmSigned))
	assert.Nil(t, err)

	mv, err := preValidateTransaction(tx, Pool{}, 1, pub, pv, emReader)
	assert.Nil(t, err)
	assert.Equal(t, pub, mv.ProofOfWork())
	assert.Equal(t, pub, mv.Validation().PublicKey())
	assert.Equal(t, chain.ValidationOK, mv.Validation().Status())
	ok, err := mv.Validation().IsValid()
	assert.True(t, ok)
	assert.Nil(t, err)
}

/*
Scenario: Lead transaction mining
	Given a valid transaction
	When I want to lead its mining
	Then the transaction is mined and stored
*/
func TestLeadMining(t *testing.T) {

	pub, pv := crypto.GenerateKeys()
	emReader := &mockEmitterReader{}
	emKP, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("pvKey")), pub)
	emReader.emKeys = append(emReader.emKeys, emKP)

	prop, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("pvKey")), pub)

	data := map[string]string{
		"encrypted_aes_key":         hex.EncodeToString([]byte("aesKey")),
		"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
		"encrypted_address_by_id":   hex.EncodeToString([]byte("addr")),
	}

	txRaw, _ := json.Marshal(map[string]interface{}{
		"addr":                    crypto.HashString("addr"),
		"data":                    data,
		"timestamp":               time.Now().Unix(),
		"type":                    chain.KeychainTransactionType,
		"public_key":              pub,
		"em_shared_keys_proposal": prop,
	})
	sig, _ := crypto.Sign(string(txRaw), pv)
	txSigned, _ := json.Marshal(map[string]interface{}{
		"addr":                    crypto.HashString("addr"),
		"data":                    data,
		"timestamp":               time.Now().Unix(),
		"type":                    chain.KeychainTransactionType,
		"public_key":              pub,
		"em_shared_keys_proposal": prop,
		"signature":               sig,
	})
	emSig, _ := crypto.Sign(string(txSigned), pv)
	txEmSigned, _ := json.Marshal(map[string]interface{}{
		"addr":                    crypto.HashString("addr"),
		"data":                    data,
		"timestamp":               time.Now().Unix(),
		"type":                    chain.KeychainTransactionType,
		"public_key":              pub,
		"em_shared_keys_proposal": prop,
		"signature":               sig,
		"em_signature":            emSig,
	})
	tx, err := chain.NewTransaction(crypto.HashString("addr"), chain.KeychainTransactionType, data, time.Now(), pub, prop, sig, emSig, crypto.HashBytes(txEmSigned))
	assert.Nil(t, err)
	poolR := &mockPoolRequester{}
	assert.Nil(t, LeadMining(tx, 1, poolR, pub, pv, emReader))

	time.Sleep(1 * time.Second)

	assert.Len(t, poolR.stores, 1)
}

/*
Scenario: Find pool for transaction mining
	Given a transaction
	When I want to find the pools
	Then I get the last validation pool, the validation pool and the storage pool
*/
func TestFindPools(t *testing.T) {
	poolR := &mockPoolRequester{}
	lastVPool, validPool, storagePool, err := findPools(chain.Transaction{}, poolR)

	assert.Nil(t, err)
	assert.Empty(t, lastVPool)
	assert.Equal(t, "127.0.0.1", validPool[0].IP().String())
	assert.Equal(t, "127.0.0.1", storagePool[0].IP().String())
}

type mockPoolRequester struct {
	stores []chain.Transaction
	ko     []chain.Transaction
}

func (pr mockPoolRequester) RequestLastTransaction(pool Pool, txAddr string, txType chain.TransactionType) (*chain.Transaction, error) {
	return nil, nil
}

func (pr mockPoolRequester) RequestTransactionLock(pool Pool, txHash string, txAddr string, masterPublicKey string) error {
	return nil
}

func (pr mockPoolRequester) RequestTransactionValidations(pool Pool, tx chain.Transaction, minValid int, masterValid chain.MasterValidation) ([]chain.Validation, error) {
	pub, pv := crypto.GenerateKeys()

	v, _ := buildValidation(chain.ValidationOK, pub, pv)
	return []chain.Validation{v}, nil
}

func (pr *mockPoolRequester) RequestTransactionStorage(pool Pool, minReplicas int, tx chain.Transaction) error {
	pr.stores = append(pr.stores, tx)
	return nil
}

type mockEmitterReader struct {
	emKeys shared.EmitterKeys
}

func (r mockEmitterReader) EmitterKeys() (shared.EmitterKeys, error) {
	return r.emKeys, nil
}

/*
Scenario: Get the minimum validation number
	Given a transaction hash
	When I want to get the validation required number
	Then I get a number  valid
	//TODO: to improve when the implementation will be defined
*/
func TestGetMinimumTransactionValidation(t *testing.T) {
	assert.Equal(t, 1, GetMinimumValidation(""))
}

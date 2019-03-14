package consensus

import (
	"crypto/rand"
	"log"
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
	pv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
	pubB, _ := pub.Marshal()

	v, _ := buildValidation(chain.ValidationOK, pub, pv)
	wHeaders := []chain.NodeHeader{chain.NewNodeHeader(pub, false, false, 0, true)}
	vHeaders := []chain.NodeHeader{chain.NewNodeHeader(pub, false, false, 0, true)}
	sHeaders := []chain.NodeHeader{chain.NewNodeHeader(pub, false, false, 0, true)}
	mv, _ := chain.NewMasterValidation([]crypto.PublicKey{}, pub, v, wHeaders, vHeaders, sHeaders)

	prop, _ := shared.NewEmitterCrossKeyPair([]byte("pvkey"), pub)

	data := map[string][]byte{
		"encrypted_aes_key":         []byte("aesKey"),
		"encrypted_address_by_node": []byte("addr"),
		"encrypted_address_by_id":   []byte("addr"),
	}

	txRaw, _ := json.Marshal(map[string]interface{}{
		"addr": hex.EncodeToString(crypto.Hash([]byte("addr"))),
		"data": map[string]string{
			"encrypted_aes_key":         hex.EncodeToString([]byte("aesKey")),
			"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
			"encrypted_address_by_id":   hex.EncodeToString([]byte("addr")),
		},
		"timestamp":  time.Now().Unix(),
		"type":       chain.KeychainTransactionType,
		"public_key": hex.EncodeToString(pubB),
		"em_shared_keys_proposal": map[string]string{
			"encrypted_private_key": hex.EncodeToString([]byte("pvkey")),
			"public_key":            hex.EncodeToString(pubB),
		},
	})
	sig, _ := pv.Sign(txRaw)
	txSigned, _ := json.Marshal(map[string]interface{}{
		"addr": hex.EncodeToString(crypto.Hash([]byte("addr"))),
		"data": map[string]string{
			"encrypted_aes_key":         hex.EncodeToString([]byte("aesKey")),
			"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
			"encrypted_address_by_id":   hex.EncodeToString([]byte("addr")),
		},
		"timestamp":  time.Now().Unix(),
		"type":       chain.KeychainTransactionType,
		"public_key": hex.EncodeToString(pubB),
		"em_shared_keys_proposal": map[string]string{
			"encrypted_private_key": hex.EncodeToString([]byte("pvkey")),
			"public_key":            hex.EncodeToString(pubB),
		},
		"signature": hex.EncodeToString(sig),
	})
	emSig, _ := pv.Sign(txSigned)
	txEmSigned, _ := json.Marshal(map[string]interface{}{
		"addr": hex.EncodeToString(crypto.Hash([]byte("addr"))),
		"data": map[string]string{
			"encrypted_aes_key":         hex.EncodeToString([]byte("aesKey")),
			"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
			"encrypted_address_by_id":   hex.EncodeToString([]byte("addr")),
		},
		"timestamp":  time.Now().Unix(),
		"type":       chain.KeychainTransactionType,
		"public_key": hex.EncodeToString(pubB),
		"em_shared_keys_proposal": map[string]string{
			"encrypted_private_key": hex.EncodeToString([]byte("pvkey")),
			"public_key":            hex.EncodeToString(pubB),
		},
		"signature":    hex.EncodeToString(sig),
		"em_signature": hex.EncodeToString(emSig),
	})
	tx, err := chain.NewTransaction(crypto.Hash([]byte("addr")), chain.KeychainTransactionType, data, time.Now(), pub, prop, sig, emSig, crypto.Hash(txEmSigned))

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
	pv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

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
	pv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
	pubB, _ := pub.Marshal()

	prop, _ := shared.NewEmitterCrossKeyPair([]byte("pvkey"), pub)

	data := map[string][]byte{
		"encrypted_aes_key":         []byte("aesKey"),
		"encrypted_address_by_node": []byte("addr"),
		"encrypted_address_by_id":   []byte("addr"),
	}

	txRaw, _ := json.Marshal(map[string]interface{}{
		"addr": hex.EncodeToString(crypto.Hash([]byte("addr"))),
		"data": map[string]string{
			"encrypted_aes_key":         hex.EncodeToString([]byte("aesKey")),
			"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
			"encrypted_address_by_id":   hex.EncodeToString([]byte("addr")),
		},
		"timestamp":  time.Now().Unix(),
		"type":       chain.KeychainTransactionType,
		"public_key": hex.EncodeToString(pubB),
		"em_shared_keys_proposal": map[string]string{
			"encrypted_private_key": hex.EncodeToString([]byte("pvkey")),
			"public_key":            hex.EncodeToString(pubB),
		},
	})

	log.Print(string(txRaw))

	sig, _ := pv.Sign(txRaw)
	txSigned, _ := json.Marshal(map[string]interface{}{
		"addr": hex.EncodeToString(crypto.Hash([]byte("addr"))),
		"data": map[string]string{
			"encrypted_aes_key":         hex.EncodeToString([]byte("aesKey")),
			"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
			"encrypted_address_by_id":   hex.EncodeToString([]byte("addr")),
		},
		"timestamp":  time.Now().Unix(),
		"type":       chain.KeychainTransactionType,
		"public_key": hex.EncodeToString(pubB),
		"em_shared_keys_proposal": map[string]string{
			"encrypted_private_key": hex.EncodeToString([]byte("pvkey")),
			"public_key":            hex.EncodeToString(pubB),
		},
		"signature": hex.EncodeToString(sig),
	})
	emSig, _ := pv.Sign(txSigned)
	txEmSigned, _ := json.Marshal(map[string]interface{}{
		"addr": hex.EncodeToString(crypto.Hash([]byte("addr"))),
		"data": map[string]string{
			"encrypted_aes_key":         hex.EncodeToString([]byte("aesKey")),
			"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
			"encrypted_address_by_id":   hex.EncodeToString([]byte("addr")),
		},
		"timestamp":  time.Now().Unix(),
		"type":       chain.KeychainTransactionType,
		"public_key": hex.EncodeToString(pubB),
		"em_shared_keys_proposal": map[string]string{
			"encrypted_private_key": hex.EncodeToString([]byte("pvkey")),
			"public_key":            hex.EncodeToString(pubB),
		},
		"signature":    hex.EncodeToString(sig),
		"em_signature": hex.EncodeToString(emSig),
	})
	tx, err := chain.NewTransaction(crypto.Hash([]byte("addr")), chain.KeychainTransactionType, data, time.Now(), pub, prop, sig, emSig, crypto.Hash(txEmSigned))
	assert.Nil(t, err)

	v, _ := buildValidation(chain.ValidationOK, pub, pv)
	wHeaders := []chain.NodeHeader{chain.NewNodeHeader(pub, false, false, 0, true)}
	vHeaders := []chain.NodeHeader{chain.NewNodeHeader(pub, false, false, 0, true)}
	sHeaders := []chain.NodeHeader{chain.NewNodeHeader(pub, false, false, 0, true)}
	mv, _ := chain.NewMasterValidation([]crypto.PublicKey{}, pub, v, wHeaders, vHeaders, sHeaders)

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
	pv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

	prop, _ := shared.NewEmitterCrossKeyPair([]byte("pvkey"), pub)

	data := map[string][]byte{
		"encrypted_aes_key":         []byte("aesKey"),
		"encrypted_address_by_node": []byte("addr"),
		"encrypted_address_by_id":   []byte("addr"),
	}

	sig, _ := pv.Sign([]byte("tx"))
	tx, _ := chain.NewTransaction(crypto.Hash([]byte("addr")), chain.IDTransactionType, data, time.Now(), pub, prop, sig, sig, crypto.Hash([]byte("hash")))

	v, _ := buildValidation(chain.ValidationOK, pub, pv)
	wHeaders := []chain.NodeHeader{chain.NewNodeHeader(pub, false, false, 0, true)}
	vHeaders := []chain.NodeHeader{chain.NewNodeHeader(pub, false, false, 0, true)}
	sHeaders := []chain.NodeHeader{chain.NewNodeHeader(pub, false, false, 0, true)}
	mv, _ := chain.NewMasterValidation([]crypto.PublicKey{}, pub, v, wHeaders, vHeaders, sHeaders)
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

	pv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
	pubB, _ := pub.Marshal()

	keyReader := &mockSharedKeyReader{}
	emKP, _ := shared.NewEmitterCrossKeyPair([]byte("pvKey"), pub)
	keyReader.crossEmitterKeys = append(keyReader.crossEmitterKeys, emKP)

	prop, _ := shared.NewEmitterCrossKeyPair([]byte("pvkey"), pub)

	data := map[string][]byte{
		"encrypted_aes_key":         []byte("aesKey"),
		"encrypted_address_by_node": []byte("addr"),
		"encrypted_address_by_id":   []byte("addr"),
	}

	txRaw, _ := json.Marshal(map[string]interface{}{
		"addr": hex.EncodeToString(crypto.Hash([]byte("addr"))),
		"data": map[string]string{
			"encrypted_aes_key":         hex.EncodeToString([]byte("aesKey")),
			"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
			"encrypted_address_by_id":   hex.EncodeToString([]byte("addr")),
		},
		"timestamp":  time.Now().Unix(),
		"type":       chain.KeychainTransactionType,
		"public_key": hex.EncodeToString(pubB),
		"em_shared_keys_proposal": map[string]string{
			"encrypted_private_key": hex.EncodeToString([]byte("pvkey")),
			"public_key":            hex.EncodeToString(pubB),
		},
	})
	sig, _ := pv.Sign(txRaw)
	txSigned, _ := json.Marshal(map[string]interface{}{
		"addr": hex.EncodeToString(crypto.Hash([]byte("addr"))),
		"data": map[string]string{
			"encrypted_aes_key":         hex.EncodeToString([]byte("aesKey")),
			"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
			"encrypted_address_by_id":   hex.EncodeToString([]byte("addr")),
		},
		"timestamp":  time.Now().Unix(),
		"type":       chain.KeychainTransactionType,
		"public_key": hex.EncodeToString(pubB),
		"em_shared_keys_proposal": map[string]string{
			"encrypted_private_key": hex.EncodeToString([]byte("pvkey")),
			"public_key":            hex.EncodeToString(pubB),
		},
		"signature": hex.EncodeToString(sig),
	})
	emSig, _ := pv.Sign(txSigned)
	txEmSigned, _ := json.Marshal(map[string]interface{}{
		"addr": hex.EncodeToString(crypto.Hash([]byte("addr"))),
		"data": map[string]string{
			"encrypted_aes_key":         hex.EncodeToString([]byte("aesKey")),
			"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
			"encrypted_address_by_id":   hex.EncodeToString([]byte("addr")),
		},
		"timestamp":  time.Now().Unix(),
		"type":       chain.KeychainTransactionType,
		"public_key": hex.EncodeToString(pubB),
		"em_shared_keys_proposal": map[string]string{
			"encrypted_private_key": hex.EncodeToString([]byte("pvkey")),
			"public_key":            hex.EncodeToString(pubB),
		},
		"signature":    hex.EncodeToString(sig),
		"em_signature": hex.EncodeToString(emSig),
	})
	tx, err := chain.NewTransaction(crypto.Hash([]byte("addr")), chain.KeychainTransactionType, data, time.Now(), pub, prop, sig, emSig, crypto.Hash(txEmSigned))
	assert.Nil(t, err)

	pow, err := proofOfWork(tx, keyReader)
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

	pv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
	pubB, _ := pub.Marshal()
	sharedKeyReader := &mockSharedKeyReader{}
	emKP, _ := shared.NewEmitterCrossKeyPair([]byte("pvkey"), pub)
	sharedKeyReader.crossEmitterKeys = append(sharedKeyReader.crossEmitterKeys, emKP)

	prop, _ := shared.NewEmitterCrossKeyPair([]byte("pvkey"), pub)

	data := map[string][]byte{
		"encrypted_aes_key":         []byte("aesKey"),
		"encrypted_address_by_node": []byte("addr"),
		"encrypted_address_by_id":   []byte("addr"),
	}

	txRaw, _ := json.Marshal(map[string]interface{}{
		"addr": hex.EncodeToString(crypto.Hash([]byte("addr"))),
		"data": map[string]string{
			"encrypted_aes_key":         hex.EncodeToString([]byte("aesKey")),
			"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
			"encrypted_address_by_id":   hex.EncodeToString([]byte("addr")),
		},
		"timestamp":  time.Now().Unix(),
		"type":       chain.KeychainTransactionType,
		"public_key": hex.EncodeToString(pubB),
		"em_shared_keys_proposal": map[string]string{
			"encrypted_private_key": hex.EncodeToString([]byte("pvkey")),
			"public_key":            hex.EncodeToString(pubB),
		},
	})
	sig, _ := pv.Sign(txRaw)
	txSigned, _ := json.Marshal(map[string]interface{}{
		"addr": hex.EncodeToString(crypto.Hash([]byte("addr"))),
		"data": map[string]string{
			"encrypted_aes_key":         hex.EncodeToString([]byte("aesKey")),
			"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
			"encrypted_address_by_id":   hex.EncodeToString([]byte("addr")),
		},
		"timestamp":  time.Now().Unix(),
		"type":       chain.KeychainTransactionType,
		"public_key": hex.EncodeToString(pubB),
		"em_shared_keys_proposal": map[string]string{
			"encrypted_private_key": hex.EncodeToString([]byte("pvkey")),
			"public_key":            hex.EncodeToString(pubB),
		},
		"signature": hex.EncodeToString(sig),
	})
	emSig, _ := pv.Sign(txSigned)
	txEmSigned, _ := json.Marshal(map[string]interface{}{
		"addr": hex.EncodeToString(crypto.Hash([]byte("addr"))),
		"data": map[string]string{
			"encrypted_aes_key":         hex.EncodeToString([]byte("aesKey")),
			"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
			"encrypted_address_by_id":   hex.EncodeToString([]byte("addr")),
		},
		"timestamp":  time.Now().Unix(),
		"type":       chain.KeychainTransactionType,
		"public_key": hex.EncodeToString(pubB),
		"em_shared_keys_proposal": map[string]string{
			"encrypted_private_key": hex.EncodeToString([]byte("pvkey")),
			"public_key":            hex.EncodeToString(pubB),
		},
		"signature":    hex.EncodeToString(sig),
		"em_signature": hex.EncodeToString(emSig),
	})
	tx, err := chain.NewTransaction(crypto.Hash([]byte("addr")), chain.KeychainTransactionType, data, time.Now(), pub, prop, sig, emSig, crypto.Hash(txEmSigned))
	assert.Nil(t, err)

	wHeaders := []chain.NodeHeader{chain.NewNodeHeader(pub, false, false, 0, true)}
	mv, err := preValidateTransaction(tx, wHeaders, Pool{Node{publicKey: pub}}, Pool{Node{publicKey: pub}}, Pool{}, 1, pub, pv, sharedKeyReader)
	assert.Nil(t, err)
	assert.Equal(t, pub, mv.ProofOfWork())
	assert.EqualValues(t, pub, mv.Validation().PublicKey())
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

	pv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
	pubB, _ := pub.Marshal()
	sharedKeyReader := &mockSharedKeyReader{}
	emKP, _ := shared.NewEmitterCrossKeyPair([]byte("pvkey"), pub)
	sharedKeyReader.crossEmitterKeys = append(sharedKeyReader.crossEmitterKeys, emKP)

	prop, _ := shared.NewEmitterCrossKeyPair([]byte("pvkey"), pub)

	data := map[string][]byte{
		"encrypted_aes_key":         []byte("aesKey"),
		"encrypted_address_by_node": []byte("addr"),
		"encrypted_address_by_id":   []byte("addr"),
	}

	txRaw, _ := json.Marshal(map[string]interface{}{
		"addr": hex.EncodeToString(crypto.Hash([]byte("addr"))),
		"data": map[string]string{
			"encrypted_aes_key":         hex.EncodeToString([]byte("aesKey")),
			"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
			"encrypted_address_by_id":   hex.EncodeToString([]byte("addr")),
		},
		"timestamp":  time.Now().Unix(),
		"type":       chain.KeychainTransactionType,
		"public_key": hex.EncodeToString(pubB),
		"em_shared_keys_proposal": map[string]string{
			"encrypted_private_key": hex.EncodeToString([]byte("pvkey")),
			"public_key":            hex.EncodeToString(pubB),
		},
	})
	sig, _ := pv.Sign(txRaw)
	txSigned, _ := json.Marshal(map[string]interface{}{
		"addr": hex.EncodeToString(crypto.Hash([]byte("addr"))),
		"data": map[string]string{
			"encrypted_aes_key":         hex.EncodeToString([]byte("aesKey")),
			"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
			"encrypted_address_by_id":   hex.EncodeToString([]byte("addr")),
		},
		"timestamp":  time.Now().Unix(),
		"type":       chain.KeychainTransactionType,
		"public_key": hex.EncodeToString(pubB),
		"em_shared_keys_proposal": map[string]string{
			"encrypted_private_key": hex.EncodeToString([]byte("pvkey")),
			"public_key":            hex.EncodeToString(pubB),
		},
		"signature": hex.EncodeToString(sig),
	})
	emSig, _ := pv.Sign(txSigned)
	txEmSigned, _ := json.Marshal(map[string]interface{}{
		"addr": hex.EncodeToString(crypto.Hash([]byte("addr"))),
		"data": map[string]string{
			"encrypted_aes_key":         hex.EncodeToString([]byte("aesKey")),
			"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
			"encrypted_address_by_id":   hex.EncodeToString([]byte("addr")),
		},
		"timestamp":  time.Now().Unix(),
		"type":       chain.KeychainTransactionType,
		"public_key": hex.EncodeToString(pubB),
		"em_shared_keys_proposal": map[string]string{
			"encrypted_private_key": hex.EncodeToString([]byte("pvkey")),
			"public_key":            hex.EncodeToString(pubB),
		},
		"signature":    hex.EncodeToString(sig),
		"em_signature": hex.EncodeToString(emSig),
	})
	tx, err := chain.NewTransaction(crypto.Hash([]byte("addr")), chain.KeychainTransactionType, data, time.Now(), pub, prop, sig, emSig, crypto.Hash(txEmSigned))
	assert.Nil(t, err)
	poolR := &mockPoolRequester{}
	wHeaders := []chain.NodeHeader{chain.NewNodeHeader(pub, false, false, 0, true)}
	assert.Nil(t, LeadMining(tx, 1, wHeaders, poolR, pub, pv, sharedKeyReader))

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

	pv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

	prop, _ := shared.NewEmitterCrossKeyPair([]byte("encPvKey"), pub)

	addr := crypto.Hash([]byte("addr"))
	data := map[string][]byte{
		"encrypted_address_by_node": []byte("addr"),
		"encrypted_address_by_id":   []byte("addr"),
		"encrypted_aes_key":         []byte("aesKey"),
	}
	hash := crypto.Hash([]byte("hash"))

	sig, _ := pv.Sign([]byte("data"))

	tx, _ := chain.NewTransaction(addr, chain.KeychainTransactionType, data, time.Now(), pub, prop, sig, sig, hash)

	lastVPool, validPool, storagePool, err := findPools(tx, poolR)

	assert.Nil(t, err)
	assert.Empty(t, lastVPool)
	assert.Equal(t, "127.0.0.1", validPool[0].IP().String())
	assert.Equal(t, "127.0.0.1", storagePool[0].IP().String())
}

type mockPoolRequester struct {
	stores []chain.Transaction
	ko     []chain.Transaction
}

func (pr mockPoolRequester) RequestLastTransaction(pool Pool, txAddr crypto.VersionnedHash, txType chain.TransactionType) (*chain.Transaction, error) {
	return nil, nil
}

func (pr mockPoolRequester) RequestTransactionTimeLock(pool Pool, txHash crypto.VersionnedHash, txAddr crypto.VersionnedHash, masterPublicKey crypto.PublicKey) error {
	return nil
}

func (pr mockPoolRequester) RequestTransactionValidations(pool Pool, tx chain.Transaction, minValid int, masterValid chain.MasterValidation) ([]chain.Validation, error) {
	pv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

	v, _ := buildValidation(chain.ValidationOK, pub, pv)
	return []chain.Validation{v}, nil
}

func (pr *mockPoolRequester) RequestTransactionStorage(pool Pool, minReplicas int, tx chain.Transaction) error {
	pr.stores = append(pr.stores, tx)
	return nil
}

/*
Scenario: Get the minimum validation number
	Given a transaction hash
	When I want to get the validation required number
	Then I get a number  valid
	//TODO: to improve when the implementation will be defined
*/
func TestGetMinimumTransactionValidation(t *testing.T) {
	assert.Equal(t, 1, GetMinimumValidation([]byte("")))
}

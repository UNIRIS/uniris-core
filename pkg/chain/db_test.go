package chain

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/uniris/uniris-core/pkg/crypto"
	"github.com/uniris/uniris-core/pkg/shared"
)

/*
Scenario: Get transaction keychain by its hash
	Given a keychain tx stored
	When I dbant to retrieve the transaction by only its hash
	Then I can get it
*/
func TestReadKeychainByHash(t *testing.T) {
	chainDB := &mockChainDB{}

	pv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

	sig, _ := pv.Sign([]byte("hello"))

	prop, _ := shared.NewEmitterCrossKeyPair([]byte("pvkey"), pub)
	data := map[string][]byte{
		"encrypted_address_by_node": []byte("addr"),
		"encrypted_wallet":          []byte("wallet"),
	}
	tx := Transaction{
		addr:      crypto.Hash([]byte("addr")),
		data:      data,
		txType:    KeychainTransactionType,
		timestamp: time.Now(),
		pubKey:    pub,
		prop:      prop,
	}
	txBytesBeforeSig, _ := tx.MarshalBeforeSignature()
	sig, _ = pv.Sign(txBytesBeforeSig)
	tx.emSig = sig
	tx.sig = sig
	txBytes, _ := tx.MarshalHash()
	hash := crypto.Hash(txBytes)
	tx.hash = hash
	vBytes, _ := json.Marshal(Validation{
		nodePubk:  pub,
		status:    ValidationOK,
		timestamp: time.Now(),
	})
	vSig, _ := pv.Sign(vBytes)
	v, _ := NewValidation(ValidationOK, time.Now(), pub, vSig)
	wHeaders := []NodeHeader{NewNodeHeader(pub, false, false, 0, true)}
	vHeaders := []NodeHeader{NewNodeHeader(pub, false, false, 0, true)}
	sHeaders := []NodeHeader{NewNodeHeader(pub, false, false, 0, true)}
	mv, _ := NewMasterValidation([]crypto.PublicKey{}, pub, v, wHeaders, vHeaders, sHeaders)

	tx.Mined(mv, []Validation{v})
	keychain, err := NewKeychain(tx)
	chainDB.keychains = append(chainDB.keychains, keychain)

	txKeychain, err := getTransactionByHash(chainDB, tx.hash)
	assert.Nil(t, err)
	assert.Equal(t, KeychainTransactionType, txKeychain.txType)
	assert.Equal(t, pub, txKeychain.PublicKey())
}

/*
Scenario: Get transaction ID by its hash
	Given a ID tx stored
	When I dbant to retrieve the transaction by only its hash
	Then I can get it
*/
func TestGetIDTransactionByHash(t *testing.T) {
	chainDB := &mockChainDB{}

	pv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

	sig, _ := pv.Sign([]byte("hello"))

	prop, _ := shared.NewEmitterCrossKeyPair([]byte("pvkey"), pub)

	data := map[string][]byte{
		"encrypted_address_by_node": []byte("addr"),
		"encrypted_address_by_id":   []byte("addr"),
		"encrypted_aes_key":         []byte("aesKey"),
	}

	tx := Transaction{
		addr:      crypto.Hash([]byte("addr")),
		data:      data,
		txType:    IDTransactionType,
		timestamp: time.Now(),
		pubKey:    pub,
		prop:      prop,
	}
	txBytesBeforeSig, _ := tx.MarshalBeforeSignature()
	sig, _ = pv.Sign(txBytesBeforeSig)
	tx.emSig = sig
	tx.sig = sig
	txBytes, _ := tx.MarshalHash()
	hash := crypto.Hash(txBytes)
	tx.hash = hash
	vBytes, _ := json.Marshal(Validation{
		nodePubk:  pub,
		status:    ValidationOK,
		timestamp: time.Now(),
	})
	vSig, _ := pv.Sign(vBytes)
	v, _ := NewValidation(ValidationOK, time.Now(), pub, vSig)
	wHeaders := []NodeHeader{NewNodeHeader(pub, false, false, 0, true)}
	vHeaders := []NodeHeader{NewNodeHeader(pub, false, false, 0, true)}
	sHeaders := []NodeHeader{NewNodeHeader(pub, false, false, 0, true)}
	mv, _ := NewMasterValidation([]crypto.PublicKey{}, pub, v, wHeaders, vHeaders, sHeaders)

	tx.Mined(mv, []Validation{v})
	id, _ := NewID(tx)

	chainDB.ids = append(chainDB.ids, id)

	txID, err := getTransactionByHash(chainDB, tx.hash)
	assert.Nil(t, err)
	assert.Equal(t, IDTransactionType, txID.txType)
	assert.Equal(t, pub, txID.PublicKey())
}

/*
Scenario: Get unknown transaction by its hash
	Given no tx stored
	When I dbant to retrieve the transaction by only its hash
	Then I can get an error
*/
func TestGetUnknodbnTransactionByHash(t *testing.T) {
	chainDB := &mockChainDB{}

	_, err := getTransactionByHash(chainDB, crypto.Hash([]byte("hash")))
	assert.EqualError(t, err, "unknown transaction")
}

/*
Scenario: Get transaction status in progress
	Given a transaction stored in in progress
	When I dbant to get its status
	Then I get in progress
*/
func TestGetTransactionStatusInProgress(t *testing.T) {
	chainDB := &mockChainDB{}
	_, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
	TimeLockTransaction(crypto.Hash([]byte("hash")), crypto.Hash([]byte("addr")), pub)

	status, err := GetTransactionStatus(chainDB, crypto.Hash([]byte("hash")))
	assert.Nil(t, err)
	assert.Equal(t, TransactionStatusInProgress, status)
	removeTimeLock(crypto.Hash([]byte("hash")), crypto.Hash([]byte("addr")))
}

/*
Scenario: Get transaction status KO
	Given a transaction stored in KO
	When I dbant to get its status
	Then I get failure
*/
func TestGetTransactionStatusFailure(t *testing.T) {
	chainDB := &mockChainDB{
		kos: []Transaction{
			Transaction{
				hash: crypto.Hash([]byte("hash")),
			},
		},
	}

	status, err := GetTransactionStatus(chainDB, crypto.Hash([]byte("hash")))
	assert.Nil(t, err)
	assert.Equal(t, TransactionStatusFailure, status)
}

/*
Scenario: Get transaction status unknown
	Given a transaction stored in KO
	When I dbant to get its status
	Then I get failure
*/
func TestGetTransactionStatusUnknown(t *testing.T) {
	chainDB := &mockChainDB{}

	status, err := GetTransactionStatus(chainDB, crypto.Hash([]byte("hash")))
	assert.Nil(t, err)
	assert.Equal(t, TransactionStatusUnknown, status)
}

/*
Scenario: Get transaction status success
	Given a transaction stored
	When I dbant to get its status
	Then I get success
*/
func TestGetTransactionStatusSuccess(t *testing.T) {
	chainDB := &mockChainDB{}

	pv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

	sig, _ := pv.Sign([]byte("hello"))

	prop, _ := shared.NewEmitterCrossKeyPair([]byte("pvkey"), pub)

	data := map[string][]byte{
		"encrypted_address_by_node": []byte("addr"),
		"encrypted_address_by_id":   []byte("addr"),
		"encrypted_aes_key":         []byte("aesKey"),
	}

	tx := Transaction{
		addr:      crypto.Hash([]byte("addr")),
		data:      data,
		txType:    IDTransactionType,
		timestamp: time.Now(),
		pubKey:    pub,
		prop:      prop,
	}
	txBytesBeforeSig, _ := tx.MarshalBeforeSignature()
	sig, _ = pv.Sign(txBytesBeforeSig)
	tx.emSig = sig
	tx.sig = sig
	txBytes, _ := tx.MarshalHash()
	hash := crypto.Hash(txBytes)
	tx.hash = hash
	vBytes, _ := json.Marshal(Validation{
		nodePubk:  pub,
		status:    ValidationOK,
		timestamp: time.Now(),
	})
	vSig, _ := pv.Sign(vBytes)
	v, _ := NewValidation(ValidationOK, time.Now(), pub, vSig)
	wHeaders := []NodeHeader{NewNodeHeader(pub, false, false, 0, true)}
	vHeaders := []NodeHeader{NewNodeHeader(pub, false, false, 0, true)}
	sHeaders := []NodeHeader{NewNodeHeader(pub, false, false, 0, true)}
	mv, _ := NewMasterValidation([]crypto.PublicKey{}, pub, v, wHeaders, vHeaders, sHeaders)

	tx.Mined(mv, []Validation{v})
	id, _ := NewID(tx)

	chainDB.ids = append(chainDB.ids, id)

	status, err := GetTransactionStatus(chainDB, tx.hash)
	assert.Nil(t, err)
	assert.Equal(t, TransactionStatusSuccess, status)
}

/*
Scenario: Get the last keychain transaction
	Given tdbo keychain transaction chained
	When I dbant to get the last
	Then I get the 2nd
*/
func TestReadLastKeychainTransaction(t *testing.T) {
	chainDB := &mockChainDB{}

	pv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

	sig, _ := pv.Sign([]byte("hello"))

	prop, _ := shared.NewEmitterCrossKeyPair([]byte("pvkey"), pub)

	data := map[string][]byte{
		"encrypted_address_by_node": []byte("addr"),
		"encrypted_wallet":          []byte("wallet"),
	}

	tx := Transaction{
		addr:      crypto.Hash([]byte("addr")),
		data:      data,
		txType:    KeychainTransactionType,
		timestamp: time.Now(),
		pubKey:    pub,
		prop:      prop,
	}
	txBytesBeforeSig, _ := tx.MarshalBeforeSignature()
	sig, _ = pv.Sign(txBytesBeforeSig)
	tx.emSig = sig
	tx.sig = sig
	txBytes, _ := tx.MarshalHash()
	hash := crypto.Hash(txBytes)
	tx.hash = hash

	vBytes, _ := json.Marshal(Validation{
		nodePubk:  pub,
		status:    ValidationOK,
		timestamp: time.Now(),
	})
	vSig, _ := pv.Sign(vBytes)
	v, _ := NewValidation(ValidationOK, time.Now(), pub, vSig)
	wHeaders := []NodeHeader{NewNodeHeader(pub, false, false, 0, true)}
	vHeaders := []NodeHeader{NewNodeHeader(pub, false, false, 0, true)}
	sHeaders := []NodeHeader{NewNodeHeader(pub, false, false, 0, true)}
	mv, _ := NewMasterValidation([]crypto.PublicKey{}, pub, v, wHeaders, vHeaders, sHeaders)

	tx.Mined(mv, []Validation{v})
	keychain1, _ := NewKeychain(tx)

	chainDB.keychains = append(chainDB.keychains, keychain1)

	tx2 := Transaction{
		addr:      crypto.Hash([]byte("addr")),
		data:      data,
		txType:    KeychainTransactionType,
		timestamp: time.Now().Add(2 * time.Second),
		pubKey:    pub,
		prop:      prop,
	}
	txBytesBeforeSig2, _ := tx2.MarshalBeforeSignature()
	sig2, _ := pv.Sign(txBytesBeforeSig2)
	tx2.emSig = sig2
	tx2.sig = sig2
	txBytes2, _ := tx2.MarshalHash()
	hash2 := crypto.Hash(txBytes2)
	tx2.hash = hash2

	tx2.Mined(mv, []Validation{v})
	keychain2, err := NewKeychain(tx2)

	chainDB.keychains = append(chainDB.keychains, keychain2)
	assert.Len(t, chainDB.keychains, 2)

	lastTx, err := LastTransaction(chainDB, crypto.Hash([]byte("addr")), KeychainTransactionType)
	assert.Nil(t, err)
	assert.NotNil(t, lastTx)
	assert.Equal(t, KeychainTransactionType, lastTx.txType)
	assert.Equal(t, hash2, lastTx.hash)

}

/*
Scenario: Get the last ID transaction
	Given tdbo ID transaction
	When I dbant to get the last
	Then I get the one I reached (because ID are not chained)
*/
func TestGetLastIDTransaction(t *testing.T) {
	chainDB := &mockChainDB{}

	pv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

	sig, _ := pv.Sign([]byte("hello"))

	prop, _ := shared.NewEmitterCrossKeyPair([]byte("pvkey"), pub)

	data := map[string][]byte{
		"encrypted_address_by_node": []byte("addr"),
		"encrypted_address_by_id":   []byte("addr"),
		"encrypted_aes_key":         []byte("aesKey"),
	}

	tx := Transaction{
		addr:      crypto.Hash([]byte("addr")),
		data:      data,
		txType:    IDTransactionType,
		timestamp: time.Now(),
		pubKey:    pub,
		prop:      prop,
	}
	txBytesBeforeSig, _ := tx.MarshalBeforeSignature()
	sig, _ = pv.Sign(txBytesBeforeSig)
	tx.emSig = sig
	tx.sig = sig
	txBytes, _ := tx.MarshalHash()
	hash := crypto.Hash(txBytes)
	tx.hash = hash

	vBytes, _ := json.Marshal(Validation{
		nodePubk:  pub,
		status:    ValidationOK,
		timestamp: time.Now(),
	})
	vSig, _ := pv.Sign(vBytes)
	v, _ := NewValidation(ValidationOK, time.Now(), pub, vSig)
	wHeaders := []NodeHeader{NewNodeHeader(pub, false, false, 0, true)}
	vHeaders := []NodeHeader{NewNodeHeader(pub, false, false, 0, true)}
	sHeaders := []NodeHeader{NewNodeHeader(pub, false, false, 0, true)}
	mv, _ := NewMasterValidation([]crypto.PublicKey{}, pub, v, wHeaders, vHeaders, sHeaders)

	tx.Mined(mv, []Validation{v})
	id1, _ := NewID(tx)

	chainDB.ids = append(chainDB.ids, id1)

	tx2 := Transaction{
		addr:      crypto.Hash([]byte("addr2")),
		data:      data,
		txType:    IDTransactionType,
		timestamp: time.Now().Add(2 * time.Second),
		pubKey:    pub,
		prop:      prop,
	}
	txBytesBeforeSig2, _ := tx.MarshalBeforeSignature()
	sig2, _ := pv.Sign(txBytesBeforeSig2)
	tx2.emSig = sig2
	tx2.sig = sig2
	txBytes2, _ := tx2.MarshalHash()
	hash2 := crypto.Hash(txBytes2)
	tx2.hash = hash2
	tx2.Mined(mv, []Validation{v})
	id2, err := NewID(tx2)

	chainDB.ids = append(chainDB.ids, id2)
	assert.Len(t, chainDB.ids, 2)

	lastTx, err := LastTransaction(chainDB, crypto.Hash([]byte("addr")), IDTransactionType)
	assert.Nil(t, err)
	assert.NotNil(t, lastTx)
	assert.Equal(t, IDTransactionType, lastTx.txType)
	assert.Equal(t, hash, lastTx.hash)

	lastTx, err = LastTransaction(chainDB, crypto.Hash([]byte("addr2")), IDTransactionType)
	assert.Nil(t, err)
	assert.NotNil(t, lastTx)
	assert.Equal(t, IDTransactionType, lastTx.txType)
	assert.Equal(t, hash2, lastTx.hash)

}

/*
Scenario: Get the chain of a transaction
	Given transactions chained
	When I dbant to get the chain
	Then I get recursively the transactions linked to each other
*/
func TestGetTransactionChain(t *testing.T) {
	chainDB := &mockChainDB{}

	pv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

	prop, _ := shared.NewEmitterCrossKeyPair([]byte("pvkey"), pub)

	data := map[string][]byte{
		"encrypted_address_by_node": []byte("addr"),
		"encrypted_wallet":          []byte("wallet"),
	}

	tx := Transaction{
		addr:      crypto.Hash([]byte("addr")),
		data:      data,
		txType:    KeychainTransactionType,
		timestamp: time.Now(),
		pubKey:    pub,
		prop:      prop,
	}
	txBytesBeforeSig, _ := tx.MarshalBeforeSignature()
	sig, _ := pv.Sign(txBytesBeforeSig)
	tx.emSig = sig
	tx.sig = sig
	txBytes, _ := tx.MarshalHash()
	hash := crypto.Hash(txBytes)
	tx.hash = hash

	vBytes, _ := json.Marshal(Validation{
		nodePubk:  pub,
		status:    ValidationOK,
		timestamp: time.Now(),
	})
	vSig, _ := pv.Sign(vBytes)
	v, _ := NewValidation(ValidationOK, time.Now(), pub, vSig)
	wHeaders := []NodeHeader{NewNodeHeader(pub, false, false, 0, true)}
	vHeaders := []NodeHeader{NewNodeHeader(pub, false, false, 0, true)}
	sHeaders := []NodeHeader{NewNodeHeader(pub, false, false, 0, true)}
	mv, _ := NewMasterValidation([]crypto.PublicKey{}, pub, v, wHeaders, vHeaders, sHeaders)

	tx.Mined(mv, []Validation{v})
	keychain1, _ := NewKeychain(tx)

	chainDB.keychains = append(chainDB.keychains, keychain1)

	tx2 := Transaction{
		addr:      crypto.Hash([]byte("addr")),
		data:      data,
		txType:    KeychainTransactionType,
		timestamp: time.Now().Add(2 * time.Second),
		pubKey:    pub,
		prop:      prop,
	}
	txBytesBeforeSig2, _ := tx2.MarshalBeforeSignature()
	sig2, _ := pv.Sign(txBytesBeforeSig2)
	tx2.emSig = sig2
	tx2.sig = sig2
	txBytes2, _ := tx2.MarshalHash()
	hash2 := crypto.Hash(txBytes2)
	tx2.hash = hash2

	tx2.Mined(mv, []Validation{v})
	keychain2, _ := NewKeychain(tx2)
	assert.Nil(t, keychain2.Chain(&tx))

	chainDB.keychains = append(chainDB.keychains, keychain2)
	assert.Len(t, chainDB.keychains, 2)

	chain, err := getFullChain(chainDB, crypto.Hash([]byte("addr")), KeychainTransactionType)
	assert.Nil(t, err)
	assert.NotNil(t, chain)
	assert.Equal(t, hash2, chain.hash)
	assert.NotNil(t, chain.prevTx)
	assert.Equal(t, hash, chain.prevTx.hash)
}

/*
Scenario: Check a valid transaction before store
	Given a transaction
	When I dbant to check before storage
	Then I get not error
*/
func TestCheckTransactionBeforeStore(t *testing.T) {
	pv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

	prop, _ := shared.NewEmitterCrossKeyPair([]byte("pvkey"), pub)

	data := map[string][]byte{
		"encrypted_address_by_node": []byte("addr"),
		"encrypted_wallet":          []byte("wallet"),
	}

	tx := Transaction{
		addr:      crypto.Hash([]byte("addr")),
		data:      data,
		txType:    KeychainTransactionType,
		timestamp: time.Now(),
		pubKey:    pub,
		prop:      prop,
	}
	txBytesBeforeSig, _ := tx.MarshalBeforeSignature()
	sig, _ := pv.Sign(txBytesBeforeSig)
	tx.sig = sig
	txBytesBeforeEmSig, _ := tx.MarshalBeforeEmitterSignature()
	emSig, _ := pv.Sign(txBytesBeforeEmSig)
	tx.emSig = emSig
	txBytes, _ := tx.MarshalHash()
	hash := crypto.Hash(txBytes)
	tx.hash = hash

	vBytes, _ := json.Marshal(Validation{
		nodePubk:  pub,
		status:    ValidationOK,
		timestamp: time.Now(),
	})
	vSig, _ := pv.Sign(vBytes)
	v, _ := NewValidation(ValidationOK, time.Now(), pub, vSig)
	wHeaders := []NodeHeader{NewNodeHeader(pub, false, false, 0, true)}
	vHeaders := []NodeHeader{NewNodeHeader(pub, false, false, 0, true)}
	sHeaders := []NodeHeader{NewNodeHeader(pub, false, false, 0, true)}
	mv, _ := NewMasterValidation([]crypto.PublicKey{}, pub, v, wHeaders, vHeaders, sHeaders)
	tx.Mined(mv, []Validation{v})

	assert.Nil(t, checkTransactionBeforeStorage(tx, 1))
}

/*
Scenario: Check a transaction before store dbith misssing validations
	Given a transaction dbith missing confirmations validations
	When I dbant to check before storage
	Then I get an error
*/
func TestCheckTransactionBeforeStoreWithMissingValidations(t *testing.T) {
	pv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

	prop, _ := shared.NewEmitterCrossKeyPair([]byte("pvkey"), pub)

	data := map[string][]byte{
		"encrypted_address_by_node": []byte("addr"),
		"encrypted_wallet":          []byte("wallet"),
	}

	tx := Transaction{
		addr:      crypto.Hash([]byte("addr")),
		data:      data,
		txType:    KeychainTransactionType,
		timestamp: time.Now(),
		pubKey:    pub,
		prop:      prop,
	}
	txBytesBeforeSig, _ := tx.MarshalBeforeSignature()
	sig, _ := pv.Sign(txBytesBeforeSig)
	tx.sig = sig
	txBytesBeforeEmSig, _ := tx.MarshalBeforeEmitterSignature()
	emSig, _ := pv.Sign(txBytesBeforeEmSig)
	tx.emSig = emSig
	txBytes, _ := tx.MarshalHash()
	hash := crypto.Hash(txBytes)
	tx.hash = hash

	vBytes, _ := json.Marshal(Validation{
		nodePubk:  pub,
		status:    ValidationOK,
		timestamp: time.Now(),
	})
	vSig, _ := pv.Sign(vBytes)
	v, _ := NewValidation(ValidationOK, time.Now(), pub, vSig)
	wHeaders := []NodeHeader{NewNodeHeader(pub, false, false, 0, true)}
	vHeaders := []NodeHeader{NewNodeHeader(pub, false, false, 0, true)}
	sHeaders := []NodeHeader{NewNodeHeader(pub, false, false, 0, true)}
	mv, _ := NewMasterValidation([]crypto.PublicKey{}, pub, v, wHeaders, vHeaders, sHeaders)
	tx.Mined(mv, []Validation{})

	assert.EqualError(t, checkTransactionBeforeStorage(tx, 1), "transaction: invalid number of validations")
}

/*
Scenario: Store a KO transaction
	Given a transaction not valid
	When I dbant to store it
	Then the transaction is stored on the KO db
*/
func TestStoreKOTransaction(t *testing.T) {

	pv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

	vBytes, _ := json.Marshal(Validation{
		nodePubk:  pub,
		status:    ValidationKO,
		timestamp: time.Now(),
	})
	vSig, _ := pv.Sign(vBytes)
	v, _ := NewValidation(ValidationKO, time.Now(), pub, vSig)

	vBytes2, _ := json.Marshal(Validation{
		nodePubk:  pub,
		status:    ValidationOK,
		timestamp: time.Now(),
	})
	vSig2, _ := pv.Sign(vBytes2)
	v2, _ := NewValidation(ValidationOK, time.Now(), pub, vSig2)
	wHeaders := []NodeHeader{NewNodeHeader(pub, false, false, 0, true)}
	vHeaders := []NodeHeader{NewNodeHeader(pub, false, false, 0, true)}
	sHeaders := []NodeHeader{NewNodeHeader(pub, false, false, 0, true)}
	mv, _ := NewMasterValidation([]crypto.PublicKey{}, pub, v, wHeaders, vHeaders, sHeaders)

	prop, _ := shared.NewEmitterCrossKeyPair([]byte("pvkey"), pub)

	data := map[string][]byte{
		"encrypted_address_by_node": []byte("addr"),
		"encrypted_wallet":          []byte("wallet"),
	}

	tx := Transaction{
		addr:      crypto.Hash([]byte("addr")),
		data:      data,
		txType:    KeychainTransactionType,
		timestamp: time.Now(),
		pubKey:    pub,
		prop:      prop,
	}
	txBytesBeforeSig, _ := tx.MarshalBeforeSignature()
	sig, _ := pv.Sign(txBytesBeforeSig)
	tx.sig = sig
	txBytesBeforeEmSig, _ := tx.MarshalBeforeEmitterSignature()
	emSig, _ := pv.Sign(txBytesBeforeEmSig)
	tx.emSig = emSig
	txBytes, _ := tx.MarshalHash()
	hash := crypto.Hash(txBytes)
	tx.hash = hash

	tx.Mined(mv, []Validation{v2})

	chainDB := &mockChainDB{}

	assert.Nil(t, WriteTransaction(chainDB, tx, 1))
	assert.Len(t, chainDB.kos, 1)
	assert.Equal(t, crypto.Hash(txBytes), chainDB.kos[0].hash)
}

/*
Scenario: Store a Keychain transaction
	Given a keychain transaction
	When I dbant to store it
	Then the transaction is stored on the keychain db
*/
func TestStoreKeychainTransaction(t *testing.T) {

	pv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

	vBytes, _ := json.Marshal(Validation{
		nodePubk:  pub,
		status:    ValidationOK,
		timestamp: time.Now(),
	})
	vSig, _ := pv.Sign(vBytes)
	v, _ := NewValidation(ValidationOK, time.Now(), pub, vSig)

	vBytes2, _ := json.Marshal(Validation{
		nodePubk:  pub,
		status:    ValidationOK,
		timestamp: time.Now(),
	})
	vSig2, _ := pv.Sign(vBytes2)
	v2, _ := NewValidation(ValidationOK, time.Now(), pub, vSig2)
	wHeaders := []NodeHeader{NewNodeHeader(pub, false, false, 0, true)}
	vHeaders := []NodeHeader{NewNodeHeader(pub, false, false, 0, true)}
	sHeaders := []NodeHeader{NewNodeHeader(pub, false, false, 0, true)}
	mv, _ := NewMasterValidation([]crypto.PublicKey{}, pub, v, wHeaders, vHeaders, sHeaders)

	prop, _ := shared.NewEmitterCrossKeyPair([]byte("pvkey"), pub)

	data := map[string][]byte{
		"encrypted_address_by_node": []byte("addr"),
		"encrypted_wallet":          []byte("wallet"),
	}

	tx := Transaction{
		addr:      crypto.Hash([]byte("addr")),
		data:      data,
		txType:    KeychainTransactionType,
		timestamp: time.Now(),
		pubKey:    pub,
		prop:      prop,
	}
	txBytesBeforeSig, _ := tx.MarshalBeforeSignature()
	sig, _ := pv.Sign(txBytesBeforeSig)
	tx.sig = sig
	txBytesBeforeEmSig, _ := tx.MarshalBeforeEmitterSignature()
	emSig, _ := pv.Sign(txBytesBeforeEmSig)
	tx.emSig = emSig
	txBytes, _ := tx.MarshalHash()
	hash := crypto.Hash(txBytes)
	tx.hash = hash

	tx.Mined(mv, []Validation{v2})

	chainDB := &mockChainDB{}

	assert.Nil(t, WriteTransaction(chainDB, tx, 1))
	assert.Len(t, chainDB.keychains, 1)
	assert.Equal(t, crypto.Hash(txBytes), chainDB.keychains[0].hash)
}

/*
Scenario: Store a ID transaction
	Given a ID transaction
	When I dbant to store it
	Then the transaction is stored on the ID db
*/
func TestStoreIDTransaction(t *testing.T) {

	pv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

	vBytes, _ := json.Marshal(Validation{
		nodePubk:  pub,
		status:    ValidationOK,
		timestamp: time.Now(),
	})
	vSig, _ := pv.Sign(vBytes)
	v, _ := NewValidation(ValidationOK, time.Now(), pub, vSig)

	vBytes2, _ := json.Marshal(Validation{
		nodePubk:  pub,
		status:    ValidationOK,
		timestamp: time.Now(),
	})
	vSig2, _ := pv.Sign(vBytes2)
	v2, _ := NewValidation(ValidationOK, time.Now(), pub, vSig2)
	wHeaders := []NodeHeader{NewNodeHeader(pub, false, false, 0, true)}
	vHeaders := []NodeHeader{NewNodeHeader(pub, false, false, 0, true)}
	sHeaders := []NodeHeader{NewNodeHeader(pub, false, false, 0, true)}
	mv, _ := NewMasterValidation([]crypto.PublicKey{}, pub, v, wHeaders, vHeaders, sHeaders)

	prop, _ := shared.NewEmitterCrossKeyPair([]byte("pvkey"), pub)

	data := map[string][]byte{
		"encrypted_address_by_node": []byte("addr"),
		"encrypted_address_by_id":   []byte("addr"),
		"encrypted_aes_key":         []byte("aesKey"),
	}

	tx := Transaction{
		addr:      crypto.Hash([]byte("addr")),
		data:      data,
		txType:    IDTransactionType,
		timestamp: time.Now(),
		pubKey:    pub,
		prop:      prop,
	}
	txBytesBeforeSig, _ := tx.MarshalBeforeSignature()
	sig, _ := pv.Sign(txBytesBeforeSig)
	tx.sig = sig
	txBytesBeforeEmSig, _ := tx.MarshalBeforeEmitterSignature()
	emSig, _ := pv.Sign(txBytesBeforeEmSig)
	tx.emSig = emSig
	txBytes, _ := tx.MarshalHash()
	hash := crypto.Hash(txBytes)
	tx.hash = hash

	tx.Mined(mv, []Validation{v2})

	chainDB := &mockChainDB{}

	assert.Nil(t, WriteTransaction(chainDB, tx, 1))
	assert.Len(t, chainDB.ids, 1)
	assert.Equal(t, crypto.Hash(txBytes), chainDB.ids[0].hash)
}

type mockChainDB struct {
	kos       []Transaction
	keychains []Keychain
	ids       []ID
}

func (db *mockChainDB) WriteKeychain(kc Keychain) error {
	db.keychains = append(db.keychains, kc)
	return nil
}

func (db *mockChainDB) WriteID(id ID) error {
	db.ids = append(db.ids, id)
	return nil
}

func (db *mockChainDB) WriteKO(tx Transaction) error {
	db.kos = append(db.kos, tx)
	return nil
}

func (db mockChainDB) FullKeychain(txAddr crypto.VersionnedHash) (*Keychain, error) {
	sort.Slice(db.keychains, func(i, j int) bool {
		return db.keychains[i].Timestamp().Unix() > db.keychains[j].Timestamp().Unix()
	})

	if len(db.keychains) > 0 {
		return &db.keychains[0], nil
	}
	return nil, nil
}

func (db mockChainDB) LastKeychain(txAddr crypto.VersionnedHash) (*Keychain, error) {
	sort.Slice(db.keychains, func(i, j int) bool {
		return db.keychains[i].Timestamp().Unix() > db.keychains[j].Timestamp().Unix()
	})

	if len(db.keychains) > 0 {
		return &db.keychains[0], nil
	}
	return nil, nil
}

func (db mockChainDB) KeychainByHash(hash crypto.VersionnedHash) (*Keychain, error) {
	for _, tx := range db.keychains {
		if bytes.Equal(tx.hash, hash) {
			return &tx, nil
		}
	}
	return nil, nil
}

func (db mockChainDB) IDByHash(hash crypto.VersionnedHash) (*ID, error) {
	for _, tx := range db.ids {
		if bytes.Equal(tx.hash, hash) {
			return &tx, nil
		}
	}
	return nil, nil
}

func (db mockChainDB) ID(addr crypto.VersionnedHash) (*ID, error) {
	for _, tx := range db.ids {
		if bytes.Equal(tx.Address(), addr) {
			return &tx, nil
		}
	}
	return nil, nil
}

func (db mockChainDB) KOByHash(hash crypto.VersionnedHash) (*Transaction, error) {
	for _, tx := range db.kos {
		if bytes.Equal(tx.hash, hash) {
			return &tx, nil
		}
	}
	return nil, nil
}

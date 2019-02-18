package chain

import (
	"encoding/hex"
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

	pub, pv := crypto.GenerateKeys()

	sig, _ := crypto.Sign("hello", pv)

	prop, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("pvKey")), pub)
	data := map[string]string{
		"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
		"encrypted_wallet":          hex.EncodeToString([]byte("wallet")),
	}
	tx := Transaction{
		addr:      crypto.HashString("addr"),
		data:      data,
		txType:    KeychainTransactionType,
		timestamp: time.Now(),
		pubKey:    pub,
		prop:      prop,
	}
	txBytesBeforeSig, _ := tx.MarshalBeforeSignature()
	sig, _ = crypto.Sign(string(txBytesBeforeSig), pv)
	tx.emSig = sig
	tx.sig = sig
	txBytes, _ := tx.MarshalHash()
	hash := crypto.HashBytes(txBytes)
	tx.hash = hash
	vBytes, _ := json.Marshal(Validation{
		nodePubk:  pub,
		status:    ValidationOK,
		timestamp: time.Now(),
	})
	vSig, _ := crypto.Sign(string(vBytes), pv)
	v, _ := NewValidation(ValidationOK, time.Now(), pub, vSig)
	mv, _ := NewMasterValidation([]string{}, pub, v)

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

	pub, pv := crypto.GenerateKeys()

	sig, _ := crypto.Sign("hello", pv)

	prop, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("pvKey")), pub)

	data := map[string]string{
		"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
		"encrypted_address_by_id":   hex.EncodeToString([]byte("addr")),
		"encrypted_aes_key":         hex.EncodeToString([]byte("aesKey")),
	}

	tx := Transaction{
		addr:      crypto.HashString("addr"),
		data:      data,
		txType:    IDTransactionType,
		timestamp: time.Now(),
		pubKey:    pub,
		prop:      prop,
	}
	txBytesBeforeSig, _ := tx.MarshalBeforeSignature()
	sig, _ = crypto.Sign(string(txBytesBeforeSig), pv)
	tx.emSig = sig
	tx.sig = sig
	txBytes, _ := tx.MarshalHash()
	hash := crypto.HashBytes(txBytes)
	tx.hash = hash
	vBytes, _ := json.Marshal(Validation{
		nodePubk:  pub,
		status:    ValidationOK,
		timestamp: time.Now(),
	})
	vSig, _ := crypto.Sign(string(vBytes), pv)
	v, _ := NewValidation(ValidationOK, time.Now(), pub, vSig)
	mv, _ := NewMasterValidation([]string{}, pub, v)

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

	_, err := getTransactionByHash(chainDB, crypto.HashString("hash"))
	assert.EqualError(t, err, "unknown transaction")
}

/*
Scenario: Get transaction status in progress
	Given a transaction stored in in progress
	When I dbant to get its status
	Then I get in progress
*/
func TestGetTransactionStatusInProgress(t *testing.T) {
	chainDB := &mockChainDB{
		inprogress: []Transaction{
			Transaction{
				hash: crypto.HashString("hash"),
			},
		},
	}

	status, err := GetTransactionStatus(chainDB, crypto.HashString("hash"))
	assert.Nil(t, err)
	assert.Equal(t, TransactionStatusInProgress, status)
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
				hash: crypto.HashString("hash"),
			},
		},
	}

	status, err := GetTransactionStatus(chainDB, crypto.HashString("hash"))
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

	status, err := GetTransactionStatus(chainDB, crypto.HashString("hash"))
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

	pub, pv := crypto.GenerateKeys()

	sig, _ := crypto.Sign("hello", pv)

	prop, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("pvKey")), pub)

	data := map[string]string{
		"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
		"encrypted_address_by_id":   hex.EncodeToString([]byte("addr")),
		"encrypted_aes_key":         hex.EncodeToString([]byte("aesKey")),
	}

	tx := Transaction{
		addr:      crypto.HashString("addr"),
		data:      data,
		txType:    IDTransactionType,
		timestamp: time.Now(),
		pubKey:    pub,
		prop:      prop,
	}
	txBytesBeforeSig, _ := tx.MarshalBeforeSignature()
	sig, _ = crypto.Sign(string(txBytesBeforeSig), pv)
	tx.emSig = sig
	tx.sig = sig
	txBytes, _ := tx.MarshalHash()
	hash := crypto.HashBytes(txBytes)
	tx.hash = hash
	vBytes, _ := json.Marshal(Validation{
		nodePubk:  pub,
		status:    ValidationOK,
		timestamp: time.Now(),
	})
	vSig, _ := crypto.Sign(string(vBytes), pv)
	v, _ := NewValidation(ValidationOK, time.Now(), pub, vSig)
	mv, _ := NewMasterValidation([]string{}, pub, v)

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

	pub, pv := crypto.GenerateKeys()

	sig, _ := crypto.Sign("hello", pv)

	prop, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("pvKey")), pub)

	data := map[string]string{
		"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
		"encrypted_wallet":          hex.EncodeToString([]byte("wallet")),
	}

	tx := Transaction{
		addr:      crypto.HashString("addr"),
		data:      data,
		txType:    KeychainTransactionType,
		timestamp: time.Now(),
		pubKey:    pub,
		prop:      prop,
	}
	txBytesBeforeSig, _ := tx.MarshalBeforeSignature()
	sig, _ = crypto.Sign(string(txBytesBeforeSig), pv)
	tx.emSig = sig
	tx.sig = sig
	txBytes, _ := tx.MarshalHash()
	hash := crypto.HashBytes(txBytes)
	tx.hash = hash

	vBytes, _ := json.Marshal(Validation{
		nodePubk:  pub,
		status:    ValidationOK,
		timestamp: time.Now(),
	})
	vSig, _ := crypto.Sign(string(vBytes), pv)
	v, _ := NewValidation(ValidationOK, time.Now(), pub, vSig)
	mv, _ := NewMasterValidation([]string{}, pub, v)

	tx.Mined(mv, []Validation{v})
	keychain1, _ := NewKeychain(tx)

	chainDB.keychains = append(chainDB.keychains, keychain1)

	time.Sleep(1 * time.Second)

	tx2 := Transaction{
		addr:      crypto.HashString("addr"),
		data:      data,
		txType:    KeychainTransactionType,
		timestamp: time.Now(),
		pubKey:    pub,
		prop:      prop,
	}
	txBytesBeforeSig2, _ := tx2.MarshalBeforeSignature()
	sig2, _ := crypto.Sign(string(txBytesBeforeSig2), pv)
	tx2.emSig = sig2
	tx2.sig = sig2
	txBytes2, _ := tx2.MarshalHash()
	hash2 := crypto.HashBytes(txBytes2)
	tx2.hash = hash2

	tx2.Mined(mv, []Validation{v})
	keychain2, err := NewKeychain(tx2)

	chainDB.keychains = append(chainDB.keychains, keychain2)
	assert.Len(t, chainDB.keychains, 2)

	lastTx, err := LastTransaction(chainDB, crypto.HashString("addr"), KeychainTransactionType)
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

	pub, pv := crypto.GenerateKeys()

	sig, _ := crypto.Sign("hello", pv)

	prop, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("pvKey")), pub)

	data := map[string]string{
		"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
		"encrypted_address_by_id":   hex.EncodeToString([]byte("addr")),
		"encrypted_aes_key":         hex.EncodeToString([]byte("aesKey")),
	}

	tx := Transaction{
		addr:      crypto.HashString("addr"),
		data:      data,
		txType:    IDTransactionType,
		timestamp: time.Now(),
		pubKey:    pub,
		prop:      prop,
	}
	txBytesBeforeSig, _ := tx.MarshalBeforeSignature()
	sig, _ = crypto.Sign(string(txBytesBeforeSig), pv)
	tx.emSig = sig
	tx.sig = sig
	txBytes, _ := tx.MarshalHash()
	hash := crypto.HashBytes(txBytes)
	tx.hash = hash

	vBytes, _ := json.Marshal(Validation{
		nodePubk:  pub,
		status:    ValidationOK,
		timestamp: time.Now(),
	})
	vSig, _ := crypto.Sign(string(vBytes), pv)
	v, _ := NewValidation(ValidationOK, time.Now(), pub, vSig)
	mv, _ := NewMasterValidation([]string{}, pub, v)

	tx.Mined(mv, []Validation{v})
	id1, _ := NewID(tx)

	chainDB.ids = append(chainDB.ids, id1)

	time.Sleep(1 * time.Second)

	tx2 := Transaction{
		addr:      crypto.HashString("addr2"),
		data:      data,
		txType:    IDTransactionType,
		timestamp: time.Now(),
		pubKey:    pub,
		prop:      prop,
	}
	txBytesBeforeSig2, _ := tx.MarshalBeforeSignature()
	sig2, _ := crypto.Sign(string(txBytesBeforeSig2), pv)
	tx2.emSig = sig2
	tx2.sig = sig2
	txBytes2, _ := tx2.MarshalHash()
	hash2 := crypto.HashBytes(txBytes2)
	tx2.hash = hash2
	tx2.Mined(mv, []Validation{v})
	id2, err := NewID(tx2)

	chainDB.ids = append(chainDB.ids, id2)
	assert.Len(t, chainDB.ids, 2)

	lastTx, err := LastTransaction(chainDB, crypto.HashString("addr"), IDTransactionType)
	assert.Nil(t, err)
	assert.NotNil(t, lastTx)
	assert.Equal(t, IDTransactionType, lastTx.txType)
	assert.Equal(t, hash, lastTx.hash)

	lastTx, err = LastTransaction(chainDB, crypto.HashString("addr2"), IDTransactionType)
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

	pub, pv := crypto.GenerateKeys()

	prop, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("pvKey")), pub)

	data := map[string]string{
		"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
		"encrypted_wallet":          hex.EncodeToString([]byte("wallet")),
	}

	tx := Transaction{
		addr:      crypto.HashString("addr"),
		data:      data,
		txType:    KeychainTransactionType,
		timestamp: time.Now(),
		pubKey:    pub,
		prop:      prop,
	}
	txBytesBeforeSig, _ := tx.MarshalBeforeSignature()
	sig, _ := crypto.Sign(string(txBytesBeforeSig), pv)
	tx.emSig = sig
	tx.sig = sig
	txBytes, _ := tx.MarshalHash()
	hash := crypto.HashBytes(txBytes)
	tx.hash = hash

	vBytes, _ := json.Marshal(Validation{
		nodePubk:  pub,
		status:    ValidationOK,
		timestamp: time.Now(),
	})
	vSig, _ := crypto.Sign(string(vBytes), pv)
	v, _ := NewValidation(ValidationOK, time.Now(), pub, vSig)
	mv, _ := NewMasterValidation([]string{}, pub, v)

	tx.Mined(mv, []Validation{v})
	keychain1, _ := NewKeychain(tx)

	chainDB.keychains = append(chainDB.keychains, keychain1)

	time.Sleep(1 * time.Second)

	tx2 := Transaction{
		addr:      crypto.HashString("addr"),
		data:      data,
		txType:    KeychainTransactionType,
		timestamp: time.Now(),
		pubKey:    pub,
		prop:      prop,
	}
	txBytesBeforeSig2, _ := tx2.MarshalBeforeSignature()
	sig2, _ := crypto.Sign(string(txBytesBeforeSig2), pv)
	tx2.emSig = sig2
	tx2.sig = sig2
	txBytes2, _ := tx2.MarshalHash()
	hash2 := crypto.HashBytes(txBytes2)
	tx2.hash = hash2

	tx2.Mined(mv, []Validation{v})
	keychain2, _ := NewKeychain(tx2)
	assert.Nil(t, keychain2.Chain(&tx))

	chainDB.keychains = append(chainDB.keychains, keychain2)
	assert.Len(t, chainDB.keychains, 2)

	chain, err := getFullChain(chainDB, crypto.HashString("addr"), KeychainTransactionType)
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
	pub, pv := crypto.GenerateKeys()

	prop, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("pvKey")), pub)

	data := map[string]string{
		"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
		"encrypted_wallet":          hex.EncodeToString([]byte("wallet")),
	}

	tx := Transaction{
		addr:      crypto.HashString("addr"),
		data:      data,
		txType:    KeychainTransactionType,
		timestamp: time.Now(),
		pubKey:    pub,
		prop:      prop,
	}
	txBytesBeforeSig, _ := tx.MarshalBeforeSignature()
	sig, _ := crypto.Sign(string(txBytesBeforeSig), pv)
	tx.sig = sig
	txBytesBeforeEmSig, _ := tx.MarshalBeforeEmitterSignature()
	emSig, _ := crypto.Sign(string(txBytesBeforeEmSig), pv)
	tx.emSig = emSig
	txBytes, _ := tx.MarshalHash()
	hash := crypto.HashBytes(txBytes)
	tx.hash = hash

	vBytes, _ := json.Marshal(Validation{
		nodePubk:  pub,
		status:    ValidationOK,
		timestamp: time.Now(),
	})
	vSig, _ := crypto.Sign(string(vBytes), pv)
	v, _ := NewValidation(ValidationOK, time.Now(), pub, vSig)
	mv, _ := NewMasterValidation([]string{}, pub, v)
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
	pub, pv := crypto.GenerateKeys()

	prop, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("pvKey")), pub)

	data := map[string]string{
		"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
		"encrypted_wallet":          hex.EncodeToString([]byte("wallet")),
	}

	tx := Transaction{
		addr:      crypto.HashString("addr"),
		data:      data,
		txType:    KeychainTransactionType,
		timestamp: time.Now(),
		pubKey:    pub,
		prop:      prop,
	}
	txBytesBeforeSig, _ := tx.MarshalBeforeSignature()
	sig, _ := crypto.Sign(string(txBytesBeforeSig), pv)
	tx.sig = sig
	txBytesBeforeEmSig, _ := tx.MarshalBeforeEmitterSignature()
	emSig, _ := crypto.Sign(string(txBytesBeforeEmSig), pv)
	tx.emSig = emSig
	txBytes, _ := tx.MarshalHash()
	hash := crypto.HashBytes(txBytes)
	tx.hash = hash

	vBytes, _ := json.Marshal(Validation{
		nodePubk:  pub,
		status:    ValidationOK,
		timestamp: time.Now(),
	})
	vSig, _ := crypto.Sign(string(vBytes), pv)
	v, _ := NewValidation(ValidationOK, time.Now(), pub, vSig)
	mv, _ := NewMasterValidation([]string{}, pub, v)
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

	pub, pv := crypto.GenerateKeys()

	vBytes, _ := json.Marshal(Validation{
		nodePubk:  pub,
		status:    ValidationKO,
		timestamp: time.Now(),
	})
	vSig, _ := crypto.Sign(string(vBytes), pv)
	v, _ := NewValidation(ValidationKO, time.Now(), pub, vSig)

	vBytes2, _ := json.Marshal(Validation{
		nodePubk:  pub,
		status:    ValidationOK,
		timestamp: time.Now(),
	})
	vSig2, _ := crypto.Sign(string(vBytes2), pv)
	v2, _ := NewValidation(ValidationOK, time.Now(), pub, vSig2)
	mv, _ := NewMasterValidation([]string{}, pub, v)

	prop, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("pvKey")), pub)

	data := map[string]string{
		"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
		"encrypted_wallet":          hex.EncodeToString([]byte("wallet")),
	}

	tx := Transaction{
		addr:      crypto.HashString("addr"),
		data:      data,
		txType:    KeychainTransactionType,
		timestamp: time.Now(),
		pubKey:    pub,
		prop:      prop,
	}
	txBytesBeforeSig, _ := tx.MarshalBeforeSignature()
	sig, _ := crypto.Sign(string(txBytesBeforeSig), pv)
	tx.sig = sig
	txBytesBeforeEmSig, _ := tx.MarshalBeforeEmitterSignature()
	emSig, _ := crypto.Sign(string(txBytesBeforeEmSig), pv)
	tx.emSig = emSig
	txBytes, _ := tx.MarshalHash()
	hash := crypto.HashBytes(txBytes)
	tx.hash = hash

	tx.Mined(mv, []Validation{v2})

	chainDB := &mockChainDB{}

	assert.Nil(t, WriteTransaction(chainDB, &mockLocker{}, tx, 1))
	assert.Len(t, chainDB.kos, 1)
	assert.Equal(t, crypto.HashBytes(txBytes), chainDB.kos[0].hash)
}

/*
Scenario: Store a Keychain transaction
	Given a keychain transaction
	When I dbant to store it
	Then the transaction is stored on the keychain db
*/
func TestStoreKeychainTransaction(t *testing.T) {

	pub, pv := crypto.GenerateKeys()

	vBytes, _ := json.Marshal(Validation{
		nodePubk:  pub,
		status:    ValidationOK,
		timestamp: time.Now(),
	})
	vSig, _ := crypto.Sign(string(vBytes), pv)
	v, _ := NewValidation(ValidationOK, time.Now(), pub, vSig)

	vBytes2, _ := json.Marshal(Validation{
		nodePubk:  pub,
		status:    ValidationOK,
		timestamp: time.Now(),
	})
	vSig2, _ := crypto.Sign(string(vBytes2), pv)
	v2, _ := NewValidation(ValidationOK, time.Now(), pub, vSig2)
	mv, _ := NewMasterValidation([]string{}, pub, v)

	prop, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("pvKey")), pub)

	data := map[string]string{
		"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
		"encrypted_wallet":          hex.EncodeToString([]byte("wallet")),
	}

	tx := Transaction{
		addr:      crypto.HashString("addr"),
		data:      data,
		txType:    KeychainTransactionType,
		timestamp: time.Now(),
		pubKey:    pub,
		prop:      prop,
	}
	txBytesBeforeSig, _ := tx.MarshalBeforeSignature()
	sig, _ := crypto.Sign(string(txBytesBeforeSig), pv)
	tx.sig = sig
	txBytesBeforeEmSig, _ := tx.MarshalBeforeEmitterSignature()
	emSig, _ := crypto.Sign(string(txBytesBeforeEmSig), pv)
	tx.emSig = emSig
	txBytes, _ := tx.MarshalHash()
	hash := crypto.HashBytes(txBytes)
	tx.hash = hash

	tx.Mined(mv, []Validation{v2})

	chainDB := &mockChainDB{}

	assert.Nil(t, WriteTransaction(chainDB, &mockLocker{}, tx, 1))
	assert.Len(t, chainDB.keychains, 1)
	assert.Equal(t, crypto.HashBytes(txBytes), chainDB.keychains[0].hash)
}

/*
Scenario: Store a ID transaction
	Given a ID transaction
	When I dbant to store it
	Then the transaction is stored on the ID db
*/
func TestStoreIDTransaction(t *testing.T) {

	pub, pv := crypto.GenerateKeys()

	vBytes, _ := json.Marshal(Validation{
		nodePubk:  pub,
		status:    ValidationOK,
		timestamp: time.Now(),
	})
	vSig, _ := crypto.Sign(string(vBytes), pv)
	v, _ := NewValidation(ValidationOK, time.Now(), pub, vSig)

	vBytes2, _ := json.Marshal(Validation{
		nodePubk:  pub,
		status:    ValidationOK,
		timestamp: time.Now(),
	})
	vSig2, _ := crypto.Sign(string(vBytes2), pv)
	v2, _ := NewValidation(ValidationOK, time.Now(), pub, vSig2)
	mv, _ := NewMasterValidation([]string{}, pub, v)

	prop, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("pvKey")), pub)

	data := map[string]string{
		"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
		"encrypted_address_by_id":   hex.EncodeToString([]byte("addr")),
		"encrypted_aes_key":         hex.EncodeToString([]byte("aesKey")),
	}

	tx := Transaction{
		addr:      crypto.HashString("addr"),
		data:      data,
		txType:    IDTransactionType,
		timestamp: time.Now(),
		pubKey:    pub,
		prop:      prop,
	}
	txBytesBeforeSig, _ := tx.MarshalBeforeSignature()
	sig, _ := crypto.Sign(string(txBytesBeforeSig), pv)
	tx.sig = sig
	txBytesBeforeEmSig, _ := tx.MarshalBeforeEmitterSignature()
	emSig, _ := crypto.Sign(string(txBytesBeforeEmSig), pv)
	tx.emSig = emSig
	txBytes, _ := tx.MarshalHash()
	hash := crypto.HashBytes(txBytes)
	tx.hash = hash

	tx.Mined(mv, []Validation{v2})

	chainDB := &mockChainDB{}

	assert.Nil(t, WriteTransaction(chainDB, &mockLocker{}, tx, 1))
	assert.Len(t, chainDB.ids, 1)
	assert.Equal(t, crypto.HashBytes(txBytes), chainDB.ids[0].hash)
}

type mockChainDB struct {
	inprogress []Transaction
	kos        []Transaction
	keychains  []Keychain
	ids        []ID
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

func (db *mockChainDB) WriteInProgress(tx Transaction) error {
	db.inprogress = append(db.inprogress, tx)
	return nil
}

func (db mockChainDB) InProgressByHash(hash string) (*Transaction, error) {
	for _, tx := range db.inprogress {
		if tx.hash == hash {
			return &tx, nil
		}
	}
	return nil, nil
}

func (db mockChainDB) FullKeychain(txAddr string) (*Keychain, error) {
	sort.Slice(db.keychains, func(i, j int) bool {
		return db.keychains[i].Timestamp().Unix() > db.keychains[j].Timestamp().Unix()
	})

	if len(db.keychains) > 0 {
		return &db.keychains[0], nil
	}
	return nil, nil
}

func (db mockChainDB) LastKeychain(txAddr string) (*Keychain, error) {
	sort.Slice(db.keychains, func(i, j int) bool {
		return db.keychains[i].Timestamp().Unix() > db.keychains[j].Timestamp().Unix()
	})

	if len(db.keychains) > 0 {
		return &db.keychains[0], nil
	}
	return nil, nil
}

func (db mockChainDB) KeychainByHash(hash string) (*Keychain, error) {
	for _, tx := range db.keychains {
		if tx.hash == hash {
			return &tx, nil
		}
	}
	return nil, nil
}

func (db mockChainDB) IDByHash(hash string) (*ID, error) {
	for _, tx := range db.ids {
		if tx.hash == hash {
			return &tx, nil
		}
	}
	return nil, nil
}

func (db mockChainDB) ID(addr string) (*ID, error) {
	for _, tx := range db.ids {
		if tx.Address() == addr {
			return &tx, nil
		}
	}
	return nil, nil
}

func (db mockChainDB) KOByHash(hash string) (*Transaction, error) {
	for _, tx := range db.kos {
		if tx.hash == hash {
			return &tx, nil
		}
	}
	return nil, nil
}

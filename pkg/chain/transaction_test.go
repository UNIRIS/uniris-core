package chain

import (
	"crypto/rand"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/uniris/uniris-core/pkg/crypto"
	"github.com/uniris/uniris-core/pkg/shared"
)

/*
Scenario: Create a new transaction
	Given transaction data (addr, hash, public key, signature, emSig, prop, timestamp)
	When I want to create the transaction
	Then I get it
*/
func TestNewTransactionTransaction(t *testing.T) {
	pv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

	prop, _ := shared.NewEmitterKeyPair([]byte("encPvKey"), pub)

	addr := crypto.Hash([]byte("addr"))
	data := map[string][]byte{
		"encrypted_address_by_node": []byte("addr"),
		"encrypted_address_by_id":   []byte("addr"),
		"encrypted_aes_key":         []byte("aesKey"),
	}
	hash := crypto.Hash([]byte("hash"))

	sig, _ := pv.Sign([]byte("data"))

	tx, err := NewTransaction(addr, KeychainTransactionType, data, time.Now(), pub, prop, sig, sig, hash)
	assert.Nil(t, err)
	assert.Equal(t, addr, tx.Address())
	assert.Equal(t, data, tx.Data())
	assert.Equal(t, KeychainTransactionType, tx.TransactionType())
	assert.Equal(t, sig, tx.Signature())
	assert.Equal(t, pub, tx.EmitterSharedKeyProposal().PublicKey())
}

/*
Scenario: Create a new transaction with an invalid addr
	Given a invalid addr hash, empty or not in he
	When I want to create the transaction
	Then I get an error
*/
func TestNewTransactionWithInvalidAddress(t *testing.T) {
	_, err := NewTransaction([]byte(""), KeychainTransactionType, map[string][]byte{}, time.Now(), nil, shared.EmitterKeyPair{}, []byte(""), []byte(""), []byte(""))
	assert.EqualError(t, err, "transaction address is missing")

	_, err = NewTransaction([]byte("abc"), KeychainTransactionType, map[string][]byte{}, time.Now(), nil, shared.EmitterKeyPair{}, []byte(""), []byte(""), []byte(""))
	assert.EqualError(t, err, "transaction address is not a valid hash")
}

/*
Scenario: Create a new transaction without public key
	Given a transaction without public key
	When I want to create the transaction
	Then I get an error
*/
func TestNewTransactionWithoutPublicKey(t *testing.T) {

	_, err := NewTransaction(crypto.Hash([]byte("addr")), KeychainTransactionType, map[string][]byte{}, time.Now(), nil, shared.EmitterKeyPair{}, []byte("fake sig"), []byte("fake sig"), []byte(""))
	assert.EqualError(t, err, "transaction public key is missing")
}

/*
Scenario: Create a new transaction without signature
	Given a transaction without signature
	When I want to create the transaction
	Then I get an error
*/
func TestNewTransactionWithoutSignature(t *testing.T) {
	_, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
	_, err := NewTransaction(crypto.Hash([]byte("addr")), KeychainTransactionType, map[string][]byte{}, time.Now(), pub, shared.EmitterKeyPair{}, nil, nil, []byte(""))
	assert.EqualError(t, err, "transaction signature is missing")
}

/*
Scenario: Create a new transaction without emitter signature
	Given a transaction without emitt ersignature
	When I want to create the transaction
	Then I get an error
*/
func TestNewTransactionWithoutEmitterSignature(t *testing.T) {
	_, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
	_, err := NewTransaction(crypto.Hash([]byte("addr")), KeychainTransactionType, map[string][]byte{}, time.Now(), pub, shared.EmitterKeyPair{}, []byte("fake sig"), nil, []byte(""))
	assert.EqualError(t, err, "transaction emitter signature is missing")
}

/*
Scenario: Create a new transaction with an invalid transaction hash
	Given a invalid transaction hash, empty or not in he
	When I want to create the transaction
	Then I get an error
*/
func TestNewTransactionWithInvalidHash(t *testing.T) {

	_, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

	_, err := NewTransaction(crypto.Hash([]byte("addr")), KeychainTransactionType, map[string][]byte{}, time.Now(), pub, shared.EmitterKeyPair{}, []byte("fake sig"), []byte("fake sig"), []byte(""))
	assert.EqualError(t, err, "transaction hash is missing")

	_, err = NewTransaction(crypto.Hash([]byte("addr")), KeychainTransactionType, map[string][]byte{}, time.Now(), pub, shared.EmitterKeyPair{}, []byte("fake sig"), []byte("fake sig"), []byte("abc"))
	assert.EqualError(t, err, "transaction hash is not a valid hash")
}

/*
Scenario: Create a new transaction with an invalid transaction data
	Given a empty data
	When I want to create the transaction
	Then I get an error
*/
func TestNewTransactionWithInvalidData(t *testing.T) {

	_, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
	_, err := NewTransaction(crypto.Hash([]byte("addr")), KeychainTransactionType, map[string][]byte{}, time.Now(), pub, shared.EmitterKeyPair{}, []byte("fake sig"), []byte("fake sig"), crypto.Hash([]byte("addr")))
	assert.EqualError(t, err, "transaction data is missing")
}

/*
Scenario: Create a new transaction with an invalid transaction timestamp (more than the current timestamp)
	Given a transaction timestamp (now + 2 seconds)
	When I want to create the transaction
	Then I get an error
*/
func TestNewTransactionWithInvalidTimestamp(t *testing.T) {
	_, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
	_, err := NewTransaction(crypto.Hash([]byte("addr")), KeychainTransactionType, map[string][]byte{
		"encrypted_address_by_node": []byte("addr"),
		"encrypted_wallet":          []byte("wallet"),
	}, time.Now().Add(2*time.Second), pub, shared.EmitterKeyPair{}, []byte("fake sig"), []byte("fake sig"), crypto.Hash([]byte("addr")))
	assert.EqualError(t, err, "transaction timestamp must be greater lower than now")
}

/*
Scenario: Create a new transaction with invalid transaction type
	Given an invalid type
	When I want to create the transaction
	Then I get an error
*/
func TestNewTransactionWithInvalidType(t *testing.T) {
	pv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

	sig, _ := pv.Sign([]byte("sig"))

	_, err := NewTransaction(crypto.Hash([]byte("addr")), 10, map[string][]byte{
		"encrypted_address_by_node": []byte("addr"),
		"encrypted_wallet":          []byte("wallet"),
	}, time.Now(), pub, shared.EmitterKeyPair{}, sig, sig, crypto.Hash([]byte("hello")))
	assert.EqualError(t, err, "transaction type is not allowed")
}

/*
Scenario: Create a new transaction without proposal
	Given a transaction without proposal
	When I want to create the transaction
	Then I get an error
*/
func TestNewTransactionWithoutProposal(t *testing.T) {
	pv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

	sig, _ := pv.Sign([]byte("sig"))

	_, err := NewTransaction(crypto.Hash([]byte("addr")), KeychainTransactionType, map[string][]byte{
		"encrypted_address_by_node": []byte("addr"),
		"encrypted_wallet":          []byte("wallet"),
	}, time.Now(), pub, shared.EmitterKeyPair{}, sig, sig, crypto.Hash([]byte("hello")))
	assert.EqualError(t, err, "transaction proposal private key is missing")

	prop, _ := shared.NewEmitterKeyPair([]byte("encPv"), nil)
	_, err = NewTransaction(crypto.Hash([]byte("addr")), KeychainTransactionType, map[string][]byte{
		"encrypted_address_by_node": []byte("addr"),
		"encrypted_wallet":          []byte("wallet"),
	}, time.Now(), pub, prop, sig, sig, crypto.Hash([]byte("hello")))
	assert.EqualError(t, err, "transaction proposal public key is missing")
}

/*
Scenario: Check the transaction integrity
	Given a transaction with a valid hash and valid signature
	When I want to check its intergrity, its check the transaction hash and the signature
	Then I get not error
*/
func TestCheckTransactionIntegrity(t *testing.T) {
	pv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

	prop, _ := shared.NewEmitterKeyPair([]byte("encPvKey"), pub)

	txRaw := Transaction{
		addr: crypto.Hash([]byte("addr")),
		data: map[string][]byte{
			"encrypted_address_by_node": []byte("addr"),
			"encrypted_wallet":          []byte("wallet"),
		},
		txType:    KeychainTransactionType,
		timestamp: time.Now(),
		pubKey:    pub,
		prop:      prop,
	}
	txBytesBeforeSig, _ := txRaw.MarshalBeforeSignature()
	sig, _ := pv.Sign(txBytesBeforeSig)
	txRaw.emSig = sig
	txRaw.sig = sig
	txBytes, _ := txRaw.MarshalHash()

	hash := crypto.Hash(txBytes)

	tx, _ := NewTransaction(txRaw.addr, KeychainTransactionType, txRaw.data, txRaw.timestamp, txRaw.pubKey, txRaw.prop, txRaw.sig, txRaw.emSig, hash)
	assert.Nil(t, tx.checkTransactionIntegrity())
}

/*
Scenario: Check the transaction integrity with invalid hash
	Given a transaction with a invalid hash
	When I want to check its intergrity, its check the transaction hash and the signature
	Then I get an error
*/
func TestCheckTransactionIntegrityWithInvalidHash(t *testing.T) {
	pv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

	prop, _ := shared.NewEmitterKeyPair([]byte("encPvKey"), pub)

	raw, _ := json.Marshal(Transaction{
		addr: crypto.Hash([]byte("addr")),
		data: map[string][]byte{
			"encrypted_address_by_node": []byte("addr"),
			"encrypted_wallet":          []byte("wallet"),
		},
		txType:    KeychainTransactionType,
		timestamp: time.Now(),
		pubKey:    pub,
		prop:      prop,
	})
	sig, _ := pv.Sign(raw)
	hash := crypto.Hash([]byte("abc"))

	tx, _ := NewTransaction(crypto.Hash([]byte("addr")), KeychainTransactionType, map[string][]byte{
		"encrypted_address_by_node": []byte("addr"),
		"encrypted_wallet":          []byte("wallet"),
	}, time.Now(), pub, prop, sig, sig, hash)
	assert.EqualError(t, tx.checkTransactionIntegrity(), "transaction integrity violated")
}

/*
Scenario: Check the transaction integrity with invalid signature
	Given a transaction with a valid hash and invalid signature
	When I want to check its intergrity, its check the transaction hash and the signature
	Then I get not error
*/
func TestCheckTransactionIntegrityWithInvalidSignature(t *testing.T) {
	pv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

	prop, _ := shared.NewEmitterKeyPair([]byte("encPvKey"), pub)

	txRaw := Transaction{
		addr: crypto.Hash([]byte("addr")),
		data: map[string][]byte{
			"encrypted_address_by_node": []byte("addr"),
			"encrypted_wallet":          []byte("wallet"),
		},
		txType:    KeychainTransactionType,
		timestamp: time.Now(),
		pubKey:    pub,
		prop:      prop,
	}
	sig, _ := pv.Sign([]byte("fakesig"))
	txRaw.emSig = sig
	txRaw.sig = sig

	txBytes, _ := txRaw.MarshalHash()

	hash := crypto.Hash(txBytes)

	tx, _ := NewTransaction(crypto.Hash([]byte("addr")), KeychainTransactionType, map[string][]byte{
		"encrypted_address_by_node": []byte("addr"),
		"encrypted_wallet":          []byte("wallet"),
	}, time.Now(), pub, prop, sig, sig, hash)
	assert.EqualError(t, tx.checkTransactionIntegrity(), "transaction signature invalid")
}

/*
Scenario: Add mining information to a transaction
	Given a transaction
	When I want to add master validation and confirmation validations
	Then I can retrieve in inside the transaction
*/
func TestMined(t *testing.T) {
	pv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

	prop, _ := shared.NewEmitterKeyPair([]byte("encPvKey"), pub)

	raw, _ := json.Marshal(Transaction{
		addr: crypto.Hash([]byte("addr")),
		data: map[string][]byte{
			"encrypted_address_by_node": []byte("addr"),
			"encrypted_wallet":          []byte("wallet"),
		},
		txType:    KeychainTransactionType,
		timestamp: time.Now(),
		pubKey:    pub,
		prop:      prop,
	})
	sig, _ := pv.Sign(raw)
	hash := crypto.Hash(raw)
	tx, _ := NewTransaction(crypto.Hash([]byte("addr")), KeychainTransactionType, map[string][]byte{
		"encrypted_address_by_node": []byte("addr"),
		"encrypted_wallet":          []byte("wallet"),
	}, time.Now(), pub, prop, sig, sig, hash)

	b, _ := json.Marshal(Validation{
		nodePubk:  pub,
		status:    ValidationOK,
		timestamp: time.Now(),
	})
	sig, _ = pv.Sign(b)
	v, _ := NewValidation(ValidationOK, time.Now(), pub, sig)

	wHeaders := []NodeHeader{NewNodeHeader(pub, false, false, 0, true)}
	vHeaders := []NodeHeader{NewNodeHeader(pub, false, false, 0, true)}
	sHeaders := []NodeHeader{NewNodeHeader(pub, false, false, 0, true)}
	masterValid, _ := NewMasterValidation([]crypto.PublicKey{}, pub, v, wHeaders, vHeaders, sHeaders)
	assert.Nil(t, tx.Mined(masterValid, []Validation{v}))

	assert.Equal(t, sig, tx.MasterValidation().Validation().Signature())
	assert.Len(t, tx.ConfirmationsValidations(), 1)
}

/*
Scenario: Add mining information to a transaction without confirmations
	Given a transaction
	When I want to add master validation
	Then I get an error
*/
func TestMinedWithoutConfirmations(t *testing.T) {
	pv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

	prop, _ := shared.NewEmitterKeyPair([]byte("encPvKey"), pub)

	raw, _ := json.Marshal(Transaction{
		addr: crypto.Hash([]byte("addr")),
		data: map[string][]byte{
			"encrypted_address_by_node": []byte("addr"),
			"encrypted_wallet":          []byte("wallet"),
		},
		txType:    KeychainTransactionType,
		timestamp: time.Now(),
		pubKey:    pub,
		prop:      prop,
	})
	sig, _ := pv.Sign(raw)
	hash := crypto.Hash(raw)
	tx, _ := NewTransaction(crypto.Hash([]byte("addr")), KeychainTransactionType, map[string][]byte{
		"encrypted_address_by_node": []byte("addr"),
		"encrypted_wallet":          []byte("wallet"),
	}, time.Now(), pub, prop, sig, sig, hash)

	b, _ := json.Marshal(Validation{
		nodePubk:  pub,
		status:    ValidationOK,
		timestamp: time.Now(),
	})
	sig, _ = pv.Sign(b)
	v, _ := NewValidation(ValidationOK, time.Now(), pub, sig)
	wHeaders := []NodeHeader{NewNodeHeader(pub, false, false, 0, true)}
	vHeaders := []NodeHeader{NewNodeHeader(pub, false, false, 0, true)}
	sHeaders := []NodeHeader{NewNodeHeader(pub, false, false, 0, true)}
	masterValid, _ := NewMasterValidation([]crypto.PublicKey{}, pub, v, wHeaders, vHeaders, sHeaders)
	assert.EqualError(t, tx.Mined(masterValid, []Validation{}), "confirmation validations of the transaction are missing")
}

/*
Scenario: Check the integrity of a transaction chain
	Given 3 transactions chained
	When I want to check their the chain integrity
	Then I get not error
*/
func TestCheckChainIntegrity(t *testing.T) {
	pv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

	prop, _ := shared.NewEmitterKeyPair([]byte("encPvKey"), pub)

	tx1 := Transaction{
		addr: crypto.Hash([]byte("addr")),
		data: map[string][]byte{
			"encrypted_address_by_node": []byte("addr"),
			"encrypted_wallet":          []byte("wallet"),
		},
		txType:    KeychainTransactionType,
		timestamp: time.Now(),
		pubKey:    pub,
		prop:      prop,
	}
	txBytesBeforeSig, _ := tx1.MarshalBeforeSignature()
	sig, _ := pv.Sign(txBytesBeforeSig)
	tx1.emSig = sig
	tx1.sig = sig
	txBytes1, _ := tx1.MarshalHash()
	txHash1 := crypto.Hash(txBytes1)
	tx1.hash = txHash1

	tx2 := Transaction{
		addr: crypto.Hash([]byte("addr")),
		data: map[string][]byte{
			"encrypted_address_by_node": []byte("addr"),
			"encrypted_wallet":          []byte("wallet"),
		},
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
	txHash2 := crypto.Hash(txBytes2)
	tx2.hash = txHash2
	tx2.prevTx = &tx1

	tx3 := Transaction{
		addr: crypto.Hash([]byte("addr")),
		data: map[string][]byte{
			"encrypted_address_by_node": []byte("addr"),
			"encrypted_wallet":          []byte("wallet"),
		},
		txType:    KeychainTransactionType,
		timestamp: time.Now().Add(3 * time.Second),
		pubKey:    pub,
		prop:      prop,
	}
	txBytesBeforeSig3, _ := tx3.MarshalBeforeSignature()
	sig3, _ := pv.Sign(txBytesBeforeSig3)
	tx3.emSig = sig3
	tx3.sig = sig3
	txBytes3, _ := tx3.MarshalHash()
	txHash3 := crypto.Hash(txBytes3)
	tx3.hash = txHash3
	tx3.prevTx = &tx2

	assert.Nil(t, tx3.CheckChainTransactionIntegrity())
}

/*
Scenario: Check chain integrity with invalid timestamp
	Given a two transaction with the same timestamp
	When I want to check their chain
	Then I get an error
*/
func TestCheckChainIntegrityWithInvalidTime(t *testing.T) {
	pv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

	prop, _ := shared.NewEmitterKeyPair([]byte("encPvKey"), pub)

	tx0 := Transaction{
		addr: crypto.Hash([]byte("addr")),
		data: map[string][]byte{
			"encrypted_address_by_node": []byte("addr"),
			"encrypted_wallet":          []byte("wallet"),
		},
		txType:    KeychainTransactionType,
		timestamp: time.Now(),
		pubKey:    pub,
		prop:      prop,
	}

	b, _ := json.Marshal(tx0)
	hash := crypto.Hash(b)
	sig, _ := pv.Sign(b)
	tx0.sig = sig
	tx0.hash = hash

	tx1 := Transaction{
		addr: crypto.Hash([]byte("addr")),
		data: map[string][]byte{
			"encrypted_address_by_node": []byte("addr"),
			"encrypted_wallet":          []byte("wallet"),
		},
		txType:    KeychainTransactionType,
		timestamp: time.Now(),
		pubKey:    pub,
		prop:      prop,
		prevTx:    &tx0,
	}

	assert.EqualError(t, tx1.CheckChainTransactionIntegrity(), "previous chained transaction must be anterior to the current transaction")
}

/*
Scenario: Chain a transaction to another one
	Given two transaction
	When I want to chain them
	Then I get not error
*/
func TestChainTransaction(t *testing.T) {

	pv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

	prop, _ := shared.NewEmitterKeyPair([]byte("encPvKey"), pub)

	tx1 := Transaction{
		addr: crypto.Hash([]byte("addr")),
		data: map[string][]byte{
			"encrypted_address_by_node": []byte("addr"),
			"encrypted_wallet":          []byte("wallet"),
		},
		txType:    KeychainTransactionType,
		timestamp: time.Now(),
		pubKey:    pub,
		prop:      prop,
	}
	txBytesBeforeSig, _ := tx1.MarshalBeforeSignature()
	sig, _ := pv.Sign(txBytesBeforeSig)
	tx1.emSig = sig
	tx1.sig = sig
	txBytes1, _ := tx1.MarshalHash()
	txHash1 := crypto.Hash(txBytes1)
	tx1.hash = txHash1

	tx2 := Transaction{
		addr: crypto.Hash([]byte("addr")),
		data: map[string][]byte{
			"encrypted_address_by_node": []byte("addr"),
			"encrypted_wallet":          []byte("wallet"),
		},
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
	txHash2 := crypto.Hash(txBytes2)
	tx2.hash = txHash2

	assert.Nil(t, tx2.Chain(&tx1))
	assert.NotNil(t, tx2.PreviousTransaction())
	assert.Equal(t, tx1.hash, tx2.PreviousTransaction().TransactionHash())
}

/*
Scenario: Chain a transaction with same timestamp for the both
	Given two transaction with the same timestamp
	When I want to chain them
	I got an error
*/
func TestChainTransactionWithInvalidTimestamp(t *testing.T) {

	pv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

	prop, _ := shared.NewEmitterKeyPair([]byte("encPvKey"), pub)

	txTime1 := time.Now()
	raw1, _ := json.Marshal(Transaction{
		addr: crypto.Hash([]byte("addr")),
		data: map[string][]byte{
			"encrypted_address_by_node": []byte("addr"),
			"encrypted_wallet":          []byte("wallet"),
		},
		txType:    KeychainTransactionType,
		timestamp: txTime1,
		pubKey:    pub,
		prop:      prop,
	})
	sig1, _ := pv.Sign(raw1)
	hash1 := crypto.Hash(raw1)

	tx1, _ := NewTransaction(crypto.Hash([]byte("addr")), KeychainTransactionType, map[string][]byte{
		"encrypted_address_by_node": []byte("addr"),
		"encrypted_wallet":          []byte("wallet"),
	}, txTime1, pub, prop, sig1, sig1, hash1)

	raw2, _ := json.Marshal(Transaction{
		addr: crypto.Hash([]byte("addr")),
		data: map[string][]byte{
			"encrypted_address_by_node": []byte("addr"),
			"encrypted_wallet":          []byte("wallet"),
		},
		txType:    KeychainTransactionType,
		timestamp: txTime1,
		pubKey:    pub,
		prop:      prop,
	})
	sig2, _ := pv.Sign(raw2)
	hash2 := crypto.Hash(raw2)

	tx2, _ := NewTransaction(crypto.Hash([]byte("addr")), KeychainTransactionType, map[string][]byte{
		"encrypted_address_by_node": []byte("addr"),
		"encrypted_wallet":          []byte("wallet"),
	}, txTime1, pub, prop, sig2, sig2, hash2)

	assert.EqualError(t, tx2.Chain(&tx1), "previous chained transaction must be anterior to the current transaction")
}

/*
Scenario: Check master validation
	Given a transaction with a master validation
	When I want to check the master validation and the POW
	Then I get not error
*/
func TestCheckMasterValidation(t *testing.T) {

	pv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

	v := Validation{
		nodePubk:  pub,
		status:    ValidationOK,
		timestamp: time.Now(),
	}
	b, _ := json.Marshal(v)
	sig, _ := pv.Sign(b)
	v.nodeSig = sig

	prop, _ := shared.NewEmitterKeyPair([]byte("encPvKey"), pub)

	tx := Transaction{
		addr: crypto.Hash([]byte("addr")),
		data: map[string][]byte{
			"encrypted_address_by_node": []byte("addr"),
			"encrypted_wallet":          []byte("wallet"),
		},
		txType:    KeychainTransactionType,
		timestamp: time.Now(),
		pubKey:    pub,
		prop:      prop,
	}
	txBytesBeforeSig, _ := tx.MarshalBeforeSignature()
	sig, _ = pv.Sign(txBytesBeforeSig)
	tx.sig = sig
	txBytesBeforeEmSig, _ := tx.MarshalBeforeEmitterSignature()
	emSig, _ := pv.Sign(txBytesBeforeEmSig)
	tx.emSig = emSig
	txBytes, _ := tx.MarshalHash()
	txHash := crypto.Hash(txBytes)
	tx.hash = txHash

	wHeaders := []NodeHeader{NewNodeHeader(pub, false, false, 0, true)}
	vHeaders := []NodeHeader{NewNodeHeader(pub, false, false, 0, true)}
	sHeaders := []NodeHeader{NewNodeHeader(pub, false, false, 0, true)}
	masterValid, _ := NewMasterValidation([]crypto.PublicKey{}, pub, v, wHeaders, vHeaders, sHeaders)
	tx.masterV = masterValid

	assert.Nil(t, tx.CheckMasterValidation())
}

/*
Scenario: Check master validation with invalid POW
	Given a transaction with a master validation including bad POW
	When I want to check the master validation and the POW
	Then I get an error
*/
func TestCheckMasterValidationWithInvalidPOW(t *testing.T) {

	pv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

	v := Validation{
		nodePubk:  pub,
		status:    ValidationOK,
		timestamp: time.Now(),
	}
	b, _ := json.Marshal(v)
	sig, _ := pv.Sign(b)
	v.nodeSig = sig

	prop, _ := shared.NewEmitterKeyPair([]byte("encPvKey"), pub)

	raw, _ := json.Marshal(Transaction{
		addr: crypto.Hash([]byte("addr")),
		data: map[string][]byte{
			"encrypted_address_by_node": []byte("addr"),
			"encrypted_wallet":          []byte("wallet"),
		},
		txType:    KeychainTransactionType,
		timestamp: time.Now(),
		pubKey:    pub,
		prop:      prop,
	})
	sigTx, _ := pv.Sign(raw)
	hash := crypto.Hash(raw)

	wHeaders := []NodeHeader{NewNodeHeader(pub, false, false, 0, true)}
	vHeaders := []NodeHeader{NewNodeHeader(pub, false, false, 0, true)}
	sHeaders := []NodeHeader{NewNodeHeader(pub, false, false, 0, true)}
	masterValid, _ := NewMasterValidation([]crypto.PublicKey{}, pub, v, wHeaders, vHeaders, sHeaders)
	tx := Transaction{
		masterV: masterValid,
		addr:    crypto.Hash([]byte("addr")),
		data: map[string][]byte{
			"encrypted_address_by_node": []byte("addr"),
			"encrypted_wallet":          []byte("wallet"),
		},
		txType:    KeychainTransactionType,
		timestamp: time.Now(),
		pubKey:    pub,
		prop:      prop,
		emSig:     sigTx,
		sig:       sigTx,
		hash:      hash,
	}

	assert.EqualError(t, tx.CheckMasterValidation(), "invalid proof of work")
}

/*
Scenario: Create a new node validation
	Given a public key, a status, a timestamp and signature
	When I want to create a node validation
	Then I get the validation
*/
func TestNewTransactionValidation(t *testing.T) {

	pv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

	b, _ := json.Marshal(Validation{
		nodePubk:  pub,
		status:    ValidationOK,
		timestamp: time.Now(),
	})
	sig, _ := pv.Sign(b)

	v, err := NewValidation(ValidationOK, time.Now(), pub, sig)
	assert.Nil(t, err)
	assert.Equal(t, ValidationOK, v.Status())
	assert.Equal(t, time.Now().Unix(), v.Timestamp().Unix())
	assert.Equal(t, pub, v.PublicKey())
	assert.Equal(t, sig, v.Signature())
}

/*
Scenario: Create a new node validation with a timestamp later than now
	Given a public key, a status and a timestamp (now + 2 sec)
	When I want to create a node validation
	Then I get an error
*/
func TestNewTransactionValidationWithInvalidTimestamp(t *testing.T) {
	_, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

	_, err := NewValidation(ValidationOK, time.Now().Add(2*time.Second), pub, []byte("sig"))
	assert.EqualError(t, err, "validation timestamp must be anterior or equal to now")
}

/*
Scenario: Create a new node validation with invalid public key
	Given no public key or no hex or not valid public key
	When I want to create a node validation
	Then I get an error
*/
func TestNewTransactionValidationWithInvalidPublicKey(t *testing.T) {
	_, err := NewValidation(ValidationOK, time.Now(), nil, []byte("sig"))
	assert.EqualError(t, err, "validation public key is missing")
}

/*
Scenario: Create a new node validation with invalid signature
	Given no hex or not valid signature
	When I want to create a node validation
	Then I get an error
*/
func TestNewTransactionValidationWithInvalidSignature(t *testing.T) {

	pv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

	_, err := NewValidation(ValidationOK, time.Now(), pub, nil)
	assert.EqualError(t, err, "validation signature is missing")

	_, err = NewValidation(ValidationOK, time.Now(), pub, []byte("sig"))
	assert.EqualError(t, err, "validation signature is not valid")

	sig, _ := pv.Sign([]byte("hello"))
	_, err = NewValidation(ValidationOK, time.Now(), pub, sig)
	assert.EqualError(t, err, "validation signature is not valid")
}

/*
Scenario: Create a new node validation with an invalid status
	Given public key, signature, timestamp and an invalid validation status
	When I want to create a node validation
	Then I get an error
*/
func TestNewTransactionValidationWithInvalidStatus(t *testing.T) {
	pv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

	sig, _ := pv.Sign([]byte("hello"))

	_, err := NewValidation(10, time.Now(), pub, sig)
	assert.EqualError(t, err, "validation status is not allowed")
}

/*
Scenario: Create a new master validation
	Given a proof of work and node validation
	When I want to create the master validation
	Then I get it
*/
func TestNewTransactionMasterValidation(t *testing.T) {
	pv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

	b, _ := json.Marshal(Validation{
		nodePubk:  pub,
		status:    ValidationOK,
		timestamp: time.Now(),
	})
	sig, _ := pv.Sign(b)

	v, _ := NewValidation(ValidationOK, time.Now(), pub, sig)

	wHeaders := []NodeHeader{NewNodeHeader(pub, false, false, 0, true)}
	vHeaders := []NodeHeader{NewNodeHeader(pub, false, false, 0, true)}
	sHeaders := []NodeHeader{NewNodeHeader(pub, false, false, 0, true)}
	masterValid, err := NewMasterValidation([]crypto.PublicKey{}, pub, v, wHeaders, vHeaders, sHeaders)
	assert.Nil(t, err)
	assert.Equal(t, pub, masterValid.ProofOfWork())
	assert.Equal(t, v.PublicKey(), masterValid.Validation().PublicKey())
	assert.Equal(t, v.Timestamp(), masterValid.Validation().Timestamp())
	assert.Empty(t, masterValid.PreviousValidationNodes())
}

/*
Scenario: Create a master validation with POW invalid
	Given a no POW or not hex or invalid public key
	When I want to create master validation
	Then I get an error
*/
func TestCreateMasterWithInvalidPOW(t *testing.T) {

	_, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

	wHeaders := []NodeHeader{NewNodeHeader(pub, false, false, 0, true)}
	vHeaders := []NodeHeader{NewNodeHeader(pub, false, false, 0, true)}
	sHeaders := []NodeHeader{NewNodeHeader(pub, false, false, 0, true)}

	_, err := NewMasterValidation([]crypto.PublicKey{}, nil, Validation{}, wHeaders, vHeaders, sHeaders)
	assert.EqualError(t, err, "proof of work is missing")
}

/*
Scenario: Create a master validation without node validation
	Given a no validation
	When I want to create master validation
	Then I get an error
*/
func TestCreateMasterWithoutValidation(t *testing.T) {
	_, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
	wHeaders := []NodeHeader{NewNodeHeader(pub, false, false, 0, true)}
	vHeaders := []NodeHeader{NewNodeHeader(pub, false, false, 0, true)}
	sHeaders := []NodeHeader{NewNodeHeader(pub, false, false, 0, true)}
	_, err := NewMasterValidation([]crypto.PublicKey{}, pub, Validation{}, wHeaders, vHeaders, sHeaders)
	assert.Contains(t, err.Error(), "master validation is not valid:")
}

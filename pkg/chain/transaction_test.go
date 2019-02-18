package chain

import (
	"encoding/hex"
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
func TestNewTransaction(t *testing.T) {
	pub, pv := crypto.GenerateKeys()

	prop, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("encPvKey")), pub)

	addr := crypto.HashString("addr")
	data := map[string]string{
		"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
		"encrypted_address_by_id":   hex.EncodeToString([]byte("addr")),
		"encrypted_aes_key":         hex.EncodeToString([]byte("aesKey")),
	}
	hash := crypto.HashString("hash")

	sig, _ := crypto.Sign("data", pv)

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
func TestNewWithInvalidAddress(t *testing.T) {
	_, err := NewTransaction("", KeychainTransactionType, map[string]string{}, time.Now(), "", shared.EmitterKeyPair{}, "", "", "")
	assert.EqualError(t, err, "transaction: addr hash is empty")

	_, err = NewTransaction("abc", KeychainTransactionType, map[string]string{}, time.Now(), "", shared.EmitterKeyPair{}, "", "", "")
	assert.EqualError(t, err, "transaction: addr hash is not in hexadecimal format")

	_, err = NewTransaction(hex.EncodeToString([]byte("abc")), KeychainTransactionType, map[string]string{}, time.Now(), "", shared.EmitterKeyPair{}, "", "", "")
	assert.EqualError(t, err, "transaction: addr hash is not valid")
}

/*
Scenario: Create a new transaction with an invalid transaction hash
	Given a invalid transaction hash, empty or not in he
	When I want to create the transaction
	Then I get an error
*/
func TestNewWithInvalidHash(t *testing.T) {

	_, err := NewTransaction(crypto.HashString("addr"), KeychainTransactionType, map[string]string{}, time.Now(), "", shared.EmitterKeyPair{}, "", "", "")
	assert.EqualError(t, err, "transaction: hash is empty")

	_, err = NewTransaction(crypto.HashString("addr"), KeychainTransactionType, map[string]string{}, time.Now(), "", shared.EmitterKeyPair{}, "", "", "abc")
	assert.EqualError(t, err, "transaction: hash is not in hexadecimal format")

	_, err = NewTransaction(crypto.HashString("addr"), KeychainTransactionType, map[string]string{}, time.Now(), "", shared.EmitterKeyPair{}, "", "", hex.EncodeToString([]byte("abc")))
	assert.EqualError(t, err, "transaction: hash is not valid")
}

/*
Scenario: Create a new transaction with an invalid transaction data
	Given a empty data
	When I want to create the transaction
	Then I get an error
*/
func TestNewWithInvalidData(t *testing.T) {

	_, err := NewTransaction(crypto.HashString("addr"), KeychainTransactionType, map[string]string{}, time.Now(), "", shared.EmitterKeyPair{}, "", "", crypto.HashString("addr"))
	assert.EqualError(t, err, "transaction: data is empty")
}

/*
Scenario: Create a new transaction with an invalid transaction timestamp (more than the current timestamp)
	Given a transaction timestamp (now + 2 seconds)
	When I want to create the transaction
	Then I get an error
*/
func TestNewWithInvalidTimestamp(t *testing.T) {
	_, err := NewTransaction(crypto.HashString("addr"), KeychainTransactionType, map[string]string{
		"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
		"encrypted_wallet":          hex.EncodeToString([]byte("wallet")),
	}, time.Now().Add(2*time.Second), "", shared.EmitterKeyPair{}, "", "", crypto.HashString("addr"))
	assert.EqualError(t, err, "transaction: timestamp must be greater lower than now")
}

/*
Scenario: Create a new transaction with an invalid transaction public key
	Given an invalid public key: empty or not hex or not a key
	When I want to create the transaction
	Then I get an error
*/
func TestNewWithInvalidPublicKey(t *testing.T) {
	_, err := NewTransaction(crypto.HashString("addr"), KeychainTransactionType, map[string]string{
		"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
		"encrypted_wallet":          hex.EncodeToString([]byte("wallet")),
	}, time.Now(), "", shared.EmitterKeyPair{}, "", "", crypto.HashString("addr"))
	assert.EqualError(t, err, "transaction: public key is empty")

	_, err = NewTransaction(crypto.HashString("addr"), KeychainTransactionType, map[string]string{
		"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
		"encrypted_wallet":          hex.EncodeToString([]byte("wallet")),
	}, time.Now(), "abc", shared.EmitterKeyPair{}, "", "", crypto.HashString("addr"))
	assert.EqualError(t, err, "transaction: public key is not in hexadecimal format")

	_, err = NewTransaction(crypto.HashString("addr"), KeychainTransactionType, map[string]string{
		"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
		"encrypted_wallet":          hex.EncodeToString([]byte("wallet")),
	}, time.Now(), hex.EncodeToString([]byte("abc")), shared.EmitterKeyPair{}, "", "", crypto.HashString("addr"))
	assert.EqualError(t, err, "transaction: public key is not valid")
}

/*
Scenario: Create a new transaction with an invalid transaction signature
	Given an invalid signature: empty or not hex or not a signature
	When I want to create the transaction
	Then I get an error
*/
func TestNewWithInvalidSignature(t *testing.T) {
	pub, _ := crypto.GenerateKeys()

	_, err := NewTransaction(crypto.HashString("addr"), KeychainTransactionType, map[string]string{
		"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
		"encrypted_wallet":          hex.EncodeToString([]byte("wallet")),
	}, time.Now(), pub, shared.EmitterKeyPair{}, "", "", crypto.HashString(("hello")))
	assert.EqualError(t, err, "transaction: signature is empty")

	_, err = NewTransaction(crypto.HashString("addr"), KeychainTransactionType, map[string]string{
		"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
		"encrypted_wallet":          hex.EncodeToString([]byte("wallet")),
	}, time.Now(), pub, shared.EmitterKeyPair{}, "abc", "", crypto.HashString("addr"))
	assert.EqualError(t, err, "transaction: signature is not in hexadecimal format")

	_, err = NewTransaction(crypto.HashString("addr"), KeychainTransactionType, map[string]string{
		"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
		"encrypted_wallet":          hex.EncodeToString([]byte("wallet")),
	}, time.Now(), pub, shared.EmitterKeyPair{}, hex.EncodeToString([]byte("abc")), "", crypto.HashString("addr"))
	assert.EqualError(t, err, "transaction: signature is not valid")
}

/*
Scenario: Create a new transaction with invalid transaction type
	Given an invalid type
	When I want to create the transaction
	Then I get an error
*/
func TestNewWithInvalidType(t *testing.T) {
	pub, pv := crypto.GenerateKeys()

	sig, _ := crypto.Sign("sig", pv)

	_, err := NewTransaction(crypto.HashString("addr"), 10, map[string]string{
		"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
		"encrypted_wallet":          hex.EncodeToString([]byte("wallet")),
	}, time.Now(), pub, shared.EmitterKeyPair{}, sig, sig, crypto.HashString(("hello")))
	assert.EqualError(t, err, "transaction: type not allowed")
}

/*
Scenario: Create a new transaction without proposal
	Given a transaction without proposal
	When I want to create the transaction
	Then I get an error
*/
func TestNewWithoutProposal(t *testing.T) {
	pub, pv := crypto.GenerateKeys()

	sig, _ := crypto.Sign("sig", pv)

	_, err := NewTransaction(crypto.HashString("addr"), KeychainTransactionType, map[string]string{
		"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
		"encrypted_wallet":          hex.EncodeToString([]byte("wallet")),
	}, time.Now(), pub, shared.EmitterKeyPair{}, sig, sig, crypto.HashString(("hello")))
	assert.EqualError(t, err, "transaction: proposal is missing")
}

/*
Scenario: Check the transaction integrity
	Given a transaction with a valid hash and valid signature
	When I want to check its intergrity, its check the transaction hash and the signature
	Then I get not error
*/
func TestCheckTransactionIntegrity(t *testing.T) {
	pub, pv := crypto.GenerateKeys()

	prop, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("encpvkey")), pub)

	txRaw := Transaction{
		addr: crypto.HashString("addr"),
		data: map[string]string{
			"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
			"encrypted_wallet":          hex.EncodeToString([]byte("wallet")),
		},
		txType:    KeychainTransactionType,
		timestamp: time.Now(),
		pubKey:    pub,
		prop:      prop,
	}
	txBytesBeforeSig, _ := txRaw.MarshalBeforeSignature()
	sig, _ := crypto.Sign(string(txBytesBeforeSig), pv)
	txRaw.emSig = sig
	txRaw.sig = sig
	txBytes, _ := txRaw.MarshalHash()

	hash := crypto.HashBytes(txBytes)

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
	pub, pv := crypto.GenerateKeys()

	prop, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("encpvkey")), pub)

	raw, _ := json.Marshal(Transaction{
		addr: crypto.HashString("addr"),
		data: map[string]string{
			"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
			"encrypted_wallet":          hex.EncodeToString([]byte("wallet")),
		},
		txType:    KeychainTransactionType,
		timestamp: time.Now(),
		pubKey:    pub,
		prop:      prop,
	})
	sig, _ := crypto.Sign(string(raw), pv)
	hash := "abc"

	tx, _ := NewTransaction(crypto.HashString("addr"), KeychainTransactionType, map[string]string{
		"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
		"encrypted_wallet":          hex.EncodeToString([]byte("wallet")),
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
	pub, pv := crypto.GenerateKeys()

	prop, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("encpvkey")), pub)

	txRaw := Transaction{
		addr: crypto.HashString("addr"),
		data: map[string]string{
			"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
			"encrypted_wallet":          hex.EncodeToString([]byte("wallet")),
		},
		txType:    KeychainTransactionType,
		timestamp: time.Now(),
		pubKey:    pub,
		prop:      prop,
	}
	sig, _ := crypto.Sign(string("fake sig"), pv)
	txRaw.emSig = sig
	txRaw.sig = sig

	txBytes, _ := txRaw.MarshalHash()

	hash := crypto.HashBytes(txBytes)

	tx, _ := NewTransaction(crypto.HashString("addr"), KeychainTransactionType, map[string]string{
		"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
		"encrypted_wallet":          hex.EncodeToString([]byte("wallet")),
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
	pub, pv := crypto.GenerateKeys()

	prop, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("encpvkey")), pub)

	raw, _ := json.Marshal(Transaction{
		addr: crypto.HashString("addr"),
		data: map[string]string{
			"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
			"encrypted_wallet":          hex.EncodeToString([]byte("wallet")),
		},
		txType:    KeychainTransactionType,
		timestamp: time.Now(),
		pubKey:    pub,
		prop:      prop,
	})
	sig, _ := crypto.Sign(string(raw), pv)
	hash := crypto.HashBytes(raw)
	tx, _ := NewTransaction(crypto.HashString("addr"), KeychainTransactionType, map[string]string{
		"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
		"encrypted_wallet":          hex.EncodeToString([]byte("wallet")),
	}, time.Now(), pub, prop, sig, sig, hash)

	b, _ := json.Marshal(Validation{
		nodePubk:  pub,
		status:    ValidationOK,
		timestamp: time.Now(),
	})
	sig, _ = crypto.Sign(string(b), pv)
	v, _ := NewValidation(ValidationOK, time.Now(), pub, sig)

	masterValid, _ := NewMasterValidation([]string{}, pub, v)
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
	pub, pv := crypto.GenerateKeys()

	prop, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("encpvkey")), pub)

	raw, _ := json.Marshal(Transaction{
		addr: crypto.HashString("addr"),
		data: map[string]string{
			"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
			"encrypted_wallet":          hex.EncodeToString([]byte("wallet")),
		},
		txType:    KeychainTransactionType,
		timestamp: time.Now(),
		pubKey:    pub,
		prop:      prop,
	})
	sig, _ := crypto.Sign(string(raw), pv)
	hash := crypto.HashBytes(raw)
	tx, _ := NewTransaction(crypto.HashString("addr"), KeychainTransactionType, map[string]string{
		"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
		"encrypted_wallet":          hex.EncodeToString([]byte("wallet")),
	}, time.Now(), pub, prop, sig, sig, hash)

	b, _ := json.Marshal(Validation{
		nodePubk:  pub,
		status:    ValidationOK,
		timestamp: time.Now(),
	})
	sig, _ = crypto.Sign(string(b), pv)
	v, _ := NewValidation(ValidationOK, time.Now(), pub, sig)

	masterValid, _ := NewMasterValidation([]string{}, pub, v)
	assert.EqualError(t, tx.Mined(masterValid, []Validation{}), "transaction: missing confirmation validations")
}

/*
Scenario: Check the integrity of a transaction chain
	Given 3 transactions chained
	When I want to check their the chain integrity
	Then I get not error
*/
func TestCheckChainIntegrity(t *testing.T) {
	pub, pv := crypto.GenerateKeys()

	prop, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("encpvkey")), pub)

	tx1 := Transaction{
		addr: crypto.HashString("addr"),
		data: map[string]string{
			"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
			"encrypted_wallet":          hex.EncodeToString([]byte("wallet")),
		},
		txType:    KeychainTransactionType,
		timestamp: time.Now(),
		pubKey:    pub,
		prop:      prop,
	}
	txBytesBeforeSig, _ := tx1.MarshalBeforeSignature()
	sig, _ := crypto.Sign(string(txBytesBeforeSig), pv)
	tx1.emSig = sig
	tx1.sig = sig
	txBytes1, _ := tx1.MarshalHash()
	txHash1 := crypto.HashBytes(txBytes1)
	tx1.hash = txHash1

	tx2 := Transaction{
		addr: crypto.HashString("addr"),
		data: map[string]string{
			"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
			"encrypted_wallet":          hex.EncodeToString([]byte("wallet")),
		},
		txType:    KeychainTransactionType,
		timestamp: time.Now().Add(2 * time.Second),
		pubKey:    pub,
		prop:      prop,
	}
	txBytesBeforeSig2, _ := tx2.MarshalBeforeSignature()
	sig2, _ := crypto.Sign(string(txBytesBeforeSig2), pv)
	tx2.emSig = sig2
	tx2.sig = sig2
	txBytes2, _ := tx2.MarshalHash()
	txHash2 := crypto.HashBytes(txBytes2)
	tx2.hash = txHash2
	tx2.prevTx = &tx1

	tx3 := Transaction{
		addr: crypto.HashString("hello3"),
		data: map[string]string{
			"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
			"encrypted_wallet":          hex.EncodeToString([]byte("wallet")),
		},
		txType:    KeychainTransactionType,
		timestamp: time.Now().Add(3 * time.Second),
		pubKey:    pub,
		prop:      prop,
	}
	txBytesBeforeSig3, _ := tx3.MarshalBeforeSignature()
	sig3, _ := crypto.Sign(string(txBytesBeforeSig3), pv)
	tx3.emSig = sig3
	tx3.sig = sig3
	txBytes3, _ := tx3.MarshalHash()
	txHash3 := crypto.HashBytes(txBytes3)
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
	pub, pv := crypto.GenerateKeys()

	prop, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("encpvkey")), pub)

	tx0 := Transaction{
		addr: crypto.HashString("addr"),
		data: map[string]string{
			"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
			"encrypted_wallet":          hex.EncodeToString([]byte("wallet")),
		},
		txType:    KeychainTransactionType,
		timestamp: time.Now(),
		pubKey:    pub,
		prop:      prop,
	}

	b, _ := json.Marshal(tx0)
	hash := crypto.HashBytes(b)
	sig, _ := crypto.Sign(string(b), pv)
	tx0.sig = sig
	tx0.hash = hash

	tx1 := Transaction{
		addr: crypto.HashString("hello2"),
		data: map[string]string{
			"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
			"encrypted_wallet":          hex.EncodeToString([]byte("wallet")),
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

	pub, pv := crypto.GenerateKeys()

	prop, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("encpvkey")), pub)

	tx1 := Transaction{
		addr: crypto.HashString("addr"),
		data: map[string]string{
			"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
			"encrypted_wallet":          hex.EncodeToString([]byte("wallet")),
		},
		txType:    KeychainTransactionType,
		timestamp: time.Now(),
		pubKey:    pub,
		prop:      prop,
	}
	txBytesBeforeSig, _ := tx1.MarshalBeforeSignature()
	sig, _ := crypto.Sign(string(txBytesBeforeSig), pv)
	tx1.emSig = sig
	tx1.sig = sig
	txBytes1, _ := tx1.MarshalHash()
	txHash1 := crypto.HashBytes(txBytes1)
	tx1.hash = txHash1

	time.Sleep(1 * time.Second)

	tx2 := Transaction{
		addr: crypto.HashString("addr"),
		data: map[string]string{
			"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
			"encrypted_wallet":          hex.EncodeToString([]byte("wallet")),
		},
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
	txHash2 := crypto.HashBytes(txBytes2)
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

	pub, pv := crypto.GenerateKeys()

	prop, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("encpvkey")), pub)

	txTime1 := time.Now()
	raw1, _ := json.Marshal(Transaction{
		addr: crypto.HashString("addr"),
		data: map[string]string{
			"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
			"encrypted_wallet":          hex.EncodeToString([]byte("wallet")),
		},
		txType:    KeychainTransactionType,
		timestamp: txTime1,
		pubKey:    pub,
		prop:      prop,
	})
	sig1, _ := crypto.Sign(string(raw1), pv)
	hash1 := crypto.HashBytes(raw1)

	tx1, _ := NewTransaction(crypto.HashString("addr"), KeychainTransactionType, map[string]string{
		"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
		"encrypted_wallet":          hex.EncodeToString([]byte("wallet")),
	}, txTime1, pub, prop, sig1, sig1, hash1)

	raw2, _ := json.Marshal(Transaction{
		addr: crypto.HashString("addr"),
		data: map[string]string{
			"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
			"encrypted_wallet":          hex.EncodeToString([]byte("wallet")),
		},
		txType:    KeychainTransactionType,
		timestamp: txTime1,
		pubKey:    pub,
		prop:      prop,
	})
	sig2, _ := crypto.Sign(string(raw2), pv)
	hash2 := crypto.HashBytes(raw2)

	tx2, _ := NewTransaction(crypto.HashString("addr"), KeychainTransactionType, map[string]string{
		"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
		"encrypted_wallet":          hex.EncodeToString([]byte("wallet")),
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

	pub, pv := crypto.GenerateKeys()

	v := Validation{
		nodePubk:  pub,
		status:    ValidationOK,
		timestamp: time.Now(),
	}
	b, _ := json.Marshal(v)
	sig, _ := crypto.Sign(string(b), pv)
	v.nodeSig = sig

	prop, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("encpvkey")), pub)

	tx := Transaction{
		addr: crypto.HashString("addr"),
		data: map[string]string{
			"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
			"encrypted_wallet":          hex.EncodeToString([]byte("wallet")),
		},
		txType:    KeychainTransactionType,
		timestamp: time.Now(),
		pubKey:    pub,
		prop:      prop,
	}
	txBytesBeforeSig, _ := tx.MarshalBeforeSignature()
	sig, _ = crypto.Sign(string(txBytesBeforeSig), pv)
	tx.sig = sig
	txBytesBeforeEmSig, _ := tx.MarshalBeforeEmitterSignature()
	emSig, _ := crypto.Sign(string(txBytesBeforeEmSig), pv)
	tx.emSig = emSig
	txBytes, _ := tx.MarshalHash()
	txHash := crypto.HashBytes(txBytes)
	tx.hash = txHash
	tx.masterV = MasterValidation{
		pow:        pub,
		validation: v,
	}

	assert.Nil(t, tx.CheckMasterValidation())
}

/*
Scenario: Check master validation with invalid POW
	Given a transaction with a master validation including bad POW
	When I want to check the master validation and the POW
	Then I get an error
*/
func TestCheckMasterValidationWithInvalidPOW(t *testing.T) {

	pub, pv := crypto.GenerateKeys()
	pub2, _ := crypto.GenerateKeys()

	v := Validation{
		nodePubk:  pub,
		status:    ValidationOK,
		timestamp: time.Now(),
	}
	b, _ := json.Marshal(v)
	sig, _ := crypto.Sign(string(b), pv)
	v.nodeSig = sig

	prop, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("encpvkey")), pub)

	raw, _ := json.Marshal(Transaction{
		addr: crypto.HashString("addr"),
		data: map[string]string{
			"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
			"encrypted_wallet":          hex.EncodeToString([]byte("wallet")),
		},
		txType:    KeychainTransactionType,
		timestamp: time.Now(),
		pubKey:    pub,
		prop:      prop,
	})
	sigTx, _ := crypto.Sign(string(raw), pv)
	hash := crypto.HashBytes(raw)

	tx := Transaction{
		masterV: MasterValidation{
			pow:        pub2,
			validation: v,
		},
		addr: crypto.HashString("addr"),
		data: map[string]string{
			"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
			"encrypted_wallet":          hex.EncodeToString([]byte("wallet")),
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
func TestNewValidation(t *testing.T) {

	pub, pv := crypto.GenerateKeys()

	b, _ := json.Marshal(Validation{
		nodePubk:  pub,
		status:    ValidationOK,
		timestamp: time.Now(),
	})
	sig, _ := crypto.Sign(string(b), pv)

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
func TestNewValidationWithInvalidTimestamp(t *testing.T) {
	_, err := NewValidation(ValidationOK, time.Now().Add(2*time.Second), "", "")
	assert.EqualError(t, err, "node validation: timestamp must be anterior or equal to now")
}

/*
Scenario: Create a new node validation with invalid public key
	Given no public key or no hex or not valid public key
	When I want to create a node validation
	Then I get an error
*/
func TestNewValidationWithInvalidPublicKey(t *testing.T) {
	_, err := NewValidation(ValidationOK, time.Now(), "", "sig")
	assert.EqualError(t, err, "node validation: public key is empty")

	_, err = NewValidation(ValidationOK, time.Now(), "key", "sig")
	assert.EqualError(t, err, "node validation: public key is not in hexadecimal format")

	_, err = NewValidation(ValidationOK, time.Now(), hex.EncodeToString([]byte("key")), "sig")
	assert.EqualError(t, err, "node validation: public key is not valid")
}

/*
Scenario: Create a new node validation with invalid signature
	Given no hex or not valid signature
	When I want to create a node validation
	Then I get an error
*/
func TestNewValidationWithInvalidSignature(t *testing.T) {

	pub, pv := crypto.GenerateKeys()

	_, err := NewValidation(ValidationOK, time.Now(), pub, "sig")
	assert.EqualError(t, err, "node validation: signature is not in hexadecimal format")

	_, err = NewValidation(ValidationOK, time.Now(), pub, hex.EncodeToString([]byte("sig")))
	assert.EqualError(t, err, "node validation: signature is not valid")

	sig, _ := crypto.Sign("hello", pv)
	_, err = NewValidation(ValidationOK, time.Now(), pub, sig)
	assert.EqualError(t, err, "node validation: signature is invalid")
}

/*
Scenario: Create a new node validation with an invalid status
	Given public key, signature, timestamp and an invalid validation status
	When I want to create a node validation
	Then I get an error
*/
func TestNewValidationWithInvalidStatus(t *testing.T) {
	pub, pv := crypto.GenerateKeys()

	sig, _ := crypto.Sign("hello", pv)

	_, err := NewValidation(10, time.Now(), pub, sig)
	assert.EqualError(t, err, "node validation: status not allowed")
}

/*
Scenario: Create a new master validation
	Given a proof of work and node validation
	When I want to create the master validation
	Then I get it
*/
func TestNewMasterValidation(t *testing.T) {
	pub, pv := crypto.GenerateKeys()

	b, _ := json.Marshal(Validation{
		nodePubk:  pub,
		status:    ValidationOK,
		timestamp: time.Now(),
	})
	sig, _ := crypto.Sign(string(b), pv)

	v, _ := NewValidation(ValidationOK, time.Now(), pub, sig)
	mv, err := NewMasterValidation([]string{}, pub, v)
	assert.Nil(t, err)
	assert.Equal(t, pub, mv.ProofOfWork())
	assert.Equal(t, v.PublicKey(), mv.Validation().PublicKey())
	assert.Equal(t, v.Timestamp(), mv.Validation().Timestamp())
	assert.Empty(t, mv.PreviousTransactionNodes())
}

/*
Scenario: Create a master validation with POW invalid
	Given a no POW or not hex or invalid public key
	When I want to create master validation
	Then I get an error
*/
func TestCreateMasterWithInvalidPOW(t *testing.T) {
	_, err := NewMasterValidation([]string{}, "", Validation{})
	assert.EqualError(t, err, "master validation POW: public key is empty")

	_, err = NewMasterValidation([]string{}, "key", Validation{})
	assert.EqualError(t, err, "master validation POW: public key is not in hexadecimal format")

	_, err = NewMasterValidation([]string{}, hex.EncodeToString([]byte("key")), Validation{})
	assert.EqualError(t, err, "master validation POW: public key is not valid")
}

/*
Scenario: Create a master validation without node validation
	Given a no validation
	When I want to create master validation
	Then I get an error
*/
func TestCreateMasterWithoutValidation(t *testing.T) {

	pub, _ := crypto.GenerateKeys()

	_, err := NewMasterValidation([]string{}, pub, Validation{})
	assert.EqualError(t, err, "master validation: node validation: public key is empty")
}

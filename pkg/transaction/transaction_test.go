package transaction

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/uniris/uniris-core/pkg/shared"

	"github.com/stretchr/testify/assert"
	"github.com/uniris/uniris-core/pkg/crypto"
)

/*
Scenario: Create a new transaction proposal
	Given a shared key pair
	When I want to create a transaction proposal
	Then I get a proposal and I can retrieve the shared keys
*/
func TestNewProposal(t *testing.T) {
	pvKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	key, _ := x509.MarshalPKIXPublicKey(pvKey.Public())

	kp, _ := shared.NewKeyPair(hex.EncodeToString([]byte("encPvKey")), hex.EncodeToString(key))
	prop, err := NewProposal(kp)
	assert.Nil(t, err)
	assert.Equal(t, hex.EncodeToString([]byte("encPvKey")), prop.SharedEmitterKeyPair().EncryptedPrivateKey())
	assert.Equal(t, hex.EncodeToString(key), prop.SharedEmitterKeyPair().PublicKey())
}

/*
Scenario: Create a new transction proposal with an empty shared keypair
	Given an empty shared key pari
	When I want to create a transaction proposal
	Then I get an error
*/
func TestNewEmptyProposal(t *testing.T) {
	_, err := NewProposal(shared.KeyPair{})
	assert.Error(t, err, "transaction proposal: missing shared keys")
}

/*
Scenario: Marshal into a JSON a transaction proposal
	Given a transaction propoal
	When I want to marshal it into a JSON
	Then I get a valid JSON
*/
func TestMarshalProposal(t *testing.T) {
	pvKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	key, _ := x509.MarshalPKIXPublicKey(pvKey.Public())

	kp, _ := shared.NewKeyPair(hex.EncodeToString([]byte("encPvKey")), hex.EncodeToString(key))
	prop, _ := NewProposal(kp)
	b, err := json.Marshal(prop)
	assert.Nil(t, err)
	assert.Equal(t, fmt.Sprintf("{\"shared_emitter_keys\":{\"encrypted_private_key\":\"%s\",\"public_key\":\"%s\"}}", hex.EncodeToString([]byte("encPvKey")), hex.EncodeToString(key)), string(b))
}

/*
Scenario: Create a new transaction
	Given transaction data (address, hash, public key, signature, emSig, prop, timestamp)
	When I want to create the transaction
	Then I get it
*/
func TestNew(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pub, _ := x509.MarshalPKIXPublicKey(key.Public())
	pv, _ := x509.MarshalECPrivateKey(key)

	kp, _ := shared.NewKeyPair(hex.EncodeToString([]byte("encPvKey")), hex.EncodeToString(pub))
	prop, _ := NewProposal(kp)

	addr := crypto.HashString("address")
	data := map[string]string{
		"encrypted_address_by_robot": hex.EncodeToString([]byte("addr")),
		"encrypted_address_by_id":    hex.EncodeToString([]byte("addr")),
		"encrypted_aes_key":          hex.EncodeToString([]byte("aesKey")),
	}
	hash := crypto.HashString("hash")

	sig, _ := crypto.Sign("data", hex.EncodeToString(pv))

	tx, err := New(addr, KeychainType, data, time.Now(), hex.EncodeToString(pub), sig, sig, prop, hash)
	assert.Nil(t, err)
	assert.Equal(t, addr, tx.Address())
	assert.Equal(t, data, tx.Data())
	assert.Equal(t, KeychainType, tx.Type())
	assert.Equal(t, sig, tx.Signature())
	assert.Equal(t, hex.EncodeToString(pub), tx.Proposal().SharedEmitterKeyPair().PublicKey())
}

/*
Scenario: Create a new transaction with an invalid address
	Given a invalid address hash, empty or not in he
	When I want to create the transaction
	Then I get an error
*/
func TestNewWithInvalidAddress(t *testing.T) {
	_, err := New("", KeychainType, map[string]string{}, time.Now(), "", "", "", Proposal{}, "")
	assert.EqualError(t, err, "transaction: hash is empty")

	_, err = New("abc", KeychainType, map[string]string{}, time.Now(), "", "", "", Proposal{}, "")
	assert.EqualError(t, err, "transaction: hash is not in hexadecimal format")

	_, err = New(hex.EncodeToString([]byte("abc")), KeychainType, map[string]string{}, time.Now(), "", "", "", Proposal{}, "")
	assert.EqualError(t, err, "transaction: hash is not valid")
}

/*
Scenario: Create a new transaction with an invalid transaction hash
	Given a invalid transaction hash, empty or not in he
	When I want to create the transaction
	Then I get an error
*/
func TestNewWithInvalidHash(t *testing.T) {

	_, err := New(crypto.HashString("addr"), KeychainType, map[string]string{}, time.Now(), "", "", "", Proposal{}, "")
	assert.EqualError(t, err, "transaction: hash is empty")

	_, err = New(crypto.HashString("addr"), KeychainType, map[string]string{}, time.Now(), "", "", "", Proposal{}, "abc")
	assert.EqualError(t, err, "transaction: hash is not in hexadecimal format")

	_, err = New(crypto.HashString("addr"), KeychainType, map[string]string{}, time.Now(), "", "", "", Proposal{}, hex.EncodeToString([]byte("abc")))
	assert.EqualError(t, err, "transaction: hash is not valid")
}

/*
Scenario: Create a new transaction with an invalid transaction data
	Given a empty data
	When I want to create the transaction
	Then I get an error
*/
func TestNewWithInvalidData(t *testing.T) {

	_, err := New(crypto.HashString("addr"), KeychainType, map[string]string{}, time.Now(), "", "", "", Proposal{}, crypto.HashString("addr"))
	assert.EqualError(t, err, "transaction: data is empty")
}

/*
Scenario: Create a new transaction with an invalid transaction timestamp (more than the current timestamp)
	Given a transaction timestamp (now + 2 seconds)
	When I want to create the transaction
	Then I get an error
*/
func TestNewWithInvalidTimestamp(t *testing.T) {
	_, err := New(crypto.HashString("addr"), KeychainType, map[string]string{
		"encrypted_address": hex.EncodeToString([]byte("addr")),
		"encrypted_wallet":  hex.EncodeToString([]byte("wallet")),
	}, time.Now().Add(2*time.Second), "", "", "", Proposal{}, crypto.HashString("addr"))
	assert.EqualError(t, err, "transaction: timestamp must be greater lower than now")
}

/*
Scenario: Create a new transaction with an invalid transaction public key
	Given an invalid public key: empty or not hex or not a key
	When I want to create the transaction
	Then I get an error
*/
func TestNewWithInvalidPublicKey(t *testing.T) {
	_, err := New(crypto.HashString("addr"), KeychainType, map[string]string{
		"encrypted_address": hex.EncodeToString([]byte("addr")),
		"encrypted_wallet":  hex.EncodeToString([]byte("wallet")),
	}, time.Now(), "", "", "", Proposal{}, crypto.HashString("addr"))
	assert.EqualError(t, err, "transaction: public key is empty")

	_, err = New(crypto.HashString("addr"), KeychainType, map[string]string{
		"encrypted_address": hex.EncodeToString([]byte("addr")),
		"encrypted_wallet":  hex.EncodeToString([]byte("wallet")),
	}, time.Now(), "abc", "", "", Proposal{}, crypto.HashString("addr"))
	assert.EqualError(t, err, "transaction: public key is not in hexadecimal format")

	_, err = New(crypto.HashString("addr"), KeychainType, map[string]string{
		"encrypted_address": hex.EncodeToString([]byte("addr")),
		"encrypted_wallet":  hex.EncodeToString([]byte("wallet")),
	}, time.Now(), hex.EncodeToString([]byte("abc")), "", "", Proposal{}, crypto.HashString("addr"))
	assert.EqualError(t, err, "transaction: public key is not valid")
}

/*
Scenario: Create a new transaction with an invalid transaction signature
	Given an invalid signature: empty or not hex or not a signature
	When I want to create the transaction
	Then I get an error
*/
func TestNewWithInvalidSignature(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pub, _ := x509.MarshalPKIXPublicKey(key.Public())

	_, err := New(crypto.HashString("addr"), KeychainType, map[string]string{
		"encrypted_address": hex.EncodeToString([]byte("addr")),
		"encrypted_wallet":  hex.EncodeToString([]byte("wallet")),
	}, time.Now(), hex.EncodeToString(pub), "", "", Proposal{}, crypto.HashString(("hello")))
	assert.EqualError(t, err, "transaction: signature is empty")

	_, err = New(crypto.HashString("addr"), KeychainType, map[string]string{
		"encrypted_address": hex.EncodeToString([]byte("addr")),
		"encrypted_wallet":  hex.EncodeToString([]byte("wallet")),
	}, time.Now(), hex.EncodeToString(pub), "abc", "", Proposal{}, crypto.HashString("addr"))
	assert.EqualError(t, err, "transaction: signature is not in hexadecimal format")

	_, err = New(crypto.HashString("addr"), KeychainType, map[string]string{
		"encrypted_address": hex.EncodeToString([]byte("addr")),
		"encrypted_wallet":  hex.EncodeToString([]byte("wallet")),
	}, time.Now(), hex.EncodeToString(pub), hex.EncodeToString([]byte("abc")), "", Proposal{}, crypto.HashString("addr"))
	assert.EqualError(t, err, "transaction: signature is not valid")
}

/*
Scenario: Create a new transaction with invalid transaction type
	Given an invalid type
	When I want to create the transaction
	Then I get an error
*/
func TestNewWithInvalidType(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pub, _ := x509.MarshalPKIXPublicKey(key.Public())
	pv, _ := x509.MarshalECPrivateKey(key)

	sig, _ := crypto.Sign("sig", hex.EncodeToString(pv))

	_, err := New(crypto.HashString("addr"), 10, map[string]string{
		"encrypted_address": hex.EncodeToString([]byte("addr")),
		"encrypted_wallet":  hex.EncodeToString([]byte("wallet")),
	}, time.Now(), hex.EncodeToString(pub), sig, sig, Proposal{}, crypto.HashString(("hello")))
	assert.EqualError(t, err, "transaction: type not allowed")
}

/*
Scenario: Create a new transaction without proposal
	Given a transaction without proposal
	When I want to create the transaction
	Then I get an error
*/
func TestNewWithoutProposal(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pub, _ := x509.MarshalPKIXPublicKey(key.Public())
	pv, _ := x509.MarshalECPrivateKey(key)

	sig, _ := crypto.Sign("sig", hex.EncodeToString(pv))

	_, err := New(crypto.HashString("addr"), KeychainType, map[string]string{
		"encrypted_address": hex.EncodeToString([]byte("addr")),
		"encrypted_wallet":  hex.EncodeToString([]byte("wallet")),
	}, time.Now(), hex.EncodeToString(pub), sig, sig, Proposal{}, crypto.HashString(("hello")))
	assert.EqualError(t, err, "transaction: proposal is missing")
}

/*
Scenario: Check the transaction integrity
	Given a transaction with a valid hash and valid signature
	When I want to check its intergrity, its check the transaction hash and the signature
	Then I get not error
*/
func TestCheckTransactionIntegrity(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pub, _ := x509.MarshalPKIXPublicKey(key.Public())
	pv, _ := x509.MarshalECPrivateKey(key)

	sk, _ := shared.NewKeyPair(hex.EncodeToString([]byte("encpvkey")), hex.EncodeToString(pub))
	prop, _ := NewProposal(sk)

	txRaw := Transaction{
		address: crypto.HashString("addr"),
		data: map[string]string{
			"encrypted_address": hex.EncodeToString([]byte("addr")),
			"encrypted_wallet":  hex.EncodeToString([]byte("wallet")),
		},
		txType:    KeychainType,
		timestamp: time.Now(),
		pubKey:    hex.EncodeToString(pub),
		prop:      prop,
	}
	txBytesBeforeSig, _ := txRaw.MarshalBeforeSignature()
	sig, _ := crypto.Sign(string(txBytesBeforeSig), hex.EncodeToString(pv))
	txRaw.emSig = sig
	txRaw.sig = sig
	txBytes, _ := txRaw.MarshalHash()

	hash := crypto.HashBytes(txBytes)

	tx, _ := New(txRaw.address, KeychainType, txRaw.data, txRaw.timestamp, txRaw.pubKey, txRaw.sig, txRaw.emSig, txRaw.prop, hash)
	assert.Nil(t, tx.checkTransactionIntegrity())
}

/*
Scenario: Check the transaction integrity with invalid hash
	Given a transaction with a invalid hash
	When I want to check its intergrity, its check the transaction hash and the signature
	Then I get an error
*/
func TestCheckTransactionIntegrityWithInvalidHash(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pub, _ := x509.MarshalPKIXPublicKey(key.Public())
	pv, _ := x509.MarshalECPrivateKey(key)

	sk, _ := shared.NewKeyPair(hex.EncodeToString([]byte("encpvkey")), hex.EncodeToString(pub))
	prop, _ := NewProposal(sk)

	raw, _ := json.Marshal(Transaction{
		address: crypto.HashString("addr"),
		data: map[string]string{
			"encrypted_address": hex.EncodeToString([]byte("addr")),
			"encrypted_wallet":  hex.EncodeToString([]byte("wallet")),
		},
		txType:    KeychainType,
		timestamp: time.Now(),
		pubKey:    hex.EncodeToString(pub),
		prop:      prop,
	})
	sig, _ := crypto.Sign(string(raw), hex.EncodeToString(pv))
	hash := "abc"

	tx, _ := New(crypto.HashString("addr"), KeychainType, map[string]string{
		"encrypted_address": hex.EncodeToString([]byte("addr")),
		"encrypted_wallet":  hex.EncodeToString([]byte("wallet")),
	}, time.Now(), hex.EncodeToString(pub), sig, sig, prop, hash)
	assert.EqualError(t, tx.checkTransactionIntegrity(), "transaction integrity violated")
}

/*
Scenario: Check the transaction integrity with invalid signature
	Given a transaction with a valid hash and invalid signature
	When I want to check its intergrity, its check the transaction hash and the signature
	Then I get not error
*/
func TestCheckTransactionIntegrityWithInvalidSignature(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pub, _ := x509.MarshalPKIXPublicKey(key.Public())
	pv, _ := x509.MarshalECPrivateKey(key)

	sk, _ := shared.NewKeyPair(hex.EncodeToString([]byte("encpvkey")), hex.EncodeToString(pub))
	prop, _ := NewProposal(sk)

	txRaw := Transaction{
		address: crypto.HashString("addr"),
		data: map[string]string{
			"encrypted_address": hex.EncodeToString([]byte("addr")),
			"encrypted_wallet":  hex.EncodeToString([]byte("wallet")),
		},
		txType:    KeychainType,
		timestamp: time.Now(),
		pubKey:    hex.EncodeToString(pub),
		prop:      prop,
	}
	sig, _ := crypto.Sign(string("fake sig"), hex.EncodeToString(pv))
	txRaw.emSig = sig
	txRaw.sig = sig

	txBytes, _ := txRaw.MarshalHash()

	hash := crypto.HashBytes(txBytes)

	tx, _ := New(crypto.HashString("addr"), KeychainType, map[string]string{
		"encrypted_address": hex.EncodeToString([]byte("addr")),
		"encrypted_wallet":  hex.EncodeToString([]byte("wallet")),
	}, time.Now(), hex.EncodeToString(pub), sig, sig, prop, hash)
	assert.EqualError(t, tx.checkTransactionIntegrity(), "transaction signature invalid")
}

/*
Scenario: Add mining information to a transaction
	Given a transaction
	When I want to add master validation and confirmation validations
	Then I can retrieve in inside the transaction
*/
func TestAddMining(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pub, _ := x509.MarshalPKIXPublicKey(key.Public())
	pv, _ := x509.MarshalECPrivateKey(key)

	sk, _ := shared.NewKeyPair(hex.EncodeToString([]byte("encpvkey")), hex.EncodeToString(pub))
	prop, _ := NewProposal(sk)

	raw, _ := json.Marshal(Transaction{
		address: crypto.HashString("addr"),
		data: map[string]string{
			"encrypted_address": hex.EncodeToString([]byte("addr")),
			"encrypted_wallet":  hex.EncodeToString([]byte("wallet")),
		},
		txType:    KeychainType,
		timestamp: time.Now(),
		pubKey:    hex.EncodeToString(pub),
		prop:      prop,
	})
	sig, _ := crypto.Sign(string(raw), hex.EncodeToString(pv))
	hash := crypto.HashBytes(raw)
	tx, _ := New(crypto.HashString("addr"), KeychainType, map[string]string{
		"encrypted_address": hex.EncodeToString([]byte("addr")),
		"encrypted_wallet":  hex.EncodeToString([]byte("wallet")),
	}, time.Now(), hex.EncodeToString(pub), sig, sig, prop, hash)

	b, _ := json.Marshal(MinerValidation{
		minerPubk: hex.EncodeToString(pub),
		status:    ValidationOK,
		timestamp: time.Now(),
	})
	sig, _ = crypto.Sign(string(b), hex.EncodeToString(pv))
	v, _ := NewMinerValidation(ValidationOK, time.Now(), hex.EncodeToString(pub), sig)

	masterValid, _ := NewMasterValidation(Pool{}, hex.EncodeToString(pub), v)
	assert.Nil(t, tx.AddMining(masterValid, []MinerValidation{v}))

	assert.Equal(t, sig, tx.MasterValidation().Validation().MinerSignature())
	assert.Len(t, tx.ConfirmationsValidations(), 1)
}

/*
Scenario: Add mining information to a transaction without confirmations
	Given a transaction
	When I want to add master validation
	Then I get an error
*/
func TestAddMiningWithoutConfirmations(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pub, _ := x509.MarshalPKIXPublicKey(key.Public())
	pv, _ := x509.MarshalECPrivateKey(key)

	sk, _ := shared.NewKeyPair(hex.EncodeToString([]byte("encpvkey")), hex.EncodeToString(pub))
	prop, _ := NewProposal(sk)

	raw, _ := json.Marshal(Transaction{
		address: crypto.HashString("addr"),
		data: map[string]string{
			"encrypted_address": hex.EncodeToString([]byte("addr")),
			"encrypted_wallet":  hex.EncodeToString([]byte("wallet")),
		},
		txType:    KeychainType,
		timestamp: time.Now(),
		pubKey:    hex.EncodeToString(pub),
		prop:      prop,
	})
	sig, _ := crypto.Sign(string(raw), hex.EncodeToString(pv))
	hash := crypto.HashBytes(raw)
	tx, _ := New(crypto.HashString("addr"), KeychainType, map[string]string{
		"encrypted_address": hex.EncodeToString([]byte("addr")),
		"encrypted_wallet":  hex.EncodeToString([]byte("wallet")),
	}, time.Now(), hex.EncodeToString(pub), sig, sig, prop, hash)

	b, _ := json.Marshal(MinerValidation{
		minerPubk: hex.EncodeToString(pub),
		status:    ValidationOK,
		timestamp: time.Now(),
	})
	sig, _ = crypto.Sign(string(b), hex.EncodeToString(pv))
	v, _ := NewMinerValidation(ValidationOK, time.Now(), hex.EncodeToString(pub), sig)

	masterValid, _ := NewMasterValidation(Pool{}, hex.EncodeToString(pub), v)
	assert.EqualError(t, tx.AddMining(masterValid, []MinerValidation{}), "transaction: missing confirmation validations")
}

/*
Scenario: Check the integrity of a transaction chain
	Given 3 transactions chained
	When I want to check their the chain integrity
	Then I get not error
*/
func TestCheckChainIntegrity(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pub, _ := x509.MarshalPKIXPublicKey(key.Public())
	pv, _ := x509.MarshalECPrivateKey(key)

	sk, _ := shared.NewKeyPair(hex.EncodeToString([]byte("encpvkey")), hex.EncodeToString(pub))
	prop, _ := NewProposal(sk)

	tx1 := Transaction{
		address: crypto.HashString("addr"),
		data: map[string]string{
			"encrypted_address": hex.EncodeToString([]byte("addr")),
			"encrypted_wallet":  hex.EncodeToString([]byte("wallet")),
		},
		txType:    KeychainType,
		timestamp: time.Now(),
		pubKey:    hex.EncodeToString(pub),
		prop:      prop,
	}
	txBytesBeforeSig, _ := tx1.MarshalBeforeSignature()
	sig, _ := crypto.Sign(string(txBytesBeforeSig), hex.EncodeToString(pv))
	tx1.emSig = sig
	tx1.sig = sig
	txBytes1, _ := tx1.MarshalHash()
	txHash1 := crypto.HashBytes(txBytes1)
	tx1.txHash = txHash1

	tx2 := Transaction{
		address: crypto.HashString("addr"),
		data: map[string]string{
			"encrypted_address": hex.EncodeToString([]byte("addr")),
			"encrypted_wallet":  hex.EncodeToString([]byte("wallet")),
		},
		txType:    KeychainType,
		timestamp: time.Now().Add(2 * time.Second),
		pubKey:    hex.EncodeToString(pub),
		prop:      prop,
	}
	txBytesBeforeSig2, _ := tx2.MarshalBeforeSignature()
	sig2, _ := crypto.Sign(string(txBytesBeforeSig2), hex.EncodeToString(pv))
	tx2.emSig = sig2
	tx2.sig = sig2
	txBytes2, _ := tx2.MarshalHash()
	txHash2 := crypto.HashBytes(txBytes2)
	tx2.txHash = txHash2
	tx2.prevTx = &tx1

	tx3 := Transaction{
		address: crypto.HashString("hello3"),
		data: map[string]string{
			"encrypted_address": hex.EncodeToString([]byte("addr")),
			"encrypted_wallet":  hex.EncodeToString([]byte("wallet")),
		},
		txType:    KeychainType,
		timestamp: time.Now().Add(3 * time.Second),
		pubKey:    hex.EncodeToString(pub),
		prop:      prop,
	}
	txBytesBeforeSig3, _ := tx3.MarshalBeforeSignature()
	sig3, _ := crypto.Sign(string(txBytesBeforeSig3), hex.EncodeToString(pv))
	tx3.emSig = sig3
	tx3.sig = sig3
	txBytes3, _ := tx3.MarshalHash()
	txHash3 := crypto.HashBytes(txBytes3)
	tx3.txHash = txHash3
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
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pub, _ := x509.MarshalPKIXPublicKey(key.Public())
	pv, _ := x509.MarshalECPrivateKey(key)

	sk, _ := shared.NewKeyPair(hex.EncodeToString([]byte("encpvkey")), hex.EncodeToString(pub))
	prop, _ := NewProposal(sk)

	tx0 := Transaction{
		address: crypto.HashString("addr"),
		data: map[string]string{
			"encrypted_address": hex.EncodeToString([]byte("addr")),
			"encrypted_wallet":  hex.EncodeToString([]byte("wallet")),
		},
		txType:    KeychainType,
		timestamp: time.Now(),
		pubKey:    hex.EncodeToString(pub),
		prop:      prop,
	}

	b, _ := json.Marshal(tx0)
	hash := crypto.HashBytes(b)
	sig, _ := crypto.Sign(string(b), hex.EncodeToString(pv))
	tx0.sig = sig
	tx0.txHash = hash

	tx1 := Transaction{
		address: crypto.HashString("hello2"),
		data: map[string]string{
			"encrypted_address": hex.EncodeToString([]byte("addr")),
			"encrypted_wallet":  hex.EncodeToString([]byte("wallet")),
		},
		txType:    KeychainType,
		timestamp: time.Now(),
		pubKey:    hex.EncodeToString(pub),
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

	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pub, _ := x509.MarshalPKIXPublicKey(key.Public())
	pv, _ := x509.MarshalECPrivateKey(key)

	sk, _ := shared.NewKeyPair(hex.EncodeToString([]byte("encpvkey")), hex.EncodeToString(pub))
	prop, _ := NewProposal(sk)

	tx1 := Transaction{
		address: crypto.HashString("addr"),
		data: map[string]string{
			"encrypted_address": hex.EncodeToString([]byte("addr")),
			"encrypted_wallet":  hex.EncodeToString([]byte("wallet")),
		},
		txType:    KeychainType,
		timestamp: time.Now(),
		pubKey:    hex.EncodeToString(pub),
		prop:      prop,
	}
	txBytesBeforeSig, _ := tx1.MarshalBeforeSignature()
	sig, _ := crypto.Sign(string(txBytesBeforeSig), hex.EncodeToString(pv))
	tx1.emSig = sig
	tx1.sig = sig
	txBytes1, _ := tx1.MarshalHash()
	txHash1 := crypto.HashBytes(txBytes1)
	tx1.txHash = txHash1

	time.Sleep(1 * time.Second)

	tx2 := Transaction{
		address: crypto.HashString("addr"),
		data: map[string]string{
			"encrypted_address": hex.EncodeToString([]byte("addr")),
			"encrypted_wallet":  hex.EncodeToString([]byte("wallet")),
		},
		txType:    KeychainType,
		timestamp: time.Now(),
		pubKey:    hex.EncodeToString(pub),
		prop:      prop,
	}
	txBytesBeforeSig2, _ := tx2.MarshalBeforeSignature()
	sig2, _ := crypto.Sign(string(txBytesBeforeSig2), hex.EncodeToString(pv))
	tx2.emSig = sig2
	tx2.sig = sig2
	txBytes2, _ := tx2.MarshalHash()
	txHash2 := crypto.HashBytes(txBytes2)
	tx2.txHash = txHash2

	assert.Nil(t, tx2.Chain(&tx1))
	assert.NotNil(t, tx2.PreviousTransaction())
	assert.Equal(t, tx1.txHash, tx2.PreviousTransaction().TransactionHash())
}

/*
Scenario: Chain a transaction with same timestamp for the both
	Given two transaction with the same timestamp
	When I want to chain them
	I got an error
*/
func TestChainTransactionWithInvalidTimestamp(t *testing.T) {

	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pub, _ := x509.MarshalPKIXPublicKey(key.Public())
	pv, _ := x509.MarshalECPrivateKey(key)

	sk, _ := shared.NewKeyPair(hex.EncodeToString([]byte("encpvkey")), hex.EncodeToString(pub))
	prop, _ := NewProposal(sk)

	txTime1 := time.Now()
	raw1, _ := json.Marshal(Transaction{
		address: crypto.HashString("addr"),
		data: map[string]string{
			"encrypted_address": hex.EncodeToString([]byte("addr")),
			"encrypted_wallet":  hex.EncodeToString([]byte("wallet")),
		},
		txType:    KeychainType,
		timestamp: txTime1,
		pubKey:    hex.EncodeToString(pub),
		prop:      prop,
	})
	sig1, _ := crypto.Sign(string(raw1), hex.EncodeToString(pv))
	hash1 := crypto.HashBytes(raw1)

	tx1, _ := New(crypto.HashString("addr"), KeychainType, map[string]string{
		"encrypted_address": hex.EncodeToString([]byte("addr")),
		"encrypted_wallet":  hex.EncodeToString([]byte("wallet")),
	}, txTime1, hex.EncodeToString(pub), sig1, sig1, prop, hash1)

	raw2, _ := json.Marshal(Transaction{
		address: crypto.HashString("addr"),
		data: map[string]string{
			"encrypted_address": hex.EncodeToString([]byte("addr")),
			"encrypted_wallet":  hex.EncodeToString([]byte("wallet")),
		},
		txType:    KeychainType,
		timestamp: txTime1,
		pubKey:    hex.EncodeToString(pub),
		prop:      prop,
	})
	sig2, _ := crypto.Sign(string(raw2), hex.EncodeToString(pv))
	hash2 := crypto.HashBytes(raw2)

	tx2, _ := New(crypto.HashString("addr"), KeychainType, map[string]string{
		"encrypted_address": hex.EncodeToString([]byte("addr")),
		"encrypted_wallet":  hex.EncodeToString([]byte("wallet")),
	}, txTime1, hex.EncodeToString(pub), sig2, sig2, prop, hash2)

	assert.EqualError(t, tx2.Chain(&tx1), "previous chained transaction must be anterior to the current transaction")
}

/*
Scenario: Check master validation
	Given a transaction with a master validation
	When I want to check the master validation and the POW
	Then I get not error
*/
func TestCheckMasterValidation(t *testing.T) {

	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pub, _ := x509.MarshalPKIXPublicKey(key.Public())
	pv, _ := x509.MarshalECPrivateKey(key)

	v := MinerValidation{
		minerPubk: hex.EncodeToString(pub),
		status:    ValidationOK,
		timestamp: time.Now(),
	}
	b, _ := json.Marshal(v)
	sig, _ := crypto.Sign(string(b), hex.EncodeToString(pv))
	v.minerSig = sig

	sk, _ := shared.NewKeyPair(hex.EncodeToString([]byte("encpvkey")), hex.EncodeToString(pub))
	prop, _ := NewProposal(sk)

	tx := Transaction{
		address: crypto.HashString("addr"),
		data: map[string]string{
			"encrypted_address": hex.EncodeToString([]byte("addr")),
			"encrypted_wallet":  hex.EncodeToString([]byte("wallet")),
		},
		txType:    KeychainType,
		timestamp: time.Now(),
		pubKey:    hex.EncodeToString(pub),
		prop:      prop,
	}
	txBytesBeforeSig, _ := tx.MarshalBeforeSignature()
	sig, _ = crypto.Sign(string(txBytesBeforeSig), hex.EncodeToString(pv))
	tx.emSig = sig
	tx.sig = sig
	txBytes, _ := tx.MarshalHash()
	txHash := crypto.HashBytes(txBytes)
	tx.txHash = txHash
	tx.masterV = MasterValidation{
		pow:        hex.EncodeToString(pub),
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

	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pub, _ := x509.MarshalPKIXPublicKey(key.Public())
	pv, _ := x509.MarshalECPrivateKey(key)

	key2, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pub2, _ := x509.MarshalPKIXPublicKey(key2.Public())

	v := MinerValidation{
		minerPubk: hex.EncodeToString(pub),
		status:    ValidationOK,
		timestamp: time.Now(),
	}
	b, _ := json.Marshal(v)
	sig, _ := crypto.Sign(string(b), hex.EncodeToString(pv))
	v.minerSig = sig

	sk, _ := shared.NewKeyPair(hex.EncodeToString([]byte("encpvkey")), hex.EncodeToString(pub))
	prop, _ := NewProposal(sk)

	raw, _ := json.Marshal(Transaction{
		address: crypto.HashString("addr"),
		data: map[string]string{
			"encrypted_address": hex.EncodeToString([]byte("addr")),
			"encrypted_wallet":  hex.EncodeToString([]byte("wallet")),
		},
		txType:    KeychainType,
		timestamp: time.Now(),
		pubKey:    hex.EncodeToString(pub),
		prop:      prop,
	})
	sigTx, _ := crypto.Sign(string(raw), hex.EncodeToString(pv))
	hash := crypto.HashBytes(raw)

	tx := Transaction{
		masterV: MasterValidation{
			pow:        hex.EncodeToString(pub2),
			validation: v,
		},
		address: crypto.HashString("addr"),
		data: map[string]string{
			"encrypted_address": hex.EncodeToString([]byte("addr")),
			"encrypted_wallet":  hex.EncodeToString([]byte("wallet")),
		},
		txType:    KeychainType,
		timestamp: time.Now(),
		pubKey:    hex.EncodeToString(pub),
		prop:      prop,
		emSig:     sigTx,
		sig:       sigTx,
		txHash:    hash,
	}

	assert.EqualError(t, tx.CheckMasterValidation(), "invalid proof of work")
}

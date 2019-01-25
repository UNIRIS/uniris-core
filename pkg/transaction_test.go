package uniris

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

	"github.com/stretchr/testify/assert"
	"github.com/uniris/uniris-core/pkg/crypto"
)

/*
Scenario: Create a new transaction proposal
	Given a shared key pair
	When I want to create a transaction proposal
	Then I get a proposal and I can retrieve the shared keys
*/
func TestNewTransactionProposal(t *testing.T) {
	pvKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	key, _ := x509.MarshalPKIXPublicKey(pvKey.Public())

	kp, _ := NewSharedKeyPair(hex.EncodeToString([]byte("encPvKey")), hex.EncodeToString(key))
	prop, err := NewTransactionProposal(kp)
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
func TestNewEmptyTransactionProposal(t *testing.T) {
	_, err := NewTransactionProposal(SharedKeys{})
	assert.Error(t, err, "Transaction proposal: missing shared keys")
}

/*
Scenario: Marshal into a JSON a transaction proposal
	Given a transaction propoal
	When I want to marshal it into a JSON
	Then I get a valid JSON
*/
func TestMarshalTransactionProposal(t *testing.T) {
	pvKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	key, _ := x509.MarshalPKIXPublicKey(pvKey.Public())

	kp, _ := NewSharedKeyPair(hex.EncodeToString([]byte("encPvKey")), hex.EncodeToString(key))
	prop, _ := NewTransactionProposal(kp)
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
func TestNewTransaction(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pub, _ := x509.MarshalPKIXPublicKey(key.Public())
	pv, _ := x509.MarshalECPrivateKey(key)

	kp, _ := NewSharedKeyPair(hex.EncodeToString([]byte("encPvKey")), hex.EncodeToString(pub))
	prop, _ := NewTransactionProposal(kp)

	addr := crypto.HashString("address")
	data := hex.EncodeToString([]byte("data"))
	hash := crypto.HashString("hash")

	sig, _ := crypto.Sign("data", hex.EncodeToString(pv))

	tx, err := NewTransaction(addr, KeychainTransactionType, data, time.Now(), hex.EncodeToString(pub), sig, sig, prop, hash)
	assert.Nil(t, err)
	assert.Equal(t, addr, tx.Address())
	assert.Equal(t, data, tx.Data())
	assert.Equal(t, KeychainTransactionType, tx.Type())
	assert.Equal(t, sig, tx.Signature())
	assert.Equal(t, hex.EncodeToString(pub), tx.Proposal().SharedEmitterKeyPair().PublicKey())
}

/*
Scenario: Create a new transaction with an invalid address
	Given a invalid address hash, empty or not in he
	When I want to create the transaction
	Then I get an error
*/
func TestNewTransactionWithInvalidAddress(t *testing.T) {
	_, err := NewTransaction("", KeychainTransactionType, "", time.Now(), "", "", "", TransactionProposal{}, "")
	assert.EqualError(t, err, "Transaction: Hash is empty")

	_, err = NewTransaction("abc", KeychainTransactionType, "", time.Now(), "", "", "", TransactionProposal{}, "")
	assert.EqualError(t, err, "Transaction: Hash is not in hexadecimal format")

	_, err = NewTransaction(hex.EncodeToString([]byte("abc")), KeychainTransactionType, "", time.Now(), "", "", "", TransactionProposal{}, "")
	assert.EqualError(t, err, "Transaction: Hash is not valid")
}

/*
Scenario: Create a new transaction with an invalid transaction hash
	Given a invalid transaction hash, empty or not in he
	When I want to create the transaction
	Then I get an error
*/
func TestNewTransactionWithInvalidHash(t *testing.T) {

	_, err := NewTransaction(crypto.HashString("hello"), KeychainTransactionType, "", time.Now(), "", "", "", TransactionProposal{}, "")
	assert.EqualError(t, err, "Transaction: Hash is empty")

	_, err = NewTransaction(crypto.HashString("hello"), KeychainTransactionType, "", time.Now(), "", "", "", TransactionProposal{}, "abc")
	assert.EqualError(t, err, "Transaction: Hash is not in hexadecimal format")

	_, err = NewTransaction(crypto.HashString("hello"), KeychainTransactionType, "", time.Now(), "", "", "", TransactionProposal{}, hex.EncodeToString([]byte("abc")))
	assert.EqualError(t, err, "Transaction: Hash is not valid")
}

/*
Scenario: Create a new transaction with an invalid transaction data
	Given a empty or not in data
	When I want to create the transaction
	Then I get an error
*/
func TestNewTransactionWithInvalidData(t *testing.T) {

	_, err := NewTransaction(crypto.HashString("hello"), KeychainTransactionType, "", time.Now(), "", "", "", TransactionProposal{}, crypto.HashString("hello"))
	assert.EqualError(t, err, "Transaction: data is empty")

	_, err = NewTransaction(crypto.HashString("hello"), KeychainTransactionType, "abc", time.Now(), "", "", "", TransactionProposal{}, crypto.HashString("hello"))
	assert.EqualError(t, err, "Transaction: data is not in hexadecimal format")
}

/*
Scenario: Create a new transaction with an invalid transaction timestamp (more than the current timestamp)
	Given a transaction timestamp (now + 2 seconds)
	When I want to create the transaction
	Then I get an error
*/
func TestNewTransactionWithInvalidTimestamp(t *testing.T) {
	_, err := NewTransaction(crypto.HashString("hello"), KeychainTransactionType, hex.EncodeToString([]byte("abc")), time.Now().Add(2*time.Second), "", "", "", TransactionProposal{}, crypto.HashString("hello"))
	assert.EqualError(t, err, "Transaction: timestamp must be greater lower than now")
}

/*
Scenario: Create a new transaction with an invalid transaction public key
	Given an invalid public key: empty or not hex or not a key
	When I want to create the transaction
	Then I get an error
*/
func TestNewTransactionWithInvalidPublicKey(t *testing.T) {
	_, err := NewTransaction(crypto.HashString("hello"), KeychainTransactionType, hex.EncodeToString([]byte("abc")), time.Now(), "", "", "", TransactionProposal{}, crypto.HashString("hello"))
	assert.EqualError(t, err, "Transaction: Public key is empty")

	_, err = NewTransaction(crypto.HashString("hello"), KeychainTransactionType, hex.EncodeToString([]byte("abc")), time.Now(), "abc", "", "", TransactionProposal{}, crypto.HashString("hello"))
	assert.EqualError(t, err, "Transaction: Public key is not in hexadecimal format")

	_, err = NewTransaction(crypto.HashString("hello"), KeychainTransactionType, hex.EncodeToString([]byte("abc")), time.Now(), hex.EncodeToString([]byte("abc")), "", "", TransactionProposal{}, crypto.HashString("hello"))
	assert.EqualError(t, err, "Transaction: Public key is not valid")
}

/*
Scenario: Create a new transaction with an invalid transaction signature
	Given an invalid signature: empty or not hex or not a signature
	When I want to create the transaction
	Then I get an error
*/
func TestNewTransactionWithInvalidSignature(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pub, _ := x509.MarshalPKIXPublicKey(key.Public())

	_, err := NewTransaction(crypto.HashString("hello"), KeychainTransactionType, hex.EncodeToString([]byte("abc")), time.Now(), hex.EncodeToString(pub), "", "", TransactionProposal{}, crypto.HashString(("hello")))
	assert.EqualError(t, err, "Transaction: Signature is empty")

	_, err = NewTransaction(crypto.HashString("hello"), KeychainTransactionType, hex.EncodeToString([]byte("abc")), time.Now(), hex.EncodeToString(pub), "abc", "", TransactionProposal{}, crypto.HashString("hello"))
	assert.EqualError(t, err, "Transaction: Signature is not in hexadecimal format")

	_, err = NewTransaction(crypto.HashString("hello"), KeychainTransactionType, hex.EncodeToString([]byte("abc")), time.Now(), hex.EncodeToString(pub), hex.EncodeToString([]byte("abc")), "", TransactionProposal{}, crypto.HashString("hello"))
	assert.EqualError(t, err, "Transaction: Signature is not valid")
}

/*
Scenario: Create a new transaction with invalid transaction type
	Given an invalid type
	When I want to create the transaction
	Then I get an error
*/
func TestNewTransactionWithInvalidType(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pub, _ := x509.MarshalPKIXPublicKey(key.Public())
	pv, _ := x509.MarshalECPrivateKey(key)

	sig, _ := crypto.Sign("sig", hex.EncodeToString(pv))

	_, err := NewTransaction(crypto.HashString("hello"), 10, hex.EncodeToString([]byte("abc")), time.Now(), hex.EncodeToString(pub), sig, sig, TransactionProposal{}, crypto.HashString(("hello")))
	assert.EqualError(t, err, "Transaction: type not allowed")
}

/*
Scenario: Create a new transaction without proposal
	Given a transaction without proposal
	When I want to create the transaction
	Then I get an error
*/
func TestNewTransactionWithoutProposal(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pub, _ := x509.MarshalPKIXPublicKey(key.Public())
	pv, _ := x509.MarshalECPrivateKey(key)

	sig, _ := crypto.Sign("sig", hex.EncodeToString(pv))

	_, err := NewTransaction(crypto.HashString("hello"), KeychainTransactionType, hex.EncodeToString([]byte("abc")), time.Now(), hex.EncodeToString(pub), sig, sig, TransactionProposal{}, crypto.HashString(("hello")))
	assert.EqualError(t, err, "Transaction: proposal is missing")
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

	sk, _ := NewSharedKeyPair(hex.EncodeToString([]byte("encpvkey")), hex.EncodeToString(pub))
	prop, _ := NewTransactionProposal(sk)

	raw, _ := json.Marshal(struct {
		Address   string
		Data      string
		Type      TransactionType
		PublicKey string
		Proposal  TransactionProposal
	}{
		Address:   crypto.HashString("hello"),
		Data:      hex.EncodeToString([]byte("abc")),
		Type:      KeychainTransactionType,
		PublicKey: hex.EncodeToString(pub),
		Proposal:  prop,
	})
	sig, _ := crypto.Sign(string(raw), hex.EncodeToString(pv))
	hash := crypto.HashBytes(raw)

	tx, _ := NewTransaction(crypto.HashString("hello"), KeychainTransactionType, hex.EncodeToString([]byte("abc")), time.Now(), hex.EncodeToString(pub), sig, sig, prop, hash)
	assert.Nil(t, tx.CheckTransactionIntegrity())
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

	sk, _ := NewSharedKeyPair(hex.EncodeToString([]byte("encpvkey")), hex.EncodeToString(pub))
	prop, _ := NewTransactionProposal(sk)

	raw, _ := json.Marshal(struct {
		Address   string
		Data      string
		Type      TransactionType
		PublicKey string
		Proposal  TransactionProposal
	}{
		Address:   crypto.HashString("hello"),
		Data:      hex.EncodeToString([]byte("abc")),
		Type:      KeychainTransactionType,
		PublicKey: hex.EncodeToString(pub),
		Proposal:  prop,
	})
	sig, _ := crypto.Sign(string(raw), hex.EncodeToString(pv))
	hash := "abc"

	tx, _ := NewTransaction(crypto.HashString("hello"), KeychainTransactionType, hex.EncodeToString([]byte("abc")), time.Now(), hex.EncodeToString(pub), sig, sig, prop, hash)
	assert.EqualError(t, tx.CheckTransactionIntegrity(), "Transaction integrity violated")
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

	sk, _ := NewSharedKeyPair(hex.EncodeToString([]byte("encpvkey")), hex.EncodeToString(pub))
	prop, _ := NewTransactionProposal(sk)

	raw, _ := json.Marshal(struct {
		Address   string              `json:"address"`
		Data      string              `json:"data"`
		Type      TransactionType     `json:"type"`
		PublicKey string              `json:"public_key"`
		Proposal  TransactionProposal `json:"proposal"`
	}{
		Address:   crypto.HashString("hello"),
		Data:      hex.EncodeToString([]byte("abc")),
		Type:      KeychainTransactionType,
		PublicKey: hex.EncodeToString(pub),
		Proposal:  prop,
	})
	sig, _ := crypto.Sign(string("fake sig"), hex.EncodeToString(pv))
	hash := crypto.HashBytes(raw)

	tx, _ := NewTransaction(crypto.HashString("hello"), KeychainTransactionType, hex.EncodeToString([]byte("abc")), time.Now(), hex.EncodeToString(pub), sig, sig, prop, hash)
	assert.EqualError(t, tx.CheckTransactionIntegrity(), "Transaction signature invalid")
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

	sk, _ := NewSharedKeyPair(hex.EncodeToString([]byte("encpvkey")), hex.EncodeToString(pub))
	prop, _ := NewTransactionProposal(sk)

	raw, _ := json.Marshal(struct {
		Address   string
		Data      string
		Type      TransactionType
		PublicKey string
		Proposal  TransactionProposal
	}{
		Address:   crypto.HashString("hello"),
		Data:      hex.EncodeToString([]byte("abc")),
		Type:      KeychainTransactionType,
		PublicKey: hex.EncodeToString(pub),
		Proposal:  prop,
	})
	sig, _ := crypto.Sign(string(raw), hex.EncodeToString(pv))
	hash := crypto.HashBytes(raw)
	tx, _ := NewTransaction(crypto.HashString("hello"), KeychainTransactionType, hex.EncodeToString([]byte("abc")), time.Now(), hex.EncodeToString(pub), sig, sig, prop, hash)

	v, _ := NewMinerValidation(ValidationOK, time.Now(), hex.EncodeToString(pub), "")
	b, _ := json.Marshal(v)
	sig, _ = crypto.Sign(string(b), hex.EncodeToString(pv))
	v, _ = NewMinerValidation(ValidationOK, time.Now(), hex.EncodeToString(pub), sig)

	masterValid, _ := NewMasterValidation([]PeerIdentity{}, hex.EncodeToString(pub), v)
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

	sk, _ := NewSharedKeyPair(hex.EncodeToString([]byte("encpvkey")), hex.EncodeToString(pub))
	prop, _ := NewTransactionProposal(sk)

	raw, _ := json.Marshal(struct {
		Address   string
		Data      string
		Type      TransactionType
		PublicKey string
		Proposal  TransactionProposal
	}{
		Address:   crypto.HashString("hello"),
		Data:      hex.EncodeToString([]byte("abc")),
		Type:      KeychainTransactionType,
		PublicKey: hex.EncodeToString(pub),
		Proposal:  prop,
	})
	sig, _ := crypto.Sign(string(raw), hex.EncodeToString(pv))
	hash := crypto.HashBytes(raw)
	tx, _ := NewTransaction(crypto.HashString("hello"), KeychainTransactionType, hex.EncodeToString([]byte("abc")), time.Now(), hex.EncodeToString(pub), sig, sig, prop, hash)

	v, _ := NewMinerValidation(ValidationOK, time.Now(), hex.EncodeToString(pub), "")
	b, _ := json.Marshal(v)
	sig, _ = crypto.Sign(string(b), hex.EncodeToString(pv))
	v, _ = NewMinerValidation(ValidationOK, time.Now(), hex.EncodeToString(pub), sig)

	masterValid, _ := NewMasterValidation([]PeerIdentity{}, hex.EncodeToString(pub), v)
	assert.EqualError(t, tx.AddMining(masterValid, []MinerValidation{}), "Transaction: Missing confirmation validations")
}

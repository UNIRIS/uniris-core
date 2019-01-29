package transaction

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/uniris/uniris-core/pkg/crypto"
	"github.com/uniris/uniris-core/pkg/shared"
)

/*
Scenario: Create a new ID transaction
	Given an transaction with ID type
	When I want to format it to an ID transaction
	Then I get it with extract of the data fields
*/
func TestNewID(t *testing.T) {

	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pub, _ := x509.MarshalPKIXPublicKey(key.Public())
	pv, _ := x509.MarshalECPrivateKey(key)

	kp, _ := shared.NewKeyPair(hex.EncodeToString([]byte("encPvKey")), hex.EncodeToString(pub))
	prop, _ := NewProposal(kp)

	addr := crypto.HashString("address")

	dataBytes, _ := json.Marshal(idData{
		EncryptedAESKey:         hex.EncodeToString([]byte("aesKey")),
		EncryptedAddressByID:    hex.EncodeToString([]byte("addr")),
		EncryptedAddressByRobot: hex.EncodeToString([]byte("addr")),
	})

	data := hex.EncodeToString(dataBytes)
	hash := crypto.HashString("hash")

	sig, _ := crypto.Sign("data", hex.EncodeToString(pv))

	tx, err := New(addr, IDType, data, time.Now(), hex.EncodeToString(pub), sig, sig, prop, hash)
	assert.Nil(t, err)

	id, err := NewID(tx)
	assert.Nil(t, err)

	assert.Equal(t, hex.EncodeToString([]byte("aesKey")), id.EncryptedAESKey())
	assert.Equal(t, hex.EncodeToString([]byte("addr")), id.EncryptedAddrByRobot())
	assert.Equal(t, hex.EncodeToString([]byte("addr")), id.EncryptedAddrByID())

}

/*
Scenario: Create a new ID transaction with another type of transaction
	Given an transaction with Keychain type
	When I want to format it to an ID transaction
	Then I get an error
*/
func TestNewIDWithInvalidType(t *testing.T) {

	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pub, _ := x509.MarshalPKIXPublicKey(key.Public())
	pv, _ := x509.MarshalECPrivateKey(key)

	kp, _ := shared.NewKeyPair(hex.EncodeToString([]byte("encPvKey")), hex.EncodeToString(pub))
	prop, _ := NewProposal(kp)

	addr := crypto.HashString("address")

	dataBytes, _ := json.Marshal(idData{
		EncryptedAESKey:         hex.EncodeToString([]byte("aesKey")),
		EncryptedAddressByID:    hex.EncodeToString([]byte("addr")),
		EncryptedAddressByRobot: hex.EncodeToString([]byte("addr")),
	})

	data := hex.EncodeToString(dataBytes)
	hash := crypto.HashString("hash")

	sig, _ := crypto.Sign("data", hex.EncodeToString(pv))

	tx, err := New(addr, KeychainType, data, time.Now(), hex.EncodeToString(pub), sig, sig, prop, hash)
	assert.Nil(t, err)

	_, err = NewID(tx)
	assert.EqualError(t, err, "transaction: invalid type of transaction")

}

/*
Scenario: Create a new ID transaction with missing data fields
	Given an transaction with ID type and missing data fields
	When I want to format it to an ID transaction
	Then I get an error
*/
func TestNewIDWithMissingDataFields(t *testing.T) {

	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pub, _ := x509.MarshalPKIXPublicKey(key.Public())
	pv, _ := x509.MarshalECPrivateKey(key)

	kp, _ := shared.NewKeyPair(hex.EncodeToString([]byte("encPvKey")), hex.EncodeToString(pub))
	prop, _ := NewProposal(kp)

	addr := crypto.HashString("address")

	dataBytes, _ := json.Marshal(idData{
		EncryptedAddressByID:    hex.EncodeToString([]byte("addr")),
		EncryptedAddressByRobot: hex.EncodeToString([]byte("addr")),
	})

	data := hex.EncodeToString(dataBytes)
	hash := crypto.HashString("hash")

	sig, _ := crypto.Sign("data", hex.EncodeToString(pv))

	tx, err := New(addr, IDType, data, time.Now(), hex.EncodeToString(pub), sig, sig, prop, hash)
	assert.Nil(t, err)

	_, err = NewID(tx)
	assert.EqualError(t, err, "transaction: missing id transaction data")
}

/*
Scenario: Create a new ID transaction with data fields not in hex
	Given an transaction with ID type and data fields with non hexadecimal
	When I want to format it to an ID transaction
	Then I get an error
*/
func TestNewIDWithNotHexDataFields(t *testing.T) {

	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pub, _ := x509.MarshalPKIXPublicKey(key.Public())
	pv, _ := x509.MarshalECPrivateKey(key)

	kp, _ := shared.NewKeyPair(hex.EncodeToString([]byte("encPvKey")), hex.EncodeToString(pub))
	prop, _ := NewProposal(kp)

	addr := crypto.HashString("address")

	dataBytes1, _ := json.Marshal(idData{
		EncryptedAESKey:         "aesKey",
		EncryptedAddressByID:    "addr",
		EncryptedAddressByRobot: "addr",
	})

	data1 := hex.EncodeToString(dataBytes1)
	hash := crypto.HashString("hash")

	sig, _ := crypto.Sign("data", hex.EncodeToString(pv))

	tx, _ := New(addr, IDType, data1, time.Now(), hex.EncodeToString(pub), sig, sig, prop, hash)
	_, err := NewID(tx)
	assert.EqualError(t, err, "transaction: id encrypted aes key is not in hexadecimal format")

	dataBytes2, _ := json.Marshal(idData{
		EncryptedAESKey:         hex.EncodeToString([]byte("aesKey")),
		EncryptedAddressByID:    "addr",
		EncryptedAddressByRobot: "addr",
	})
	data2 := hex.EncodeToString(dataBytes2)
	tx, _ = New(addr, IDType, data2, time.Now(), hex.EncodeToString(pub), sig, sig, prop, hash)
	_, err = NewID(tx)
	assert.EqualError(t, err, "transaction: id encrypted address for id is not in hexadecimal format")

	dataBytes3, _ := json.Marshal(idData{
		EncryptedAESKey:         hex.EncodeToString([]byte("aesKey")),
		EncryptedAddressByID:    hex.EncodeToString([]byte("addr")),
		EncryptedAddressByRobot: "addr",
	})
	data3 := hex.EncodeToString(dataBytes3)
	tx, _ = New(addr, IDType, data3, time.Now(), hex.EncodeToString(pub), sig, sig, prop, hash)
	_, err = NewID(tx)
	assert.EqualError(t, err, "transaction: id encrypted address for robot is not in hexadecimal format")
}

/*
Scenario: Convert back a ID to its parent Transaction
	Given an ID transaction struct
	When I want to convert back to its parent
	Then I get a transaction struct
*/
func TestIDToTransaction(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pub, _ := x509.MarshalPKIXPublicKey(key.Public())
	pv, _ := x509.MarshalECPrivateKey(key)

	kp, _ := shared.NewKeyPair(hex.EncodeToString([]byte("encPvKey")), hex.EncodeToString(pub))
	prop, _ := NewProposal(kp)

	addr := crypto.HashString("address")

	dataBytes, _ := json.Marshal(idData{
		EncryptedAESKey:         hex.EncodeToString([]byte("aesKey")),
		EncryptedAddressByID:    hex.EncodeToString([]byte("addr")),
		EncryptedAddressByRobot: hex.EncodeToString([]byte("addr")),
	})

	data := hex.EncodeToString(dataBytes)
	hash := crypto.HashString("hash")

	sig, _ := crypto.Sign("data", hex.EncodeToString(pv))

	tx, err := New(addr, IDType, data, time.Now(), hex.EncodeToString(pub), sig, sig, prop, hash)
	assert.Nil(t, err)

	id, err := NewID(tx)
	assert.Nil(t, err)

	tx, err = id.ToTransaction()
	assert.Nil(t, err)
	assert.Equal(t, data, tx.Data())

	b, _ := json.Marshal(MinerValidation{
		minerPubk: hex.EncodeToString(pub),
		status:    ValidationOK,
		timestamp: time.Now(),
	})
	sig, _ = crypto.Sign(string(b), hex.EncodeToString(pv))
	v, _ := NewMinerValidation(ValidationOK, time.Now(), hex.EncodeToString(pub), sig)

	masterValid, _ := NewMasterValidation(Pool{}, hex.EncodeToString(pub), v)

	tx.AddMining(masterValid, []MinerValidation{v})
	assert.Equal(t, ValidationOK, tx.MasterValidation().Validation().Status())
	assert.Len(t, tx.ConfirmationsValidations(), 1)
}

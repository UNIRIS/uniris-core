package transaction

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
Scenario: Create a new Keychain transaction
	Given an transaction with Keychain type
	When I want to format it to an Keychain transaction
	Then I get it with extract of the data fields
*/
func TestNewKeychain(t *testing.T) {

	pub, pv := crypto.GenerateKeys()

	kp, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("encPvKey")), pub)
	prop, _ := NewProposal(kp)

	addr := crypto.HashString("address")

	hash := crypto.HashString("hash")
	sig, _ := crypto.Sign("data", pv)

	tx, err := New(addr, KeychainType, map[string]string{
		"encrypted_address": hex.EncodeToString([]byte("addr")),
		"encrypted_wallet":  hex.EncodeToString([]byte("wallet")),
	}, time.Now(), pub, sig, sig, prop, hash)
	assert.Nil(t, err)

	keychain, err := NewKeychain(tx)
	assert.Nil(t, err)

	assert.Equal(t, hex.EncodeToString([]byte("addr")), keychain.EncryptedAddrByRobot())
	assert.Equal(t, hex.EncodeToString([]byte("wallet")), keychain.EncryptedWallet())

}

/*
Scenario: Create a new Keychain transaction with another type of transaction
	Given an transaction with Keychain type
	When I want to format it to an Keychain transaction
	Then I get an error
*/
func TestNewKeychainWithInvalkeychainType(t *testing.T) {

	pub, pv := crypto.GenerateKeys()

	kp, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("encPvKey")), pub)
	prop, _ := NewProposal(kp)

	addr := crypto.HashString("address")

	hash := crypto.HashString("hash")

	sig, _ := crypto.Sign("data", pv)

	tx, err := New(addr, IDType, map[string]string{
		"encrypted_address": hex.EncodeToString([]byte("addr")),
		"encrypted_wallet":  hex.EncodeToString([]byte("wallet")),
	}, time.Now(), pub, sig, sig, prop, hash)
	assert.Nil(t, err)

	_, err = NewKeychain(tx)
	assert.EqualError(t, err, "transaction: invalid type of transaction")

}

/*
Scenario: Create a new Keychain transaction with missing data fields
	Given an transaction with Keychain type and missing data fields
	When I want to format it to an Keychain transaction
	Then I get an error
*/
func TestNewKeychainWithMissingDataFields(t *testing.T) {

	pub, pv := crypto.GenerateKeys()

	kp, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("encPvKey")), pub)
	prop, _ := NewProposal(kp)

	addr := crypto.HashString("address")

	hash := crypto.HashString("hash")

	sig, _ := crypto.Sign("data", pv)

	tx, err := New(addr, KeychainType, map[string]string{
		"encrypted_wallet": hex.EncodeToString([]byte("wallet")),
	}, time.Now(), pub, sig, sig, prop, hash)
	assert.Nil(t, err)

	_, err = NewKeychain(tx)
	assert.EqualError(t, err, "transaction: missing data keychain: 'encrypted_address'")

	tx, err = New(addr, KeychainType, map[string]string{
		"encrypted_address": hex.EncodeToString([]byte("wallet")),
	}, time.Now(), pub, sig, sig, prop, hash)
	assert.Nil(t, err)

	_, err = NewKeychain(tx)
	assert.EqualError(t, err, "transaction: missing data keychain: 'encrypted_wallet'")
}

/*
Scenario: Create a new Keychain transaction with data fields not in hex
	Given an transaction with Keychain type and data fields with non hexadecimal
	When I want to format it to an Keychain transaction
	Then I get an error
*/
func TestNewKeychainWithNotHexDataFields(t *testing.T) {

	pub, pv := crypto.GenerateKeys()

	kp, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("encPvKey")), pub)
	prop, _ := NewProposal(kp)

	addr := crypto.HashString("address")

	hash := crypto.HashString("hash")

	sig, _ := crypto.Sign("data", pv)

	tx, _ := New(addr, KeychainType, map[string]string{
		"encrypted_address": "addr",
		"encrypted_wallet":  hex.EncodeToString([]byte("wallet")),
	}, time.Now(), pub, sig, sig, prop, hash)
	_, err := NewKeychain(tx)
	assert.EqualError(t, err, "transaction: keychain encrypted address is not in hexadecimal format")

	tx, _ = New(addr, KeychainType, map[string]string{
		"encrypted_address": hex.EncodeToString([]byte("addr")),
		"encrypted_wallet":  "wallet",
	}, time.Now(), pub, sig, sig, prop, hash)
	_, err = NewKeychain(tx)
	assert.EqualError(t, err, "transaction: keychain encrypted wallet is not in hexadecimal format")
}

/*
Scenario: Convert back a Keychain to its parent Transaction
	Given an Keychain transaction struct
	When I want to convert back to its parent
	Then I get a transaction struct
*/
func TestKeychainToTransaction(t *testing.T) {
	pub, pv := crypto.GenerateKeys()

	kp, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("encPvKey")), pub)
	prop, _ := NewProposal(kp)

	addr := crypto.HashString("address")

	hash := crypto.HashString("hash")

	sig, _ := crypto.Sign("data", pv)

	tx, err := New(addr, KeychainType, map[string]string{
		"encrypted_address": hex.EncodeToString([]byte("addr")),
		"encrypted_wallet":  hex.EncodeToString([]byte("wallet")),
	}, time.Now(), pub, sig, sig, prop, hash)
	assert.Nil(t, err)

	keychain, err := NewKeychain(tx)
	assert.Nil(t, err)

	tx, err = keychain.ToTransaction()
	assert.Nil(t, err)

	assert.Equal(t, map[string]string{
		"encrypted_address": hex.EncodeToString([]byte("addr")),
		"encrypted_wallet":  hex.EncodeToString([]byte("wallet")),
	}, tx.Data())

	b, _ := json.Marshal(MinerValidation{
		minerPubk: pub,
		status:    ValidationOK,
		timestamp: time.Now(),
	})
	sig, _ = crypto.Sign(string(b), pv)
	v, _ := NewMinerValidation(ValidationOK, time.Now(), pub, sig)

	masterValkeychain, _ := NewMasterValidation(Pool{}, pub, v)

	tx.AddMining(masterValkeychain, []MinerValidation{v})
	assert.Equal(t, ValidationOK, tx.MasterValidation().Validation().Status())
	assert.Len(t, tx.ConfirmationsValidations(), 1)
}

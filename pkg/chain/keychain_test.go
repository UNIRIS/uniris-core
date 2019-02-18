package chain

import (
	"encoding/hex"
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

	prop, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("encPvKey")), pub)

	addr := crypto.HashString("address")

	hash := crypto.HashString("hash")
	sig, _ := crypto.Sign("data", pv)

	tx, err := NewTransaction(addr, KeychainTransactionType, map[string]string{
		"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
		"encrypted_wallet":           hex.EncodeToString([]byte("wallet")),
	}, time.Now(), pub, prop, sig, sig, hash)
	assert.Nil(t, err)

	keychain, err := NewKeychain(tx)
	assert.Nil(t, err)

	assert.Equal(t, hex.EncodeToString([]byte("addr")), keychain.EncryptedAddrBy())
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

	prop, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("encPvKey")), pub)

	addr := crypto.HashString("address")

	hash := crypto.HashString("hash")

	sig, _ := crypto.Sign("data", pv)

	tx, err := NewTransaction(addr, IDTransactionType, map[string]string{
		"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
		"encrypted_wallet":           hex.EncodeToString([]byte("wallet")),
	}, time.Now(), pub, prop, sig, sig, hash)
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

	prop, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("encPvKey")), pub)

	addr := crypto.HashString("address")

	hash := crypto.HashString("hash")

	sig, _ := crypto.Sign("data", pv)

	tx, err := NewTransaction(addr, KeychainTransactionType, map[string]string{
		"encrypted_wallet": hex.EncodeToString([]byte("wallet")),
	}, time.Now(), pub, prop, sig, sig, hash)
	assert.Nil(t, err)

	_, err = NewKeychain(tx)
	assert.EqualError(t, err, "transaction: missing data keychain: 'encrypted_address_by_node'")

	tx, err = NewTransaction(addr, KeychainTransactionType, map[string]string{
		"encrypted_address_by_node": hex.EncodeToString([]byte("wallet")),
	}, time.Now(), pub, prop, sig, sig, hash)
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

	prop, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("encPvKey")), pub)

	addr := crypto.HashString("address")

	hash := crypto.HashString("hash")

	sig, _ := crypto.Sign("data", pv)

	tx, _ := NewTransaction(addr, KeychainTransactionType, map[string]string{
		"encrypted_address_by_node": "addr",
		"encrypted_wallet":           hex.EncodeToString([]byte("wallet")),
	}, time.Now(), pub, prop, sig, sig, hash)
	_, err := NewKeychain(tx)
	assert.EqualError(t, err, "transaction: keychain encrypted address for node is not in hexadecimal format")

	tx, _ = NewTransaction(addr, KeychainTransactionType, map[string]string{
		"encrypted_address_by_node": hex.EncodeToString([]byte("addr")),
		"encrypted_wallet":           "wallet",
	}, time.Now(), pub, prop, sig, sig, hash)
	_, err = NewKeychain(tx)
	assert.EqualError(t, err, "transaction: keychain encrypted wallet is not in hexadecimal format")
}

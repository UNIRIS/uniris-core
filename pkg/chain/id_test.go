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
Scenario: Create a new ID transaction
	Given an transaction with ID type
	When I want to format it to an ID transaction
	Then I get it with extract of the data fields
*/
func TestNewID(t *testing.T) {

	pub, pv := crypto.GenerateKeys()

	prop, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("encPvKey")), pub)

	addr := crypto.HashString("address")

	hash := crypto.HashString("hash")

	sig, _ := crypto.Sign("data", pv)

	tx, err := NewTransaction(addr, IDTransactionType, map[string]string{
		"encrypted_aes_key":          hex.EncodeToString([]byte("aesKey")),
		"encrypted_address_by_miner": hex.EncodeToString([]byte("addr")),
		"encrypted_address_by_id":    hex.EncodeToString([]byte("addr")),
	}, time.Now(), pub, prop, sig, sig, hash)
	assert.Nil(t, err)

	id, err := NewID(tx)
	assert.Nil(t, err)

	assert.Equal(t, hex.EncodeToString([]byte("aesKey")), id.EncryptedAESKey())
	assert.Equal(t, hex.EncodeToString([]byte("addr")), id.EncryptedAddrByMiner())
	assert.Equal(t, hex.EncodeToString([]byte("addr")), id.EncryptedAddrByID())

}

/*
Scenario: Create a new ID transaction with another type of transaction
	Given an transaction with Keychain type
	When I want to format it to an ID transaction
	Then I get an error
*/
func TestNewIDWithInvalidType(t *testing.T) {

	pub, pv := crypto.GenerateKeys()

	prop, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("encPvKey")), pub)

	addr := crypto.HashString("address")

	hash := crypto.HashString("hash")

	sig, _ := crypto.Sign("data", pv)

	tx, err := NewTransaction(addr, KeychainTransactionType, map[string]string{
		"encrypted_aes_key":          hex.EncodeToString([]byte("aesKey")),
		"encrypted_address_by_miner": hex.EncodeToString([]byte("addr")),
		"encrypted_address_by_id":    hex.EncodeToString([]byte("addr")),
	}, time.Now(), pub, prop, sig, sig, hash)
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

	pub, pv := crypto.GenerateKeys()

	prop, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("encPvKey")), pub)

	addr := crypto.HashString("address")

	hash := crypto.HashString("hash")

	sig, _ := crypto.Sign("data", pv)

	tx, err := NewTransaction(addr, IDTransactionType, map[string]string{
		"encrypted_address_by_id": hex.EncodeToString([]byte("addr")),
		"encrypted_aes_key":       hex.EncodeToString([]byte("aesKey")),
	}, time.Now(), pub, prop, sig, sig, hash)
	assert.Nil(t, err)

	_, err = NewID(tx)
	assert.EqualError(t, err, "transaction: missing data ID 'encrypted_address_by_miner'")

	tx, err = NewTransaction(addr, IDTransactionType, map[string]string{
		"encrypted_address_by_miner": hex.EncodeToString([]byte("addr")),
		"encrypted_aes_key":          hex.EncodeToString([]byte("aesKey")),
	}, time.Now(), pub, prop, sig, sig, hash)
	assert.Nil(t, err)

	_, err = NewID(tx)
	assert.EqualError(t, err, "transaction: missing data ID 'encrypted_address_by_id'")

	tx, err = NewTransaction(addr, IDTransactionType, map[string]string{
		"encrypted_address_by_miner": hex.EncodeToString([]byte("addr")),
		"encrypted_address_by_id":    hex.EncodeToString([]byte("addr")),
	}, time.Now(), pub, prop, sig, sig, hash)
	assert.Nil(t, err)

	_, err = NewID(tx)
	assert.EqualError(t, err, "transaction: missing data ID 'encrypted_aes_key'")
}

/*
Scenario: Create a new ID transaction with data fields not in hex
	Given an transaction with ID type and data fields with non hexadecimal
	When I want to format it to an ID transaction
	Then I get an error
*/
func TestNewIDWithNotHexDataFields(t *testing.T) {

	pub, pv := crypto.GenerateKeys()

	prop, _ := shared.NewEmitterKeyPair(hex.EncodeToString([]byte("encPvKey")), pub)

	addr := crypto.HashString("address")

	hash := crypto.HashString("hash")

	sig, _ := crypto.Sign("data", pv)

	tx, _ := NewTransaction(addr, IDTransactionType, map[string]string{
		"encrypted_aes_key":          "aes key",
		"encrypted_address_by_miner": hex.EncodeToString([]byte("addr")),
		"encrypted_address_by_id":    hex.EncodeToString([]byte("addr")),
	}, time.Now(), pub, prop, sig, sig, hash)
	_, err := NewID(tx)
	assert.EqualError(t, err, "transaction: id encrypted aes key is not in hexadecimal format")

	tx, _ = NewTransaction(addr, IDTransactionType, map[string]string{
		"encrypted_aes_key":          hex.EncodeToString([]byte("aesKey")),
		"encrypted_address_by_id":    "addr",
		"encrypted_address_by_miner": hex.EncodeToString([]byte("addr")),
	}, time.Now(), pub, prop, sig, sig, hash)
	_, err = NewID(tx)
	assert.EqualError(t, err, "transaction: id encrypted address for id is not in hexadecimal format")

	tx, _ = NewTransaction(addr, IDTransactionType, map[string]string{
		"encrypted_aes_key":          hex.EncodeToString([]byte("aesKey")),
		"encrypted_address_by_id":    hex.EncodeToString([]byte("addr")),
		"encrypted_address_by_miner": "addr",
	}, time.Now(), pub, prop, sig, sig, hash)
	_, err = NewID(tx)
	assert.EqualError(t, err, "transaction: id encrypted address for miner is not in hexadecimal format")
}

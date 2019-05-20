package chain

// import (
// 	"crypto/rand"
// 	"testing"
// 	"time"

// 	"github.com/stretchr/testify/assert"

// 	"github.com/uniris/uniris-core/pkg/crypto"
// 	"github.com/uniris/uniris-core/pkg/shared"
// )

// /*
// Scenario: Create a new ID transaction
// 	Given an transaction with ID type
// 	When I want to format it to an ID transaction
// 	Then I get it with extract of the data fields
// */
// func TestNewID(t *testing.T) {

// 	pv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

// 	prop, _ := shared.NewEmitterCrossKeyPair([]byte("encPvKey"), pub)

// 	addr := crypto.Hash([]byte("address"))

// 	hash := crypto.Hash([]byte("hash"))

// 	sig, _ := pv.Sign([]byte("data"))

// 	tx, err := NewTransaction(addr, IDTransactionType, map[string][]byte{
// 		"encrypted_aes_key":         []byte("aesKey"),
// 		"encrypted_address_by_node": []byte("addr"),
// 		"encrypted_address_by_id":   []byte("addr"),
// 	}, time.Now(), pub, prop, sig, sig, hash)
// 	assert.Nil(t, err)

// 	id, err := NewID(tx)
// 	assert.Nil(t, err)

// 	assert.Equal(t, []byte("aesKey"), id.EncryptedAESKey())
// 	assert.Equal(t, []byte("addr"), id.EncryptedAddrBy())
// 	assert.Equal(t, []byte("addr"), id.EncryptedAddrByID())

// }

// /*
// Scenario: Create a new ID transaction with another type of transaction
// 	Given an transaction with Keychain type
// 	When I want to format it to an ID transaction
// 	Then I get an error
// */
// func TestNewIDWithInvalidType(t *testing.T) {

// 	pv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

// 	prop, _ := shared.NewEmitterCrossKeyPair([]byte("encPvKey"), pub)

// 	addr := crypto.Hash([]byte("address"))

// 	hash := crypto.Hash([]byte("hash"))

// 	sig, _ := pv.Sign([]byte("data"))

// 	tx, err := NewTransaction(addr, KeychainTransactionType, map[string][]byte{
// 		"encrypted_aes_key":         []byte("aesKey"),
// 		"encrypted_address_by_node": []byte("addr"),
// 		"encrypted_address_by_id":   []byte("addr"),
// 	}, time.Now(), pub, prop, sig, sig, hash)
// 	assert.Nil(t, err)

// 	_, err = NewID(tx)
// 	assert.EqualError(t, err, "invalid type of transaction")

// }

// /*
// Scenario: Create a new ID transaction with missing data fields
// 	Given an transaction with ID type and missing data fields
// 	When I want to format it to an ID transaction
// 	Then I get an error
// */
// func TestNewIDWithMissingDataFields(t *testing.T) {

// 	pv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

// 	prop, _ := shared.NewEmitterCrossKeyPair([]byte("encPvKey"), pub)

// 	addr := crypto.Hash([]byte("address"))

// 	hash := crypto.Hash([]byte("hash"))

// 	sig, _ := pv.Sign([]byte("data"))

// 	tx, err := NewTransaction(addr, IDTransactionType, map[string][]byte{
// 		"encrypted_address_by_id": []byte("addr"),
// 		"encrypted_aes_key":       []byte("aesKey"),
// 	}, time.Now(), pub, prop, sig, sig, hash)
// 	assert.Nil(t, err)

// 	_, err = NewID(tx)
// 	assert.EqualError(t, err, "missing ID data: 'encrypted_address_by_node'")

// 	tx, err = NewTransaction(addr, IDTransactionType, map[string][]byte{
// 		"encrypted_address_by_node": []byte("addr"),
// 		"encrypted_aes_key":         []byte("aesKey"),
// 	}, time.Now(), pub, prop, sig, sig, hash)
// 	assert.Nil(t, err)

// 	_, err = NewID(tx)
// 	assert.EqualError(t, err, "missing ID data: 'encrypted_address_by_id'")

// 	tx, err = NewTransaction(addr, IDTransactionType, map[string][]byte{
// 		"encrypted_address_by_node": []byte("addr"),
// 		"encrypted_address_by_id":   []byte("addr"),
// 	}, time.Now(), pub, prop, sig, sig, hash)
// 	assert.Nil(t, err)

// 	_, err = NewID(tx)
// 	assert.EqualError(t, err, "missing ID data: 'encrypted_aes_key'")
// }

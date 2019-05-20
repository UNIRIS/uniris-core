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
// Scenario: Create a new Keychain transaction
// 	Given an transaction with Keychain type
// 	When I want to format it to an Keychain transaction
// 	Then I get it with extract of the data fields
// */
// func TestNewKeychain(t *testing.T) {

// 	pv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

// 	prop, _ := shared.NewEmitterCrossKeyPair([]byte("encPvKey"), pub)

// 	addr := crypto.Hash([]byte("address"))

// 	hash := crypto.Hash([]byte("hash"))
// 	sig, _ := pv.Sign([]byte("data"))

// 	tx, err := NewTransaction(addr, KeychainTransactionType, map[string][]byte{
// 		"encrypted_address_by_node": []byte("addr"),
// 		"encrypted_wallet":          []byte("wallet"),
// 	}, time.Now(), pub, prop, sig, sig, hash)
// 	assert.Nil(t, err)

// 	keychain, err := NewKeychain(tx)
// 	assert.Nil(t, err)

// 	assert.Equal(t, []byte("addr"), keychain.EncryptedAddrBy())
// 	assert.Equal(t, []byte("wallet"), keychain.EncryptedWallet())

// }

// /*
// Scenario: Create a new Keychain transaction with another type of transaction
// 	Given an transaction with Keychain type
// 	When I want to format it to an Keychain transaction
// 	Then I get an error
// */
// func TestNewKeychainWithInvalkeychainType(t *testing.T) {

// 	pv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

// 	prop, _ := shared.NewEmitterCrossKeyPair([]byte("encPvKey"), pub)

// 	addr := crypto.Hash([]byte("address"))

// 	hash := crypto.Hash([]byte("hash"))

// 	sig, _ := pv.Sign([]byte("data"))

// 	tx, err := NewTransaction(addr, IDTransactionType, map[string][]byte{
// 		"encrypted_address_by_node": []byte("addr"),
// 		"encrypted_wallet":          []byte("wallet"),
// 	}, time.Now(), pub, prop, sig, sig, hash)
// 	assert.Nil(t, err)

// 	_, err = NewKeychain(tx)
// 	assert.EqualError(t, err, "invalid type of transaction")

// }

// /*
// Scenario: Create a new Keychain transaction with missing data fields
// 	Given an transaction with Keychain type and missing data fields
// 	When I want to format it to an Keychain transaction
// 	Then I get an error
// */
// func TestNewKeychainWithMissingDataFields(t *testing.T) {

// 	pv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

// 	prop, _ := shared.NewEmitterCrossKeyPair([]byte("encPvKey"), pub)

// 	addr := crypto.Hash([]byte("address"))

// 	hash := crypto.Hash([]byte("hash"))

// 	sig, _ := pv.Sign([]byte("data"))

// 	tx, err := NewTransaction(addr, KeychainTransactionType, map[string][]byte{
// 		"encrypted_wallet": []byte("wallet"),
// 	}, time.Now(), pub, prop, sig, sig, hash)
// 	assert.Nil(t, err)

// 	_, err = NewKeychain(tx)
// 	assert.EqualError(t, err, "missing keychain data: 'encrypted_address_by_node'")

// 	tx, err = NewTransaction(addr, KeychainTransactionType, map[string][]byte{
// 		"encrypted_address_by_node": []byte("wallet"),
// 	}, time.Now(), pub, prop, sig, sig, hash)
// 	assert.Nil(t, err)

// 	_, err = NewKeychain(tx)
// 	assert.EqualError(t, err, "missing keychain data: 'encrypted_wallet'")
// }

package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/uniris/uniris-core/api/pkg/adding"
	"github.com/uniris/uniris-core/api/pkg/listing"
)

/*
Scenario: Sign account creation result
	Given a keypair and a result for an account creation
	When I want sign it
	Then the signature is inserted and valid
*/
func TestSignAccountCreationResult(t *testing.T) {

	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pvKey, _ := x509.MarshalECPrivateKey(key)
	pubKey, _ := x509.MarshalPKIXPublicKey(key.Public())

	txID := adding.NewTransactionResult("hash", "ip", "sig")
	txKeychain := adding.NewTransactionResult("hash", "ip", "sig")
	txRes := adding.NewAccountCreationTransactionResult(txID, txKeychain)
	res := adding.NewAccountCreationResult(txRes, "")
	res, err := NewSigner().SignAccountCreationResult(res, hex.EncodeToString(pvKey))
	assert.NotEmpty(t, res.Signature())

	assert.Nil(t, err)

	b, _ := json.Marshal(accountCreationResult{
		Transactions: accountCreationTransactionsResult{
			ID: transactionResult{
				MasterPeerIP:    "ip",
				Signature:       "sig",
				TransactionHash: "hash",
			},
			Keychain: transactionResult{
				MasterPeerIP:    "ip",
				Signature:       "sig",
				TransactionHash: "hash",
			},
		},
	})

	assert.Nil(t, verifySignature(hex.EncodeToString(pubKey), string(b), res.Signature()))
}

/*
Scenario: Sign account search result
	Given a keypair and a result for an account search
	When I want sign it
	Then the signature is inserted and valid
*/
func TestSignAccountResult(t *testing.T) {

	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pvKey, _ := x509.MarshalECPrivateKey(key)
	pubKey, _ := x509.MarshalPKIXPublicKey(key.Public())

	res := listing.NewAccountResult("enc aes key", "enc wallet", "enc addr", "")

	b, err := json.Marshal(accountResult{
		EncryptedAESKey:  "enc aes key",
		EncryptedAddress: "enc addr",
		EncryptedWallet:  "enc wallet",
	})
	sig, err := sign(hex.EncodeToString(pvKey), string(b))
	assert.Nil(t, err)

	res = listing.NewAccountResult("enc aes key", "enc wallet", "enc addr", sig)
	assert.Nil(t, NewSigner().VerifyAccountResultSignature(res, hex.EncodeToString(pubKey)))
}

/*
Scenario: Verify creation result signature
	Given a keypair and signed creation result signature
	When I want to verify it
	Then I get not error
*/
func TestVerifyCreationResultSignature(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pvKey, _ := x509.MarshalECPrivateKey(key)
	pubKey, _ := x509.MarshalPKIXPublicKey(key.Public())

	b, _ := json.Marshal(transactionResult{
		MasterPeerIP:    "ip",
		TransactionHash: "hash",
	})
	sig, _ := sign(hex.EncodeToString(pvKey), string(b))
	res := adding.NewTransactionResult("hash", "ip", sig)

	assert.Nil(t, NewSigner().VerifyCreationTransactionResultSignature(res, hex.EncodeToString(pubKey)))
}

/*
Scenario: Verify account creation request signature
	Given a keypair and signed account creation request signature
	When I want to verify it
	Then I get not error
*/
func TestVerifyAccountCreationRequestSignature(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pvKey, _ := x509.MarshalECPrivateKey(key)
	pubKey, _ := x509.MarshalPKIXPublicKey(key.Public())

	b, _ := json.Marshal(accountCreationRequest{
		EncryptedID:       "enc id",
		EncryptedKeychain: "enc keychain",
	})
	sig, _ := sign(hex.EncodeToString(pvKey), string(b))
	req := adding.NewAccountCreationRequest("enc id", "enc keychain", sig)

	assert.Nil(t, NewSigner().VerifyAccountCreationRequestSignature(req, hex.EncodeToString(pubKey)))
}

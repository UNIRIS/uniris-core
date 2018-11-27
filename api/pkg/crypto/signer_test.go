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

	res := adding.AccountCreationResult{
		Transactions: adding.AccountCreationTransactionsResult{
			ID: adding.TransactionResult{
				MasterPeerIP:    "ip",
				Signature:       "sig",
				TransactionHash: "hash",
			},
		},
	}
	err := NewSigner().SignAccountCreationResult(&res, hex.EncodeToString(pvKey))
	assert.NotEmpty(t, res.Signature)

	assert.Nil(t, err)

	res2 := &adding.AccountCreationResult{
		Transactions: adding.AccountCreationTransactionsResult{
			ID: adding.TransactionResult{
				MasterPeerIP:    "ip",
				Signature:       "sig",
				TransactionHash: "hash",
			},
		},
	}

	b, _ := json.Marshal(res2)

	assert.Nil(t, verifySignature(hex.EncodeToString(pubKey), string(b), res.Signature))
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

	res := &listing.AccountResult{
		EncryptedAddress: "enc addr",
		EncryptedAESKey:  "enc aes key",
		EncryptedWallet:  "enc wallet",
	}

	b, err := json.Marshal(res)
	sig, err := sign(hex.EncodeToString(pvKey), string(b))
	assert.Nil(t, err)

	res.Signature = sig
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

	res := adding.TransactionResult{
		MasterPeerIP:    "ip",
		TransactionHash: "hash",
	}

	b, _ := json.Marshal(res)
	sig, _ := sign(hex.EncodeToString(pvKey), string(b))
	res.Signature = sig

	assert.Nil(t, NewSigner().VerifyCreationTransactionResultSignature(res, hex.EncodeToString(pubKey)))
}

/*
Scenario: Verify account search result signature
	Given a keypair and signed account search result signature
	When I want to verify it
	Then I get not error
*/
func TestVerifyAccountSearchResultSignature(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pvKey, _ := x509.MarshalECPrivateKey(key)
	pubKey, _ := x509.MarshalPKIXPublicKey(key.Public())

	res := &listing.AccountResult{
		EncryptedAddress: "enc address",
		EncryptedAESKey:  "enc aes",
		EncryptedWallet:  "wallet",
	}

	b, _ := json.Marshal(res)
	sig, _ := sign(hex.EncodeToString(pvKey), string(b))
	res.Signature = sig

	assert.Nil(t, NewSigner().VerifyAccountResultSignature(res, hex.EncodeToString(pubKey)))
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

	req := &adding.AccountCreationRequest{
		EncryptedID:       "enc id",
		EncryptedKeychain: "enc keychain",
	}

	b, _ := json.Marshal(req)
	sig, _ := sign(hex.EncodeToString(pvKey), string(b))
	req.Signature = sig

	assert.Nil(t, NewSigner().VerifyAccountCreationRequestSignature(req, hex.EncodeToString(pubKey)))
}

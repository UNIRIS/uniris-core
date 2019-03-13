package crypto

import (
	"crypto/elliptic"
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
Scenario: Version key by build a byte slice with the curve and the key marshalled
	Given a key marshalled and a curve
	When I want to version it
	Then I get a slice with the curve and the key marshalled
*/
func TestVersionKey(t *testing.T) {
	pv, _, _ := generateEd25519Keys(rand.Reader)
	bytes := pv.bytes()
	key := versionKey(pv.curve(), bytes)
	assert.Equal(t, Ed25519Curve, key.Curve())
	assert.Equal(t, pv.bytes(), key.Marshalling())
}

/*
Scenario: Parse a marshalled versionned ECDSA public key
	Given a key which has been versionned and marshalled (curve NIST P256)
	When I want to parse it
	Then I get an ECDSA public key with the good curve
*/
func TestParseEcdsaPublicKey(t *testing.T) {
	_, pub, _ := generateECDSAKeys(rand.Reader, elliptic.P256())
	versionnedKey, _ := pub.Marshal()
	publicKey, err := ParsePublicKey(versionnedKey)
	assert.Nil(t, err)
	assert.Equal(t, publicKey.curve(), P256Curve)
	assert.Equal(t, pub.bytes(), publicKey.bytes())
}

/*
Scenario: Parse a marshalled versionned Ed25519 public key
	Given a key which has been versionned and marshalled (curve ed25519)
	When I want to parse it
	Then I get an ed25519 public key
*/
func TestParseEd25519PublicKey(t *testing.T) {
	_, pub, _ := generateEd25519Keys(rand.Reader)
	versionnedKey, _ := pub.Marshal()
	publicKey, err := ParsePublicKey(versionnedKey)
	assert.Nil(t, err)
	assert.Equal(t, publicKey.curve(), Ed25519Curve)
	assert.Equal(t, pub.bytes(), publicKey.bytes())
}

/*
Scenario: Parse a marshalled versionned ECDSA private key
	Given a key which has been versionned and marshalled (curve NIST P256)
	When I want to parse it
	Then I get an ECDSA private key with the good curve
*/
func TestParseEcdsaPrivateKey(t *testing.T) {
	pv, _, _ := generateECDSAKeys(rand.Reader, elliptic.P256())
	versionnedKey, _ := pv.Marshal()
	privateKey, err := ParsePrivateKey(versionnedKey)
	assert.Nil(t, err)
	assert.Equal(t, privateKey.curve(), P256Curve)
	assert.Equal(t, pv.bytes(), privateKey.bytes())
}

/*
Scenario: Parse a marshalled versionned Ed25519 private key
	Given a key which has been versionned and marshalled (curve ed25519)
	When I want to parse it
	Then I get an ed25519 private key
*/
func TestParseEd25519PrivateKey(t *testing.T) {
	pv, _, _ := generateEd25519Keys(rand.Reader)
	versionnedKey, _ := pv.Marshal()
	privateKey, err := ParsePrivateKey(versionnedKey)
	assert.Nil(t, err)
	assert.Equal(t, privateKey.curve(), Ed25519Curve)
	assert.Equal(t, pv.bytes(), privateKey.bytes())
}

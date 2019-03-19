package crypto

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
Scenario: Get bytes from a ECDSA private key
	Given a ecdsa private key
	When I want to get bytes
	Then  I get a X509 marshalling of the key
*/
func TestGetBytesFromEcdsaPrivateKey(t *testing.T) {
	pv, _, _ := generateECDSAKeys(rand.Reader, elliptic.P256())
	b := pv.bytes()
	assert.NotEmpty(t, b)

	ecdsaKey := pv.(ecdsaPrivateKey)

	x, y := elliptic.Unmarshal(elliptic.P256(), b)
	assert.Equal(t, ecdsaKey.priv.X, x)
	assert.Equal(t, ecdsaKey.priv.Y, y)
}

/*
Scenario: Get curve from a ECDSA private key
	Given a ecdsa private key
	When I want to get curve
	Then  I get the curve associated to this key
*/
func TestGetCurveFromEcdsaPrivateKey(t *testing.T) {
	pv, _, _ := generateECDSAKeys(rand.Reader, elliptic.P256())
	c := pv.curve()
	assert.Equal(t, P256Curve, c)
}

/*
Scenario: Get curve from a ECDSA private key not supported
	Given a ecdsa private key
	When I want to get curve
	Then  I get the curve associated to this key
*/
func TestGetUnsupportedCurveFromEcdsaPrivateKey(t *testing.T) {
	pv, _, _ := generateECDSAKeys(rand.Reader, elliptic.P224())
	c := pv.curve()
	assert.Equal(t, Curve(-1), c)
}

/*
Scenario: Marshal an ECDSA private key preceed by its curve identity
	Given an ECDSA private key
	When I want to marshal it
	Then I get curve identity + key bytes (x509)
*/
func TestMarshalEcdsaPrivateKey(t *testing.T) {
	pv, _, _ := generateECDSAKeys(rand.Reader, elliptic.P256())
	b, err := pv.Marshal()
	assert.Nil(t, err)
	assert.NotEmpty(t, b)
	assert.Equal(t, b[0], byte(P256Curve))

	pvKey := pv.(ecdsaPrivateKey)
	b2, _ := x509.MarshalECPrivateKey(pvKey.priv)
	assert.Equal(t, P256Curve, b.Curve())
	assert.Equal(t, b2, b.Marshalling())

	ecdsaKey := pv.(ecdsaPrivateKey)
	key, err := x509.ParseECPrivateKey(b.Marshalling())
	assert.Nil(t, err)
	assert.Equal(t, ecdsaKey.priv.X, key.X)
	assert.Equal(t, ecdsaKey.priv.Y, key.Y)
}

/*
Scenario: Sign with an ECDSA private key
	Given an ECDSA private key and some data
	When I want to sign this data
	Then I get signature valid by ASN1
*/
func TestEcdsaSignWithPrivateKey(t *testing.T) {
	pv, pub, _ := generateECDSAKeys(rand.Reader, elliptic.P256())
	sig, err := pv.Sign([]byte("hello"))
	assert.Nil(t, err)
	assert.NotEmpty(t, sig)

	ecdsaSig := new(ecdsaSignature)
	_, err = asn1.Unmarshal(sig, ecdsaSig)
	assert.Nil(t, err)
	assert.NotNil(t, sig)

	hash := sha256.Sum256([]byte("hello"))
	ecdsaPub := pub.(ecdsaPublicKey)
	assert.True(t, ecdsa.Verify(ecdsaPub.pub, hash[:], ecdsaSig.R, ecdsaSig.S))
}

/*
Scenario: Marshal an ECDSA public key preceed by its curve identity
	Given an ECDSA public key
	When I want to marshal it
	Then I get curve identity + key bytes (x509)
*/
func TestMarshalEcdsaPublicKey(t *testing.T) {
	_, pub, _ := generateECDSAKeys(rand.Reader, elliptic.P256())
	b, err := pub.Marshal()
	assert.Nil(t, err)
	assert.NotEmpty(t, b)

	assert.Equal(t, P256Curve, b.Curve())

	ecdsaKey := pub.(ecdsaPublicKey)
	b2, _ := x509.MarshalPKIXPublicKey(ecdsaKey.pub)
	assert.Equal(t, P256Curve, b.Curve())
	assert.Equal(t, b2, b.Marshalling())

	key, err := x509.ParsePKIXPublicKey(b.Marshalling())
	ecdsaPub := key.(*ecdsa.PublicKey)
	assert.Nil(t, err)
	assert.Equal(t, ecdsaKey.pub.X, ecdsaPub.X)
	assert.Equal(t, ecdsaKey.pub.Y, ecdsaPub.Y)

}

/*
Scenario: Verify a signature with an ECDSA public key
	Given an ECDSA public key, a signature and the related data
	When I want to verify the signature
	Then I get not error and return true
*/
func TestEcdsaVerifyWithPublicKey(t *testing.T) {
	pv, pub, _ := generateECDSAKeys(rand.Reader, elliptic.P256())
	sig, _ := pv.Sign([]byte("hello"))
	assert.True(t, pub.Verify([]byte("hello"), sig))
}

/*
Scenario: Verify a signature with an ECDSA public key
	Given an ECDSA public key, a signature and the related data
	When I want to verify the signature
	Then I get not error and return true
*/
func TestEcdsaVerifyWithInvalidSignature(t *testing.T) {
	_, pub, _ := generateECDSAKeys(rand.Reader, elliptic.P256())
	assert.False(t, pub.Verify([]byte("hello"), []byte("fakesig")))
}

/*
Scenario: Generate an ECDSA shared key secret
	Given a public key and a private key
	When I want to generated a shared key secret
	Then I get the bytes of this secret and can obtain the same by inverting the public key and private key
*/
func TestGeneratedEcdsaShared(t *testing.T) {
	pv1, pub1, _ := generateECDSAKeys(rand.Reader, elliptic.P256())
	pv2, pub2, _ := generateECDSAKeys(rand.Reader, elliptic.P256())
	sharedSecret1, err := ecdsaGenerateShared(pub2, pv1)
	assert.Nil(t, err)
	assert.NotEmpty(t, sharedSecret1)

	sharedSecret2, err := ecdsaGenerateShared(pub1, pv2)
	assert.Nil(t, err)
	assert.NotEmpty(t, sharedSecret2)

	assert.Equal(t, sharedSecret1, sharedSecret2)
}

/*
Scenario: Extract the random public key from the cipher data
	Given a cipher data
	When I want to extract by unmarshal the public key
	Then I get the public key
*/
func TestExtractEcdsaRandomPubKeyFromCipher(t *testing.T) {
	_, pub, _ := generateECDSAKeys(rand.Reader, elliptic.P256())
	b := pub.bytes()

	rPub, pos, err := ecdsaExtractRandomPublicKey(elliptic.P256())(b)
	assert.Nil(t, err)
	assert.NotEmpty(t, rPub)
	assert.NotEqual(t, 0, pos)

	assert.Equal(t, pub, rPub)
}

/*
Scenario: Compare two ecdsa public keys
	Given a two ed25519 generated with the same secret
	When I want to compare
	Then I get a true response which said they are the same
*/
func TestEqualEcdsaPublicKey(t *testing.T) {
	_, pub, _ := generateECDSAKeys(bytes.NewBufferString("helloaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"), elliptic.P256())
	_, pub2, _ := generateECDSAKeys(bytes.NewBufferString("helloaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"), elliptic.P256())
	assert.True(t, pub.Equals(pub2))

	_, pub3, _ := generateECDSAKeys(bytes.NewBufferString("abcccoaaaaaaaaaaaaaabbbbbbaaaaaaaaaaaaaaaaaaaaaaa"), elliptic.P256())
	assert.False(t, pub.Equals(pub3))
}

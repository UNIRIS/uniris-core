package crypto

import (
	"bytes"
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/ed25519"
)

/*
Scenario: Get bytes from a Ed25519 private key
	Given a Ed25519 private key
	When I want to get bytes
	Then  I get a X509 marshalling of the key
*/
func TestGetBytesFromEd25519PrivateKey(t *testing.T) {
	pv, _, _ := generateEd25519Keys(rand.Reader)

	b := pv.Bytes()
	assert.NotEmpty(t, b)
	assert.Equal(t, len(b), ed25519.PrivateKeySize)
}

/*
Scenario: Get curve from a Ed25519 private key
	Given a Ed25519 private key
	When I want to get curve
	Then  I get the curve associated to this key
*/
func TestGetCurveFromEd25519PrivateKey(t *testing.T) {
	pv, _, _ := generateEd25519Keys(rand.Reader)
	c := pv.curve()
	assert.Equal(t, Ed25519Curve, c)
}

/*
Scenario: Sign with an Ed25519 private key
	Given an Ed25519 private key and some data
	When I want to sign this data
	Then I get signature valid
*/
func TestEd25519SignWithPrivateKey(t *testing.T) {
	pv, pub, _ := generateEd25519Keys(rand.Reader)
	sig, err := pv.Sign([]byte("hello"))
	assert.Nil(t, err)

	ed25519key := pub.(ed25519PublicKey)
	assert.True(t, ed25519.Verify(ed25519key.pub, []byte("hello"), sig))
}

/*
Scenario: Marshal an Ed25519 private key preceed by its curve identity
	Given an Ed25519 private key
	When I want to marshal it
	Then I get curve identity + key bytes (x509)
*/
func TestMarshalEd25519PrivateKey(t *testing.T) {
	pv, _, _ := generateEd25519Keys(rand.Reader)
	vKey, err := pv.Marshal()
	assert.Nil(t, err)
	assert.NotEmpty(t, vKey)
	assert.Equal(t, Ed25519Curve, vKey.Curve())

	pvKey := ed25519.PrivateKey(vKey.Marshalling())
	assert.Equal(t, []byte(pvKey), vKey.Marshalling())
}

/*
Scenario: Verify a signature with an Ed25519 public key
	Given an Ed25519 public key, a signature and the related data
	When I want to verify the signature
	Then I get not error and return true
*/
func TestEd25519VerifyWithPublicKey(t *testing.T) {
	pv, pub, _ := generateEd25519Keys(rand.Reader)
	sig, _ := pv.Sign([]byte("hello"))
	assert.True(t, pub.Verify([]byte("hello"), sig))
}

/*
Scenario: Verify a signature with an Ed25519 public key
	Given an Ed25519 public key, a signature and the related data
	When I want to verify the signature
	Then I get not error and return true
*/
func TestEd25519VerifyWithInvalidSignature(t *testing.T) {
	_, pub, _ := generateEd25519Keys(rand.Reader)
	assert.False(t, pub.Verify([]byte("hello"), []byte("fakesig")))
}

/*
Scenario: Marshal an Ed25519 public key preceed by its curve identity
	Given an Ed25519 public key
	When I want to marshal it
	Then I get curve identity + key bytes (x509)
*/
func TestMarshalEd25519PublicKey(t *testing.T) {
	_, pb, _ := generateEd25519Keys(rand.Reader)
	vKey, err := pb.Marshal()
	assert.Nil(t, err)
	assert.NotEmpty(t, vKey)
	assert.Equal(t, Ed25519Curve, vKey.Curve())

	pubKey := ed25519.PublicKey(vKey.Marshalling())
	assert.Equal(t, []byte(pubKey), vKey.Marshalling())
}

/*
Scenario: Generate an Ed25519 shared key secret
	Given a public key and a private key
	When I want to generated a shared key secret
	Then I get the bytes of this secret and can obtain the same by inverting the public key and private key
*/
func TestGeneratedEd25519Shared(t *testing.T) {
	pv1, pub1, _ := generateEd25519Keys(rand.Reader)
	pv2, pub2, _ := generateEd25519Keys(rand.Reader)
	sharedSecret1, err := ed25519GenerateShared(pub2, pv1)
	assert.Nil(t, err)
	assert.NotEmpty(t, sharedSecret1)

	sharedSecret2, err := ed25519GenerateShared(pub1, pv2)
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
func TestExtractEd25519RandomPubKeyFromCipher(t *testing.T) {
	_, pub, _ := generateEd25519Keys(rand.Reader)
	b := pub.Bytes()
	rPub, pos, err := ed25519ExtractRandomPublicKey(b)
	assert.Nil(t, err)
	assert.NotEmpty(t, rPub)
	assert.Equal(t, rPub, pub)
	assert.NotEqual(t, 0, pos)

}

/*
Scenario: Compare two ed25519 public keys
	Given a two ed25519 generated with the same secret
	When I want to compare
	Then I get a true response which said they are the same
*/
func TestEqualEd25519PublicKey(t *testing.T) {
	_, pub, _ := generateEd25519Keys(bytes.NewBufferString("helloaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"))
	_, pub2, _ := generateEd25519Keys(bytes.NewBufferString("helloaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"))
	_, pub3, _ := generateEd25519Keys(bytes.NewBufferString("helloaaaaaaaaagggggaaaaaaaaaaaaaaaaaaaaaaaaaaaa"))
	assert.True(t, pub.Equals(pub2))
	assert.False(t, pub.Equals(pub3))
}

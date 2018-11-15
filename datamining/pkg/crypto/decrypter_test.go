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
	"github.com/uniris/ecies/pkg"
)

/*
Scenario: Decrypt data
	Given encrypted bytes
	When I want decrypt it with ECIES
	Then I get the clear data
*/
func TestDecryptText(t *testing.T) {
	superKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pbKey, _ := x509.MarshalPKIXPublicKey(superKey.Public())
	pu, err := x509.ParsePKIXPublicKey(pbKey)
	robotEciesKey := ecies.ImportECDSAPublic(pu.(*ecdsa.PublicKey))
	cipher, err := ecies.Encrypt(rand.Reader, robotEciesKey, []byte("uniris"), nil, nil)
	assert.Nil(t, err)
	assert.NotEmpty(t, cipher)

	pvkey, _ := x509.MarshalECPrivateKey(superKey)

	clear, _ := decrypt(hex.EncodeToString(pvkey), hex.EncodeToString(cipher))
	assert.Equal(t, "uniris", string(clear))
}

/*
Scenario: Decrypt encrypted hash
	Given an hash address
	When I want decrypt it
	Then I get the hash
*/
func TestDecryptHash(t *testing.T) {
	superKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pbKey, _ := x509.MarshalPKIXPublicKey(superKey.Public())
	pu, _ := x509.ParsePKIXPublicKey(pbKey)
	robotEciesKey := ecies.ImportECDSAPublic(pu.(*ecdsa.PublicKey))
	cipher, _ := ecies.Encrypt(rand.Reader, robotEciesKey, []byte("hash"), nil, nil)

	pvkey, _ := x509.MarshalECPrivateKey(superKey)

	hash, err := NewDecrypter().DecryptHash(hex.EncodeToString(cipher), hex.EncodeToString(pvkey))
	assert.Nil(t, err)
	assert.Equal(t, "hash", hash)
}

/*
Scenario: Decrypt biometric data
	Given an encrypted biometric data
	When I want to decrypt it
	Then it serialize json and build biometric data without signatures
*/
func TestDecryptBiometricData(t *testing.T) {
	bioRaw := biometricRaw{
		PersonHash:          "hash",
		EncryptedAddrPerson: "addr",
		EncryptedAddrRobot:  "addr",
		EncryptedAESKey:     "enc aes key",
		PersonPublicKey:     "pub",
	}

	superKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pbKey, _ := x509.MarshalPKIXPublicKey(superKey.Public())
	pu, _ := x509.ParsePKIXPublicKey(pbKey)
	robotEciesKey := ecies.ImportECDSAPublic(pu.(*ecdsa.PublicKey))

	b, _ := json.Marshal(bioRaw)

	cipher, _ := ecies.Encrypt(rand.Reader, robotEciesKey, b, nil, nil)

	pvkey, _ := x509.MarshalECPrivateKey(superKey)

	bioData, err := NewDecrypter().DecryptBiometricData(hex.EncodeToString(cipher), hex.EncodeToString(pvkey))
	assert.Nil(t, err)
	assert.Equal(t, "hash", bioData.PersonHash())
	assert.Equal(t, "addr", bioData.CipherAddrPerson())
	assert.Equal(t, "addr", bioData.CipherAddrRobot())
	assert.Equal(t, "enc aes key", bioData.CipherAESKey())
	assert.Equal(t, "pub", bioData.PersonPublicKey())
}

/*
Scenario: Decrypt keychain data
	Given an encrypted keychain data
	When I want to decrypt it
	Then it serialize json and build biometric data without signatures
*/
func TestDecryptKeychaincData(t *testing.T) {
	bioRaw := keychainRaw{
		EncryptedWallet:    "enc wallet",
		EncryptedAddrRobot: "addr",
		PersonPublicKey:    "pub",
	}

	superKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pbKey, _ := x509.MarshalPKIXPublicKey(superKey.Public())
	pu, _ := x509.ParsePKIXPublicKey(pbKey)
	robotEciesKey := ecies.ImportECDSAPublic(pu.(*ecdsa.PublicKey))

	b, _ := json.Marshal(bioRaw)

	cipher, _ := ecies.Encrypt(rand.Reader, robotEciesKey, b, nil, nil)

	pvkey, _ := x509.MarshalECPrivateKey(superKey)

	keychainData, err := NewDecrypter().DecryptKeychainData(hex.EncodeToString(cipher), hex.EncodeToString(pvkey))
	assert.Nil(t, err)
	assert.Equal(t, "enc wallet", keychainData.CipherWallet())
	assert.Equal(t, "addr", keychainData.CipherAddrRobot())
	assert.Equal(t, "pub", keychainData.PersonPublicKey())
}

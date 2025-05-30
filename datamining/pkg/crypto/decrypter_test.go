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
Scenario: Decrypt ID
	Given an encrypted ID
	When I want to decrypt it
	Then it serialize json and build ID
*/
func TestDecryptID(t *testing.T) {
	id := id{
		Hash:                 "hash",
		EncryptedAddrByID:    "addr",
		EncryptedAddrByRobot: "addr",
		EncryptedAESKey:      "enc aes key",
		PublicKey:            "pub",
		Proposal: proposal{
			SharedEmitterKeyPair: proposalKeypair{
				EncryptedPrivateKey: "enc pv key",
				PublicKey:           "pub",
			},
		},
		IDSignature:      "sig",
		EmitterSignature: "sig",
	}

	superKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pbKey, _ := x509.MarshalPKIXPublicKey(superKey.Public())
	pu, _ := x509.ParsePKIXPublicKey(pbKey)
	robotEciesKey := ecies.ImportECDSAPublic(pu.(*ecdsa.PublicKey))

	b, _ := json.Marshal(id)

	cipher, _ := ecies.Encrypt(rand.Reader, robotEciesKey, b, nil, nil)

	pvkey, _ := x509.MarshalECPrivateKey(superKey)

	newID, err := NewDecrypter().DecryptID(hex.EncodeToString(cipher), hex.EncodeToString(pvkey))
	assert.Nil(t, err)
	assert.Equal(t, "hash", newID.Hash())
	assert.Equal(t, "addr", newID.EncryptedAddrByID())
	assert.Equal(t, "addr", newID.EncryptedAddrByRobot())
	assert.Equal(t, "enc aes key", newID.EncryptedAESKey())
	assert.Equal(t, "pub", newID.PublicKey())
	assert.Equal(t, "sig", newID.IDSignature())
	assert.Equal(t, "sig", newID.EmitterSignature())
}

/*
Scenario: Decrypt keychain data
	Given an encrypted keychain data
	When I want to decrypt it
	Then it serialize json and build ID data without signatures
*/
func TestDecryptKeychaincData(t *testing.T) {
	kc := keychain{
		EncryptedWallet:      "enc wallet",
		EncryptedAddrByRobot: "addr",
		IDPublicKey:          "pub",
		Proposal: proposal{
			SharedEmitterKeyPair: proposalKeypair{
				EncryptedPrivateKey: "enc pv key",
				PublicKey:           "pub",
			},
		},
		IDSignature:      "sig",
		EmitterSignature: "sig",
	}

	superKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pbKey, _ := x509.MarshalPKIXPublicKey(superKey.Public())
	pu, _ := x509.ParsePKIXPublicKey(pbKey)
	robotEciesKey := ecies.ImportECDSAPublic(pu.(*ecdsa.PublicKey))

	b, _ := json.Marshal(kc)

	cipher, _ := ecies.Encrypt(rand.Reader, robotEciesKey, b, nil, nil)

	pvkey, _ := x509.MarshalECPrivateKey(superKey)

	keychain, err := NewDecrypter().DecryptKeychain(hex.EncodeToString(cipher), hex.EncodeToString(pvkey))
	assert.Nil(t, err)
	assert.Equal(t, "enc wallet", keychain.EncryptedWallet())
	assert.Equal(t, "addr", keychain.EncryptedAddrByRobot())
	assert.Equal(t, "pub", keychain.IDPublicKey())
	assert.Equal(t, "sig", keychain.EmitterSignature())
	assert.Equal(t, "sig", keychain.IDSignature())
}

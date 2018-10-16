package internalrpc

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	internalrpc "github.com/uniris/uniris-core/datamining/pkg/transport/rpc/internalrpc"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/uniris/ecies/pkg"

	api "github.com/uniris/uniris-core/datamining/api/protobuf-spec"
)

/*
Scenerio: Decrypted the sent encrypted wallet data
	Given a wallet containing encrypted datat
	When I want to decrypt it
	Then I get the cleared data
*/
func TestDecryptData(t *testing.T) {

	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	eciesKey := ecies.ImportECDSA(key)
	pvKey, _ := x509.MarshalECPrivateKey(key)

	encBioData, _ := ecies.Encrypt(rand.Reader, &eciesKey.PublicKey, []byte("bio data"), nil, nil)
	encWalData, _ := ecies.Encrypt(rand.Reader, &eciesKey.PublicKey, []byte("wallet data"), nil, nil)

	w := &api.Wallet{
		EncryptedBioData:    []byte(hex.EncodeToString(encBioData)),
		EncryptedWalletData: []byte(hex.EncodeToString(encWalData)),
	}

	decodeBio, decodeWal, err := DecryptWallet(w, []byte(hex.EncodeToString(pvKey)))
	assert.Nil(t, err)
	assert.NotNil(t, decodeBio)
	assert.NotNil(t, decodeWal)

	assert.Equal(t, "bio data", string(decodeBio))
	assert.Equal(t, "wallet data", string(decodeWal))
}

/*
Scenario: Verify bio data signatures
	Given a bio data and a signature
	When I want to check the signature
	Then no error is returned and the data is valid
*/
func TestVerifyBioSignatures(t *testing.T) {

	bioKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	bioPubKey, _ := x509.MarshalPKIXPublicKey(bioKey.Public())
	bioPvKey, _ := x509.MarshalECPrivateKey(bioKey)
	personKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	personPubKey, _ := x509.MarshalPKIXPublicKey(personKey.Public())
	personPvKey, _ := x509.MarshalECPrivateKey(personKey)

	bioJSON := BioDataFromJSON{
		BiodPublicKey:       hex.EncodeToString(bioPubKey),
		EncryptedAddrPerson: "encrypted addr",
		EncryptedAddrRobot:  "encrypted addr",
		EncryptedAESKey:     "encrypted aes key",
		PersonHash:          "person hash",
		PersonPublicKey:     hex.EncodeToString(personPubKey),
	}

	b, _ := json.Marshal(bioJSON)

	sigBio, _ := crypto.Sign([]byte(hex.EncodeToString(bioPvKey)), b)
	sigPerson, _ := crypto.Sign([]byte(hex.EncodeToString(personPvKey)), b)

	assert.Nil(t, VerifyBioSignatures(bioJSON, &api.Signature{
		Biod:   sigBio,
		Person: sigPerson,
	}))
}

/*
Scenario: Verify wallet data signatures
	Given a wallet data and a signature
	When I want to check the signature
	Then no error is returned and the data is valid
*/
func TestVerifyWalletSignatures(t *testing.T) {

	bioKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	bioPubKey, _ := x509.MarshalPKIXPublicKey(bioKey.Public())
	bioPvKey, _ := x509.MarshalECPrivateKey(bioKey)
	personKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	personPubKey, _ := x509.MarshalPKIXPublicKey(personKey.Public())
	personPvKey, _ := x509.MarshalECPrivateKey(personKey)

	walJSON := WalletDataFromJSON{
		BiodPublicKey:      hex.EncodeToString(bioPubKey),
		EncryptedAddrRobot: "encrypted addr",
		EncryptedWallet:    "encrypted wallet",
		PersonPublicKey:    hex.EncodeToString(personPubKey),
	}

	b, _ := json.Marshal(walJSON)

	sigBio, _ := crypto.Sign([]byte(hex.EncodeToString(bioPvKey)), b)
	sigPerson, _ := crypto.Sign([]byte(hex.EncodeToString(personPvKey)), b)

	assert.Nil(t, VerifyWalSignatures(walJSON, &api.Signature{
		Biod:   sigBio,
		Person: sigPerson,
	}))
}

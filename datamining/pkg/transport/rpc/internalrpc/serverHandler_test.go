package internalrpc

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"testing"
	"time"

	mockcrypto "github.com/uniris/uniris-core/datamining/pkg/crypto/mock"
	"github.com/uniris/uniris-core/datamining/pkg/mining/master"
	mockstorage "github.com/uniris/uniris-core/datamining/pkg/storage/mock"
	mocktransport "github.com/uniris/uniris-core/datamining/pkg/transport/mock"

	"github.com/uniris/uniris-core/datamining/pkg/system"

	"github.com/uniris/uniris-core/datamining/pkg/listing"

	"github.com/stretchr/testify/assert"
	api "github.com/uniris/uniris-core/datamining/api/protobuf-spec"
	datamining "github.com/uniris/uniris-core/datamining/pkg"
	"github.com/uniris/uniris-core/datamining/pkg/crypto"
)

/*
Scenario: Retrieve an account from a bio hash
	Given a biometric stored and a bio hash
	When i want to retrieve the account associated
	Then i can retrieve the account stored
*/
func TestGetAccount(t *testing.T) {

	repo := mockstorage.NewDatabase()

	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pbKey, _ := x509.MarshalPKIXPublicKey(key.Public())
	pvKey, _ := x509.MarshalECPrivateKey(key)

	cipherAddr, _ := crypto.Encrypt(hex.EncodeToString(pbKey), "addr")
	cipherBhash, _ := crypto.Encrypt(hex.EncodeToString(pbKey), "hash")

	bdata := &datamining.BioData{
		PersonHash:      "hash",
		CipherAddrRobot: cipherAddr,
		CipherAESKey:    "encrypted_aes_key",
	}
	endors := datamining.NewEndorsement(
		time.Now(),
		"hello",
		&datamining.MasterValidation{},
		[]datamining.Validation{},
	)
	repo.StoreBiometric(datamining.NewBiometric(bdata, endors))

	kdata := &datamining.KeyChainData{
		WalletAddr:      "addr",
		CipherAddrRobot: cipherAddr,
	}
	repo.StoreKeychain(datamining.NewKeychain(kdata, endors, ""))

	list := listing.NewService(repo)
	errors := system.DataMininingErrors{}

	master := master.NewService(
		mocktransport.NewPoolFinder(),
		mocktransport.NewPoolRequester(),
		mocktransport.NewNotifier(),
		mockcrypto.NewSigner(),
		mockcrypto.NewHasher(),
		list,
		"robotPubKey",
		"robotPvKey",
	)

	h := NewInternalServerHandler(list, master, hex.EncodeToString(pvKey), errors)
	res, err := h.GetAccount(context.TODO(), &api.AccountSearchRequest{
		EncryptedHashPerson: cipherBhash,
	})
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "encrypted_aes_key", res.EncryptedAESkey)
}

/*
Scenario: Create an account
	Given some keychain data and biometric data encrypted
	When I want to store the account
	Then the account is stored and the transactions hash are returned
*/
func TestCreateAccount(t *testing.T) {
	repo := mockstorage.NewDatabase()

	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	pbKey, _ := x509.MarshalPKIXPublicKey(key.Public())
	pvKey, _ := x509.MarshalECPrivateKey(key)

	list := listing.NewService(repo)
	leading := master.NewService(
		mocktransport.NewPoolFinder(),
		mocktransport.NewPoolRequester(),
		mocktransport.NewNotifier(),
		mockcrypto.NewSigner(),
		mockcrypto.NewHasher(),
		list,
		"robotPubKey",
		"robotPvKey",
	)
	errors := system.DataMininingErrors{}

	h := NewInternalServerHandler(list, leading, hex.EncodeToString(pvKey), errors)

	cipherAddr, _ := crypto.Encrypt(hex.EncodeToString(pbKey), "encrypted_addr_robot")

	bioData := BioDataFromJSON{
		BiodPublicKey:       hex.EncodeToString(pbKey),
		EncryptedAddrPerson: "encrypted_person_addr",
		EncryptedAddrRobot:  cipherAddr,
		EncryptedAESKey:     "encrypted_aes_key",
		PersonHash:          "person hash",
		PersonPublicKey:     hex.EncodeToString(pbKey),
	}

	keychainData := KeychainDataFromJSON{
		BiodPublicKey:      hex.EncodeToString(pbKey),
		EncryptedAddrRobot: cipherAddr,
		EncryptedWallet:    "encrypted_wallet",
		PersonPublicKey:    hex.EncodeToString(pbKey),
	}

	bioBytes, _ := json.Marshal(bioData)
	cipherBio, _ := crypto.Encrypt(hex.EncodeToString(pbKey), string(bioBytes))
	sigBio, _ := crypto.Sign(hex.EncodeToString(pvKey), string(bioBytes))

	keychainBytes, _ := json.Marshal(keychainData)
	cipherWallet, _ := crypto.Encrypt(hex.EncodeToString(pbKey), string(keychainBytes))
	sigKeychain, _ := crypto.Sign(hex.EncodeToString(pvKey), string(keychainBytes))

	req := &api.AccountCreationRequest{
		EncryptedBioData:      cipherBio,
		EncryptedKeychainData: cipherWallet,
		SignatureBioData: &api.Signature{
			Biod:   sigBio,
			Person: sigBio,
		},
		SignatureKeychainData: &api.Signature{
			Biod:   sigKeychain,
			Person: sigKeychain,
		},
	}

	res, err := h.CreateAccount(context.TODO(), req)
	assert.Nil(t, err)
	assert.NotNil(t, res)

	assert.NotNil(t, res.BioTransactionHash)
	assert.NotNil(t, res.KeychainTransactionHash)
}

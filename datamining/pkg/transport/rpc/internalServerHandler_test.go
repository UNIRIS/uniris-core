package rpc

import (
	"context"
	"net"
	"testing"

	"github.com/uniris/uniris-core/datamining/pkg"

	"github.com/stretchr/testify/assert"

	api "github.com/uniris/uniris-core/datamining/api/protobuf-spec"
	"github.com/uniris/uniris-core/datamining/pkg/account"
	"github.com/uniris/uniris-core/datamining/pkg/system"
)

/*
Scenario: Get account
	Given a person hash
	When I want to get the account
	Then I get the encrypted keychain and biometric data
*/
func TestGetAccount(t *testing.T) {
	conf := system.UnirisConfig{}
	crypto := Crypto{
		hasher:    mockHasher{},
		decrypter: mockDecrypter{},
		signer:    mockSigner{},
	}
	db := &mockDatabase{}

	db.StoreBiometric(
		account.NewBiometric(
			account.NewBiometricData("hash", "enc addr", "enc addr", "enc aes key", "pub", "pub", account.NewSignatures("sig", "sig")),
			nil,
		),
	)

	db.StoreKeychain(
		account.NewKeychain(
			"hash",
			account.NewKeychainData("enc addr", "enc wallet", "pub", "pub", account.NewSignatures("sig", "sig")),
			nil,
		),
	)

	srvHandler := NewInternalServerHandler(mockPoolRequester{
		repo: db,
	}, mockAIClient{}, crypto, conf)

	res, err := srvHandler.GetAccount(context.TODO(), &api.AccountSearchRequest{
		EncryptedHashPerson: "enc person hash",
	})
	assert.Nil(t, err)
	assert.Equal(t, "enc wallet", res.EncryptedWallet)
	assert.Equal(t, "enc aes key", res.EncryptedAESkey)
	assert.Equal(t, "enc addr", res.EncryptedAddress)
	assert.Equal(t, "sig", res.Signature)
}

/*
Scenario: Create keychain
	Given a keychain creation request
	When I want create it
	Then the mining process started and the keychain is stored
*/
func TestCreateKeychain(t *testing.T) {
	conf := system.UnirisConfig{}
	crypto := Crypto{
		hasher:    mockHasher{},
		decrypter: mockDecrypter{},
		signer:    mockSigner{},
	}
	srvHandler := NewInternalServerHandler(mockPoolRequester{}, mockAIClient{}, crypto, conf)

	res, err := srvHandler.CreateKeychain(context.TODO(), &api.KeychainCreationRequest{
		EncryptedKeychainData: "cipher data",
		SignatureKeychainData: &api.Signature{
			Biod:   "sig",
			Person: "sig",
		},
	})
	assert.Nil(t, err)
	assert.Equal(t, "hash", res.TransactionHash)
	assert.Equal(t, "127.0.0.1", res.MasterPeerIP)
	assert.Equal(t, "sig", res.Signature)
}

/*
Scenario: Create biometric
	Given a biometric creation request
	When I want create it
	Then the mining process started and the keychain is stored
*/
func TestCreateBiometric(t *testing.T) {

	conf := system.UnirisConfig{}
	crypto := Crypto{
		hasher:    mockHasher{},
		decrypter: mockDecrypter{},
		signer:    mockSigner{},
	}
	srvHandler := NewInternalServerHandler(mockPoolRequester{}, mockAIClient{}, crypto, conf)

	res, err := srvHandler.CreateBiometric(context.TODO(), &api.BiometricCreationRequest{
		EncryptedBiometricData: "cipher data",
		SignatureBiometricData: &api.Signature{
			Biod:   "sig",
			Person: "sig",
		},
	})
	assert.Nil(t, err)
	assert.Equal(t, "hash", res.TransactionHash)
	assert.Equal(t, "127.0.0.1", res.MasterPeerIP)
	assert.Equal(t, "sig", res.Signature)
}

type mockAIClient struct {
}

func (c mockAIClient) GetStoragePool(personHash string) (datamining.Pool, error) {
	return datamining.NewPool(datamining.Peer{IP: net.ParseIP("127.0.0.1")}), nil
}

func (c mockAIClient) GetMasterPeer(txHash string) (datamining.Peer, error) {
	return datamining.Peer{IP: net.ParseIP("127.0.0.1")}, nil
}

func (c mockAIClient) GetValidationPool(txHash string) (datamining.Pool, error) {
	return datamining.NewPool(datamining.Peer{IP: net.ParseIP("127.0.0.1")}), nil
}

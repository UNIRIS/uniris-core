package internalrpc

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/uniris/uniris-core/datamining/pkg"

	"github.com/stretchr/testify/assert"

	api "github.com/uniris/uniris-core/datamining/api/protobuf-spec"
	"github.com/uniris/uniris-core/datamining/pkg/account"
	"github.com/uniris/uniris-core/datamining/pkg/system"

	"github.com/uniris/uniris-core/datamining/pkg/account/listing"
	"github.com/uniris/uniris-core/datamining/pkg/mock"
)

/*
Scenario: Get account
	Given a person hash
	When I want to get the account
	Then I get the encrypted keychain and biometric data
*/
func TestGetAccount(t *testing.T) {
	db := mock.NewDatabase()
	accountLister := listing.NewService(db)
	conf := system.UnirisConfig{}
	srvHandler := NewInternalServerHandler(accountLister, mockHasher{}, mockDecrypter{}, conf)

	db.StoreBiometric(account.NewBiometric(
		&account.BioData{
			CipherAddrBio: "enc address",
			CipherAESKey:  "cipher aes",
			PersonHash:    "hash",
		},
		datamining.NewEndorsement("", "hash", nil, nil),
	))

	db.StoreKeychain(account.NewKeychain(
		&account.KeyChainData{
			CipherWallet: "cipher wallet",
			WalletAddr:   "address",
		},
		datamining.NewEndorsement("", "hash", nil, nil),
	))

	res, err := srvHandler.GetAccount(context.TODO(), &api.AccountSearchRequest{
		EncryptedHashPerson: "enc person hash",
	})
	assert.Nil(t, err)
	assert.Equal(t, "cipher wallet", res.EncryptedWallet)
	assert.Equal(t, "cipher aes", res.EncryptedAESkey)
	assert.Equal(t, "enc address", res.EncryptedAddress)
}

/*
Scenario: Create keychain
	Given a keychain creation request
	When I want create it
	Then the mining process started and the keychain is stored
*/
func TestCreateKeychain(t *testing.T) {
	db := mock.NewDatabase()
	accountLister := listing.NewService(db)
	conf := system.UnirisConfig{}
	srvHandler := NewInternalServerHandler(accountLister, mockHasher{}, mockKeychainDecrypter{}, conf)
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
}

/*
Scenario: Create biometric
	Given a biometric creation request
	When I want create it
	Then the mining process started and the keychain is stored
*/
func TestCreateBiometric(t *testing.T) {
	db := mock.NewDatabase()
	accountLister := listing.NewService(db)
	conf := system.UnirisConfig{}
	srvHandler := NewInternalServerHandler(accountLister, mockHasher{}, mockBiometricDecrypter{}, conf)
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
}

type mockHasher struct{}

func (h mockHasher) HashKeychainJSON(*KeychainDataFromJSON) (string, error) {
	return "hash", nil
}
func (h mockHasher) HashBiometricJSON(*BioDataFromJSON) (string, error) {
	return "hash", nil

}

type mockKeychainDecrypter struct{}

func (d mockKeychainDecrypter) DecryptHashPerson(hash string, pvKey string) (string, error) {
	return "hash", nil
}

func (d mockKeychainDecrypter) DecryptCipherAddress(cipherAddr string, pvKey string) (string, error) {
	return "address", nil
}
func (d mockKeychainDecrypter) DecryptTransactionData(data string, pvKey string) (string, error) {
	keychainJSON := KeychainDataFromJSON{
		EncryptedWallet:    "cipher wallet",
		BiodPublicKey:      "pubk",
		PersonPublicKey:    "pubk",
		EncryptedAddrRobot: "cipher addr",
	}
	b, _ := json.Marshal(keychainJSON)
	return string(b), nil
}

type mockBiometricDecrypter struct{}

func (d mockBiometricDecrypter) DecryptHashPerson(hash string, pvKey string) (string, error) {
	return "hash", nil
}

func (d mockBiometricDecrypter) DecryptCipherAddress(cipherAddr string, pvKey string) (string, error) {
	return "address", nil
}
func (d mockBiometricDecrypter) DecryptTransactionData(data string, pvKey string) (string, error) {
	biometricJSON := BioDataFromJSON{
		EncryptedAddrPerson: "cipher addr",
		EncryptedAESKey:     "cipher aes",
		PersonHash:          "person hash",
		BiodPublicKey:       "pubk",
		PersonPublicKey:     "pubk",
		EncryptedAddrRobot:  "cipher addr",
	}
	b, _ := json.Marshal(biometricJSON)
	return string(b), nil
}

type mockDecrypter struct{}

func (d mockDecrypter) DecryptHashPerson(hash string, pvKey string) (string, error) {
	return "hash", nil
}

func (d mockDecrypter) DecryptCipherAddress(cipherAddr string, pvKey string) (string, error) {
	return "address", nil
}
func (d mockDecrypter) DecryptTransactionData(data string, pvKey string) (string, error) {
	return "data", nil
}

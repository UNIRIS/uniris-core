package rpc

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	api "github.com/uniris/uniris-core/datamining/api/protobuf-spec"
	"github.com/uniris/uniris-core/datamining/pkg/account"
	mockcrypto "github.com/uniris/uniris-core/datamining/pkg/crypto/mock"
	mockstorage "github.com/uniris/uniris-core/datamining/pkg/storage/mock"
	"github.com/uniris/uniris-core/datamining/pkg/system"
	mocktransport "github.com/uniris/uniris-core/datamining/pkg/transport/mock"
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
		decrypter: mockcrypto.NewDecrypter(),
		signer:    mockcrypto.NewSigner(),
		hasher:    mockcrypto.NewHasher(),
	}
	db := mockstorage.NewDatabase()
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

	cli := mocktransport.NewExternalClient(db)
	poolR := mocktransport.NewPoolRequester(cli)
	srvHandler := NewInternalServerHandler(poolR, mocktransport.NewAIClient(), crypto, conf)

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
		decrypter: mockcrypto.NewDecrypter(),
		signer:    mockcrypto.NewSigner(),
		hasher:    mockcrypto.NewHasher(),
	}
	db := mockstorage.NewDatabase()
	cli := mocktransport.NewExternalClient(db)
	poolR := mocktransport.NewPoolRequester(cli)
	aiCli := mocktransport.NewAIClient()
	srvHandler := NewInternalServerHandler(poolR, aiCli, crypto, conf)

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
		decrypter: mockcrypto.NewDecrypter(),
		signer:    mockcrypto.NewSigner(),
		hasher:    mockcrypto.NewHasher(),
	}
	db := mockstorage.NewDatabase()
	cli := mocktransport.NewExternalClient(db)
	poolR := mocktransport.NewPoolRequester(cli)
	aiCli := mocktransport.NewAIClient()
	srvHandler := NewInternalServerHandler(poolR, aiCli, crypto, conf)

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

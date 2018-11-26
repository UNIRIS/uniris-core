package rpc

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	api "github.com/uniris/uniris-core/datamining/api/protobuf-spec"
	datamining "github.com/uniris/uniris-core/datamining/pkg"
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
	Then I get the encrypted keychain and ID data
*/
func TestGetAccount(t *testing.T) {
	conf := system.UnirisConfig{}
	crypto := Crypto{
		decrypter: mockcrypto.NewDecrypter(),
		signer:    mockcrypto.NewSigner(),
		hasher:    mockcrypto.NewHasher(),
	}
	prop := datamining.NewProposal(
		datamining.NewProposedKeyPair("enc pv key", "pub key"),
	)

	db := mockstorage.NewDatabase()
	db.StoreID(
		account.NewEndorsedID(
			account.NewID("hash", "enc addr", "enc addr", "enc aes key", "id pub", "id sig", "em sig", prop),
			nil,
		),
	)

	db.StoreKeychain(
		account.NewEndorsedKeychain(
			"hash",
			account.NewKeychain("enc addr", "enc wallet", "id pub", "id sig", "em sig", prop),
			nil,
		),
	)

	cli := mocktransport.NewExternalClient(db)
	poolR := mocktransport.NewPoolRequester(cli)
	srvHandler := NewInternalServerHandler(poolR, mocktransport.NewAIClient(), crypto, conf)

	res, err := srvHandler.GetAccount(context.TODO(), &api.AccountSearchRequest{
		EncryptedIDHash: "enc id hash",
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
		EncryptedKeychain: "cipher data",
	})
	assert.Nil(t, err)
	assert.Equal(t, "hash", res.TransactionHash)
	assert.Equal(t, "127.0.0.1", res.MasterPeerIP)
	assert.Equal(t, "sig", res.Signature)
}

/*
Scenario: Create ID
	Given a ID creation request
	When I want create it
	Then the mining process started and the keychain is stored
*/
func TestCreateID(t *testing.T) {

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

	res, err := srvHandler.CreateID(context.TODO(), &api.IDCreationRequest{
		EncryptedID: "cipher data",
	})
	assert.Nil(t, err)
	assert.Equal(t, "hash", res.TransactionHash)
	assert.Equal(t, "127.0.0.1", res.MasterPeerIP)
	assert.Equal(t, "sig", res.Signature)
}

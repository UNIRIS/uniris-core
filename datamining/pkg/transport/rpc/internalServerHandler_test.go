package rpc

import (
	"context"
	"testing"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/uniris/uniris-core/datamining/pkg/emitter"

	"github.com/stretchr/testify/assert"

	api "github.com/uniris/uniris-core/datamining/api/protobuf-spec"
	datamining "github.com/uniris/uniris-core/datamining/pkg"
	"github.com/uniris/uniris-core/datamining/pkg/account"
	mockcrypto "github.com/uniris/uniris-core/datamining/pkg/crypto/mock"
	emlisting "github.com/uniris/uniris-core/datamining/pkg/emitter/listing"
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
			account.NewID("hash", "enc addr", "enc addr", "enc aes key", "id pub", prop, "id sig", "em sig"),
			nil,
		),
	)

	db.StoreKeychain(
		account.NewEndorsedKeychain(
			"hash",
			account.NewKeychain("enc addr", "enc wallet", "id pub", prop, "id sig", "em sig"),
			nil,
		),
	)

	cli := mocktransport.NewExternalClient(db)
	poolR := mocktransport.NewPoolRequester(cli)

	emLister := emlisting.NewService(db)
	srvHandler := NewInternalServerHandler(emLister, poolR, mocktransport.NewAIClient(), crypto, conf)

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
	emLister := emlisting.NewService(db)
	srvHandler := NewInternalServerHandler(emLister, poolR, aiCli, crypto, conf)

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

	emLister := emlisting.NewService(db)
	srvHandler := NewInternalServerHandler(emLister, poolR, aiCli, crypto, conf)

	res, err := srvHandler.CreateID(context.TODO(), &api.IDCreationRequest{
		EncryptedID: "cipher data",
	})
	assert.Nil(t, err)
	assert.Equal(t, "hash", res.TransactionHash)
	assert.Equal(t, "127.0.0.1", res.MasterPeerIP)
	assert.Equal(t, "sig", res.Signature)
}

/*
Scenario: Check if emitter is authorized
	Given a emitter public key
	When I want to check if it's authorized
	Then I get no error
*/
func TestIsAuthorized(t *testing.T) {

	conf := system.UnirisConfig{}
	crypto := Crypto{
		decrypter: mockcrypto.NewDecrypter(),
		signer:    mockcrypto.NewSigner(),
		hasher:    mockcrypto.NewHasher(),
	}

	db := mockstorage.NewDatabase()
	emLister := emlisting.NewService(db)
	cli := mocktransport.NewExternalClient(db)
	poolR := mocktransport.NewPoolRequester(cli)
	aiCli := mocktransport.NewAIClient()

	srvHandler := NewInternalServerHandler(emLister, poolR, aiCli, crypto, conf)

	_, err := srvHandler.IsEmitterAuthorized(context.TODO(), &api.AuthorizationRequest{
		PublicKey: "pubkey",
	})

	assert.Nil(t, err)

}

/*
Scenario: Get shared keys
	Given shared keys already stored
	When I want get the shared keys
	Then I get the shared keys
*/
func TestGetSharedKeys(t *testing.T) {

	conf := system.UnirisConfig{
		SharedKeys: system.SharedKeys{
			Robot: system.KeyPair{
				PrivateKey: "pv key",
				PublicKey:  "pub key",
			},
		},
	}
	crypto := Crypto{
		decrypter: mockcrypto.NewDecrypter(),
		signer:    mockcrypto.NewSigner(),
		hasher:    mockcrypto.NewHasher(),
	}

	db := mockstorage.NewDatabase()
	emLister := emlisting.NewService(db)
	cli := mocktransport.NewExternalClient(db)
	poolR := mocktransport.NewPoolRequester(cli)
	aiCli := mocktransport.NewAIClient()

	db.StoreSharedEmitterKeyPair(emitter.SharedKeyPair{
		PublicKey:           "pub key",
		EncryptedPrivateKey: "enc pv key",
	})

	srvHandler := NewInternalServerHandler(emLister, poolR, aiCli, crypto, conf)

	res, err := srvHandler.GetSharedKeys(context.TODO(), &empty.Empty{})
	assert.Nil(t, err)

	assert.Equal(t, "pv key", res.RobotPrivateKey)
	assert.Equal(t, "pub key", res.RobotPublicKey)
	assert.Equal(t, "enc pv key", res.EmitterKeys[0].EncryptedPrivateKey)
	assert.Equal(t, "pub key", res.EmitterKeys[0].PublicKey)
}

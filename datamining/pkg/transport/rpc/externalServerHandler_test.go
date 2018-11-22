package rpc

import (
	"context"
	"log"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	api "github.com/uniris/uniris-core/datamining/api/protobuf-spec"
	datamining "github.com/uniris/uniris-core/datamining/pkg"
	"github.com/uniris/uniris-core/datamining/pkg/account"
	accountadding "github.com/uniris/uniris-core/datamining/pkg/account/adding"
	accountListing "github.com/uniris/uniris-core/datamining/pkg/account/listing"
	accountMining "github.com/uniris/uniris-core/datamining/pkg/account/mining"
	mockcrypto "github.com/uniris/uniris-core/datamining/pkg/crypto/mock"
	emListing "github.com/uniris/uniris-core/datamining/pkg/emitter/listing"
	"github.com/uniris/uniris-core/datamining/pkg/lock"
	"github.com/uniris/uniris-core/datamining/pkg/mining"
	mockstorage "github.com/uniris/uniris-core/datamining/pkg/storage/mock"
	"github.com/uniris/uniris-core/datamining/pkg/system"
	mocktransport "github.com/uniris/uniris-core/datamining/pkg/transport/mock"
)

/*
Scenario: Get ID
	Given a id hash and ID already stored
	When i want to retrive it
	Then I get it from the db
*/
func TestGetID(t *testing.T) {
	db := mockstorage.NewDatabase()
	accLister := accountListing.NewService(db)

	id := account.NewID("hash", "enc addr", "enc addr", "enc aes key", "id pub", "id sig", "em sig")
	endors := mining.NewEndorsement("", "hash",
		mining.NewMasterValidation([]string{"hash"}, "robotkey", mining.NewValidation(mining.ValidationOK, time.Now(), "pub key", "sig")),
		[]mining.Validation{})

	db.StoreID(account.NewEndorsedID(id, endors))

	srv := Services{accLister: accLister}
	crypto := Crypto{
		decrypter: mockcrypto.NewDecrypter(),
		signer:    mockcrypto.NewSigner(),
	}
	h := NewExternalServerHandler(srv, crypto, system.UnirisConfig{})

	res, err := h.GetID(context.TODO(), &api.IDRequest{
		EncryptedIDHash: "enc hash",
	})
	assert.Nil(t, err)
	assert.Equal(t, "sig", res.Signature)
	assert.Equal(t, "enc aes key", res.Data.EncryptedAESKey)
}

/*
Scenario: Get keychain
	Given an address and keychain already stored
	When i want to retrive it
	Then I get it from the db
*/
func TestGetKeychain(t *testing.T) {
	db := mockstorage.NewDatabase()
	accLister := accountListing.NewService(db)

	kc := account.NewKeychain("enc address", "enc wallet", "id pub", "id sig", "em sig")
	endors := mining.NewEndorsement("", "hash",
		mining.NewMasterValidation([]string{"hash"}, "robotkey", mining.NewValidation(mining.ValidationOK, time.Now(), "pub key", "sig")),
		[]mining.Validation{})

	db.StoreKeychain(account.NewEndorsedKeychain("hash", kc, endors))

	srv := Services{accLister: accLister}
	crypto := Crypto{
		decrypter: mockcrypto.NewDecrypter(),
		signer:    mockcrypto.NewSigner(),
	}
	h := NewExternalServerHandler(srv, crypto, system.UnirisConfig{})
	res, err := h.GetKeychain(context.TODO(), &api.KeychainRequest{
		EncryptedAddress: "enc address",
	})

	assert.Nil(t, err)
	assert.Equal(t, "sig", res.Signature)
	assert.Equal(t, "enc wallet", res.Data.EncryptedWallet)
}

/*
Scenario: Lead mining of keychain transaction
	Given a lead keychain request
	When I mine the transaction
	Then I got no error and the data is stored
*/
func TestLeadKeychainMining(t *testing.T) {
	db := mockstorage.NewDatabase()
	lockSrv := lock.NewService(db)
	notifier := mockNotifier{}
	poolF := mockPoolFinder{}
	cli := mocktransport.NewExternalClient(db)
	poolR := mocktransport.NewPoolRequester(cli)
	emLister := emListing.NewService(db)
	accLister := accountListing.NewService(db)
	aiClient := mocktransport.NewAIClient()

	txMiners := map[mining.TransactionType]mining.TransactionMiner{
		mining.KeychainTransaction: accountMining.NewKeychainMiner(mockcrypto.NewSigner(), mockcrypto.NewHasher(), accLister),
	}

	mineSrv := mining.NewService(aiClient, notifier, poolF, poolR, mockcrypto.NewSigner(), emLister, system.UnirisConfig{}, txMiners)

	accAdder := accountadding.NewService(aiClient, db, accLister, mockcrypto.NewSigner(), mockcrypto.NewHasher())

	srv := Services{accAdd: accAdder, lock: lockSrv, mining: mineSrv}
	crypto := Crypto{
		decrypter: mockcrypto.NewDecrypter(),
		signer:    mockcrypto.NewSigner(),
		hasher:    mockcrypto.NewHasher(),
	}
	h := NewExternalServerHandler(srv, crypto, system.UnirisConfig{})

	_, err := h.LeadKeychainMining(context.TODO(), &api.KeychainLeadRequest{
		EncryptedKeychain: "encrypted data",
		TransactionHash:   "hash",
		ValidatorPeerIPs:  []string{"127.0.0.1"},
	})
	assert.Nil(t, err)

	keychain, _ := db.FindLastKeychain("address")
	assert.NotNil(t, keychain)
}

/*
Scenario: Lead mining of ID transaction
	Given a lead ID request
	When I mine the transaction
	Then I got no error and the data is stored
*/
func TestLeadIDMining(t *testing.T) {
	db := mockstorage.NewDatabase()
	lockSrv := lock.NewService(db)
	notifier := mockNotifier{}
	poolF := mockPoolFinder{}
	cli := mocktransport.NewExternalClient(db)
	poolR := mocktransport.NewPoolRequester(cli)
	emLister := emListing.NewService(db)
	aiClient := mocktransport.NewAIClient()

	txMiners := map[mining.TransactionType]mining.TransactionMiner{
		mining.IDTransaction: accountMining.NewIDMiner(mockcrypto.NewSigner(), mockcrypto.NewHasher()),
	}

	mineSrv := mining.NewService(aiClient, notifier, poolF, poolR, mockcrypto.NewSigner(), emLister, system.UnirisConfig{
		SharedKeys: system.SharedKeys{
			RobotPublicKey: "robotkey",
		},
	}, txMiners)

	accLister := accountListing.NewService(db)
	accAdder := accountadding.NewService(aiClient, db, accLister, mockcrypto.NewSigner(), mockcrypto.NewHasher())

	srv := Services{accAdd: accAdder, lock: lockSrv, mining: mineSrv}
	crypto := Crypto{
		decrypter: mockcrypto.NewDecrypter(),
		signer:    mockcrypto.NewSigner(),
		hasher:    mockcrypto.NewHasher(),
	}
	h := NewExternalServerHandler(srv, crypto, system.UnirisConfig{})

	_, err := h.LeadIDMining(context.TODO(), &api.IDLeadRequest{
		EncryptedID:      "encrypted data",
		TransactionHash:  "hash",
		ValidatorPeerIPs: []string{"127.0.0.1"},
	})
	assert.Nil(t, err)

	id, _ := db.FindID("hash")
	assert.NotNil(t, id)
}

/*
Scenario: Lock a transaction
	Given a transaction to lock
	When I want to lock it
	Then the lock is stored
*/
func TestLockTransaction(t *testing.T) {
	db := mockstorage.NewDatabase()
	lockSrv := lock.NewService(db)

	srv := Services{lock: lockSrv}
	crypto := Crypto{
		decrypter: mockcrypto.NewDecrypter(),
		signer:    mockcrypto.NewSigner(),
		hasher:    mockcrypto.NewHasher(),
	}
	h := NewExternalServerHandler(srv, crypto, system.UnirisConfig{})

	ack, err := h.LockTransaction(context.TODO(), &api.LockRequest{
		MasterRobotKey:  "robotkey",
		Signature:       "sig",
		Address:         "address",
		TransactionHash: "hash",
	})
	assert.Nil(t, err)
	assert.NotNil(t, ack)
	assert.Equal(t, "sig", ack.Signature)

	assert.True(t, db.ContainsLock(lock.TransactionLock{
		Address:        "address",
		MasterRobotKey: "robotkey",
		TxHash:         "hash",
	}))
}

/*
Scenario: Unlock a transaction
	Given a locked transaction
	When I want to unlock it
	Then the lock is removed
*/
func TestUnlockTransaction(t *testing.T) {
	db := mockstorage.NewDatabase()
	lockSrv := lock.NewService(db)

	srv := Services{lock: lockSrv}
	crypto := Crypto{
		decrypter: mockcrypto.NewDecrypter(),
		signer:    mockcrypto.NewSigner(),
		hasher:    mockcrypto.NewHasher(),
	}
	h := NewExternalServerHandler(srv, crypto, system.UnirisConfig{})

	_, err := h.LockTransaction(context.TODO(), &api.LockRequest{
		MasterRobotKey:  "robotkey",
		Signature:       "sig",
		Address:         "address",
		TransactionHash: "hash",
	})
	assert.Nil(t, err)
	assert.True(t, db.ContainsLock(lock.TransactionLock{
		Address:        "address",
		MasterRobotKey: "robotkey",
		TxHash:         "hash",
	}))

	ack, err := h.UnlockTransaction(context.TODO(), &api.LockRequest{
		MasterRobotKey:  "robotkey",
		Signature:       "sig",
		Address:         "address",
		TransactionHash: "hash",
	})
	assert.Nil(t, err)
	assert.NotNil(t, ack)
	assert.Equal(t, "sig", ack.Signature)

	assert.False(t, db.ContainsLock(lock.TransactionLock{
		Address:        "address",
		MasterRobotKey: "robotkey",
		TxHash:         "hash",
	}))

}

/*
Scenario: Validate keychain as slave
	Given keychain transaction
	When I want to validate it
	Then I get a validation
*/
func TestValidateKeychain(t *testing.T) {
	txMiners := map[mining.TransactionType]mining.TransactionMiner{
		mining.KeychainTransaction: accountMining.NewKeychainMiner(mockcrypto.NewSigner(), mockcrypto.NewHasher(), nil),
	}

	mineSrv := mining.NewService(nil, nil, nil, nil, mockcrypto.NewSigner(), nil, system.UnirisConfig{
		SharedKeys: system.SharedKeys{
			RobotPublicKey: "robotkey",
		},
	}, txMiners)

	services := Services{mining: mineSrv}
	crypto := Crypto{
		decrypter: mockcrypto.NewDecrypter(),
		signer:    mockcrypto.NewSigner(),
	}
	h := NewExternalServerHandler(services, crypto, system.UnirisConfig{})

	valid, err := h.ValidateKeychain(context.TODO(), &api.KeychainValidationRequest{
		Data: &api.Keychain{
			EncryptedAddrByRobot: "encrypted addr",
			EncryptedWallet:      "cipher wallet",
			IDPublicKey:          "pubk",
			IDSignature:          "sig",
			EmitterSignature:     "sig",
		},
		TransactionHash: "hash",
	})
	assert.Nil(t, err)
	assert.Equal(t, api.Validation_OK, valid.Validation.Status)
	assert.Equal(t, "sig", valid.Validation.Signature)
	assert.Equal(t, "robotkey", valid.Validation.PublicKey)
	assert.Equal(t, "sig", valid.Signature)
}

/*
Scenario: Validate ID as slave
	Given ID transaction
	When I want to validate it
	Then I get a validation
*/
func TestValidateID(t *testing.T) {

	txMiners := map[mining.TransactionType]mining.TransactionMiner{
		mining.IDTransaction: accountMining.NewIDMiner(mockcrypto.NewSigner(), mockcrypto.NewHasher()),
	}

	mineSrv := mining.NewService(nil, nil, nil, nil, mockcrypto.NewSigner(), nil, system.UnirisConfig{
		SharedKeys: system.SharedKeys{
			RobotPublicKey: "robotkey",
		},
	}, txMiners)

	services := Services{mining: mineSrv}
	crypto := Crypto{
		decrypter: mockcrypto.NewDecrypter(),
		signer:    mockcrypto.NewSigner(),
		hasher:    mockcrypto.NewHasher(),
	}
	h := NewExternalServerHandler(services, crypto, system.UnirisConfig{})

	valid, err := h.ValidateID(context.TODO(), &api.IDValidationRequest{
		Data: &api.ID{
			EncryptedAddrByRobot: "encrypted addr",
			EncryptedAddrByID:    "encrypted addr",
			EncryptedAESKey:      "cipher aes",
			PublicKey:            "pubk",
			IDSignature:          "sig",
			EmitterSignature:     "sig",
		},
		TransactionHash: "hash",
	})
	assert.Nil(t, err)
	assert.Equal(t, api.Validation_OK, valid.Validation.Status)
	assert.Equal(t, "sig", valid.Validation.Signature)
	assert.Equal(t, "robotkey", valid.Validation.PublicKey)
	assert.Equal(t, "sig", valid.Signature)
}

/*
Scenario: Store keychain transaction
	Given a keychain transaction
	When I want to store it
	Then I get retrieve it in the db
*/
func TestStoreKeychain(t *testing.T) {
	db := mockstorage.NewDatabase()
	accLister := accountListing.NewService(db)
	aiClient := mocktransport.NewAIClient()
	accAdder := accountadding.NewService(aiClient, db, accLister, mockcrypto.NewSigner(), mockcrypto.NewHasher())

	services := Services{accAdd: accAdder}
	crypto := Crypto{
		decrypter: mockcrypto.NewDecrypter(),
		signer:    mockcrypto.NewSigner(),
		hasher:    mockcrypto.NewHasher(),
	}
	h := NewExternalServerHandler(services, crypto, system.UnirisConfig{})

	ack, err := h.StoreKeychain(context.TODO(), &api.KeychainStorageRequest{
		Data: &api.Keychain{
			EncryptedWallet:      "encrypted addr",
			EncryptedAddrByRobot: "encrypted addr",
			IDPublicKey:          "pubk",
			IDSignature:          "sig",
			EmitterSignature:     "sig",
		},
		Endorsement: &api.Endorsement{
			LastTransactionHash: "",
			TransactionHash:     "hash",
			MasterValidation: &api.MasterValidation{
				LastTransactionMiners: []string{"hash", "hash"},
				ProofOfWorkKey:        "robotkey",
				ProofOfWorkValidation: &api.Validation{
					PublicKey: "robotkey",
					Signature: "sig",
					Status:    api.Validation_OK,
					Timestamp: time.Now().Unix(),
				},
			},
			Validations: []*api.Validation{
				&api.Validation{
					PublicKey: "robotkey",
					Signature: "sig",
					Status:    api.Validation_OK,
					Timestamp: time.Now().Unix(),
				},
			},
		},
	})
	assert.Nil(t, err)
	assert.NotNil(t, ack)
	assert.Equal(t, "sig", ack.Signature)

	keychain, err := db.FindLastKeychain("hash")
	assert.Nil(t, err)
	assert.NotNil(t, keychain.Endorsement())
	assert.Equal(t, mining.ValidationOK, keychain.Endorsement().MasterValidation().ProofOfWorkValidation().Status())
}

/*
Scenario: Store ID transaction
	Given a ID transaction
	When I want to store it
	Then I get retrieve it in the db
*/
func TestStoreID(t *testing.T) {
	db := mockstorage.NewDatabase()

	accLister := accountListing.NewService(db)
	aiClient := mocktransport.NewAIClient()
	accAdder := accountadding.NewService(aiClient, db, accLister, mockcrypto.NewSigner(), mockcrypto.NewHasher())

	services := Services{accAdd: accAdder}
	crypto := Crypto{
		decrypter: mockcrypto.NewDecrypter(),
		signer:    mockcrypto.NewSigner(),
		hasher:    mockcrypto.NewHasher(),
	}
	h := NewExternalServerHandler(services, crypto, system.UnirisConfig{})

	ack, err := h.StoreID(context.TODO(), &api.IDStorageRequest{
		Data: &api.ID{
			EncryptedAESKey:      "encrypted aes key",
			Hash:                 "hash",
			EncryptedAddrByRobot: "encrypted addr",
			PublicKey:            "pubk",
			IDSignature:          "sig",
			EmitterSignature:     "sig",
		},
		Endorsement: &api.Endorsement{
			LastTransactionHash: "",
			TransactionHash:     "hash",
			MasterValidation: &api.MasterValidation{
				LastTransactionMiners: []string{"hash", "hash"},
				ProofOfWorkKey:        "robotkey",
				ProofOfWorkValidation: &api.Validation{
					PublicKey: "robotkey",
					Signature: "sig",
					Status:    api.Validation_OK,
					Timestamp: time.Now().Unix(),
				},
			},
			Validations: []*api.Validation{
				&api.Validation{
					PublicKey: "robotkey",
					Signature: "sig",
					Status:    api.Validation_OK,
					Timestamp: time.Now().Unix(),
				},
			},
		},
	})

	assert.Nil(t, err)
	assert.NotNil(t, ack)
	assert.Equal(t, "sig", ack.Signature)

	id, err := db.FindID("hash")
	assert.Nil(t, err)
	assert.Equal(t, "encrypted aes key", id.EncryptedAESKey())
	assert.NotNil(t, id.Endorsement())
	assert.Equal(t, mining.ValidationOK, id.Endorsement().MasterValidation().ProofOfWorkValidation().Status())
}

type mockPoolFinder struct{}

func (p mockPoolFinder) FindLastValidationPool(addr string) (datamining.Pool, error) {
	return datamining.NewPool(datamining.Peer{
		IP:        net.ParseIP("127.0.0.1"),
		PublicKey: "key",
	}), nil
}

func (p mockPoolFinder) FindValidationPool() (datamining.Pool, error) {
	return datamining.NewPool(datamining.Peer{
		IP:        net.ParseIP("127.0.0.1"),
		PublicKey: "key",
	}), nil
}

func (p mockPoolFinder) FindStoragePool(addr string) (datamining.Pool, error) {
	return datamining.NewPool(datamining.Peer{
		IP:        net.ParseIP("127.0.0.1"),
		PublicKey: "key",
	}), nil
}

type mockNotifier struct{}

func (n mockNotifier) NotifyTransactionStatus(tx string, status mining.TransactionStatus) error {
	log.Printf("Transaction %s with status %s", tx, status.String())
	return nil
}

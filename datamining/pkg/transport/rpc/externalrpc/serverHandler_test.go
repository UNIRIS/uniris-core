package externalrpc

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/golang/protobuf/ptypes/any"

	datamining "github.com/uniris/uniris-core/datamining/pkg"
	"github.com/uniris/uniris-core/datamining/pkg/adding"
	"github.com/uniris/uniris-core/datamining/pkg/mining/slave"
	"github.com/uniris/uniris-core/datamining/pkg/system"

	"github.com/stretchr/testify/assert"

	api "github.com/uniris/uniris-core/datamining/api/protobuf-spec"
	"github.com/uniris/uniris-core/datamining/pkg/listing"
	"github.com/uniris/uniris-core/datamining/pkg/locking"

	mockcrypto "github.com/uniris/uniris-core/datamining/pkg/crypto/mock"
	mockstorage "github.com/uniris/uniris-core/datamining/pkg/storage/mock"
)

/*
Scenario: Lock a transaction
	Given a lock request
	When I want to lock a transaction
	Then the transaction is stored as locked
*/
func TestLockTransaction(t *testing.T) {

	locker := mockstorage.NewTransactionLocker()
	srv := externalSrvHandler{
		lock: locking.NewService(locker),
	}

	_, err := srv.LockTransaction(context.TODO(), &api.LockRequest{
		MasterRobotKey:  "robotkey",
		TransactionHash: "txhash",
	})
	assert.Nil(t, err)

	assert.True(t, locker.ContainsLock(locking.TransactionLock{
		MasterRobotKey: "robotkey",
		TxHash:         "txhash",
	}))
}

/*
Scenario: Lock a already locked transaction
	Given a locked transaction
	When I want to relock the transaction
	Then I get an error
*/
func TestAlreadyLockTransaction(t *testing.T) {

	locker := mockstorage.NewTransactionLocker()
	srv := externalSrvHandler{
		lock: locking.NewService(locker),
	}

	srv.LockTransaction(context.TODO(), &api.LockRequest{
		MasterRobotKey:  "robotkey",
		TransactionHash: "txhash",
	})

	_, err := srv.LockTransaction(context.TODO(), &api.LockRequest{
		MasterRobotKey:  "robotkey",
		TransactionHash: "txhash",
	})
	assert.Equal(t, err, locking.ErrLockExisting)
}

/*
Scenario: UNLock a transaction
	Given locked transaction
	When I want to unlock this transaction
	Then the transaction is not stored as locked
*/
func TestUnLockTransaction(t *testing.T) {

	locker := mockstorage.NewTransactionLocker()
	srv := externalSrvHandler{
		lock: locking.NewService(locker),
	}

	srv.LockTransaction(context.TODO(), &api.LockRequest{
		MasterRobotKey:  "robotkey",
		TransactionHash: "txhash",
	})

	_, err := srv.UnlockTransaction(context.TODO(), &api.LockRequest{
		MasterRobotKey:  "robotkey",
		TransactionHash: "txhash",
	})
	assert.Nil(t, err)
	assert.False(t, locker.ContainsLock(locking.TransactionLock{
		MasterRobotKey: "robotkey",
		TxHash:         "txhash",
	}))
}

/*
Scenario: Validate an incoming transaction
	Given a validation request and a valid transaction
	When I want to valid it
	Then I get a validation with OK status
*/
func TestValidTransaction(t *testing.T) {

	db := mockstorage.NewDatabase()
	l := listing.NewService(db)
	s := slave.NewService(
		mockcrypto.NewSigner(),
		"robotkey",
		"robotpvKey",
	)

	bd := &datamining.BioData{
		BiodPubk:   "biopubkey",
		PersonPubk: "personpubkey",
		Sigs: datamining.Signatures{
			BiodSig:   "fake sig",
			PersonSig: "fake sig",
		},
	}

	b, _ := json.Marshal(bd)

	srv := NewExternalServerHandler(l, nil, s, nil, "", system.DataMininingErrors{})
	resp, err := srv.Validate(context.TODO(), &api.ValidationRequest{
		TransactionType: api.TransactionType_CreateBio,
		Data:            &any.Any{Value: b},
	})

	assert.Nil(t, err)
	assert.Equal(t, resp.Validation.Status, api.Validation_OK)

}

/*
Scenario: Store an biometric creation transaction
	Given a biometric storage request
	When I want to store it
	Then I the data is stored
*/
func TestStoreBioData(t *testing.T) {

	db := mockstorage.NewDatabase()
	l := listing.NewService(db)
	a := adding.NewService(db)

	bd := &datamining.BioData{
		PersonHash: "hash",
		BiodPubk:   "biopubkey",
		PersonPubk: "personpubkey",
		Sigs: datamining.Signatures{
			BiodSig:   "fake sig",
			PersonSig: "fake sig",
		},
	}

	b := datamining.NewBiometric(bd, nil)

	bytes, _ := json.Marshal(b)

	srv := NewExternalServerHandler(l, a, nil, nil, "", system.DataMininingErrors{})
	_, err := srv.Store(context.TODO(), &api.StorageRequest{
		TransactionType: api.TransactionType_CreateBio,
		Data:            &any.Any{Value: bytes},
	})

	assert.Nil(t, err)

	bio, _ := db.FindBiometric("hash")
	assert.NotNil(t, bio)
}

/*
Scenario: Store an keychain creation transaction
	Given a keychain storage request
	When I want to store it
	Then I the data is stored
*/
func TestStoreKeychainData(t *testing.T) {

	db := mockstorage.NewDatabase()
	l := listing.NewService(db)
	a := adding.NewService(db)

	kcData := &datamining.KeyChainData{
		CipherAddrRobot: "addr",
		CipherWallet:    "addr",
		WalletAddr:      "addr",
		BiodPubk:        "biopubkey",
		PersonPubk:      "personpubkey",
		Sigs: datamining.Signatures{
			BiodSig:   "fake sig",
			PersonSig: "fake sig",
		},
	}

	kc := datamining.NewKeychain(kcData, nil, "")

	bytes, _ := json.Marshal(kc)

	srv := NewExternalServerHandler(l, a, nil, nil, "", system.DataMininingErrors{})
	_, err := srv.Store(context.TODO(), &api.StorageRequest{
		TransactionType: api.TransactionType_CreateKeychain,
		Data:            &any.Any{Value: bytes},
	})

	assert.Nil(t, err)

	keychain, _ := db.FindKeychain("addr")
	assert.NotNil(t, keychain)
}

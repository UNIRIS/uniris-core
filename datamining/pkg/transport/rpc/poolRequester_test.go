package rpc

import (
	"net"
	"testing"
	"time"

	"github.com/uniris/uniris-core/datamining/pkg/lock"
	"github.com/uniris/uniris-core/datamining/pkg/mining"

	"github.com/stretchr/testify/assert"

	"github.com/uniris/uniris-core/datamining/pkg"

	"github.com/uniris/uniris-core/datamining/pkg/account"
	mockcrypto "github.com/uniris/uniris-core/datamining/pkg/crypto/mock"
	mockstorage "github.com/uniris/uniris-core/datamining/pkg/storage/mock"
	"github.com/uniris/uniris-core/datamining/pkg/system"
	"github.com/uniris/uniris-core/datamining/pkg/transport/mock"
)

/*
Scenario: Request biometric data
	Given an encypted person hash and storage pool
	When I want to retrieve the biometric data related to this hash
	Then it launches a pool of goroutines and requests the biometric data without error
*/
func TestRequestBiometric(t *testing.T) {
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

	cli := mock.NewExternalClient(db)

	pr := NewPoolRequester(cli, conf, crypto)
	pool := datamining.NewPool(
		datamining.Peer{IP: net.ParseIP("127.0.0.1")},
		datamining.Peer{IP: net.ParseIP("127.0.0.1")})
	bio, err := pr.RequestBiometric(pool, "enc hash")
	assert.Nil(t, err)
	assert.NotNil(t, bio)

	assert.Equal(t, "hash", bio.PersonHash())
	assert.Equal(t, "enc addr", bio.CipherAddrRobot())
}

/*
Scenario: Request keychain data
	Given an encypted address hash and storage pool
	When I want to retrieve the keychain data related to this address
	Then it launches a pool of goroutines and requests the keychain without error
*/
func TestRequestKeychain(t *testing.T) {
	conf := system.UnirisConfig{}
	crypto := Crypto{
		decrypter: mockcrypto.NewDecrypter(),
		signer:    mockcrypto.NewSigner(),
		hasher:    mockcrypto.NewHasher(),
	}

	db := mockstorage.NewDatabase()

	db.StoreKeychain(
		account.NewKeychain(
			"hash",
			account.NewKeychainData("enc addr", "enc wallet", "pub", "pub", account.NewSignatures("sig", "sig")),
			nil,
		),
	)

	cli := mock.NewExternalClient(db)

	pr := NewPoolRequester(cli, conf, crypto)
	pool := datamining.NewPool(
		datamining.Peer{IP: net.ParseIP("127.0.0.1")},
		datamining.Peer{IP: net.ParseIP("127.0.0.1")})
	kc, err := pr.RequestKeychain(pool, "enc hash")
	assert.Nil(t, err)
	assert.NotNil(t, kc)

	assert.Equal(t, "hash", kc.Address())
	assert.Equal(t, "enc wallet", kc.CipherWallet())
	assert.Equal(t, "enc addr", kc.CipherAddrRobot())
}

/*
Scenario: Request lock transaction
	Given an lock transaction and lock pool
	When I want to lock the transaction
	Then it launches a pool of goroutines and lock the transaction without error
*/
func TestRequestLock(t *testing.T) {
	conf := system.UnirisConfig{}
	crypto := Crypto{
		decrypter: mockcrypto.NewDecrypter(),
		signer:    mockcrypto.NewSigner(),
		hasher:    mockcrypto.NewHasher(),
	}

	db := mockstorage.NewDatabase()

	cli := mock.NewExternalClient(db)

	pr := NewPoolRequester(cli, conf, crypto)
	pool := datamining.NewPool(
		datamining.Peer{IP: net.ParseIP("127.0.0.1")},
		datamining.Peer{IP: net.ParseIP("127.0.0.1")})
	err := pr.RequestLock(pool, lock.TransactionLock{
		Address:        "address",
		MasterRobotKey: "robotkey",
		TxHash:         "txhash",
	})
	assert.Nil(t, err)

	assert.True(t, db.ContainsLock(lock.TransactionLock{
		Address:        "address",
		MasterRobotKey: "robotkey",
		TxHash:         "txhash",
	}))
}

/*
Scenario: Request unlock transaction
	Given an lock transaction and lock pool
	When I want to unlock the transaction
	Then it launches a pool of goroutines and unlock the transaction without error
*/
func TestRequestUnLock(t *testing.T) {
	conf := system.UnirisConfig{}
	crypto := Crypto{
		decrypter: mockcrypto.NewDecrypter(),
		signer:    mockcrypto.NewSigner(),
		hasher:    mockcrypto.NewHasher(),
	}

	db := mockstorage.NewDatabase()

	cli := mock.NewExternalClient(db)

	pr := NewPoolRequester(cli, conf, crypto)
	pool := datamining.NewPool(
		datamining.Peer{IP: net.ParseIP("127.0.0.1")},
		datamining.Peer{IP: net.ParseIP("127.0.0.1")})
	err := pr.RequestLock(pool, lock.TransactionLock{
		Address:        "address",
		MasterRobotKey: "robotkey",
		TxHash:         "txhash",
	})
	assert.Nil(t, err)

	err = pr.RequestUnlock(pool, lock.TransactionLock{
		Address:        "address",
		MasterRobotKey: "robotkey",
		TxHash:         "txhash",
	})
	assert.Nil(t, err)

	assert.False(t, db.ContainsLock(lock.TransactionLock{
		Address:        "address",
		MasterRobotKey: "robotkey",
		TxHash:         "txhash",
	}))
}

/*
Scenario: Request transaction validations
	Given an transaction hash and data associated
	When I want to validate them
	Then it launches a pool of goroutines and validate the information without error
*/
func TestRequestValidations(t *testing.T) {
	conf := system.UnirisConfig{}
	crypto := Crypto{
		decrypter: mockcrypto.NewDecrypter(),
		signer:    mockcrypto.NewSigner(),
		hasher:    mockcrypto.NewHasher(),
	}

	db := mockstorage.NewDatabase()

	cli := mock.NewExternalClient(db)

	pr := NewPoolRequester(cli, conf, crypto)
	pool := datamining.NewPool(
		datamining.Peer{IP: net.ParseIP("127.0.0.1")},
		datamining.Peer{IP: net.ParseIP("127.0.0.1")})

	keychainData := account.NewKeychainData("enc addr", "enc wallet", "pub", "pub", account.NewSignatures("sig", "sig"))

	valids, err := pr.RequestValidations(pool, "hash", keychainData, mining.KeychainTransaction)
	assert.Nil(t, err)
	assert.NotEmpty(t, valids)
	assert.Equal(t, 2, len(valids))
}

/*
Scenario: Request transaction storages
	Given an transaction hash, data associated and endorsement
	When I want to storage them
	Then it launches a pool of goroutines and store the information without error
*/
func TestRequestStorage(t *testing.T) {
	conf := system.UnirisConfig{}
	crypto := Crypto{
		decrypter: mockcrypto.NewDecrypter(),
		signer:    mockcrypto.NewSigner(),
		hasher:    mockcrypto.NewHasher(),
	}

	db := mockstorage.NewDatabase()

	cli := mock.NewExternalClient(db)

	keychainData := account.NewKeychainData("enc addr", "enc wallet", "pub", "pub", account.NewSignatures("sig", "sig"))
	end := mining.NewEndorsement("", "hash",
		mining.NewMasterValidation([]string{""}, "key", mining.NewValidation(mining.ValidationOK, time.Now(), "pub", "sig")),
		[]mining.Validation{mining.NewValidation(mining.ValidationOK, time.Now(), "pub", "sig")},
	)

	pr := NewPoolRequester(cli, conf, crypto)
	pool := datamining.NewPool(
		datamining.Peer{IP: net.ParseIP("127.0.0.1")},
		datamining.Peer{IP: net.ParseIP("127.0.0.1")})
	err := pr.RequestStorage(pool, keychainData, end, mining.KeychainTransaction)
	assert.Nil(t, err)

	kc, _ := db.FindLastKeychain("address")
	assert.NotNil(t, kc)
}

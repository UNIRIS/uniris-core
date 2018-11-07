package externalrpc

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/uniris/uniris-core/datamining/pkg"

	"github.com/stretchr/testify/assert"

	api "github.com/uniris/uniris-core/datamining/api/protobuf-spec"
	accountadding "github.com/uniris/uniris-core/datamining/pkg/account/adding"
	accountListing "github.com/uniris/uniris-core/datamining/pkg/account/listing"
	accountMining "github.com/uniris/uniris-core/datamining/pkg/account/mining"
	biodlisting "github.com/uniris/uniris-core/datamining/pkg/biod/listing"
	"github.com/uniris/uniris-core/datamining/pkg/lock"
	"github.com/uniris/uniris-core/datamining/pkg/mining"
	"github.com/uniris/uniris-core/datamining/pkg/mock"
	"github.com/uniris/uniris-core/datamining/pkg/system"
)

/*
Scenario: Lead mining of keychain transaction
	Given a lead keychain request
	When I mine the transaction
	Then I got no error and the data is stored
*/
func TestLeadKeychainMining(t *testing.T) {
	db := mock.NewDatabase()
	lockSrv := lock.NewService(db)
	notifier := mock.NewNotifier()
	poolF := mock.NewPoolFinder()
	poolR := mock.NewPoolRequester(db)
	biodlister := biodlisting.NewService(db)
	accLister := accountListing.NewService(db)

	txMiners := map[mining.TransactionType]mining.TransactionMiner{
		mining.KeychainTransaction: accountMining.NewKeychainMiner(mockSigner{}, mockHasher{}, accLister),
	}

	mineSrv := mining.NewService(notifier, poolF, poolR, mockSigner{}, biodlister, "robotkey", "robotpbKey", txMiners)

	accSrv := accountadding.NewService(db)

	conf := system.UnirisConfig{}

	srvHandler := NewExternalServerHandler(lockSrv, mineSrv, accSrv, mockBiometricDecrypter{}, conf)

	_, err := srvHandler.LeadKeychainMining(context.TODO(), &api.KeychainLeadRequest{
		EncryptedKeychainData: "encrypted data",
		SignatureKeychainData: &api.Signature{
			Biod:   "sig",
			Person: "sig",
		},
		TransactionHash:  "hash",
		ValidatorPeerIPs: []string{"127.0.0.1"},
	})
	assert.Nil(t, err)

	keychain, _ := db.FindLastKeychain("address")
	assert.NotNil(t, keychain)
}

/*
Scenario: Lead mining of biometric transaction
	Given a lead biometric request
	When I mine the transaction
	Then I got no error and the data is stored
*/
func TestLeadBiometricMining(t *testing.T) {
	db := mock.NewDatabase()
	lockSrv := lock.NewService(db)
	notifier := mock.NewNotifier()
	poolF := mock.NewPoolFinder()
	poolR := mock.NewPoolRequester(db)
	biodlister := biodlisting.NewService(db)

	txMiners := map[mining.TransactionType]mining.TransactionMiner{
		mining.BiometricTransaction: accountMining.NewBiometricMiner(mockSigner{}, mockHasher{}),
	}

	mineSrv := mining.NewService(notifier, poolF, poolR, mockSigner{}, biodlister, "robotkey", "robotpbKey", txMiners)

	accSrv := accountadding.NewService(db)

	conf := system.UnirisConfig{}

	srvHandler := NewExternalServerHandler(lockSrv, mineSrv, accSrv, mockBiometricDecrypter{}, conf)

	_, err := srvHandler.LeadBiometricMining(context.TODO(), &api.BiometricLeadRequest{
		EncryptedBioData: "encrypted data",
		SignatureBioData: &api.Signature{
			Biod:   "sig",
			Person: "sig",
		},
		TransactionHash:  "hash",
		ValidatorPeerIPs: []string{"127.0.0.1"},
	})
	assert.Nil(t, err)

	biometric, _ := db.FindBiometric("person hash")
	assert.NotNil(t, biometric)
}

/*
Scenario: Lock a transaction
	Given a transaction to lock
	When I want to lock it
	Then the lock is stored
*/
func TestLockTransaction(t *testing.T) {
	db := mock.NewDatabase()
	lockSrv := lock.NewService(db)

	srvHandler := externalSrvHandler{
		lock: lockSrv,
	}

	_, err := srvHandler.LockTransaction(context.TODO(), &api.LockRequest{
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
}

/*
Scenario: Unlock a transaction
	Given a locked transaction
	When I want to unlock it
	Then the lock is removed
*/
func TestUnlockTransaction(t *testing.T) {
	db := mock.NewDatabase()
	lockSrv := lock.NewService(db)

	srvHandler := externalSrvHandler{
		lock: lockSrv,
	}

	_, err := srvHandler.LockTransaction(context.TODO(), &api.LockRequest{
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

	_, err = srvHandler.UnlockTransaction(context.TODO(), &api.LockRequest{
		MasterRobotKey:  "robotkey",
		Signature:       "sig",
		Address:         "address",
		TransactionHash: "hash",
	})
	assert.Nil(t, err)

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
		mining.KeychainTransaction: accountMining.NewKeychainMiner(mockSigner{}, mockHasher{}, nil),
	}

	mineSrv := mining.NewService(nil, nil, nil, mockSigner{}, nil, "robotkey", "robotpbKey", txMiners)

	srvHandler := NewExternalServerHandler(nil, mineSrv, nil, mockBiometricDecrypter{}, system.UnirisConfig{})

	valid, err := srvHandler.ValidateKeychain(context.TODO(), &api.KeychainValidationRequest{
		Data: &api.KeychainData{
			BiodPubk:        "pubk",
			CipherAddrRobot: "encrypted addr",
			CipherWallet:    "cipher wallet",
			PersonPubk:      "pubk",
			Signature: &api.Signature{
				Biod:   "sig",
				Person: "sig",
			},
		},
		TransactionHash: "hash",
	})
	assert.Nil(t, err)
	assert.Equal(t, api.Validation_OK, valid.Validation.Status)
	assert.Equal(t, "sign", valid.Validation.Signature)
	assert.Equal(t, "robotkey", valid.Validation.PublicKey)
}

/*
Scenario: Validate biometric as slave
	Given biometric transaction
	When I want to validate it
	Then I get a validation
*/
func TestValidateBiometric(t *testing.T) {
	txMiners := map[mining.TransactionType]mining.TransactionMiner{
		mining.BiometricTransaction: accountMining.NewBiometricMiner(mockSigner{}, mockHasher{}),
	}

	mineSrv := mining.NewService(nil, nil, nil, mockSigner{}, nil, "robotkey", "robotpbKey", txMiners)

	srvHandler := NewExternalServerHandler(nil, mineSrv, nil, mockBiometricDecrypter{}, system.UnirisConfig{})

	valid, err := srvHandler.ValidateBiometric(context.TODO(), &api.BiometricValidationRequest{
		Data: &api.BiometricData{
			BiodPubk:        "pubk",
			CipherAddrRobot: "encrypted addr",
			CipherAddrBio:   "encrypted addr",
			CipherAESKey:    "cipher aes",
			PersonPubk:      "pubk",
			Signature: &api.Signature{
				Biod:   "sig",
				Person: "sig",
			},
		},
		TransactionHash: "hash",
	})
	assert.Nil(t, err)
	assert.Equal(t, api.Validation_OK, valid.Validation.Status)
	assert.Equal(t, "sign", valid.Validation.Signature)
	assert.Equal(t, "robotkey", valid.Validation.PublicKey)
}

/*
Scenario: Store keychain transaction
	Given a keychain transaction
	When I want to store it
	Then I get retrieve it in the db
*/
func TestStoreKeychain(t *testing.T) {
	db := mock.NewDatabase()
	accAdder := accountadding.NewService(db)
	srvHandler := NewExternalServerHandler(nil, nil, accAdder, mockKeychainDecrypter{}, system.UnirisConfig{})

	_, err := srvHandler.StoreKeychain(context.TODO(), &api.KeychainStorageRequest{
		Data: &api.KeychainData{
			CipherWallet:    "encrypted addr",
			BiodPubk:        "pubk",
			CipherAddrRobot: "encrypted addr",
			PersonPubk:      "pubk",
			Signature: &api.Signature{
				Biod:   "sig",
				Person: "sig",
			},
		},
		Endorsement: &api.Endorsement{
			LastTransactionHash: "",
			TransactionHash:     "hash",
			MasterValidation: &api.MasterValidation{
				LastTransactionMiners: []string{"hash", "hash"},
				ProofOfWorkRobotKey:   "robot key",
				ProofOfWorkValidation: &api.Validation{
					PublicKey: "robotkey",
					Signature: "sig",
					Status:    api.Validation_OK,
					Timestamp: time.Now().Unix(),
				},
			},
		},
		TransactionHash: "hash",
	})

	assert.Nil(t, err)
	keychain, err := db.FindLastKeychain("address")
	assert.Nil(t, err)
	assert.NotNil(t, keychain.Endorsement())
	assert.Equal(t, datamining.ValidationOK, keychain.Endorsement().MasterValidation().ProofOfWorkValidation().Status())
}

/*
Scenario: Store biometric transaction
	Given a biometric transaction
	When I want to store it
	Then I get retrieve it in the db
*/
func TestStoreBiometric(t *testing.T) {
	db := mock.NewDatabase()
	accAdder := accountadding.NewService(db)
	srvHandler := NewExternalServerHandler(nil, nil, accAdder, mockBiometricDecrypter{}, system.UnirisConfig{})

	_, err := srvHandler.StoreBiometric(context.TODO(), &api.BiometricStorageRequest{
		Data: &api.BiometricData{
			CipherAESKey:    "encrypted aes key",
			PersonHash:      "hash",
			BiodPubk:        "pubk",
			CipherAddrRobot: "encrypted addr",
			PersonPubk:      "pubk",
			Signature: &api.Signature{
				Biod:   "sig",
				Person: "sig",
			},
		},
		Endorsement: &api.Endorsement{
			LastTransactionHash: "",
			TransactionHash:     "hash",
			MasterValidation: &api.MasterValidation{
				LastTransactionMiners: []string{"hash", "hash"},
				ProofOfWorkRobotKey:   "robot key",
				ProofOfWorkValidation: &api.Validation{
					PublicKey: "robotkey",
					Signature: "sig",
					Status:    api.Validation_OK,
					Timestamp: time.Now().Unix(),
				},
			},
		},
		TransactionHash: "hash",
	})

	assert.Nil(t, err)
	biometric, err := db.FindBiometric("hash")
	assert.Nil(t, err)
	assert.Equal(t, "encrypted aes key", biometric.CipherAESKey())
	assert.NotNil(t, biometric.Endorsement())
	assert.Equal(t, datamining.ValidationOK, biometric.Endorsement().MasterValidation().ProofOfWorkValidation().Status())
}

type mockBiometricDecrypter struct{}

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

type mockSigner struct{}

func (s mockSigner) CheckTransactionSignature(pubk string, txHash string, sig string) error {
	return nil
}

func (s mockSigner) CheckBiometricSignature(pubk string, data accountMining.UnsignedBiometricData, sig string) error {
	return nil
}

func (s mockSigner) CheckKeychainSignature(pubk string, data accountMining.UnsignedKeychainData, sig string) error {
	return nil
}

func (s mockSigner) SignValidation(v mining.UnsignedValidation, pvKey string) (string, error) {
	return "sign", nil
}

func (s mockSigner) SignLock(txLock lock.TransactionLock, pvKey string) (string, error) {
	return "sign", nil
}

type mockHasher struct{}

func (h mockHasher) HashUnsignedBiometricData(data accountMining.UnsignedBiometricData) (string, error) {
	return "hash", nil
}

func (h mockHasher) HashUnsignedKeychainData(data accountMining.UnsignedKeychainData) (string, error) {
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

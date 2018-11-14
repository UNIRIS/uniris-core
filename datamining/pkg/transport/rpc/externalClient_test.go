package rpc

import (
	"fmt"
	"log"
	"net"
	"testing"
	"time"

	"github.com/uniris/uniris-core/datamining/pkg/lock"
	"github.com/uniris/uniris-core/datamining/pkg/mining"
	mocktransport "github.com/uniris/uniris-core/datamining/pkg/transport/mock"

	"github.com/stretchr/testify/assert"

	"github.com/uniris/uniris-core/datamining/pkg/account"
	mockcrypto "github.com/uniris/uniris-core/datamining/pkg/crypto/mock"
	"github.com/uniris/uniris-core/datamining/pkg/system"

	api "github.com/uniris/uniris-core/datamining/api/protobuf-spec"
	accountAdding "github.com/uniris/uniris-core/datamining/pkg/account/adding"
	accountListing "github.com/uniris/uniris-core/datamining/pkg/account/listing"
	accountMining "github.com/uniris/uniris-core/datamining/pkg/account/mining"
	mockstorage "github.com/uniris/uniris-core/datamining/pkg/storage/mock"
	"google.golang.org/grpc"
)

/*
Scenario: Call RequestBiometric GRPC endpoint
	Given a biometric stored and a encrypted person hash
	When I want to retrieve it, the client call the GRPC endpoint
	Then I retrieve the biometric data stored
*/
func TestRequestBiometricClient(t *testing.T) {

	db := mockstorage.NewDatabase()
	accLister := accountListing.NewService(db)

	crypto := Crypto{
		decrypter: mockcrypto.NewDecrypter(),
		signer:    mockcrypto.NewSigner(),
	}

	port := 2000
	conf := system.UnirisConfig{
		Datamining: system.DataMiningConfiguration{
			ExternalPort: port,
		},
	}

	grpcServer := grpc.NewServer()
	defer grpcServer.Stop()

	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
		assert.Nil(t, err)

		services := Services{
			accLister: accLister,
		}

		handler := NewExternalServerHandler(services, crypto, conf)
		api.RegisterExternalServer(grpcServer, handler)
		log.Printf("External grpc Server listening on 127.0.0.1:%d", port)
		err = grpcServer.Serve(lis)
		assert.Nil(t, err)
	}()

	time.Sleep(1 * time.Second)

	db.StoreBiometric(
		account.NewBiometric(
			account.NewBiometricData("hash", "enc addr", "enc addr", "enc aes key", "pub", "pub", account.NewSignatures("sig", "sig")),
			mining.NewEndorsement("", "hash",
				mining.NewMasterValidation([]string{"hash"}, "key", mining.NewValidation(mining.ValidationOK, time.Now(), "pub", "sig")),
				[]mining.Validation{mining.NewValidation(mining.ValidationOK, time.Now(), "pub", "sig")}),
		),
	)

	cli := NewExternalClient(crypto, conf)
	bio, err := cli.RequestBiometric("127.0.0.1", "hash")
	assert.Nil(t, err)
	assert.NotNil(t, bio)
	assert.Equal(t, "enc aes key", bio.CipherAESKey())
}

/*
Scenario: Call RequestKeychain GRPC endpoint
	Given a keychain stored and a encrypted address
	When I want to retrieve it, the client call the GRPC endpoint
	Then I retrieve the biometric data stored
*/
func TestRequestKeychainClient(t *testing.T) {

	db := mockstorage.NewDatabase()
	accLister := accountListing.NewService(db)

	crypto := Crypto{
		decrypter: mockcrypto.NewDecrypter(),
		signer:    mockcrypto.NewSigner(),
	}

	port := 2001
	conf := system.UnirisConfig{
		Datamining: system.DataMiningConfiguration{
			ExternalPort: port,
		},
	}

	grpcServer := grpc.NewServer()
	defer grpcServer.Stop()

	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
		assert.Nil(t, err)

		services := Services{
			accLister: accLister,
		}

		handler := NewExternalServerHandler(services, crypto, conf)
		api.RegisterExternalServer(grpcServer, handler)
		log.Printf("External grpc Server listening on 127.0.0.1:%d", port)
		err = grpcServer.Serve(lis)
		assert.Nil(t, err)

	}()

	time.Sleep(1 * time.Second)

	db.StoreKeychain(
		account.NewKeychain(
			"hash",
			account.NewKeychainData("enc addr", "enc wallet", "pub", "pub", account.NewSignatures("sig", "sig")),
			mining.NewEndorsement("", "hash",
				mining.NewMasterValidation([]string{"hash"}, "key", mining.NewValidation(mining.ValidationOK, time.Now(), "pub", "sig")),
				[]mining.Validation{mining.NewValidation(mining.ValidationOK, time.Now(), "pub", "sig")}),
		),
	)

	cli := NewExternalClient(crypto, conf)
	kc, err := cli.RequestKeychain("127.0.0.1", "hash")
	assert.Nil(t, err)
	assert.NotNil(t, kc)
	assert.Equal(t, "enc wallet", kc.CipherWallet())
}

/*
Scenario: Call RequestLock GRPC endpoint
	Given a transaction
	When I want to lock it, the client call the GRPC endpoint
	Then I retrieve the lock is stored
*/
func TestRequestLockClient(t *testing.T) {

	db := mockstorage.NewDatabase()
	locker := lock.NewService(db)

	crypto := Crypto{
		decrypter: mockcrypto.NewDecrypter(),
		signer:    mockcrypto.NewSigner(),
		hasher:    mockcrypto.NewHasher(),
	}

	port := 2002
	conf := system.UnirisConfig{
		Datamining: system.DataMiningConfiguration{
			ExternalPort: port,
		},
	}

	grpcServer := grpc.NewServer()
	defer grpcServer.Stop()

	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
		assert.Nil(t, err)

		services := Services{
			lock: locker,
		}

		handler := NewExternalServerHandler(services, crypto, conf)
		api.RegisterExternalServer(grpcServer, handler)
		log.Printf("External grpc Server listening on 127.0.0.1:%d", port)
		err = grpcServer.Serve(lis)
		assert.Nil(t, err)

	}()

	time.Sleep(1 * time.Second)

	cli := NewExternalClient(crypto, conf)
	err := cli.RequestLock("127.0.0.1", lock.TransactionLock{
		Address:        "address",
		MasterRobotKey: "robotkey",
		TxHash:         "hash",
	})
	assert.Nil(t, err)

	assert.True(t, db.ContainsLock(lock.TransactionLock{
		Address:        "address",
		MasterRobotKey: "robotkey",
		TxHash:         "hash",
	}))
}

/*
Scenario: Call RequestUnLock GRPC endpoint
	Given a locked transaction
	When I want to unlock it, the client call the GRPC endpoint
	Then I retrieve the lock is removed
*/
func TestRequestUnLockClient(t *testing.T) {

	db := mockstorage.NewDatabase()
	locker := lock.NewService(db)

	crypto := Crypto{
		decrypter: mockcrypto.NewDecrypter(),
		signer:    mockcrypto.NewSigner(),
		hasher:    mockcrypto.NewHasher(),
	}

	port := 2003
	conf := system.UnirisConfig{
		Datamining: system.DataMiningConfiguration{
			ExternalPort: port,
		},
	}

	grpcServer := grpc.NewServer()
	defer grpcServer.Stop()

	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
		assert.Nil(t, err)

		services := Services{
			lock: locker,
		}

		handler := NewExternalServerHandler(services, crypto, conf)
		api.RegisterExternalServer(grpcServer, handler)
		log.Printf("External grpc Server listening on 127.0.0.1:%d", port)
		err = grpcServer.Serve(lis)
		assert.Nil(t, err)

	}()

	time.Sleep(1 * time.Second)

	cli := NewExternalClient(crypto, conf)
	err := cli.RequestLock("127.0.0.1", lock.TransactionLock{
		Address:        "address",
		MasterRobotKey: "robotkey",
		TxHash:         "hash",
	})
	assert.Nil(t, err)

	err = cli.RequestUnlock("127.0.0.1", lock.TransactionLock{
		Address:        "address",
		MasterRobotKey: "robotkey",
		TxHash:         "hash",
	})
	assert.Nil(t, err)

	assert.False(t, db.ContainsLock(lock.TransactionLock{
		Address:        "address",
		MasterRobotKey: "robotkey",
		TxHash:         "hash",
	}))
}

/*
Scenario: Call RequestValidation GRPC endpoint for a keychain
	Given a transaction hash, keychain data
	When I want to validate it, the client call the GRPC endpoint
	Then I get the validation processed
*/
func TestRequestKeychainValidationClient(t *testing.T) {

	db := mockstorage.NewDatabase()
	accLister := accountListing.NewService(db)

	txMiners := map[mining.TransactionType]mining.TransactionMiner{
		mining.KeychainTransaction: accountMining.NewKeychainMiner(mockcrypto.NewSigner(), mockcrypto.NewHasher(), accLister),
	}

	crypto := Crypto{
		decrypter: mockcrypto.NewDecrypter(),
		signer:    mockcrypto.NewSigner(),
		hasher:    mockcrypto.NewHasher(),
	}

	port := 2004
	conf := system.UnirisConfig{
		Datamining: system.DataMiningConfiguration{
			ExternalPort: port,
		},
	}

	aiClient := mocktransport.NewAIClient()
	miner := mining.NewService(aiClient, nil, nil, nil, mockcrypto.NewSigner(), nil, conf, txMiners)

	grpcServer := grpc.NewServer()
	defer grpcServer.Stop()

	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
		assert.Nil(t, err)

		services := Services{
			mining: miner,
		}

		handler := NewExternalServerHandler(services, crypto, conf)
		api.RegisterExternalServer(grpcServer, handler)
		log.Printf("External grpc Server listening on 127.0.0.1:%d", port)
		err = grpcServer.Serve(lis)
		assert.Nil(t, err)

	}()

	time.Sleep(1 * time.Second)

	keychainData := account.NewKeychainData("enc addr", "enc wallet", "pub", "pub", account.NewSignatures("sig", "sig"))

	cli := NewExternalClient(crypto, conf)
	valid, err := cli.RequestValidation("127.0.0.1", mining.KeychainTransaction, "hash", keychainData)
	assert.Nil(t, err)
	assert.NotNil(t, valid)
}

/*
Scenario: Call RequestValidation GRPC endpoint for a biometric
	Given a transaction hash, biometric data
	When I want to validate it, the client call the GRPC endpoint
	Then I get the validation processed
*/
func TestRequestBiometricValidationClient(t *testing.T) {
	txMiners := map[mining.TransactionType]mining.TransactionMiner{
		mining.BiometricTransaction: accountMining.NewBiometricMiner(mockcrypto.NewSigner(), mockcrypto.NewHasher()),
	}

	crypto := Crypto{
		decrypter: mockcrypto.NewDecrypter(),
		signer:    mockcrypto.NewSigner(),
		hasher:    mockcrypto.NewHasher(),
	}

	port := 2005
	conf := system.UnirisConfig{
		Datamining: system.DataMiningConfiguration{
			ExternalPort: port,
		},
	}

	aiClient := mocktransport.NewAIClient()
	miner := mining.NewService(aiClient, nil, nil, nil, mockcrypto.NewSigner(), nil, conf, txMiners)

	grpcServer := grpc.NewServer()
	defer grpcServer.Stop()

	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
		assert.Nil(t, err)

		services := Services{
			mining: miner,
		}

		handler := NewExternalServerHandler(services, crypto, conf)
		api.RegisterExternalServer(grpcServer, handler)
		log.Printf("External grpc Server listening on 127.0.0.1:%d", port)
		err = grpcServer.Serve(lis)
		assert.Nil(t, err)

	}()

	time.Sleep(1 * time.Second)

	bioData := account.NewBiometricData("hash", "enc addr", "enc addr", "enc aes key", "pub", "pub", account.NewSignatures("sig", "sig"))

	cli := NewExternalClient(crypto, conf)
	valid, err := cli.RequestValidation("127.0.0.1", mining.BiometricTransaction, "hash", bioData)
	assert.Nil(t, err)
	assert.NotNil(t, valid)
}

/*
Scenario: Call RequestStorage GRPC endpoint for keychain
	Given a keychain data and its endorsement
	When I want to store it, the client call the GRPC endpoint
	Then the data stored
*/
func TestRequestKeychainStorageClient(t *testing.T) {
	db := mockstorage.NewDatabase()
	accLister := accountListing.NewService(db)
	aiClient := mocktransport.NewAIClient()
	accAdder := accountAdding.NewService(aiClient, db, accLister, mockcrypto.NewSigner(), mockcrypto.NewHasher())

	crypto := Crypto{
		decrypter: mockcrypto.NewDecrypter(),
		signer:    mockcrypto.NewSigner(),
		hasher:    mockcrypto.NewHasher(),
	}

	port := 2006
	conf := system.UnirisConfig{
		Datamining: system.DataMiningConfiguration{
			ExternalPort: port,
		},
	}

	grpcServer := grpc.NewServer()
	defer grpcServer.Stop()

	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
		assert.Nil(t, err)

		services := Services{
			accAdd: accAdder,
		}

		handler := NewExternalServerHandler(services, crypto, conf)
		api.RegisterExternalServer(grpcServer, handler)
		log.Printf("External grpc Server listening on 127.0.0.1:%d", port)
		err = grpcServer.Serve(lis)
		assert.Nil(t, err)

	}()

	time.Sleep(1 * time.Second)

	keychainData := account.NewKeychainData("enc addr", "enc wallet", "pub", "pub", account.NewSignatures("sig", "sig"))
	end := mining.NewEndorsement("", "hash",
		mining.NewMasterValidation([]string{""}, "robotkey", mining.NewValidation(mining.ValidationOK, time.Now(), "robotkey", "sig")),
		[]mining.Validation{mining.NewValidation(mining.ValidationOK, time.Now(), "pub", "sig")},
	)

	cli := NewExternalClient(crypto, conf)
	err := cli.RequestStorage("127.0.0.1", mining.KeychainTransaction, keychainData, end)
	assert.Nil(t, err)

	kc, _ := db.FindLastKeychain("hash")
	assert.NotNil(t, kc)
}

/*
Scenario: Call RequestStorage GRPC endpoint for biometric
	Given a biometric data and its endorsement
	When I want to store it, the client call the GRPC endpoint
	Then the data stored
*/
func TestRequestBiometricStorageClient(t *testing.T) {
	db := mockstorage.NewDatabase()
	accLister := accountListing.NewService(db)
	aiClient := mocktransport.NewAIClient()
	accAdder := accountAdding.NewService(aiClient, db, accLister, mockcrypto.NewSigner(), mockcrypto.NewHasher())

	crypto := Crypto{
		decrypter: mockcrypto.NewDecrypter(),
		signer:    mockcrypto.NewSigner(),
		hasher:    mockcrypto.NewHasher(),
	}

	port := 2007
	conf := system.UnirisConfig{
		Datamining: system.DataMiningConfiguration{
			ExternalPort: port,
		},
	}

	grpcServer := grpc.NewServer()
	defer grpcServer.Stop()

	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
		assert.Nil(t, err)

		services := Services{
			accAdd: accAdder,
		}

		handler := NewExternalServerHandler(services, crypto, conf)
		api.RegisterExternalServer(grpcServer, handler)
		log.Printf("External grpc Server listening on 127.0.0.1:%d", port)
		err = grpcServer.Serve(lis)
		assert.Nil(t, err)

	}()

	time.Sleep(1 * time.Second)

	bioData := account.NewBiometricData("hash", "enc addr", "enc addr", "enc aes key", "pub", "pub", account.NewSignatures("sig", "sig"))
	end := mining.NewEndorsement("", "hash",
		mining.NewMasterValidation([]string{""}, "robotkey", mining.NewValidation(mining.ValidationOK, time.Now(), "robotkey", "sig")),
		[]mining.Validation{mining.NewValidation(mining.ValidationOK, time.Now(), "pub", "sig")},
	)

	cli := NewExternalClient(crypto, conf)
	err := cli.RequestStorage("127.0.0.1", mining.BiometricTransaction, bioData, end)
	assert.Nil(t, err)

	kc, _ := db.FindBiometric("hash")
	assert.NotNil(t, kc)
}

package mining

import (
	"testing"

	"github.com/uniris/uniris-core/datamining/pkg"

	"github.com/stretchr/testify/assert"

	"github.com/uniris/uniris-core/datamining/pkg/account"
	"github.com/uniris/uniris-core/datamining/pkg/account/listing"
	"github.com/uniris/uniris-core/datamining/pkg/mining"
	"github.com/uniris/uniris-core/datamining/pkg/mock"
)

/*
Scenario: Get last transaction hash
	Given an address
	When I want to retrieve the last transaction hash
	Then I get it
*/
func TestKeychainGetLastTransactionHash(t *testing.T) {
	db := mock.NewDatabase()
	db.StoreKeychain(account.NewKeychain(&account.KeyChainData{
		WalletAddr: "address",
	}, datamining.NewEndorsement("", "hash", nil, nil)))
	miner := keychainMiner{accLister: listing.NewService(db)}
	lastHash, err := miner.GetLastTransactionHash("address")
	assert.Nil(t, err)
	assert.Equal(t, "hash", lastHash)
}

/*
Scenario: Checks the keychain data integrity
	Given a transaction hash and keychain data
	When I want to check if the data match the transaction
	Then I get no errors
*/
func TestKeychainIntegrity(t *testing.T) {
	miner := keychainMiner{hasher: mockKeychainHasher{}}
	err := miner.checkDataIntegrity("hash", &account.KeyChainData{})
	assert.Nil(t, err)
}

/*
Scenario: Checks the keychain data integrity
	Given a invalid transaction hash for a keychain data
	When I want to check if the data match the transaction
	Then I get an errors
*/
func TestInvalidKeychainIntegrity(t *testing.T) {
	miner := keychainMiner{hasher: mockBadKeychainHasher{}}
	err := miner.checkDataIntegrity("hash", &account.KeyChainData{})
	assert.Equal(t, mining.ErrInvalidTransaction, err)
}

/*
Scenario: Checks the keychain data signature
	Given keychain data
	When I want to check if the signature match the transaction
	Then I get no errors
*/
func TestKeychainSignature(t *testing.T) {
	miner := keychainMiner{signer: mockKeychainSigner{}}
	err := miner.checkDataSignature(&account.KeyChainData{})
	assert.Nil(t, err)
}

/*
Scenario: Check keychain data as master node
	Given a transaction hash and keychain data
	When I want to check it as master
	Then I get not error
*/
func TestKeychainMasterCheck(t *testing.T) {
	miner := NewKeychainMiner(mockKeychainSigner{}, mockKeychainHasher{}, nil)
	err := miner.CheckAsMaster("hash", &account.KeyChainData{})
	assert.Nil(t, err)
}

/*
Scenario: Check keychain data as slave node
	Given a transaction hash and keychain data
	When I want to check it as slave
	Then I get not error
*/
func TestKeychainSlaveCheck(t *testing.T) {
	miner := NewKeychainMiner(mockKeychainSigner{}, mockKeychainHasher{}, nil)
	err := miner.CheckAsSlave("hash", &account.KeyChainData{})
	assert.Nil(t, err)
}

type mockKeychainHasher struct{}

func (h mockKeychainHasher) HashUnsignedKeychainData(data UnsignedKeychainData) (string, error) {
	return "hash", nil
}

type mockBadKeychainHasher struct{}

func (h mockBadKeychainHasher) HashUnsignedKeychainData(data UnsignedKeychainData) (string, error) {
	return "other hash", nil
}

type mockKeychainSigner struct{}

func (s mockKeychainSigner) CheckKeychainSignature(pubK string, data UnsignedKeychainData, sig string) error {
	return nil
}

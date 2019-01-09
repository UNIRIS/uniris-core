package mining

import (
	"testing"

	"github.com/stretchr/testify/assert"

	datamining "github.com/uniris/uniris-core/datamining/pkg"
	"github.com/uniris/uniris-core/datamining/pkg/account"
	"github.com/uniris/uniris-core/datamining/pkg/account/listing"
	"github.com/uniris/uniris-core/datamining/pkg/mining"
	"github.com/uniris/uniris-core/datamining/pkg/storage/mock"
)

/*
Scenario: Get last transaction hash
	Given an address
	When I want to retrieve the last transaction hash
	Then I get it
*/
func TestKeychainGetLastTransactionHash(t *testing.T) {
	db := mock.NewDatabase()

	prop := datamining.NewProposal(datamining.NewProposedKeyPair("enc pv key", "pub key"))
	kc := account.NewKeychain("enc addr", "enc wallet", "id pub", prop, "id sig", "em sig")
	eKc := account.NewEndorsedKeychain("address", kc, mining.NewEndorsement("", "hash", nil, nil))
	db.StoreKeychain(eKc)
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
	prop := datamining.NewProposal(datamining.NewProposedKeyPair("enc pv key", "pub key"))
	kc := account.NewKeychain("enc addr", "enc wallet", "id pub", prop, "id sig", "em sig")
	err := miner.checkDataIntegrity("hash", kc)
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
	prop := datamining.NewProposal(datamining.NewProposedKeyPair("enc pv key", "pub key"))

	kc := account.NewKeychain("enc addr", "enc wallet", "id pub", prop, "id sig", "em sig")
	err := miner.checkDataIntegrity("hash", kc)
	assert.Equal(t, mining.ErrInvalidTransaction, err)
}

/*
Scenario: Check keychain data as master peer
	Given a transaction hash and keychain data
	When I want to check it as master
	Then I get not error
*/
func TestKeychainMasterCheck(t *testing.T) {
	miner := NewKeychainMiner(mockKeychainSigner{}, mockKeychainHasher{}, nil)

	prop := datamining.NewProposal(datamining.NewProposedKeyPair("enc pv key", "pub key"))

	kc := account.NewKeychain("enc addr", "enc wallet", "id pub", prop, "id sig", "em sig")
	err := miner.CheckAsMaster("hash", kc)
	assert.Nil(t, err)
}

/*
Scenario: Check keychain data as slave peer
	Given a transaction hash and keychain data
	When I want to check it as slave
	Then I get not error
*/
func TestKeychainSlaveCheck(t *testing.T) {
	miner := NewKeychainMiner(mockKeychainSigner{}, mockKeychainHasher{}, nil)
	prop := datamining.NewProposal(datamining.NewProposedKeyPair("enc pv key", "pub key"))

	data := account.NewKeychain("enc addr", "enc wallet", "id pub", prop, "id sig", "em sig")
	err := miner.CheckAsSlave("hash", data)
	assert.Nil(t, err)
}

type mockKeychainHasher struct{}

func (h mockKeychainHasher) HashKeychain(data account.Keychain) (string, error) {
	return "hash", nil
}

func (h mockKeychainHasher) HashEndorsedKeychain(data account.EndorsedKeychain) (string, error) {
	return "hash", nil
}

type mockBadKeychainHasher struct{}

func (h mockBadKeychainHasher) HashKeychain(data account.Keychain) (string, error) {
	return "other hash", nil
}

func (h mockBadKeychainHasher) HashEndorsedKeychain(data account.EndorsedKeychain) (string, error) {
	return "other hash", nil
}

type mockKeychainSigner struct{}

func (s mockKeychainSigner) VerifyKeychainSignatures(account.Keychain) error {
	return nil
}

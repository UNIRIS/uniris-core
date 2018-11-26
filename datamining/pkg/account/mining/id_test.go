package mining

import (
	"testing"

	"github.com/stretchr/testify/assert"

	datamining "github.com/uniris/uniris-core/datamining/pkg"
	"github.com/uniris/uniris-core/datamining/pkg/account"
	"github.com/uniris/uniris-core/datamining/pkg/mining"
)

/*
Scenario: Checks the ID data integrity
	Given a transaction hash and ID data
	When I want to check if the data match the transaction
	Then I get no errors
*/
func TestIDIntegrity(t *testing.T) {
	miner := idMiner{hasher: mockIDHasher{}}
	prop := datamining.NewProposal(datamining.NewProposedKeyPair("enc pv key", "pub key"))

	id := account.NewID("hash", "enc addr", "enc addr", "enc aes key", "id pub", prop, "id sig", "em sig")
	err := miner.checkDataIntegrity("hash", id)
	assert.Nil(t, err)
}

/*
Scenario: Checks the ID data integrity
	Given a invalid transaction hash for a ID data
	When I want to check if the data match the transaction
	Then I get an errors
*/
func TestInvalidIDIntegrity(t *testing.T) {
	miner := idMiner{hasher: mockBadIDHasher{}}
	prop := datamining.NewProposal(datamining.NewProposedKeyPair("enc pv key", "pub key"))

	id := account.NewID("hash", "enc addr", "enc addr", "enc aes key", "id pub", prop, "id sig", "em sig")
	err := miner.checkDataIntegrity("hash", id)
	assert.Equal(t, mining.ErrInvalidTransaction, err)
}

/*
Scenario: Check ID data as master peer
	Given a transaction hash and ID data
	When I want to check it as master
	Then I get not error
*/
func TestIDMasterCheck(t *testing.T) {
	miner := NewIDMiner(mockIDSigner{}, mockIDHasher{})
	prop := datamining.NewProposal(datamining.NewProposedKeyPair("enc pv key", "pub key"))

	id := account.NewID("hash", "enc addr", "enc addr", "enc aes key", "id pub", prop, "id sig", "em sig")
	err := miner.CheckAsMaster("hash", id)
	assert.Nil(t, err)
}

/*
Scenario: Check ID data as slave peer
	Given a transaction hash and ID data
	When I want to check it as slave
	Then I get not error
*/
func TestIDSlaveCheck(t *testing.T) {
	miner := NewIDMiner(mockIDSigner{}, mockIDHasher{})
	prop := datamining.NewProposal(datamining.NewProposedKeyPair("enc pv key", "pub key"))

	id := account.NewID("hash", "enc addr", "enc addr", "enc aes key", "id pub", prop, "id sig", "em sig")
	err := miner.CheckAsSlave("hash", id)
	assert.Nil(t, err)
}

type mockIDHasher struct{}

func (h mockIDHasher) HashID(account.ID) (string, error) {
	return "hash", nil
}

func (h mockIDHasher) HashEndorsedID(account.EndorsedID) (string, error) {
	return "hash", nil
}

type mockBadIDHasher struct{}

func (h mockBadIDHasher) HashID(account.ID) (string, error) {
	return "other hash", nil
}

func (h mockBadIDHasher) HashEndorsedID(account.EndorsedID) (string, error) {
	return "hash", nil
}

type mockIDSigner struct{}

func (s mockIDSigner) VerifyIDSignatures(account.ID) error {
	return nil
}

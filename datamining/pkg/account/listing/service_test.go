package listing

import (
	"sort"
	"testing"
	"time"

	"github.com/uniris/uniris-core/datamining/pkg"

	"github.com/stretchr/testify/assert"

	"github.com/uniris/uniris-core/datamining/pkg/account"
	"github.com/uniris/uniris-core/datamining/pkg/mining"
)

/*
Scenario: TestGetLastKeychain keychain
	Given a empty database
	When I add two keychain
	Then I get the last keychain
*/
func TestGetLastKeychain(t *testing.T) {
	db := new(databasemock)
	s := NewService(db)

	prop := datamining.NewProposal(datamining.NewProposedKeyPair("enc pv key", "pub key"))
	keychain := account.NewKeychain("enc addr", "enc wallet", "id pub", "id sig", "em sig", prop)

	endors1 := mining.NewEndorsement(
		"", "hash1",
		mining.NewMasterValidation([]string{}, "key1", mining.NewValidation(mining.ValidationOK, time.Now(), "robotkey", "sig")),
		[]mining.Validation{
			mining.NewValidation(mining.ValidationOK, time.Now(), "pub", "sig"),
		},
	)

	eKc1 := account.NewEndorsedKeychain("address", keychain, endors1)

	time.Sleep(1 * time.Second)
	endors2 := mining.NewEndorsement(
		"hash1", "hash2",
		mining.NewMasterValidation([]string{}, "key1", mining.NewValidation(mining.ValidationOK, time.Now(), "robotkey", "sig")),
		[]mining.Validation{
			mining.NewValidation(mining.ValidationOK, time.Now(), "pub", "sig"),
		},
	)

	eKc2 := account.NewEndorsedKeychain("address", keychain, endors2)

	db.Keychains = append(db.Keychains, eKc1)
	db.Keychains = append(db.Keychains, eKc2)

	kc, err := s.GetLastKeychain("address")
	assert.Nil(t, err)
	assert.NotNil(t, keychain)
	assert.Equal(t, "hash2", kc.Endorsement().TransactionHash())
	assert.Equal(t, "hash1", kc.Endorsement().LastTransactionHash())
}

/*
Scenario: Get ID
	Given a empty database
	When I add an ID
	Then return values of a GetID are the exepeted ones
*/
func TestGetID(t *testing.T) {

	db := new(databasemock)
	s := NewService(db)

	prop := datamining.NewProposal(datamining.NewProposedKeyPair("enc pv key", "pub key"))
	id := account.NewID("hash1", "enc addr robot", "enc addr person", "enc aes key", "id pub key", "id sig", "em sig", prop)

	endors := mining.NewEndorsement(
		"", "hash1",
		mining.NewMasterValidation([]string{}, "key1", mining.NewValidation(mining.ValidationOK, time.Now(), "robotkey", "sig")),
		[]mining.Validation{
			mining.NewValidation(mining.ValidationOK, time.Now(), "pub", "sig"),
		},
	)

	eID := account.NewEndorsedID(id, endors)

	db.IDs = append(db.IDs, eID)
	wa, err := s.GetID("hash1")
	assert.Nil(t, err)
	assert.NotNil(t, wa)
}

type databasemock struct {
	IDs       []account.EndorsedID
	Keychains []account.EndorsedKeychain
}

func (d *databasemock) FindID(idHash string) (account.EndorsedID, error) {
	for _, id := range d.IDs {
		if id.Hash() == idHash {
			return id, nil
		}
	}
	return nil, nil
}

func (d *databasemock) FindLastKeychain(addr string) (account.EndorsedKeychain, error) {
	sort.Slice(d.Keychains, func(i, j int) bool {
		iTimestamp := d.Keychains[i].Endorsement().MasterValidation().ProofOfWorkValidation().Timestamp().Unix()
		jTimestamp := d.Keychains[j].Endorsement().MasterValidation().ProofOfWorkValidation().Timestamp().Unix()
		return iTimestamp > jTimestamp
	})

	for _, b := range d.Keychains {
		if b.Address() == addr {
			return b, nil
		}
	}
	return nil, nil
}

package listing

import (
	"sort"
	"testing"
	"time"

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

	sigs := account.NewSignatures("sig1", "sig2")
	keychainData := account.NewKeychainData("xxx", "xxx", "xxx", "xxx", sigs)

	masterValid1 := mining.NewMasterValidation([]string{}, "robotKey", mining.NewValidation(mining.ValidationOK, time.Now(), "pubkey", "signature"))
	endors1 := mining.NewEndorsement("", "hash1", masterValid1, []mining.Validation{})

	kc1 := account.NewKeychain("address", keychainData, endors1)

	time.Sleep(1 * time.Second)

	masterValid2 := mining.NewMasterValidation([]string{}, "robotKey", mining.NewValidation(mining.ValidationOK, time.Now(), "pubkey", "signature"))
	endors2 := mining.NewEndorsement("hash1", "hash2", masterValid2, []mining.Validation{})

	kc2 := account.NewKeychain("address", keychainData, endors2)

	db.Keychains = append(db.Keychains, kc1)
	db.Keychains = append(db.Keychains, kc2)

	keychain, err := s.GetLastKeychain("address")
	assert.Nil(t, err)
	assert.NotNil(t, keychain)
	assert.Equal(t, "hash2", keychain.Endorsement().TransactionHash())
	assert.Equal(t, "hash1", keychain.Endorsement().LastTransactionHash())
}

/*
Scenario: Get biometric
	Given a empty database
	When I add a Biometric
	Then return values of a GetBiometric are the exepeted ones
*/
func TestGetBiometric(t *testing.T) {

	db := new(databasemock)
	s := NewService(db)

	sigs := account.NewSignatures("sig1", "sig2")
	data := account.NewBiometricData("hash1", "xxx", "xxx", "xxx", "xxx", "xxx", sigs)

	masterValid := mining.NewMasterValidation([]string{}, "robotKey", mining.NewValidation(mining.ValidationOK, time.Now(), "pubkey", "signature"))
	endors := mining.NewEndorsement("", "xxx", masterValid, []mining.Validation{})

	b := account.NewBiometric(data, endors)

	db.Biometrics = append(db.Biometrics, b)
	wa, err := s.GetBiometric("hash1")
	assert.Nil(t, err)
	assert.NotNil(t, wa)
}

type databasemock struct {
	Biometrics []account.Biometric
	Keychains  []account.Keychain
}

func (d *databasemock) FindBiometric(bh string) (account.Biometric, error) {
	for _, b := range d.Biometrics {
		if b.PersonHash() == bh {
			return b, nil
		}
	}
	return nil, nil
}

func (d *databasemock) FindLastKeychain(addr string) (account.Keychain, error) {
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

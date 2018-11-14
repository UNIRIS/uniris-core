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
	keychainData := account.NewKeychainData("cipher addr", "cipher wallet", "person pub", sigs)

	endors1 := mining.NewEndorsement(
		"", "hash1",
		mining.NewMasterValidation([]string{}, "key1", mining.NewValidation(mining.ValidationOK, time.Now(), "robotkey", "sig")),
		[]mining.Validation{
			mining.NewValidation(mining.ValidationOK, time.Now(), "pub", "sig"),
		},
	)

	kc1 := account.NewKeychain("address", keychainData, endors1)

	time.Sleep(1 * time.Second)
	endors2 := mining.NewEndorsement(
		"hash1", "hash2",
		mining.NewMasterValidation([]string{}, "key1", mining.NewValidation(mining.ValidationOK, time.Now(), "robotkey", "sig")),
		[]mining.Validation{
			mining.NewValidation(mining.ValidationOK, time.Now(), "pub", "sig"),
		},
	)

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
	data := account.NewBiometricData("hash1", "cipher addr robot", "cipher addr person", "cipher aes key", "person key", sigs)

	endors := mining.NewEndorsement(
		"", "hash1",
		mining.NewMasterValidation([]string{}, "key1", mining.NewValidation(mining.ValidationOK, time.Now(), "robotkey", "sig")),
		[]mining.Validation{
			mining.NewValidation(mining.ValidationOK, time.Now(), "pub", "sig"),
		},
	)

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

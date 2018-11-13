package adding

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/uniris/uniris-core/datamining/pkg/account"
	"github.com/uniris/uniris-core/datamining/pkg/mining"
)

/*
Scenario: Store a keychain
	Given a data data
	When I want to store a keychain
	Then the wallet is stored on the database
*/
func TestStoreKeychain(t *testing.T) {
	repo := &databasemock{}
	s := NewService(repo)

	sigs := account.NewSignatures("sig1", "sig2")

	data := account.NewKeychainData("xxx", "xxx", "xxx", "xxx", sigs)
	kc := account.NewKeychain("addr", data, mining.NewEndorsement("", "hash", nil, nil))

	err := s.StoreKeychain(kc)
	assert.Nil(t, err)
	l := len(repo.keychains)
	assert.Equal(t, 1, l)
	assert.Equal(t, "addr", repo.keychains[0].Address())
	assert.Equal(t, "hash", repo.keychains[0].Endorsement().TransactionHash())
}

/*
Scenario: Stroe a biometric
	Given a bio data
	When I want to store a biometric data
	Then the bio data are stored on the database
*/
func TestStoreBiometric(t *testing.T) {
	repo := &databasemock{}
	s := NewService(repo)

	sigs := account.NewSignatures("sig1", "sig2")

	data := account.NewBiometricData("pHash", "xxx", "xxx", "xxx", "xxx", "xxx", sigs)
	bio := account.NewBiometric(data, mining.NewEndorsement("", "hash", nil, nil))
	err := s.StoreBiometric(bio)
	assert.Nil(t, err)
	l := len(repo.biometrics)
	assert.Equal(t, 1, l)
	assert.Equal(t, 1, l)
	assert.Equal(t, "pHash", repo.biometrics[0].PersonHash())
}

type databasemock struct {
	biometrics []account.Biometric
	keychains  []account.Keychain
}

func (d *databasemock) StoreKeychain(kc account.Keychain) error {
	d.keychains = append(d.keychains, kc)
	return nil
}

func (d *databasemock) StoreBiometric(b account.Biometric) error {
	d.biometrics = append(d.biometrics, b)
	return nil
}

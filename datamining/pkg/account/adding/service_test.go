package adding

import (
	"testing"

	"github.com/stretchr/testify/assert"

	datamining "github.com/uniris/uniris-core/datamining/pkg"
	"github.com/uniris/uniris-core/datamining/pkg/account"
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

	sigs := datamining.Signatures{
		BiodSig:   "sig1",
		PersonSig: "sig2",
	}

	data := &account.KeyChainData{
		WalletAddr:      "addr1",
		CipherAddrRobot: "xxxx",
		CipherWallet:    "xxxx",
		PersonPubk:      "xxxx",
		BiodPubk:        "xxxx",
		Sigs:            sigs,
	}

	err := s.StoreKeychain(data, datamining.NewEndorsement(
		"",
		"hash",
		nil,
		nil,
	))
	assert.Nil(t, err)
	l := len(repo.keychains)
	assert.Equal(t, 1, l)
	assert.Equal(t, "addr1", repo.keychains[0].WalletAddr())
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

	sigs := datamining.Signatures{
		BiodSig:   "sig1",
		PersonSig: "sig2",
	}

	bdata := &account.BioData{
		PersonHash:      "hash1",
		BiodPubk:        "xxxx",
		CipherAddrBio:   "xxxx",
		CipherAddrRobot: "xxxx",
		CipherAESKey:    "xxxx",
		PersonPubk:      "xxxx",
		Sigs:            sigs,
	}

	err := s.StoreBiometric(bdata, datamining.NewEndorsement(
		"",
		"hash",
		nil,
		nil,
	))
	assert.Nil(t, err)
	l := len(repo.biometrics)
	assert.Equal(t, 1, l)
	assert.Equal(t, 1, l)
	assert.Equal(t, "hash1", repo.biometrics[0].PersonHash())
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

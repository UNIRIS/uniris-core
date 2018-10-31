package adding

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	datamining "github.com/uniris/uniris-core/datamining/pkg"
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

	wdata := &datamining.KeyChainData{
		WalletAddr:      "addr1",
		CipherAddrRobot: "xxxx",
		CipherWallet:    "xxxx",
		PersonPubk:      "xxxx",
		BiodPubk:        "xxxx",
		Sigs:            sigs,
	}

	w := datamining.NewKeychain(wdata, datamining.NewEndorsement(
		time.Now(),
		"hash",
		nil,
		nil,
	), "old hash")
	err := s.StoreKeychain(w)
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

	bdata := &datamining.BioData{
		PersonHash:      "hash1",
		BiodPubk:        "xxxx",
		CipherAddrBio:   "xxxx",
		CipherAddrRobot: "xxxx",
		CipherAESKey:    "xxxx",
		PersonPubk:      "xxxx",
		Sigs:            sigs,
	}

	b := datamining.NewBiometric(bdata, datamining.NewEndorsement(
		time.Now(),
		"hash",
		nil,
		nil,
	))

	err := s.StoreBiometric(b)
	assert.Nil(t, err)
	l := len(repo.biometrics)
	assert.Equal(t, 1, l)
	assert.Equal(t, 1, l)
	assert.Equal(t, "hash1", repo.biometrics[0].PersonHash())
}

type databasemock struct {
	biometrics []*datamining.Biometric
	keychains  []*datamining.Keychain
}

func (d *databasemock) StoreKeychain(w *datamining.Keychain) error {
	d.keychains = append(d.keychains, w)
	return nil
}

func (d *databasemock) StoreBiometric(bw *datamining.Biometric) error {
	d.biometrics = append(d.biometrics, bw)
	return nil
}

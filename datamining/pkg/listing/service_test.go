package listing

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	datamining "github.com/uniris/uniris-core/datamining/pkg"
)

/*
Scenario: Get keychain
	Given a empty database
	When I add a keychain
	Then return values of a GetKeychain  are the exepeted ones
*/
func TestGetKeychain(t *testing.T) {

	db := new(databasemock)
	s := NewService(db)

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

	oldTxnHash := "xxx"

	endors := datamining.NewEndorsement(time.Now(), "xxx", &datamining.MasterValidation{}, []datamining.Validation{})

	w := datamining.NewKeychain(wdata, endors, oldTxnHash)

	db.StoreKeychain(w)
	wa, err := s.GetKeychain("addr1")
	assert.Nil(t, err)
	assert.NotNil(t, wa)
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

	endors := datamining.NewEndorsement(time.Now(), "xxxx", &datamining.MasterValidation{}, []datamining.Validation{})

	bw := datamining.NewBiometric(bdata, endors)

	db.StoreBiometric(bw)
	wa, err := s.GetBiometric("hash1")
	assert.Nil(t, err)
	assert.NotNil(t, wa)
}

type databasemock struct {
	biometrics []*datamining.Biometric
	keychains  []*datamining.Keychain
}

func (d *databasemock) FindBiometric(bh string) (*datamining.Biometric, error) {
	for _, b := range d.biometrics {
		if b.PersonHash() == bh {
			return b, nil
		}
	}
	return nil, nil
}

func (d *databasemock) FindKeychain(addr string) (*datamining.Keychain, error) {
	for _, b := range d.keychains {
		if b.WalletAddr() == addr {
			return b, nil
		}
	}
	return nil, nil
}

func (d *databasemock) ListBiodPubKeys() ([]string, error) {
	return []string{}, nil
}

func (d *databasemock) StoreKeychain(k *datamining.Keychain) error {
	d.keychains = append(d.keychains, k)
	return nil
}

func (d *databasemock) StoreBiometric(b *datamining.Biometric) error {
	d.biometrics = append(d.biometrics, b)
	return nil
}

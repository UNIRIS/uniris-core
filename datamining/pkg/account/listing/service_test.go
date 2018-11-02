package listing

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	datamining "github.com/uniris/uniris-core/datamining/pkg"
	"github.com/uniris/uniris-core/datamining/pkg/account"
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

	wdata := &account.KeyChainData{
		WalletAddr:      "addr1",
		CipherAddrRobot: "xxxx",
		CipherWallet:    "xxxx",
		PersonPubk:      "xxxx",
		BiodPubk:        "xxxx",
		Sigs:            sigs,
	}

	oldTxnHash := "xxx"

	masterValid := datamining.NewMasterValidation([]string{}, "robotKey", datamining.NewValidation(datamining.ValidationOK, time.Now(), "pubkey", "signature"))
	endors := datamining.NewEndorsement(time.Now(), "xxx", masterValid, []datamining.Validation{})

	kc := account.NewKeychain(wdata, endors, oldTxnHash)

	db.Keychains = append(db.Keychains, kc)
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

	bdata := &account.BioData{
		PersonHash:      "hash1",
		BiodPubk:        "xxxx",
		CipherAddrBio:   "xxxx",
		CipherAddrRobot: "xxxx",
		CipherAESKey:    "xxxx",
		PersonPubk:      "xxxx",
		Sigs:            sigs,
	}

	masterValid := datamining.NewMasterValidation([]string{}, "robotKey", datamining.NewValidation(datamining.ValidationOK, time.Now(), "pubkey", "signature"))
	endors := datamining.NewEndorsement(time.Now(), "xxx", masterValid, []datamining.Validation{})

	b := account.NewBiometric(bdata, endors)

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

func (d *databasemock) FindKeychain(addr string) (account.Keychain, error) {
	for _, b := range d.Keychains {
		if b.WalletAddr() == addr {
			return b, nil
		}
	}
	return nil, nil
}

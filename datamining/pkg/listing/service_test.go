package listing

import (
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	datamining "github.com/uniris/uniris-core/datamining/pkg"
)

/*
Scenario: List Wallet
	Given a empty database
	When I add a Wallet
	Then return values of a GetWallet  are the exepeted ones
*/
func TestGetWallet(t *testing.T) {

	db := new(databasemock)
	s := NewService(db)

	sigs := datamining.Signatures{
		BiodSig: "sig1",
		EmSig:   "sig2",
	}

	wdata := &datamining.WalletData{
		WalletAddr:      "addr1",
		CipherAddrRobot: "xxxx",
		CipherWallet:    "xxxx",
		EmPubk:          "xxxx",
		BiodPubk:        "xxxx",
		Sigs:            sigs,
	}

	oldTxnHash := "xxx"

	endors := datamining.NewEndorsement(time.Now(), "xxx", &datamining.MasterValidation{}, []datamining.Validation{})

	w := datamining.NewWallet(wdata, endors, oldTxnHash)

	db.StoreWallet(w)
	wa, err := s.GetWallet("addr1")
	assert.Nil(t, err)
	assert.NotNil(t, wa)
}

/*
Scenario: List BioWallet
	Given a empty database
	When I add a BioWallet
	Then return values of a GetBioWallet are the exepeted ones
*/
func TestGetBioWallet(t *testing.T) {

	db := new(databasemock)
	s := NewService(db)

	sigs := datamining.Signatures{
		BiodSig: "sig1",
		EmSig:   "sig2",
	}

	bdata := &datamining.BioData{
		BHash:           "hash1",
		BiodPubk:        "xxxx",
		CipherAddrBio:   "xxxx",
		CipherAddrRobot: "xxxx",
		CipherAESKey:    "xxxx",
		EmPubk:          "xxxx",
		Sigs:            sigs,
	}

	endors := datamining.NewEndorsement(time.Now(), "xxxx", &datamining.MasterValidation{}, []datamining.Validation{})

	bw := datamining.NewBioWallet(bdata, endors)

	db.StoreBioWallet(bw)
	wa, err := s.GetBioWallet("hash1")
	assert.Nil(t, err)
	assert.NotNil(t, wa)
}

type databasemock struct {
	bioWallets []*datamining.BioWallet
	wallets    []*datamining.Wallet
}

func (d *databasemock) FindBioWallet(bh string) (*datamining.BioWallet, error) {
	for _, b := range d.bioWallets {
		if b.Bhash() == bh {
			return b, nil
		}
	}
	return nil, nil
}

func (d *databasemock) FindWallet(addr string) (*datamining.Wallet, error) {
	for _, b := range d.wallets {
		log.Print(b.WalletAddr())
		if b.WalletAddr() == addr {
			return b, nil
		}
	}
	return nil, nil
}

func (d *databasemock) StoreWallet(w *datamining.Wallet) error {
	d.wallets = append(d.wallets, w)
	return nil
}

func (d *databasemock) StoreBioWallet(bw *datamining.BioWallet) error {
	d.bioWallets = append(d.bioWallets, bw)
	return nil
}

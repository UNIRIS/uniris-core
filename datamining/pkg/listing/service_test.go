package listing

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

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
		BiodSig: []byte("sig1"),
		EmSig:   []byte("sig2"),
	}

	wdata := datamining.WalletData{
		WalletAddr:      []byte("addr1"),
		CipherAddrRobot: []byte("xxxx"),
		CipherWallet:    []byte("xxxx"),
		EmPubk:          []byte("xxxx"),
		BiodPubk:        []byte("xxxx"),
		Sigs:            sigs,
	}

	oldTxnHash := datamining.Hash([]byte("xxxx"))

	endors := datamining.NewEndorsement(datamining.Timestamp(time.Now()), datamining.Hash([]byte("xxxx")), datamining.MasterValidation{}, []datamining.Validation{})

	w := datamining.NewWallet(wdata, endors, oldTxnHash)

	db.StoreWallet(w)
	wa, err := s.GetWallet([]byte("addr1"))
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
		BiodSig: []byte("sig1"),
		EmSig:   []byte("sig2"),
	}

	bwdata := datamining.BioData{
		BHash:           []byte("hash1"),
		BiodPubk:        []byte("xxxx"),
		CipherAddrBio:   []byte("xxxx"),
		CipherAddrRobot: []byte("xxxx"),
		CipherAESKey:    []byte("xxxx"),
		EmPubk:          []byte("xxxx"),
		Sigs:            sigs,
	}

	endors := datamining.NewEndorsement(datamining.Timestamp(time.Now()), datamining.Hash([]byte("xxxx")), datamining.MasterValidation{}, []datamining.Validation{})

	bw := datamining.NewBioWallet(bwdata, endors)

	db.StoreBioWallet(bw)
	wa, err := s.GetBioWallet([]byte("hash1"))
	assert.Nil(t, err)
	assert.NotNil(t, wa)
}

type databasemock struct {
	bioWallets []datamining.BioWallet
	wallets    []datamining.Wallet
}

func (d *databasemock) FindBioWallet(bh datamining.BioHash) (b datamining.BioWallet, err error) {
	for _, b := range d.bioWallets {
		if string(b.Bhash()) == string(bh) {
			return b, nil
		}
	}
	return
}

func (d *databasemock) FindWallet(addr datamining.WalletAddr) (b datamining.Wallet, err error) {
	for _, b := range d.wallets {
		if string(b.WalletAddr()) == string(addr) {
			return b, nil
		}
	}
	return
}

func (d *databasemock) StoreWallet(w datamining.Wallet) error {
	d.wallets = append(d.wallets, w)
	return nil
}

func (d *databasemock) StoreBioWallet(bw datamining.BioWallet) error {
	d.bioWallets = append(d.bioWallets, bw)
	return nil
}

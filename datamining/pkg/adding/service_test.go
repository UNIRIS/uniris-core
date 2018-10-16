package adding

import (
	"testing"

	"github.com/stretchr/testify/assert"

	datamining "github.com/uniris/uniris-core/datamining/pkg"
	validating "github.com/uniris/uniris-core/datamining/pkg/validating"
)

/*
Scenario: Add Wallet
	Given a empty database
	When I add a Wallet to the database
	Then the length of the database Wallet elements is 1
*/
func TestAddWallet(t *testing.T) {

	db := new(databasemock)
	v := validating.NewService(mockSigner{}, mockValiationRequester{})
	s := NewService(db, v)

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

	err := s.AddWallet(wdata)
	assert.Nil(t, err)
	l := len(db.wallets)
	assert.Equal(t, 1, l)
}

/*
Scenario: Add BioWallet
	Given a empty database
	When I add a BioWallet to the database
	Then the length of the database BioWallet elements is 1
*/
func TestAddBioWallet(t *testing.T) {

	db := new(databasemock)
	v := validating.NewService(mockSigner{}, mockValiationRequester{})
	s := NewService(db, v)

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

	err := s.AddBioWallet(bwdata)
	assert.Nil(t, err)
	l := len(db.bioWallets)
	assert.Equal(t, 1, l)
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

type mockSigner struct{}

func (s mockSigner) CheckSignature(pubk []byte, data interface{}, der []byte) error {
	return nil
}

type mockValiationRequester struct{}

func (v mockValiationRequester) RequestWalletValidation(validating.Peer, datamining.WalletData) (datamining.Validation, error) {
	return datamining.Validation{}, nil
}

func (v mockValiationRequester) RequestBioValidation(validating.Peer, datamining.BioData) (datamining.Validation, error) {
	return datamining.Validation{}, nil
}

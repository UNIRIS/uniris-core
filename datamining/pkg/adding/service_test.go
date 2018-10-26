package adding

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	datamining "github.com/uniris/uniris-core/datamining/pkg"
)

/*
Scenario: Store a data wallet
	Given a data data
	When I want to store a data wallet
	Then the wallet is stored on the database
*/
func TestStoreWallet(t *testing.T) {
	repo := &databasemock{}
	s := NewService(repo)

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

	w := datamining.NewWallet(wdata, datamining.NewEndorsement(
		time.Now(),
		"hash",
		nil,
		nil,
	), "old hash")
	err := s.StoreWallet(w)
	assert.Nil(t, err)
	l := len(repo.wallets)
	assert.Equal(t, 1, l)
	assert.Equal(t, "addr1", repo.wallets[0].WalletAddr())
	assert.Equal(t, "hash", repo.wallets[0].Endorsement().TransactionHash())
}

/*
Scenario: Stroe a bio wallet
	Given a bio data
	When I want to store a bio wallet
	Then the bio data are stored on the database
*/
func TestStoreBioWallet(t *testing.T) {
	repo := &databasemock{}
	s := NewService(repo)

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

	b := datamining.NewBioWallet(bdata, datamining.NewEndorsement(
		time.Now(),
		"hash",
		nil,
		nil,
	))

	err := s.StoreBioWallet(b)
	assert.Nil(t, err)
	l := len(repo.bioWallets)
	assert.Equal(t, 1, l)
	assert.Equal(t, 1, l)
	assert.Equal(t, "hash1", repo.bioWallets[0].Bhash())
}

type databasemock struct {
	bioWallets []*datamining.BioWallet
	wallets    []*datamining.Wallet
}

func (d *databasemock) StoreWallet(w *datamining.Wallet) error {
	d.wallets = append(d.wallets, w)
	return nil
}

func (d *databasemock) StoreBioWallet(bw *datamining.BioWallet) error {
	d.bioWallets = append(d.bioWallets, bw)
	return nil
}

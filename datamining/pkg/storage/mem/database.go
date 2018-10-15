package mem

import (
	datamining "github.com/uniris/uniris-core/datamining/pkg"
	"github.com/uniris/uniris-core/datamining/pkg/adding"
	"github.com/uniris/uniris-core/datamining/pkg/listing"
)

//Database represents a database
type Database interface {
	adding.Repository
	listing.Repository
}

type db struct {
	bioWallets []datamining.BioWallet
	wallets    []datamining.Wallet
}

//NewDatabase implements the database in memory
func NewDatabase() Database {
	return &db{}
}

//FindBioWallet return the bio wallet
func (d *db) FindBioWallet(bh datamining.BioHash) (b datamining.BioWallet, err error) {
	for _, b := range d.bioWallets {
		if string(b.Bhash()) == string(bh) {
			return b, nil
		}
	}
	return
}

//FindWallet return the wallet
func (d *db) FindWallet(addr datamining.WalletAddr) (b datamining.Wallet, err error) {
	for _, b := range d.wallets {
		if string(b.WalletAddr()) == string(addr) {
			return b, nil
		}
	}
	return
}

//AddWallet add a wallet line to the database
func (d *db) AddWallet(w datamining.Wallet) error {
	d.wallets = append(d.wallets, w)
	return nil
}

//AddBioWallet add a biowallet line to the database
func (d *db) AddBioWallet(bw datamining.BioWallet) error {
	d.bioWallets = append(d.bioWallets, bw)
	return nil
}

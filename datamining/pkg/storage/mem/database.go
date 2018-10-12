package mem

import (
	robot "github.com/uniris/uniris-core/datamining/pkg"
	"github.com/uniris/uniris-core/datamining/pkg/wallet"
)

type db struct {
	bioWallets []wallet.BioWallet
	wallets    []wallet.Wallet
}

//NewDatabase implements the database in memory
func NewDatabase() wallet.Database {
	return &db{}
}

//GetEncWalletAddr return the encrypted wallet address
func (d *db) GetBioWallet(bh robot.BioHash) (b wallet.BioWallet, err error) {
	for _, b := range d.bioWallets {
		if string(b.Bhash()) == string(bh) {
			return b, nil
		}
	}
	return
}

//GetEncWallet return the encrypted wallet
func (d *db) GetWallet(addr []byte) (b wallet.Wallet, err error) {
	for _, b := range d.wallets {
		if string(b.WalletAddr()) == string(addr) {
			return b, nil
		}
	}
	return
}

//AddWallet add a wallet line to the database
func (d *db) AddWallet(w wallet.Wallet) error {
	d.wallets = append(d.wallets, w)
	return nil
}

//AddBioWallet add a biowallet line to the database
func (d *db) AddBioWallet(bw wallet.BioWallet) error {
	d.bioWallets = append(d.bioWallets, bw)
	return nil
}

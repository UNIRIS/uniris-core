package mem

import (
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
func (d *db) GetEncWalletAddr(bh wallet.BioHash) ([]byte, error) {
	for _, b := range d.bioWallets {
		if b.Bhash() == bh {
			return b.CipherAddrRobot(), nil
		}
	}
	return nil, nil
}

//GetEncWallet return the encrypted wallet
func (d *db) GetEncWallet(addr []byte) (wallet.CipherWallet, error) {
	for _, b := range d.wallets {
		if string(b.WalletAddr()) == string(addr) {
			return b.WalletAddr(), nil
		}
	}
	return nil, nil
}

package mem

import (
	datamining "github.com/uniris/uniris-core/datamining/pkg"
	"github.com/uniris/uniris-core/datamining/pkg/adding"
	"github.com/uniris/uniris-core/datamining/pkg/leading"
	"github.com/uniris/uniris-core/datamining/pkg/listing"
)

//Database represents a database
type Database interface {
	adding.Repository
	listing.Repository
	leading.TechRepository
}

type db struct {
	bioWallets   []*datamining.BioWallet
	wallets      []*datamining.Wallet
	sharedBioKey string
}

//NewDatabase implements the database in memory
func NewDatabase(sharedBioPubKey string) Database {
	return &db{
		sharedBioKey: sharedBioPubKey,
	}
}

//FindBioWallet return the bio wallet
func (d *db) FindBioWallet(bioHash string) (*datamining.BioWallet, error) {
	for _, b := range d.bioWallets {
		if b.Bhash() == bioHash {
			return b, nil
		}
	}
	return nil, nil
}

//FindWallet return the wallet
func (d *db) FindWallet(addr string) (*datamining.Wallet, error) {
	for _, b := range d.wallets {
		if b.WalletAddr() == addr {
			return b, nil
		}
	}
	return nil, nil
}

//AddWallet add a wallet line to the database
func (d *db) StoreWallet(w *datamining.Wallet) error {
	d.wallets = append(d.wallets, w)
	return nil
}

//AddBioWallet add a biowallet line to the database
func (d *db) StoreBioWallet(bw *datamining.BioWallet) error {
	d.bioWallets = append(d.bioWallets, bw)
	return nil
}

func (d *db) ListBiodPubKeys() ([]string, error) {
	return []string{d.sharedBioKey}, nil
}

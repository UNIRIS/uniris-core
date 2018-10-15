package listing

import (
	"github.com/uniris/uniris-core/datamining/pkg"
)

//Repository defines methods to get data from the database
type Repository interface {
	FindBioWallet(bh datamining.BioHash) (datamining.BioWallet, error)
	FindWallet(addr datamining.WalletAddr) (datamining.Wallet, error)
}

//Service defines methods for the listing service
type Service interface {
	GetBioWallet(bh datamining.BioHash) (datamining.BioWallet, error)
	GetWallet(addr datamining.WalletAddr) (datamining.Wallet, error)
}

type service struct {
	repo Repository
}

//NewService creates a new listing service
func NewService(repo Repository) Service {
	return service{repo}
}

func (s service) GetBioWallet(bh datamining.BioHash) (w datamining.BioWallet, err error) {
	w, err = s.repo.FindBioWallet(bh)
	if err != nil {
		return
	}
	return
}

func (s service) GetWallet(addr datamining.WalletAddr) (w datamining.Wallet, err error) {
	w, err = s.repo.FindWallet(addr)
	if err != nil {
		return
	}
	return
}

package listing

import (
	"github.com/uniris/uniris-core/datamining/pkg"
)

//Repository defines methods to get data from the database
type Repository interface {
	FindBioWallet(bioHash string) (*datamining.BioWallet, error)
	FindWallet(addr string) (*datamining.Wallet, error)
	ListBiodPubKeys() ([]string, error)
}

//Service defines method for the listing service
type Service interface {
	GetBioWallet(bioHash string) (*datamining.BioWallet, error)
	GetWallet(addr string) (*datamining.Wallet, error)
	ListBiodPubKeys() ([]string, error)
}

type service struct {
	repo Repository
}

//NewService creates a new listing service
func NewService(repo Repository) Service {
	return service{repo}
}

func (s service) GetBioWallet(bioHash string) (*datamining.BioWallet, error) {
	w, err := s.repo.FindBioWallet(bioHash)
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (s service) GetWallet(addr string) (*datamining.Wallet, error) {
	w, err := s.repo.FindWallet(addr)
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (s service) ListBiodPubKeys() ([]string, error) {
	return s.repo.ListBiodPubKeys()
}

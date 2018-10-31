package listing

import (
	"github.com/uniris/uniris-core/datamining/pkg"
)

//repository defines wrap repository methods
type repository interface {
	AccountRepository
	TechRepository
}

//AccountRepository defines methods to get data from the account database
type AccountRepository interface {
	FindBiometric(bioHash string) (*datamining.Biometric, error)
	FindKeychain(addr string) (*datamining.Keychain, error)
}

//TechRepository defines mtehods to get data from the tech database
type TechRepository interface {
	ListBiodPubKeys() ([]string, error)
}

//Service defines method for the listing service
type Service interface {
	GetBiometric(bioHash string) (*datamining.Biometric, error)
	GetKeychain(addr string) (*datamining.Keychain, error)
	ListBiodPubKeys() ([]string, error)
}

type service struct {
	repo repository
}

//NewService creates a new listing service
func NewService(repo repository) Service {
	return service{repo}
}

func (s service) GetBiometric(bioHash string) (*datamining.Biometric, error) {
	w, err := s.repo.FindBiometric(bioHash)
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (s service) GetKeychain(addr string) (*datamining.Keychain, error) {
	w, err := s.repo.FindKeychain(addr)
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (s service) ListBiodPubKeys() ([]string, error) {
	return s.repo.ListBiodPubKeys()
}

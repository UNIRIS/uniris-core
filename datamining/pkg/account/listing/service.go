package listing

import (
	"github.com/uniris/uniris-core/datamining/pkg/account"
)

//Repository defines methods to get data from the account database
type Repository interface {
	FindBiometric(bioHash string) (account.Biometric, error)
	FindKeychain(addr string) (account.Keychain, error)
}

//Service defines method for the listing service
type Service interface {
	GetBiometric(bioHash string) (account.Biometric, error)
	GetKeychain(addr string) (account.Keychain, error)
}

type service struct {
	repo Repository
}

//NewService creates a new listing service
func NewService(repo Repository) Service {
	return service{repo}
}

func (s service) GetBiometric(bioHash string) (account.Biometric, error) {
	w, err := s.repo.FindBiometric(bioHash)
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (s service) GetKeychain(addr string) (account.Keychain, error) {
	w, err := s.repo.FindKeychain(addr)
	if err != nil {
		return nil, err
	}
	return w, nil
}

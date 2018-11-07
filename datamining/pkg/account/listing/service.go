package listing

import (
	"github.com/uniris/uniris-core/datamining/pkg/account"
)

//Repository defines methods to get data from the account database
type Repository interface {
	FindBiometric(bioHash string) (account.Biometric, error)
	FindLastKeychain(addr string) (account.Keychain, error)
}

//Service defines method for the listing service
type Service interface {
	GetBiometric(bioHash string) (account.Biometric, error)
	GetLastKeychain(addr string) (account.Keychain, error)
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

func (s service) GetLastKeychain(addr string) (account.Keychain, error) {
	w, err := s.repo.FindLastKeychain(addr)
	if err != nil {
		return nil, err
	}
	return w, nil
}

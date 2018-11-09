package listing

import (
	"github.com/uniris/uniris-core/datamining/pkg/account"
)

//Repository defines methods to get data from the account database
type Repository interface {

	//FindBiometric retrieve a biometric from a given person hash
	FindBiometric(hash string) (account.Biometric, error)

	//FindLastKeychain retrieve the last keychain from a given account's address
	FindLastKeychain(addr string) (account.Keychain, error)
}

//Service defines method for the listing service
type Service interface {

	//GetBiometric retrieve a biometric from a given person hash
	GetBiometric(personHash string) (account.Biometric, error)

	//GetLastKeychain retrieve the last keychain from a given account's address
	GetLastKeychain(addr string) (account.Keychain, error)
}

type service struct {
	repo Repository
}

//NewService creates a new listing service
func NewService(repo Repository) Service {
	return service{repo}
}

func (s service) GetBiometric(personHash string) (account.Biometric, error) {
	w, err := s.repo.FindBiometric(personHash)
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

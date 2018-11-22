package listing

import (
	"github.com/uniris/uniris-core/datamining/pkg/account"
)

//Repository defines methods to get data from the account database
type Repository interface {

	//FindID retrieve a ID from a given hash
	FindID(idHash string) (account.EndorsedID, error)

	//FindLastKeychain retrieve the last keychain from a given account's address
	FindLastKeychain(addr string) (account.EndorsedKeychain, error)
}

//Service defines method for the listing service
type Service interface {

	//GetID retrieve an ID from a given hash
	GetID(idHash string) (account.EndorsedID, error)

	//GetLastKeychain retrieve the last keychain from a given account's address
	GetLastKeychain(addr string) (account.EndorsedKeychain, error)
}

type service struct {
	repo Repository
}

//NewService creates a new listing service
func NewService(repo Repository) Service {
	return service{repo}
}

func (s service) GetID(idHash string) (account.EndorsedID, error) {
	id, err := s.repo.FindID(idHash)
	if err != nil {
		return nil, err
	}
	return id, nil
}

func (s service) GetLastKeychain(addr string) (account.EndorsedKeychain, error) {
	kc, err := s.repo.FindLastKeychain(addr)
	if err != nil {
		return nil, err
	}
	return kc, nil
}

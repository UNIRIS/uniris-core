package adding

import (
	"github.com/uniris/uniris-core/datamining/pkg/account"
)

//Repository handles account adding storage
type Repository interface {
	StoreKeychain(account.Keychain) error
	StoreBiometric(account.Biometric) error
}

//Service is the interface that provide methods for wallets transactions on robot side
type Service interface {
	StoreKeychain(account.Keychain) error
	StoreBiometric(account.Biometric) error
}

type service struct {
	repo Repository
}

//NewService creates a new adding service
func NewService(repo Repository) Service {
	return service{repo}
}

func (s service) StoreKeychain(kc account.Keychain) error {

	//TODO: check integrity of keychain

	//TODO: handle store pending/ko
	return s.repo.StoreKeychain(kc)
}

func (s service) StoreBiometric(b account.Biometric) error {

	//TODO: check integrity of biometric

	//TODO: handle store pending/ko
	return s.repo.StoreBiometric(b)
}

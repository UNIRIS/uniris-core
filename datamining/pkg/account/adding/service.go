package adding

import (
	"github.com/uniris/uniris-core/datamining/pkg"
	"github.com/uniris/uniris-core/datamining/pkg/account"
)

//Repository handles account adding storage
type Repository interface {
	StoreKeychain(account.Keychain) error
	StoreBiometric(account.Biometric) error
}

//Service is the interface that provide methods for wallets transactions on robot side
type Service interface {
	StoreKeychain(data *account.KeyChainData, endorsement datamining.Endorsement) error
	StoreBiometric(data *account.BioData, endorsement datamining.Endorsement) error
}

type service struct {
	repo Repository
}

//NewService creates a new adding service
func NewService(repo Repository) Service {
	return service{repo}
}

func (s service) StoreKeychain(data *account.KeyChainData, end datamining.Endorsement) error {

	//TODO: check integrity of keychain

	kc := account.NewKeychain(data, end)

	//TODO: handle store pending/ko
	return s.repo.StoreKeychain(kc)
}

func (s service) StoreBiometric(data *account.BioData, end datamining.Endorsement) error {

	//TODO: check integrity of biometric

	b := account.NewBiometric(data, end)

	//TODO: handle store pending/ko
	return s.repo.StoreBiometric(b)
}

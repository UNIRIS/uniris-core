package adding

import (
	"github.com/uniris/uniris-core/datamining/pkg"
)

//AccountRepository defines methods to add account data into the database
type AccountRepository interface {
	StoreKeychain(*datamining.Keychain) error
	StoreBiometric(*datamining.Biometric) error
}

//Service is the interface that provide methods for wallets transactions on robot side
type Service interface {
	StoreKeychain(w *datamining.Keychain) error
	StoreBiometric(bw *datamining.Biometric) error
}

type service struct {
	accRepo AccountRepository
}

//NewService creates a new adding service
func NewService(accRepo AccountRepository) Service {
	return service{accRepo}
}

func (s service) StoreKeychain(w *datamining.Keychain) error {

	//TODO: check integrity of keychain

	//TODO: handle store pending/ko
	return s.accRepo.StoreKeychain(w)
}

func (s service) StoreBiometric(b *datamining.Biometric) error {

	//TODO: check integrity of biometric

	//TODO: handle store pending/ko
	return s.accRepo.StoreBiometric(b)
}

package adding

import (
	"github.com/uniris/uniris-core/datamining/pkg"
)

//Repository defines methods to add data into the database
type Repository interface {
	StoreWallet(*datamining.Wallet) error
	StoreBioWallet(*datamining.BioWallet) error
}

//Service is the interface that provide methods for wallets transactions on robot side
type Service interface {
	StoreWallet(w *datamining.Wallet) error
	StoreBioWallet(bw *datamining.BioWallet) error
}

type service struct {
	repo Repository
}

//NewService creates a new adding service
func NewService(repo Repository) Service {
	return service{repo}
}

func (s service) StoreWallet(w *datamining.Wallet) error {
	//TODO: handle store pending/ko
	return s.repo.StoreWallet(w)
}

func (s service) StoreBioWallet(bw *datamining.BioWallet) error {
	//TODO: handle store pending/ko
	return s.repo.StoreBioWallet(bw)
}

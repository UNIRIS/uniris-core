package walletadding

import (
	"github.com/uniris/uniris-core/datamining/pkg"
	"github.com/uniris/uniris-core/datamining/pkg/validating"
)

//Repository defines methods to add data into the database
type Repository interface {
	AddWallet(w datamining.Wallet) error
	AddBioWallet(bw datamining.BioWallet) error
}

//Service is the interface that provide methods for wallets transactions on robot side
type Service interface {
	AddWallet(w datamining.WalletData) error
	AddBioWallet(bw datamining.BioData) error
}

type service struct {
	repo  Repository
	valid validating.Service
}

//NewService creates a new adding service
func NewService(repo Repository, valid validating.Service) Service {
	return service{repo, valid}
}

func (s service) AddWallet(data datamining.WalletData) error {
	t, oth, th, mv, err := s.valid.EndorseWalletAsMaster(data)
	if err != nil {
		return err
	}

	rv, err := s.valid.EndorseWallet(data)
	if err != nil {
		return err
	}

	e := datamining.NewEndorsement(t, th, mv, rv)
	w := datamining.NewWallet(data, e, oth)

	if err := s.repo.AddWallet(w); err != nil {
		return err
	}

	return nil
}

func (s service) AddBioWallet(data datamining.BioData) error {
	t, th, mv, err := s.valid.EndorseBioWalletAsMaster(data)
	if err != nil {
		return err
	}
	rv, err := s.valid.EndorseBioWallet(data)
	if err != nil {
		return err
	}

	e := datamining.NewEndorsement(t, th, mv, rv)
	w := datamining.NewBioWallet(data, e)

	if err := s.repo.AddBioWallet(w); err != nil {
		return err
	}

	return nil
}

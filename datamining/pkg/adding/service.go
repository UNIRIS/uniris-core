package adding

import (
	"github.com/uniris/uniris-core/datamining/pkg"
	"github.com/uniris/uniris-core/datamining/pkg/validating"
)

//Repository defines methods to add data into the database
type Repository interface {
	StoreWallet(w datamining.Wallet) error
	StoreBioWallet(bw datamining.BioWallet) error
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
	walEndor, oldTxHash, err := s.valid.EndorseWalletAsMaster(data)
	if err != nil {
		return err
	}

	peerValidators := make([]validating.Peer, 0)
	//TODO: retrieve the peers validators

	validations, err := s.valid.AskWalletValidations(peerValidators, data)
	if err != nil {
		return err
	}

	e := datamining.NewEndorsement(walEndor.Timestamp(),
		walEndor.TransactionHash(),
		walEndor.MasterValidation(),
		validations)

	w := datamining.NewWallet(data, e, oldTxHash)

	if err := s.repo.StoreWallet(w); err != nil {
		return err
	}

	return nil
}

func (s service) AddBioWallet(data datamining.BioData) error {
	masterEndors, err := s.valid.EndorseBioWalletAsMaster(data)
	if err != nil {
		return err
	}

	peerValidators := make([]validating.Peer, 0)
	//TODO: retrieve the peers validators

	validations, err := s.valid.AskBioWalletValidations(peerValidators, data)
	if err != nil {
		return err
	}

	e := datamining.NewEndorsement(masterEndors.Timestamp(),
		masterEndors.TransactionHash(),
		masterEndors.MasterValidation(),
		validations)

	w := datamining.NewBioWallet(data, e)
	if err := s.repo.StoreBioWallet(w); err != nil {
		return err
	}

	return nil
}

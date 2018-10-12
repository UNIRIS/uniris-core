package validating

import (
	"time"

	"github.com/uniris/uniris-core/datamining/pkg"
)

//Service is the interface that provide methods for wallets validation
type Service interface {
	EndorseWalletAsMaster(datamining.WalletData) (datamining.Timestamp, datamining.Hash, datamining.Hash, datamining.MasterValidation, error)
	EndorseWallet(datamining.WalletData) ([]datamining.Validation, error)
	EndorseBioWalletAsMaster(datamining.BioData) (datamining.Timestamp, datamining.Hash, datamining.MasterValidation, error)
	EndorseBioWallet(datamining.BioData) ([]datamining.Validation, error)
}

type service struct {
}

//NewService creates a approving service
func NewService() Service {
	return &service{}
}

func (s service) EndorseWalletAsMaster(w datamining.WalletData) (datamining.Timestamp, datamining.Hash, datamining.Hash, datamining.MasterValidation, error) {
	t := datamining.Timestamp(time.Now())
	mrv := datamining.MasterValidation{}
	return t, nil, nil, mrv, nil
}

func (s service) EndorseWallet(w datamining.WalletData) ([]datamining.Validation, error) {
	rv := []datamining.Validation{}
	return rv, nil
}

func (s service) EndorseBioWalletAsMaster(bw datamining.BioData) (datamining.Timestamp, datamining.Hash, datamining.MasterValidation, error) {
	t := datamining.Timestamp(time.Now())
	mrv := datamining.MasterValidation{}
	return t, nil, mrv, nil
}

func (s service) EndorseBioWallet(bw datamining.BioData) ([]datamining.Validation, error) {
	rv := []datamining.Validation{}
	return rv, nil
}

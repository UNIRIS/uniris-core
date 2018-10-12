package validation

import (
	formater "github.com/uniris/uniris-core/datamining/pkg/walletformating"
	"hash"
	"time"
)

//Service is the interface that provide methods for wallets validation
type Service interface {
	WalletValidateAsMaster(formater.FormatedWallet) (time.Time, hash.Hash64, hash.Hash64, MasterRobotValidation, error)
	WalletValidate(formater.FormatedWallet) ([]RobotValidation, error)
	BioWalletValidateAsMaster(fbw formater.FormatedBioWallet) (time.Time, hash.Hash64, MasterRobotValidation, error)
	BioWalletValidate(fbw formater.FormatedBioWallet) ([]RobotValidation, error)
}

type service struct {
}

func (s service) WalletValidateAsMaster(fw formater.FormatedWallet) (time.Time, hash.Hash64, hash.Hash64, MasterRobotValidation, error) {
	t := time.Now()
	mrv := MasterRobotValidation{}
	return t, nil, nil, mrv, nil
}

func (s service) WalletValidate(fw formater.FormatedWallet) ([]RobotValidation, error) {
	rv := []RobotValidation{}
	return rv, nil
}

func (s service) BioWalletValidateAsMaster(fbw formater.FormatedBioWallet) (time.Time, hash.Hash64, MasterRobotValidation, error) {
	t := time.Now()
	mrv := MasterRobotValidation{}
	return t, nil, mrv, nil
}

func (s service) BioWalletValidate(fbw formater.FormatedBioWallet) ([]RobotValidation, error) {
	rv := []RobotValidation{}
	return rv, nil
}

//NewService creates a validation service
func NewService() Service {
	return &service{}
}

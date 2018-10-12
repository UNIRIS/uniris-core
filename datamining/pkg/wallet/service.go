package wallet

import (
	validation "github.com/uniris/uniris-core/datamining/pkg/validation"
	formater "github.com/uniris/uniris-core/datamining/pkg/walletformating"
)

//Service is the interface that provide methods for wallets transactions on robot side
type Service interface {
	NewWallet() (Wallet, error)
	NewBioWallet() (BioWallet, error)
	WriteWallet(w Wallet) error
	WriteBioWallet(bw BioWallet) error
}

type service struct {
	fw  formater.FormatedWallet
	fbw formater.FormatedBioWallet
	db  Database
}

func (s service) NewWallet() (w Wallet, err error) {
	v := validation.NewService()
	t, oth, th, mv, err := v.WalletValidateAsMaster(s.fw)
	if err != nil {
		return
	}
	rv, err := v.WalletValidate(s.fw)
	if err != nil {
		return
	}
	w = Wallet{
		fw:           s.fw,
		timeStamp:    t,
		oldTxnHash:   oth,
		txnHash:      th,
		masterRobotv: mv,
		robotsv:      rv,
	}
	return
}

func (s service) NewBioWallet() (bw BioWallet, err error) {
	v := validation.NewService()
	t, th, mv, err := v.BioWalletValidateAsMaster(s.fbw)
	if err != nil {
		return
	}
	rv, err := v.BioWalletValidate(s.fbw)
	if err != nil {
		return
	}
	bw = BioWallet{
		fbw:          s.fbw,
		timeStamp:    t,
		txnHash:      th,
		masterRobotv: mv,
		robotsv:      rv,
	}
	return bw, nil
}

func (s service) WriteWallet(w Wallet) error {
	err := s.db.AddWallet(w)
	return err
}

func (s service) WriteBioWallet(bw BioWallet) error {
	err := s.db.AddBioWallet(bw)
	return err
}

//NewService creates a Wallet service
func NewService(d Database, fw formater.FormatedWallet, fbw formater.FormatedBioWallet) Service {
	return &service{
		fw:  fw,
		fbw: fbw,
		db:  d,
	}
}

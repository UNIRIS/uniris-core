package formater

import ()

//Service is the interface that provide methods for wallets json formating
type Service interface {
	FormatWallet([]byte) (FormatedWallet, error)
	FormatBioWallet([]byte) (FormatedBioWallet, error)
}

type service struct {
}

func (s service) FormatWallet([]byte) (FormatedWallet, error) {
	//read from grpc / verify sign / decrypt data
	fw := FormatedWallet{}
	return fw, nil
}

func (s service) FormatBioWallet([]byte) (FormatedBioWallet, error) {
	//read from grpc / verify sign / decrypt data
	fw := FormatedBioWallet{}
	return fw, nil
}

//NewService creates a formater service
func NewService() Service {
	return &service{}
}

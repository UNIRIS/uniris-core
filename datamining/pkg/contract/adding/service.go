package contract

import "github.com/uniris/uniris-core/datamining/pkg/contract"

type Repository interface {
	StoreEndorsedContract(contract.EndorsedContract) error
	StoreEndorsedMessage(contract.EndorsedMessage) error
}

type Service interface {
	StoreEndorsedContract(contract.EndorsedContract) error
	StoreEndorsedMessage(contract.EndorsedMessage) error
}

type addService struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return addService{repo}
}

func (s addService) StoreEndorsedContract(c contract.EndorsedContract) error {
	//REMOVE STORAGE CHECKS FOR DEMO PURPOSE

	return s.repo.StoreEndorsedContract(c)
}

func (s addService) StoreEndorsedMessage(m contract.EndorsedMessage) error {
	//REMOVE STORAGE CHECKS FOR DEMO PURPOSE

	return s.repo.StoreEndorsedMessage(m)
}

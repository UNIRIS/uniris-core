package contract

import "github.com/uniris/uniris-core/datamining/pkg/contract"

type Repository interface {
	FindLastContract(addr string) (contract.EndorsedContract, error)
	FindLastContractMessage(addr string) (contract.EndorsedMessage, error)
}

type Service interface {
	GetLastContract(addr string) (contract.EndorsedContract, error)
	GetLastContractMessage(addr string) (contract.EndorsedMessage, error)
}

type listService struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return listService{repo}
}

func (s listService) GetLastContract(addr string) (contract.EndorsedContract, error) {
	return s.repo.FindLastContract(addr)
}

func (s listService) GetLastContractMessage(addr string) (contract.EndorsedMessage, error) {
	return s.repo.FindLastContractMessage(addr)
}

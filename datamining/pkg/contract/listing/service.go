package contract

import (
	"errors"
	"fmt"
	"strings"

	"github.com/uniris/uniris-core/datamining/pkg/contract"
	uniris "github.com/uniris/uniris-interpreter/pkg"
)

type Repository interface {
	FindLastContract(addr string) (contract.EndorsedContract, error)
	FindLastContractMessage(addr string) (contract.EndorsedMessage, error)
	FindContractByAddressAndTransactionHash(addr string, hash string) (contract.EndorsedContract, error)
	FindMessagesByContract(addr string) ([]contract.EndorsedMessage, error)
	FindMessageByAddressAndTransactionHash(addr string, hash string) (contract.EndorsedMessage, error)
}

type Service interface {
	GetLastContract(addr string) (contract.EndorsedContract, error)
	GetLastContractMessage(addr string) (contract.EndorsedMessage, error)
	GetContractState(addr string) (string, error)
	GetContractByAddressAndTransaction(addr string, hash string) (contract.EndorsedContract, error)
	GetContractMessageByContractAndTransaction(addr string, hash string) (contract.EndorsedMessage, error)
	GetContractMessages(addr string) ([]contract.EndorsedMessage, error)
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

func (s listService) GetContractMessages(addr string) ([]contract.EndorsedMessage, error) {
	return s.repo.FindMessagesByContract(addr)
}

func (s listService) GetContractState(addr string) (string, error) {
	contract, err := s.repo.FindLastContract(addr)
	if err != nil {
		return "", err
	}
	if contract == nil {
		return "", errors.New("Unknown contract")
	}

	messages, err := s.repo.FindMessagesByContract(addr)
	if err != nil {
		return "", err
	}

	code := contract.Code()
	env := uniris.NewEnvironment(nil)
	_, err = uniris.Interpret(code, env)
	if err != nil {
		return "", nil
	}

	for _, m := range messages {
		params := strings.Join(m.Parameters(), ",")
		_, err := uniris.Interpret(fmt.Sprintf("%v(%v)", m.Method(), params), env)
		if err != nil {
			return "", nil
		}
	}

	res, err := uniris.Interpret("getState()", env)
	if err != nil {
		return "", err
	}

	return res, nil
}

func (s listService) GetContractByAddressAndTransaction(addr string, hash string) (contract.EndorsedContract, error) {
	return s.repo.FindContractByAddressAndTransactionHash(addr, hash)
}

func (s listService) GetContractMessageByContractAndTransaction(addr string, hash string) (contract.EndorsedMessage, error) {
	return s.repo.FindMessageByAddressAndTransactionHash(addr, hash)
}

package contract

import (
	"fmt"
	"strings"

	"github.com/uniris/uniris-core/datamining/pkg/contract"
	uniris "github.com/uniris/uniris-interpreter/pkg"
)

type Repository interface {
	FindLastContract(addr string) (contract.EndorsedContract, error)
	FindLastContractMessage(addr string) (contract.EndorsedMessage, error)
	FindMessagesByContract(addr string) ([]contract.EndorsedMessage, error)
}

type Service interface {
	GetLastContract(addr string) (contract.EndorsedContract, error)
	GetLastContractMessage(addr string) (contract.EndorsedMessage, error)
	GetContractState(addr string) (string, error)
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

func (s listService) GetContractState(addr string) (string, error) {
	contract, err := s.repo.FindLastContract(addr)
	if err != nil {
		return "", err
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
		env.Set("messagePublicKey", m.PublicKey())

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

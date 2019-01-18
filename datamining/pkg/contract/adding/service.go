package contract

import (
	"errors"
	"fmt"
	"strings"

	"github.com/uniris/uniris-core/datamining/pkg/contract"
	contractListing "github.com/uniris/uniris-core/datamining/pkg/contract/listing"
	uniris "github.com/uniris/uniris-interpreter/pkg"
)

type Repository interface {
	StoreEndorsedContract(contract.EndorsedContract) error
	StoreEndorsedMessage(contract.EndorsedMessage) error
}

type Service interface {
	StoreEndorsedContract(contract.EndorsedContract) error
	StoreEndorsedMessage(contract.EndorsedMessage) error
}

type addService struct {
	repo   Repository
	lister contractListing.Service
}

func NewService(repo Repository, lister contractListing.Service) Service {
	return addService{repo, lister}
}

func (s addService) StoreEndorsedContract(c contract.EndorsedContract) error {
	//REMOVE STORAGE CHECKS FOR DEMO PURPOSE

	return s.repo.StoreEndorsedContract(c)
}

func (s addService) StoreEndorsedMessage(m contract.EndorsedMessage) error {
	//REMOVE STORAGE CHECKS FOR DEMO PURPOSE

	if err := s.executeContract(m); err != nil {
		return err
	}
	return s.repo.StoreEndorsedMessage(m)
}

func (s addService) executeContract(m contract.EndorsedMessage) error {
	contract, err := s.lister.GetLastContract(m.ContractAddress())
	if err != nil {
		return err
	}
	if contract == nil {
		return errors.New("Unknown smart contract")
	}

	code := contract.Code()
	env := uniris.NewEnvironment(nil)
	_, err = uniris.Interpret(code, env)
	if err != nil {
		return err
	}

	prevMsgs, err := s.lister.GetContractMessages(m.ContractAddress())
	if err != nil {
		return err
	}
	for _, m := range prevMsgs {
		params := strings.Join(m.Parameters(), ",")
		_, err := uniris.Interpret(fmt.Sprintf("%v(%v)", m.Method(), params), env)
		if err != nil {
			return err
		}
	}

	params := strings.Join(m.Parameters(), ",")
	_, err = uniris.Interpret(fmt.Sprintf("%v(%v)", m.Method(), params), env)
	if err != nil {
		return err
	}

	return nil
}

package contract

import (
	"github.com/uniris/uniris-core/datamining/pkg/contract"
	contractListing "github.com/uniris/uniris-core/datamining/pkg/contract/listing"
	"github.com/uniris/uniris-core/datamining/pkg/mining"
)

type contractMessageMiner struct {
	sigVerifier contract.SignatureVerifier
	hasher      contract.Hasher
	listService contractListing.Service
}

func NewContractMessageMiner(sigVerif contract.SignatureVerifier, hash contract.Hasher, lister contractListing.Service) mining.TransactionMiner {
	return contractMiner{
		sigVerifier: sigVerif,
		hasher:      hash,
		listService: lister,
	}
}

func (m contractMessageMiner) GetLastTransactionHash(addr string) (string, error) {
	contract, err := m.listService.GetLastContractMessage(addr)
	if err != nil {
		return "", err
	}
	if contract == nil {
		return "", nil
	}
	return contract.Endorsement().TransactionHash(), nil
}

func (m contractMessageMiner) CheckAsMaster(txHash string, data interface{}) error {
	msg := data.(contract.Message)
	if err := m.checkDataIntegrity(txHash, msg); err != nil {
		return err
	}
	if err := m.sigVerifier.VerifyContractMessageSignature(msg); err != nil {
		return err
	}

	return nil
}
func (m contractMessageMiner) CheckAsSlave(txHash string, data interface{}) error {
	msg := data.(contract.Message)
	if err := m.checkDataIntegrity(txHash, msg); err != nil {
		return err
	}
	if err := m.sigVerifier.VerifyContractMessageSignature(msg); err != nil {
		return err
	}

	return nil
}

func (m contractMessageMiner) checkDataIntegrity(txHash string, msg contract.Message) error {
	hash, err := m.hasher.HashContractMessage(msg)

	if err != nil {
		return err
	}
	if hash != txHash {
		return mining.ErrInvalidTransaction
	}
	return nil
}

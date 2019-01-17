package contract

import (
	"github.com/uniris/uniris-core/datamining/pkg/contract"
	contractListing "github.com/uniris/uniris-core/datamining/pkg/contract/listing"
	"github.com/uniris/uniris-core/datamining/pkg/mining"
)

type contractMiner struct {
	sigVerifier contract.SignatureVerifier
	hasher      contract.Hasher
	listService contractListing.Service
}

func NewContractMiner(sigVerif contract.SignatureVerifier, hash contract.Hasher, lister contractListing.Service) mining.TransactionMiner {
	return contractMiner{
		sigVerifier: sigVerif,
		hasher:      hash,
		listService: lister,
	}
}

func (m contractMiner) GetLastTransactionHash(addr string) (string, error) {
	contract, err := m.listService.GetLastContract(addr)
	if err != nil {
		return "", err
	}
	if contract == nil {
		return "", nil
	}
	return contract.Endorsement().TransactionHash(), nil
}

func (m contractMiner) CheckAsMaster(txHash string, data interface{}) error {
	contract := data.(contract.Contract)
	if err := m.checkDataIntegrity(txHash, contract); err != nil {
		return err
	}
	if err := m.sigVerifier.VerifyContractSignature(contract); err != nil {
		return err
	}

	return nil
}
func (m contractMiner) CheckAsSlave(txHash string, data interface{}) error {
	contract := data.(contract.Contract)
	if err := m.checkDataIntegrity(txHash, contract); err != nil {
		return err
	}
	if err := m.sigVerifier.VerifyContractSignature(contract); err != nil {
		return err
	}

	return nil
}

func (m contractMiner) checkDataIntegrity(txHash string, c contract.Contract) error {
	hash, err := m.hasher.HashContract(c)

	if err != nil {
		return err
	}
	if hash != txHash {
		return mining.ErrInvalidTransaction
	}
	return nil
}

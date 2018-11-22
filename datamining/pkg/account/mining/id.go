package mining

import (
	"github.com/uniris/uniris-core/datamining/pkg/account"
	"github.com/uniris/uniris-core/datamining/pkg/mining"
)

type idMiner struct {
	sigVerifier account.IDSignatureVerifier
	hasher      account.IDHasher
}

//NewIDMiner creates a miner for the biometric transaction
func NewIDMiner(sigVerifier account.IDSignatureVerifier, hasher account.IDHasher) mining.TransactionMiner {
	return idMiner{sigVerifier, hasher}
}

func (m idMiner) GetLastTransactionHash(addr string) (string, error) {
	return "", nil
}

func (m idMiner) CheckAsMaster(txHash string, data interface{}) error {
	id := data.(account.ID)
	if err := m.checkDataIntegrity(txHash, id); err != nil {
		return err
	}
	if err := m.sigVerifier.VerifyIDSignatures(id); err != nil {
		return err
	}

	return nil
}

func (m idMiner) CheckAsSlave(txHash string, data interface{}) error {
	id := data.(account.ID)
	if err := m.checkDataIntegrity(txHash, id); err != nil {
		return err
	}
	if err := m.sigVerifier.VerifyIDSignatures(id); err != nil {
		return err
	}

	return nil
}

func (m idMiner) checkDataIntegrity(txHash string, id account.ID) error {
	hash, err := m.hasher.HashID(id)
	if err != nil {
		return err
	}
	if hash != txHash {
		return mining.ErrInvalidTransaction
	}
	return nil
}

package mining

import (
	"github.com/uniris/uniris-core/datamining/pkg/account"
	"github.com/uniris/uniris-core/datamining/pkg/mining"
)

type biometricMiner struct {
	signer account.BiometricSigner
	hasher account.BiometricHasher
}

//NewBiometricMiner creates a miner for the biometric transaction
func NewBiometricMiner(signer account.BiometricSigner, hasher account.BiometricHasher) mining.TransactionMiner {
	return biometricMiner{signer, hasher}
}

func (m biometricMiner) GetLastTransactionHash(addr string) (string, error) {
	return "", nil
}

func (m biometricMiner) CheckAsMaster(txHash string, data interface{}) error {
	biometric := data.(account.BiometricData)
	if err := m.checkDataIntegrity(txHash, biometric); err != nil {
		return err
	}
	if err := m.signer.VerifyBiometricDataSignatures(biometric); err != nil {
		return err
	}

	return nil
}

func (m biometricMiner) CheckAsSlave(txHash string, data interface{}) error {
	biometric := data.(account.BiometricData)
	if err := m.checkDataIntegrity(txHash, biometric); err != nil {
		return err
	}
	if err := m.signer.VerifyBiometricDataSignatures(biometric); err != nil {
		return err
	}

	return nil
}

func (m biometricMiner) checkDataIntegrity(txHash string, data account.BiometricData) error {
	hash, err := m.hasher.HashBiometricData(data)
	if err != nil {
		return err
	}
	if hash != txHash {
		return mining.ErrInvalidTransaction
	}
	return nil
}

package mining

import (
	"github.com/uniris/uniris-core/datamining/pkg/account"
	"github.com/uniris/uniris-core/datamining/pkg/mining"
)

type biometricMiner struct {
	sigVerifier account.BiometricSignatureVerifier
	hasher      account.BiometricHasher
}

//NewBiometricMiner creates a miner for the biometric transaction
func NewBiometricMiner(sigVerifier account.BiometricSignatureVerifier, hasher account.BiometricHasher) mining.TransactionMiner {
	return biometricMiner{sigVerifier, hasher}
}

func (m biometricMiner) GetLastTransactionHash(addr string) (string, error) {
	return "", nil
}

func (m biometricMiner) CheckAsMaster(txHash string, data interface{}) error {
	biometric := data.(account.BiometricData)
	if err := m.checkDataIntegrity(txHash, biometric); err != nil {
		return err
	}
	if err := m.sigVerifier.VerifyBiometricDataSignatures(biometric); err != nil {
		return err
	}

	return nil
}

func (m biometricMiner) CheckAsSlave(txHash string, data interface{}) error {
	biometric := data.(account.BiometricData)
	if err := m.checkDataIntegrity(txHash, biometric); err != nil {
		return err
	}
	if err := m.sigVerifier.VerifyBiometricDataSignatures(biometric); err != nil {
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

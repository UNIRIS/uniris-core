package mining

import (
	"github.com/uniris/uniris-core/datamining/pkg/account"
	"github.com/uniris/uniris-core/datamining/pkg/mining"
)

//BiometricHasher define methods to hash biometric data
type BiometricHasher interface {

	//NewBiometricDataHash creates a hash of the biometric data
	NewBiometricDataHash(account.BiometricData) (string, error)
}

//BiometricSigner define methods to handle signatures
type BiometricSigner interface {

	//CheckBiometricDataSignature checks the signature of the biometric data using the shared robot public key
	CheckBiometricDataSignature(pubKey string, data account.BiometricData, sig string) error
}

type biometricMiner struct {
	signer BiometricSigner
	hasher BiometricHasher
}

//NewBiometricMiner creates a miner for the biometric transaction
func NewBiometricMiner(signer BiometricSigner, hasher BiometricHasher) mining.TransactionMiner {
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
	if err := m.checkDataSignature(biometric); err != nil {
		return err
	}

	return nil
}

func (m biometricMiner) CheckAsSlave(txHash string, data interface{}) error {
	biometric := data.(account.BiometricData)
	if err := m.checkDataIntegrity(txHash, biometric); err != nil {
		return err
	}
	if err := m.checkDataSignature(biometric); err != nil {
		return err
	}

	return nil
}

func (m biometricMiner) checkDataIntegrity(txHash string, data account.BiometricData) error {
	hash, err := m.hasher.NewBiometricDataHash(data)
	if err != nil {
		return err
	}
	if hash != txHash {
		return mining.ErrInvalidTransaction
	}
	return nil
}

func (m biometricMiner) checkDataSignature(data account.BiometricData) error {
	if err := m.signer.CheckBiometricDataSignature(data.BiodPublicKey(), data, data.Signatures().Biod()); err != nil {
		return err
	}

	if err := m.signer.CheckBiometricDataSignature(data.PersonPublicKey(), data, data.Signatures().Person()); err != nil {
		return err
	}

	return nil
}

package checks

import (
	"errors"

	"github.com/uniris/uniris-core/datamining/pkg/account"
	"github.com/uniris/uniris-core/datamining/pkg/mining"
)

type integrityChecker struct {
	h TransactionDataHasher
}

//NewIntegrityChecker creates an intergrity checker
func NewIntegrityChecker(h TransactionDataHasher) Handler {
	return integrityChecker{h}
}

func (c integrityChecker) CheckData(data interface{}, txHash string) error {

	var hash string
	var err error

	switch data.(type) {
	case *account.KeyChainData:
		rawData := c.getKeychainData(data.(*account.KeyChainData))
		hash, err = c.h.HashTransactionData(rawData)
		if err != nil {
			return err
		}
		return c.checkHash(hash, txHash)
	case *account.BioData:
		rawData := c.getBiometricData(data.(*account.BioData))
		hash, err = c.h.HashTransactionData(rawData)
		if err != nil {
			return err
		}
		return c.checkHash(hash, txHash)
	}

	return errors.New("Unsupported data")
}

func (c integrityChecker) checkHash(hash, txHash string) error {
	if hash != txHash {
		return mining.ErrInvalidTransaction
	}
	return nil
}

func (c integrityChecker) getKeychainData(keychainData *account.KeyChainData) rawKeychainData {
	return rawKeychainData{
		BiodPublicKey:      keychainData.BiodPubk,
		EncryptedAddrRobot: keychainData.CipherAddrRobot,
		EncryptedWallet:    keychainData.CipherWallet,
		PersonPublicKey:    keychainData.PersonPubk,
	}
}

func (c integrityChecker) getBiometricData(biometricData *account.BioData) rawBiometricData {
	return rawBiometricData{
		BiodPublicKey:       biometricData.BiodPubk,
		EncryptedAddrRobot:  biometricData.CipherAddrRobot,
		EncryptedAddrPerson: biometricData.CipherAddrBio,
		PersonHash:          biometricData.PersonHash,
		EncryptedAESKey:     biometricData.CipherAESKey,
		PersonPublicKey:     biometricData.PersonPubk,
	}
}

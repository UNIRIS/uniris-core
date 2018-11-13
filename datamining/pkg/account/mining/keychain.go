package mining

import (
	"github.com/uniris/uniris-core/datamining/pkg/account"
	listing "github.com/uniris/uniris-core/datamining/pkg/account/listing"
	"github.com/uniris/uniris-core/datamining/pkg/mining"
)

type keychainMiner struct {
	signer    account.KeychainSigner
	hasher    account.KeychainHasher
	accLister listing.Service
}

//NewKeychainMiner creates a miner for the keychain transaction
func NewKeychainMiner(signer account.KeychainSigner, hasher account.KeychainHasher, accLister listing.Service) mining.TransactionMiner {
	return keychainMiner{signer, hasher, accLister}
}

func (m keychainMiner) GetLastTransactionHash(addr string) (string, error) {
	keychain, err := m.accLister.GetLastKeychain(addr)
	if err != nil {
		return "", err
	}
	if keychain == nil {
		return "", nil
	}
	return keychain.Endorsement().TransactionHash(), nil
}

func (m keychainMiner) CheckAsMaster(txHash string, data interface{}) error {
	keychain := data.(account.KeychainData)
	if err := m.checkDataIntegrity(txHash, keychain); err != nil {
		return err
	}
	if err := m.signer.VerifyKeychainDataSignatures(keychain); err != nil {
		return err
	}

	return nil
}

func (m keychainMiner) CheckAsSlave(txHash string, data interface{}) error {
	keychain := data.(account.KeychainData)
	if err := m.checkDataIntegrity(txHash, keychain); err != nil {
		return err
	}
	if err := m.signer.VerifyKeychainDataSignatures(keychain); err != nil {
		return err
	}

	return nil
}

func (m keychainMiner) checkDataIntegrity(txHash string, data account.KeychainData) error {
	hash, err := m.hasher.HashKeychainData(data)
	if err != nil {
		return err
	}
	if hash != txHash {
		return mining.ErrInvalidTransaction
	}
	return nil
}

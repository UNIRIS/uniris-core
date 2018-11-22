package mining

import (
	"github.com/uniris/uniris-core/datamining/pkg/account"
	listing "github.com/uniris/uniris-core/datamining/pkg/account/listing"
	"github.com/uniris/uniris-core/datamining/pkg/mining"
)

type keychainMiner struct {
	sigVerifier account.KeychainSignatureVerifier
	hasher      account.KeychainHasher
	accLister   listing.Service
}

//NewKeychainMiner creates a miner for the keychain transaction
func NewKeychainMiner(sigVerifier account.KeychainSignatureVerifier, hasher account.KeychainHasher, accLister listing.Service) mining.TransactionMiner {
	return keychainMiner{sigVerifier, hasher, accLister}
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
	keychain := data.(account.Keychain)
	if err := m.checkDataIntegrity(txHash, keychain); err != nil {
		return err
	}
	if err := m.sigVerifier.VerifyKeychainSignatures(keychain); err != nil {
		return err
	}

	return nil
}

func (m keychainMiner) CheckAsSlave(txHash string, data interface{}) error {
	keychain := data.(account.Keychain)
	if err := m.checkDataIntegrity(txHash, keychain); err != nil {
		return err
	}
	if err := m.sigVerifier.VerifyKeychainSignatures(keychain); err != nil {
		return err
	}

	return nil
}

func (m keychainMiner) checkDataIntegrity(txHash string, kc account.Keychain) error {
	hash, err := m.hasher.HashKeychain(kc)

	if err != nil {
		return err
	}
	if hash != txHash {
		return mining.ErrInvalidTransaction
	}
	return nil
}

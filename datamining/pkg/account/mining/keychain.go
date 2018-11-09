package mining

import (
	"github.com/uniris/uniris-core/datamining/pkg/account"
	listing "github.com/uniris/uniris-core/datamining/pkg/account/listing"
	"github.com/uniris/uniris-core/datamining/pkg/mining"
)

//KeychainSigner define methods to handle keychain signature
type KeychainSigner interface {

	//CheckKeychainDataSignature checks the signature of the keychain data using the shared robot public key
	CheckKeychainDataSignature(pubKey string, data account.KeychainData, sig string) error
}

//KeychainHasher define methods to hash keychain data
type KeychainHasher interface {

	//NewKeychainDataHash creates a hash of the keychain data
	NewKeychainDataHash(data account.KeychainData) (string, error)
}

type keychainMiner struct {
	signer    KeychainSigner
	hasher    KeychainHasher
	accLister listing.Service
}

//NewKeychainMiner creates a miner for the keychain transaction
func NewKeychainMiner(signer KeychainSigner, hasher KeychainHasher, accLister listing.Service) mining.TransactionMiner {
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
	if err := m.checkDataSignature(keychain); err != nil {
		return err
	}

	return nil
}

func (m keychainMiner) CheckAsSlave(txHash string, data interface{}) error {
	keychain := data.(account.KeychainData)
	if err := m.checkDataIntegrity(txHash, keychain); err != nil {
		return err
	}
	if err := m.checkDataSignature(keychain); err != nil {
		return err
	}

	return nil
}

func (m keychainMiner) checkDataIntegrity(txHash string, data account.KeychainData) error {
	hash, err := m.hasher.NewKeychainDataHash(data)
	if err != nil {
		return err
	}
	if hash != txHash {
		return mining.ErrInvalidTransaction
	}
	return nil
}

func (m keychainMiner) checkDataSignature(data account.KeychainData) error {
	if err := m.signer.CheckKeychainDataSignature(data.BiodPublicKey(), data, data.Signatures().Biod()); err != nil {
		return err
	}

	if err := m.signer.CheckKeychainDataSignature(data.PersonPublicKey(), data, data.Signatures().Person()); err != nil {
		return err
	}
	return nil
}

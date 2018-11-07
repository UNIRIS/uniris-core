package mining

import (
	"github.com/uniris/uniris-core/datamining/pkg/account"
	listing "github.com/uniris/uniris-core/datamining/pkg/account/listing"
	"github.com/uniris/uniris-core/datamining/pkg/mining"
)

//UnsignedKeychainData defines keychain data before signature
type UnsignedKeychainData struct {
	PersonPublicKey    string `json:"person_pubk"`
	BiodPublicKey      string `json:"biod_pubk"`
	EncryptedWallet    string `json:"encrypted_wal"`
	EncryptedAddrRobot string `json:"encrypted_addr_robot"`
}

//KeychainSigner define methods to handle keychain signature
type KeychainSigner interface {
	CheckKeychainSignature(pubKey string, data UnsignedKeychainData, sig string) error
}

//KeychainHasher define methods to hash keychain data
type KeychainHasher interface {
	HashUnsignedKeychainData(data UnsignedKeychainData) (string, error)
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
	keychain := data.(*account.KeyChainData)
	if err := m.checkDataIntegrity(txHash, keychain); err != nil {
		return err
	}
	if err := m.checkDataSignature(keychain); err != nil {
		return err
	}

	return nil
}

func (m keychainMiner) CheckAsSlave(txHash string, data interface{}) error {
	keychain := data.(*account.KeyChainData)
	if err := m.checkDataIntegrity(txHash, keychain); err != nil {
		return err
	}
	if err := m.checkDataSignature(keychain); err != nil {
		return err
	}

	return nil
}

func (m keychainMiner) checkDataIntegrity(txHash string, data *account.KeyChainData) error {
	unsignedData := m.buildUnsignedData(data)
	hash, err := m.hasher.HashUnsignedKeychainData(unsignedData)
	if err != nil {
		return err
	}
	if hash != txHash {
		return mining.ErrInvalidTransaction
	}
	return nil
}

func (m keychainMiner) checkDataSignature(data *account.KeyChainData) error {
	unsignedData := m.buildUnsignedData(data)

	if err := m.signer.CheckKeychainSignature(data.BiodPubk, unsignedData, data.Sigs.BiodSig); err != nil {
		return err
	}

	if err := m.signer.CheckKeychainSignature(data.PersonPubk, unsignedData, data.Sigs.PersonSig); err != nil {
		return err
	}
	return nil
}

func (m keychainMiner) buildUnsignedData(data *account.KeyChainData) UnsignedKeychainData {
	return UnsignedKeychainData{
		BiodPublicKey:      data.BiodPubk,
		EncryptedAddrRobot: data.CipherAddrRobot,
		EncryptedWallet:    data.CipherWallet,
		PersonPublicKey:    data.PersonPubk,
	}
}

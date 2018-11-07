package mining

import (
	"github.com/uniris/uniris-core/datamining/pkg/account"
	"github.com/uniris/uniris-core/datamining/pkg/mining"
)

//UnsignedBiometricData defines biometric data before signature
type UnsignedBiometricData struct {
	PersonPublicKey     string `json:"person_pubk"`
	BiodPublicKey       string `json:"biod_pubk"`
	PersonHash          string `json:"person_hash"`
	EncryptedAESKey     string `json:"encrypted_aes_key"`
	EncryptedAddrPerson string `json:"encrypted_addr_person"`
	EncryptedAddrRobot  string `json:"encrypted_addr_robot"`
}

//BiometricHasher define methods to hash biometric data
type BiometricHasher interface {
	HashUnsignedBiometricData(data UnsignedBiometricData) (string, error)
}

//BiometricSigner define methods to handle signatures
type BiometricSigner interface {
	CheckBiometricSignature(pubKey string, data UnsignedBiometricData, sig string) error
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
	biometric := data.(*account.BioData)
	if err := m.checkDataIntegrity(txHash, biometric); err != nil {
		return err
	}
	if err := m.checkDataSignature(biometric); err != nil {
		return err
	}

	return nil
}

func (m biometricMiner) CheckAsSlave(txHash string, data interface{}) error {
	biometric := data.(*account.BioData)
	if err := m.checkDataIntegrity(txHash, biometric); err != nil {
		return err
	}
	if err := m.checkDataSignature(biometric); err != nil {
		return err
	}

	return nil
}

func (m biometricMiner) checkDataIntegrity(txHash string, data *account.BioData) error {
	unsignedData := m.buildUnsignedData(data)

	hash, err := m.hasher.HashUnsignedBiometricData(unsignedData)
	if err != nil {
		return err
	}
	if hash != txHash {
		return mining.ErrInvalidTransaction
	}
	return nil
}

func (m biometricMiner) checkDataSignature(data *account.BioData) error {
	unsignedData := m.buildUnsignedData(data)

	if err := m.signer.CheckBiometricSignature(data.BiodPubk, unsignedData, data.Sigs.BiodSig); err != nil {
		return err
	}

	if err := m.signer.CheckBiometricSignature(data.PersonPubk, unsignedData, data.Sigs.PersonSig); err != nil {
		return err
	}

	return nil
}

func (m biometricMiner) buildUnsignedData(data *account.BioData) UnsignedBiometricData {
	return UnsignedBiometricData{
		BiodPublicKey:       data.BiodPubk,
		EncryptedAddrPerson: data.CipherAddrBio,
		EncryptedAddrRobot:  data.CipherAddrRobot,
		EncryptedAESKey:     data.CipherAESKey,
		PersonHash:          data.PersonHash,
		PersonPublicKey:     data.PersonPubk,
	}
}

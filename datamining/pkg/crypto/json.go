package crypto

import (
	"github.com/uniris/uniris-core/datamining/pkg/mining"
)

type endorsement struct {
	LastTxHash       string           `json:"last_transaction_hash"`
	TxHash           string           `json:"transaction_hash"`
	MasterValidation masterValidation `json:"master_validation"`
	Validations      []validation     `json:"validations"`
}

type masterValidation struct {
	LastTxRvk   []string   `json:"last_transaction_miners"`
	PowRobotKey string     `json:"pow_robot_key"`
	PowValid    validation `json:"pow_validation"`
}

type validation struct {
	Status    mining.ValidationStatus `json:"status"`
	Timestamp int64                   `json:"timestamp"`
	Pubk      string                  `json:"public_key"`
	Sig       string                  `json:"signature"`
}

type biometric struct {
	PersonPublicKey     string      `json:"person_pubk"`
	BIODPublicKey       string      `json:"biod_pubk"`
	PersonHash          string      `json:"person_hash"`
	EncryptedAESKey     string      `json:"encrypted_aes_key"`
	EncryptedAddrPerson string      `json:"encrypted_addr_person"`
	EncryptedAddrRobot  string      `json:"encrypted_addr_robot"`
	BIODSignature       string      `json:"biod_sig"`
	PersonSignature     string      `json:"person_sig"`
	Endorsement         endorsement `json:"endorsement"`
}

type biometricRaw struct {
	PersonPublicKey     string `json:"person_pubk"`
	BIODPublicKey       string `json:"biod_pubk"`
	PersonHash          string `json:"person_hash"`
	EncryptedAESKey     string `json:"encrypted_aes_key"`
	EncryptedAddrPerson string `json:"encrypted_addr_person"`
	EncryptedAddrRobot  string `json:"encrypted_addr_robot"`
}

type biometricData struct {
	PersonPublicKey     string `json:"person_pubk"`
	BIODPublicKey       string `json:"biod_pubk"`
	PersonHash          string `json:"person_hash"`
	EncryptedAESKey     string `json:"encrypted_aes_key"`
	EncryptedAddrPerson string `json:"encrypted_addr_person"`
	EncryptedAddrRobot  string `json:"encrypted_addr_robot"`
	BIODSignature       string `json:"biod_sig"`
	PersonSignature     string `json:"person_sig"`
}

type keychain struct {
	Address            string      `json:"address"`
	PersonPublicKey    string      `json:"person_pubk"`
	BIODPublicKey      string      `json:"biod_pubk"`
	EncryptedWallet    string      `json:"encrypted_wal"`
	EncryptedAddrRobot string      `json:"encrypted_addr_robot"`
	BIODSignature      string      `json:"biod_sig"`
	PersonSignature    string      `json:"person_sig"`
	Endorsement        endorsement `json:"endorsement"`
}

type keychainData struct {
	PersonPublicKey    string `json:"person_pubk"`
	BIODPublicKey      string `json:"biod_pubk"`
	EncryptedWallet    string `json:"encrypted_wal"`
	EncryptedAddrRobot string `json:"encrypted_addr_robot"`
	BIODSignature      string `json:"biod_sig"`
	PersonSignature    string `json:"person_sig"`
}

type keychainRaw struct {
	PersonPublicKey    string `json:"person_pubk"`
	BIODPublicKey      string `json:"biod_pubk"`
	EncryptedWallet    string `json:"encrypted_wal"`
	EncryptedAddrRobot string `json:"encrypted_addr_robot"`
}

type lockRaw struct {
	TxHash         string `json:"transaction_hash"`
	MasterRobotKey string `json:"master_robot_key"`
	Address        string `json:"address"`
}

type validationRaw struct {
	Status    mining.ValidationStatus `json:"status"`
	Timestamp int64                   `json:"timestamp"`
	PublicKey string                  `json:"public_key"`
}

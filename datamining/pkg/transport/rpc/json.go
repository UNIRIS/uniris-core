package rpc

import (
	api "github.com/uniris/uniris-core/datamining/api/protobuf-spec"
	datamining "github.com/uniris/uniris-core/datamining/pkg"
	"github.com/uniris/uniris-core/datamining/pkg/account"
)

//KeychainJSON represents a keychain as JSON
type KeychainJSON struct {
	Data        KeychainDataJSON `json:"data"`
	Endorsement EndorsementJSON  `json:"endorsement"`
}

//BiometricJSON represents a biometric as JSON
type BiometricJSON struct {
	Data        BioDataJSON     `json:"data"`
	Endorsement EndorsementJSON `json:"endorsement"`
}

//EndorsementJSON represents a endorsement as JSON
type EndorsementJSON struct {
	LastTransactionHash string               `json:"last_transaction_hash"`
	TransactionHash     string               `json:"transaction_hash"`
	MasterValidation    MasterValidationJSON `json:"master_validation"`
	Validations         []ValidationJSON     `json:"validations"`
}

//MasterValidationJSON represents a master validation as JSON
type MasterValidationJSON struct {
	LastTransactionMiners []string       `json:"last_transaction_miners"`
	ProofOfWorkRobotKey   string         `json:"proof_of_work_robot_key"`
	ProofOfWorkValidation ValidationJSON `json:"proof_of_work_validation"`
}

//ValidationJSON represents a validation as JSON
type ValidationJSON struct {
	PublicKey string `json:"public_key"`
	Signature string `json:"signature"`
	Status    string `json:"status"`
	Timestamp int64  `json:"timestamp"`
}

//KeychainDataJSON represents keychain data JSON
type KeychainDataJSON struct {
	PersonPublicKey    string `json:"person_pubk"`
	BiodPublicKey      string `json:"biod_pubk"`
	EncryptedWallet    string `json:"encrypted_wal"`
	EncryptedAddrRobot string `json:"encrypted_addr_robot"`
}

//BioDataJSON represents bio data JSON
type BioDataJSON struct {
	PersonPublicKey     string `json:"person_pubk"`
	BiodPublicKey       string `json:"biod_pubk"`
	PersonHash          string `json:"person_hash"`
	EncryptedAESKey     string `json:"encrypted_aes_key"`
	EncryptedAddrPerson string `json:"encrypted_addr_person"`
	EncryptedAddrRobot  string `json:"encrypted_addr_robot"`
}

//Signatures represents signatures JSON
type Signatures struct {
	Person string `json:"person_sig"`
	Biod   string `json:"biod_sig"`
}

//BuildBioDataFromJSON convert JSON to biometric data
func BuildBioDataFromJSON(bioData *BioDataJSON, sig *api.Signature) *account.BioData {
	return &account.BioData{
		PersonHash:      bioData.PersonHash,
		BiodPubk:        bioData.BiodPublicKey,
		CipherAddrBio:   bioData.EncryptedAddrPerson,
		CipherAddrRobot: bioData.EncryptedAddrRobot,
		CipherAESKey:    bioData.EncryptedAESKey,
		PersonPubk:      bioData.PersonPublicKey,
		Sigs: account.Signatures{
			BiodSig:   sig.Biod,
			PersonSig: sig.Person,
		},
	}

}

//BuildKeychainDataFromJSON convert json to keychain data
func BuildKeychainDataFromJSON(keychainData *KeychainDataJSON, sig *api.Signature, clearAddr string) *account.KeyChainData {
	return &account.KeyChainData{
		WalletAddr:      clearAddr,
		BiodPubk:        keychainData.BiodPublicKey,
		CipherAddrRobot: keychainData.EncryptedAddrRobot,
		CipherWallet:    keychainData.EncryptedWallet,
		PersonPubk:      keychainData.PersonPublicKey,
		Sigs: account.Signatures{
			BiodSig:   sig.Biod,
			PersonSig: sig.Person,
		},
	}
}

//BuildKeychainJSON convert keychain to JSON
func BuildKeychainJSON(keychain account.Keychain) KeychainJSON {
	return KeychainJSON{
		Data: KeychainDataJSON{
			BiodPublicKey:      keychain.BiodPublicKey(),
			EncryptedAddrRobot: keychain.CipherAddrRobot(),
			EncryptedWallet:    keychain.CipherWallet(),
			PersonPublicKey:    keychain.PersonPublicKey(),
		},
		Endorsement: BuildEndorsementJSON(keychain.Endorsement()),
	}
}

//BuildBiometricJSON convert biometric to JSON
func BuildBiometricJSON(biometric account.Biometric) BiometricJSON {
	return BiometricJSON{
		Data: BioDataJSON{
			BiodPublicKey:       biometric.BiodPublicKey(),
			EncryptedAddrRobot:  biometric.CipherAddrRobot(),
			EncryptedAddrPerson: biometric.CipherAddrBio(),
			PersonHash:          biometric.PersonHash(),
			EncryptedAESKey:     biometric.CipherAESKey(),
			PersonPublicKey:     biometric.PersonPublicKey(),
		},
		Endorsement: BuildEndorsementJSON(biometric.Endorsement()),
	}
}

//BuildEndorsementJSON convert endorsement to JSON
func BuildEndorsementJSON(end datamining.Endorsement) EndorsementJSON {
	valids := make([]ValidationJSON, 0)
	for _, v := range end.Validations() {
		valids = append(valids, BuildValidationJSON(v))
	}
	return EndorsementJSON{
		LastTransactionHash: end.LastTransactionHash(),
		TransactionHash:     end.TransactionHash(),
		MasterValidation:    BuildMasterValidationJSON(end.MasterValidation()),
		Validations:         valids,
	}
}

//BuildMasterValidationJSON convert master validation to JSON
func BuildMasterValidationJSON(mv datamining.MasterValidation) MasterValidationJSON {
	return MasterValidationJSON{
		LastTransactionMiners: mv.LastTransactionMiners(),
		ProofOfWorkRobotKey:   mv.ProofOfWorkRobotKey(),
		ProofOfWorkValidation: BuildValidationJSON(mv.ProofOfWorkValidation()),
	}
}

//BuildValidationJSON convert validation to JSON
func BuildValidationJSON(v datamining.Validation) ValidationJSON {
	return ValidationJSON{
		PublicKey: v.PublicKey(),
		Signature: v.Signature(),
		Status:    v.Status().String(),
		Timestamp: v.Timestamp().Unix(),
	}
}

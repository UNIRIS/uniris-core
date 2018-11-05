package externalrpc

import (
	api "github.com/uniris/uniris-core/datamining/api/protobuf-spec"
	datamining "github.com/uniris/uniris-core/datamining/pkg"
	"github.com/uniris/uniris-core/datamining/pkg/account"
)

//BuildBioData transform json bio data into bio data
func BuildBioData(bioData *BioDataFromJSON, sig *api.Signature) *account.BioData {
	return &account.BioData{
		PersonHash:      bioData.PersonHash,
		BiodPubk:        bioData.BiodPublicKey,
		CipherAddrBio:   bioData.EncryptedAddrPerson,
		CipherAddrRobot: bioData.EncryptedAddrRobot,
		CipherAESKey:    bioData.EncryptedAESKey,
		PersonPubk:      bioData.PersonPublicKey,
		Sigs: datamining.Signatures{
			BiodSig:   sig.Biod,
			PersonSig: sig.Person,
		},
	}

}

//BuildKeychainData transform json keychain data into keychain data
func BuildKeychainData(keychainData *KeychainDataFromJSON, sig *api.Signature, clearAddr string) *account.KeyChainData {
	return &account.KeyChainData{
		WalletAddr:      clearAddr,
		BiodPubk:        keychainData.BiodPublicKey,
		CipherAddrRobot: keychainData.EncryptedAddrRobot,
		CipherWallet:    keychainData.EncryptedWallet,
		PersonPubk:      keychainData.PersonPublicKey,
		Sigs: datamining.Signatures{
			BiodSig:   sig.Biod,
			PersonSig: sig.Person,
		},
	}
}

//KeychainDataFromJSON represents keychain data JSON
type KeychainDataFromJSON struct {
	PersonPublicKey    string `json:"person_pubk"`
	BiodPublicKey      string `json:"biod_pubk"`
	EncryptedWallet    string `json:"encrypted_wal"`
	EncryptedAddrRobot string `json:"encrypted_addr_robot"`
}

//BioDataFromJSON represents bio data JSON
type BioDataFromJSON struct {
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

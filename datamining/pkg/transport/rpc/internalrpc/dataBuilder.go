package internalrpc

import (
	api "github.com/uniris/uniris-core/datamining/api/protobuf-spec"
	"github.com/uniris/uniris-core/datamining/pkg"
)

//BuildWalletResult constitue details for the WalletResult rpc command
func BuildWalletResult(w *datamining.Wallet, bioW *datamining.BioWallet) *api.WalletResult {
	return &api.WalletResult{
		EncryptedAESkey:        bioW.CipherAESKey(),
		EncryptedWallet:        w.CipherWallet(),
		EncryptedWalletAddress: bioW.CipherAddrBio(),
	}
}

//BuildBioData transform json bio data into bio data
func BuildBioData(bioData BioDataFromJSON, sig *api.Signature) *datamining.BioData {
	return &datamining.BioData{
		BHash:           bioData.PersonHash,
		BiodPubk:        bioData.BiodPublicKey,
		CipherAddrBio:   bioData.EncryptedAddrPerson,
		CipherAddrRobot: bioData.EncryptedAddrRobot,
		CipherAESKey:    bioData.EncryptedAESKey,
		EmPubk:          bioData.PersonPublicKey,
		Sigs: datamining.Signatures{
			BiodSig: sig.Biod,
			EmSig:   sig.Person,
		},
	}

}

//BuildWalletData transform json wallet data into wallet data
func BuildWalletData(walletData *WalletDataFromJSON, sig *api.Signature, clearAddr string) *datamining.WalletData {
	return &datamining.WalletData{
		WalletAddr:      clearAddr,
		BiodPubk:        walletData.BiodPublicKey,
		CipherAddrRobot: walletData.EncryptedAddrRobot,
		CipherWallet:    walletData.EncryptedWallet,
		EmPubk:          walletData.PersonPublicKey,
		Sigs: datamining.Signatures{
			BiodSig: sig.Biod,
			EmSig:   sig.Person,
		},
	}
}

//WalletDataFromJSON represents wallet data JSON
type WalletDataFromJSON struct {
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

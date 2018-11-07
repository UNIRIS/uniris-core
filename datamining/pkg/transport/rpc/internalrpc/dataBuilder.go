package internalrpc

import (
	api "github.com/uniris/uniris-core/datamining/api/protobuf-spec"
	"github.com/uniris/uniris-core/datamining/pkg/account"
)

func buildAccountSearchResult(kc account.Keychain, biometric account.Biometric) *api.AccountSearchResult {
	return &api.AccountSearchResult{
		EncryptedAESkey:  biometric.CipherAESKey(),
		EncryptedWallet:  kc.CipherWallet(),
		EncryptedAddress: biometric.CipherAddrBio(),
	}
}

func buildBioData(bioData BioDataJSON, sig *api.Signature) *account.BioData {
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

func buildKeychainData(keychainData *KeychainDataJSON, sig *api.Signature, clearAddr string) *account.KeyChainData {
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

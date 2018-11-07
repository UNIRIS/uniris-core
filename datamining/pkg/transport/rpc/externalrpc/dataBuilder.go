package externalrpc

import (
	"time"

	api "github.com/uniris/uniris-core/datamining/api/protobuf-spec"
	datamining "github.com/uniris/uniris-core/datamining/pkg"
	"github.com/uniris/uniris-core/datamining/pkg/account"
)

func buildBioDataFromJSON(bioData *BioDataJSON, sig *api.Signature) *account.BioData {
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

func buildKeychainDataFromJSON(keychainData *KeychainDataJSON, sig *api.Signature, clearAddr string) *account.KeyChainData {
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

func buildKeychainAPIResponse(keychain account.Keychain, sig string) *api.KeychainResponse {
	return &api.KeychainResponse{
		Data: &api.KeychainData{
			BiodPubk:        keychain.BiodPublicKey(),
			CipherAddrRobot: keychain.CipherAddrRobot(),
			CipherWallet:    keychain.CipherWallet(),
			PersonPubk:      keychain.PersonPublicKey(),
			Signature: &api.Signature{
				Biod:   keychain.Signatures().BiodSig,
				Person: keychain.Signatures().PersonSig,
			},
		},
		Endorsement: buildAPIEndorsement(keychain.Endorsement()),
		Signature:   sig,
	}
}

func buildBiometricAPIResponse(biometric account.Biometric, sig string) *api.BiometricResponse {
	return &api.BiometricResponse{
		Data: &api.BiometricData{
			BiodPubk:        biometric.BiodPublicKey(),
			CipherAddrBio:   biometric.CipherAddrBio(),
			CipherAddrRobot: biometric.CipherAddrRobot(),
			CipherAESKey:    biometric.CipherAESKey(),
			PersonHash:      biometric.PersonHash(),
			PersonPubk:      biometric.PersonPublicKey(),
			Signature: &api.Signature{
				Biod:   biometric.Signatures().BiodSig,
				Person: biometric.Signatures().PersonSig,
			},
		},
		Endorsement: buildAPIEndorsement(biometric.Endorsement()),
		Signature:   sig,
	}
}

func buildBiometricAPIData(data *account.BioData) *api.BiometricData {
	return &api.BiometricData{
		BiodPubk:        data.BiodPubk,
		CipherAddrBio:   data.CipherAddrBio,
		CipherAddrRobot: data.CipherAddrRobot,
		CipherAESKey:    data.CipherAESKey,
		PersonHash:      data.PersonHash,
		PersonPubk:      data.PersonPubk,
		Signature: &api.Signature{
			Biod:   data.Sigs.BiodSig,
			Person: data.Sigs.PersonSig,
		},
	}
}

func buildAPIEndorsement(end datamining.Endorsement) *api.Endorsement {
	valids := make([]*api.Validation, 0)
	for _, v := range end.Validations() {
		valids = append(valids, buildAPIValidation(v))
	}
	return &api.Endorsement{
		LastTransactionHash: end.LastTransactionHash(),
		MasterValidation:    buildAPIMasterValidation(end.MasterValidation()),
		TransactionHash:     end.TransactionHash(),
		Validations:         valids,
	}
}

func buildAPIMasterValidation(mv datamining.MasterValidation) *api.MasterValidation {
	return &api.MasterValidation{
		LastTransactionMiners: mv.LastTransactionMiners(),
		ProofOfWorkRobotKey:   mv.ProofOfWorkRobotKey(),
		ProofOfWorkValidation: buildAPIValidation(mv.ProofOfWorkValidation()),
	}
}

func buildAPIValidation(v datamining.Validation) *api.Validation {
	return &api.Validation{
		PublicKey: v.PublicKey(),
		Signature: v.Signature(),
		Status:    api.Validation_ValidationStatus(v.Status()),
		Timestamp: v.Timestamp().Unix(),
	}
}

func buildKeychainJSON(keychain account.Keychain) KeychainJSON {
	return KeychainJSON{
		Data: KeychainDataJSON{
			BiodPublicKey:      keychain.BiodPublicKey(),
			EncryptedAddrRobot: keychain.CipherAddrRobot(),
			EncryptedWallet:    keychain.CipherWallet(),
			PersonPublicKey:    keychain.PersonPublicKey(),
		},
		Endorsement: buildEndorsementJSON(keychain.Endorsement()),
	}
}

func buildBiometricJSON(biometric account.Biometric) BiometricJSON {
	return BiometricJSON{
		Data: BioDataJSON{
			BiodPublicKey:       biometric.BiodPublicKey(),
			EncryptedAddrRobot:  biometric.CipherAddrRobot(),
			EncryptedAddrPerson: biometric.CipherAddrBio(),
			PersonHash:          biometric.PersonHash(),
			EncryptedAESKey:     biometric.CipherAESKey(),
			PersonPublicKey:     biometric.PersonPublicKey(),
		},
		Endorsement: buildEndorsementJSON(biometric.Endorsement()),
	}
}

func buildEndorsementJSON(end datamining.Endorsement) EndorsementJSON {
	valids := make([]ValidationJSON, 0)
	for _, v := range end.Validations() {
		valids = append(valids, buildValidationJSON(v))
	}
	return EndorsementJSON{
		LastTransactionHash: end.LastTransactionHash(),
		TransactionHash:     end.TransactionHash(),
		MasterValidation:    buildMasterValidationJSON(end.MasterValidation()),
		Validations:         valids,
	}
}

func buildMasterValidationJSON(mv datamining.MasterValidation) MasterValidationJSON {
	return MasterValidationJSON{
		LastTransactionMiners: mv.LastTransactionMiners(),
		ProofOfWorkRobotKey:   mv.ProofOfWorkRobotKey(),
		ProofOfWorkValidation: buildValidationJSON(mv.ProofOfWorkValidation()),
	}
}

func buildValidationJSON(v datamining.Validation) ValidationJSON {
	return ValidationJSON{
		PublicKey: v.PublicKey(),
		Signature: v.Signature(),
		Status:    v.Status().String(),
		Timestamp: v.Timestamp().Unix(),
	}
}

func buildBiometricFromResponse(res *api.BiometricResponse) account.Biometric {
	return account.NewBiometric(buildBiometricDataFromAPI(res.Data), buildEndorsementFromAPI(res.Endorsement))
}

func buildBiometricDataFromAPI(data *api.BiometricData) *account.BioData {
	return &account.BioData{
		BiodPubk:        data.BiodPubk,
		CipherAddrBio:   data.CipherAddrBio,
		CipherAddrRobot: data.CipherAddrRobot,
		CipherAESKey:    data.CipherAESKey,
		PersonHash:      data.PersonHash,
		PersonPubk:      data.PersonPubk,
		Sigs: account.Signatures{
			BiodSig:   data.Signature.Biod,
			PersonSig: data.Signature.Person,
		},
	}
}

func buildKeychainFromResponse(res *api.KeychainResponse) account.Keychain {
	return account.NewKeychain(buildKeychainDataFromAPI(res.Data), buildEndorsementFromAPI(res.Endorsement))
}

func buildKeychainDataFromAPI(data *api.KeychainData) *account.KeyChainData {
	return &account.KeyChainData{
		BiodPubk:        data.BiodPubk,
		CipherAddrRobot: data.CipherAddrRobot,
		CipherWallet:    data.CipherWallet,
		PersonPubk:      data.PersonPubk,
		Sigs: account.Signatures{
			BiodSig:   data.Signature.Biod,
			PersonSig: data.Signature.Person,
		},
	}
}

func buildEndorsementFromAPI(end *api.Endorsement) datamining.Endorsement {
	valids := make([]datamining.Validation, 0)
	for _, v := range end.Validations {
		valids = append(valids, buildValidationFromAPI(v))
	}
	return datamining.NewEndorsement(
		end.LastTransactionHash,
		end.TransactionHash,
		buildMasterValidationFromAPI(end.MasterValidation),
		valids,
	)
}

func buildMasterValidationFromAPI(mv *api.MasterValidation) datamining.MasterValidation {
	return datamining.NewMasterValidation(
		mv.LastTransactionMiners,
		mv.ProofOfWorkRobotKey,
		buildValidationFromAPI(mv.ProofOfWorkValidation),
	)
}

func buildValidationFromAPI(valid *api.Validation) datamining.Validation {
	return datamining.NewValidation(
		datamining.ValidationStatus(valid.Status),
		time.Unix(valid.Timestamp, 0),
		valid.PublicKey,
		valid.Signature,
	)
}

package externalrpc

import (
	"time"

	api "github.com/uniris/uniris-core/datamining/api/protobuf-spec"
	"github.com/uniris/uniris-core/datamining/pkg"
	"github.com/uniris/uniris-core/datamining/pkg/account"
)

func buildEndorsementAPI(data datamining.Endorsement) *api.Endorsement {
	return &api.Endorsement{
		MasterValidation: &api.MasterValidation{
			LastTransactionMiners: data.MasterValidation().LastTransactionMiners(),
			ProofOfWorkRobotKey:   data.MasterValidation().ProofOfWorkRobotKey(),
			ProofOfWorkValidation: &api.Validation{
				PublicKey: data.MasterValidation().ProofOfWorkValidation().PublicKey(),
				Signature: data.MasterValidation().ProofOfWorkValidation().Signature(),
				Status:    api.Validation_ValidationStatus(data.MasterValidation().ProofOfWorkValidation().Status()),
				Timestamp: data.MasterValidation().ProofOfWorkValidation().Timestamp().Unix(),
			},
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

func buildKeychainAPIData(data *account.KeyChainData) *api.KeychainData {
	return &api.KeychainData{
		BiodPubk:        data.BiodPubk,
		CipherAddrRobot: data.CipherAddrRobot,
		CipherWallet:    data.CipherWallet,
		PersonPubk:      data.PersonPubk,
		Signature: &api.Signature{
			Biod:   data.Sigs.BiodSig,
			Person: data.Sigs.PersonSig,
		},
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
	return account.NewKeychain(buildKeychainDataFromAPI(res.Data, ""), buildEndorsementFromAPI(res.Endorsement))
}

func buildKeychainDataFromAPI(data *api.KeychainData, clearAddr string) *account.KeyChainData {
	return &account.KeyChainData{
		WalletAddr:      clearAddr,
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

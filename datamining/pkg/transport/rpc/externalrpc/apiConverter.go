package externalrpc

import (
	"time"

	api "github.com/uniris/uniris-core/datamining/api/protobuf-spec"
	"github.com/uniris/uniris-core/datamining/pkg"
	"github.com/uniris/uniris-core/datamining/pkg/account"
)

func createKeychainData(data *account.KeyChainData) *api.KeychainData {
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

func createBiometricData(data *account.BioData) *api.BiometricData {
	return &api.BiometricData{
		BiodPubk:        data.BiodPubk,
		CipherAddrRobot: data.CipherAddrRobot,
		CipherAddrBio:   data.CipherAddrBio,
		CipherAESKey:    data.CipherAESKey,
		PersonPubk:      data.PersonPubk,
		PersonHash:      data.PersonHash,
		Signature: &api.Signature{
			Biod:   data.Sigs.BiodSig,
			Person: data.Sigs.PersonSig,
		},
	}
}

func createEndorsement(data datamining.Endorsement) *api.Endorsement {
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

func formatKeychainDataAPI(data *api.KeychainData, addr string) *account.KeyChainData {
	return &account.KeyChainData{
		WalletAddr:      addr,
		BiodPubk:        data.BiodPubk,
		CipherAddrRobot: data.CipherAddrRobot,
		CipherWallet:    data.CipherWallet,
		PersonPubk:      data.PersonPubk,
		Sigs: datamining.Signatures{
			BiodSig:   data.Signature.Biod,
			PersonSig: data.Signature.Person,
		},
	}
}

func formatBiometricDataAPI(data *api.BiometricData) *account.BioData {
	return &account.BioData{
		BiodPubk:        data.BiodPubk,
		CipherAddrRobot: data.CipherAddrRobot,
		CipherAddrBio:   data.CipherAddrBio,
		CipherAESKey:    data.CipherAESKey,
		PersonPubk:      data.PersonPubk,
		PersonHash:      data.PersonHash,
		Sigs: datamining.Signatures{
			BiodSig:   data.Signature.Biod,
			PersonSig: data.Signature.Person,
		},
	}
}

func formatEndorsementAPI(end *api.Endorsement) datamining.Endorsement {
	valids := make([]datamining.Validation, 0)
	for _, v := range end.Validations {
		valids = append(valids, formatValidationAPI(v))
	}
	return datamining.NewEndorsement(
		end.LastTransactionHash,
		end.TransactionHash,
		formatMasterValidationAPI(end.MasterValidation),
		valids,
	)
}

func formatMasterValidationAPI(mv *api.MasterValidation) datamining.MasterValidation {
	return datamining.NewMasterValidation(
		mv.LastTransactionMiners,
		mv.ProofOfWorkRobotKey,
		formatValidationAPI(mv.ProofOfWorkValidation))
}

func formatValidationAPI(v *api.Validation) datamining.Validation {
	return datamining.NewValidation(
		datamining.ValidationStatus(v.Status),
		time.Unix(v.Timestamp, 0),
		v.PublicKey,
		v.Signature,
	)
}

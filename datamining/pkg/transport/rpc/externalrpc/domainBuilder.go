package externalrpc

import (
	"time"

	api "github.com/uniris/uniris-core/datamining/api/protobuf-spec"
	"github.com/uniris/uniris-core/datamining/pkg"
)

//BuillDomainWallet creates wallet from a GRPC response
func BuillDomainWallet(w *api.Wallet) *datamining.Wallet {
	return datamining.NewWallet(
		&datamining.WalletData{
			BiodPubk:        w.Data.BiodPublicKey,
			CipherAddrRobot: w.Data.EncryptedAddrRobot,
			CipherWallet:    w.Data.EncryptedWallet,
			EmPubk:          w.Data.PersonPubKey,
			Sigs: datamining.Signatures{
				BiodSig: w.Data.Signatures.Biod,
				EmSig:   w.Data.Signatures.Person,
			},
		},
		BuildDomainEndorsement(w.Endorsement),
		w.LastTransactionHash,
	)
}

//BuildWalletDataFromStorageRequest create data wallet from GRPC request
func BuildWalletDataFromStorageRequest(req *api.WalletStoreRequest) *datamining.Wallet {
	wData := &datamining.WalletData{
		BiodPubk:        req.Data.BiodPublicKey,
		EmPubk:          req.Data.PersonPubKey,
		CipherAddrRobot: req.Data.EncryptedAddrRobot,
		CipherWallet:    req.Data.EncryptedWallet,
		Sigs: datamining.Signatures{
			BiodSig: req.Data.Signatures.Biod,
			EmSig:   req.Data.Signatures.Person,
		},
	}

	return datamining.NewWallet(wData, BuildDomainEndorsement(req.Endorsement), req.LastTransactionHash)
}

//BuildWalletFromValidation create wallet data from wallet validation request
func BuildWalletFromValidation(req *api.WalletValidationRequest) *datamining.WalletData {
	return &datamining.WalletData{
		BiodPubk:        req.BiodPublicKey,
		EmPubk:          req.PersonPubKey,
		CipherAddrRobot: req.EncryptedAddrRobot,
		CipherWallet:    req.EncryptedWallet,
		Sigs: datamining.Signatures{
			BiodSig: req.Signatures.Biod,
			EmSig:   req.Signatures.Person,
		},
	}
}

//BuildDomainEndorsement create endorsment from GRPC request
func BuildDomainEndorsement(end *api.Endorsement) *datamining.Endorsement {

	valids := make([]datamining.Validation, 0)
	for _, v := range end.Validations {
		valids = append(valids, BuildDomainValidation(v))
	}

	return datamining.NewEndorsement(
		time.Unix(end.Timestamp, 0),
		end.TransactionHash,
		BuildDomainMasterValidation(end.MasterValidation),
		valids,
	)
}

//BuildDomainMasterValidation create a master validation for GRPC response
func BuildDomainMasterValidation(mv *api.MasterValidation) *datamining.MasterValidation {
	return datamining.NewMasterValidation(
		mv.LastTransactionMiners,
		mv.PowMasterKey,
		BuildDomainValidation(mv.PowValidation),
	)
}

//BuildDomainValidation create validation from GRPC request
func BuildDomainValidation(v *api.Validation) datamining.Validation {
	return datamining.NewValidation(
		datamining.ValidationStatus(v.Status),
		time.Unix(v.Timestamp, 0),
		v.PublicKey,
		v.Signature)
}

//BuilBioDataFromStoreRequest create bio wallet from GRPC request
func BuilBioDataFromStoreRequest(req *api.BioStorageRequest) *datamining.BioWallet {
	bioData := &datamining.BioData{
		BHash:           req.Data.BiometricHash,
		BiodPubk:        req.Data.BiodPubKey,
		EmPubk:          req.Data.PersonPubKey,
		CipherAddrBio:   req.Data.EncryptedAddrBiod,
		CipherAddrRobot: req.Data.EncryptedAddrRobot,
		CipherAESKey:    req.Data.EncryptedAESKey,
		Sigs: datamining.Signatures{
			BiodSig: req.Data.Signatures.Biod,
			EmSig:   req.Data.Signatures.Person,
		},
	}

	return datamining.NewBioWallet(bioData, BuildDomainEndorsement(req.Endorsement))
}

//BuildBioDataFromValidation create bio data from bio validation request
func BuildBioDataFromValidation(req *api.BioValidationRequest) *datamining.BioData {
	return &datamining.BioData{
		BHash:           req.BiometricHash,
		BiodPubk:        req.BiodPubKey,
		EmPubk:          req.PersonPubKey,
		CipherAddrBio:   req.EncryptedAddrBiod,
		CipherAddrRobot: req.EncryptedAddrRobot,
		CipherAESKey:    req.EncryptedAESKey,
		Sigs: datamining.Signatures{
			BiodSig: req.Signatures.Biod,
			EmSig:   req.Signatures.Person,
		},
	}
}

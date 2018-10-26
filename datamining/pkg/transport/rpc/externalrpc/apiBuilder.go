package externalrpc

import (
	api "github.com/uniris/uniris-core/datamining/api/protobuf-spec"
	datamining "github.com/uniris/uniris-core/datamining/pkg"
)

//BuildAPIWalletStoreRequest create a wallet storage request for GRPC response
func BuildAPIWalletStoreRequest(w *datamining.Wallet) *api.WalletStoreRequest {
	apiWallet := BuildAPIWallet(w)
	return &api.WalletStoreRequest{
		Data:                apiWallet.Data,
		Endorsement:         apiWallet.Endorsement,
		LastTransactionHash: apiWallet.LastTransactionHash,
	}
}

//BuildAPIBioStoreRequest create a wallet storage request for GRPC response
func BuildAPIBioStoreRequest(bw *datamining.BioWallet) *api.BioStorageRequest {
	return &api.BioStorageRequest{
		Data:        BuildAPIBioData(bw),
		Endorsement: BuildAPIEndorsement(bw.Endorsement()),
	}
}

//BuildAPIWallet create a wallet for GRPC response
func BuildAPIWallet(w *datamining.Wallet) *api.Wallet {
	return &api.Wallet{
		Data: &api.WalletData{
			BiodPublicKey:      w.BiodPublicKey(),
			EncryptedAddrRobot: w.CipherAddrRobot(),
			EncryptedWallet:    w.CipherWallet(),
			PersonPubKey:       w.PersonPublicKey(),
			Signatures: &api.Signature{
				Biod:   w.Signatures().BiodSig,
				Person: w.Signatures().EmSig,
			},
		},
		Endorsement:         BuildAPIEndorsement(w.Endorsement()),
		LastTransactionHash: w.OldTransactionHash(),
	}
}

//BuildAPIBioData create a bio data for GRPC response
func BuildAPIBioData(b *datamining.BioWallet) *api.BioData {
	return &api.BioData{
		BiodPubKey:         b.BiodPublicKey(),
		BiometricHash:      b.Bhash(),
		EncryptedAddrBiod:  b.CipherAddrBio(),
		EncryptedAddrRobot: b.CipherAddrRobot(),
		EncryptedAESKey:    b.CipherAESKey(),
		PersonPubKey:       b.PersonPublicKey(),
		Signatures: &api.Signature{
			Biod:   b.Signatures().BiodSig,
			Person: b.Signatures().EmSig,
		},
	}
}

//BuildAPIMasterValidation create a master validation for GRPC response
func BuildAPIMasterValidation(mv *datamining.MasterValidation) *api.MasterValidation {
	return &api.MasterValidation{
		LastTransactionMiners: mv.LastTransactionMiners(),
		PowMasterKey:          mv.ProofOfWorkRobotKey(),
		PowValidation:         BuildAPIValidation(mv.ProofOfWorkValidation()),
	}
}

//BuildAPIEndorsement create a endorsement for GRPC response
func BuildAPIEndorsement(end *datamining.Endorsement) *api.Endorsement {
	valids := make([]*api.Validation, 0)
	for _, v := range end.Validations() {
		valids = append(valids, BuildAPIValidation(v))
	}

	return &api.Endorsement{
		MasterValidation: BuildAPIMasterValidation(end.MasterValidation()),
		Timestamp:        end.Timestamp().Unix(),
		TransactionHash:  end.TransactionHash(),
		Validations:      valids,
	}
}

//BuildAPIValidation creates a validation for GRPC response
func BuildAPIValidation(v datamining.Validation) *api.Validation {
	return &api.Validation{
		PublicKey: v.PublicKey(),
		Signature: v.Signature(),
		Status:    api.Validation_ValidationStatus(v.Status()),
		Timestamp: v.Timestamp().Unix(),
	}
}

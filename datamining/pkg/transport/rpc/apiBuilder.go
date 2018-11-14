package rpc

import (
	api "github.com/uniris/uniris-core/datamining/api/protobuf-spec"
	"github.com/uniris/uniris-core/datamining/pkg/account"
	"github.com/uniris/uniris-core/datamining/pkg/mining"
)

type apiBuilder struct{}

func (b apiBuilder) buildBiometricData(data account.BiometricData) *api.BiometricData {
	return &api.BiometricData{
		CipherAddrBio:   data.CipherAddrPerson(),
		CipherAddrRobot: data.CipherAddrRobot(),
		CipherAESKey:    data.CipherAESKey(),
		PersonHash:      data.PersonHash(),
		PersonPubk:      data.PersonPublicKey(),
		Signature: &api.Signature{
			Biod:   data.Signatures().Biod(),
			Person: data.Signatures().Person(),
		},
	}
}

func (b apiBuilder) buildKeychainData(data account.KeychainData) *api.KeychainData {
	return &api.KeychainData{
		CipherAddrRobot: data.CipherAddrRobot(),
		CipherWallet:    data.CipherWallet(),
		PersonPubk:      data.PersonPublicKey(),
		Signature: &api.Signature{
			Biod:   data.Signatures().Biod(),
			Person: data.Signatures().Person(),
		},
	}
}

func buildKeychainDataForAPI(data account.KeychainData) *api.KeychainData {
	return &api.KeychainData{
		CipherAddrRobot: data.CipherAddrRobot(),
		CipherWallet:    data.CipherWallet(),
		PersonPubk:      data.PersonPublicKey(),
		Signature: &api.Signature{
			Biod:   data.Signatures().Biod(),
			Person: data.Signatures().Person(),
		},
	}
}

func (b apiBuilder) buildEndorsement(end mining.Endorsement) *api.Endorsement {
	valids := make([]*api.Validation, 0)
	for _, v := range end.Validations() {
		valids = append(valids, b.buildValidation(v))
	}
	return &api.Endorsement{
		LastTransactionHash: end.LastTransactionHash(),
		MasterValidation:    b.buildMasterValidation(end.MasterValidation()),
		TransactionHash:     end.TransactionHash(),
		Validations:         valids,
	}
}

func (b apiBuilder) buildMasterValidation(mv mining.MasterValidation) *api.MasterValidation {
	return &api.MasterValidation{
		LastTransactionMiners: mv.LastTransactionMiners(),
		ProofOfWorkKey:        mv.ProofOfWorkKey(),
		ProofOfWorkValidation: b.buildValidation(mv.ProofOfWorkValidation()),
	}
}

func (b apiBuilder) buildValidation(v mining.Validation) *api.Validation {
	return &api.Validation{
		PublicKey: v.PublicKey(),
		Signature: v.Signature(),
		Status:    api.Validation_ValidationStatus(v.Status()),
		Timestamp: v.Timestamp().Unix(),
	}
}

package rpc

import (
	api "github.com/uniris/uniris-core/datamining/api/protobuf-spec"
	"github.com/uniris/uniris-core/datamining/pkg/account"
	"github.com/uniris/uniris-core/datamining/pkg/mining"
)

type apiBuilder struct{}

func (b apiBuilder) buildID(id account.ID) *api.ID {
	return &api.ID{
		EncryptedAddrByID:    id.EncryptedAddrByID(),
		EncryptedAddrByRobot: id.EncryptedAddrByRobot(),
		EncryptedAESKey:      id.EncryptedAESKey(),
		Hash:                 id.Hash(),
		PublicKey:            id.PublicKey(),
		IDSignature:          id.IDSignature(),
		EmitterSignature:     id.EmitterSignature(),
		Proposal: &api.Proposal{
			SharedEmitterKeyPair: &api.KeyPairProposal{
				EncryptedPrivateKey: id.Proposal().SharedEmitterKeyPair().EncryptedPrivateKey(),
				PublicKey:           id.Proposal().SharedEmitterKeyPair().PublicKey(),
			},
		},
	}
}

func (b apiBuilder) buildKeychain(kc account.Keychain) *api.Keychain {
	return &api.Keychain{
		EncryptedAddrByRobot: kc.EncryptedAddrByRobot(),
		EncryptedWallet:      kc.EncryptedWallet(),
		IDPublicKey:          kc.IDPublicKey(),
		EmitterSignature:     kc.EmitterSignature(),
		IDSignature:          kc.IDSignature(),
		Proposal: &api.Proposal{
			SharedEmitterKeyPair: &api.KeyPairProposal{
				EncryptedPrivateKey: kc.Proposal().SharedEmitterKeyPair().EncryptedPrivateKey(),
				PublicKey:           kc.Proposal().SharedEmitterKeyPair().PublicKey(),
			},
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

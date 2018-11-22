package rpc

import (
	"time"

	api "github.com/uniris/uniris-core/datamining/api/protobuf-spec"
	"github.com/uniris/uniris-core/datamining/pkg/account"
	"github.com/uniris/uniris-core/datamining/pkg/mining"
)

type dataBuilder struct{}

func (b dataBuilder) buildID(id *api.ID) account.ID {
	return account.NewID(id.Hash,
		id.EncryptedAddrByRobot,
		id.EncryptedAddrByID,
		id.EncryptedAESKey,
		id.PublicKey,
		id.IDSignature,
		id.EmitterSignature)
}

func (b dataBuilder) buildKeychain(kc *api.Keychain) account.Keychain {
	return account.NewKeychain(
		kc.EncryptedAddrByRobot,
		kc.EncryptedWallet,
		kc.IDPublicKey,
		kc.IDSignature,
		kc.EmitterSignature,
	)
}

func (b dataBuilder) buildEndorsement(data *api.Endorsement) mining.Endorsement {
	valids := make([]mining.Validation, 0)
	for _, v := range data.Validations {
		valids = append(valids, b.buildValidation(v))
	}

	return mining.NewEndorsement(
		data.LastTransactionHash,
		data.TransactionHash,
		b.buildMasterValidation(data.MasterValidation),
		valids,
	)
}

func (b dataBuilder) buildMasterValidation(mv *api.MasterValidation) mining.MasterValidation {
	return mining.NewMasterValidation(
		mv.LastTransactionMiners,
		mv.ProofOfWorkKey,
		b.buildValidation(mv.ProofOfWorkValidation),
	)
}

func (b dataBuilder) buildValidation(v *api.Validation) mining.Validation {
	return mining.NewValidation(
		mining.ValidationStatus(v.Status),
		time.Unix(v.Timestamp, 0),
		v.PublicKey,
		v.Signature)
}

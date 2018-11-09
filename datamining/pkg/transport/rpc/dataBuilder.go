package rpc

import (
	"time"

	api "github.com/uniris/uniris-core/datamining/api/protobuf-spec"
	"github.com/uniris/uniris-core/datamining/pkg/account"
	"github.com/uniris/uniris-core/datamining/pkg/mining"
)

type dataBuilder struct{}

func (b dataBuilder) buildBiometricData(data *api.BiometricData) account.BiometricData {
	sigs := account.NewSignatures(data.Signature.Biod, data.Signature.Person)
	return account.NewBiometricData(
		data.PersonHash,
		data.CipherAddrRobot,
		data.CipherAddrBio,
		data.CipherAESKey,
		data.PersonPubk,
		data.BiodPubk,
		sigs)
}

func (b dataBuilder) buildKeychainData(data *api.KeychainData) account.KeychainData {
	sigs := account.NewSignatures(data.Signature.Biod, data.Signature.Person)
	return account.NewKeychainData(
		data.CipherAddrRobot,
		data.CipherWallet,
		data.PersonPubk,
		data.BiodPubk,
		sigs,
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
		mv.ProofOfWorkRobotKey,
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

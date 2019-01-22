package rpc

import (
	"time"

	api "github.com/uniris/uniris-core/api/protobuf-spec"
	uniris "github.com/uniris/uniris-core/pkg"
)

func formatTransaction(tx *api.Transaction) uniris.Transaction {
	prop := uniris.NewTransactionProposal(
		uniris.NewSharedKeyPair(tx.Proposal.SharedEmitterKeys.EncryptedPrivateKey, tx.Proposal.SharedEmitterKeys.PublicKey),
	)

	return uniris.NewTransactionBase(tx.Address, uniris.TransactionType(tx.Type), tx.Data,
		time.Unix(tx.Timestamp, 0),
		tx.PublicKey,
		tx.Signature,
		tx.EmitterSignature,
		prop,
		tx.TransactionHash)
}

func formatMinedTransaction(tx *api.Transaction, mv *api.MasterValidation, valids []*api.MinerValidation) uniris.Transaction {
	masterValidation := uniris.NewMasterValidation(mv.PreviousTransactionMiners, mv.ProofOfWork, formatValidation(mv.PreValidation))

	confValids := make([]uniris.MinerValidation, 0)
	for _, v := range valids {
		confValids = append(confValids, formatValidation(v))
	}

	return uniris.NewMinedTransaction(formatTransaction(tx), masterValidation, confValids)
}

func formatAPIValidation(v uniris.MinerValidation) *api.MinerValidation {
	return &api.MinerValidation{
		PublicKey: v.MinerPublicKey(),
		Signature: v.MinerSignature(),
		Status:    api.MinerValidation_ValidationStatus(v.Status()),
		Timestamp: v.Timestamp().Unix(),
	}
}

func formatValidation(v *api.MinerValidation) uniris.MinerValidation {
	return uniris.NewMinerValidation(uniris.ValidationStatus(v.Status), time.Unix(v.Timestamp, 0), v.PublicKey, v.Signature)
}

func formatAPITransaction(tx uniris.Transaction) *api.Transaction {
	return &api.Transaction{
		Address:          tx.Address(),
		Data:             tx.Data(),
		Type:             api.TransactionType(tx.Type()),
		PublicKey:        tx.PublicKey(),
		Signature:        tx.Signature(),
		EmitterSignature: tx.EmitterSignature(),
		Timestamp:        tx.Timestamp().Unix(),
		TransactionHash:  tx.TransactionHash(),
		Proposal: &api.TransactionProposal{
			SharedEmitterKeys: &api.SharedKeys{
				EncryptedPrivateKey: tx.Proposal().SharedEmitterKeyPair().EncryptedPrivateKey(),
				PublicKey:           tx.Proposal().SharedEmitterKeyPair().PublicKey(),
			},
		},
	}
}

func formatAPIMasterValidationAPI(masterValid uniris.MasterValidation) *api.MasterValidation {
	return &api.MasterValidation{
		ProofOfWork:               masterValid.ProofOfWork(),
		PreviousTransactionMiners: masterValid.PreviousTransactionMiners(),
		PreValidation: &api.MinerValidation{
			PublicKey: masterValid.Validation().MinerPublicKey(),
			Signature: masterValid.Validation().MinerSignature(),
			Status:    api.MinerValidation_ValidationStatus(masterValid.Validation().Status()),
			Timestamp: masterValid.Validation().Timestamp().Unix(),
		},
	}
}

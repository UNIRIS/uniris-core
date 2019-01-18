package rpc

import (
	"context"
	"fmt"
	"time"

	uniris "github.com/uniris/uniris-core/pkg"
	"github.com/uniris/uniris-core/pkg/locking"

	api "github.com/uniris/uniris-core/api/protobuf-spec"
	"github.com/uniris/uniris-core/pkg/adding"
	"github.com/uniris/uniris-core/pkg/listing"
	"github.com/uniris/uniris-core/pkg/mining"
)

type transactionSrv struct {
	adder  adding.Service
	lister listing.Service
	miner  mining.Service
	locker locking.Service
}

//NewTransactionServer creates a new GRPC transaction server
func NewTransactionServer() api.TransactionServiceServer {
	return transactionSrv{}
}

func (s transactionSrv) GetTransactionStatus(ctx context.Context, req *api.TransactionStatusRequest) (*api.TransactionStatusResponse, error) {
	fmt.Printf("GET TRANSACTION STATUS REQUEST - %s\n", time.Unix(req.Timestamp, 0).String())

	//TODO: verify signature request

	status, err := s.lister.GetTransactionStatus(req.TransactionHash)
	if err != nil {
		return nil, err
	}

	return &api.TransactionStatusResponse{
		Status:            api.TransactionStatusResponse_TransactionStatus(status),
		Timestamp:         time.Now().Unix(),
		SignatureResponse: "", //TODO
	}, nil
}

func (s transactionSrv) LockTransaction(ctx context.Context, req *api.LockRequest) (*api.LockResponse, error) {
	fmt.Printf("LOCK TRANSACTION REQUEST - %s\n", time.Unix(req.Timestamp, 0).String())

	//TODO: verify signature request

	if err := s.locker.LockTransaction(req.TransactionHash, req.Address, req.MasterPeerIp); err != nil {
		return nil, err
	}

	return &api.LockResponse{
		LockHash:          "", //TODO
		Timestamp:         time.Now().Unix(),
		SignatureResponse: "", //TODO
	}, nil
}

func (s transactionSrv) PreValidateTransaction(ctx context.Context, req *api.PreValidationRequest) (*api.PreValidationResponse, error) {
	fmt.Printf("PRE VALIDATE TRANSACTION REQUEST - %s\n", time.Unix(req.Timestamp, 0).String())

	//TODO: verify signature request

	s.miner.LeadTransactionValidation(formatTransaction(req.Transaction), int(req.MinimumValidations))

	return &api.PreValidationResponse{
		PreValidationHash: "", //TODO
		Timestamp:         time.Now().Unix(),
		SignatureResponse: "", //TODO
	}, nil
}

func (s transactionSrv) ConfirmTransactionValidation(ctx context.Context, req *api.ConfirmValidationRequest) (*api.ConfirmValidationResponse, error) {
	fmt.Printf("CONFIRM VALIDATION TRANSACTION REQUEST - %s\n", time.Unix(req.Timestamp, 0).String())

	//TODO: verify signature request

	valid, err := s.miner.ConfirmTransactionValidation(formatMinedTransaction(req.Transaction, req.MasterValidation, nil))
	if err != nil {
		return nil, err
	}

	return &api.ConfirmValidationResponse{
		Validation:        formatAPIValidation(valid),
		Timestamp:         time.Now().Unix(),
		SignatureResponse: "", //TODO
	}, nil
}

func (s transactionSrv) StoreTransaction(ctx context.Context, req *api.StoreRequest) (*api.StoreResponse, error) {
	fmt.Printf("STORE TRANSACTION REQUEST - %s\n", time.Unix(req.Timestamp, 0).String())

	//TODO: verify signature request

	tx := formatMinedTransaction(req.MinedTransaction.Transaction, req.MinedTransaction.MasterValidation, req.MinedTransaction.ConfirmValidations)
	if err := s.adder.StoreTransaction(tx); err != nil {
		return nil, err
	}

	return &api.StoreResponse{
		MinedTransactionHash: "", //TODO
		SignatureResponse:    "", //TODO
		Timestamp:            time.Now().Unix(),
	}, nil
}

func formatTransaction(tx *api.Transaction) uniris.Transaction {
	prop := uniris.NewProposal(
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

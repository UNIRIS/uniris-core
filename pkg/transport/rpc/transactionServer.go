package rpc

import (
	"context"
	"fmt"
	"time"

	uniris "github.com/uniris/uniris-core/pkg"

	api "github.com/uniris/uniris-core/api/protobuf-spec"
	"github.com/uniris/uniris-core/pkg/adding"
	"github.com/uniris/uniris-core/pkg/listing"
	"github.com/uniris/uniris-core/pkg/mining"
)

type transactionSrv struct {
	adder        adding.Service
	lister       listing.Service
	miner        mining.Service
	sigHandler   SignatureHandler
	sharedPubKey string
	sharedPvKey  string
}

//NewTransactionServer creates a new GRPC transaction server
func NewTransactionServer(a adding.Service, l listing.Service, m mining.Service, sigHandler SignatureHandler, sharedPubk, sharedPvk string) api.TransactionServiceServer {
	return transactionSrv{
		adder:        a,
		lister:       l,
		miner:        m,
		sigHandler:   sigHandler,
		sharedPubKey: sharedPubk,
		sharedPvKey:  sharedPvk,
	}
}

func (s transactionSrv) GetTransactionStatus(ctx context.Context, req *api.TransactionStatusRequest) (*api.TransactionStatusResponse, error) {
	fmt.Printf("GET TRANSACTION STATUS REQUEST - %s\n", time.Unix(req.Timestamp, 0).String())

	if err := s.sigHandler.VerifyTransactionStatusRequestSignature(req, s.sharedPubKey); err != nil {
		return nil, err
	}

	status, err := s.lister.GetTransactionStatus(req.TransactionHash)
	if err != nil {
		return nil, err
	}

	res := &api.TransactionStatusResponse{
		Status:    api.TransactionStatusResponse_TransactionStatus(status),
		Timestamp: time.Now().Unix(),
	}
	if err := s.sigHandler.SignTransactionStatusResponse(res, s.sharedPvKey); err != nil {
		return nil, err
	}

	return res, nil
}

func (s transactionSrv) LockTransaction(ctx context.Context, req *api.LockRequest) (*api.LockResponse, error) {
	fmt.Printf("LOCK TRANSACTION REQUEST - %s\n", time.Unix(req.Timestamp, 0).String())

	if err := s.sigHandler.VerifyLockRequestSignature(req, s.sharedPubKey); err != nil {
		return nil, err
	}

	if err := s.adder.StoreLock(uniris.NewLock(req.TransactionHash, req.Address, req.MasterPeerIp)); err != nil {
		return nil, err
	}

	res := &api.LockResponse{
		Timestamp: time.Now().Unix(),
	}
	if err := s.sigHandler.SignLockResponse(res, s.sharedPvKey); err != nil {
		return nil, err
	}
	return res, nil
}

func (s transactionSrv) PreValidateTransaction(ctx context.Context, req *api.PreValidationRequest) (*api.PreValidationResponse, error) {
	fmt.Printf("PRE VALIDATE TRANSACTION REQUEST - %s\n", time.Unix(req.Timestamp, 0).String())

	if err := s.sigHandler.VerifyPreValidateRequestSignature(req, s.sharedPubKey); err != nil {
		return nil, err
	}

	s.miner.LeadTransactionValidation(formatTransaction(req.Transaction), int(req.MinimumValidations))

	res := &api.PreValidationResponse{
		Timestamp: time.Now().Unix(),
	}
	if err := s.sigHandler.SignPreValidationResponse(res, s.sharedPvKey); err != nil {
		return nil, err
	}
	return res, nil
}

func (s transactionSrv) ConfirmTransactionValidation(ctx context.Context, req *api.ConfirmValidationRequest) (*api.ConfirmValidationResponse, error) {
	fmt.Printf("CONFIRM VALIDATION TRANSACTION REQUEST - %s\n", time.Unix(req.Timestamp, 0).String())

	if err := s.sigHandler.VerifyConfirmValidationRequestSignature(req, s.sharedPubKey); err != nil {
		return nil, err
	}

	valid, err := s.miner.ConfirmTransactionValidation(formatMinedTransaction(req.Transaction, req.MasterValidation, nil))
	if err != nil {
		return nil, err
	}

	res := &api.ConfirmValidationResponse{
		Validation: formatAPIValidation(valid),
		Timestamp:  time.Now().Unix(),
	}
	if err := s.sigHandler.SignConfirmValidationResponse(res, s.sharedPvKey); err != nil {
		return nil, err
	}
	return res, nil
}

func (s transactionSrv) StoreTransaction(ctx context.Context, req *api.StoreRequest) (*api.StoreResponse, error) {
	fmt.Printf("STORE TRANSACTION REQUEST - %s\n", time.Unix(req.Timestamp, 0).String())

	if err := s.sigHandler.VerifyStoreRequest(req, s.sharedPubKey); err != nil {
		return nil, err
	}

	tx := formatMinedTransaction(req.MinedTransaction.Transaction, req.MinedTransaction.MasterValidation, req.MinedTransaction.ConfirmValidations)
	if err := s.adder.StoreTransaction(tx); err != nil {
		return nil, err
	}

	res := &api.StoreResponse{
		Timestamp: time.Now().Unix(),
	}
	if err := s.sigHandler.SignStoreResponse(res, s.sharedPvKey); err != nil {
		return nil, err
	}
	return res, nil
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

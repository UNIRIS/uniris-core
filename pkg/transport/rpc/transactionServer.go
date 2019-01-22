package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	uniris "github.com/uniris/uniris-core/pkg"

	api "github.com/uniris/uniris-core/api/protobuf-spec"
	"github.com/uniris/uniris-core/pkg/adding"
	"github.com/uniris/uniris-core/pkg/crypto"
	"github.com/uniris/uniris-core/pkg/listing"
	"github.com/uniris/uniris-core/pkg/mining"
)

type transactionSrv struct {
	adder        adding.Service
	lister       listing.Service
	miner        mining.Service
	sharedPubKey string
	sharedPvKey  string
}

//NewTransactionServer creates a new GRPC transaction server
func NewTransactionServer(a adding.Service, l listing.Service, m mining.Service, sharedPubk, sharedPvk string) api.TransactionServiceServer {
	return transactionSrv{
		adder:        a,
		lister:       l,
		miner:        m,
		sharedPubKey: sharedPubk,
		sharedPvKey:  sharedPvk,
	}
}

func (s transactionSrv) GetTransactionStatus(ctx context.Context, req *api.TransactionStatusRequest) (*api.TransactionStatusResponse, error) {
	fmt.Printf("GET TRANSACTION STATUS REQUEST - %s\n", time.Unix(req.Timestamp, 0).String())

	reqBytes, err := json.Marshal(&api.TransactionStatusRequest{
		TransactionHash: req.TransactionHash,
		Timestamp:       req.Timestamp,
	})
	if err != nil {
		return nil, err
	}
	if err := crypto.VerifySignature(string(reqBytes), s.sharedPubKey, req.SignatureRequest); err != nil {
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
	resBytes, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}
	sig, err := crypto.Sign(string(resBytes), s.sharedPvKey)
	if err != nil {
		return nil, err
	}
	res.SignatureResponse = sig

	return res, nil
}

func (s transactionSrv) LockTransaction(ctx context.Context, req *api.LockRequest) (*api.LockResponse, error) {
	fmt.Printf("LOCK TRANSACTION REQUEST - %s\n", time.Unix(req.Timestamp, 0).String())

	reqBytes, err := json.Marshal(&api.LockRequest{
		TransactionHash: req.TransactionHash,
		MasterPeerIp:    req.MasterPeerIp,
		Timestamp:       req.Timestamp,
	})
	if err := crypto.VerifySignature(string(reqBytes), s.sharedPubKey, req.SignatureRequest); err != nil {
		return nil, err
	}

	if err := s.adder.StoreLock(uniris.NewLock(req.TransactionHash, req.Address, req.MasterPeerIp)); err != nil {
		return nil, err
	}

	res := &api.LockResponse{
		Timestamp: time.Now().Unix(),
	}
	resBytes, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}
	sig, err := crypto.Sign(string(resBytes), s.sharedPvKey)
	if err != nil {
		return nil, err
	}
	res.SignatureResponse = sig
	return res, nil
}

func (s transactionSrv) PreValidateTransaction(ctx context.Context, req *api.PreValidationRequest) (*api.PreValidationResponse, error) {
	fmt.Printf("PRE VALIDATE TRANSACTION REQUEST - %s\n", time.Unix(req.Timestamp, 0).String())

	reqBytes, err := json.Marshal(&api.PreValidationRequest{
		Transaction:        req.Transaction,
		MinimumValidations: req.MinimumValidations,
		Timestamp:          req.Timestamp,
	})
	if err := crypto.VerifySignature(string(reqBytes), s.sharedPubKey, req.SignatureRequest); err != nil {
		return nil, err
	}

	s.miner.LeadTransactionValidation(formatTransaction(req.Transaction), int(req.MinimumValidations))

	res := &api.PreValidationResponse{
		Timestamp: time.Now().Unix(),
	}
	resBytes, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}
	sig, err := crypto.Sign(string(resBytes), s.sharedPvKey)
	if err != nil {
		return nil, err
	}
	res.SignatureResponse = sig
	return res, nil
}

func (s transactionSrv) ConfirmTransactionValidation(ctx context.Context, req *api.ConfirmValidationRequest) (*api.ConfirmValidationResponse, error) {
	fmt.Printf("CONFIRM VALIDATION TRANSACTION REQUEST - %s\n", time.Unix(req.Timestamp, 0).String())

	reqBytes, err := json.Marshal(&api.ConfirmValidationRequest{
		Transaction:      req.Transaction,
		MasterValidation: req.MasterValidation,
		Timestamp:        req.Timestamp,
	})
	if err != nil {
		return nil, err
	}

	if err := crypto.VerifySignature(string(reqBytes), s.sharedPubKey, req.SignatureRequest); err != nil {
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
	resBytes, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}
	sig, err := crypto.Sign(string(resBytes), s.sharedPvKey)
	if err != nil {
		return nil, err
	}
	res.SignatureResponse = sig
	return res, nil
}

func (s transactionSrv) StoreTransaction(ctx context.Context, req *api.StoreRequest) (*api.StoreResponse, error) {
	fmt.Printf("STORE TRANSACTION REQUEST - %s\n", time.Unix(req.Timestamp, 0).String())

	reqBytes, err := json.Marshal(&api.StoreRequest{
		MinedTransaction: req.MinedTransaction,
		Timestamp:        req.Timestamp,
	})
	if err := crypto.VerifySignature(string(reqBytes), s.sharedPubKey, req.SignatureRequest); err != nil {
		return nil, err
	}

	tx := formatMinedTransaction(req.MinedTransaction.Transaction, req.MinedTransaction.MasterValidation, req.MinedTransaction.ConfirmValidations)
	if err := s.adder.StoreTransaction(tx); err != nil {
		return nil, err
	}

	res := &api.StoreResponse{
		Timestamp: time.Now().Unix(),
	}
	resBytes, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}
	sig, err := crypto.Sign(string(resBytes), s.sharedPvKey)
	if err != nil {
		return nil, err
	}
	res.SignatureResponse = sig
	return res, nil
}

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

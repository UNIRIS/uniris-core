package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	api "github.com/uniris/uniris-core/api/protobuf-spec"
	"github.com/uniris/uniris-core/pkg/crypto"
	"github.com/uniris/uniris-core/pkg/transaction"
)

type transactionSrv struct {
	sharedPubKey string
	sharedPvKey  string
	storeSrv     transaction.StorageService
	lockSrv      transaction.LockService
	miningSrv    transaction.MiningService
}

//NewTransactionServer creates a new GRPC transaction server
func NewTransactionServer(storeSrv transaction.StorageService, lockSrv transaction.LockService, miningSrv transaction.MiningService, sharedPubk, sharedPvk string) api.TransactionServiceServer {
	return transactionSrv{
		storeSrv:     storeSrv,
		lockSrv:      lockSrv,
		miningSrv:    miningSrv,
		sharedPubKey: sharedPubk,
		sharedPvKey:  sharedPvk,
	}
}

func (s transactionSrv) GetLastTransaction(ctx context.Context, req *api.LastTransactionRequest) (*api.LastTransactionResponse, error) {
	fmt.Printf("GET LAST TRANSACTION REQUEST - %s\n", time.Unix(req.Timestamp, 0).String())

	reqBytes, err := json.Marshal(&api.LastTransactionRequest{
		TransactionAddress: req.TransactionAddress,
		Type:               req.Type,
		Timestamp:          req.Timestamp,
	})
	if err != nil {
		return nil, err
	}
	if err := crypto.VerifySignature(string(reqBytes), s.sharedPubKey, req.SignatureRequest); err != nil {
		return nil, err
	}

	tx, err := s.storeSrv.GetLastTransaction(req.TransactionAddress, transaction.Type(req.Type))
	if err != nil {
		return nil, err
	}

	res := &api.LastTransactionResponse{
		Timestamp:   time.Now().Unix(),
		Transaction: formatAPITransaction(*tx),
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

	status, err := s.storeSrv.GetTransactionStatus(req.TransactionHash)
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

	lock, err := transaction.NewLock(req.TransactionHash, req.Address, req.MasterPeerIp)
	if err != nil {
		return nil, err
	}
	if err := s.lockSrv.StoreLock(lock); err != nil {
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

	tx, err := formatTransaction(req.Transaction)
	if err != nil {
		return nil, err
	}
	s.miningSrv.LeadTransactionValidation(tx, int(req.MinimumValidations))

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

	tx, err := formatTransaction(req.Transaction)
	if err != nil {
		return nil, err
	}
	masterValid, err := formatMasterValidation(req.MasterValidation)
	if err != nil {
		return nil, err
	}
	valid, err := s.miningSrv.ValidateTransaction(tx, masterValid)
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

	tx, err := formatMinedTransaction(req.MinedTransaction.Transaction, req.MinedTransaction.MasterValidation, req.MinedTransaction.ConfirmValidations)
	if err != nil {
		return nil, err
	}

	if err := s.storeSrv.StoreTransaction(tx); err != nil {
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

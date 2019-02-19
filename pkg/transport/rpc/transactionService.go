package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	api "github.com/uniris/uniris-core/api/protobuf-spec"
	"github.com/uniris/uniris-core/pkg/chain"
	"github.com/uniris/uniris-core/pkg/consensus"
	"github.com/uniris/uniris-core/pkg/crypto"
	"github.com/uniris/uniris-core/pkg/shared"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type txSrv struct {
	chainDB        chain.Database
	locker         chain.Locker
	techDB         shared.TechDatabaseReader
	poolR          consensus.PoolRequester
	nodePublicKey  string
	nodePrivateKey string
}

//NewTransactionService creates service handler for the GRPC Transaction service
func NewTransactionService(cDB chain.Database, l chain.Locker, tDB shared.TechDatabaseReader, pR consensus.PoolRequester, nodePublicKeyk, nodePrivateKeyk string) api.TransactionServiceServer {
	return txSrv{
		chainDB:        cDB,
		locker:         l,
		techDB:         tDB,
		poolR:          pR,
		nodePublicKey:  nodePublicKeyk,
		nodePrivateKey: nodePrivateKeyk,
	}
}

func (s txSrv) GetLastTransaction(ctx context.Context, req *api.GetLastTransactionRequest) (*api.GetLastTransactionResponse, error) {
	fmt.Printf("GET LAST TRANSACTION REQUEST - %s\n", time.Unix(req.Timestamp, 0).String())

	reqBytes, err := json.Marshal(&api.GetLastTransactionRequest{
		TransactionAddress: req.TransactionAddress,
		Type:               req.Type,
		Timestamp:          req.Timestamp,
	})
	if err != nil {
		return nil, status.New(codes.InvalidArgument, err.Error()).Err()
	}

	nodeLastKeys, err := s.techDB.NodeLastKeys()
	if err != nil {
		return nil, err
	}

	if err := crypto.VerifySignature(string(reqBytes), nodeLastKeys.PublicKey(), req.SignatureRequest); err != nil {
		return nil, status.New(codes.InvalidArgument, err.Error()).Err()
	}

	tx, err := chain.LastTransaction(s.chainDB, req.TransactionAddress, chain.TransactionType(req.Type))
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	if tx == nil {
		return nil, status.New(codes.NotFound, "transaction does not exist").Err()
	}

	res := &api.GetLastTransactionResponse{
		Timestamp:   time.Now().Unix(),
		Transaction: formatAPITransaction(*tx),
	}
	resBytes, err := json.Marshal(res)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	sig, err := crypto.Sign(string(resBytes), nodeLastKeys.PrivateKey())
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	res.SignatureResponse = sig

	return res, nil
}

func (s txSrv) GetTransactionStatus(ctx context.Context, req *api.GetTransactionStatusRequest) (*api.GetTransactionStatusResponse, error) {
	fmt.Printf("GET TRANSACTION STATUS REQUEST - %s\n", time.Unix(req.Timestamp, 0).String())

	reqBytes, err := json.Marshal(&api.GetTransactionStatusRequest{
		TransactionHash: req.TransactionHash,
		Timestamp:       req.Timestamp,
	})
	if err != nil {
		return nil, status.New(codes.InvalidArgument, err.Error()).Err()
	}
	nodeLastKeys, err := s.techDB.NodeLastKeys()
	if err != nil {
		return nil, err
	}
	if err := crypto.VerifySignature(string(reqBytes), nodeLastKeys.PublicKey(), req.SignatureRequest); err != nil {
		return nil, status.New(codes.InvalidArgument, err.Error()).Err()
	}

	txStatus, err := chain.GetTransactionStatus(s.chainDB, req.TransactionHash)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	res := &api.GetTransactionStatusResponse{
		Status:    api.TransactionStatus(txStatus),
		Timestamp: time.Now().Unix(),
	}
	resBytes, err := json.Marshal(res)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	sig, err := crypto.Sign(string(resBytes), nodeLastKeys.PrivateKey())
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	res.SignatureResponse = sig

	return res, nil
}

func (s txSrv) StoreTransaction(ctx context.Context, req *api.StoreTransactionRequest) (*api.StoreTransactionResponse, error) {
	fmt.Printf("STORE TRANSACTION REQUEST - %s\n", time.Unix(req.Timestamp, 0).String())

	reqBytes, err := json.Marshal(&api.StoreTransactionRequest{
		MinedTransaction: req.MinedTransaction,
		Timestamp:        req.Timestamp,
	})
	nodeLastKeys, err := s.techDB.NodeLastKeys()
	if err != nil {
		return nil, err
	}
	if err := crypto.VerifySignature(string(reqBytes), nodeLastKeys.PublicKey(), req.SignatureRequest); err != nil {
		return nil, status.New(codes.InvalidArgument, err.Error()).Err()
	}

	tx, err := formatMinedTransaction(req.MinedTransaction.Transaction, req.MinedTransaction.MasterValidation, req.MinedTransaction.ConfirmValidations)
	if err != nil {
		return nil, status.New(codes.InvalidArgument, err.Error()).Err()
	}

	if !consensus.IsAuthorizedToStoreTx(tx) {
		return nil, status.New(codes.PermissionDenied, "not authorized to store this data").Err()
	}

	if err := chain.WriteTransaction(s.chainDB, s.locker, tx, int(req.MinimumValidations)); err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	res := &api.StoreTransactionResponse{
		Timestamp: time.Now().Unix(),
	}
	resBytes, err := json.Marshal(res)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	sig, err := crypto.Sign(string(resBytes), nodeLastKeys.PrivateKey())
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	res.SignatureResponse = sig
	return res, nil
}

func (s txSrv) LockTransaction(ctx context.Context, req *api.LockTransactionRequest) (*api.LockTransactionResponse, error) {
	fmt.Printf("LOCK TRANSACTION REQUEST - %s\n", time.Unix(req.Timestamp, 0).String())

	reqBytes, err := json.Marshal(&api.LockTransactionRequest{
		TransactionHash:     req.TransactionHash,
		MasterNodePublicKey: req.MasterNodePublicKey,
		Timestamp:           req.Timestamp,
		Address:             req.Address,
	})
	nodeLastKeys, err := s.techDB.NodeLastKeys()
	if err != nil {
		return nil, err
	}
	if err := crypto.VerifySignature(string(reqBytes), nodeLastKeys.PublicKey(), req.SignatureRequest); err != nil {
		return nil, status.New(codes.InvalidArgument, err.Error()).Err()
	}

	if err := chain.LockTransaction(s.locker, req.TransactionHash, req.Address, req.MasterNodePublicKey); err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	res := &api.LockTransactionResponse{
		Timestamp: time.Now().Unix(),
	}
	resBytes, err := json.Marshal(res)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	sig, err := crypto.Sign(string(resBytes), nodeLastKeys.PrivateKey())
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	res.SignatureResponse = sig
	return res, nil
}

func (s txSrv) LeadTransactionMining(ctx context.Context, req *api.LeadTransactionMiningRequest) (*api.LeadTransactionMiningResponse, error) {
	fmt.Printf("LEAD TRANSACTION MINING REQUEST - %s\n", time.Unix(req.Timestamp, 0).String())

	reqBytes, err := json.Marshal(&api.LeadTransactionMiningRequest{
		Transaction:        req.Transaction,
		MinimumValidations: req.MinimumValidations,
		Timestamp:          req.Timestamp,
	})
	nodeLastKeys, err := s.techDB.NodeLastKeys()
	if err != nil {
		return nil, err
	}
	if err := crypto.VerifySignature(string(reqBytes), nodeLastKeys.PublicKey(), req.SignatureRequest); err != nil {
		return nil, status.New(codes.InvalidArgument, err.Error()).Err()
	}

	tx, err := formatTransaction(req.Transaction)
	if err != nil {
		return nil, status.New(codes.InvalidArgument, err.Error()).Err()
	}

	if err := consensus.LeadMining(tx, int(req.MinimumValidations), s.poolR, s.nodePublicKey, s.nodePrivateKey, s.techDB); err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	res := &api.LeadTransactionMiningResponse{
		Timestamp: time.Now().Unix(),
	}
	resBytes, err := json.Marshal(res)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	sig, err := crypto.Sign(string(resBytes), nodeLastKeys.PrivateKey())
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	res.SignatureResponse = sig
	return res, nil
}

func (s txSrv) ConfirmTransactionValidation(ctx context.Context, req *api.ConfirmTransactionValidationRequest) (*api.ConfirmTransactionValidationResponse, error) {
	fmt.Printf("CONFIRM VALIDATION TRANSACTION REQUEST - %s\n", time.Unix(req.Timestamp, 0).String())

	reqBytes, err := json.Marshal(&api.ConfirmTransactionValidationRequest{
		Transaction:      req.Transaction,
		MasterValidation: req.MasterValidation,
		Timestamp:        req.Timestamp,
	})
	if err != nil {
		return nil, status.New(codes.InvalidArgument, err.Error()).Err()
	}

	nodeLastKeys, err := s.techDB.NodeLastKeys()
	if err != nil {
		return nil, err
	}
	if err := crypto.VerifySignature(string(reqBytes), nodeLastKeys.PublicKey(), req.SignatureRequest); err != nil {
		return nil, status.New(codes.InvalidArgument, err.Error()).Err()
	}

	tx, err := formatTransaction(req.Transaction)
	if err != nil {
		return nil, status.New(codes.InvalidArgument, err.Error()).Err()
	}
	masterValid, err := formatMasterValidation(req.MasterValidation)
	if err != nil {
		return nil, status.New(codes.InvalidArgument, err.Error()).Err()
	}
	valid, err := consensus.ConfirmTransactionValidation(tx, masterValid, s.nodePublicKey, s.nodePrivateKey)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	res := &api.ConfirmTransactionValidationResponse{
		Validation: formatAPIValidation(valid),
		Timestamp:  time.Now().Unix(),
	}
	resBytes, err := json.Marshal(res)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	sig, err := crypto.Sign(string(resBytes), nodeLastKeys.PrivateKey())
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	res.SignatureResponse = sig
	return res, nil
}

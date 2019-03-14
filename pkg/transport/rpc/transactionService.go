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
	chainDB         chain.Database
	sharedKeyReader shared.KeyReader
	poolR           consensus.PoolRequester
	nodePublicKey   crypto.PublicKey
	nodePrivateKey  crypto.PrivateKey
}

//NewTransactionService creates service handler for the GRPC Transaction service
func NewTransactionService(cDB chain.Database, skr shared.KeyReader, pR consensus.PoolRequester, nodePublicKeyk crypto.PublicKey, nodePrivateKeyk crypto.PrivateKey) api.TransactionServiceServer {
	return txSrv{
		chainDB:         cDB,
		sharedKeyReader: skr,
		poolR:           pR,
		nodePublicKey:   nodePublicKeyk,
		nodePrivateKey:  nodePrivateKeyk,
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

	nodeLastKeys, err := s.sharedKeyReader.LastNodeCrossKeypair()
	if err != nil {
		return nil, err
	}

	if !nodeLastKeys.PublicKey().Verify(reqBytes, req.SignatureRequest) {
		return nil, status.New(codes.InvalidArgument, "invalid signature").Err()
	}

	tx, err := chain.LastTransaction(s.chainDB, req.TransactionAddress, chain.TransactionType(req.Type))
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	if tx == nil {
		return nil, status.New(codes.NotFound, "transaction does not exist").Err()
	}

	tvf, err := formatAPITransaction(*tx)
	if err != nil {
		return nil, err
	}
	res := &api.GetLastTransactionResponse{
		Timestamp:   time.Now().Unix(),
		Transaction: tvf,
	}
	resBytes, err := json.Marshal(res)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	sig, err := nodeLastKeys.PrivateKey().Sign(resBytes)
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
	nodeLastKeys, err := s.sharedKeyReader.LastNodeCrossKeypair()
	if err != nil {
		return nil, err
	}
	if !nodeLastKeys.PublicKey().Verify(reqBytes, req.SignatureRequest) {
		return nil, status.New(codes.InvalidArgument, "invalid signature").Err()
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
	sig, err := nodeLastKeys.PrivateKey().Sign(resBytes)
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
	nodeLastKeys, err := s.sharedKeyReader.LastNodeCrossKeypair()
	if err != nil {
		return nil, err
	}
	if !nodeLastKeys.PublicKey().Verify(reqBytes, req.SignatureRequest) {
		return nil, status.New(codes.InvalidArgument, "invalid signature").Err()
	}

	tx, err := formatMinedTransaction(req.MinedTransaction.Transaction, req.MinedTransaction.MasterValidation, req.MinedTransaction.ConfirmValidations)
	if err != nil {
		return nil, status.New(codes.InvalidArgument, err.Error()).Err()
	}

	if !consensus.IsAuthorizedToStoreTx(tx) {
		return nil, status.New(codes.PermissionDenied, "not authorized to store this data").Err()
	}

	if err := chain.WriteTransaction(s.chainDB, tx, int(req.MinimumValidations)); err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	res := &api.StoreTransactionResponse{
		Timestamp: time.Now().Unix(),
	}
	resBytes, err := json.Marshal(res)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	sig, err := nodeLastKeys.PrivateKey().Sign(resBytes)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	res.SignatureResponse = sig
	return res, nil
}

func (s txSrv) TimeLockTransaction(ctx context.Context, req *api.TimeLockTransactionRequest) (*api.TimeLockTransactionResponse, error) {
	fmt.Printf("TIMELOCK TRANSACTION REQUEST - %s\n", time.Unix(req.Timestamp, 0).String())

	reqBytes, err := json.Marshal(&api.TimeLockTransactionRequest{
		TransactionHash:     req.TransactionHash,
		Address:             req.Address,
		MasterNodePublicKey: req.MasterNodePublicKey,
		Timestamp:           req.Timestamp,
	})
	nodeLastKeys, err := s.sharedKeyReader.LastNodeCrossKeypair()
	if err != nil {
		return nil, err
	}
	if !nodeLastKeys.PublicKey().Verify(reqBytes, req.SignatureRequest) {
		return nil, status.New(codes.InvalidArgument, "invalid signature").Err()
	}

	masterKey, err := crypto.ParsePublicKey(req.MasterNodePublicKey)
	if err != nil {
		return nil, status.New(codes.InvalidArgument, "invalid master public key").Err()
	}

	if err := chain.TimeLockTransaction(req.TransactionHash, req.Address, masterKey); err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	res := &api.TimeLockTransactionResponse{
		Timestamp: time.Now().Unix(),
	}
	resBytes, err := json.Marshal(res)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	sig, err := nodeLastKeys.PrivateKey().Sign(resBytes)
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
		WelcomeHeaders:     req.WelcomeHeaders,
	})
	nodeLastKeys, err := s.sharedKeyReader.LastNodeCrossKeypair()
	if err != nil {
		return nil, err
	}
	if !nodeLastKeys.PublicKey().Verify(reqBytes, req.SignatureRequest) {
		return nil, status.New(codes.InvalidArgument, "invalid signature").Err()
	}

	tx, err := formatTransaction(req.Transaction)
	if err != nil {
		return nil, status.New(codes.InvalidArgument, err.Error()).Err()
	}

	wHeaders, err := formatNodeHeaders(req.WelcomeHeaders)
	if err != nil {
		return nil, status.New(codes.InvalidArgument, err.Error()).Err()
	}

	if err := consensus.LeadMining(tx, int(req.MinimumValidations), wHeaders, s.poolR, s.nodePublicKey, s.nodePrivateKey, s.sharedKeyReader); err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	res := &api.LeadTransactionMiningResponse{
		Timestamp: time.Now().Unix(),
	}
	resBytes, err := json.Marshal(res)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	sig, err := nodeLastKeys.PrivateKey().Sign(resBytes)
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

	nodeLastKeys, err := s.sharedKeyReader.LastNodeCrossKeypair()
	if err != nil {
		return nil, err
	}
	if !nodeLastKeys.PublicKey().Verify(reqBytes, req.SignatureRequest) {
		return nil, status.New(codes.InvalidArgument, "invalid signature").Err()
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

	v, err := formatAPIValidation(valid)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	res := &api.ConfirmTransactionValidationResponse{
		Validation: v,
		Timestamp:  time.Now().Unix(),
	}
	resBytes, err := json.Marshal(res)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	sig, err := nodeLastKeys.PrivateKey().Sign(resBytes)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	res.SignatureResponse = sig
	return res, nil
}

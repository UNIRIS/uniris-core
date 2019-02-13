package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/uniris/uniris-core/pkg/chain"
	"github.com/uniris/uniris-core/pkg/consensus"
	"github.com/uniris/uniris-core/pkg/shared"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	api "github.com/uniris/uniris-core/api/protobuf-spec"
	"github.com/uniris/uniris-core/pkg/crypto"
)

type chainSrv struct {
	db     chain.Database
	techDB shared.TechDatabaseReader
	poolR  consensus.PoolRequester
}

//NewChainServer creates a new GPRC server for the chain service
func NewChainServer(db chain.Database, tDB shared.TechDatabaseReader, poolR consensus.PoolRequester) api.ChainServiceServer {
	return &chainSrv{
		db:     db,
		techDB: tDB,
		poolR:  poolR,
	}
}

func (s chainSrv) GetLastTransaction(ctx context.Context, req *api.LastTransactionRequest) (*api.LastTransactionResponse, error) {
	fmt.Printf("GET LAST TRANSACTION REQUEST - %s\n", time.Unix(req.Timestamp, 0).String())

	reqBytes, err := json.Marshal(&api.LastTransactionRequest{
		TransactionAddress: req.TransactionAddress,
		Type:               req.Type,
		Timestamp:          req.Timestamp,
	})
	if err != nil {
		return nil, status.New(codes.InvalidArgument, err.Error()).Err()
	}

	lastMinerKeys, err := s.techDB.LastMinerKeys()
	if err != nil {
		return nil, err
	}

	if err := crypto.VerifySignature(string(reqBytes), lastMinerKeys.PublicKey(), req.SignatureRequest); err != nil {
		return nil, status.New(codes.InvalidArgument, err.Error()).Err()
	}

	tx, err := chain.LastTransaction(s.db, req.TransactionAddress, chain.TransactionType(req.Type))
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	if tx == nil {
		return nil, status.New(codes.NotFound, "transaction does not exist").Err()
	}

	res := &api.LastTransactionResponse{
		Timestamp:   time.Now().Unix(),
		Transaction: formatAPITransaction(*tx),
	}
	resBytes, err := json.Marshal(res)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	sig, err := crypto.Sign(string(resBytes), lastMinerKeys.PrivateKey())
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	res.SignatureResponse = sig

	return res, nil
}

func (s chainSrv) GetTransactionStatus(ctx context.Context, req *api.TransactionStatusRequest) (*api.TransactionStatusResponse, error) {
	fmt.Printf("GET TRANSACTION STATUS REQUEST - %s\n", time.Unix(req.Timestamp, 0).String())

	reqBytes, err := json.Marshal(&api.TransactionStatusRequest{
		TransactionHash: req.TransactionHash,
		Timestamp:       req.Timestamp,
	})
	if err != nil {
		return nil, status.New(codes.InvalidArgument, err.Error()).Err()
	}
	lastMinerKeys, err := s.techDB.LastMinerKeys()
	if err != nil {
		return nil, err
	}
	if err := crypto.VerifySignature(string(reqBytes), lastMinerKeys.PublicKey(), req.SignatureRequest); err != nil {
		return nil, status.New(codes.InvalidArgument, err.Error()).Err()
	}

	txStatus, err := chain.GetTransactionStatus(s.db, req.TransactionHash)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	res := &api.TransactionStatusResponse{
		Status:    api.TransactionStatusResponse_TransactionStatus(txStatus),
		Timestamp: time.Now().Unix(),
	}
	resBytes, err := json.Marshal(res)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	sig, err := crypto.Sign(string(resBytes), lastMinerKeys.PrivateKey())
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	res.SignatureResponse = sig

	return res, nil
}

func (s chainSrv) StoreTransaction(ctx context.Context, req *api.StoreRequest) (*api.StoreResponse, error) {
	fmt.Printf("STORE TRANSACTION REQUEST - %s\n", time.Unix(req.Timestamp, 0).String())

	reqBytes, err := json.Marshal(&api.StoreRequest{
		MinedTransaction: req.MinedTransaction,
		Timestamp:        req.Timestamp,
	})
	lastMinerKeys, err := s.techDB.LastMinerKeys()
	if err != nil {
		return nil, err
	}
	if err := crypto.VerifySignature(string(reqBytes), lastMinerKeys.PublicKey(), req.SignatureRequest); err != nil {
		return nil, status.New(codes.InvalidArgument, err.Error()).Err()
	}

	tx, err := formatMinedTransaction(req.MinedTransaction.Transaction, req.MinedTransaction.MasterValidation, req.MinedTransaction.ConfirmValidations)
	if err != nil {
		return nil, status.New(codes.InvalidArgument, err.Error()).Err()
	}

	if !consensus.IsAuthorizedToStoreTx(tx) {
		return nil, status.New(codes.PermissionDenied, "not authorized to store this data").Err()
	}

	if err := chain.WriteTransaction(s.db, tx, int(req.MinimumValidations)); err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	res := &api.StoreResponse{
		Timestamp: time.Now().Unix(),
	}
	resBytes, err := json.Marshal(res)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	sig, err := crypto.Sign(string(resBytes), lastMinerKeys.PrivateKey())
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}
	res.SignatureResponse = sig
	return res, nil
}

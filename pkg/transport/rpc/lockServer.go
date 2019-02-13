package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/uniris/uniris-core/pkg/shared"

	api "github.com/uniris/uniris-core/api/protobuf-spec"
	"github.com/uniris/uniris-core/pkg/consensus"
	"github.com/uniris/uniris-core/pkg/crypto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type lockSrv struct {
	db     consensus.LockDatabase
	techDB shared.TechDatabaseReader
}

//NewLockServer creates anew GRPC service for the lock service
func NewLockServer(db consensus.LockDatabase, tDB shared.TechDatabaseReader) api.LockServiceServer {
	return &lockSrv{
		db:     db,
		techDB: tDB,
	}
}

func (s lockSrv) LockTransaction(ctx context.Context, req *api.LockRequest) (*api.LockResponse, error) {
	fmt.Printf("LOCK TRANSACTION REQUEST - %s\n", time.Unix(req.Timestamp, 0).String())

	reqBytes, err := json.Marshal(&api.LockRequest{
		TransactionHash:     req.TransactionHash,
		MasterPeerPublicKey: req.MasterPeerPublicKey,
		Timestamp:           req.Timestamp,
		Address:             req.Address,
	})
	lastMinerKeys, err := s.techDB.LastMinerKeys()
	if err != nil {
		return nil, err
	}
	if err := crypto.VerifySignature(string(reqBytes), lastMinerKeys.PublicKey(), req.SignatureRequest); err != nil {
		return nil, status.New(codes.InvalidArgument, err.Error()).Err()
	}

	if err := consensus.LockTransaction(s.db, req.TransactionHash, req.Address, req.MasterPeerPublicKey); err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	res := &api.LockResponse{
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

func (s lockSrv) UnlockTransaction(ctx context.Context, req *api.LockRequest) (*api.LockResponse, error) {
	fmt.Printf("UNLOCK TRANSACTION REQUEST - %s\n", time.Unix(req.Timestamp, 0).String())

	reqBytes, err := json.Marshal(&api.LockRequest{
		TransactionHash:     req.TransactionHash,
		MasterPeerPublicKey: req.MasterPeerPublicKey,
		Timestamp:           req.Timestamp,
		Address:             req.Address,
	})
	lastMinerKeys, err := s.techDB.LastMinerKeys()
	if err != nil {
		return nil, err
	}
	if err := crypto.VerifySignature(string(reqBytes), lastMinerKeys.PublicKey(), req.SignatureRequest); err != nil {
		return nil, status.New(codes.InvalidArgument, err.Error()).Err()
	}

	if err := consensus.UnlockTransaction(s.db, req.TransactionHash, req.Address); err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	res := &api.LockResponse{
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

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

type miningSrv struct {
	techDB   shared.TechDatabaseReader
	poolR    consensus.PoolRequester
	minerPub string
	minerPv  string
}

//NewMiningServer creates a new GRPC server for the mining service
func NewMiningServer(tDB shared.TechDatabaseReader, pr consensus.PoolRequester, pubKey, pvKey string) api.MiningServiceServer {
	return &miningSrv{
		techDB:   tDB,
		poolR:    pr,
		minerPub: pubKey,
		minerPv:  pvKey,
	}
}

func (s miningSrv) LeadTransactionMining(ctx context.Context, req *api.LeadMiningRequest) (*api.LeadMiningResponse, error) {
	fmt.Printf("LEAD TRANSACTION MINING REQUEST - %s\n", time.Unix(req.Timestamp, 0).String())

	reqBytes, err := json.Marshal(&api.LeadMiningRequest{
		Transaction:        req.Transaction,
		MinimumValidations: req.MinimumValidations,
		Timestamp:          req.Timestamp,
	})
	lastMinerKeys, err := s.techDB.LastMinerKeys()
	if err != nil {
		return nil, err
	}
	if err := crypto.VerifySignature(string(reqBytes), lastMinerKeys.PublicKey(), req.SignatureRequest); err != nil {
		return nil, status.New(codes.InvalidArgument, err.Error()).Err()
	}

	tx, err := formatTransaction(req.Transaction)
	if err != nil {
		return nil, status.New(codes.InvalidArgument, err.Error()).Err()
	}

	if err := consensus.LeadMining(tx, int(req.MinimumValidations), s.poolR, s.minerPub, s.minerPv, s.techDB); err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	res := &api.LeadMiningResponse{
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

func (s miningSrv) ConfirmTransactionValidation(ctx context.Context, req *api.ConfirmValidationRequest) (*api.ConfirmValidationResponse, error) {
	fmt.Printf("CONFIRM VALIDATION TRANSACTION REQUEST - %s\n", time.Unix(req.Timestamp, 0).String())

	reqBytes, err := json.Marshal(&api.ConfirmValidationRequest{
		Transaction:      req.Transaction,
		MasterValidation: req.MasterValidation,
		Timestamp:        req.Timestamp,
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

	tx, err := formatTransaction(req.Transaction)
	if err != nil {
		return nil, status.New(codes.InvalidArgument, err.Error()).Err()
	}
	masterValid, err := formatMasterValidation(req.MasterValidation)
	if err != nil {
		return nil, status.New(codes.InvalidArgument, err.Error()).Err()
	}
	valid, err := consensus.ConfirmTransactionValidation(tx, masterValid, s.minerPub, s.minerPv)
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	res := &api.ConfirmValidationResponse{
		Validation: formatAPIValidation(valid),
		Timestamp:  time.Now().Unix(),
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

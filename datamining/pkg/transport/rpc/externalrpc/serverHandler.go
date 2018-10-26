package externalrpc

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/uniris/uniris-core/datamining/pkg"

	"github.com/uniris/uniris-core/datamining/pkg/mining"
	"github.com/uniris/uniris-core/datamining/pkg/mining/lock"
	"github.com/uniris/uniris-core/datamining/pkg/mining/transactions"

	"github.com/golang/protobuf/ptypes/empty"
	api "github.com/uniris/uniris-core/datamining/api/protobuf-spec"
	"github.com/uniris/uniris-core/datamining/pkg/adding"
	"github.com/uniris/uniris-core/datamining/pkg/listing"
	"github.com/uniris/uniris-core/datamining/pkg/system"
)

type externalSrvHandler struct {
	list              listing.Service
	add               adding.Service
	mining            mining.Service
	sharedRobotPubKey string
	errors            system.DataMininingErrors
}

//NewExternalServerHandler creates a new External GRPC handler
func NewExternalServerHandler(list listing.Service, add adding.Service, mine mining.Service, sharedRobotPubKey string, errors system.DataMininingErrors) api.ExternalServer {
	return externalSrvHandler{list, add, mine, sharedRobotPubKey, errors}
}

func (s externalSrvHandler) LockTransaction(ctx context.Context, req *api.LockRequest) (*empty.Empty, error) {
	//TODO: verify signature

	if err := s.mining.LockTransaction(lock.TransactionLock{
		TxHash:         req.TransactionHash,
		MasterRobotKey: req.MasterRobotKey,
	}); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s externalSrvHandler) UnlockTransaction(ctx context.Context, req *api.LockRequest) (*empty.Empty, error) {
	//TODO: verify signature

	if err := s.mining.UnlockTransaction(lock.TransactionLock{
		TxHash:         req.TransactionHash,
		MasterRobotKey: req.MasterRobotKey,
	}); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s externalSrvHandler) Validate(ctx context.Context, req *api.ValidationRequest) (*api.ValidationResponse, error) {
	var data interface{}
	err := json.Unmarshal(req.Data.Value, &data)
	if err != nil {
		return nil, err
	}

	valid, err := s.mining.Validate(data, transactions.Type(req.TransactionType))
	if err != nil {
		return nil, err
	}

	return &api.ValidationResponse{
		Validation: BuildAPIValidation(valid),
	}, nil
}

func (s externalSrvHandler) Store(ctx context.Context, req *api.StorageRequest) (*empty.Empty, error) {

	//TODO: verify signatures

	switch req.TransactionType {
	case api.TransactionType_CreateWallet:
		w := &datamining.Wallet{}
		if err := w.UnmarshalJSON(req.Data.Value); err != nil {
			return nil, err
		}
		if err := s.add.StoreWallet(w); err != nil {
			return nil, err
		}
		return &empty.Empty{}, nil
	case api.TransactionType_CreateBio:
		bw := &datamining.BioWallet{}
		if err := bw.UnmarshalJSON(req.Data.Value); err != nil {
			return nil, err
		}
		if err := s.add.StoreBioWallet(bw); err != nil {
			return nil, err
		}
		return &empty.Empty{}, nil
	}

	return nil, errors.New("Unsupported operation")
}

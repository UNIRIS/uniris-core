package externalrpc

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/uniris/uniris-core/datamining/pkg/mining/slave"

	"github.com/uniris/uniris-core/datamining/pkg"

	"github.com/golang/protobuf/ptypes/empty"
	api "github.com/uniris/uniris-core/datamining/api/protobuf-spec"
	"github.com/uniris/uniris-core/datamining/pkg/adding"
	"github.com/uniris/uniris-core/datamining/pkg/listing"
	"github.com/uniris/uniris-core/datamining/pkg/system"
)

type externalSrvHandler struct {
	list              listing.Service
	add               adding.Service
	mining            slave.Service
	sharedRobotPubKey string
	errors            system.DataMininingErrors
}

//NewExternalServerHandler creates a new External GRPC handler
func NewExternalServerHandler(list listing.Service, add adding.Service, mine slave.Service, sharedRobotPubKey string, errors system.DataMininingErrors) api.ExternalServer {
	return externalSrvHandler{list, add, mine, sharedRobotPubKey, errors}
}

func (s externalSrvHandler) LockTransaction(ctx context.Context, req *api.LockRequest) (*empty.Empty, error) {
	//TODO: verify signature

	if err := s.mining.LockTransaction(slave.TransactionLock{
		TxHash:         req.TransactionHash,
		MasterRobotKey: req.MasterRobotKey,
	}); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s externalSrvHandler) UnlockTransaction(ctx context.Context, req *api.LockRequest) (*empty.Empty, error) {
	//TODO: verify signature

	if err := s.mining.UnlockTransaction(slave.TransactionLock{
		TxHash:         req.TransactionHash,
		MasterRobotKey: req.MasterRobotKey,
	}); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s externalSrvHandler) Validate(ctx context.Context, req *api.ValidationRequest) (*api.ValidationResponse, error) {

	//TODO: verify signatures

	var data interface{}
	err := json.Unmarshal(req.Data.Value, &data)
	if err != nil {
		return nil, err
	}

	valid, err := s.mining.Validate(data, datamining.TransactionType(req.TransactionType))
	if err != nil {
		return nil, err
	}

	return &api.ValidationResponse{
		Validation: &api.Validation{
			PublicKey: valid.PublicKey(),
			Signature: valid.Signature(),
			Status:    api.Validation_ValidationStatus(valid.Status()),
			Timestamp: valid.Timestamp().Unix(),
		},
	}, nil
}

func (s externalSrvHandler) Store(ctx context.Context, req *api.StorageRequest) (*empty.Empty, error) {

	//TODO: verify signatures

	switch req.TransactionType {
	case api.TransactionType_CreateKeychain:
		w := &datamining.Keychain{}
		if err := w.UnmarshalJSON(req.Data.Value); err != nil {
			return nil, err
		}
		if err := s.add.StoreKeychain(w); err != nil {
			return nil, err
		}
		return &empty.Empty{}, nil
	case api.TransactionType_CreateBio:
		bw := &datamining.Biometric{}
		if err := bw.UnmarshalJSON(req.Data.Value); err != nil {
			return nil, err
		}
		if err := s.add.StoreBiometric(bw); err != nil {
			return nil, err
		}
		return &empty.Empty{}, nil
	}

	return nil, errors.New("Unsupported operation")
}

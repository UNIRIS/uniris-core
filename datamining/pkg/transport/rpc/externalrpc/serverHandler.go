package externalrpc

import (
	"context"

	"github.com/uniris/uniris-core/datamining/pkg/validating"

	"github.com/golang/protobuf/ptypes/empty"
	api "github.com/uniris/uniris-core/datamining/api/protobuf-spec"
	"github.com/uniris/uniris-core/datamining/pkg/adding"
	"github.com/uniris/uniris-core/datamining/pkg/listing"
	"github.com/uniris/uniris-core/datamining/pkg/system"
)

type externalSrvHandler struct {
	list              listing.Service
	add               adding.Service
	valid             validating.Service
	sharedRobotPubKey string
	errors            system.DataMininingErrors
}

//NewExternalServerHandler creates a new External GRPC handler
func NewExternalServerHandler(list listing.Service, add adding.Service, valid validating.Service, sharedRobotPubKey string, errors system.DataMininingErrors) api.ExternalServer {
	return externalSrvHandler{list, add, valid, sharedRobotPubKey, errors}
}

func (s externalSrvHandler) GetLastWallet(ctx context.Context, req *api.LastWalletRequest) (*api.LastWalletResponse, error) {
	wallet, err := s.list.GetWallet(req.Address)
	if err != nil {
		return nil, err
	}
	return &api.LastWalletResponse{
		Wallet: BuildAPIWallet(wallet),
	}, nil
}

func (s externalSrvHandler) LockTransaction(ctx context.Context, req *api.LockRequest) (*empty.Empty, error) {
	//TODO: verify signature

	if err := s.valid.LockTransaction(validating.TransactionLock{
		TxHash:         req.TransactionHash,
		MasterRobotKey: req.MasterRobotKey,
	}); err != nil {
		return nil, err
	}

	return nil, nil
}

func (s externalSrvHandler) UnlockTransaction(ctx context.Context, req *api.LockRequest) (*empty.Empty, error) {
	//TODO: verify signature

	if err := s.valid.UnlockTransaction(validating.TransactionLock{
		TxHash:         req.TransactionHash,
		MasterRobotKey: req.MasterRobotKey,
	}); err != nil {
		return nil, err
	}

	return nil, nil
}

func (s externalSrvHandler) ValidateBio(ctx context.Context, req *api.BioValidationRequest) (*api.ValidationResponse, error) {
	v, err := s.valid.ValidateBioData(BuildBioDataFromValidation(req))
	if err != nil {
		return nil, err
	}

	return &api.ValidationResponse{
		Validation: BuildAPIValidation(v),
	}, nil
}

func (s externalSrvHandler) ValidateWallet(ctx context.Context, req *api.WalletValidationRequest) (*api.ValidationResponse, error) {
	v, err := s.valid.ValidateWalletData(BuildWalletFromValidation(req))
	if err != nil {
		return nil, err
	}

	return &api.ValidationResponse{
		Validation: BuildAPIValidation(v),
	}, nil
}

func (s externalSrvHandler) StoreBio(ctx context.Context, req *api.BioStorageRequest) (*empty.Empty, error) {
	w := BuilBioDataFromStoreRequest(req)

	//TODO: verify signatures

	if err := s.add.StoreBioWallet(w); err != nil {
		return nil, err
	}

	//TODO: handle store pending/ko

	return nil, nil
}

func (s externalSrvHandler) StoreWallet(ctx context.Context, req *api.WalletStoreRequest) (*empty.Empty, error) {
	w := BuildWalletDataFromStorageRequest(req)

	//TODO: verify signatures

	if err := s.add.StoreDataWallet(w); err != nil {
		return nil, err
	}

	//TODO: handle store pending/ko
	return nil, nil
}

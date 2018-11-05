package externalrpc

import (
	"context"
	"encoding/json"
	"fmt"
	"net"

	"github.com/uniris/uniris-core/datamining/pkg/crypto"
	"github.com/uniris/uniris-core/datamining/pkg/mining"

	"github.com/uniris/uniris-core/datamining/pkg/lock"

	"github.com/uniris/uniris-core/datamining/pkg/system"

	"github.com/golang/protobuf/ptypes/empty"
	api "github.com/uniris/uniris-core/datamining/api/protobuf-spec"

	accAdding "github.com/uniris/uniris-core/datamining/pkg/account/adding"
)

type externalSrvHandler struct {
	lock              lock.Service
	mining            mining.Service
	accAdd            accAdding.Service
	sharedRobotPubKey string
	sharedRobotPvKey  string
	errors            system.DataMininingErrors
}

//NewExternalServerHandler creates a new External GRPC handler
func NewExternalServerHandler(lock lock.Service, mining mining.Service, accAdd accAdding.Service, sharedRobotPubKey, sharedRobotPvKey string, errors system.DataMininingErrors) api.ExternalServer {
	return externalSrvHandler{lock, mining, accAdd, sharedRobotPubKey, sharedRobotPvKey, errors}
}

func (s externalSrvHandler) LeadKeychainMining(ctx context.Context, req *api.KeychainLeadRequest) (*empty.Empty, error) {
	keychainRawData, err := crypto.Decrypt(s.sharedRobotPvKey, req.EncryptedKeychainData)
	if err != nil {
		return nil, fmt.Errorf("Cannot decrypt the wallet data - %s", err.Error())
	}

	var keychain *KeychainDataFromJSON
	err = json.Unmarshal([]byte(keychainRawData), &keychain)
	if err != nil {
		return nil, err
	}

	clearaddr, err := crypto.Decrypt(s.sharedRobotPvKey, keychain.EncryptedAddrRobot)
	if err != nil {
		return nil, fmt.Errorf("Cannot decrypt the address - %s", err.Error())
	}

	keychainData := BuildKeychainData(keychain, req.SignatureKeychainData, clearaddr)
	pp := make([]mining.Peer, 0)
	for _, p := range req.ValidatorPeerIPs {
		pp = append(pp, mining.Peer{IP: net.ParseIP(p)})
	}
	vPool := mining.NewPool(pp...)
	if err := s.mining.LeadMining(req.TransactionHash, clearaddr, keychainData, vPool, mining.CreateKeychainTransaction, keychainData.Sigs.BiodSig); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s externalSrvHandler) LeadBiometricMining(ctx context.Context, req *api.BiometricLeadRequest) (*empty.Empty, error) {
	biometricRawData, err := crypto.Decrypt(s.sharedRobotPvKey, req.EncryptedBioData)
	if err != nil {
		return nil, fmt.Errorf("Cannot decrypt the wallet data - %s", err.Error())
	}

	var bio *BioDataFromJSON
	err = json.Unmarshal([]byte(biometricRawData), &bio)
	if err != nil {
		return nil, err
	}

	clearaddr, err := crypto.Decrypt(s.sharedRobotPvKey, bio.EncryptedAddrRobot)
	if err != nil {
		return nil, fmt.Errorf("Cannot decrypt the address - %s", err.Error())
	}

	bioData := BuildBioData(bio, req.SignatureBioData)

	pp := make([]mining.Peer, 0)
	for _, p := range req.ValidatorPeerIPs {
		pp = append(pp, mining.Peer{IP: net.ParseIP(p)})
	}
	vPool := mining.NewPool(pp...)
	if err := s.mining.LeadMining(req.TransactionHash, clearaddr, bioData, vPool, mining.CreateBioTransaction, bioData.Sigs.BiodSig); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s externalSrvHandler) LockTransaction(ctx context.Context, req *api.LockRequest) (*empty.Empty, error) {
	//TODO: verify signature

	if err := s.lock.LockTransaction(lock.TransactionLock{
		TxHash:         req.TransactionHash,
		MasterRobotKey: req.MasterRobotKey,
	}); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s externalSrvHandler) UnlockTransaction(ctx context.Context, req *api.LockRequest) (*empty.Empty, error) {
	//TODO: verify signature

	if err := s.lock.UnlockTransaction(lock.TransactionLock{
		TxHash:         req.TransactionHash,
		MasterRobotKey: req.MasterRobotKey,
	}); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s externalSrvHandler) ValidateKeychain(ctx context.Context, req *api.KeychainValidationRequest) (*api.ValidationResponse, error) {
	//TODO: verify signatures

	clearaddr, err := crypto.Decrypt(s.sharedRobotPvKey, req.Data.CipherAddrRobot)
	if err != nil {
		return nil, fmt.Errorf("Cannot decrypt the address - %s", err.Error())
	}
	data := formatKeychainDataAPI(req.Data, clearaddr)

	valid, err := s.mining.Validate(req.TransactionHash, data, mining.CreateKeychainTransaction)
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

func (s externalSrvHandler) ValidateBiometric(ctx context.Context, req *api.BiometricValidationRequest) (*api.ValidationResponse, error) {
	//TODO: verify signatures
	data := formatBiometricDataAPI(req.Data)

	valid, err := s.mining.Validate(req.TransactionHash, data, mining.CreateBioTransaction)
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

func (s externalSrvHandler) StoreKeychain(ctx context.Context, req *api.KeychainStorageRequest) (*empty.Empty, error) {
	//TODO: verify signatures

	clearaddr, err := crypto.Decrypt(s.sharedRobotPvKey, req.Data.CipherAddrRobot)
	if err != nil {
		return nil, fmt.Errorf("Cannot decrypt the address - %s", err.Error())
	}

	if err := s.accAdd.StoreKeychain(formatKeychainDataAPI(req.Data, clearaddr), formatEndorsementAPI(req.Endorsement)); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s externalSrvHandler) StoreBiometric(ctx context.Context, req *api.BiometricStorageRequest) (*empty.Empty, error) {
	//TODO: verify signatures
	if err := s.accAdd.StoreBiometric(formatBiometricDataAPI(req.Data), formatEndorsementAPI(req.Endorsement)); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

package externalrpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"

	datamining "github.com/uniris/uniris-core/datamining/pkg"
	"github.com/uniris/uniris-core/datamining/pkg/mining"
	"github.com/uniris/uniris-core/datamining/pkg/transport/rpc"

	"github.com/uniris/uniris-core/datamining/pkg/lock"

	"github.com/uniris/uniris-core/datamining/pkg/system"

	"github.com/golang/protobuf/ptypes/empty"
	api "github.com/uniris/uniris-core/datamining/api/protobuf-spec"

	accAdding "github.com/uniris/uniris-core/datamining/pkg/account/adding"
	accListing "github.com/uniris/uniris-core/datamining/pkg/account/listing"
)

type externalSrvHandler struct {
	lock      lock.Service
	mining    mining.Service
	accAdd    accAdding.Service
	accLister accListing.Service
	decrypt   rpc.Decrypter
	signer    rpc.Signer
	conf      system.UnirisConfig
}

//NewExternalServerHandler creates a new External GRPC handler
func NewExternalServerHandler(lock lock.Service, mining mining.Service, accAdd accAdding.Service, accLister accListing.Service, decrypt rpc.Decrypter, signer rpc.Signer, conf system.UnirisConfig) api.ExternalServer {
	return externalSrvHandler{lock, mining, accAdd, accLister, decrypt, signer, conf}
}

func (s externalSrvHandler) GetBiometric(ctxt context.Context, req *api.BiometricRequest) (*api.BiometricResponse, error) {
	//TODO: check signature

	personHash, err := s.decrypt.DecryptHashPerson(req.PersonHash, s.conf.SharedKeys.RobotPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("Cannot decrypt person hash - %s", err.Error())
	}

	biometric, err := s.accLister.GetBiometric(personHash)
	if err != nil {
		return nil, err
	}

	if biometric == nil {
		return nil, errors.New(s.conf.Datamining.Errors.AccountNotExist)
	}

	sig, err := s.signer.SignBiometric(rpc.BuildBiometricJSON(biometric), s.conf.SharedKeys.RobotPrivateKey)
	if err != nil {
		return nil, err
	}

	return buildBiometricAPIResponse(biometric, sig), nil
}

func (s externalSrvHandler) GetKeychain(ctxt context.Context, req *api.KeychainRequest) (*api.KeychainResponse, error) {

	//TODO: check signature

	clearaddr, err := s.decrypt.DecryptCipherAddress(req.Address, s.conf.SharedKeys.RobotPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("Cannot decrypt the address - %s", err.Error())
	}

	keychain, err := s.accLister.GetLastKeychain(clearaddr)
	if err != nil {
		return nil, err
	}

	if keychain == nil {
		return nil, errors.New(s.conf.Datamining.Errors.AccountNotExist)
	}

	sig, err := s.signer.SignKeychain(rpc.BuildKeychainJSON(keychain), s.conf.SharedKeys.RobotPrivateKey)
	if err != nil {
		return nil, err
	}

	return buildKeychainAPIResponse(keychain, sig), nil
}

func (s externalSrvHandler) LeadKeychainMining(ctx context.Context, req *api.KeychainLeadRequest) (*empty.Empty, error) {
	keychainRawData, err := s.decrypt.DecryptTransactionData(req.EncryptedKeychainData, s.conf.SharedKeys.RobotPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("Cannot decrypt the wallet data - %s", err.Error())
	}

	var keychain *rpc.KeychainDataJSON
	err = json.Unmarshal([]byte(keychainRawData), &keychain)
	if err != nil {
		return nil, err
	}

	clearaddr, err := s.decrypt.DecryptCipherAddress(keychain.EncryptedAddrRobot, s.conf.SharedKeys.RobotPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("Cannot decrypt the address - %s", err.Error())
	}

	keychainData := rpc.BuildKeychainDataFromJSON(keychain, req.SignatureKeychainData, clearaddr)
	pp := make([]datamining.Peer, 0)
	for _, p := range req.ValidatorPeerIPs {
		pp = append(pp, datamining.Peer{IP: net.ParseIP(p)})
	}
	vPool := datamining.NewPool(pp...)
	if err := s.mining.LeadMining(req.TransactionHash, clearaddr, keychainData, vPool, mining.KeychainTransaction, keychainData.Sigs.BiodSig); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s externalSrvHandler) LeadBiometricMining(ctx context.Context, req *api.BiometricLeadRequest) (*empty.Empty, error) {
	biometricRawData, err := s.decrypt.DecryptTransactionData(req.EncryptedBioData, s.conf.SharedKeys.RobotPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("Cannot decrypt the wallet data - %s", err.Error())
	}

	var bio *rpc.BioDataJSON
	err = json.Unmarshal([]byte(biometricRawData), &bio)
	if err != nil {
		return nil, err
	}

	clearaddr, err := s.decrypt.DecryptCipherAddress(bio.EncryptedAddrRobot, s.conf.SharedKeys.RobotPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("Cannot decrypt the address - %s", err.Error())
	}

	bioData := rpc.BuildBioDataFromJSON(bio, req.SignatureBioData)

	pp := make([]datamining.Peer, 0)
	for _, p := range req.ValidatorPeerIPs {
		pp = append(pp, datamining.Peer{IP: net.ParseIP(p)})
	}
	vPool := datamining.NewPool(pp...)
	if err := s.mining.LeadMining(req.TransactionHash, clearaddr, bioData, vPool, mining.BiometricTransaction, bioData.Sigs.BiodSig); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s externalSrvHandler) LockTransaction(ctx context.Context, req *api.LockRequest) (*empty.Empty, error) {
	//TODO: verify signature

	if err := s.lock.LockTransaction(lock.TransactionLock{
		TxHash:         req.TransactionHash,
		MasterRobotKey: req.MasterRobotKey,
		Address:        req.Address,
	}); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s externalSrvHandler) UnlockTransaction(ctx context.Context, req *api.LockRequest) (*empty.Empty, error) {
	//TODO: verify signature

	if err := s.lock.UnlockTransaction(lock.TransactionLock{
		Address:        req.Address,
		TxHash:         req.TransactionHash,
		MasterRobotKey: req.MasterRobotKey,
	}); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s externalSrvHandler) ValidateKeychain(ctx context.Context, req *api.KeychainValidationRequest) (*api.ValidationResponse, error) {
	//TODO: verify signatures

	clearaddr, err := s.decrypt.DecryptCipherAddress(req.Data.CipherAddrRobot, s.conf.SharedKeys.RobotPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("Cannot decrypt the address - %s", err.Error())
	}
	data := buildKeychainDataFromAPI(req.Data, clearaddr)

	valid, err := s.mining.Validate(req.TransactionHash, data, mining.KeychainTransaction)
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
	data := buildBiometricDataFromAPI(req.Data)

	valid, err := s.mining.Validate(req.TransactionHash, data, mining.BiometricTransaction)
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

	clearaddr, err := s.decrypt.DecryptCipherAddress(req.Data.CipherAddrRobot, s.conf.SharedKeys.RobotPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("Cannot decrypt the address - %s", err.Error())
	}

	if err := s.accAdd.StoreKeychain(buildKeychainDataFromAPI(req.Data, clearaddr), buildEndorsementFromAPI(req.Endorsement)); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s externalSrvHandler) StoreBiometric(ctx context.Context, req *api.BiometricStorageRequest) (*empty.Empty, error) {
	//TODO: verify signatures
	if err := s.accAdd.StoreBiometric(buildBiometricDataFromAPI(req.Data), buildEndorsementFromAPI(req.Endorsement)); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

package rpc

import (
	"context"
	"errors"
	"net"

	"github.com/golang/protobuf/ptypes/empty"
	datamining "github.com/uniris/uniris-core/datamining/pkg"
	"github.com/uniris/uniris-core/datamining/pkg/account"
	"github.com/uniris/uniris-core/datamining/pkg/contract"
	"github.com/uniris/uniris-core/datamining/pkg/mining"

	"github.com/uniris/uniris-core/datamining/pkg/lock"

	"github.com/uniris/uniris-core/datamining/pkg/system"

	api "github.com/uniris/uniris-core/datamining/api/protobuf-spec"

	accAdding "github.com/uniris/uniris-core/datamining/pkg/account/adding"
	accListing "github.com/uniris/uniris-core/datamining/pkg/account/listing"

	contractAdding "github.com/uniris/uniris-core/datamining/pkg/contract/adding"
	contractListing "github.com/uniris/uniris-core/datamining/pkg/contract/listing"
)

//Services define the required services
type Services struct {
	lock           lock.Service
	mining         mining.Service
	accAdd         accAdding.Service
	accLister      accListing.Service
	contractAdd    contractAdding.Service
	contractLister contractListing.Service
}

//NewExternalServices creates a new container of required services
func NewExternalServices(lock lock.Service, mine mining.Service, accountAdder accAdding.Service, accountLister accListing.Service, contractAdder contractAdding.Service, contractLister contractListing.Service) Services {
	return Services{
		lock:           lock,
		mining:         mine,
		accAdd:         accountAdder,
		accLister:      accountLister,
		contractAdd:    contractAdder,
		contractLister: contractLister,
	}
}

type externalSrvHandler struct {
	services Services
	crypto   Crypto
	conf     system.UnirisConfig
	api      apiBuilder
	data     dataBuilder
}

//NewExternalServerHandler creates a new External GRPC handler
func NewExternalServerHandler(srv Services, crypto Crypto, conf system.UnirisConfig) api.ExternalServer {
	return externalSrvHandler{
		services: srv,
		crypto:   crypto,
		conf:     conf,
		api:      apiBuilder{},
		data:     dataBuilder{},
	}
}

func (h externalSrvHandler) GetID(ctxt context.Context, req *api.IDRequest) (*api.IDResponse, error) {
	if err := h.crypto.signer.VerifyHashSignature(h.conf.SharedKeys.Robot.PublicKey, req.EncryptedIDHash, req.Signature); err != nil {
		return nil, ErrInvalidSignature
	}

	idHash, err := h.crypto.decrypter.DecryptHash(req.EncryptedIDHash, h.conf.SharedKeys.Robot.PrivateKey)
	if err != nil {
		return nil, ErrInvalidEncryption
	}

	id, err := h.services.accLister.GetID(idHash)
	if err != nil {
		return nil, err
	}

	if id == nil {
		return nil, errors.New(h.conf.Services.Datamining.Errors.AccountNotExist)
	}

	res := &api.IDResponse{
		Data:        h.api.buildID(id),
		Endorsement: h.api.buildEndorsement(id.Endorsement()),
	}

	if err := h.crypto.signer.SignIDResponse(res, h.conf.SharedKeys.Robot.PrivateKey); err != nil {
		return nil, err
	}

	return res, nil
}

func (h externalSrvHandler) GetKeychain(ctxt context.Context, req *api.KeychainRequest) (*api.KeychainResponse, error) {
	if err := h.crypto.signer.VerifyHashSignature(h.conf.SharedKeys.Robot.PublicKey, req.EncryptedAddress, req.Signature); err != nil {
		return nil, ErrInvalidSignature
	}

	clearAddress, err := h.crypto.decrypter.DecryptHash(req.EncryptedAddress, h.conf.SharedKeys.Robot.PrivateKey)
	if err != nil {
		return nil, ErrInvalidEncryption
	}

	keychain, err := h.services.accLister.GetLastKeychain(clearAddress)
	if err != nil {
		return nil, err
	}

	if keychain == nil {
		return nil, errors.New(h.conf.Services.Datamining.Errors.AccountNotExist)
	}

	res := &api.KeychainResponse{
		Data:        h.api.buildKeychain(keychain),
		Endorsement: h.api.buildEndorsement(keychain.Endorsement()),
	}

	if err := h.crypto.signer.SignKeychainResponse(res, h.conf.SharedKeys.Robot.PrivateKey); err != nil {
		return nil, err
	}

	return res, nil
}

func (h externalSrvHandler) LeadKeychainMining(ctx context.Context, req *api.KeychainLeadRequest) (*empty.Empty, error) {
	if err := h.crypto.signer.VerifyKeychainLeadRequestSignature(h.conf.SharedKeys.Robot.PublicKey, req); err != nil {
		return nil, ErrInvalidSignature
	}

	keychain, err := h.crypto.decrypter.DecryptKeychain(req.EncryptedKeychain, h.conf.SharedKeys.Robot.PrivateKey)
	if err != nil {
		return nil, ErrInvalidEncryption
	}

	clearaddr, err := h.crypto.decrypter.DecryptHash(keychain.EncryptedAddrByRobot(), h.conf.SharedKeys.Robot.PrivateKey)
	if err != nil {
		return nil, ErrInvalidEncryption
	}

	pp := make([]datamining.Peer, 0)
	for _, p := range req.ValidatorPeerIPs {
		pp = append(pp, datamining.Peer{IP: net.ParseIP(p)})
	}
	vPool := datamining.NewPool(pp...)
	if err := h.services.mining.LeadMining(req.TransactionHash, clearaddr, keychain, vPool, mining.KeychainTransaction, keychain.EmitterSignature()); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (h externalSrvHandler) LeadIDMining(ctx context.Context, req *api.IDLeadRequest) (*empty.Empty, error) {
	if err := h.crypto.signer.VerifyIDLeadRequestSignature(h.conf.SharedKeys.Robot.PublicKey, req); err != nil {
		return nil, ErrInvalidSignature
	}

	id, err := h.crypto.decrypter.DecryptID(req.EncryptedID, h.conf.SharedKeys.Robot.PrivateKey)
	if err != nil {
		return nil, ErrInvalidEncryption
	}

	clearaddr, err := h.crypto.decrypter.DecryptHash(id.EncryptedAddrByRobot(), h.conf.SharedKeys.Robot.PrivateKey)
	if err != nil {
		return nil, ErrInvalidEncryption
	}

	pp := make([]datamining.Peer, 0)
	for _, p := range req.ValidatorPeerIPs {
		pp = append(pp, datamining.Peer{IP: net.ParseIP(p)})
	}
	vPool := datamining.NewPool(pp...)
	if err := h.services.mining.LeadMining(req.TransactionHash, clearaddr, id, vPool, mining.IDTransaction, id.EmitterSignature()); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (h externalSrvHandler) LockTransaction(ctx context.Context, req *api.LockRequest) (*api.LockAck, error) {
	if err := h.crypto.signer.VerifyLockRequestSignature(h.conf.SharedKeys.Robot.PublicKey, req); err != nil {
		return nil, err
	}

	lockTx := lock.TransactionLock{
		TxHash:         req.TransactionHash,
		MasterRobotKey: req.MasterRobotKey,
		Address:        req.Address,
	}

	if err := h.services.lock.LockTransaction(lockTx); err != nil {
		return nil, err
	}

	lockHash, err := h.crypto.hasher.HashLock(lockTx)
	if err != nil {
		return nil, err
	}

	ack := &api.LockAck{
		LockHash: lockHash,
	}
	if err := h.crypto.signer.SignLockAck(ack, h.conf.SharedKeys.Robot.PrivateKey); err != nil {
		return nil, err
	}
	return ack, nil
}

func (h externalSrvHandler) UnlockTransaction(ctx context.Context, req *api.LockRequest) (*api.LockAck, error) {

	if err := h.crypto.signer.VerifyLockRequestSignature(h.conf.SharedKeys.Robot.PublicKey, req); err != nil {
		return nil, err
	}

	lockTx := lock.TransactionLock{
		TxHash:         req.TransactionHash,
		MasterRobotKey: req.MasterRobotKey,
		Address:        req.Address,
	}

	if err := h.services.lock.UnlockTransaction(lockTx); err != nil {
		return nil, err
	}

	lockHash, err := h.crypto.hasher.HashLock(lockTx)
	if err != nil {
		return nil, err
	}

	ack := &api.LockAck{
		LockHash: lockHash,
	}
	if err := h.crypto.signer.SignLockAck(ack, h.conf.SharedKeys.Robot.PrivateKey); err != nil {
		return nil, err
	}
	return ack, nil
}

func (h externalSrvHandler) ValidateKeychain(ctx context.Context, req *api.KeychainValidationRequest) (*api.ValidationResponse, error) {
	if err := h.crypto.signer.VerifyKeychainValidationRequestSignature(h.conf.SharedKeys.Robot.PublicKey, req); err != nil {
		return nil, err
	}

	valid, err := h.services.mining.Validate(req.TransactionHash, h.data.buildKeychain(req.Data), mining.KeychainTransaction)
	if err != nil {
		return nil, err
	}

	vRes := &api.Validation{
		PublicKey: valid.PublicKey(),
		Signature: valid.Signature(),
		Status:    api.Validation_ValidationStatus(valid.Status()),
		Timestamp: valid.Timestamp().Unix(),
	}

	res := &api.ValidationResponse{
		Validation: vRes,
	}
	if err := h.crypto.signer.SignValidationResponse(res, h.conf.SharedKeys.Robot.PrivateKey); err != nil {
		return nil, err
	}

	return res, nil
}

func (h externalSrvHandler) ValidateID(ctx context.Context, req *api.IDValidationRequest) (*api.ValidationResponse, error) {
	if err := h.crypto.signer.VerifyIDValidationRequestSignature(h.conf.SharedKeys.Robot.PublicKey, req); err != nil {
		return nil, err
	}

	valid, err := h.services.mining.Validate(req.TransactionHash, h.data.buildID(req.Data), mining.IDTransaction)
	if err != nil {
		return nil, err
	}

	vRes := &api.Validation{
		PublicKey: valid.PublicKey(),
		Signature: valid.Signature(),
		Status:    api.Validation_ValidationStatus(valid.Status()),
		Timestamp: valid.Timestamp().Unix(),
	}

	res := &api.ValidationResponse{
		Validation: vRes,
	}
	if err := h.crypto.signer.SignValidationResponse(res, h.conf.SharedKeys.Robot.PrivateKey); err != nil {
		return nil, err
	}

	return res, nil
}

func (h externalSrvHandler) StoreKeychain(ctx context.Context, req *api.KeychainStorageRequest) (*api.StorageAck, error) {
	if err := h.crypto.signer.VerifyKeychainStorageRequestSignature(h.conf.SharedKeys.Robot.PublicKey, req); err != nil {
		return nil, err
	}

	clearaddr, err := h.crypto.decrypter.DecryptHash(req.Data.EncryptedAddrByRobot, h.conf.SharedKeys.Robot.PrivateKey)
	if err != nil {
		return nil, ErrInvalidEncryption
	}

	keychain := account.NewEndorsedKeychain(clearaddr,
		h.data.buildKeychain(req.Data),
		h.data.buildEndorsement(req.Endorsement))

	if err := h.services.accAdd.StoreKeychain(keychain); err != nil {
		return nil, err
	}

	hash, err := h.crypto.hasher.HashEndorsedKeychain(keychain)
	if err != nil {
		return nil, err
	}

	ack := &api.StorageAck{
		StorageHash: hash,
	}
	if err := h.crypto.signer.SignStorageAck(ack, h.conf.SharedKeys.Robot.PrivateKey); err != nil {
		return nil, err
	}
	return ack, nil
}

func (h externalSrvHandler) StoreID(ctx context.Context, req *api.IDStorageRequest) (*api.StorageAck, error) {
	if err := h.crypto.signer.VerifyIDStorageRequestSignature(h.conf.SharedKeys.Robot.PublicKey, req); err != nil {
		return nil, err
	}

	id := account.NewEndorsedID(h.data.buildID(req.Data), h.data.buildEndorsement(req.Endorsement))
	if err := h.services.accAdd.StoreID(id); err != nil {
		return nil, err
	}

	hash, err := h.crypto.hasher.HashEndorsedID(id)
	if err != nil {
		return nil, err
	}

	ack := &api.StorageAck{
		StorageHash: hash,
	}
	if err := h.crypto.signer.SignStorageAck(ack, h.conf.SharedKeys.Robot.PrivateKey); err != nil {
		return nil, err
	}
	return ack, nil
}

func (h externalSrvHandler) GetTransactionStatus(ctx context.Context, req *api.TransactionStatusRequest) (*api.TransactionStatusResponse, error) {

	addr, err := h.crypto.decrypter.DecryptHash(req.Address, h.conf.SharedKeys.Robot.PrivateKey)
	if err != nil {
		return nil, err
	}

	//TODO: search in the pending database

	//Search if the transaction is a keychain and is stored
	kc, err := h.services.accLister.GetKeychain(addr, req.Hash)
	if err != nil {
		return nil, err
	}
	if kc != nil {
		return &api.TransactionStatusResponse{
			Status: api.TransactionStatusResponse_TransactionStatus(kc.Endorsement().GetStatus()),
		}, nil
	}

	//Search if the transaction is a ID and is stored
	id, err := h.services.accLister.GetIDByTransaction(req.Hash)
	if err != nil {
		return nil, err
	}
	if id != nil {
		return &api.TransactionStatusResponse{
			Status: api.TransactionStatusResponse_TransactionStatus(id.Endorsement().GetStatus()),
		}, nil
	}

	//TODO: //Search if the transaction is an IRIS exchange and is stored

	//TODO: //Search if the transaction is a smartcontract and is stored
	contract, err := h.services.contractLister.GetContractByAddressAndTransaction(addr, req.Hash)
	if err != nil {
		return nil, err
	}
	if contract != nil {
		return &api.TransactionStatusResponse{
			Status: api.TransactionStatusResponse_TransactionStatus(contract.Endorsement().GetStatus()),
		}, nil
	}

	contractMsg, err := h.services.contractLister.GetContractMessageByContractAndTransaction(addr, req.Hash)
	if err != nil {
		return nil, err
	}
	if contractMsg != nil {
		return &api.TransactionStatusResponse{
			Status: api.TransactionStatusResponse_TransactionStatus(contractMsg.Endorsement().GetStatus()),
		}, nil
	}

	return &api.TransactionStatusResponse{
		Status: api.TransactionStatusResponse_Unknown,
	}, nil
}

func (h externalSrvHandler) LeadContractMining(ctx context.Context, req *api.ContractLeadRequest) (*empty.Empty, error) {
	if err := h.crypto.signer.VerifyContractLeadRequestSignature(h.conf.SharedKeys.Robot.PublicKey, req); err != nil {
		return nil, ErrInvalidSignature
	}

	contract := contract.New(req.Contract.Address, req.Contract.Code, req.Contract.Event, req.Contract.PublicKey, req.Contract.Signature, req.Contract.EmitterSignature)

	pp := make([]datamining.Peer, 0)
	for _, p := range req.ValidatorPeerIPs {
		pp = append(pp, datamining.Peer{IP: net.ParseIP(p)})
	}
	vPool := datamining.NewPool(pp...)
	if err := h.services.mining.LeadMining(req.TransactionHash, req.Contract.Address, contract, vPool, mining.ContractTransaction, req.Contract.EmitterSignature); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (h externalSrvHandler) ValidateContract(ctx context.Context, req *api.ContractValidationRequest) (*api.ValidationResponse, error) {
	if err := h.crypto.signer.VerifyContractValidationRequestSignature(h.conf.SharedKeys.Robot.PublicKey, req); err != nil {
		return nil, err
	}

	contract := contract.New(req.Contract.Address, req.Contract.Code, req.Contract.Event, req.Contract.PublicKey, req.Contract.Signature, req.Contract.EmitterSignature)

	valid, err := h.services.mining.Validate(req.TransactionHash, contract, mining.ContractTransaction)
	if err != nil {
		return nil, err
	}

	vRes := &api.Validation{
		PublicKey: valid.PublicKey(),
		Signature: valid.Signature(),
		Status:    api.Validation_ValidationStatus(valid.Status()),
		Timestamp: valid.Timestamp().Unix(),
	}

	res := &api.ValidationResponse{
		Validation: vRes,
	}
	if err := h.crypto.signer.SignValidationResponse(res, h.conf.SharedKeys.Robot.PrivateKey); err != nil {
		return nil, err
	}

	return res, nil
}

func (h externalSrvHandler) StoreContract(ctx context.Context, req *api.ContractStorageRequest) (*api.StorageAck, error) {
	if err := h.crypto.signer.VerifyContractStorageRequestSignature(h.conf.SharedKeys.Robot.PublicKey, req); err != nil {
		return nil, err
	}

	contract := contract.NewEndorsedContract(
		contract.New(req.Contract.Address, req.Contract.Code, req.Contract.Event, req.Contract.PublicKey, req.Contract.Signature, req.Contract.EmitterSignature),
		h.data.buildEndorsement(req.Endorsement))

	if err := h.services.contractAdd.StoreEndorsedContract(contract); err != nil {
		return nil, err
	}

	hash, err := h.crypto.hasher.HashEndorsedContract(contract)
	if err != nil {
		return nil, err
	}

	ack := &api.StorageAck{
		StorageHash: hash,
	}
	if err := h.crypto.signer.SignStorageAck(ack, h.conf.SharedKeys.Robot.PrivateKey); err != nil {
		return nil, err
	}
	return ack, nil
}

func (h externalSrvHandler) LeadContractMessageMining(ctx context.Context, req *api.ContractMessageLeadRequest) (*empty.Empty, error) {
	if err := h.crypto.signer.VerifyContractMessageLeadRequestSignature(h.conf.SharedKeys.Robot.PublicKey, req); err != nil {
		return nil, ErrInvalidSignature
	}

	msg := contract.NewMessage(req.ContractMessage.ContractAddress, req.ContractMessage.Method, req.ContractMessage.Parameters, req.ContractMessage.PublicKey, req.ContractMessage.Signature, req.ContractMessage.EmitterSignature)

	pp := make([]datamining.Peer, 0)
	for _, p := range req.ValidatorPeerIPs {
		pp = append(pp, datamining.Peer{IP: net.ParseIP(p)})
	}
	vPool := datamining.NewPool(pp...)
	if err := h.services.mining.LeadMining(req.TransactionHash, req.ContractMessage.ContractAddress, msg, vPool, mining.ContractMessageTransaction, req.ContractMessage.EmitterSignature); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (h externalSrvHandler) ValidateContractMessage(ctx context.Context, req *api.ContractMessageValidationRequest) (*api.ValidationResponse, error) {
	if err := h.crypto.signer.VerifyContractMessageValidationRequestSignature(h.conf.SharedKeys.Robot.PublicKey, req); err != nil {
		return nil, err
	}

	msg := contract.NewMessage(req.ContractMessage.ContractAddress, req.ContractMessage.Method, req.ContractMessage.Parameters, req.ContractMessage.PublicKey, req.ContractMessage.Signature, req.ContractMessage.EmitterSignature)

	valid, err := h.services.mining.Validate(req.TransactionHash, msg, mining.ContractMessageTransaction)
	if err != nil {
		return nil, err
	}

	vRes := &api.Validation{
		PublicKey: valid.PublicKey(),
		Signature: valid.Signature(),
		Status:    api.Validation_ValidationStatus(valid.Status()),
		Timestamp: valid.Timestamp().Unix(),
	}

	res := &api.ValidationResponse{
		Validation: vRes,
	}
	if err := h.crypto.signer.SignValidationResponse(res, h.conf.SharedKeys.Robot.PrivateKey); err != nil {
		return nil, err
	}

	return res, nil
}

func (h externalSrvHandler) StoreContractMessage(ctx context.Context, req *api.ContractMessageStorageRequest) (*api.StorageAck, error) {
	if err := h.crypto.signer.VerifyContractMessageStorageRequestSignature(h.conf.SharedKeys.Robot.PublicKey, req); err != nil {
		return nil, err
	}

	msg := contract.NewEndorsedContractMessage(
		contract.NewMessage(req.ContractMessage.ContractAddress, req.ContractMessage.Method, req.ContractMessage.Parameters, req.ContractMessage.PublicKey, req.ContractMessage.Signature, req.ContractMessage.EmitterSignature),
		h.data.buildEndorsement(req.Endorsement))

	if err := h.services.contractAdd.StoreEndorsedMessage(msg); err != nil {
		return nil, err
	}

	hash, err := h.crypto.hasher.HashEndorsedContractMessage(msg)
	if err != nil {
		return nil, err
	}

	ack := &api.StorageAck{
		StorageHash: hash,
	}
	if err := h.crypto.signer.SignStorageAck(ack, h.conf.SharedKeys.Robot.PrivateKey); err != nil {
		return nil, err
	}
	return ack, nil
}

func (h externalSrvHandler) GetContractState(ctx context.Context, req *api.ContractStateRequest) (*api.ContractStateResponse, error) {
	state, err := h.services.contractLister.GetContractState(req.ContractAddress)
	if err != nil {
		return nil, err
	}

	res := &api.ContractStateResponse{
		Data: state,
	}
	if err := h.crypto.signer.SignContractState(res, h.conf.SharedKeys.Robot.PrivateKey); err != nil {
		return nil, err
	}

	return res, nil
}

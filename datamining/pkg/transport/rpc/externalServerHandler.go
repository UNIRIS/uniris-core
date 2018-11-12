package rpc

import (
	"context"
	"errors"
	"net"

	"github.com/golang/protobuf/ptypes/empty"
	datamining "github.com/uniris/uniris-core/datamining/pkg"
	"github.com/uniris/uniris-core/datamining/pkg/account"
	"github.com/uniris/uniris-core/datamining/pkg/mining"

	"github.com/uniris/uniris-core/datamining/pkg/lock"

	"github.com/uniris/uniris-core/datamining/pkg/system"

	api "github.com/uniris/uniris-core/datamining/api/protobuf-spec"

	accAdding "github.com/uniris/uniris-core/datamining/pkg/account/adding"
	accListing "github.com/uniris/uniris-core/datamining/pkg/account/listing"
)

//Services define the required services
type Services struct {
	lock      lock.Service
	mining    mining.Service
	accAdd    accAdding.Service
	accLister accListing.Service
}

//NewExternalServices creates a new container of required services
func NewExternalServices(lock lock.Service, mine mining.Service, accountAdder accAdding.Service, accountLister accListing.Service) Services {
	return Services{
		lock:      lock,
		mining:    mine,
		accAdd:    accountAdder,
		accLister: accountLister,
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

func (h externalSrvHandler) GetBiometric(ctxt context.Context, req *api.BiometricRequest) (*api.BiometricResponse, error) {
	if err := h.crypto.signer.CheckHashSignature(h.conf.SharedKeys.RobotPublicKey, req.EncryptedPersonHash, req.Signature); err != nil {
		return nil, ErrInvalidSignature
	}

	personHash, err := h.crypto.decrypter.DecryptHash(req.EncryptedPersonHash, h.conf.SharedKeys.RobotPrivateKey)
	if err != nil {
		return nil, ErrInvalidEncryption
	}

	biometric, err := h.services.accLister.GetBiometric(personHash)
	if err != nil {
		return nil, err
	}

	if biometric == nil {
		return nil, errors.New(h.conf.Datamining.Errors.AccountNotExist)
	}

	res := &api.BiometricResponse{
		Data:        h.api.buildBiometricData(biometric),
		Endorsement: h.api.buildEndorsement(biometric.Endorsement()),
	}

	if err := h.crypto.signer.SignBiometricResponse(res, h.conf.SharedKeys.RobotPrivateKey); err != nil {
		return nil, err
	}

	return res, nil
}

func (h externalSrvHandler) GetKeychain(ctxt context.Context, req *api.KeychainRequest) (*api.KeychainResponse, error) {
	if err := h.crypto.signer.CheckHashSignature(h.conf.SharedKeys.RobotPublicKey, req.EncryptedAddress, req.Signature); err != nil {
		return nil, ErrInvalidSignature
	}

	clearAddress, err := h.crypto.decrypter.DecryptHash(req.EncryptedAddress, h.conf.SharedKeys.RobotPrivateKey)
	if err != nil {
		return nil, ErrInvalidEncryption
	}

	keychain, err := h.services.accLister.GetLastKeychain(clearAddress)
	if err != nil {
		return nil, err
	}

	if keychain == nil {
		return nil, errors.New(h.conf.Datamining.Errors.AccountNotExist)
	}

	res := &api.KeychainResponse{
		Data:        h.api.buildKeychainData(keychain),
		Endorsement: h.api.buildEndorsement(keychain.Endorsement()),
	}

	if err := h.crypto.signer.SignKeychainResponse(res, h.conf.SharedKeys.RobotPrivateKey); err != nil {
		return nil, err
	}

	return res, nil
}

func (h externalSrvHandler) LeadKeychainMining(ctx context.Context, req *api.KeychainLeadRequest) (*empty.Empty, error) {
	if err := h.crypto.signer.CheckKeychainLeadRequestSignature(h.conf.SharedKeys.RobotPublicKey, req); err != nil {
		return nil, ErrInvalidSignature
	}

	keychainSigLess, err := h.crypto.decrypter.DecryptKeychainData(req.EncryptedKeychainData, h.conf.SharedKeys.RobotPrivateKey)
	if err != nil {
		return nil, ErrInvalidEncryption
	}

	keychain := account.NewKeychainData(
		keychainSigLess.CipherAddrRobot(),
		keychainSigLess.CipherWallet(),
		keychainSigLess.PersonPublicKey(),
		keychainSigLess.BiodPublicKey(),
		account.NewSignatures(req.SignatureKeychainData.Biod, req.SignatureKeychainData.Person),
	)

	clearaddr, err := h.crypto.decrypter.DecryptHash(keychain.CipherAddrRobot(), h.conf.SharedKeys.RobotPrivateKey)
	if err != nil {
		return nil, ErrInvalidEncryption
	}

	pp := make([]datamining.Peer, 0)
	for _, p := range req.ValidatorPeerIPs {
		pp = append(pp, datamining.Peer{IP: net.ParseIP(p)})
	}
	vPool := datamining.NewPool(pp...)
	if err := h.services.mining.LeadMining(req.TransactionHash, clearaddr, keychain, vPool, mining.KeychainTransaction, keychain.Signatures().Biod()); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (h externalSrvHandler) LeadBiometricMining(ctx context.Context, req *api.BiometricLeadRequest) (*empty.Empty, error) {
	if err := h.crypto.signer.CheckBiometricLeadRequestSignature(h.conf.SharedKeys.RobotPublicKey, req); err != nil {
		return nil, ErrInvalidSignature
	}

	bioSigLess, err := h.crypto.decrypter.DecryptBiometricData(req.EncryptedBioData, h.conf.SharedKeys.RobotPrivateKey)
	if err != nil {
		return nil, ErrInvalidEncryption
	}

	biometricData := account.NewBiometricData(
		bioSigLess.PersonHash(),
		bioSigLess.CipherAddrRobot(),
		bioSigLess.CipherAddrPerson(),
		bioSigLess.CipherAESKey(),
		bioSigLess.PersonPublicKey(),
		bioSigLess.BiodPublicKey(),
		account.NewSignatures(req.SignatureBioData.Biod, req.SignatureBioData.Person),
	)

	clearaddr, err := h.crypto.decrypter.DecryptHash(biometricData.CipherAddrRobot(), h.conf.SharedKeys.RobotPrivateKey)
	if err != nil {
		return nil, ErrInvalidEncryption
	}

	pp := make([]datamining.Peer, 0)
	for _, p := range req.ValidatorPeerIPs {
		pp = append(pp, datamining.Peer{IP: net.ParseIP(p)})
	}
	vPool := datamining.NewPool(pp...)
	if err := h.services.mining.LeadMining(req.TransactionHash, clearaddr, biometricData, vPool, mining.BiometricTransaction, biometricData.Signatures().Biod()); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (h externalSrvHandler) LockTransaction(ctx context.Context, req *api.LockRequest) (*api.LockAck, error) {
	if err := h.crypto.signer.CheckLockRequestSignature(h.conf.SharedKeys.RobotPublicKey, req); err != nil {
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
	if err := h.crypto.signer.SignLockAck(ack, h.conf.SharedKeys.RobotPrivateKey); err != nil {
		return nil, err
	}
	return ack, nil
}

func (h externalSrvHandler) UnlockTransaction(ctx context.Context, req *api.LockRequest) (*api.LockAck, error) {

	if err := h.crypto.signer.CheckLockRequestSignature(h.conf.SharedKeys.RobotPublicKey, req); err != nil {
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
	if err := h.crypto.signer.SignLockAck(ack, h.conf.SharedKeys.RobotPrivateKey); err != nil {
		return nil, err
	}
	return ack, nil
}

func (h externalSrvHandler) ValidateKeychain(ctx context.Context, req *api.KeychainValidationRequest) (*api.ValidationResponse, error) {
	if err := h.crypto.signer.CheckKeychainValidationRequestSignature(h.conf.SharedKeys.RobotPublicKey, req); err != nil {
		return nil, err
	}

	valid, err := h.services.mining.Validate(req.TransactionHash, h.data.buildKeychainData(req.Data), mining.KeychainTransaction)
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
	if err := h.crypto.signer.SignValidationResponse(res, h.conf.SharedKeys.RobotPrivateKey); err != nil {
		return nil, err
	}

	return res, nil
}

func (h externalSrvHandler) ValidateBiometric(ctx context.Context, req *api.BiometricValidationRequest) (*api.ValidationResponse, error) {
	if err := h.crypto.signer.CheckBiometricValidationRequestSignature(h.conf.SharedKeys.RobotPublicKey, req); err != nil {
		return nil, err
	}

	valid, err := h.services.mining.Validate(req.TransactionHash, h.data.buildBiometricData(req.Data), mining.BiometricTransaction)
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
	if err := h.crypto.signer.SignValidationResponse(res, h.conf.SharedKeys.RobotPrivateKey); err != nil {
		return nil, err
	}

	return res, nil
}

func (h externalSrvHandler) StoreKeychain(ctx context.Context, req *api.KeychainStorageRequest) (*api.StorageAck, error) {
	if err := h.crypto.signer.CheckKeychainStorageRequestSignature(h.conf.SharedKeys.RobotPublicKey, req); err != nil {
		return nil, err
	}

	clearaddr, err := h.crypto.decrypter.DecryptHash(req.Data.CipherAddrRobot, h.conf.SharedKeys.RobotPrivateKey)
	if err != nil {
		return nil, ErrInvalidEncryption
	}

	keychain := account.NewKeychain(clearaddr,
		h.data.buildKeychainData(req.Data),
		h.data.buildEndorsement(req.Endorsement))

	if err := h.services.accAdd.StoreKeychain(keychain); err != nil {
		return nil, err
	}

	hash, err := h.crypto.hasher.HashKeychain(keychain)
	if err != nil {
		return nil, err
	}

	ack := &api.StorageAck{
		StorageHash: hash,
	}
	if err := h.crypto.signer.SignStorageAck(ack, h.conf.SharedKeys.RobotPrivateKey); err != nil {
		return nil, err
	}
	return ack, nil
}

func (h externalSrvHandler) StoreBiometric(ctx context.Context, req *api.BiometricStorageRequest) (*api.StorageAck, error) {
	if err := h.crypto.signer.CheckBiometricStorageRequestSignature(h.conf.SharedKeys.RobotPublicKey, req); err != nil {
		return nil, err
	}

	biometric := account.NewBiometric(h.data.buildBiometricData(req.Data), h.data.buildEndorsement(req.Endorsement))
	if err := h.services.accAdd.StoreBiometric(biometric); err != nil {
		return nil, err
	}

	hash, err := h.crypto.hasher.HashBiometric(biometric)
	if err != nil {
		return nil, err
	}

	ack := &api.StorageAck{
		StorageHash: hash,
	}
	if err := h.crypto.signer.SignStorageAck(ack, h.conf.SharedKeys.RobotPrivateKey); err != nil {
		return nil, err
	}
	return ack, nil
}

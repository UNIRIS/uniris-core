package rpc

import (
	"errors"

	"github.com/uniris/uniris-core/datamining/pkg/account"
	"github.com/uniris/uniris-core/datamining/pkg/system"
	"golang.org/x/net/context"

	api "github.com/uniris/uniris-core/datamining/api/protobuf-spec"
)

type internalSrvHandler struct {
	pR       account.PoolRequester
	aiClient AIClient
	crypto   Crypto
	conf     system.UnirisConfig
}

//NewInternalServerHandler create a new GRPC server handler for account
func NewInternalServerHandler(pR PoolRequester, aiClient AIClient, crypto Crypto, conf system.UnirisConfig) api.InternalServer {
	return internalSrvHandler{
		pR:       pR,
		aiClient: aiClient,
		crypto:   crypto,
		conf:     conf,
	}
}

//GetAccount implements the protobuf GetAccount request handler
func (s internalSrvHandler) GetAccount(ctx context.Context, req *api.AccountSearchRequest) (*api.AccountSearchResult, error) {
	idHash, err := s.crypto.decrypter.DecryptHash(req.EncryptedIDHash, s.conf.SharedKeys.RobotPrivateKey)
	if err != nil {
		return nil, ErrInvalidEncryption
	}

	idPool, err := s.aiClient.GetStoragePool(idHash)
	if err != nil {
		return nil, err
	}

	biometric, err := s.pR.RequestID(idPool, req.EncryptedIDHash)
	if err != nil {
		return nil, err
	}

	clearAddr, err := s.crypto.decrypter.DecryptHash(biometric.EncryptedAddrByRobot(), s.conf.SharedKeys.RobotPrivateKey)
	if err != nil {
		return nil, ErrInvalidEncryption
	}

	keychainPool, err := s.aiClient.GetStoragePool(clearAddr)
	if err != nil {
		return nil, err
	}

	keychain, err := s.pR.RequestKeychain(keychainPool, biometric.EncryptedAddrByRobot())
	if err != nil {
		return nil, err
	}

	if keychain == nil {
		return nil, errors.New(s.conf.Datamining.Errors.AccountNotExist)
	}

	res := &api.AccountSearchResult{
		EncryptedAESkey:  biometric.EncryptedAESKey(),
		EncryptedWallet:  keychain.EncryptedWallet(),
		EncryptedAddress: biometric.EncryptedAddrByID(),
	}
	if err := s.crypto.signer.SignAccountSearchResult(res, s.conf.SharedKeys.RobotPrivateKey); err != nil {
		return nil, err
	}

	return res, nil
}

func (s internalSrvHandler) CreateKeychain(ctx context.Context, req *api.KeychainCreationRequest) (*api.CreationResult, error) {
	keychain, err := s.crypto.decrypter.DecryptKeychain(req.EncryptedKeychain, s.conf.SharedKeys.RobotPrivateKey)
	if err != nil {
		return nil, err
	}

	txHash, err := s.crypto.hasher.HashKeychain(keychain)
	if err != nil {
		return nil, err
	}

	master, err := s.aiClient.GetMasterPeer(txHash)
	if err != nil {
		return nil, err
	}
	validPool, err := s.aiClient.GetValidationPool(txHash)
	if err != nil {
		return nil, err
	}

	extCli := NewExternalClient(s.crypto, s.conf)
	go extCli.LeadKeychainMining(master.IP.String(), txHash, req.EncryptedKeychain, validPool.Peers().IPs())

	res := &api.CreationResult{
		TransactionHash: txHash,
		MasterPeerIP:    master.IP.String(),
	}

	if err := s.crypto.signer.SignCreationResult(res, s.conf.SharedKeys.RobotPrivateKey); err != nil {
		return nil, err
	}

	return res, nil
}

func (s internalSrvHandler) CreateID(ctx context.Context, req *api.IDCreationRequest) (*api.CreationResult, error) {
	id, err := s.crypto.decrypter.DecryptID(req.EncryptedID, s.conf.SharedKeys.RobotPrivateKey)
	if err != nil {
		return nil, err
	}

	txHash, err := s.crypto.hasher.HashID(id)
	if err != nil {
		return nil, err
	}

	master, err := s.aiClient.GetMasterPeer(txHash)
	if err != nil {
		return nil, err
	}
	validPool, err := s.aiClient.GetValidationPool(txHash)
	if err != nil {
		return nil, err
	}

	extCli := NewExternalClient(s.crypto, s.conf)
	go extCli.LeadIDMining(master.IP.String(), txHash, req.EncryptedID, validPool.Peers().IPs())

	res := &api.CreationResult{
		TransactionHash: txHash,
		MasterPeerIP:    master.IP.String(),
	}

	if err := s.crypto.signer.SignCreationResult(res, s.conf.SharedKeys.RobotPrivateKey); err != nil {
		return nil, err
	}

	return res, nil
}

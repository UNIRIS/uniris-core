package rpc

import (
	"errors"
	"log"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/uniris/uniris-core/datamining/pkg/account"
	"github.com/uniris/uniris-core/datamining/pkg/mining"
	"github.com/uniris/uniris-core/datamining/pkg/system"
	"golang.org/x/net/context"

	api "github.com/uniris/uniris-core/datamining/api/protobuf-spec"

	emListing "github.com/uniris/uniris-core/datamining/pkg/emitter/listing"
)

type internalSrvHandler struct {
	pR       account.PoolRequester
	aiClient AIClient
	extCli   ExternalClient
	crypto   Crypto
	conf     system.UnirisConfig
	emLister emListing.Service
	poolF    mining.PoolFinder
}

//NewInternalServerHandler create a new GRPC server handler for account
func NewInternalServerHandler(emLister emListing.Service, pR PoolRequester, pF mining.PoolFinder, aiClient AIClient, extCli ExternalClient, crypto Crypto, conf system.UnirisConfig) api.InternalServer {
	return internalSrvHandler{
		emLister: emLister,
		pR:       pR,
		poolF:    pF,
		aiClient: aiClient,
		extCli:   extCli,
		crypto:   crypto,
		conf:     conf,
	}
}

//GetAccount implements the protobuf GetAccount request handler
func (s internalSrvHandler) GetAccount(ctx context.Context, req *api.AccountSearchRequest) (*api.AccountSearchResult, error) {
	idHash, err := s.crypto.decrypter.DecryptHash(req.EncryptedIDHash, s.conf.SharedKeys.Robot.PrivateKey)
	if err != nil {
		return nil, ErrInvalidEncryption
	}

	idPool, err := s.aiClient.GetStoragePool(idHash)
	if err != nil {
		return nil, err
	}

	id, err := s.pR.RequestID(idPool, req.EncryptedIDHash)
	if err != nil {
		return nil, err
	}

	clearAddr, err := s.crypto.decrypter.DecryptHash(id.EncryptedAddrByRobot(), s.conf.SharedKeys.Robot.PrivateKey)
	if err != nil {
		return nil, ErrInvalidEncryption
	}

	keychainPool, err := s.aiClient.GetStoragePool(clearAddr)
	if err != nil {
		return nil, err
	}

	keychain, err := s.pR.RequestKeychain(keychainPool, id.EncryptedAddrByRobot())
	if err != nil {
		return nil, err
	}

	if keychain == nil {
		return nil, errors.New(s.conf.Services.Datamining.Errors.AccountNotExist)
	}

	res := &api.AccountSearchResult{
		EncryptedAESkey:  id.EncryptedAESKey(),
		EncryptedWallet:  keychain.EncryptedWallet(),
		EncryptedAddress: id.EncryptedAddrByID(),
	}

	if err := s.crypto.signer.SignAccountSearchResult(res, s.conf.SharedKeys.Robot.PrivateKey); err != nil {
		return nil, err
	}

	return res, nil
}

func (s internalSrvHandler) CreateKeychain(ctx context.Context, req *api.KeychainCreationRequest) (*api.CreationResult, error) {
	keychain, err := s.crypto.decrypter.DecryptKeychain(req.EncryptedKeychain, s.conf.SharedKeys.Robot.PrivateKey)
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

	if err := s.extCli.LeadKeychainMining(master.IP.String(), txHash, req.EncryptedKeychain, validPool.Peers().IPs()); err != nil {
		log.Print(err.Error())
		return nil, err
	}

	res := &api.CreationResult{
		TransactionHash: txHash,
		MasterPeerIP:    master.IP.String(),
	}

	if err := s.crypto.signer.SignCreationResult(res, s.conf.SharedKeys.Robot.PrivateKey); err != nil {
		return nil, err
	}

	return res, nil
}

func (s internalSrvHandler) CreateID(ctx context.Context, req *api.IDCreationRequest) (*api.CreationResult, error) {
	id, err := s.crypto.decrypter.DecryptID(req.EncryptedID, s.conf.SharedKeys.Robot.PrivateKey)
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

	if err := s.extCli.LeadIDMining(master.IP.String(), txHash, req.EncryptedID, validPool.Peers().IPs()); err != nil {
		return nil, err
	}

	res := &api.CreationResult{
		TransactionHash: txHash,
		MasterPeerIP:    master.IP.String(),
	}

	if err := s.crypto.signer.SignCreationResult(res, s.conf.SharedKeys.Robot.PrivateKey); err != nil {
		return nil, err
	}

	return res, nil
}

func (s internalSrvHandler) IsEmitterAuthorized(ctx context.Context, req *api.AuthorizationRequest) (*api.AuthorizationResponse, error) {
	if err := s.emLister.IsEmitterAuthorized(req.PublicKey); err != nil {
		if err == emListing.ErrUnauthorizedEmitter {
			return &api.AuthorizationResponse{
				Status: false,
			}, nil
		}
	}

	return &api.AuthorizationResponse{
		Status: true,
	}, nil
}

func (s internalSrvHandler) GetSharedKeys(ctx context.Context, req *empty.Empty) (*api.SharedKeysResult, error) {
	kps, err := s.emLister.ListSharedEmitterKeyPairs()
	if err != nil {
		return nil, err
	}

	emiterKeys := make([]*api.SharedKeyPair, 0)
	for _, kp := range kps {
		emiterKeys = append(emiterKeys, &api.SharedKeyPair{
			EncryptedPrivateKey: kp.EncryptedPrivateKey,
			PublicKey:           kp.PublicKey,
		})
	}

	return &api.SharedKeysResult{
		EmitterKeys:     emiterKeys,
		RobotPublicKey:  s.conf.SharedKeys.Robot.PublicKey,
		RobotPrivateKey: s.conf.SharedKeys.Robot.PrivateKey,
	}, nil
}

func (s internalSrvHandler) GetTransactionStatus(ctx context.Context, req *api.TransactionStatusRequest) (*api.TransactionStatusResponse, error) {
	addr, err := s.crypto.decrypter.DecryptHash(req.Address, s.conf.SharedKeys.Robot.PrivateKey)
	if err != nil {
		return nil, err
	}

	storagePool, err := s.poolF.FindStoragePool(addr)
	if err != nil {
		return nil, err
	}

	statuses := make([]mining.TransactionStatus, 0)
	for _, p := range storagePool.Peers().IPs() {
		status, err := s.extCli.GetTransactionStatus(p, req.Address, req.Hash)
		if err != nil {
			log.Print(err.Error())
		}
		statuses = append(statuses, status)
	}

	if len(statuses) > 0 {
		//TODO: Provide consensus of the data retrieval
		return &api.TransactionStatusResponse{
			Status: api.TransactionStatusResponse_TransactionStatus(statuses[0]),
		}, nil
	}

	return nil, nil

}

func (s internalSrvHandler) CreateContract(ctx context.Context, req *api.ContractCreationRequest) (*api.CreationResult, error) {
	txHash, err := s.crypto.hasher.HashAPIContract(req.Contract)
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

	if err := s.extCli.LeadContractMining(master.IP.String(), txHash, req.Contract, validPool.Peers().IPs()); err != nil {
		log.Print(err.Error())
		return nil, err
	}

	res := &api.CreationResult{
		TransactionHash: txHash,
		MasterPeerIP:    master.IP.String(),
	}

	if err := s.crypto.signer.SignCreationResult(res, s.conf.SharedKeys.Robot.PrivateKey); err != nil {
		return nil, err
	}

	return res, nil
}

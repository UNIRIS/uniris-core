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
	personHash, err := s.crypto.decrypter.DecryptHash(req.EncryptedHashPerson, s.conf.SharedKeys.RobotPrivateKey)
	if err != nil {
		return nil, ErrInvalidEncryption
	}

	biometricPool, err := s.aiClient.GetStoragePool(personHash)
	if err != nil {
		return nil, err
	}

	biometric, err := s.pR.RequestBiometric(biometricPool, req.EncryptedHashPerson)
	if err != nil {
		return nil, err
	}

	clearAddr, err := s.crypto.decrypter.DecryptHash(biometric.CipherAddrRobot(), s.conf.SharedKeys.RobotPrivateKey)
	if err != nil {
		return nil, ErrInvalidEncryption
	}

	keychainPool, err := s.aiClient.GetStoragePool(clearAddr)
	if err != nil {
		return nil, err
	}

	keychain, err := s.pR.RequestKeychain(keychainPool, biometric.CipherAddrRobot())
	if err != nil {
		return nil, err
	}

	if keychain == nil {
		return nil, errors.New(s.conf.Datamining.Errors.AccountNotExist)
	}

	res := &api.AccountSearchResult{
		EncryptedAESkey:  biometric.CipherAESKey(),
		EncryptedWallet:  keychain.CipherWallet(),
		EncryptedAddress: biometric.CipherAddrPerson(),
	}
	if err := s.crypto.signer.SignAccountSearchResult(res, s.conf.SharedKeys.RobotPrivateKey); err != nil {
		return nil, err
	}

	return res, nil
}

func (s internalSrvHandler) CreateKeychain(ctx context.Context, req *api.KeychainCreationRequest) (*api.CreationResult, error) {
	cKeychain, err := s.crypto.decrypter.DecryptKeychainData(req.EncryptedKeychainData, s.conf.SharedKeys.RobotPrivateKey)
	if err != nil {
		return nil, err
	}

	keychainData := account.NewKeychainData(
		cKeychain.CipherAddrRobot(),
		cKeychain.CipherWallet(),
		cKeychain.PersonPublicKey(),
		cKeychain.BiodPublicKey(),
		account.NewSignatures(req.SignatureKeychainData.Biod, req.SignatureKeychainData.Person),
	)

	txHash, err := s.crypto.hasher.HashKeychainData(keychainData)
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
	go extCli.LeadKeychainMining(master.IP.String(), txHash, req.EncryptedKeychainData, req.SignatureKeychainData, validPool.Peers().IPs())

	res := &api.CreationResult{
		TransactionHash: txHash,
		MasterPeerIP:    master.IP.String(),
	}

	if err := s.crypto.signer.SignCreationResult(res, s.conf.SharedKeys.RobotPrivateKey); err != nil {
		return nil, err
	}

	return res, nil
}

func (s internalSrvHandler) CreateBiometric(ctx context.Context, req *api.BiometricCreationRequest) (*api.CreationResult, error) {
	cBio, err := s.crypto.decrypter.DecryptBiometricData(req.EncryptedBiometricData, s.conf.SharedKeys.RobotPrivateKey)
	if err != nil {
		return nil, err
	}

	bio := account.NewBiometricData(
		cBio.PersonHash(),
		cBio.CipherAddrRobot(),
		cBio.CipherAddrPerson(),
		cBio.CipherAESKey(),
		cBio.PersonPublicKey(),
		cBio.BiodPublicKey(),
		account.NewSignatures(req.SignatureBiometricData.Biod, req.SignatureBiometricData.Person),
	)

	txHash, err := s.crypto.hasher.HashBiometricData(bio)
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
	go extCli.LeadBiometricMining(master.IP.String(), txHash, req.EncryptedBiometricData, req.SignatureBiometricData, validPool.Peers().IPs())

	res := &api.CreationResult{
		TransactionHash: txHash,
		MasterPeerIP:    master.IP.String(),
	}

	if err := s.crypto.signer.SignCreationResult(res, s.conf.SharedKeys.RobotPrivateKey); err != nil {
		return nil, err
	}

	return res, nil
}

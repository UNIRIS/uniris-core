package rpc

import (
	"errors"
	"fmt"

	"github.com/uniris/uniris-core/datamining/pkg/contract"
	"github.com/uniris/uniris-core/datamining/pkg/lock"

	api "github.com/uniris/uniris-core/datamining/api/protobuf-spec"
	"github.com/uniris/uniris-core/datamining/pkg/account"
	"github.com/uniris/uniris-core/datamining/pkg/mining"
	"github.com/uniris/uniris-core/datamining/pkg/system"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

//ExternalClient define methods the client for the External GRPC has to define
type ExternalClient interface {

	//LeadKeychainMining process the keychain mining as master robot
	LeadKeychainMining(ip string, txHash string, encData string, validators []string) error

	//LeadIDMining process the ID mining as master robot
	LeadIDMining(ip string, txHash string, encData string, validators []string) error

	//RequestID requests a peer to retrive a ID data from a given encrypted ID hash
	RequestID(ip string, encPersonHash string) (account.EndorsedID, error)

	//RequestKeychain requests a peer to retrieve the last keychain from a given encrypted address
	RequestKeychain(ip string, encAddress string) (account.EndorsedKeychain, error)

	//RequestLock requests a peer to lock a given transaction
	RequestLock(ip string, txLock lock.TransactionLock) error

	//RequestLock requests a peer to unlock a given transaction
	RequestUnlock(ip string, txLock lock.TransactionLock) error

	//RequestValidations requests a peer to process validation/mining as a slave robot
	RequestValidation(ip string, txType mining.TransactionType, txHash string, data interface{}) (mining.Validation, error)

	//RequestStorage requests a peer to store the transaction
	RequestStorage(ip string, txType mining.TransactionType, data interface{}, end mining.Endorsement) error

	//GetTransactionStatus requests a peer to retrieve transaction status
	GetTransactionStatus(ip string, addr string, txHash string) (mining.TransactionStatus, error)

	LeadContractMining(ip string, txHash string, contract *api.Contract, validators []string) error
	LeadContractMessageMining(ip string, txHash string, message *api.ContractMessage, validators []string) error
	GetContractState(ip string, address string) (*api.ContractStateResponse, error)
}

type externalClient struct {
	crypto Crypto
	conf   system.UnirisConfig
	data   dataBuilder
	api    apiBuilder
}

//NewExternalClient create a GRPC implementation of the external client
func NewExternalClient(crypto Crypto, conf system.UnirisConfig) ExternalClient {
	return externalClient{
		crypto: crypto,
		conf:   conf,
		data:   dataBuilder{},
		api:    apiBuilder{},
	}
}

func (c externalClient) LeadKeychainMining(ip string, txHash string, encData string, validators []string) error {
	serverAddr := fmt.Sprintf("%s:%d", ip, c.conf.Services.Datamining.ExternalPort)
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	defer conn.Close()

	if err != nil {
		return err
	}

	client := api.NewExternalClient(conn)

	req := &api.KeychainLeadRequest{
		EncryptedKeychain: encData,
		TransactionHash:   txHash,
		ValidatorPeerIPs:  validators,
	}
	if err := c.crypto.signer.SignKeychainLeadRequest(req, c.conf.SharedKeys.Robot.PrivateKey); err != nil {
		return err
	}

	_, err = client.LeadKeychainMining(context.Background(), req)
	if err != nil {
		s, _ := status.FromError(err)
		return errors.New(s.Message())
	}
	return nil
}

func (c externalClient) LeadIDMining(ip string, txHash string, encData string, validators []string) error {
	serverAddr := fmt.Sprintf("%s:%d", ip, c.conf.Services.Datamining.ExternalPort)
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	defer conn.Close()

	if err != nil {
		return err
	}

	client := api.NewExternalClient(conn)

	req := &api.IDLeadRequest{
		EncryptedID:      encData,
		TransactionHash:  txHash,
		ValidatorPeerIPs: validators,
	}
	if err := c.crypto.signer.SignIDLeadRequest(req, c.conf.SharedKeys.Robot.PrivateKey); err != nil {
		return err
	}
	_, err = client.LeadIDMining(context.Background(), req)
	if err != nil {
		s, _ := status.FromError(err)
		return errors.New(s.Message())
	}
	return nil
}

func (c externalClient) RequestID(ip string, encIDHash string) (account.EndorsedID, error) {
	serverAddr := fmt.Sprintf("%s:%d", ip, c.conf.Services.Datamining.ExternalPort)
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	defer conn.Close()

	if err != nil {
		return nil, err
	}

	client := api.NewExternalClient(conn)

	sigReq, err := c.crypto.signer.SignHash(encIDHash, c.conf.SharedKeys.Robot.PrivateKey)
	if err != nil {
		return nil, err
	}

	res, err := client.GetID(context.Background(), &api.IDRequest{
		EncryptedIDHash: encIDHash,
		Signature:       sigReq,
	})
	if err != nil {
		s, _ := status.FromError(err)
		return nil, errors.New(s.Message())
	}

	if err := c.crypto.signer.VerifyIDResponseSignature(c.conf.SharedKeys.Robot.PublicKey, res); err != nil {
		return nil, err
	}

	return account.NewEndorsedID(
		c.data.buildID(res.Data),
		c.data.buildEndorsement(res.Endorsement)), nil
}

func (c externalClient) RequestKeychain(ip string, encAddress string) (account.EndorsedKeychain, error) {
	serverAddr := fmt.Sprintf("%s:%d", ip, c.conf.Services.Datamining.ExternalPort)
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	defer conn.Close()

	if err != nil {
		return nil, err
	}

	client := api.NewExternalClient(conn)

	sigReq, err := c.crypto.signer.SignHash(encAddress, c.conf.SharedKeys.Robot.PrivateKey)
	if err != nil {
		return nil, err
	}

	res, err := client.GetKeychain(context.Background(), &api.KeychainRequest{
		EncryptedAddress: encAddress,
		Signature:        sigReq,
	})
	if err != nil {
		s, _ := status.FromError(err)
		return nil, errors.New(s.Message())
	}

	clearaddr, err := c.crypto.decrypter.DecryptHash(res.Data.EncryptedAddrByRobot, c.conf.SharedKeys.Robot.PrivateKey)
	if err != nil {
		return nil, err
	}

	if err := c.crypto.signer.VerifyKeychainResponseSignature(c.conf.SharedKeys.Robot.PublicKey, res); err != nil {
		return nil, err
	}

	return account.NewEndorsedKeychain(
		clearaddr,
		c.data.buildKeychain(res.Data),
		c.data.buildEndorsement(res.Endorsement),
	), nil
}

func (c externalClient) RequestLock(ip string, txLock lock.TransactionLock) error {
	serverAddr := fmt.Sprintf("%s:%d", ip, c.conf.Services.Datamining.ExternalPort)
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	defer conn.Close()

	if err != nil {
		return err
	}

	client := api.NewExternalClient(conn)

	lockReq := &api.LockRequest{
		MasterRobotKey:  txLock.MasterRobotKey,
		TransactionHash: txLock.TxHash,
		Address:         txLock.Address,
	}
	if err := c.crypto.signer.SignLockRequest(lockReq, c.conf.SharedKeys.Robot.PrivateKey); err != nil {
		return err
	}

	res, err := client.LockTransaction(context.Background(), lockReq)
	if err != nil {
		s, _ := status.FromError(err)
		return errors.New(s.Message())
	}

	if err := c.crypto.signer.VerifyLockAckSignature(c.conf.SharedKeys.Robot.PublicKey, res); err != nil {
		return err
	}

	//TODO: Verify res.LockHash

	return nil
}

func (c externalClient) RequestUnlock(ip string, txLock lock.TransactionLock) error {
	serverAddr := fmt.Sprintf("%s:%d", ip, c.conf.Services.Datamining.ExternalPort)
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	defer conn.Close()

	if err != nil {
		return err
	}

	client := api.NewExternalClient(conn)

	lockReq := &api.LockRequest{
		MasterRobotKey:  txLock.MasterRobotKey,
		TransactionHash: txLock.TxHash,
		Address:         txLock.Address,
	}
	if err := c.crypto.signer.SignLockRequest(lockReq, c.conf.SharedKeys.Robot.PrivateKey); err != nil {
		return err
	}

	res, err := client.UnlockTransaction(context.Background(), lockReq)
	if err != nil {
		s, _ := status.FromError(err)
		return errors.New(s.Message())
	}

	if err := c.crypto.signer.VerifyLockAckSignature(c.conf.SharedKeys.Robot.PublicKey, res); err != nil {
		return err
	}

	//TODO: Verify res.LockHash

	return nil
}

func (c externalClient) RequestValidation(ip string, txType mining.TransactionType, txHash string, data interface{}) (mining.Validation, error) {
	serverAddr := fmt.Sprintf("%s:%d", ip, c.conf.Services.Datamining.ExternalPort)
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	defer conn.Close()

	if err != nil {
		return nil, err
	}

	client := api.NewExternalClient(conn)

	switch txType {
	case mining.KeychainTransaction:
		return c.validateKeychain(client, txHash, c.api.buildKeychain(data.(account.Keychain)))
	case mining.IDTransaction:
		return c.validateID(client, txHash, c.api.buildID(data.(account.ID)))
	case mining.ContractTransaction:
		contract := data.(contract.Contract)
		return c.validateContract(client, txHash, &api.Contract{
			Address:          contract.Address(),
			Code:             contract.Code(),
			Event:            contract.Event(),
			PublicKey:        contract.PublicKey(),
			Signature:        contract.Signature(),
			EmitterSignature: contract.EmitterSignature(),
		})
	case mining.ContractMessageTransaction:
		contractMsg := data.(contract.Message)
		return c.validateContractMessage(client, txHash, &api.ContractMessage{
			ContractAddress:  contractMsg.ContractAddress(),
			Method:           contractMsg.Method(),
			Parameters:       contractMsg.Parameters(),
			PublicKey:        contractMsg.PublicKey(),
			Signature:        contractMsg.Signature(),
			EmitterSignature: contractMsg.EmitterSignature(),
		})
	}

	return nil, errors.New("Unsupported transaction type")
}

func (c externalClient) RequestStorage(ip string, txType mining.TransactionType, data interface{}, end mining.Endorsement) error {
	serverAddr := fmt.Sprintf("%s:%d", ip, c.conf.Services.Datamining.ExternalPort)
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	defer conn.Close()

	if err != nil {
		return err
	}

	client := api.NewExternalClient(conn)

	endorsement := c.api.buildEndorsement(end)

	switch txType {
	case mining.KeychainTransaction:
		return c.storeKeychain(client, c.api.buildKeychain(data.(account.Keychain)), endorsement)
	case mining.IDTransaction:
		return c.storeID(client, c.api.buildID(data.(account.ID)), endorsement)
	case mining.ContractTransaction:
		contract := data.(contract.Contract)
		return c.storeContract(client, &api.Contract{
			Address:          contract.Address(),
			Code:             contract.Code(),
			Event:            contract.Event(),
			PublicKey:        contract.PublicKey(),
			Signature:        contract.Signature(),
			EmitterSignature: contract.EmitterSignature(),
		}, endorsement)
	case mining.ContractMessageTransaction:
		contractMsg := data.(contract.Message)
		return c.storeContractMessage(client, &api.ContractMessage{
			ContractAddress:  contractMsg.ContractAddress(),
			Method:           contractMsg.Method(),
			Parameters:       contractMsg.Parameters(),
			PublicKey:        contractMsg.PublicKey(),
			Signature:        contractMsg.Signature(),
			EmitterSignature: contractMsg.EmitterSignature(),
		}, endorsement)
	}

	return nil
}

func (c externalClient) validateKeychain(client api.ExternalClient, txHash string, kc *api.Keychain) (mining.Validation, error) {
	req := &api.KeychainValidationRequest{
		Data:            kc,
		TransactionHash: txHash,
	}
	if err := c.crypto.signer.SignKeychainValidationRequestSignature(req, c.conf.SharedKeys.Robot.PrivateKey); err != nil {
		return nil, err
	}

	res, err := client.ValidateKeychain(context.Background(), req)
	if err != nil {
		s, _ := status.FromError(err)
		return nil, errors.New(s.Message())
	}

	if err := c.crypto.signer.VerifyValidationResponseSignature(c.conf.SharedKeys.Robot.PublicKey, res); err != nil {
		return nil, err
	}

	return c.data.buildValidation(res.Validation), nil
}

func (c externalClient) validateID(client api.ExternalClient, txHash string, id *api.ID) (mining.Validation, error) {
	req := &api.IDValidationRequest{
		Data:            id,
		TransactionHash: txHash,
	}
	if err := c.crypto.signer.SignIDValidationRequestSignature(req, c.conf.SharedKeys.Robot.PrivateKey); err != nil {
		return nil, err
	}
	res, err := client.ValidateID(context.Background(), req)
	if err != nil {
		s, _ := status.FromError(err)
		return nil, errors.New(s.Message())
	}

	if err := c.crypto.signer.VerifyValidationResponseSignature(c.conf.SharedKeys.Robot.PublicKey, res); err != nil {
		return nil, err
	}

	return c.data.buildValidation(res.Validation), nil
}

func (c externalClient) validateContract(client api.ExternalClient, txHash string, contract *api.Contract) (mining.Validation, error) {
	req := &api.ContractValidationRequest{
		Contract:        contract,
		TransactionHash: txHash,
	}
	if err := c.crypto.signer.SignContractValidationRequestSignature(req, c.conf.SharedKeys.Robot.PrivateKey); err != nil {
		return nil, err
	}
	res, err := client.ValidateContract(context.Background(), req)
	if err != nil {
		s, _ := status.FromError(err)
		return nil, errors.New(s.Message())
	}

	if err := c.crypto.signer.VerifyValidationResponseSignature(c.conf.SharedKeys.Robot.PublicKey, res); err != nil {
		return nil, err
	}

	return c.data.buildValidation(res.Validation), nil
}

func (c externalClient) validateContractMessage(client api.ExternalClient, txHash string, contractMsg *api.ContractMessage) (mining.Validation, error) {
	req := &api.ContractMessageValidationRequest{
		ContractMessage: contractMsg,
		TransactionHash: txHash,
	}
	if err := c.crypto.signer.SignContractMessageValidationRequestSignature(req, c.conf.SharedKeys.Robot.PrivateKey); err != nil {
		return nil, err
	}
	res, err := client.ValidateContractMessage(context.Background(), req)
	if err != nil {
		s, _ := status.FromError(err)
		return nil, errors.New(s.Message())
	}

	if err := c.crypto.signer.VerifyValidationResponseSignature(c.conf.SharedKeys.Robot.PublicKey, res); err != nil {
		return nil, err
	}

	return c.data.buildValidation(res.Validation), nil
}

func (c externalClient) storeKeychain(client api.ExternalClient, kc *api.Keychain, end *api.Endorsement) error {
	req := &api.KeychainStorageRequest{
		Data:        kc,
		Endorsement: end,
	}
	if err := c.crypto.signer.SignKeychainStorageRequestSignature(req, c.conf.SharedKeys.Robot.PrivateKey); err != nil {
		return err
	}

	res, err := client.StoreKeychain(context.Background(), req)
	if err != nil {
		s, _ := status.FromError(err)
		return errors.New(s.Message())
	}

	if err := c.crypto.signer.VerifyStorageAckSignature(c.conf.SharedKeys.Robot.PublicKey, res); err != nil {
		return err
	}

	//TODO: Verify res.StorageHash

	return nil
}

func (c externalClient) storeID(client api.ExternalClient, id *api.ID, end *api.Endorsement) error {
	req := &api.IDStorageRequest{
		Data:        id,
		Endorsement: end,
	}
	if err := c.crypto.signer.SignIDStorageRequestSignature(req, c.conf.SharedKeys.Robot.PrivateKey); err != nil {
		return err
	}

	res, err := client.StoreID(context.Background(), req)
	if err != nil {
		s, _ := status.FromError(err)
		return errors.New(s.Message())
	}

	if err := c.crypto.signer.VerifyStorageAckSignature(c.conf.SharedKeys.Robot.PublicKey, res); err != nil {
		return err
	}

	//TODO: Verify res.StorageHash

	return nil
}

func (c externalClient) storeContract(client api.ExternalClient, contract *api.Contract, end *api.Endorsement) error {
	req := &api.ContractStorageRequest{
		Contract:    contract,
		Endorsement: end,
	}
	if err := c.crypto.signer.SignContractStorageRequestSignature(req, c.conf.SharedKeys.Robot.PrivateKey); err != nil {
		return err
	}

	res, err := client.StoreContract(context.Background(), req)
	if err != nil {
		s, _ := status.FromError(err)
		return errors.New(s.Message())
	}

	if err := c.crypto.signer.VerifyStorageAckSignature(c.conf.SharedKeys.Robot.PublicKey, res); err != nil {
		return err
	}

	//TODO: Verify res.StorageHash

	return nil
}

func (c externalClient) storeContractMessage(client api.ExternalClient, contractMsg *api.ContractMessage, end *api.Endorsement) error {
	req := &api.ContractMessageStorageRequest{
		ContractMessage: contractMsg,
		Endorsement:     end,
	}
	if err := c.crypto.signer.SignContractMessageStorageRequestSignature(req, c.conf.SharedKeys.Robot.PrivateKey); err != nil {
		return err
	}

	res, err := client.StoreContractMessage(context.Background(), req)
	if err != nil {
		s, _ := status.FromError(err)
		return errors.New(s.Message())
	}

	if err := c.crypto.signer.VerifyStorageAckSignature(c.conf.SharedKeys.Robot.PublicKey, res); err != nil {
		return err
	}

	//TODO: Verify res.StorageHash

	return nil
}

func (c externalClient) GetTransactionStatus(ip string, addr string, txHash string) (mining.TransactionStatus, error) {

	req := &api.TransactionStatusRequest{
		Address: addr,
		Hash:    txHash,
	}

	serverAddr := fmt.Sprintf("%s:%d", ip, c.conf.Services.Datamining.ExternalPort)
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	defer conn.Close()

	if err != nil {
		return mining.TransactionFailure, err
	}

	client := api.NewExternalClient(conn)

	res, err := client.GetTransactionStatus(context.Background(), req)
	if err != nil {
		s, _ := status.FromError(err)
		return -1, errors.New(s.Message())
	}

	return mining.TransactionStatus(res.Status), nil
}

func (c externalClient) LeadContractMining(ip string, txHash string, contract *api.Contract, validators []string) error {
	serverAddr := fmt.Sprintf("%s:%d", ip, c.conf.Services.Datamining.ExternalPort)
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	defer conn.Close()

	if err != nil {
		return err
	}

	client := api.NewExternalClient(conn)

	req := &api.ContractLeadRequest{
		Contract:         contract,
		TransactionHash:  txHash,
		ValidatorPeerIPs: validators,
	}
	if err := c.crypto.signer.SignContractLeadRequest(req, c.conf.SharedKeys.Robot.PrivateKey); err != nil {
		return err
	}
	_, err = client.LeadContractMining(context.Background(), req)
	if err != nil {
		s, _ := status.FromError(err)
		return errors.New(s.Message())
	}
	return nil
}

func (c externalClient) LeadContractMessageMining(ip string, txHash string, message *api.ContractMessage, validators []string) error {
	serverAddr := fmt.Sprintf("%s:%d", ip, c.conf.Services.Datamining.ExternalPort)
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	defer conn.Close()

	if err != nil {
		return err
	}

	client := api.NewExternalClient(conn)

	req := &api.ContractMessageLeadRequest{
		ContractMessage:  message,
		TransactionHash:  txHash,
		ValidatorPeerIPs: validators,
	}
	if err := c.crypto.signer.SignContractMessageLeadRequest(req, c.conf.SharedKeys.Robot.PrivateKey); err != nil {
		return err
	}
	_, err = client.LeadContractMessageMining(context.Background(), req)
	if err != nil {
		s, _ := status.FromError(err)
		return errors.New(s.Message())
	}
	return nil
}

func (c externalClient) GetContractState(ip string, contractAddress string) (*api.ContractStateResponse, error) {
	serverAddr := fmt.Sprintf("%s:%d", ip, c.conf.Services.Datamining.ExternalPort)
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	defer conn.Close()

	if err != nil {
		return nil, err
	}

	client := api.NewExternalClient(conn)

	state, err := client.GetContractState(context.Background(), &api.ContractStateRequest{
		ContractAddress: contractAddress,
	})
	if err != nil {
		return nil, err
	}

	if err := c.crypto.signer.VerifyContractStateSignature(c.conf.SharedKeys.Robot.PublicKey, state); err != nil {
		return nil, err
	}

	return state, nil
}

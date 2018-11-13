package rpc

import (
	"errors"
	"fmt"

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
	LeadKeychainMining(ip string, txHash string, encData string, sig *api.Signature, validators []string) error

	//LeadBiometricMining process the biometric mining as master robot
	LeadBiometricMining(ip string, txHash string, encData string, sig *api.Signature, validators []string) error

	//RequestBiometric requests a peer to retrive a biometric data from a given encrypted person hash
	RequestBiometric(ip string, encPersonHash string) (account.Biometric, error)

	//RequestKeychain requests a peer to retrieve the last keychain from a given encrypted address
	RequestKeychain(ip string, encAddress string) (account.Keychain, error)

	//RequestLock requests a peer to lock a given transaction
	RequestLock(ip string, txLock lock.TransactionLock) error

	//RequestLock requests a peer to unlock a given transaction
	RequestUnlock(ip string, txLock lock.TransactionLock) error

	//RequestValidations requests a peer to process validation/mining as a slave robot
	RequestValidation(ip string, txType mining.TransactionType, txHash string, data interface{}) (mining.Validation, error)

	//RequestStorage requests a peer to store the transaction
	RequestStorage(ip string, txType mining.TransactionType, data interface{}, end mining.Endorsement) error
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

func (c externalClient) LeadKeychainMining(ip string, txHash string, encData string, sig *api.Signature, validators []string) error {
	serverAddr := fmt.Sprintf("%s:%d", ip, c.conf.Datamining.ExternalPort)
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	defer conn.Close()

	if err != nil {
		return err
	}

	client := api.NewExternalClient(conn)

	req := &api.KeychainLeadRequest{
		EncryptedKeychainData: encData,
		SignatureKeychainData: sig,
		TransactionHash:       txHash,
		ValidatorPeerIPs:      validators,
	}
	if err := c.crypto.signer.SignKeychainLeadRequest(req, c.conf.SharedKeys.RobotPrivateKey); err != nil {
		return err
	}

	_, err = client.LeadKeychainMining(context.Background(), req)
	if err != nil {
		s, _ := status.FromError(err)
		return errors.New(s.Message())
	}
	return nil
}

func (c externalClient) LeadBiometricMining(ip string, txHash string, encData string, sig *api.Signature, validators []string) error {
	serverAddr := fmt.Sprintf("%s:%d", ip, c.conf.Datamining.ExternalPort)
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	defer conn.Close()

	if err != nil {
		return err
	}

	client := api.NewExternalClient(conn)

	req := &api.BiometricLeadRequest{
		EncryptedBioData: encData,
		SignatureBioData: sig,
		TransactionHash:  txHash,
		ValidatorPeerIPs: validators,
	}
	if err := c.crypto.signer.SignBiometricLeadRequest(req, c.conf.SharedKeys.RobotPrivateKey); err != nil {
		return err
	}
	_, err = client.LeadBiometricMining(context.Background(), req)
	if err != nil {
		s, _ := status.FromError(err)
		return errors.New(s.Message())
	}
	return nil
}

func (c externalClient) RequestBiometric(ip string, encPersonHash string) (account.Biometric, error) {
	serverAddr := fmt.Sprintf("%s:%d", ip, c.conf.Datamining.ExternalPort)
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	defer conn.Close()

	if err != nil {
		return nil, err
	}

	client := api.NewExternalClient(conn)

	sigReq, err := c.crypto.signer.SignHash(encPersonHash, c.conf.SharedKeys.RobotPrivateKey)
	if err != nil {
		return nil, err
	}

	res, err := client.GetBiometric(context.Background(), &api.BiometricRequest{
		EncryptedPersonHash: encPersonHash,
		Signature:           sigReq,
	})
	if err != nil {
		s, _ := status.FromError(err)
		return nil, errors.New(s.Message())
	}

	if err := c.crypto.signer.VerifyBiometricResponseSignature(c.conf.SharedKeys.RobotPublicKey, res); err != nil {
		return nil, err
	}

	return account.NewBiometric(
		c.data.buildBiometricData(res.Data),
		c.data.buildEndorsement(res.Endorsement)), nil
}

func (c externalClient) RequestKeychain(ip string, encAddress string) (account.Keychain, error) {
	serverAddr := fmt.Sprintf("%s:%d", ip, c.conf.Datamining.ExternalPort)
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	defer conn.Close()

	if err != nil {
		return nil, err
	}

	client := api.NewExternalClient(conn)

	sigReq, err := c.crypto.signer.SignHash(encAddress, c.conf.SharedKeys.RobotPrivateKey)
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

	clearaddr, err := c.crypto.decrypter.DecryptHash(res.Data.CipherAddrRobot, c.conf.SharedKeys.RobotPrivateKey)
	if err != nil {
		return nil, err
	}

	if err := c.crypto.signer.VerifyKeychainResponseSignature(c.conf.SharedKeys.RobotPublicKey, res); err != nil {
		return nil, err
	}

	return account.NewKeychain(
		clearaddr,
		c.data.buildKeychainData(res.Data),
		c.data.buildEndorsement(res.Endorsement),
	), nil
}

func (c externalClient) RequestLock(ip string, txLock lock.TransactionLock) error {
	serverAddr := fmt.Sprintf("%s:%d", ip, c.conf.Datamining.ExternalPort)
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
	if err := c.crypto.signer.SignLockRequest(lockReq, c.conf.SharedKeys.RobotPrivateKey); err != nil {
		return err
	}

	res, err := client.LockTransaction(context.Background(), lockReq)
	if err != nil {
		s, _ := status.FromError(err)
		return errors.New(s.Message())
	}

	if err := c.crypto.signer.VerifyLockAckSignature(c.conf.SharedKeys.RobotPublicKey, res); err != nil {
		return err
	}

	//TODO: Verify res.LockHash

	return nil
}

func (c externalClient) RequestUnlock(ip string, txLock lock.TransactionLock) error {
	serverAddr := fmt.Sprintf("%s:%d", ip, c.conf.Datamining.ExternalPort)
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
	if err := c.crypto.signer.SignLockRequest(lockReq, c.conf.SharedKeys.RobotPrivateKey); err != nil {
		return err
	}

	res, err := client.UnlockTransaction(context.Background(), lockReq)
	if err != nil {
		s, _ := status.FromError(err)
		return errors.New(s.Message())
	}

	if err := c.crypto.signer.VerifyLockAckSignature(c.conf.SharedKeys.RobotPublicKey, res); err != nil {
		return err
	}

	//TODO: Verify res.LockHash

	return nil
}

func (c externalClient) RequestValidation(ip string, txType mining.TransactionType, txHash string, data interface{}) (mining.Validation, error) {
	serverAddr := fmt.Sprintf("%s:%d", ip, c.conf.Datamining.ExternalPort)
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	defer conn.Close()

	if err != nil {
		return nil, err
	}

	client := api.NewExternalClient(conn)

	switch txType {
	case mining.KeychainTransaction:
		return c.validateKeychain(client, txHash, c.api.buildKeychainData(data.(account.KeychainData)))
	case mining.BiometricTransaction:
		return c.validateBiometric(client, txHash, c.api.buildBiometricData(data.(account.BiometricData)))
	}

	return nil, errors.New("Unsupported transaction type")
}

func (c externalClient) RequestStorage(ip string, txType mining.TransactionType, data interface{}, end mining.Endorsement) error {
	serverAddr := fmt.Sprintf("%s:%d", ip, c.conf.Datamining.ExternalPort)
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	defer conn.Close()

	if err != nil {
		return err
	}

	client := api.NewExternalClient(conn)

	endorsement := c.api.buildEndorsement(end)

	switch txType {
	case mining.KeychainTransaction:
		return c.storeKeychain(client, c.api.buildKeychainData(data.(account.KeychainData)), endorsement)
	case mining.BiometricTransaction:
		return c.storeBiometric(client, c.api.buildBiometricData(data.(account.BiometricData)), endorsement)
	}

	return nil
}

func (c externalClient) validateKeychain(client api.ExternalClient, txHash string, data *api.KeychainData) (mining.Validation, error) {
	req := &api.KeychainValidationRequest{
		Data:            data,
		TransactionHash: txHash,
	}
	if err := c.crypto.signer.SignKeychainValidationRequestSignature(req, c.conf.SharedKeys.RobotPrivateKey); err != nil {
		return nil, err
	}

	res, err := client.ValidateKeychain(context.Background(), req)
	if err != nil {
		s, _ := status.FromError(err)
		return nil, errors.New(s.Message())
	}

	if err := c.crypto.signer.VerifyValidationResponseSignature(c.conf.SharedKeys.RobotPublicKey, res); err != nil {
		return nil, err
	}

	return c.data.buildValidation(res.Validation), nil
}

func (c externalClient) validateBiometric(client api.ExternalClient, txHash string, data *api.BiometricData) (mining.Validation, error) {
	req := &api.BiometricValidationRequest{
		Data:            data,
		TransactionHash: txHash,
	}
	if err := c.crypto.signer.SignBiometricValidationRequestSignature(req, c.conf.SharedKeys.RobotPrivateKey); err != nil {
		return nil, err
	}
	res, err := client.ValidateBiometric(context.Background(), req)
	if err != nil {
		s, _ := status.FromError(err)
		return nil, errors.New(s.Message())
	}

	if err := c.crypto.signer.VerifyValidationResponseSignature(c.conf.SharedKeys.RobotPublicKey, res); err != nil {
		return nil, err
	}

	return c.data.buildValidation(res.Validation), nil
}

func (c externalClient) storeKeychain(client api.ExternalClient, data *api.KeychainData, end *api.Endorsement) error {
	req := &api.KeychainStorageRequest{
		Data:        data,
		Endorsement: end,
	}
	if err := c.crypto.signer.SignKeychainStorageRequestSignature(req, c.conf.SharedKeys.RobotPrivateKey); err != nil {
		return err
	}

	res, err := client.StoreKeychain(context.Background(), req)
	if err != nil {
		s, _ := status.FromError(err)
		return errors.New(s.Message())
	}

	if err := c.crypto.signer.VerifyStorageAckSignature(c.conf.SharedKeys.RobotPublicKey, res); err != nil {
		return err
	}

	//TODO: Verify res.StorageHash

	return nil
}

func (c externalClient) storeBiometric(client api.ExternalClient, data *api.BiometricData, end *api.Endorsement) error {
	req := &api.BiometricStorageRequest{
		Data:        data,
		Endorsement: end,
	}
	if err := c.crypto.signer.SignBiometricStorageRequestSignature(req, c.conf.SharedKeys.RobotPrivateKey); err != nil {
		return err
	}

	res, err := client.StoreBiometric(context.Background(), req)
	if err != nil {
		return err
	}

<<<<<<< HEAD
<<<<<<< HEAD
	if err := c.crypto.signer.VerifyStorageAckSignature(c.conf.SharedKeys.RobotPublicKey, res); err != nil {
		return err
=======
	if err := c.crypto.signer.CheckStorageAckSignature(c.conf.SharedKeys.RobotPublicKey, res); err != nil {
		s, _ := status.FromError(err)
		return errors.New(s.Message())
>>>>>>> Use goroutine to dispatch on pool the requests
=======
	if err := c.crypto.signer.CheckStorageAckSignature(c.conf.SharedKeys.RobotPublicKey, res); err != nil {
		s, _ := status.FromError(err)
		return errors.New(s.Message())
>>>>>>> b917fa18c00850ba6fc404e8506c6778b962c01c
	}

	//TODO: Verify res.StorageHash

	return nil
}

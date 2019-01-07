package listing

import (
	"errors"
)

//ErrAccountNotExist is returned when the requested account not exist
var ErrAccountNotExist = errors.New("Account does not exist")

//ErrUnauthorized is returned when the emitter is not authorized
var ErrUnauthorized = errors.New("Unauthorized")

//RobotClient define methods to interfact with the robot
type RobotClient interface {

	//GetAccount asks the datamining service to get the account based on the encrypted ID hash
	GetAccount(encHash string) (AccountResult, error)

	//GetSharedKeys asks the datamining service to get the latest shared emitter keys
	GetSharedKeys() (SharedKeys, error)

	//IsEmitterAuthorized asks the datamining service if the public key is related to an authorized emitter
	IsEmitterAuthorized(emPubKey string) error

	//GetTransactionStatus asks the datamining service to get the transaction status
	GetTransactionStatus(addr string, txHash string) (string, error)
}

//SignatureVerifier defines methods to handle signature verification
type SignatureVerifier interface {

	//VerifyHashSignature checks the hash signature
	VerifyHashSignature(hashedData string, key string, sig string) error

	//VerifyAccountResultSignature checks the account result signature
	VerifyAccountResultSignature(res AccountResult, pubKey string) error
}

//Service define methods for the listing feature
type Service interface {

	//GetSafeSharedKeys gets the latest shared keys (used internally)
	GetSafeSharedKeys() (SharedKeys, error)

	//GetSharedKeys gets the latest shared keys
	GetSharedKeys(emPubKey string, sig string) (SharedKeys, error)

	//ExistAccount checks if an account is related to an encrypted ID hash
	ExistAccount(encryptedIDHash string, sig string) error

	//GetAccount gets an account related to the encrypted ID hash
	GetAccount(encryptedIDHash string, sig string) (AccountResult, error)

	//GetTransactionStatus gets the transaction status
	GetTransactionStatus(addr, txHash string) (string, error)
}

type service struct {
	client RobotClient
	sig    SignatureVerifier
}

//NewService creates a new listing service
func NewService(client RobotClient, sig SignatureVerifier) Service {
	return service{
		client: client,
		sig:    sig,
	}
}

func (s service) ExistAccount(encryptedIDHash string, sig string) error {

	_, err := s.GetAccount(encryptedIDHash, sig)
	if err != nil {
		return err
	}
	return nil
}

func (s service) GetSharedKeys(emPubKey string, sig string) (SharedKeys, error) {
	if err := s.sig.VerifyHashSignature(emPubKey, emPubKey, sig); err != nil {
		return nil, err
	}

	if err := s.client.IsEmitterAuthorized(emPubKey); err != nil {
		return nil, err
	}

	keys, err := s.client.GetSharedKeys()
	if err != nil {
		return nil, err
	}

	return keys, nil
}

func (s service) GetSafeSharedKeys() (SharedKeys, error) {
	keys, err := s.client.GetSharedKeys()
	if err != nil {
		return nil, err
	}

	return keys, nil
}

func (s service) GetAccount(encryptedIDHash string, sig string) (AccountResult, error) {
	keys, err := s.client.GetSharedKeys()
	if err != nil {
		return nil, err
	}

	if err := s.sig.VerifyHashSignature(encryptedIDHash, keys.RequestPublicKey(), sig); err != nil {
		return nil, err
	}

	res, err := s.client.GetAccount(encryptedIDHash)
	if err != nil {
		return nil, err
	}

	if err := s.sig.VerifyAccountResultSignature(res, keys.RobotPublicKey()); err != nil {
		return nil, err
	}

	return res, nil
}

func (s service) GetTransactionStatus(addr string, txHash string) (string, error) {
	return s.client.GetTransactionStatus(addr, txHash)
}

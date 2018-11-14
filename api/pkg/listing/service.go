package listing

import (
	"errors"
)

//ErrAccountNotExist is returned when the requested account not exist
var ErrAccountNotExist = errors.New("Account doest not exist")

//RobotClient define methods to interfact with the robot
type RobotClient interface {

	//GetAccount dispatches the account search request to the Datamining service
	GetAccount(encHash string) (*AccountResult, error)
}

//SignatureVerifier defines methods to verify signature requests
type SignatureVerifier interface {

	//VerifyHashSignatures verifies the signature of the hash
	VerifyHashSignature(hashedData string, key string, sig string) error
}

//Service define methods for the listing feature
type Service interface {

	//ExistAccount checks it the encrypted person hash matches any account
	ExistAccount(encryptedHash string, sig string) error

	//GetAccount retrieves an account from the encrypted person hash
	GetAccount(encryptedHash string, sig string) (*AccountResult, error)
}

type service struct {
	client       RobotClient
	sigVerif     SignatureVerifier
	sharedBioPub string
}

//NewService creates a new listing service
func NewService(sharedBioPub string, client RobotClient, sigVerif SignatureVerifier) Service {
	return service{
		sharedBioPub: sharedBioPub,
		client:       client,
		sigVerif:     sigVerif,
	}
}

func (s service) ExistAccount(encryptedHash string, sig string) error {
	_, err := s.GetAccount(encryptedHash, sig)
	if err != nil {
		return err
	}
	return nil
}

func (s service) GetAccount(encryptedHash string, sig string) (*AccountResult, error) {
	if err := s.sigVerif.VerifyHashSignature(encryptedHash, s.sharedBioPub, sig); err != nil {
		return nil, err
	}

	return s.client.GetAccount(encryptedHash)
}

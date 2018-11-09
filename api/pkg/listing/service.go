package listing

import (
	"errors"
)

//ErrAccountNotExist is returned when the requested account not exist
var ErrAccountNotExist = errors.New("Account doest not exist")

//RobotClient define methods to interfact with the robot
type RobotClient interface {
	GetAccount(encHash string) (*AccountResult, error)
}

//SignatureChecker defines methods to validate signature requests
type SignatureChecker interface {
	CheckHashSignature(hashedData string, key string, sig string) error
}

//Service define methods for the listing feature
type Service interface {
	ExistAccount(encryptedHash string, sig string) error
	GetAccount(encryptedHash string, sig string) (*AccountResult, error)
}

type service struct {
	client       RobotClient
	sigChecker   SignatureChecker
	sharedBioPub string
}

//NewService creates a new listing service
func NewService(sharedBioPub string, client RobotClient, sigChecker SignatureChecker) Service {
	return service{
		sharedBioPub: sharedBioPub,
		client:       client,
		sigChecker:   sigChecker,
	}
}

func (s service) ExistAccount(encryptedHash string, sig string) error {
	if err := s.sigChecker.CheckHashSignature(encryptedHash, s.sharedBioPub, sig); err != nil {
		return err
	}

	_, err := s.client.GetAccount(encryptedHash)
	if err != nil {
		return err
	}
	return nil
}

func (s service) GetAccount(encryptedHash string, sig string) (*AccountResult, error) {
	if err := s.sigChecker.CheckHashSignature(encryptedHash, s.sharedBioPub, sig); err != nil {
		return nil, err
	}

	return s.client.GetAccount(encryptedHash)
}

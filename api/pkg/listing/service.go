package listing

import (
	"errors"
)

//ErrInvalidSignature is returned when the request contains invalid signatures
var ErrInvalidSignature = errors.New("Invalid signature")

//ErrAccountNotExist is returned when the requested account not exist
var ErrAccountNotExist = errors.New("Account doest not exist")

//RobotClient define methods to interfact with the robot
type RobotClient interface {
	GetAccount(encHash string) (*SignedAccountResult, error)
	GetMasterPeer() (MasterPeer, error)
}

//RequestValidator defines methods to validate requests
type RequestValidator interface {
	CheckRawSignature(hashedData string, key string, sig string) (bool, error)
}

//Service define methods for the listing feature
type Service interface {
	ExistAccount(encryptedHash string, sig string) error
	GetAccount(encryptedHash string, sig string) (*SignedAccountResult, error)
	GetMasterPeer() (MasterPeer, error)
}

type service struct {
	client       RobotClient
	val          RequestValidator
	sharedBioPub string
}

//NewService creates a new listing service
func NewService(sharedBioPub string, client RobotClient, val RequestValidator) Service {
	return service{
		sharedBioPub: sharedBioPub,
		client:       client,
		val:          val,
	}
}

func (s service) ExistAccount(encryptedHash string, sig string) error {
	valid, err := s.val.CheckRawSignature(encryptedHash, s.sharedBioPub, sig)
	if err != nil {
		return err
	}
	if !valid {
		return ErrInvalidSignature
	}

	_, err = s.client.GetAccount(encryptedHash)
	if err != nil {
		return err
	}
	return nil
}

func (s service) GetAccount(encryptedHash string, sig string) (*SignedAccountResult, error) {
	valid, err := s.val.CheckRawSignature(encryptedHash, s.sharedBioPub, sig)
	if err != nil {
		return nil, err
	}

	if !valid {
		return nil, ErrInvalidSignature
	}

	return s.client.GetAccount(encryptedHash)
}

func (s service) GetMasterPeer() (MasterPeer, error) {
	return s.client.GetMasterPeer()
}

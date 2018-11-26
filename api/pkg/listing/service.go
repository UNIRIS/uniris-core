package listing

import (
	"errors"

	"github.com/uniris/uniris-core/api/pkg/system"
)

//ErrAccountNotExist is returned when the requested account not exist
var ErrAccountNotExist = errors.New("Account doest not exist")

//RobotClient define methods to interfact with the robot
type RobotClient interface {
	GetAccount(encHash string) (*AccountResult, error)
}

//SignatureVerifier defines methods to verify signature requests
type SignatureVerifier interface {
	VerifyHashSignature(hashedData string, key string, sig string) error
}

//Service define methods for the listing feature
type Service interface {
	ExistAccount(encryptedHash string, sig string) error
	GetAccount(encryptedHash string, sig string) (*AccountResult, error)
}

type service struct {
	client   RobotClient
	sigVerif SignatureVerifier
	conf     system.UnirisConfig
}

//NewService creates a new listing service
func NewService(conf system.UnirisConfig, client RobotClient, sigVerif SignatureVerifier) Service {
	return service{
		conf:     conf,
		client:   client,
		sigVerif: sigVerif,
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
	verifKey := s.conf.SharedKeys.EmitterRequestKey().PublicKey
	if err := s.sigVerif.VerifyHashSignature(encryptedHash, verifKey, sig); err != nil {
		return nil, err
	}

	return s.client.GetAccount(encryptedHash)
}

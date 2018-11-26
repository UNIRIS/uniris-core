package adding

import (
	"github.com/uniris/uniris-core/api/pkg/system"
)

//Service defines methods to adding to the blockchain
type Service interface {
	AddAccount(AccountCreationRequest) (*AccountCreationResult, error)
}

//RobotClient define methods to interfact with the robot
type RobotClient interface {
	AddAccount(AccountCreationRequest) (*AccountCreationResult, error)
}

//SignatureVerifier defines methods to verify signature requests
type SignatureVerifier interface {
	VerifyAccountCreationRequestSignature(data AccountCreationRequest, key string) error
}

type service struct {
	conf     system.UnirisConfig
	client   RobotClient
	sigVerif SignatureVerifier
}

//NewService creates a new adding service
func NewService(conf system.UnirisConfig, cli RobotClient, sigVerif SignatureVerifier) Service {
	return service{conf, cli, sigVerif}
}

func (s service) AddAccount(req AccountCreationRequest) (*AccountCreationResult, error) {
	verifKey := s.conf.SharedKeys.EmitterRequestKey().PublicKey
	if err := s.sigVerif.VerifyAccountCreationRequestSignature(req, verifKey); err != nil {
		return nil, err
	}

	return s.client.AddAccount(req)
}

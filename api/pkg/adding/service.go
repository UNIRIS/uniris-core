package adding

import (
	"errors"
)

//ErrInvalidSignature is returned when the request contains invalid signatures
var ErrInvalidSignature = errors.New("Invalid signature")

//Service defines methods to adding to the blockchain
type Service interface {
	AddAccount(EnrollmentRequest) (*EnrollmentResult, error)
}

//RobotClient define methods to interfact with the robot
type RobotClient interface {
	AddAccount(EnrollmentRequest) (*EnrollmentResult, error)
}

//RequestValidator defines methods to validate requests
type RequestValidator interface {
	CheckDataSignature(data interface{}, key string, sig string) (bool, error)
}

type service struct {
	sharedBioPub string
	client       RobotClient
	val          RequestValidator
}

//NewService creates a new adding service
func NewService(sharedBioPub string, cli RobotClient, val RequestValidator) Service {
	return service{sharedBioPub, cli, val}
}

func (s service) AddAccount(req EnrollmentRequest) (*EnrollmentResult, error) {
	verifReq := EnrollmentData{
		EncryptedBioData:    req.EncryptedBioData,
		EncryptedWalletData: req.EncryptedWalletData,
		SignaturesBio:       req.SignaturesBio,
		SignaturesWallet:    req.SignaturesWallet,
	}

	valid, err := s.val.CheckDataSignature(verifReq, s.sharedBioPub, req.SignatureRequest)
	if err != nil {
		return nil, err
	}

	if !valid {
		return nil, ErrInvalidSignature
	}

	return s.client.AddAccount(req)
}

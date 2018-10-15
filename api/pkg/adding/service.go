package adding

import "errors"

//ErrInvalidSignature is returned when the request contains invalid signatures
var ErrInvalidSignature = errors.New("Request contains an invalid signature")

//Service defines methods to adding to the blockchain
type Service interface {
	AddAccount(EnrollmentRequest) (EnrollmentResult, error)
}

//RobotClient define methods to interfact with the robot
type RobotClient interface {
	AddAccount(EnrollmentRequest) (EnrollmentResult, error)
}

//RequestValidator defines methods to validate requests
type RequestValidator interface {
	CheckSignature(data interface{}, key []byte, sig []byte) (bool, error)
}

type service struct {
	sharedBioPub []byte
	client       RobotClient
	val          RequestValidator
}

//NewService creates a new adding service
func NewService(sharedBioPub []byte, cli RobotClient, val RequestValidator) Service {
	return service{sharedBioPub, cli, val}
}

func (s service) AddAccount(req EnrollmentRequest) (EnrollmentResult, error) {
	var res EnrollmentResult

	valid, err := s.val.CheckSignature(req, s.sharedBioPub, []byte(req.SignatureRequest))
	if err != nil {
		return res, err
	}

	if !valid {
		return res, ErrInvalidSignature
	}

	return s.client.AddAccount(req)
}

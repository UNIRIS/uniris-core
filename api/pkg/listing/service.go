package listing

import "errors"

//ErrInvalidSignature is returned when the request contains invalid signatures
var ErrInvalidSignature = errors.New("Invalid signature")

//RobotClient define methods to interfact with the robot
type RobotClient interface {
	GetAccount(AccountRequest) (AccountResult, error)
}

//RequestValidator defines methods to validate requests
type RequestValidator interface {
	CheckSignature(data interface{}, key []byte, sig []byte) (bool, error)
}

//Service define methods for the listing feature
type Service interface {
	GetAccount(encryptedHash []byte, sig []byte) (AccountResult, error)
}

type service struct {
	client       RobotClient
	val          RequestValidator
	sharedBioPub []byte
}

//NewService creates a new listing service
func NewService(sharedBioPub []byte, client RobotClient, val RequestValidator) Service {
	return service{
		sharedBioPub: sharedBioPub,
		client:       client,
		val:          val,
	}
}

func (s service) GetAccount(encryptedHash []byte, sig []byte) (AccountResult, error) {
	var res AccountResult

	req := AccountRequest{
		EncryptedHash:    encryptedHash,
		SignatureRequest: sig,
	}

	valid, err := s.val.CheckSignature(req, s.sharedBioPub, sig)
	if err != nil {
		return res, err
	}

	if !valid {
		return res, ErrInvalidSignature
	}

	return s.client.GetAccount(req)
}

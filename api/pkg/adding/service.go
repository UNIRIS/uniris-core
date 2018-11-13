package adding

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
	sharedBioPub string
	client       RobotClient
	sigVerif     SignatureVerifier
}

//NewService creates a new adding service
func NewService(sharedBioPub string, cli RobotClient, sigVerif SignatureVerifier) Service {
	return service{sharedBioPub, cli, sigVerif}
}

func (s service) AddAccount(req AccountCreationRequest) (*AccountCreationResult, error) {
	if err := s.sigVerif.VerifyAccountCreationRequestSignature(req, s.sharedBioPub); err != nil {
		return nil, err
	}

	return s.client.AddAccount(req)
}

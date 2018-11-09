package adding

//Service defines methods to adding to the blockchain
type Service interface {
	AddAccount(AccountCreationRequest) (*AccountCreationResult, error)
}

//RobotClient define methods to interfact with the robot
type RobotClient interface {
	AddAccount(AccountCreationRequest) (*AccountCreationResult, error)
}

//SignatureChecker defines methods to validate signature requests
type SignatureChecker interface {
	CheckAccountSignature(data AccountCreationRequest, key string) error
}

type service struct {
	sharedBioPub string
	client       RobotClient
	sigChecker   SignatureChecker
}

//NewService creates a new adding service
func NewService(sharedBioPub string, cli RobotClient, sigChecker SignatureChecker) Service {
	return service{sharedBioPub, cli, sigChecker}
}

func (s service) AddAccount(req AccountCreationRequest) (*AccountCreationResult, error) {
	if err := s.sigChecker.CheckAccountSignature(req, s.sharedBioPub); err != nil {
		return nil, err
	}

	return s.client.AddAccount(req)
}

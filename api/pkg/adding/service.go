package adding

//Service defines methods to adding to the blockchain
type Service interface {

	//RegisterBiod validates the request and dispatches the biod register request to the Datamining service
	RegisterBiod(BiodRegisterRequest) (*BiodRegisterResponse, error)

	//AddAccount validates the request and dispatch the account creation request to the Datamining service
	AddAccount(AccountCreationRequest) (*AccountCreationResult, error)
}

//RobotClient define methods to interfact with the robot
type RobotClient interface {

	//RegisterBiod dispatches the encrypted biod public key to the Datamining service
	RegisterBiod(encPubKey string) (*BiodRegisterResponse, error)

	//AddAccount dispatches the account creation request to the Datamining service
	AddAccount(AccountCreationRequest) (*AccountCreationResult, error)
}

//SignatureVerifier defines methods to verify signature requests
type SignatureVerifier interface {

	//VerifyBiodRegisteringRequestSignature verifies the signature of the biod register request
	VerifyBiodRegisteringRequestSignature(req BiodRegisterRequest, pubKey string) error

	//VerifyAccountCreationRequestSignature verifies the signature of the account creation request
	VerifyAccountCreationRequestSignature(data AccountCreationRequest, pubKey string) error
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

func (s service) RegisterBiod(req BiodRegisterRequest) (*BiodRegisterResponse, error) {
	if err := s.sigVerif.VerifyBiodRegisteringRequestSignature(req, s.sharedBioPub); err != nil {
		return nil, err
	}
	return s.client.RegisterBiod(req.EncryptedPublicKey)
}

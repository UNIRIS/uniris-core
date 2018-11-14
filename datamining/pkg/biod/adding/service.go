package adding

//Repository defines methods to store data to the tech database
type Repository interface {
	StoreBiodPublicKey(key string) error
}

//Service define methods to handle biometric devices registering
type Service interface {

	//RegisterKey stores the public key into the tech database
	RegisterKey(pubKey string) error
}

type service struct {
	repo Repository
}

//NewService creates a new service handling the registering of biometric device
func NewService(repo Repository) Service {
	return service{repo}
}

func (s service) RegisterKey(pubKey string) error {
	return s.repo.StoreBiodPublicKey(pubKey)
}

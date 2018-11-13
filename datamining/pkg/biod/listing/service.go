package listing

//Repository defines mtehods to get data from the tech database
type Repository interface {
	//ListBiodPubKeys queries the biometric device public keys registered
	ListBiodPubKeys() ([]string, error)
}

//Service define methods for the biometric devices
type Service interface {

	//ListBiodPubKeys find the biometric device public keys registered
	ListBiodPubKeys() ([]string, error)
}

type service struct {
	repo Repository
}

//NewService creates a new service for biometric devices listing
func NewService(repo Repository) Service {
	return service{repo}
}

func (s service) ListBiodPubKeys() ([]string, error) {
	return s.repo.ListBiodPubKeys()
}

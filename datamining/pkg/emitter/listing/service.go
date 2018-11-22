package listing

//Repository defines methods to handle emitters sotrage
type Repository interface {

	//ListEmitterPublicKeys retrieves the emitter public keys
	ListEmitterPublicKeys() ([]string, error)
}

//Service define methods to list emitters
type Service interface {

	//ListEmitterPublicKeys list the emitter public keys registered
	ListEmitterPublicKeys() ([]string, error)
}

type service struct {
	repo Repository
}

//NewService creates a new service for biometric devices listing
func NewService(repo Repository) Service {
	return service{repo}
}

func (s service) ListEmitterPublicKeys() ([]string, error) {
	return s.repo.ListEmitterPublicKeys()
}

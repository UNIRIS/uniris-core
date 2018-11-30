package listing

import "github.com/uniris/uniris-core/datamining/pkg/emitter"

//Repository defines methods to handle emitters sotrage
type Repository interface {

	//ListSharedEmitterKeyPairs retrieves the shared emitter keypair
	ListSharedEmitterKeyPairs() ([]emitter.SharedKeyPair, error)
}

//Service define methods to list emitters
type Service interface {

	//IsEmitterAuthorized checks if the emitter public key is authorized
	IsEmitterAuthorized(pubKey string) error

	//ListSharedEmitterKeyPairs get the shared emitter key pairs
	ListSharedEmitterKeyPairs() ([]emitter.SharedKeyPair, error)
}

type service struct {
	repo Repository
}

//NewService creates a new service for ID devices listing
func NewService(repo Repository) Service {
	return service{repo}
}

func (s service) ListSharedEmitterKeyPairs() ([]emitter.SharedKeyPair, error) {
	return s.repo.ListSharedEmitterKeyPairs()
}

func (s service) IsEmitterAuthorized(pubKey string) error {
	//TODO: request smart contract
	return nil
}

package shared

//KeyRepository handle the shared key storage
type KeyRepository interface {
	ListSharedEmitterKeyPairs() ([]KeyPair, error)
	StoreSharedEmitterKeyPair(kp KeyPair) error
}

//Service handle shared management
type Service struct {
	keyRepo KeyRepository
}

func NewService(keyR KeyRepository) Service {
	return Service{
		keyRepo: keyR,
	}
}

//ListSharedEmitterKeyPairs get the shared emitter key pairs
func (s Service) ListSharedEmitterKeyPairs() ([]KeyPair, error) {
	return s.keyRepo.ListSharedEmitterKeyPairs()
}

//StoreSharedEmitterKeyPair store emitter shared key
func (s Service) StoreSharedEmitterKeyPair(kp KeyPair) error {
	return s.keyRepo.StoreSharedEmitterKeyPair(kp)
}

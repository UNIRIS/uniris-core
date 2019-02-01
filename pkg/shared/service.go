package shared

//KeyRepository handle the shared key storage
type KeyRepository interface {
	ListSharedEmitterKeyPairs() (EmitterKeys, error)
	StoreSharedEmitterKeyPair(EmitterKeyPair) error
	GetLastSharedMinersKeyPair() (MinerKeyPair, error)
}

//Service handle shared management
type Service struct {
	keyRepo KeyRepository
}

//NewService create a new shared service
func NewService(keyR KeyRepository) Service {
	return Service{
		keyRepo: keyR,
	}
}

//IsEmitterKeyAuthorized checks if the emitter public key is authorized
func (s Service) IsEmitterKeyAuthorized(emPubKey string) (bool, error) {
	//TODO: request smart contract
	return true, nil
}

//ListSharedEmitterKeyPairs get the shared emitter key pairs
func (s Service) ListSharedEmitterKeyPairs() (EmitterKeys, error) {
	return s.keyRepo.ListSharedEmitterKeyPairs()
}

//GetSharedMinerKeys gets the shared miners keys
func (s Service) GetSharedMinerKeys() (MinerKeyPair, error) {
	return s.keyRepo.GetLastSharedMinersKeyPair()
}

//StoreSharedEmitterKeyPair store emitter shared key
func (s Service) StoreSharedEmitterKeyPair(kp EmitterKeyPair) error {
	return s.keyRepo.StoreSharedEmitterKeyPair(kp)
}

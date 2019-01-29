package emitter

type Service struct{}

//IsEmitterAuthorized checks if the emitter public key is authorized
func (s Service) IsEmitterAuthorized(emPubKey string) (bool, error) {
	//TODO: request smart contract
	return true, nil
}

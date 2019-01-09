package listing

import (
	uniris "github.com/uniris/uniris-core/pkg"
)

//Repository defines methods to handle listing storage
type Repository interface {

	//ListSharedEmitterKeyPairs retrieves the shared emitter keypair
	ListSharedEmitterKeyPairs() ([]uniris.SharedKeys, error)

	//FindID retrieve a ID from a given hash
	FindID(idHash string) (uniris.ID, error)

	//FindIDByTransaction retrieve an ID from a given transaction hash
	FindIDByTransaction(txHash string) (uniris.ID, error)

	//FindLastKeychain retrieve the last keychain from a given account's address
	FindLastKeychain(addr string) (uniris.Keychain, error)

	//FindKeychain retrieve a keychain from a given account's address and the transaction hash
	FindKeychain(addr, txHash string) (uniris.Keychain, error)
}

//Service handles data retreiving
type Service struct {
	repo Repository
}

//NewService creates a new service to retrieve data
func NewService(repo Repository) Service {
	return Service{repo}
}

//ListSharedEmitterKeyPairs get the shared emitter key pairs
func (s Service) ListSharedEmitterKeyPairs() ([]uniris.SharedKeys, error) {
	return s.repo.ListSharedEmitterKeyPairs()
}

//IsEmitterAuthorized checks if the emitter public key is authorized
func (s Service) IsEmitterAuthorized(emPubKey string) (bool, error) {
	//TODO: request smart contract
	return true, nil
}

//GetLastKeychain retrieve the last keychain from a given account's address
func (s Service) GetLastKeychain(addr string) (kc uniris.Keychain, err error) {

	return
}

//GetKeychainByTransactionHash retrieve the a keychain from a given account's address and a transaction hash
func (s Service) GetKeychainByTransactionHash(txHash string) (kc uniris.Keychain, err error) {
	return
}

//GetID retrieve an ID from a given hash
func (s Service) GetID(idHash string) (id uniris.ID, err error) {
	return
}

//GetIDByTransactionHash retrieve an ID from a given transaction hash
func (s Service) GetIDByTransactionHash(txHash string) (id uniris.ID, err error) {
	return
}

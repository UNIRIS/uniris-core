package shared

//TechDatabaseReader wraps the shared emitter and node storage
type TechDatabaseReader interface {
	EmitterReader
	NodeReader
}

//NodeReader performs queries to retrieve shared nodes information
type NodeReader interface {

	//NodeFirstKeys retrieves the first shared node keys
	NodeFirstKeys() (KeyPair, error)

	//NodeLastKeys retrieve the last shared node keys
	NodeLastKeys() (KeyPair, error)

	//AuthorizedPublicKeys retrieves the list of all the authorized node public keys
	AuthorizedPublicKeys() ([]string, error)
}

//EmitterReader handles queries for the shared emitter information
type EmitterReader interface {

	//EmitterKeys retrieve the shared emitter key from the Tech DB
	EmitterKeys() (EmitterKeys, error)
}

//IsEmitterKeyAuthorized checks if the emitter public key is authorized
func IsEmitterKeyAuthorized(emPubKey string) (bool, error) {
	//TODO: request smart contract

	return true, nil
}

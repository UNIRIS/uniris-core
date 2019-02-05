package shared

//TechDatabaseReader wraps the shared emitter and miner storage
type TechDatabaseReader interface {
	EmitterDatabaseReader

	//LastMinerKeys retrieve the last shared miner keys from the Tech DB
	LastMinerKeys() (MinerKeyPair, error)
}

//EmitterDatabaseReader handles queries for the shared emitter information
type EmitterDatabaseReader interface {
	//EmitterKeys retrieve the shared emitter key from the Tech DB
	EmitterKeys() (EmitterKeys, error)
}

//IsEmitterKeyAuthorized checks if the emitter public key is authorized
func IsEmitterKeyAuthorized(emPubKey string) (bool, error) {
	//TODO: request smart contract

	return true, nil
}

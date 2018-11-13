package account

import datamining "github.com/uniris/uniris-core/datamining/pkg"

//PoolRequester handles account requesting on a dedicated pool
type PoolRequester interface {

	//RequestBiometric ask a storage pool to retrieve a biometric based on the person hash
	RequestBiometric(sPool datamining.Pool, personHash string) (Biometric, error)

	//RequestKeychain asks a storage pool to retrieve keychain based on the account's address
	RequestKeychain(sPool datamining.Pool, addr string) (Keychain, error)
}

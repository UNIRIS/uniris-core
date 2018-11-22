package account

import datamining "github.com/uniris/uniris-core/datamining/pkg"

//PoolRequester handles account requesting on a dedicated pool
type PoolRequester interface {

	//RequestID ask a storage pool to retrieve a ID based on thehash
	RequestID(sPool datamining.Pool, hash string) (EndorsedID, error)

	//RequestKeychain asks a storage pool to retrieve keychain based on the account's address
	RequestKeychain(sPool datamining.Pool, addr string) (EndorsedKeychain, error)
}

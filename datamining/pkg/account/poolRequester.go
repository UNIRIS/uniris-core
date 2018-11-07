package account

import datamining "github.com/uniris/uniris-core/datamining/pkg"

//PoolRequester define methods for data requesting
type PoolRequester interface {
	RequestBiometric(sPool datamining.Pool, personHash string) (Biometric, error)
	RequestKeychain(sPool datamining.Pool, addr string) (Keychain, error)
}

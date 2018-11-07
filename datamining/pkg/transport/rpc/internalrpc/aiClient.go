package internalrpc

import datamining "github.com/uniris/uniris-core/datamining/pkg"

//AIClient define methods to communicate with AI
type AIClient interface {
	GetBiometricStoragePool(personHash string) (datamining.Pool, error)
	GetKeychainStoragePool(address string) (datamining.Pool, error)
	GetMasterPeer(txHash string) (datamining.Peer, error)
	GetValidationPool(txHash string) (datamining.Pool, error)
}

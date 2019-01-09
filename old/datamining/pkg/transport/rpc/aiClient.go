package rpc

import datamining "github.com/uniris/uniris-core/datamining/pkg"

//AIClient define methods to communicate with AI
type AIClient interface {

	//GetStoragePool asks the AI service to perform lookup of a storage pool based on a hash
	//The hash can be an address, or a transaction hash
	GetStoragePool(hash string) (datamining.Pool, error)

	//GetMasterPeer asks the AI service to perform lookup of a elected master peer based on a transaction hash
	GetMasterPeer(txHash string) (datamining.Peer, error)

	//GetValidationPool asks the AI service to perform a search of validation pools based on a transaction hash
	GetValidationPool(txHash string) (datamining.Pool, error)
}

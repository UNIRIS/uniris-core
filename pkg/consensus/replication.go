package consensus

import "github.com/uniris/uniris-core/pkg/chain"

//GetMinimumReplicas returns the minimum number of replication for the transaction hash
func GetMinimumReplicas(txHash string) int {
	//TODO: Implement the algorithm
	return 1
}

//ReplicateTransaction process the transaction replication inside the sharding pool
func ReplicateTransaction() {
	//TODO: Implement the algorithm
}

//IsAuthorizedToStoreTx checks if the transaction can be stored on this peer
func IsAuthorizedToStoreTx(tx chain.Transaction) bool {
	//TODO: Implement the algorithm
	return true
}

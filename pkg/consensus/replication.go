package consensus

import "github.com/uniris/uniris-core/pkg/chain"

//GetMinimumReplicas returns the minimum number of replication for the transaction hash
func GetMinimumReplicas(txHash string) int {
	return 1
}

//ReplicateTransaction process the transaction replication inside the sharding pool
func ReplicateTransaction() {
	//TODO: Implement the algorithm
}

func IsAuthorizedToStoreTx(tx chain.Transaction) bool {
	return true
}

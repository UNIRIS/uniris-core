package electing

import (
	"net"

	uniris "github.com/uniris/uniris-core/pkg"
)

func FindLastValidationPool(address string, req uniris.PoolRequester) (uniris.Pool, error) {
	sPool, err := FindStoragePool(address)
	if err != nil {
		return nil, err
	}

	tx, err := req.RequestLastTransaction(sPool, address)
	if err != nil {
		return nil, err
	}

	//TODO: if last transaction is nil (i.efirst transaction of the chain) which pool to choose ?

	return tx.MasterValidation().PreviousTransactionMiners(), nil
}

//FindValidationPool searches a validation pool from a transaction hash
func FindValidationPool(txHash string) (uniris.Pool, error) {
	//TODO: Implements AI lookups to identify the right validation pool

	peers := make([]uniris.PeerIdentity, 0)
	peers[0] = uniris.NewPeerIdentity(net.ParseIP("127.0.0.1"), 0, "")

	return peers, nil
}

//FindStoragePool searches a storage pool for the given address
func FindStoragePool(address string) (uniris.Pool, error) {
	//TODO: Implements AI lookups to identify the right storage pool
	peers := make([]uniris.PeerIdentity, 0)
	peers[0] = uniris.NewPeerIdentity(net.ParseIP("127.0.0.1"), 0, "")

	return peers, nil
}

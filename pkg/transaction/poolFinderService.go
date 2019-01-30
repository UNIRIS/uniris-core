package transaction

import "net"

//PoolRetriever define methods to get transaction from a pool
type PoolRetriever interface {

	//RequestLastTransaction asks a pool to retrieve the last transaction of an address
	RequestLastTransaction(pool Pool, txAddr string, txType Type) (*Transaction, error)
}

//PoolFindingService handles the pool finding
type PoolFindingService struct {
	pRetr PoolRetriever
}

//NewPoolFindingService creates a new transaction pool finding service
func NewPoolFindingService(pR PoolRetriever) PoolFindingService {
	return PoolFindingService{
		pRetr: pR,
	}
}

//RequestLastTransaction find the storage pool address and request the last transaction
func (s PoolFindingService) RequestLastTransaction(txAddr string, txType Type) (*Transaction, error) {
	sPool, err := s.FindStoragePool(txAddr)
	if err != nil {
		return nil, err
	}

	tx, err := s.pRetr.RequestLastTransaction(sPool, txAddr, txType)
	return tx, nil
}

//FindLastValidationPool retrieves the last validation pool for a given address
//TODO: if last transaction is nil (i.efirst transaction of the chain) which pool to choose ?
func (s PoolFindingService) FindLastValidationPool(txAddr string, txType Type) (Pool, error) {
	tx, err := s.RequestLastTransaction(txAddr, txType)
	if err != nil {
		return nil, err
	}
	if tx == nil {
		return nil, nil
	}

	return tx.MasterValidation().PreviousTransactionMiners(), nil
}

//FindValidationPool searches a validation pool from a transaction hash
//TODO: Implements AI lookups to identify the right validation pool
func (s PoolFindingService) FindValidationPool(txHash string) (Pool, error) {
	return Pool{
		PoolMember{
			ip:   net.ParseIP("127.0.0.1"),
			port: 3545,
		},
	}, nil
}

//FindStoragePool searches a storage pool for the given address
//TODO: Implements AI lookups to identify the right storage pool
func (s PoolFindingService) FindStoragePool(address string) (Pool, error) {
	return Pool{
		PoolMember{
			ip:   net.ParseIP("127.0.0.1"),
			port: 3545,
		},
	}, nil
}

//FindTransactionMasterPeer finds a master peer from a transaction hash
//TODO: To implement with AI algorithms
func (s PoolFindingService) FindTransactionMasterPeer(txHash string) (string, int) {
	return "127.0.0.1", 3545
}

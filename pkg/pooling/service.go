package pooling

//Service handle pool generation
type Service struct {
	req PoolRequester
}

//Pool represents a group of peers
type Pool struct {
	peers []string
}

//Peers return the list of peers inside a pool
func (p Pool) Peers() []string {
	return p.peers
}

func NewService(req PoolRequester) Service {
	return Service{req}
}

//FindStoragePool searches a storage pool for the given address
func (s Service) FindStoragePool(address string) (Pool, error) {
	//TODO: Implements AI lookups to identify the right storage pool
	return Pool{
		peers: []string{
			"127.0.0.1",
		},
	}, nil
}

//FindLastValidationPool searches a validation pool for the given address
func (s Service) FindLastValidationPool(address string) (p Pool, err error) {
	sPool, err := s.FindStoragePool(address)
	if err != nil {
		return
	}

	tx, err := s.req.RequestLastTransaction(sPool, address)
	if err != nil {
		return
	}

	lastPool := Pool{
		peers: make([]string, 0),
	}
	for _, miner := range tx.MasterValidation().PreviousTransactionMiners() {
		lastPool.peers = append(lastPool.peers, miner)
	}

	return lastPool, nil
}

//FindValidationPool searches a validation pool from a transaction hash
func (s Service) FindValidationPool(txHash string) (Pool, error) {
	//TODO: Implements AI lookups to identify the right validation pool
	return Pool{
		peers: []string{
			"127.0.0.1",
		},
	}, nil
}

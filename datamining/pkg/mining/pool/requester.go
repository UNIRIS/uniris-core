package pool

//Finder defines methods to find miners to perform validation
type Finder interface {
	FindValidationPool() (PeerCluster, error)
	FindStoragePool() (PeerCluster, error)
	FindLastValidationPool(addr string) (PeerCluster, error)
}

func GetPools(addr string, f Finder) (lastVPool PeerCluster, vPool PeerCluster, sPool PeerCluster, err error) {
	sPool, err = f.FindStoragePool()
	if err != nil {
		return
	}

	lastVPool, err = f.FindLastValidationPool(addr)
	if err != nil {
		return
	}

	vPool, err = f.FindValidationPool()
	if err != nil {
		return
	}

	return
}

package pool

//Finder defines methods to find miners to perform validation
type Finder interface {
	FindLastValidationPool(addr string) (PeerGroup, error)
	FindValidationPool() (PeerGroup, error)
	FindStoragePool() (PeerGroup, error)
}

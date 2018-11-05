package mock

import (
	"net"

	"github.com/uniris/uniris-core/datamining/pkg/mining"
)

//NewPoolFinder creates a new pool finder
func NewPoolFinder() mining.PoolFinder {
	return poolFinder{}
}

type poolFinder struct{}

func (p poolFinder) FindLastValidationPool(addr string) (mining.Pool, error) {
	return mining.NewPool(mining.Peer{
		IP:        net.ParseIP("127.0.0.1"),
		PublicKey: "key",
	}), nil
}

func (p poolFinder) FindValidationPool() (mining.Pool, error) {
	return mining.NewPool(mining.Peer{
		IP:        net.ParseIP("127.0.0.1"),
		PublicKey: "key",
	}), nil
}

func (p poolFinder) FindStoragePool() (mining.Pool, error) {
	return mining.NewPool(mining.Peer{
		IP:        net.ParseIP("127.0.0.1"),
		PublicKey: "key",
	}), nil
}

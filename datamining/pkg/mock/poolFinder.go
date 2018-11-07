package mock

import (
	"net"

	datamining "github.com/uniris/uniris-core/datamining/pkg"
	"github.com/uniris/uniris-core/datamining/pkg/mining"
)

//NewPoolFinder creates a new pool finder
func NewPoolFinder() mining.PoolFinder {
	return poolFinder{}
}

type poolFinder struct{}

func (p poolFinder) FindLastValidationPool(addr string) (datamining.Pool, error) {
	return datamining.NewPool(datamining.Peer{
		IP:        net.ParseIP("127.0.0.1"),
		PublicKey: "key",
	}), nil
}

func (p poolFinder) FindValidationPool() (datamining.Pool, error) {
	return datamining.NewPool(datamining.Peer{
		IP:        net.ParseIP("127.0.0.1"),
		PublicKey: "key",
	}), nil
}

func (p poolFinder) FindStoragePool() (datamining.Pool, error) {
	return datamining.NewPool(datamining.Peer{
		IP:        net.ParseIP("127.0.0.1"),
		PublicKey: "key",
	}), nil
}

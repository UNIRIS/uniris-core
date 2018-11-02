package mock

import (
	"net"

	datamining "github.com/uniris/uniris-core/datamining/pkg"
)

//NewPoolFinder creates a new pool finder
func NewPoolFinder() master.PoolFinder {
	return poolFinder{}
}

type poolFinder struct{}

func (p poolFinder) FindLastValidationPool(addr string) (datamining.Pool, error) {
	return datamining.Pool{
		Peers: []datamining.Peer{
			datamining.Peer{
				IP:        net.ParseIP("127.0.0.1"),
				PublicKey: "key",
			},
		},
	}, nil
}

func (p poolFinder) FindValidationPool() (datamining.Pool, error) {
	return datamining.Pool{
		Peers: []datamining.Peer{
			datamining.Peer{
				IP:        net.ParseIP("127.0.0.1"),
				PublicKey: "key",
			},
		},
	}, nil
}

func (p poolFinder) FindStoragePool() (datamining.Pool, error) {
	return datamining.Pool{
		Peers: []datamining.Peer{
			datamining.Peer{
				IP:        net.ParseIP("127.0.0.1"),
				PublicKey: "key",
			},
		},
	}, nil
}

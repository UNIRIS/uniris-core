package mock

import (
	"net"

	"github.com/uniris/uniris-core/datamining/pkg/leading"
)

//NewPoolFinder creates a new pool finder
func NewPoolFinder() leading.PoolFinder {
	return poolFinder{}
}

type poolFinder struct{}

func (p poolFinder) FindLastValidationPool(addr string) (leading.Pool, error) {
	return leading.Pool{
		Peers: []leading.Peer{
			leading.Peer{
				IP:        net.ParseIP("127.0.0.1"),
				PublicKey: "key",
			},
		},
	}, nil
}

func (p poolFinder) FindValidationPool() (leading.Pool, error) {
	return leading.Pool{
		Peers: []leading.Peer{
			leading.Peer{
				IP:        net.ParseIP("127.0.0.1"),
				PublicKey: "key",
			},
		},
	}, nil
}

func (p poolFinder) FindStoragePool() (leading.Pool, error) {
	return leading.Pool{
		Peers: []leading.Peer{
			leading.Peer{
				IP:        net.ParseIP("127.0.0.1"),
				PublicKey: "key",
			},
		},
	}, nil
}

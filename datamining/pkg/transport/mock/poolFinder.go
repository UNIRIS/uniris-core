package mock

import (
	"net"

	"github.com/uniris/uniris-core/datamining/pkg/mining/master/pool"
)

//NewPoolFinder creates a new pool finder
func NewPoolFinder() pool.Finder {
	return poolFinder{}
}

type poolFinder struct{}

func (p poolFinder) FindLastValidationPool(addr string) (pool.PeerGroup, error) {
	return pool.PeerGroup{
		Peers: []pool.Peer{
			pool.Peer{
				IP:        net.ParseIP("127.0.0.1"),
				PublicKey: "key",
			},
		},
	}, nil
}

func (p poolFinder) FindValidationPool() (pool.PeerGroup, error) {
	return pool.PeerGroup{
		Peers: []pool.Peer{
			pool.Peer{
				IP:        net.ParseIP("127.0.0.1"),
				PublicKey: "key",
			},
		},
	}, nil
}

func (p poolFinder) FindStoragePool() (pool.PeerGroup, error) {
	return pool.PeerGroup{
		Peers: []pool.Peer{
			pool.Peer{
				IP:        net.ParseIP("127.0.0.1"),
				PublicKey: "key",
			},
		},
	}, nil
}

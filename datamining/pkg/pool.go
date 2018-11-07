package datamining

import (
	"net"
)

//Pool represents a pool of peers
type Pool interface {
	Peers() []Peer
}

type pool struct {
	peers []Peer
}

func (p pool) Peers() []Peer {
	return p.peers
}

//NewPool creates a new pool
func NewPool(pp ...Peer) Pool {
	return pool{pp}
}

//Peer represents a peer in a pool
type Peer struct {
	IP        net.IP
	PublicKey string
}

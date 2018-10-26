package pool

import (
	"net"
)

type PeerCluster struct {
	Peers []Peer
}

//Peer represents a peer in a pool
type Peer struct {
	IP        net.IP
	PublicKey string
}

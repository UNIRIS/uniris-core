package pool

import "net"

//PeerGroup represents a group of peers in a pool
type PeerGroup struct {
	Peers []Peer
}

//Peer represents a peer in a pool
type Peer struct {
	IP        net.IP
	PublicKey string
}

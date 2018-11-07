package datamining

import (
	"net"
)

//Pool represents a pool of peers
type Pool interface {
	Peers() PeerList
}

type pool struct {
	peers PeerList
}

//PeerList define a list of peer for a pool
type PeerList []Peer

func (p pool) Peers() PeerList {
	return p.peers
}

//IPs returns the IP of the peer list
func (pL PeerList) IPs() []string {
	ips := make([]string, 0)
	for _, peer := range pL {
		ips = append(ips, peer.IP.String())
	}
	return ips
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

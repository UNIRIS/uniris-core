package domain

import (
	"fmt"
	"net"
	"time"
)

//Peer represents a peer discovered
type Peer struct {
	PublicKey      []byte
	IP             net.IP
	Port           int
	GenerationTime time.Time
	State          *PeerState
	IsOwned        bool
}

//NewPeer creates a new peer instance
func NewPeer(pbKey []byte, ip net.IP, port int, isOwned bool) Peer {
	return Peer{
		PublicKey:      pbKey,
		IP:             ip,
		Port:           port,
		GenerationTime: time.Now(),
		IsOwned:        isOwned,
	}
}

//Refresh updates peer information
func (p *Peer) Refresh(ip net.IP, port int, gen time.Time, state *PeerState) {
	p.IP = ip
	p.Port = port
	p.State = state
	p.GenerationTime = gen
}

//IsDiscovered checks if a peer contains an app state discovered
func (p *Peer) IsDiscovered() bool {
	return p.State != nil && p.State.Status == Ok
}

//GetElapsedHeartbeats computes the elasted hearbeats from the generation time
func (p Peer) GetElapsedHeartbeats() int64 {
	return time.Now().Unix() - p.GenerationTime.Unix()
}

//GetDiscoveryEndpoint returns the peer endpoint
func (p Peer) GetDiscoveryEndpoint() string {
	return fmt.Sprintf("%s:%d", p.IP.String(), p.Port)
}

//Equals checks if two peers are the same
func (p Peer) Equals(otherPeer Peer) bool {
	return string(p.PublicKey) == string(otherPeer.PublicKey)
}

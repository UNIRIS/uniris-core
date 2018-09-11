package entities

import (
	"net"
	"time"
)

//Peer represents a peer discovered
type Peer struct {
	IP        net.IP
	Port      int
	PublicKey []byte
	Heartbeat PeerHeartbeat
	AppState  PeerAppState
	Category  PeerCategory
}

type PeerCategory int

const (
	SeedCategory       PeerCategory = 1
	DiscoveredCategory PeerCategory = 2
)

//GetElapsedHeartbeats computes the elasted hearbeats from the generation time
func (p Peer) GetElapsedHeartbeats() int64 {
	return time.Now().Unix() - p.Heartbeat.GenerationTime.Unix()
}

//UpdateElapsedHeartbeats refreshes the elapsed heartbeazts
func (p *Peer) UpdateElapsedHeartbeats() {
	p.Heartbeat.ElapsedBeats = p.GetElapsedHeartbeats()
}

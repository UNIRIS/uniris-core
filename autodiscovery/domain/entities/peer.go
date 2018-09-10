package entities

import (
	"net"
	"time"
)

//Peer represents a peer discovered
type Peer struct {
	IP        net.IP
	PublicKey []byte
	Heartbeat PeerHeartbeat
	Details   PeerDetails
}

//GetElapsedHeartbeats computes the elasted hearbeats from the generation time
func (p Peer) GetElapsedHeartbeats() int64 {
	return time.Now().Unix() - p.Heartbeat.GenerationTime.Unix()
}

//SetElapsedHeartbeats stores the elapsed heartbeazts
func (p *Peer) SetElapsedHeartbeats() {
	p.Heartbeat.ElapsedBeats = p.GetElapsedHeartbeats()
}

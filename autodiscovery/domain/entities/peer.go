package entities

import (
	"net"
	"time"
)

//Peer represents a peer discovered
type Peer struct {
	IP        net.IP        `json:"ip"`
	PublicKey []byte        `json:"publicKey"`
	Heartbeat PeerHeartbeat `json:"hearbeat"`
	Details   PeerDetails   `json:"state"`
}

//RefreshHearbeat refresh the generation time of the peer's Heartbeat
func (p *Peer) RefreshHearbeat() {
	p.Heartbeat.GenerationTime = time.Now()
}

//GetElapsedHeartbeats computes of the elasted hearbeat from the generation time
func (p Peer) GetElapsedHeartbeats() int64 {
	return time.Now().Unix() - p.Heartbeat.GenerationTime.Unix()
}

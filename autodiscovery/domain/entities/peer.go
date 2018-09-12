package entities

import (
	"net"
	"time"
)

//Peer represents a peer discovered
type Peer struct {
	IP              net.IP
	Port            int
	PublicKey       []byte
	Heartbeat       PeerHeartbeat
	AppState        PeerAppState
	IsSelf          bool
	DiscoveredNodes int
}

//PeerHeartbeat represents how fresh is the peer information
type PeerHeartbeat struct {
	GenerationTime time.Time
	ElapsedBeats   int64
}

//PeerAppState defines peer caracteristics
type PeerAppState struct {
	State          PeerState
	CPULoad        string
	IOWaitRate     float64
	FreeDiskSpace  float64
	Version        string
	GeoCoordinates Coordinates
	P2PFactor      int
}

//Coordinates represents the peer coordinates
type Coordinates struct {
	Lon float64
	Lat float64
}

//PeerState defines the peer situation (Faulty, Bootstraping or Ok)
type PeerState int

const (
	//FaultyState defines if the peer is not started
	FaultyState PeerState = 0

	//BootstrapingState defines if the peer is starting
	BootstrapingState PeerState = 1

	//OkState defines if the peer is started
	OkState PeerState = 2

	//StorageOnlyState defines if the peer only accept storage request
	StorageOnlyState PeerState = 3
)

//GetElapsedHeartbeats computes the elasted hearbeats from the generation time
func (p Peer) GetElapsedHeartbeats() int64 {
	return time.Now().Unix() - p.Heartbeat.GenerationTime.Unix()
}

//UpdateElapsedHeartbeats refreshes the elapsed heartbeazts
func (p *Peer) UpdateElapsedHeartbeats() {
	p.Heartbeat.ElapsedBeats = p.GetElapsedHeartbeats()
}

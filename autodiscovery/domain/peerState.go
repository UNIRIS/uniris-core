package domain

import "net"

//PeerStatus defines the peer situation (Faulty, Bootstraping or Ok)
type PeerStatus int

const (
	//Fault defines if the peer is not started
	Fault PeerStatus = 0

	//Bootstraping defines if the peer is starting
	Bootstraping PeerStatus = 1

	//Ok defines if the peer is started
	Ok PeerStatus = 2

	//StorageOnly defines if the peer only accept storage request
	StorageOnly PeerStatus = 3
)

//PeerState defines peer app state
type PeerState struct {
	Status        PeerStatus
	CPULoad       string
	IOWaitRate    float64
	FreeDiskSpace float64
	Version       string
	GeoPosition   GeoPosition
	P2PFactor     int
}

//GeoPosition defines the geographic position
type GeoPosition struct {
	IP  net.IP
	Lat float64
	Lon float64
}

//NewPeerState creates a peer state
func NewPeerState(status PeerStatus, ver string, pos GeoPosition, p2pFactor int) *PeerState {
	return &PeerState{
		Status:      status,
		Version:     ver,
		GeoPosition: pos,
		P2PFactor:   p2pFactor,
	}
}

package domain

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
	Lat float64
	Lon float64
}

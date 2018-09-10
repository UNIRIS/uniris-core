package entities

//PeerDetails defines peer caracteristics
type PeerDetails struct {
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

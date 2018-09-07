package entities

//PeerDetails defines peer caracteristics
type PeerDetails struct {
	State          PeerState `json:"state"`
	CPULoadStatus  string    `json:"cpuLoadStatus"`
	IOWaitRate     float64   `json:"ioWaitRate"`
	FreeDiskSpace  float64   `json:"freeDiskSpace"`
	Version        string    `json:"version"`
	GeoCoordinates []float64 `json:"geoCoordinates"`
	P2PFactor      int       `json:"p2pFactor"`
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

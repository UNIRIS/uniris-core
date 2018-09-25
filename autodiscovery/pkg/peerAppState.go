package discovery

import "net"

//PeerAppState describes the state of peer and its metrics
type PeerAppState interface {
	Status() PeerStatus
	Version() string
	CPULoad() string
	FreeDiskSpace() float64
	GeoPosition() PeerPosition
	P2PFactor() int32
}

//PeerStatus defines a peer health analysis
type PeerStatus int

const (
	//BootstrapingStatus defines if the peer is starting
	BootstrapingStatus PeerStatus = iota

	//OkStatus defines if the peer is started
	OkStatus

	//FaultStatus defines if the peer is not started
	FaultStatus

	//StorageOnlyStatus defines if the peer only accept storage request
	StorageOnlyStatus
)

//PeerPosition wraps the geo coordinates of a peer
type PeerPosition struct {
	Lat float64
	Lon float64
}

type appState struct {
	ip            net.IP
	port          int32
	publicKey     PublicKey
	status        PeerStatus
	cpuLoad       string
	freeDiskSpace float64
	version       string
	geoPosition   PeerPosition
	p2pFactor     int32
}

func (a appState) IP() net.IP {
	return a.ip
}

func (a appState) Port() int32 {
	return a.port
}

func (a appState) PublicKey() PublicKey {
	return a.publicKey
}

func (a appState) Status() PeerStatus {
	return a.status
}

func (a appState) CPULoad() string {
	return a.cpuLoad
}

func (a appState) FreeDiskSpace() float64 {
	return a.freeDiskSpace
}

func (a appState) Version() string {
	return a.version
}

func (a appState) P2PFactor() int32 {
	return a.p2pFactor
}

func (a appState) GeoPosition() PeerPosition {
	return a.geoPosition
}

//NewPeerAppState creates a new peer's app state
func NewPeerAppState(ver string, stat PeerStatus, geo PeerPosition, cpu string, disk float64, p2pfactor int32) PeerAppState {
	return appState{
		version:       ver,
		status:        stat,
		geoPosition:   geo,
		cpuLoad:       cpu,
		freeDiskSpace: disk,
		p2pFactor:     p2pfactor,
	}
}

//Refresh the peer state
func (a *appState) refresh(status PeerStatus, disk float64, cpu string, p2pFactor int32) {
	a.cpuLoad = cpu
	a.status = status
	a.freeDiskSpace = disk
	a.p2pFactor = p2pFactor
}

package discovery

import (
	"fmt"
	"net"
	"time"
)

//Repository provides access to the peer repository
type Repository interface {
	ListSeedPeers() ([]Seed, error)
	ListKnownPeers() ([]Peer, error)
	AddPeer(Peer) error
	AddSeed(Seed) error
	UpdatePeer(Peer) error
}

//Seed is initial peer need to startup the discovery process
type Seed struct {
	IP   net.IP
	Port int
}

//ToPeer converts a seed into a peer
func (s Seed) ToPeer() Peer {
	return Peer{
		ip:   s.IP,
		port: s.Port,
	}
}

//Peer describes a member of the P2P network
type Peer struct {
	publicKey      []byte
	ip             net.IP
	port           int
	state          *PeerState
	generationTime time.Time
	isOwned        bool
}

//PeerState describes the state of peer and its metrics
type PeerState struct {
	status        PeerStatus
	cpuLoad       string
	ioWaitRate    float64
	freeDiskSpace float64
	version       string
	geoPosition   PeerPosition
	p2pFactor     int
}

//PeerStatus defines a peer health analysis
type PeerStatus int

const (
	//FaultStatus defines if the peer is not started
	FaultStatus PeerStatus = 0

	//BootstrapingStatus defines if the peer is starting
	BootstrapingStatus PeerStatus = 1

	//OkStatus defines if the peer is started
	OkStatus PeerStatus = 2

	//StorageOnlyStatus defines if the peer only accept storage request
	StorageOnlyStatus PeerStatus = 3
)

//PeerPosition wraps the geo coordinates of a peer
type PeerPosition struct {
	Lat float64
	Lon float64
}

//PublicKey returns the peer's public key. It's the identification of a peer among the network
func (p Peer) PublicKey() []byte {
	return p.publicKey
}

//GenerationTime returns the peer's generation time
func (p Peer) GenerationTime() time.Time {
	return p.generationTime
}

//IsOwned determinates if the peer has been created locally (by startup on this computer)
func (p Peer) IsOwned() bool {
	return p.isOwned
}

//IP returns the peer's IP
func (p Peer) IP() net.IP {
	return p.ip
}

//Port returns the peer's port
func (p Peer) Port() int {
	return p.port
}

//GeoPosition returns the peer's geo coordinates
func (p Peer) GeoPosition() *PeerPosition {
	if p.state == nil {
		return nil
	}
	return &p.state.geoPosition
}

//P2PFactor returns the peer's replication factor
func (p Peer) P2PFactor() int {
	if p.state == nil {
		return 1
	}
	return p.state.p2pFactor
}

//Status returns the peer's status
func (p Peer) Status() PeerStatus {
	if p.state == nil {
		return BootstrapingStatus
	}
	return p.state.status
}

//Version returns the peer's version of the application
func (p Peer) Version() string {
	if p.state == nil {
		return "1.0.0"
	}
	return p.state.version
}

//CPULoad returns the load on the peer's CPU
func (p Peer) CPULoad() string {
	if p.state == nil {
		return "0.0.0"
	}
	return p.state.cpuLoad
}

//FreeDiskSpace returns the available space on the peer's disk
func (p Peer) FreeDiskSpace() float64 {
	if p.state == nil {
		return 0.0
	}
	return p.state.freeDiskSpace
}

//IOWaitRate returns the rate of the peer's I/O operations
func (p Peer) IOWaitRate() float64 {
	if p.state == nil {
		return 0.0
	}
	return p.state.ioWaitRate
}

//IsOk checks if a peer is healthy
func (p Peer) IsOk() bool {
	return p.state.status == OkStatus
}

//GetElapsedHeartbeats returns the elasted hearbeats from the peer's generation time
func (p Peer) GetElapsedHeartbeats() int64 {
	return time.Now().Unix() - p.generationTime.Unix()
}

//GetEndpoint returns the peer endpoint
func (p Peer) GetEndpoint() string {
	return fmt.Sprintf("%s:%d", p.ip.String(), p.port)
}

//Refresh the peer state
func (p *Peer) Refresh(status PeerStatus, disk float64, cpu string, io float64) {
	if p.state == nil {
		p.state = &PeerState{}
	}
	p.state.cpuLoad = cpu
	p.state.status = status
	p.state.freeDiskSpace = disk
	p.state.ioWaitRate = io
}

//NewStartupPeer creates a new peer started on the peer's machine (aka owned peer)
func NewStartupPeer(pbKey []byte, ip net.IP, port int, version string, pos PeerPosition, p2Pfactor int) Peer {
	return Peer{
		ip:             ip,
		port:           port,
		publicKey:      pbKey,
		generationTime: time.Now(),
		isOwned:        true,
		state: &PeerState{
			status:      BootstrapingStatus,
			version:     version,
			geoPosition: pos,
			p2pFactor:   p2Pfactor,
		},
	}
}

//NewDiscoveredPeer creates a peer discovered on the network
func NewDiscoveredPeer(pbKey []byte, ip net.IP, port int, genTime time.Time, state *PeerState) Peer {
	return Peer{
		ip:             ip,
		port:           port,
		publicKey:      pbKey,
		generationTime: genTime,
		isOwned:        false,
		state:          state,
	}
}

//NewState creates a new peer's state
func NewState(ver string, stat PeerStatus, geo PeerPosition, cpu string, disk float64, io float64, p2pfactor int) *PeerState {
	return &PeerState{
		version:       ver,
		status:        stat,
		geoPosition:   geo,
		cpuLoad:       cpu,
		freeDiskSpace: disk,
		ioWaitRate:    io,
		p2pFactor:     p2pfactor,
	}
}

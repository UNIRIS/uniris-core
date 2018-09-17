package discovery

import (
	"fmt"
	"net"
	"time"
)

type Peer struct {
	publicKey      []byte
	ip             net.IP
	port           int
	state          *PeerState
	generationTime time.Time
	isOwned        bool
}

type PeerState struct {
	Status        PeerStatus
	CPULoad       string
	IOWaitRate    float64
	FreeDiskSpace float64
	Version       string
	GeoPosition   PeerPosition
	P2PFactor     int
}

//PeerStatus defines the peer's status according the peer inspector analysis
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

type PeerPosition struct {
	Lat float64
	Lon float64
}

func (p Peer) IP() net.IP {
	return p.ip
}

func (p Peer) Port() int {
	return p.port
}

func (p Peer) GenerationTime() time.Time {
	return p.generationTime
}

func (p Peer) IsOwned() bool {
	return p.isOwned
}

func (p Peer) State() *PeerState {
	return p.state
}

func (p Peer) PublicKey() []byte {
	return p.publicKey
}

//IsOk checks if a peer contains an app state discovered
func (p Peer) IsOk() bool {
	return p.state.Status == OkStatus
}

//ElapsedHearts returns the elasted hearbeats from the generation time
func (p Peer) ElapsedHearts() int64 {
	return time.Now().Unix() - p.generationTime.Unix()
}

//Endpoint returns the peer endpoint
func (p Peer) Endpoint() string {
	return fmt.Sprintf("%s:%d", p.ip.String(), p.port)
}

//Refresh the peer state
func (p *Peer) Refresh(status PeerStatus, disk float64, cpu string, io float64) {
	if p.State == nil {
		p.state = &PeerState{}
	}
	p.state.CPULoad = cpu
	p.state.Status = status
	p.state.FreeDiskSpace = disk
	p.state.IOWaitRate = io
}

func StartPeer(pbKey []byte, ip net.IP, port int, version string, pos PeerPosition, p2Pfactor int) Peer {
	return Peer{
		ip:             ip,
		port:           port,
		publicKey:      pbKey,
		generationTime: time.Now(),
		state: &PeerState{
			Status:      BootstrapingStatus,
			Version:     version,
			GeoPosition: pos,
			P2PFactor:   p2Pfactor,
		},
	}
}

type SeedPeer struct {
	IP   net.IP
	Port int
}

func (s SeedPeer) AsPeer() Peer {
	return Peer{
		ip:   s.IP,
		port: s.Port,
	}
}

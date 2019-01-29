package discovery

import (
	"errors"
	"fmt"
	"net"
	"time"
)

//Seed describes a seeding peer
type Seed struct {
	PeerIdentity
}

func (s Seed) AsPeer() Peer {
	return Peer{
		identity: s.PeerIdentity,
	}
}

//ErrChangeNotOwnedPeer is returned when you try to change the state of peer that you don't own
var ErrChangeNotOwnedPeer = errors.New("Cannot change a peer that you don't own")

//PeerIdentity describes the peer identification the network
type PeerIdentity struct {
	ip        net.IP
	port      int
	publicKey string
}

//NewPeerIdentity creates a new peer identity
func NewPeerIdentity(ip net.IP, port int, pbKey string) PeerIdentity {
	return PeerIdentity{
		ip:        ip,
		port:      port,
		publicKey: pbKey,
	}
}

//IP returns the peer's IP address
func (p PeerIdentity) IP() net.IP {
	return p.ip
}

//Port returns the peer's port
func (p PeerIdentity) Port() int {
	return p.port
}

//PublicKey returns the peer's public key
func (p PeerIdentity) PublicKey() string {
	return p.publicKey
}

//Peer describes a network member
type Peer struct {
	identity PeerIdentity
	hbState  PeerHeartbeatState
	appState PeerAppState
	isLocal  bool
}

//NewLocalPeer creates a new peer started on the miner's machine (aka local peer)
func NewLocalPeer(pbKey string, ip net.IP, port int, version string, lon float64, lat float64) Peer {
	return Peer{
		identity: PeerIdentity{
			ip:        ip,
			port:      port,
			publicKey: pbKey,
		},
		appState: PeerAppState{
			status:  BootstrapingPeer,
			version: version,
			geoPosition: PeerPosition{
				lon: lon,
				lat: lat,
			},
			p2pFactor: 0,
		},
		hbState: PeerHeartbeatState{
			generationTime: time.Now(),
		},
		isLocal: true,
	}
}

//NewDiscoveredPeer creates a peer when including identity, heartbeat and app state
func NewDiscoveredPeer(identity PeerIdentity, hbS PeerHeartbeatState, aS PeerAppState) Peer {
	return Peer{
		identity: identity,
		hbState:  hbS,
		appState: aS,
		isLocal:  false,
	}
}

//NewPeerDigest creates a peer with the minimum information for network transfert
func NewPeerDigest(identity PeerIdentity, hbS PeerHeartbeatState) Peer {
	return Peer{
		identity: identity,
		hbState:  hbS,
	}
}

//Identity returns the peer's identity
func (p Peer) Identity() PeerIdentity {
	return p.identity
}

//HeartbeatState returns the peer's hearbeat state
func (p Peer) HeartbeatState() PeerHeartbeatState {
	return p.hbState
}

//AppState returns the peer's app state including all the metrics
func (p Peer) AppState() PeerAppState {
	return p.appState
}

//IsLocal determinates if the peer has been created locally (by startup on this computer)
func (p Peer) IsLocal() bool {
	return p.isLocal
}

//Endpoint returns the peer endpoint
func (p Peer) Endpoint() string {
	return fmt.Sprintf("%s:%d", p.Identity().IP().String(), p.Identity().Port())
}

//Refresh a peer with metrics and updates the elapsed heartbeats
func (p *Peer) Refresh(status PeerStatus, disk float64, cpu string, p2pFactor int, discoveryPeersNb int) {
	if !p.isLocal {
		return
	}

	p.appState.refresh(status, disk, cpu, p2pFactor, discoveryPeersNb)
	p.hbState.refreshElapsedHeartbeats()
}

func (p Peer) String() string {
	return fmt.Sprintf("Endpoint: %s, Local: %t, %s, %s",
		p.Endpoint(),
		p.IsLocal(),
		p.HeartbeatState().String(),
		p.AppState().String(),
	)
}

//PeerHeartbeatState describes the living state of a peer
type PeerHeartbeatState struct {
	generationTime    time.Time
	elapsedHeartbeats int64
}

//NewPeerHeartbeatState creates a new peer's heartbeat state
func NewPeerHeartbeatState(genTime time.Time, elapsedHb int64) PeerHeartbeatState {
	return PeerHeartbeatState{
		generationTime:    genTime,
		elapsedHeartbeats: elapsedHb,
	}
}

//GenerationTime returns the peer's generation time
func (ph PeerHeartbeatState) GenerationTime() time.Time {
	return ph.generationTime
}

//ElapsedHeartbeats returns the peer's elapsed living seconds from the latest refresh
func (ph PeerHeartbeatState) ElapsedHeartbeats() int64 {
	if ph.elapsedHeartbeats == 0 {
		ph.refreshElapsedHeartbeats()
	}
	return ph.elapsedHeartbeats
}

func (ph *PeerHeartbeatState) refreshElapsedHeartbeats() {
	ph.elapsedHeartbeats = time.Now().Unix() - ph.generationTime.Unix()
}

//MoreRecentThan check if the current heartbeat state is more recent than the another heartbeat state
func (ph PeerHeartbeatState) MoreRecentThan(otherhS PeerHeartbeatState) bool {

	//more recent generation time
	if ph.generationTime.Unix() > otherhS.GenerationTime().Unix() {
		return true
	} else if ph.generationTime.Unix() == otherhS.GenerationTime().Unix() {
		if ph.elapsedHeartbeats == otherhS.ElapsedHeartbeats() {
			return false
		} else if ph.elapsedHeartbeats > otherhS.ElapsedHeartbeats() {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

func (ph PeerHeartbeatState) String() string {
	return fmt.Sprintf("Generation time: %s, Elapsed heartbeats: %d",
		ph.GenerationTime().String(),
		ph.ElapsedHeartbeats())
}

//PeerStatus defines a peer health analysis
type PeerStatus int

const (
	//BootstrapingPeer defines if the peer is starting
	BootstrapingPeer PeerStatus = iota

	//OkPeerStatus defines if the peer is started
	OkPeerStatus

	//FaultyPeer defines if the peer is not started
	FaultyPeer

	//StorageOnlyPeer defines if the peer only accept storage request
	StorageOnlyPeer
)

func (s PeerStatus) String() string {
	if s == OkPeerStatus {
		return "Ok"
	} else if s == BootstrapingPeer {
		return "Bootstraping"
	} else if s == FaultyPeer {
		return "Faulty"
	} else {
		return "StorageOnly"
	}
}

//PeerPosition wraps the geo coordinates of a peer
type PeerPosition struct {
	lat float64
	lon float64
}

func (p PeerPosition) Latitude() float64 {
	return p.lat
}

func (p PeerPosition) Longitude() float64 {
	return p.lon
}

func (pos PeerPosition) String() string {
	return fmt.Sprintf("Lat: %f, Lon: %f", pos.lat, pos.lon)
}

//PeerAppState describes the state of peer and its metrics
type PeerAppState struct {
	status                PeerStatus
	cpuLoad               string
	freeDiskSpace         float64
	version               string
	geoPosition           PeerPosition
	p2pFactor             int
	discoveredPeersNumber int
}

func (a PeerAppState) Status() PeerStatus {
	return a.status
}

func (a PeerAppState) CPULoad() string {
	return a.cpuLoad
}

func (a PeerAppState) FreeDiskSpace() float64 {
	return a.freeDiskSpace
}

func (a PeerAppState) Version() string {
	return a.version
}

func (a PeerAppState) P2PFactor() int {
	return a.p2pFactor
}

func (a PeerAppState) GeoPosition() PeerPosition {
	return a.geoPosition
}

//DiscoveredPeersNumber returns the number of discovered peers
func (a PeerAppState) DiscoveredPeersNumber() int {
	return a.discoveredPeersNumber
}

func (a PeerAppState) String() string {
	return fmt.Sprintf("Status: %s, CPU load: %s, Free disk space: %f, Version: %s, GeoPosition: %s, P2PFactor: %d, Discovered peer number: %d",
		a.Status().String(),
		a.CPULoad(),
		a.FreeDiskSpace(),
		a.Version(),
		a.GeoPosition().String(),
		a.P2PFactor(),
		a.DiscoveredPeersNumber(),
	)
}

//NewPeerAppState creates a new peer's app state
func NewPeerAppState(ver string, stat PeerStatus, lon float64, lat float64, cpu string, disk float64, p2pfactor int, discoveredPeersNumber int) PeerAppState {
	return PeerAppState{
		version: ver,
		status:  stat,
		geoPosition: PeerPosition{
			lon: lon,
			lat: lat,
		},
		cpuLoad:               cpu,
		freeDiskSpace:         disk,
		p2pFactor:             p2pfactor,
		discoveredPeersNumber: discoveredPeersNumber,
	}
}

//Refresh the peer state
func (a *PeerAppState) refresh(status PeerStatus, disk float64, cpu string, p2pFactor int, discoveryPeersNb int) {
	a.cpuLoad = cpu
	a.status = status
	a.freeDiskSpace = disk
	a.p2pFactor = p2pFactor
	a.discoveredPeersNumber = discoveryPeersNb
}

package discovery

import (
	"errors"
	"fmt"
	"net"
	"time"
)

//ErrChangeNotOwnedPeer is returned when you try to change the state of peer that you don't own
var ErrChangeNotOwnedPeer = errors.New("Cannot change a peer that you don't own")

//PublicKey describes a public key value object
type PublicKey []byte

func (k PublicKey) String() string {
	return string(k)
}

//Equals checks if two public key are the same
func (k PublicKey) Equals(key PublicKey) bool {
	return k.String() == key.String()
}

//PeerIdentity describes the peer identification the network
type PeerIdentity interface {
	IP() net.IP
	Port() int32
	PublicKey() PublicKey
}

type peerIdentity struct {
	ip        net.IP
	port      int32
	publicKey PublicKey
}

func NewPeerIdentity(ip net.IP, port int32, pbKey PublicKey) PeerIdentity {
	return peerIdentity{
		ip:        ip,
		port:      port,
		publicKey: pbKey,
	}
}

func (p peerIdentity) IP() net.IP {
	return p.ip
}

func (p peerIdentity) Port() int32 {
	return p.port
}

func (p peerIdentity) PublicKey() PublicKey {
	return p.publicKey
}

//Peer describes a network member
type Peer interface {
	Identity() PeerIdentity
	AppState() PeerAppState
	HeartbeatState() PeerHeartbeatState
	Refresh(status PeerStatus, disk float64, cpu string, p2pFactor int32) error
	Endpoint() string
	Owned() bool
}

//Repository provides access to the peer repository
type Repository interface {
	GetOwnedPeer() (Peer, error)
	ListSeedPeers() ([]Seed, error)
	ListKnownPeers() ([]Peer, error)
	AddPeer(Peer) error
	AddSeed(Seed) error
	UpdatePeer(Peer) error
}

//Peer describes a member of the P2P network
type peer struct {
	identity PeerIdentity
	hbState  PeerHeartbeatState
	appState appState
	isOwned  bool
}

func (p peer) Identity() PeerIdentity {
	return p.identity
}

func (p peer) HeartbeatState() PeerHeartbeatState {
	return p.hbState
}

func (p peer) AppState() PeerAppState {
	return p.appState
}

//Owned determinates if the peer has been created locally (by startup on this computer)
func (p peer) Owned() bool {
	return p.isOwned
}

//Endpoint returns the peer endpoint
func (p peer) Endpoint() string {
	return fmt.Sprintf("%s:%d", p.Identity().IP().String(), p.Identity().Port())
}

func (p *peer) Refresh(status PeerStatus, disk float64, cpu string, p2pFactor int32) error {
	if !p.isOwned {
		return ErrChangeNotOwnedPeer
	}
	p.appState.refresh(status, disk, cpu, p2pFactor)
	return nil
}

//NewStartupPeer creates a new peer started on the peer's machine (aka owned peer)
func NewStartupPeer(pbKey PublicKey, ip net.IP, port int32, version string, pos PeerPosition) Peer {
	return &peer{
		identity: peerIdentity{
			ip:        ip,
			port:      port,
			publicKey: pbKey,
		},
		appState: appState{
			status:      BootstrapingStatus,
			version:     version,
			geoPosition: pos,
			p2pFactor:   1,
		},
		hbState: heartbeatState{
			generationTime: time.Now(),
		},
		isOwned: true,
	}
}

//NewDiscoveredPeer creates a peer when including identity, heartbeat and app state
func NewDiscoveredPeer(identity PeerIdentity, hbS PeerHeartbeatState, aS PeerAppState) Peer {
	return &peer{
		identity: identity,
		hbState:  hbS,
		appState: aS.(appState),
		isOwned:  false,
	}
}

func NewPeerDigest(identity PeerIdentity, hbS PeerHeartbeatState) Peer {
	return &peer{
		identity: identity,
		hbState:  hbS,
	}
}

//Seed is initial peer need to startup the discovery process
type Seed struct {
	IP   net.IP
	Port int32
}

//AsPeer converts a seed into a peer
func (s Seed) AsPeer() Peer {
	return &peer{
		identity: peerIdentity{
			ip:   s.IP,
			port: s.Port,
		},
	}
}

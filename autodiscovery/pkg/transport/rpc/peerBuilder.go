package rpc

import (
	"net"
	"time"

	api "github.com/uniris/uniris-core/autodiscovery/api/protobuf-spec"
	discovery "github.com/uniris/uniris-core/autodiscovery/pkg"
)

//PeerBuilder defines methods to transform API entities for the domain layer
type PeerBuilder struct{}

//ToPeerDigest creates a peer digest from a peer
func (f PeerBuilder) ToPeerDigest(p discovery.Peer) *api.PeerDigest {
	return &api.PeerDigest{
		Identity: &api.PeerIdentity{
			IP:        p.Identity().IP().String(),
			Port:      uint32(p.Identity().Port()),
			PublicKey: p.Identity().PublicKey(),
		},
		HeartbeatState: &api.PeerHeartbeatState{
			GenerationTime:    uint64(p.HeartbeatState().GenerationTime().Unix()),
			ElapsedHeartbeats: p.HeartbeatState().ElapsedHeartbeats(),
		},
	}
}

//FromPeerDigest creates a peer from a digest peer
func (f PeerBuilder) FromPeerDigest(p *api.PeerDigest) discovery.Peer {
	return discovery.NewPeerDigest(
		discovery.NewPeerIdentity(net.ParseIP(p.Identity.IP), uint16(p.Identity.Port), p.Identity.PublicKey),
		discovery.NewPeerHeartbeatState(time.Unix(int64(p.HeartbeatState.GenerationTime), 0), uint64(p.HeartbeatState.ElapsedHeartbeats)),
	)
}

//ToPeerDiscovered creates a discovered peer from a peer
func (f PeerBuilder) ToPeerDiscovered(p discovery.Peer) *api.PeerDiscovered {
	return &api.PeerDiscovered{
		Identity: &api.PeerIdentity{
			IP:        p.Identity().IP().String(),
			Port:      uint32(p.Identity().Port()),
			PublicKey: p.Identity().PublicKey(),
		},
		HeartbeatState: &api.PeerHeartbeatState{
			GenerationTime:    uint64(p.HeartbeatState().GenerationTime().Unix()),
			ElapsedHeartbeats: p.HeartbeatState().ElapsedHeartbeats(),
		},
		AppState: &api.PeerAppState{
			CPULoad:       p.AppState().CPULoad(),
			FreeDiskSpace: float32(p.AppState().FreeDiskSpace()),
			GeoPosition: &api.PeerAppState_GeoCoordinates{
				Lat: float32(p.AppState().GeoPosition().Lat),
				Lon: float32(p.AppState().GeoPosition().Lon),
			},
			P2PFactor: uint32(p.AppState().P2PFactor()),
			Status:    api.PeerAppState_PeerStatus(p.AppState().Status()),
			Version:   p.AppState().Version(),
		},
	}
}

//FromPeerDiscovered creates a peer from a discovered peer
func (f PeerBuilder) FromPeerDiscovered(p *api.PeerDiscovered) discovery.Peer {

	id := discovery.NewPeerIdentity(net.ParseIP(p.Identity.IP), uint16(p.Identity.Port), p.Identity.PublicKey)
	hb := discovery.NewPeerHeartbeatState(time.Unix(int64(p.HeartbeatState.GenerationTime), 0), p.HeartbeatState.ElapsedHeartbeats)
	state := discovery.NewPeerAppState(
		p.AppState.Version,
		discovery.PeerStatus(p.AppState.Status),
		discovery.PeerPosition{
			Lat: float64(p.AppState.GeoPosition.Lat),
			Lon: float64(p.AppState.GeoPosition.Lon),
		},
		p.AppState.CPULoad,
		float64(p.AppState.FreeDiskSpace),
		uint8(p.AppState.P2PFactor),
	)

	return discovery.NewDiscoveredPeer(id, hb, state)
}

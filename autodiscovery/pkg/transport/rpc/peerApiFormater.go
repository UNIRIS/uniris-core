package rpc

import (
	api "github.com/uniris/uniris-core/autodiscovery/api/protobuf-spec"
	discovery "github.com/uniris/uniris-core/autodiscovery/pkg"
)

//PeerAPIFormater defines methods to transform domains entities for the API layer
type PeerAPIFormater struct{}

//BuildPeerDigest converts a domain peer into a digest peer
func (f PeerAPIFormater) BuildPeerDigest(p discovery.Peer) *api.PeerDigest {
	return &api.PeerDigest{
		IP:        p.IP().String(),
		PublicKey: p.PublicKey(),
		Port:      int32(p.Port()),
	}
}

//BuildPeerDigestCollection converts a list of domain peers into a list of digest peers
func (f PeerAPIFormater) BuildPeerDigestCollection(pp []discovery.Peer) []*api.PeerDigest {
	peers := make([]*api.PeerDigest, 0)
	if pp == nil {
		return peers
	}
	for _, p := range pp {
		peers = append(peers, f.BuildPeerDigest(p))
	}
	return peers
}

//BuildPeerDetailed converts a domain peer into a detailed peer
func (f PeerAPIFormater) BuildPeerDetailed(peer discovery.Peer) *api.PeerDetailed {
	p := &api.PeerDetailed{
		PublicKey:      peer.PublicKey(),
		IP:             peer.IP().String(),
		Port:           int32(peer.Port()),
		GenerationTime: peer.GenerationTime().Unix(),
	}

	if p.State != nil {
		p.State = &api.PeerDetailed_PeerState{
			CPULoad:       peer.CPULoad(),
			FreeDiskSpace: float32(peer.FreeDiskSpace()),
			GeoPosition: &api.PeerDetailed_PeerState_GeoCoordinates{
				Lat: float32(peer.GeoPosition().Lat),
				Lon: float32(peer.GeoPosition().Lon),
			},
			IOWaitRate:      float32(peer.IOWaitRate()),
			P2PFactor:       int32(peer.P2PFactor()),
			DiscoveredPeers: int32(peer.DiscoveredPeers()),
			Version:         peer.Version(),
			Status:          api.PeerDetailed_PeerState_PeerStatus(peer.Status()),
		}
	}
	return p
}

//BuildPeerDetailedCollection converts a list of domain peers into a list of detailed peers
func (f PeerAPIFormater) BuildPeerDetailedCollection(pp []discovery.Peer) []*api.PeerDetailed {
	peers := make([]*api.PeerDetailed, 0)
	if pp == nil {
		return peers
	}
	for _, p := range pp {
		peers = append(peers, f.BuildPeerDetailed(p))
	}
	return peers
}

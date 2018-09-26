package rpc

import (
	"net"
	"time"

	api "github.com/uniris/uniris-core/autodiscovery/api/protobuf-spec"
	discovery "github.com/uniris/uniris-core/autodiscovery/pkg"
)

//PeerDomainFormater defines methods to transform API entities for the domain layer
type PeerDomainFormater struct{}

//BuildPeerDigest converst a digest peer into a domain peer digest
func (f PeerDomainFormater) BuildPeerDigest(p *api.PeerDigest) discovery.Peer {
	return discovery.NewPeerDigest(p.PublicKey, net.ParseIP(p.IP), int(p.Port))
}

//BuildPeerDigestCollection converts a list of digest peers into a list of domain peer digest
func (f PeerDomainFormater) BuildPeerDigestCollection(pp []*api.PeerDigest) []discovery.Peer {
	peers := make([]discovery.Peer, 0)
	if pp == nil {
		return peers
	}
	for _, p := range pp {
		peers = append(peers, f.BuildPeerDigest(p))
	}
	return peers
}

//BuildPeerDetailed converts a detailed peer into a domain peer detailed
func (f PeerDomainFormater) BuildPeerDetailed(peer *api.PeerDetailed) discovery.Peer {
	var s *discovery.PeerState
	if peer.State == nil {
		s = &discovery.PeerState{}
	} else {
		s = discovery.NewState(
			peer.State.Version,
			discovery.PeerStatus(peer.State.Status),
			discovery.PeerPosition{
				Lat: float64(peer.State.GeoPosition.Lat),
				Lon: float64(peer.State.GeoPosition.Lon),
			},
			peer.State.CPULoad,
			float64(peer.State.FreeDiskSpace),
			int(peer.State.P2PFactor),
			int(peer.State.DiscoveredPeers),
		)
	}

	return discovery.NewPeerDetailed(peer.PublicKey, net.ParseIP(peer.IP), int(peer.Port), time.Unix(peer.GenerationTime, 0), s)
}

//BuildPeerDetailedCollection converts a list of detailed peer into a list of domain peer detailed
func (f PeerDomainFormater) BuildPeerDetailedCollection(pp []*api.PeerDetailed) []discovery.Peer {
	peers := make([]discovery.Peer, 0)
	if pp == nil {
		return peers
	}
	for _, p := range pp {
		peers = append(peers, f.BuildPeerDetailed(p))
	}
	return peers
}

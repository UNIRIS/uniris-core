package grpc

import (
	"net"
	"time"

	api "github.com/uniris/uniris-core/autodiscovery/api/protobuf-spec"
	discovery "github.com/uniris/uniris-core/autodiscovery/pkg"
	"github.com/uniris/uniris-core/autodiscovery/pkg/gossip"
)

//ToDomainBulk adapts a list of protobuf peer into a list of domain peer
func ToDomainBulk(protoPeers []*api.Peer) []discovery.Peer {
	domainPeers := make([]discovery.Peer, 0)
	for _, peer := range protoPeers {
		domainPeers = append(domainPeers, ToDomain(peer))
	}
	return domainPeers
}

//FromDomainBulk adapts a list of domain peer into a list of protobuf peer
func FromDomainBulk(domainPeers []discovery.Peer) []*api.Peer {
	protoPeers := make([]*api.Peer, 0)
	for _, peer := range domainPeers {
		protoPeers = append(protoPeers, FromDomain(peer))
	}
	return protoPeers
}

//FromDomain adapts a domain peer into a protobuf peer
func FromDomain(peer discovery.Peer) *api.Peer {

	p := &api.Peer{
		PublicKey:      peer.PublicKey(),
		IP:             peer.IP().String(),
		Port:           int32(peer.Port()),
		GenerationTime: peer.GenerationTime().Unix(),
	}

	if p.State != nil {
		p.State = &api.Peer_PeerState{
			CPULoad:       peer.CPULoad(),
			FreeDiskSpace: float32(peer.FreeDiskSpace()),
			GeoPosition: &api.Peer_PeerState_GeoCoordinates{
				Lat: float32(peer.GeoPosition().Lat),
				Lon: float32(peer.GeoPosition().Lon),
			},
			IOWaitRate: float32(peer.IOWaitRate()),
			P2PFactor:  int32(peer.P2PFactor()),
			Version:    peer.Version(),
			Status:     api.Peer_PeerState_PeerStatus(peer.Status()),
		}
	}
	return p
}

//ToDomain adapts a probotuf peer into a domain peer
func ToDomain(peer *api.Peer) discovery.Peer {
	s := discovery.NewState(
		peer.State.Version,
		discovery.PeerStatus(peer.State.Status),
		discovery.PeerPosition{
			Lat: float64(peer.State.GeoPosition.Lat),
			Lon: float64(peer.State.GeoPosition.Lon),
		},
		peer.State.CPULoad,
		float64(peer.State.FreeDiskSpace),
		float64(peer.State.IOWaitRate),
		int(peer.State.P2PFactor),
	)

	return discovery.NewDiscoveredPeer(peer.PublicKey, net.ParseIP(peer.IP), int(peer.Port), time.Unix(peer.GenerationTime, 0), s)
}

//BuildSynRequest formats a protobuf SynRequest
func BuildSynRequest(req gossip.SynRequest) *api.SynRequest {
	return &api.SynRequest{
		Initiator:  FromDomain(req.Initiator),
		Receiver:   FromDomain(req.Receiver),
		KnownPeers: FromDomainBulk(req.KnownPeers),
	}
}

//BuildAckRequest formats a protobuf AckRequest
func BuildAckRequest(req gossip.AckRequest) *api.AckRequest {
	return &api.AckRequest{
		Initiator:      FromDomain(req.Initiator),
		Receiver:       FromDomain(req.Receiver),
		RequestedPeers: FromDomainBulk(req.RequestedPeers),
	}
}

//BuildProtoSynAckResponse formats SynAck response for protobuf
func BuildProtoSynAckResponse(res gossip.SynAck) *api.SynAck {
	return &api.SynAck{
		Initiator:    FromDomain(res.Initiator),
		Receiver:     FromDomain(res.Receiver),
		NewPeers:     FromDomainBulk(res.NewPeers),
		UnknownPeers: FromDomainBulk(res.UnknownPeers),
	}
}

//BuildDomainSynAckResponse formats SynAck response for domain
func BuildDomainSynAckResponse(res *api.SynAck) gossip.SynAck {
	return gossip.SynAck{
		Initiator:    ToDomain(res.Initiator),
		Receiver:     ToDomain(res.Receiver),
		NewPeers:     ToDomainBulk(res.NewPeers),
		UnknownPeers: ToDomainBulk(res.UnknownPeers),
	}
}

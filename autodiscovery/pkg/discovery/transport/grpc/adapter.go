package grpc

import (
	"net"
	"time"

	api "github.com/uniris/uniris-core/autodiscovery/api/protobuf-spec"
	"github.com/uniris/uniris-core/autodiscovery/pkg/discovery"
	"github.com/uniris/uniris-core/autodiscovery/pkg/discovery/gossip"
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
		GenerationTime: peer.GenerationTime(),
	}

	if p.State != nil {
		p.State = &api.Peer_PeerState{
			CPULoad:       peer.State.CPULoad,
			FreeDiskSpace: float32(peer.State.FreeDiskSpace),
			GeoPosition: &api.Peer_PeerState_GeoCoordinates{
				Lat: float32(peer.State.GeoPosition.Lat),
				Lon: float32(peer.State.GeoPosition.Lon),
			},
			IOWaitRate: float32(peer.State.IOWaitRate),
			P2PFactor:  int32(peer.State.P2PFactor),
			Version:    peer.State.Version,
			Status:     api.Peer_PeerState_PeerStatus(peer.State.Status),
		}
	}
	return p
}

//ToDomain adapts a probotuf peer into a domain peer
func ToDomain(peer *api.Peer) discovery.Peer {
	p := discovery.Peer{
		PublicKey:      []byte(peer.PublicKey),
		IP:             net.ParseIP(peer.IP),
		Port:           int(peer.Port),
		GenerationTime: time.Unix(peer.GenerationTime, 0),
	}

	if peer.State != nil {
		p.State = &discovery.PeerState{
			CPULoad:       peer.State.CPULoad,
			FreeDiskSpace: float64(peer.State.FreeDiskSpace),
			GeoPosition: p2p.GeoPosition{
				Lat: float64(peer.State.GeoPosition.Lat),
				Lon: float64(peer.State.GeoPosition.Lon),
			},
			IOWaitRate: float64(peer.State.IOWaitRate),
			P2PFactor:  int(peer.State.P2PFactor),
			Version:    peer.State.Version,
			Status:     p2p.PeerStatus(peer.State.Status),
		}
	}
	return p
}

//BuildSynRequest formats a protobuf SynRequest
func BuildSynRequest(req gossip.SynRequest) *api.SynRequest {
	return &SynRequest{
		Initiator:  FromDomain(req.Initiator),
		Receiver:   FromDomain(req.Receiver),
		KnownPeers: FromDomainBulk(req.KnownPeers),
	}
}

//BuildAckRequest formats a protobuf AckRequest
func BuildAckRequest(req gossip.AckRequest) *api.AckRequest {
	return &AckRequest{
		Initiator:      FromDomain(req.Initiator),
		Receiver:       FromDomain(req.Receiver),
		RequestedPeers: FromDomainBulk(req.RequestedPeers),
	}
}

//BuildProtoSynAckResponse formats SynAck response for protobuf
func BuildProtoSynAckResponse(res gossip.SynAck) *api.SynAck {
	return &SynAck{
		Initiator:    FromDomain(res.Initiator),
		Receiver:     FromDomain(res.Receiver),
		NewPeers:     FromDomainBulk(res.NewPeers),
		UnknownPeers: FromDomainBulk(res.UnknownPeers),
	}
}

//BuildDomainSynAckResponse formats SynAck response for domain
func BuildDomainSynAckResponse(res *gossip.SynAck) api.SynAck {
	return gossip.SynAck{
		Initiator:    ToDomain(res.Initiator),
		Receiver:     ToDomain(res.Receiver),
		NewPeers:     ToDomainBulk(res.NewPeers),
		UnknownPeers: ToDomainBulk(res.UnknownPeers),
	}
}

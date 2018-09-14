package adapters

import (
	"net"
	"time"

	"github.com/uniris/uniris-core/autodiscovery/domain"
	"github.com/uniris/uniris-core/autodiscovery/infrastructure/proto"
)

//ToDomainBulk adapts a list of protobuf peer into a list of domain peer
func ToDomainBulk(protoPeers []*proto.Peer) []domain.Peer {
	domainPeers := make([]domain.Peer, 0)
	for _, peer := range protoPeers {
		domainPeers = append(domainPeers, ToDomain(peer))
	}
	return domainPeers
}

//FromDomainBulk adapts a list of domain peer into a list of protobuf peer
func FromDomainBulk(domainPeers []domain.Peer) []*proto.Peer {
	protoPeers := make([]*proto.Peer, 0)
	for _, peer := range domainPeers {
		protoPeers = append(protoPeers, FromDomain(peer))
	}
	return protoPeers
}

//FromDomain adapts a domain peer into a protobuf peer
func FromDomain(peer domain.Peer) *proto.Peer {
	return &proto.Peer{
		PublicKey:      peer.PublicKey,
		IP:             peer.IP.String(),
		Port:           int32(peer.Port),
		GenerationTime: peer.GenerationTime.Unix(),
		State: &proto.Peer_PeerState{
			CPULoad:       peer.State.CPULoad,
			FreeDiskSpace: float32(peer.State.FreeDiskSpace),
			GeoPosition: &proto.Peer_PeerState_GeoCoordinates{
				Lat: float32(peer.State.GeoPosition.Lat),
				Lon: float32(peer.State.GeoPosition.Lon),
			},
			IOWaitRate: float32(peer.State.IOWaitRate),
			P2PFactor:  int32(peer.State.P2PFactor),
			Version:    peer.State.Version,
			Status:     proto.Peer_PeerState_PeerStatus(peer.State.Status),
		},
	}
}

//ToDomain adapts a probotuf peer into a domain peer
func ToDomain(peer *proto.Peer) domain.Peer {
	return domain.Peer{
		PublicKey:      []byte(peer.PublicKey),
		IP:             net.ParseIP(peer.IP),
		Port:           int(peer.Port),
		GenerationTime: time.Unix(peer.GenerationTime, 0),
		State: &domain.PeerState{
			CPULoad:       peer.State.CPULoad,
			FreeDiskSpace: float64(peer.State.FreeDiskSpace),
			GeoPosition: domain.GeoPosition{
				Lat: float64(peer.State.GeoPosition.Lat),
				Lon: float64(peer.State.GeoPosition.Lon),
			},
			IOWaitRate: float64(peer.State.IOWaitRate),
			P2PFactor:  int(peer.State.P2PFactor),
			Version:    peer.State.Version,
			Status:     domain.PeerStatus(peer.State.Status),
		},
	}
}

//BuildSynRequest formats a protobuf SynRequest
func BuildSynRequest(sender domain.Peer, knownPeers []domain.Peer) *proto.SynRequest {
	protoPeers := FromDomainBulk(knownPeers)
	return &proto.SynRequest{
		Sender:     FromDomain(sender),
		KnownPeers: protoPeers,
	}
}

//BuildAckRequest formats a protobuf AckRequest
func BuildAckRequest(detailedKnownPeers []domain.Peer) *proto.AckRequest {
	protoPeers := FromDomainBulk(detailedKnownPeers)
	return &proto.AckRequest{
		DetailedKnownPeers: protoPeers,
	}
}

//BuildProtoSynAckResponse formats SynAck response for protobuf
func BuildProtoSynAckResponse(newPeers []domain.Peer, detailedPeersRequested []domain.Peer) *proto.SynAck {
	protoNewPeers := FromDomainBulk(newPeers)
	protoDetailedPeersRequested := FromDomainBulk(detailedPeersRequested)
	return &proto.SynAck{
		NewPeers:             protoNewPeers,
		DetailPeersRequested: protoDetailedPeersRequested,
	}
}

//BuildDomainSynAckResponse formats SynAck response for domain
func BuildDomainSynAckResponse(req *proto.SynAck) *domain.SynAck {
	newPeers := ToDomainBulk(req.NewPeers)
	detailedPeersRequested := ToDomainBulk(req.DetailPeersRequested)

	return &domain.SynAck{
		NewPeers:             newPeers,
		DetailPeersRequested: detailedPeersRequested,
	}
}

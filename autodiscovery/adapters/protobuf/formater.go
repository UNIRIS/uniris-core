package protobuf

import (
	"net"
	"time"

	"github.com/uniris/uniris-core/autodiscovery/core/domain"
)

//ToDomainBulk adapts a list of protobuf peer into a list of domain peer
func ToDomainBulk(protoPeers []*Peer) []domain.Peer {
	domainPeers := make([]domain.Peer, 0)
	for _, peer := range protoPeers {
		domainPeers = append(domainPeers, ToDomain(peer))
	}
	return domainPeers
}

//FromDomainBulk adapts a list of domain peer into a list of protobuf peer
func FromDomainBulk(domainPeers []domain.Peer) []*Peer {
	protoPeers := make([]*Peer, 0)
	for _, peer := range domainPeers {
		protoPeers = append(protoPeers, FromDomain(peer))
	}
	return protoPeers
}

//FromDomain adapts a domain peer into a protobuf peer
func FromDomain(peer domain.Peer) *Peer {
	p := &Peer{
		PublicKey:      peer.PublicKey,
		IP:             peer.IP.String(),
		Port:           int32(peer.Port),
		GenerationTime: peer.GenerationTime.Unix(),
	}

	if p.State != nil {
		p.State = &Peer_PeerState{
			CPULoad:       peer.State.CPULoad,
			FreeDiskSpace: float32(peer.State.FreeDiskSpace),
			GeoPosition: &Peer_PeerState_GeoCoordinates{
				Lat: float32(peer.State.GeoPosition.Lat),
				Lon: float32(peer.State.GeoPosition.Lon),
			},
			IOWaitRate: float32(peer.State.IOWaitRate),
			P2PFactor:  int32(peer.State.P2PFactor),
			Version:    peer.State.Version,
			Status:     Peer_PeerState_PeerStatus(peer.State.Status),
		}
	}
	return p
}

//ToDomain adapts a probotuf peer into a domain peer
func ToDomain(peer *Peer) domain.Peer {
	p := domain.Peer{
		PublicKey:      []byte(peer.PublicKey),
		IP:             net.ParseIP(peer.IP),
		Port:           int(peer.Port),
		GenerationTime: time.Unix(peer.GenerationTime, 0),
	}

	if peer.State != nil {
		p.State = &domain.PeerState{
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
		}
	}
	return p
}

//BuildSynRequest formats a protobuf SynRequest
func BuildSynRequest(req domain.SynRequest) *SynRequest {
	return &SynRequest{
		Initiator:  FromDomain(req.Initiator),
		Receiver:   FromDomain(req.Receiver),
		KnownPeers: FromDomainBulk(req.KnownPeers),
	}
}

//BuildAckRequest formats a protobuf AckRequest
func BuildAckRequest(req domain.AckRequest) *AckRequest {
	return &AckRequest{
		Initiator:      FromDomain(req.Initiator),
		Receiver:       FromDomain(req.Receiver),
		RequestedPeers: FromDomainBulk(req.RequestedPeers),
	}
}

//BuildProtoSynAckResponse formats SynAck response for protobuf
func BuildProtoSynAckResponse(res domain.SynAck) *SynAck {
	return &SynAck{
		Initiator:    FromDomain(res.Initiator),
		Receiver:     FromDomain(res.Receiver),
		NewPeers:     FromDomainBulk(res.NewPeers),
		UnknownPeers: FromDomainBulk(res.UnknownPeers),
	}
}

//BuildDomainSynAckResponse formats SynAck response for domain
func BuildDomainSynAckResponse(res *SynAck) domain.SynAck {
	return domain.SynAck{
		Initiator:    ToDomain(res.Initiator),
		Receiver:     ToDomain(res.Receiver),
		NewPeers:     ToDomainBulk(res.NewPeers),
		UnknownPeers: ToDomainBulk(res.UnknownPeers),
	}
}

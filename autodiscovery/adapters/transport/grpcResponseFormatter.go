package transport

import (
	"github.com/uniris/uniris-core/autodiscovery/domain/entities"
)

//FormatAcknownledgeReponseForDomain adapts the ACK GRPC message for domain
func FormatAcknownledgeReponseForDomain(ack *AcknowledgeResponse) *entities.AcknowledgeResponse {
	unknownSenderPeers := make([]*entities.Peer, 0)
	unknownReceiverPeer := make([]*entities.Peer, 0)

	for _, peer := range ack.UnknownSenderPeers {
		unknownSenderPeers = append(unknownSenderPeers, FormatPeerToDomain(*peer))
	}
	for _, peer := range ack.UnknownSenderPeers {
		unknownReceiverPeer = append(unknownReceiverPeer, FormatPeerToDomain(*peer))
	}
	return &entities.AcknowledgeResponse{
		UnknownReceiverPeers: unknownReceiverPeer,
		UnknownSenderPeers:   unknownSenderPeers,
	}
}

//FormatAcknownledgeReponseForGRPC adapts the ACK domain response for GRPC message
func FormatAcknownledgeReponseForGRPC(ack *entities.AcknowledgeResponse) *AcknowledgeResponse {
	unknownSenderPeers := make([]*Peer, 0)
	unknownReceiverPeer := make([]*Peer, 0)

	for _, peer := range ack.UnknownSenderPeers {
		unknownSenderPeers = append(unknownSenderPeers, FormatPeerToGrpc(peer))
	}
	for _, peer := range ack.UnknownReceiverPeers {
		unknownReceiverPeer = append(unknownReceiverPeer, FormatPeerToGrpc(peer))
	}

	return &AcknowledgeResponse{
		UnknownSenderPeers:   unknownSenderPeers,
		UnknownReceiverPeers: unknownReceiverPeer,
	}
}

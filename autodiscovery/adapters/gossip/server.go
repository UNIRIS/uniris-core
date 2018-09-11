package gossip

import (
	"context"

	"github.com/uniris/uniris-core/autodiscovery/domain/entities"
	"github.com/uniris/uniris-core/autodiscovery/domain/usecases"

	"github.com/uniris/uniris-core/autodiscovery/domain/repositories"
)

//Server implements the GRPC server
type Server struct {
	PeerRepo repositories.PeerRepository
}

//Synchronize implements the GRPC Synchronize handler
func (s *Server) Synchronize(ctx context.Context, req *SynchronizeRequest) (*AcknowledgeResponse, error) {
	initiatorPeers := make([]*entities.Peer, 0)
	for _, peer := range req.GetKnownPeers() {
		initiatorPeers = append(initiatorPeers, FormatPeerToDomain(*peer))
	}

	ack, err := usecases.AcknowledgeRequest(s.PeerRepo, initiatorPeers)
	if err != nil {
		return nil, err
	}

	unknownInitiatorPeers := make([]*Peer, 0)
	wishedUnknownPeers := make([]*Peer, 0)

	for _, peer := range ack.UnknownInitiatorPeers {
		unknownInitiatorPeers = append(unknownInitiatorPeers, FormatPeerToGrpc(peer))
	}
	for _, peer := range ack.WishedUnknownPeers {
		wishedUnknownPeers = append(wishedUnknownPeers, FormatPeerToGrpc(peer))
	}

	return &AcknowledgeResponse{
		UnknownInitiatorPeers: unknownInitiatorPeers,
		WishedUnknownPeers:    wishedUnknownPeers,
	}, nil
}

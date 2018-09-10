package gossip

import (
	"context"

	"github.com/uniris/uniris-core/autodiscovery/domain/entities"

	"github.com/uniris/uniris-core/autodiscovery/domain/repositories"
	"github.com/uniris/uniris-core/autodiscovery/domain/usecases"
)

//Server implements the GRPC server
type Server struct {
	PeerRepo repositories.PeerRepository
}

func (s *Server) DiscoverPeers(ctx context.Context, req *DiscoveryRequest) (*DiscoveryResponse, error) {
	newPeers := make([]*entities.Peer, 0)
	for _, peer := range req.GetKnownPeers() {
		newPeers = append(newPeers, FormatPeerToDomain(*peer))
	}
	unknownPeers, err := usecases.GetUnknownPeers(s.PeerRepo, newPeers)
	if err != nil {
		return nil, err
	}

	grpcPeers := make([]*Peer, 0)
	for _, peer := range unknownPeers {
		grpcPeers = append(grpcPeers, FormatPeerToGrpc(peer))
	}

	return &DiscoveryResponse{
		Peers: grpcPeers,
	}, nil
}

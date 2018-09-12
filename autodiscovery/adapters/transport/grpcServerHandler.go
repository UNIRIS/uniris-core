package transport

import (
	"context"

	"github.com/uniris/uniris-core/autodiscovery/domain/entities"
	"github.com/uniris/uniris-core/autodiscovery/domain/usecases"

	"github.com/uniris/uniris-core/autodiscovery/domain/repositories"
)

//GrpcServer implements the GRPC server
type GrpcServer struct {
	PeerRepo repositories.PeerRepository
}

//Synchronize implements the GRPC Synchronize handler
func (s *GrpcServer) Synchronize(ctx context.Context, req *SynchronizeRequest) (*AcknowledgeResponse, error) {
	knownSenderPeers := make([]*entities.Peer, 0)
	for _, peer := range req.GetKnownPeers() {
		knownSenderPeers = append(knownSenderPeers, FormatPeerToDomain(*peer))
	}

	receiverKnownPeers, err := s.PeerRepo.ListPeers()
	if err != nil {
		return nil, err
	}

	unknownSenderPeers := usecases.GetUnknownPeers(receiverKnownPeers, knownSenderPeers)
	unknownReceiverPeers := usecases.GetUnknownPeers(knownSenderPeers, receiverKnownPeers)

	err = usecases.SetNewPeers(s.PeerRepo, knownSenderPeers)
	if err != nil {
		return nil, err
	}

	return FormatAcknownledgeReponseForGRPC(&entities.AcknowledgeResponse{
		UnknownSenderPeers:   unknownSenderPeers,
		UnknownReceiverPeers: unknownReceiverPeers,
	}), nil
}

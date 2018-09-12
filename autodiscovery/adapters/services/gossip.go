package services

import (
	"context"
	"fmt"

	"github.com/uniris/uniris-core/autodiscovery/adapters/transport"
	"github.com/uniris/uniris-core/autodiscovery/domain/entities"
	"google.golang.org/grpc"
)

//GossipService implements the autodiscovery requests
type GossipService struct{}

//Synchronize call GRPC method to retrieve unknown peers (SYN handshake)
func (s GossipService) Synchronize(synReq *entities.SynchronizationRequest) (*entities.AcknowledgeResponse, error) {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", synReq.PeerReceiver.IP.String(), synReq.PeerReceiver.Port), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client := transport.NewGossipGrpcServiceClient(conn)

	knownSenderPeers := make([]*transport.Peer, 0)
	for _, peer := range synReq.KnownSenderPeers {
		knownSenderPeers = append(knownSenderPeers, transport.FormatPeerToGrpc(peer))
	}

	resp, err := client.Synchronize(context.Background(), &transport.SynchronizeRequest{
		KnownPeers: knownSenderPeers,
	})
	if err != nil {
		return nil, err
	}

	return transport.FormatAcknownledgeReponseForDomain(resp), nil
}

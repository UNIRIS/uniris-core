package services

import (
	"context"
	"fmt"

	"github.com/uniris/uniris-core/autodiscovery/adapters/gossip"
	"github.com/uniris/uniris-core/autodiscovery/domain/entities"
	"google.golang.org/grpc"
)

//GossipService implements the autodiscovery requests
type GossipService struct{}

//Synchronize call GRPC method to retrieve unknown peers (SYN handshake)
func (s GossipService) Synchronize(destPeer *entities.Peer, knownPeers []*entities.Peer) (*entities.Acknowledge, error) {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", destPeer.IP.String(), destPeer.Port), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client := gossip.NewGossipServiceClient(conn)
	discoveryRequestPeers := make([]*gossip.Peer, 0)
	for _, peer := range knownPeers {
		discoveryRequestPeers = append(discoveryRequestPeers, gossip.FormatPeerToGrpc(peer))
	}

	resp, err := client.Synchronize(context.Background(), &gossip.SynchronizeRequest{
		KnownPeers: discoveryRequestPeers,
	})
	if err != nil {
		return nil, err
	}

	unknownInitiatorPeers := make([]*entities.Peer, 0)
	wishedUnknownPeers := make([]*entities.Peer, 0)
	for _, peer := range resp.UnknownInitiatorPeers {
		unknownInitiatorPeers = append(unknownInitiatorPeers, gossip.FormatPeerToDomain(*peer))
	}
	for _, peer := range resp.WishedUnknownPeers {
		wishedUnknownPeers = append(wishedUnknownPeers, gossip.FormatPeerToDomain(*peer))
	}

	return &entities.Acknowledge{
		UnknownInitiatorPeers: unknownInitiatorPeers,
		WishedUnknownPeers:    wishedUnknownPeers,
	}, nil
}

package services

import (
	"context"
	"fmt"

	"github.com/uniris/uniris-core/autodiscovery/adapters/gossip"
	"github.com/uniris/uniris-core/autodiscovery/domain/entities"
	"github.com/uniris/uniris-core/autodiscovery/infrastructure"
	"google.golang.org/grpc"
)

//GossipService implements the autodiscovery requests
type GossipService struct{}

//DiscoverPeers call GRPC method to retrieve unknown peers (SYN handshake)
func (s GossipService) DiscoverPeers(destPeer entities.Peer, knownPeers []entities.Peer) ([]*entities.Peer, error) {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", destPeer.IP.String(), infrastructure.GrpcPort))
	if err != nil {
		return nil, err
	}

	client := gossip.NewGossipServiceClient(conn)
	discoveryRequestPeers := make([]*gossip.Peer, 0)
	for _, peer := range knownPeers {
		discoveryRequestPeers = append(discoveryRequestPeers, FormatPeerToGrpc(peer))
	}
	resp, err := client.DiscoverPeers(context.Background(), &DiscoveryRequest{
		KnownPeers: discoveryRequestPeers,
	})

	if err != nil {
		return nil, err
	}

	discoveryResponsePeers := make([]*entities.Peer, 0)
	for _, peer := range resp.Peers {
		discoveryResponsePeers = append(discoveryResponsePeers, FormatPeerToDomain(*peer))
	}

	return discoveryResponsePeers, nil
}

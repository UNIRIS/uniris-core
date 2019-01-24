package rpc

import (
	"context"
	"fmt"
	"time"

	"github.com/uniris/uniris-core/pkg/gossip"

	api "github.com/uniris/uniris-core/api/protobuf-spec"
	uniris "github.com/uniris/uniris-core/pkg"
)

type discoverySrv struct {
	gossip gossip.Service
}

//NewDiscoveryServer creates a new GRPC discovery server
func NewDiscoveryServer(g gossip.Service) api.DiscoveryServiceServer {
	return discoverySrv{
		gossip: g,
	}
}

func (s discoverySrv) Synchronize(ctx context.Context, req *api.SynRequest) (*api.SynResponse, error) {
	fmt.Printf("SYN REQUEST - %s\n", time.Unix(req.Timestamp, 0).String())

	reqP := make([]uniris.Peer, 0)
	for _, p := range req.KnownPeers {
		reqP = append(reqP, formatPeerDigest(p))
	}

	unknown, new, err := s.gossip.CompareSyncRequest(reqP)
	if err != nil {
		return nil, err
	}

	unknownPeers := make([]*api.PeerDigest, 0)
	for _, p := range unknown {
		unknownPeers = append(unknownPeers, formatPeerDigestAPI(p))
	}

	newPeers := make([]*api.PeerDiscovered, 0)
	for _, p := range new {
		newPeers = append(newPeers, formatPeerDiscoveredAPI(p))
	}

	return &api.SynResponse{
		Target:       req.Source,
		Source:       req.Target,
		UnknownPeers: unknownPeers,
		NewPeers:     newPeers,
		Timestamp:    time.Now().Unix(),
	}, nil
}

func (s discoverySrv) Acknowledge(ctx context.Context, req *api.AckRequest) (*api.AckResponse, error) {
	fmt.Printf("ACK REQUEST - %s\n", time.Unix(req.Timestamp, 0).String())

	newPeers := make([]uniris.Peer, 0)
	for _, p := range req.RequestedPeers {
		newPeers = append(newPeers, formatPeerDiscovered(p))
	}

	if err := s.gossip.StoreAcknowledgePeers(newPeers); err != nil {
		return nil, err
	}

	return &api.AckResponse{
		Timestamp: time.Now().Unix(),
	}, nil
}

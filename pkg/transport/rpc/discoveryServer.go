package rpc

import (
	"context"
	"fmt"
	"time"

	api "github.com/uniris/uniris-core/api/protobuf-spec"
	"github.com/uniris/uniris-core/pkg/discovery"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type discoverySrv struct {
	db    discovery.Database
	notif discovery.Notifier
}

//NewDiscoveryServer creates a new GRPC server for the discovery service
func NewDiscoveryServer(db discovery.Database, n discovery.Notifier) api.DiscoveryServiceServer {
	return &discoverySrv{
		db:    db,
		notif: n,
	}
}

func (s discoverySrv) Synchronize(ctx context.Context, req *api.SynRequest) (*api.SynResponse, error) {
	fmt.Printf("SYN REQUEST - %s\n", time.Unix(req.Timestamp, 0).String())

	receivedPeers := make([]discovery.Peer, 0)
	for _, p := range req.KnownPeers {
		receivedPeers = append(receivedPeers, formatPeerDigest(p))
	}

	localPeers, err := s.db.DiscoveredPeers()
	if err != nil {
		return nil, status.New(codes.Internal, err.Error()).Err()
	}

	//Get the unknowns peers or more recent than the local peers
	localDiscoveries := make([]*api.PeerDigest, 0)
	for _, p := range discovery.ExcludedOrRecent(localPeers, receivedPeers) {
		localDiscoveries = append(localDiscoveries, formatPeerDigestAPI(p))
	}

	//get the unknown peers or more recent than the received peers
	remoteDiscoveries := make([]*api.PeerDiscovered, 0)
	for _, p := range discovery.ExcludedOrRecent(receivedPeers, localPeers) {
		remoteDiscoveries = append(remoteDiscoveries, formatPeerDiscoveredAPI(p))
	}

	return &api.SynResponse{
		LocalDiscoveries: localDiscoveries,
		RemoteDiscoveris: remoteDiscoveries,
		Timestamp:        time.Now().Unix(),
	}, nil
}

func (s discoverySrv) Acknowledge(ctx context.Context, req *api.AckRequest) (*api.AckResponse, error) {
	fmt.Printf("ACK REQUEST - %s\n", time.Unix(req.Timestamp, 0).String())

	newPeers := make([]discovery.Peer, 0)
	for _, p := range req.RequestedPeers {
		newPeers = append(newPeers, formatPeerDiscovered(p))
	}

	for _, p := range newPeers {
		if err := s.db.WriteDiscoveredPeer(p); err != nil {
			return nil, status.New(codes.Internal, err.Error()).Err()
		}
		if err := s.notif.NotifyDiscovery(p); err != nil {
			return nil, status.New(codes.Internal, err.Error()).Err()
		}
	}

	return &api.AckResponse{
		Timestamp: time.Now().Unix(),
	}, nil
}

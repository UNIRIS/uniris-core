package grpc

import (
	"context"

	discovery "github.com/uniris/uniris-core/autodiscovery/pkg"
	"github.com/uniris/uniris-core/autodiscovery/pkg/inspecting"

	"github.com/golang/protobuf/ptypes/empty"
	api "github.com/uniris/uniris-core/autodiscovery/api/protobuf-spec"
	"github.com/uniris/uniris-core/autodiscovery/pkg/gossip"
)

type Handler interface {
	Synchronize(ctx context.Context, req *api.SynRequest) (*api.SynAck, error)
	Acknowledge(ctx context.Context, req *api.AckRequest) (*empty.Empty, error)
}

type handler struct {
	ownedPeer discovery.Peer
	repo      discovery.Repository
	pr        inspecting.PeerMonitor
}

//Synchronize implements the protobuf Synchronize request handler
func (h handler) Synchronize(ctx context.Context, req *api.SynRequest) (*api.SynAck, error) {

	initiator := ToDomain(req.Initiator)
	receiver := ToDomain(req.Receiver)
	receivedPeers := ToDomainBulk(req.KnownPeers)

	inspecting.NewService(h.repo, h.pr).RefreshPeer(&h.ownedPeer)

	g := gossip.NewService(h.repo, nil, nil)
	diff, err := g.DiffPeers(receivedPeers)
	if err != nil {
		return nil, err
	}

	synAck := gossip.NewSynAck(initiator, receiver, diff.UnknownRemotly, diff.UnknownLocally)
	return BuildProtoSynAckResponse(synAck), nil
}

//Acknowledge implements the protobuf Acknowledge request handler
func (h handler) Acknowledge(ctx context.Context, req *api.AckRequest) (*empty.Empty, error) {
	//Store the peers requested
	for _, p := range req.RequestedPeers {
		h.repo.AddPeer(ToDomain(p))
	}
	return nil, nil
}

//NewHandler create a new GRPC handler
func NewHandler(op discovery.Peer, repo discovery.Repository, pr inspecting.PeerMonitor) Handler {
	return handler{
		repo:      repo,
		ownedPeer: op,
		pr:        pr,
	}
	// grpcServer := grpc.NewServer()
	// api.RegisterDiscoveryServer(grpcServer, GrpcHandler{})
	// return grpcServer
}

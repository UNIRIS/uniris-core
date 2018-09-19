package rpc

import (
	"context"
	"log"

	discovery "github.com/uniris/uniris-core/autodiscovery/pkg"

	"github.com/golang/protobuf/ptypes/empty"
	api "github.com/uniris/uniris-core/autodiscovery/api/protobuf-spec"
	"github.com/uniris/uniris-core/autodiscovery/pkg/gossip"
)

//Handler is the interface that provides methods to handle GRPC requests
type Handler interface {
	Synchronize(ctx context.Context, req *api.SynRequest) (*api.SynAck, error)
	Acknowledge(ctx context.Context, req *api.AckRequest) (*empty.Empty, error)
}

type handler struct {
	repo         discovery.Repository
	domainFormat PeerDomainFormater
	apiFormat    PeerAPIFormater
}

//Synchronize implements the protobuf Synchronize request handler
func (h handler) Synchronize(ctx context.Context, req *api.SynRequest) (*api.SynAck, error) {
	// FOR DEBUG
	// init := h.domainFormat.BuildPeerDigest(req.Initiator)
	// log.Printf("Syn request received from %s", init.GetEndpoint())

	receivedPeers := h.domainFormat.BuildPeerDigestCollection(req.KnownPeers)

	g := gossip.NewService(h.repo, nil, nil, nil)
	diff, err := g.DiffPeers(receivedPeers)
	if err != nil {
		return nil, err
	}

	return &api.SynAck{
		Initiator:    req.Receiver,
		Receiver:     req.Initiator,
		NewPeers:     h.apiFormat.BuildPeerDetailedCollection(diff.UnknownRemotly),
		UnknownPeers: h.apiFormat.BuildPeerDigestCollection(diff.UnknownLocally),
	}, nil
}

//Acknowledge implements the protobuf Acknowledge request handler
func (h handler) Acknowledge(ctx context.Context, req *api.AckRequest) (*empty.Empty, error) {
	// FOR DEBUG
	// init := h.domainFormat.BuildPeerDigest(req.Initiator)
	// log.Printf("Ack request received from %s", init.GetEndpoint())

	//Store the peers requested
	for _, p := range req.RequestedPeers {
		log.Printf("New peer discovered %s", h.domainFormat.BuildPeerDetailed(p).GetEndpoint())
		h.repo.AddPeer(h.domainFormat.BuildPeerDetailed(p))
	}
	return new(empty.Empty), nil
}

//NewHandler create a new GRPC handler
func NewHandler(repo discovery.Repository) Handler {
	return handler{
		repo:         repo,
		domainFormat: PeerDomainFormater{},
		apiFormat:    PeerAPIFormater{},
	}
}

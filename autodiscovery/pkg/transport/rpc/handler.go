package rpc

import (
	"golang.org/x/net/context"

	discovery "github.com/uniris/uniris-core/autodiscovery/pkg"

	"github.com/golang/protobuf/ptypes/empty"
	api "github.com/uniris/uniris-core/autodiscovery/api/protobuf-spec"
	"github.com/uniris/uniris-core/autodiscovery/pkg/gossip"
	"github.com/uniris/uniris-core/autodiscovery/pkg/monitoring"
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
	gos          gossip.Service
	mon          monitoring.Service
	notif        gossip.Notifier
}

//Synchronize implements the protobuf Synchronize request handler
func (h handler) Synchronize(ctx context.Context, req *api.SynRequest) (*api.SynAck, error) {
	// FOR DEBUG
	// init := h.domainFormat.BuildPeerDigest(req.Initiator)
	// log.Printf("Syn request received from %s", init.Endpoint())

	receivedPeers := h.domainFormat.BuildPeerDigestCollection(req.KnownPeers)

	//Update metrics of own peer before to communicate known peers
	if err := h.mon.RefreshOwnedPeer(); err != nil {
		return nil, err
	}

	//Get the diff between known peers and the received peers
	diff, err := h.gos.ComparePeers(receivedPeers)
	if err != nil {
		return nil, err
	}

	return &api.SynAck{
		Initiator:    req.Target,
		Target:       req.Initiator,
		NewPeers:     h.apiFormat.BuildPeerDetailedCollection(diff.UnknownRemotly),
		UnknownPeers: h.apiFormat.BuildPeerDigestCollection(diff.UnknownLocally),
	}, nil
}

//Acknowledge implements the protobuf Acknowledge request handler
func (h handler) Acknowledge(ctx context.Context, req *api.AckRequest) (*empty.Empty, error) {
	// FOR DEBUG
	// init := h.domainFormat.BuildPeerDigest(req.Initiator)
	// log.Printf("Ack request received from %s", init.GetEndpoint())

	//Store the peers requested and notifies them
	for _, rp := range req.RequestedPeers {
		p := h.domainFormat.BuildPeerDetailed(rp)
		h.notif.Notify(p)
		h.repo.SetPeer(p)
	}
	return new(empty.Empty), nil
}

//NewHandler create a new GRPC handler
func NewHandler(repo discovery.Repository, gos gossip.Service, mon monitoring.Service, notif gossip.Notifier) Handler {
	return handler{
		repo:         repo,
		domainFormat: PeerDomainFormater{},
		apiFormat:    PeerAPIFormater{},
		gos:          gos,
		mon:          mon,
		notif:        notif,
	}
}

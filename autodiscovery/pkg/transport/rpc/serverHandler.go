package rpc

import (
	"golang.org/x/net/context"

	"github.com/uniris/uniris-core/autodiscovery/pkg/comparing"

	discovery "github.com/uniris/uniris-core/autodiscovery/pkg"

	"github.com/golang/protobuf/ptypes/empty"
	api "github.com/uniris/uniris-core/autodiscovery/api/protobuf-spec"
	"github.com/uniris/uniris-core/autodiscovery/pkg/gossip"
)

type srvHandler struct {
	repo  discovery.Repository
	notif gossip.Notifier
}

//Synchronize implements the protobuf Synchronize request handler
func (h srvHandler) Synchronize(ctx context.Context, req *api.SynRequest) (*api.SynAck, error) {
	builder := PeerBuilder{}

	reqP := make([]discovery.Peer, 0)
	for _, p := range req.KnownPeers {
		reqP = append(reqP, builder.FromPeerDigest(p))
	}

	knownPeers, err := h.repo.ListKnownPeers()
	if err != nil {
		return nil, err
	}

	newPeers := make([]*api.PeerDiscovered, 0)
	unknownPeers := make([]*api.PeerDigest, 0)

	diff := comparing.NewPeerDiffer(knownPeers)
	for _, p := range diff.UnknownPeers(reqP) {
		unknownPeers = append(unknownPeers, builder.ToPeerDigest(p))
	}
	for _, p := range diff.ProvidePeers(reqP) {
		newPeers = append(newPeers, builder.ToPeerDiscovered(p))
	}

	return &api.SynAck{
		Initiator:    req.Target,
		Target:       req.Initiator,
		NewPeers:     newPeers,
		UnknownPeers: unknownPeers,
	}, nil
}

//Acknowledge implements the protobuf Acknowledge request handler
func (h srvHandler) Acknowledge(ctx context.Context, req *api.AckRequest) (*empty.Empty, error) {
	builder := PeerBuilder{}

	//Store the peers requested
	for _, rp := range req.RequestedPeers {
		p := builder.FromPeerDiscovered(rp)
		h.notif.Notify(p)
		h.repo.SetKnownPeer(p)
	}

	return new(empty.Empty), nil
}

//NewServerHandler create a new GRPC server handler
func NewServerHandler(repo discovery.Repository, notif gossip.Notifier) api.DiscoveryServer {
	return srvHandler{
		repo:  repo,
		notif: notif,
	}
}

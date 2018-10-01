package rpc

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"

	api "github.com/uniris/uniris-core/autodiscovery/api/protobuf-spec"
	discovery "github.com/uniris/uniris-core/autodiscovery/pkg"
	"github.com/uniris/uniris-core/autodiscovery/pkg/gossip"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type gossipMessenger struct {
}

//SendSyn calls the Synchronize grpc method to retrieve unknown peers (SYN handshake)
func (g gossipMessenger) SendSyn(req gossip.SynRequest) (synAck *gossip.SynAck, err error) {
	serverAddr := fmt.Sprintf("%s", req.Receiver.Endpoint())
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	defer conn.Close()

	if err != nil {
		return nil, err
	}

	client := api.NewDiscoveryClient(conn)

	builder := PeerBuilder{}
	kp := make([]*api.PeerDigest, 0)
	for _, p := range req.KnownPeers {
		kp = append(kp, builder.ToPeerDigest(p))
	}

	res, err := client.Synchronize(context.Background(), &api.SynRequest{
		Initiator:  builder.ToPeerDigest(req.Initiator),
		Receiver:   builder.ToPeerDigest(req.Receiver),
		KnownPeers: kp,
	})
	if err != nil {
		statusCode, _ := status.FromError(err)
		if statusCode.Code() == codes.Unavailable {
			return nil, gossip.ErrUnrechablePeer
		}
		return nil, err
	}

	newP := make([]discovery.Peer, 0)
	unknown := make([]discovery.Peer, 0)
	for _, p := range res.NewPeers {
		newP = append(newP, builder.FromPeerDiscovered(p))
	}
	for _, p := range res.UnknownPeers {
		unknown = append(unknown, builder.FromPeerDigest(p))
	}

	return &gossip.SynAck{
		Initiator:    req.Receiver,
		Receiver:     req.Initiator,
		NewPeers:     newP,
		UnknownPeers: unknown,
	}, nil
}

//SendAck calls the Acknoweledge grpc method to send detailed peers requested
func (g gossipMessenger) SendAck(req gossip.AckRequest) error {
	serverAddr := fmt.Sprintf("%s", req.Receiver.Endpoint())
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	defer conn.Close()

	if err != nil {
		return err
	}

	builder := PeerBuilder{}

	client := api.NewDiscoveryClient(conn)

	reqP := make([]*api.PeerDiscovered, 0)
	for _, p := range req.RequestedPeers {
		reqP = append(reqP, builder.ToPeerDiscovered(p))
	}

	_, err = client.Acknowledge(context.Background(), &api.AckRequest{
		Initiator:      builder.ToPeerDigest(req.Initiator),
		Receiver:       builder.ToPeerDigest(req.Receiver),
		RequestedPeers: reqP,
	})
	if err != nil {
		return err
	}
	return nil
}

//NewMessenger creates a new gossip messenger using GRPC
func NewMessenger() gossip.Messenger {
	return gossipMessenger{}
}

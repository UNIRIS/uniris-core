package rpc

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"

	api "github.com/uniris/uniris-core/autodiscovery/api/protobuf-spec"
	"github.com/uniris/uniris-core/autodiscovery/pkg/gossip"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type gossipMessenger struct {
	domainFormat PeerDomainFormater
	apiFormat    PeerAPIFormater
}

//SendSyn calls the Synchronize grpc method to retrieve unknown peers (SYN handshake)
func (g gossipMessenger) SendSyn(req gossip.SynRequest) (synAck *gossip.SynAck, err error) {
	serverAddr := fmt.Sprintf("%s", req.Receiver.GetEndpoint())
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	defer conn.Close()

	if err != nil {
		return nil, err
	}

	client := api.NewDiscoveryClient(conn)
	resp, err := client.Synchronize(context.Background(), &api.SynRequest{
		Initiator:  g.apiFormat.BuildPeerDigest(req.Initiator),
		Receiver:   g.apiFormat.BuildPeerDigest(req.Receiver),
		KnownPeers: g.apiFormat.BuildPeerDigestCollection(req.KnownPeers),
	})
	if err != nil {
		statusCode, _ := status.FromError(err)
		if statusCode.Code() == codes.Unavailable {
			return nil, fmt.Errorf("Peer %s is unavailable", req.Receiver.GetEndpoint())
		}
		return nil, err
	}

	return &gossip.SynAck{
		Initiator:    req.Receiver,
		Receiver:     req.Initiator,
		NewPeers:     g.domainFormat.BuildPeerDetailedCollection(resp.NewPeers),
		UnknownPeers: g.domainFormat.BuildPeerDigestCollection(resp.UnknownPeers),
	}, nil
}

//SendAck calls the Acknoweledge grpc method to send detailed peers requested
func (g gossipMessenger) SendAck(req gossip.AckRequest) error {
	serverAddr := fmt.Sprintf("%s", req.Receiver.GetEndpoint())
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		return err
	}

	defer conn.Close()

	client := api.NewDiscoveryClient(conn)
	_, err = client.Acknowledge(context.Background(), &api.AckRequest{
		Initiator:      g.apiFormat.BuildPeerDigest(req.Initiator),
		Receiver:       g.apiFormat.BuildPeerDigest(req.Receiver),
		RequestedPeers: g.apiFormat.BuildPeerDetailedCollection(req.RequestedPeers),
	})
	if err != nil {
		return err
	}
	return nil
}

//NewMessenger creates a new gossip messenger using GRPC
func NewMessenger() gossip.Messenger {
	return gossipMessenger{
		apiFormat:    PeerAPIFormater{},
		domainFormat: PeerDomainFormater{},
	}
}

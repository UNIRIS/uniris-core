package rpc

import (
	"fmt"

	"golang.org/x/net/context"

	"google.golang.org/grpc/codes"

	api "github.com/uniris/uniris-core/autodiscovery/api/protobuf-spec"
	discovery "github.com/uniris/uniris-core/autodiscovery/pkg"
	"github.com/uniris/uniris-core/autodiscovery/pkg/gossip"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type clientMessenger struct{}

//SendSyn calls the Synchronize grpc method to retrieve unknown peers (SYN handshake)
func (m clientMessenger) SendSyn(req gossip.SynRequest) (synAck *gossip.SynAck, err error) {
	serverAddr := fmt.Sprintf("%s", req.Target.Endpoint())
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	defer conn.Close()

	if err != nil {
		return nil, err
	}

	//We initalize a GRPC client
	client := api.NewDiscoveryClient(conn)

	builder := PeerBuilder{}
	kp := make([]*api.PeerDigest, 0)
	for _, p := range req.KnownPeers {
		kp = append(kp, builder.ToPeerDigest(p))
	}

	res, err := client.Synchronize(context.Background(), &api.SynRequest{
		Initiator:  builder.ToPeerDigest(req.Initiator),
		Target:     builder.ToPeerDigest(req.Target),
		KnownPeers: kp,
	})

	if err != nil {
		//If the peer cannot be reached, we throw an ErrPeerUnreachable error
		statusCode, _ := status.FromError(err)
		if statusCode.Code() == codes.Unavailable {
			return nil, gossip.ErrUnreachablePeer
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
		Initiator:    req.Target,
		Target:       req.Initiator,
		NewPeers:     newP,
		UnknownPeers: unknown,
	}, nil
}

//SendAck calls the Acknoweledge grpc method to send detailed peers requested
func (m clientMessenger) SendAck(req gossip.AckRequest) error {
	serverAddr := fmt.Sprintf("%s", req.Target.Endpoint())
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	defer conn.Close()

	if err != nil {
		return err
	}

	builder := PeerBuilder{}

	//We initalize a GRPC client
	client := api.NewDiscoveryClient(conn)

	reqP := make([]*api.PeerDiscovered, 0)
	for _, p := range req.RequestedPeers {
		reqP = append(reqP, builder.ToPeerDiscovered(p))
	}

	_, err = client.Acknowledge(context.Background(), &api.AckRequest{
		Initiator:      builder.ToPeerDigest(req.Initiator),
		Target:         builder.ToPeerDigest(req.Target),
		RequestedPeers: reqP,
	})
	if err != nil {
		return err
	}
	return nil
}

//NewMessenger creates a new gossip messenger using GRPC
func NewMessenger() gossip.Messenger {
	return clientMessenger{}
}

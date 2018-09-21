package rpc

import (
	"fmt"
	"log"

	"golang.org/x/net/context"

	"google.golang.org/grpc/codes"

	api "github.com/uniris/uniris-core/autodiscovery/api/protobuf-spec"
	discovery "github.com/uniris/uniris-core/autodiscovery/pkg"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type gossipMessenger struct {
	domainFormat PeerDomainFormater
	apiFormat    PeerAPIFormater
}

//SendSyn calls the Synchronize grpc method to retrieve unknown peers (SYN handshake)
func (g gossipMessenger) SendSyn(req discovery.SynRequest) (synAck *discovery.SynAck, err error) {
	serverAddr := fmt.Sprintf("%s", req.Receiver.GetEndpoint())
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	defer conn.Close()

	if err != nil {
		return nil, err
	}

	//We initalize a GRPC client
	client := api.NewDiscoveryClient(conn)

	//We send the SYN request
	resp, err := client.Synchronize(context.Background(), &api.SynRequest{
		Initiator:  g.apiFormat.BuildPeerDigest(req.Initiator),
		Receiver:   g.apiFormat.BuildPeerDigest(req.Receiver),
		KnownPeers: g.apiFormat.BuildPeerDigestCollection(req.KnownPeers),
	})

	if err != nil {
		//If the peer cannot be reached, we throw an ErrPeerUnreachable error
		statusCode, _ := status.FromError(err)
		if statusCode.Code() == codes.Unavailable {
			log.Printf(discovery.ErrPeerUnreachable.Error(), req.Receiver.GetEndpoint())
			return nil, discovery.ErrPeerUnreachable
		}
		return nil, err
	}

	//We format the SYN ACK reponse for the domain
	return &discovery.SynAck{
		Initiator:    req.Receiver,
		Receiver:     req.Initiator,
		NewPeers:     g.domainFormat.BuildPeerDetailedCollection(resp.NewPeers),
		UnknownPeers: g.domainFormat.BuildPeerDigestCollection(resp.UnknownPeers),
	}, nil
}

//SendAck calls the Acknoweledge grpc method to send detailed peers requested
func (g gossipMessenger) SendAck(req discovery.AckRequest) error {
	serverAddr := fmt.Sprintf("%s", req.Receiver.GetEndpoint())
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		return err
	}

	defer conn.Close()

	//We initalize a GRPC client
	client := api.NewDiscoveryClient(conn)

	//We send the ACK request with the requested detailed peers
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
func NewMessenger() discovery.GossipCycleMessenger {
	return gossipMessenger{
		apiFormat:    PeerAPIFormater{},
		domainFormat: PeerDomainFormater{},
	}
}

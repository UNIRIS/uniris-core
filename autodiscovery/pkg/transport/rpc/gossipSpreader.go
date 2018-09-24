package rpc

import (
	"fmt"
	"log"

	"golang.org/x/net/context"

	"google.golang.org/grpc/codes"

	api "github.com/uniris/uniris-core/autodiscovery/api/protobuf-spec"
	discovery "github.com/uniris/uniris-core/autodiscovery/pkg"
	"github.com/uniris/uniris-core/autodiscovery/pkg/gossip"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type gossipSpreader struct {
	domainFormat PeerDomainFormater
	apiFormat    PeerAPIFormater
}

//SendSyn calls the Synchronize grpc method to retrieve unknown peers (SYN handshake)
func (g gossipSpreader) SendSyn(req discovery.SynRequest) (synAck *discovery.SynAck, err error) {
	serverAddr := fmt.Sprintf("%s", req.Target.Endpoint())
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
		Target:     g.apiFormat.BuildPeerDigest(req.Target),
		KnownPeers: g.apiFormat.BuildPeerDigestCollection(req.KnownPeers),
	})

	if err != nil {
		//If the peer cannot be reached, we throw an ErrPeerUnreachable error
		statusCode, _ := status.FromError(err)
		if statusCode.Code() == codes.Unavailable {
			log.Printf(gossip.ErrPeerUnreachable.Error(), req.Target.Endpoint())
			return nil, gossip.ErrPeerUnreachable
		}
		return nil, err
	}

	//We format the SYN ACK reponse for the domain
	return &discovery.SynAck{
		Initiator:    req.Target,
		Target:       req.Initiator,
		NewPeers:     g.domainFormat.BuildPeerDetailedCollection(resp.NewPeers),
		UnknownPeers: g.domainFormat.BuildPeerDigestCollection(resp.UnknownPeers),
	}, nil
}

//SendAck calls the Acknoweledge grpc method to send detailed peers requested
func (g gossipSpreader) SendAck(req discovery.AckRequest) error {
	serverAddr := fmt.Sprintf("%s", req.Target.Endpoint())
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
		Target:         g.apiFormat.BuildPeerDigest(req.Target),
		RequestedPeers: g.apiFormat.BuildPeerDetailedCollection(req.RequestedPeers),
	})
	if err != nil {
		return err
	}
	return nil
}

//NewMessenger creates a new gossip messenger using GRPC
func NewMessenger() discovery.GossipSpreader {
	return gossipSpreader{
		apiFormat:    PeerAPIFormater{},
		domainFormat: PeerDomainFormater{},
	}
}

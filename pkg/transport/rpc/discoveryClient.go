package rpc

import (
	"context"
	"fmt"
	"time"

	api "github.com/uniris/uniris-core/api/protobuf-spec"
	"github.com/uniris/uniris-core/pkg/discovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type rndMsg struct{}

//NewDiscoveryClient creates a new discovery client with GRPC
func NewDiscoveryClient() discovery.Client {
	return rndMsg{}
}

func (m rndMsg) SendSyn(target discovery.PeerIdentity, known []discovery.Peer) (unknown []discovery.Peer, new []discovery.Peer, err error) {
	serverAddr := fmt.Sprintf("%s:%d", target.IP().String(), target.Port())
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		return nil, nil, nil
	}
	defer conn.Close()

	kp := make([]*api.PeerDigest, 0)
	for _, p := range known {
		kp = append(kp, formatPeerDigestAPI(p))
	}

	client := api.NewDiscoveryServiceClient(conn)
	res, err := client.Synchronize(context.Background(), &api.SynRequest{
		KnownPeers: kp,
		Timestamp:  time.Now().Unix(),
	})
	if err != nil {
		//If the peer cannot be reached, we throw an ErrPeerUnreachable error
		statusCode, _ := status.FromError(err)
		if statusCode.Code() == codes.Unavailable {
			return nil, nil, discovery.ErrUnreachablePeer
		}
		return nil, nil, err
	}

	fmt.Printf("SYNC RESPONSE - %s\n", time.Unix(res.Timestamp, 0).String())

	for _, p := range res.NewPeers {
		new = append(new, formatPeerDiscovered(p))
	}
	for _, p := range res.UnknownPeers {
		unknown = append(unknown, formatPeerDigest(p))
	}

	return
}

func (m rndMsg) SendAck(target discovery.PeerIdentity, requested []discovery.Peer) error {
	serverAddr := fmt.Sprintf("%s:%d", target.IP().String(), target.Port())
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	defer conn.Close()

	if err != nil {
		return err
	}

	client := api.NewDiscoveryServiceClient(conn)

	reqP := make([]*api.PeerDiscovered, 0)
	for _, p := range requested {
		reqP = append(reqP, formatPeerDiscoveredAPI(p))
	}

	res, err := client.Acknowledge(context.Background(), &api.AckRequest{
		RequestedPeers: reqP,
	})
	if err != nil {
		//If the peer cannot be reached, we throw an ErrPeerUnreachable error
		statusCode, _ := status.FromError(err)
		if statusCode.Code() == codes.Unavailable {
			return discovery.ErrUnreachablePeer
		}
		return err
	}

	fmt.Printf("ACK RESPONSE - %s\n", time.Unix(res.Timestamp, 0).String())
	return nil
}

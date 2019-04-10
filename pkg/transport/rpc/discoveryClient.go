package rpc

import (
	"context"
	"fmt"
	"time"

	api "github.com/uniris/uniris-core/api/protobuf-spec"
	"github.com/uniris/uniris-core/pkg/crypto"
	"github.com/uniris/uniris-core/pkg/discovery"
	"github.com/uniris/uniris-core/pkg/logging"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type rndMsg struct {
	logger     logging.Logger
	publickey  crypto.PublicKey
	privatekey crypto.PrivateKey
}

//NewGossipRoundMessenger creates a new round messenger with GRPC
func NewGossipRoundMessenger(l logging.Logger, pubK crypto.PublicKey, privK crypto.PrivateKey) discovery.Messenger {
	return rndMsg{
		logger:     l,
		publickey:  pubK,
		privatekey: privK,
	}
}

func (m rndMsg) SendSyn(target discovery.PeerIdentity, known []discovery.Peer) (requestedPeers []discovery.PeerIdentity, discoveredPeers []discovery.Peer, err error) {
	serverAddr := fmt.Sprintf("%s:%d", target.IP().String(), target.Port())
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		return nil, nil, nil
	}
	defer conn.Close()

	kp := make([]*api.PeerDigest, 0)
	for _, p := range known {
		fp, err := formatPeerDigestAPI(p)
		if err != nil {
			return nil, nil, err
		}
		kp = append(kp, fp)
	}

	pk, err := m.publickey.Marshal()
	if err != nil {
		return nil, nil, err
	}
	sig, err := m.privatekey.Sign(pk)
	if err != nil {
		return nil, nil, err
	}

	client := api.NewDiscoveryServiceClient(conn)
	res, err := client.Synchronize(context.Background(), &api.SynRequest{
		KnownPeers: kp,
		Timestamp:  time.Now().Unix(),
		PublicKey:  pk,
		Signature:  sig,
	})
	if err != nil {
		//If the peer cannot be reached, we throw an ErrPeerUnreachable error
		statusCode, _ := status.FromError(err)
		if statusCode.Code() == codes.Unavailable {
			return nil, nil, discovery.ErrUnreachablePeer
		}
		return nil, nil, err
	}

	m.logger.Debug("SYN RESPONSE - " + time.Unix(res.Timestamp, 0).String())

	for _, p := range res.DiscoveredPeers {
		fp, err := formatPeerDiscovered(p)
		if err != nil {
			return nil, nil, err
		}
		discoveredPeers = append(discoveredPeers, fp)
	}
	for _, p := range res.RequestedPeers {
		fp, err := formatPeerIdentity(p)
		if err != nil {
			return nil, nil, err
		}
		requestedPeers = append(requestedPeers, fp)
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
		fp, err := formatPeerDiscoveredAPI(p)
		if err != nil {
			return err
		}
		reqP = append(reqP, fp)
	}

	pk, err := m.publickey.Marshal()
	if err != nil {
		return err
	}
	sig, err := m.privatekey.Sign(pk)
	if err != nil {
		return err
	}

	res, err := client.Acknowledge(context.Background(), &api.AckRequest{
		RequestedPeers: reqP,
		PublicKey:      pk,
		Signature:      sig,
	})
	if err != nil {
		//If the peer cannot be reached, we throw an ErrPeerUnreachable error
		statusCode, _ := status.FromError(err)
		if statusCode.Code() == codes.Unavailable {
			return discovery.ErrUnreachablePeer
		}
		return err
	}

	m.logger.Debug("ACK RESPONSE - " + time.Unix(res.Timestamp, 0).String())
	return nil
}

package rpc

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	api "github.com/uniris/uniris-core/api/protobuf-spec"
	"github.com/uniris/uniris-core/pkg/discovery"
	"google.golang.org/grpc"
)

/*
Scenario: Send a syn request to a peer by sending its local view
	Given a local view
	When I want to send it to a peer
	Then I get the unknown from this peers and the peers I don't known
*/
func TestSendSynRequest(t *testing.T) {

	repo := &mockDiscoveryRepo{}
	service := discovery.NewService(repo, nil, &mockDiscoveryNotifier{}, mockPeerNetworker{}, mockPeerInfo{})

	lis, _ := net.Listen("tcp", ":3000")
	defer lis.Close()
	grpcServer := grpc.NewServer()
	api.RegisterDiscoveryServiceServer(grpcServer, NewDiscoveryServer(service))
	go grpcServer.Serve(lis)

	cli := NewDiscoveryClient()
	target := discovery.NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "pubkey")
	unknown, news, err := cli.SendSyn(target, []discovery.Peer{
		discovery.NewDiscoveredPeer(
			discovery.NewPeerIdentity(net.ParseIP("10.0.0.1"), 3000, "pubkey2"),
			discovery.NewPeerHeartbeatState(time.Now(), 1000),
			discovery.NewPeerAppState("1.0", discovery.OkPeerStatus, 10.0, 20.0, "", 300, 1, 1000),
		),
	})
	assert.Nil(t, err)
	assert.NotEmpty(t, unknown)
	assert.Empty(t, news)
	assert.Equal(t, "pubkey2", unknown[0].Identity().PublicKey())
}

/*
Scenario: Send a syn request to a disconnected peer
	Given a local view
	When I want to send it to a disconnected peer
	Then I get an error as unreachable
*/
func TestSendSynRequestUnreach(t *testing.T) {
	cli := NewDiscoveryClient()
	target := discovery.NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "pubkey")
	_, _, err := cli.SendSyn(target, []discovery.Peer{
		discovery.NewDiscoveredPeer(
			discovery.NewPeerIdentity(net.ParseIP("10.0.0.1"), 3000, "pubkey2"),
			discovery.NewPeerHeartbeatState(time.Now(), 1000),
			discovery.NewPeerAppState("1.0", discovery.OkPeerStatus, 10.0, 20.0, "", 300, 1, 1000),
		),
	})
	assert.Equal(t, err, discovery.ErrUnreachablePeer)
}

/*
Scenario: Send a ack request to a peer by sending details about the requested peers
	Given details of the requested peers
	When I want to send it to a peer
	Then I get not error and the sended peer is stored
*/
func TestSendAckRequest(t *testing.T) {

	repo := &mockDiscoveryRepo{}
	service := discovery.NewService(repo, nil, &mockDiscoveryNotifier{}, mockPeerNetworker{}, mockPeerInfo{})

	lis, _ := net.Listen("tcp", ":3000")
	defer lis.Close()
	grpcServer := grpc.NewServer()
	api.RegisterDiscoveryServiceServer(grpcServer, NewDiscoveryServer(service))
	go grpcServer.Serve(lis)

	cli := NewDiscoveryClient()
	target := discovery.NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "pubkey")
	err := cli.SendAck(target, []discovery.Peer{
		discovery.NewDiscoveredPeer(
			discovery.NewPeerIdentity(net.ParseIP("10.0.0.1"), 3000, "pubkey2"),
			discovery.NewPeerHeartbeatState(time.Now(), 1000),
			discovery.NewPeerAppState("1.0", discovery.OkPeerStatus, 10.0, 20.0, "", 300, 1, 1000),
		),
	})
	assert.Nil(t, err)
	assert.NotEmpty(t, repo.knownPeers)
	assert.Equal(t, "pubkey2", repo.knownPeers[0].Identity().PublicKey())
}

/*
Scenario: Send a ack request to a disconnected peer
	Given details of the requested peers
	When I want to send it to a disconnected peer
	Then I get an error as unreachable
*/
func TestSendAckRequestUnreach(t *testing.T) {
	cli := NewDiscoveryClient()
	target := discovery.NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "pubkey")
	err := cli.SendAck(target, []discovery.Peer{
		discovery.NewDiscoveredPeer(
			discovery.NewPeerIdentity(net.ParseIP("10.0.0.1"), 3000, "pubkey2"),
			discovery.NewPeerHeartbeatState(time.Now(), 1000),
			discovery.NewPeerAppState("1.0", discovery.OkPeerStatus, 10.0, 20.0, "", 300, 1, 1000),
		),
	})
	assert.Equal(t, err, discovery.ErrUnreachablePeer)
}

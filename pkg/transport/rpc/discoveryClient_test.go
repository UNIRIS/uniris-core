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

	db := &mockDiscoveryDB{}

	lis, _ := net.Listen("tcp", ":3000")
	defer lis.Close()
	grpcServer := grpc.NewServer()
	discoverySrv := NewDiscoveryServer(db, &mockDiscoveryNotifier{})
	api.RegisterDiscoveryServiceServer(grpcServer, discoverySrv)
	go grpcServer.Serve(lis)

	rndMsg := NewGossipRoundMessenger()
	target := discovery.NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "pubkey")
	reqPeers, discoveries, err := rndMsg.SendSyn(target, []discovery.Peer{
		discovery.NewDiscoveredPeer(
			discovery.NewPeerIdentity(net.ParseIP("10.0.0.1"), 3000, "pubkey2"),
			discovery.NewPeerHeartbeatState(time.Now(), 1000),
			discovery.NewPeerAppState("1.0", discovery.OkPeerStatus, 10.0, 20.0, "", 300, 1, 1000),
		),
	})
	assert.Nil(t, err)
	assert.NotEmpty(t, reqPeers)
	assert.Empty(t, discoveries)
	assert.Equal(t, "pubkey2", reqPeers[0].PublicKey())
}

/*
Scenario: Send a syn request to a disconnected peer
	Given a local view
	When I want to send it to a disconnected peer
	Then I get an error as unreachable
*/
func TestSendSynRequestUnreach(t *testing.T) {
	rndMsg := NewGossipRoundMessenger()
	target := discovery.NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "pubkey")
	_, _, err := rndMsg.SendSyn(target, []discovery.Peer{
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

	db := &mockDiscoveryDB{}

	lis, _ := net.Listen("tcp", ":3000")
	defer lis.Close()
	grpcServer := grpc.NewServer()
	discoverySrv := NewDiscoveryServer(db, &mockDiscoveryNotifier{})
	api.RegisterDiscoveryServiceServer(grpcServer, discoverySrv)
	go grpcServer.Serve(lis)

	rndMsg := NewGossipRoundMessenger()
	target := discovery.NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "pubkey")
	err := rndMsg.SendAck(target, []discovery.Peer{
		discovery.NewDiscoveredPeer(
			discovery.NewPeerIdentity(net.ParseIP("10.0.0.1"), 3000, "pubkey2"),
			discovery.NewPeerHeartbeatState(time.Now(), 1000),
			discovery.NewPeerAppState("1.0", discovery.OkPeerStatus, 10.0, 20.0, "", 300, 1, 1000),
		),
	})
	assert.Nil(t, err)
	assert.NotEmpty(t, db.discoveredPeers)
	assert.Equal(t, "pubkey2", db.discoveredPeers[0].Identity().PublicKey())
}

/*
Scenario: Send a ack request to a disconnected peer
	Given details of the requested peers
	When I want to send it to a disconnected peer
	Then I get an error as unreachable
*/
func TestSendAckRequestUnreach(t *testing.T) {
	rndMsg := NewGossipRoundMessenger()
	target := discovery.NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "pubkey")
	err := rndMsg.SendAck(target, []discovery.Peer{
		discovery.NewDiscoveredPeer(
			discovery.NewPeerIdentity(net.ParseIP("10.0.0.1"), 3000, "pubkey2"),
			discovery.NewPeerHeartbeatState(time.Now(), 1000),
			discovery.NewPeerAppState("1.0", discovery.OkPeerStatus, 10.0, 20.0, "", 300, 1, 1000),
		),
	})
	assert.Equal(t, err, discovery.ErrUnreachablePeer)
}

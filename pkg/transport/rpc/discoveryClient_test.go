package rpc

import (
	"crypto/rand"
	"log"
	"net"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	api "github.com/uniris/uniris-core/api/protobuf-spec"
	"github.com/uniris/uniris-core/pkg/crypto"
	"github.com/uniris/uniris-core/pkg/discovery"
	"github.com/uniris/uniris-core/pkg/logging"

	"google.golang.org/grpc"
)

/*
Scenario: Send a syn request to a peer by sending its local view
	Given a local view
	When I want to send it to a peer
	Then I get the unknown from this peers and the peers I don't known
*/
func TestSendSynRequest(t *testing.T) {
	priv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
	_, pub1, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
	_, pub2, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

	sharedDB := mochSharedDB{}
	sharedDB.WriteAuthorizedNode(pub)
	sharedDB.WriteAuthorizedNode(pub1)
	sharedDB.WriteAuthorizedNode(pub2)

	db := &mockDiscoveryDB{}
	l := logging.NewLogger("stdout", log.New(os.Stdout, "", 0), "test", net.ParseIP("127.0.0.1"), logging.ErrorLogLevel)
	lis, _ := net.Listen("tcp", ":3000")
	defer lis.Close()
	grpcServer := grpc.NewServer()
	discoverySrv := NewDiscoveryServer(db, &mockDiscoveryNotifier{}, l, pub, priv, sharedDB)
	api.RegisterDiscoveryServiceServer(grpcServer, discoverySrv)
	go grpcServer.Serve(lis)

	rndMsg := NewGossipRoundMessenger(l, pub, priv)
	target := discovery.NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, pub1)
	reqPeers, discoveries, err := rndMsg.SendSyn(target, []discovery.Peer{
		discovery.NewDiscoveredPeer(
			discovery.NewPeerIdentity(net.ParseIP("10.0.0.1"), 3000, pub2),
			discovery.NewPeerHeartbeatState(time.Now(), 1000),
			discovery.NewPeerAppState("1.0", discovery.OkPeerStatus, 10.0, 20.0, "", 300, 1, 1000),
		),
	})
	assert.Nil(t, err)
	assert.NotEmpty(t, reqPeers)
	assert.Empty(t, discoveries)
	assert.Equal(t, pub2, reqPeers[0].PublicKey())
}

/*
Scenario: Send a syn request to a disconnected peer
	Given a local view
	When I want to send it to a disconnected peer
	Then I get an error as unreachable
*/
func TestSendSynRequestUnreach(t *testing.T) {
	priv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
	_, pub2, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
	_, pub3, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

	l := logging.NewLogger("stdout", log.New(os.Stdout, "", 0), "test", net.ParseIP("127.0.0.1"), logging.ErrorLogLevel)
	rndMsg := NewGossipRoundMessenger(l, pub, priv)
	target := discovery.NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, pub2)
	_, _, err := rndMsg.SendSyn(target, []discovery.Peer{
		discovery.NewDiscoveredPeer(
			discovery.NewPeerIdentity(net.ParseIP("10.0.0.1"), 3000, pub3),
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
	sharedDB := mochSharedDB{}

	l := logging.NewLogger("stdout", log.New(os.Stdout, "", 0), "test", net.ParseIP("127.0.0.1"), logging.ErrorLogLevel)

	priv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
	_, pub2, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
	_, pub3, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

	sharedDB.WriteAuthorizedNode(pub)
	sharedDB.WriteAuthorizedNode(pub2)
	sharedDB.WriteAuthorizedNode(pub3)

	lis, _ := net.Listen("tcp", ":3000")
	defer lis.Close()
	grpcServer := grpc.NewServer()
	discoverySrv := NewDiscoveryServer(db, &mockDiscoveryNotifier{}, l, pub, priv, sharedDB)
	api.RegisterDiscoveryServiceServer(grpcServer, discoverySrv)
	go grpcServer.Serve(lis)

	rndMsg := NewGossipRoundMessenger(l, pub, priv)
	target := discovery.NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, pub2)
	err := rndMsg.SendAck(target, []discovery.Peer{
		discovery.NewDiscoveredPeer(
			discovery.NewPeerIdentity(net.ParseIP("10.0.0.1"), 3000, pub3),
			discovery.NewPeerHeartbeatState(time.Now(), 1000),
			discovery.NewPeerAppState("1.0", discovery.OkPeerStatus, 10.0, 20.0, "", 300, 1, 1000),
		),
	})
	assert.Nil(t, err)
	assert.NotEmpty(t, db.discoveredPeers)
	assert.Equal(t, true, db.discoveredPeers[0].Identity().PublicKey().Equals(pub3))
}

/*
Scenario: Send a ack request to a disconnected peer
	Given details of the requested peers
	When I want to send it to a disconnected peer
	Then I get an error as unreachable
*/
func TestSendAckRequestUnreach(t *testing.T) {
	priv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
	_, pub2, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
	_, pub3, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

	l := logging.NewLogger("stdout", log.New(os.Stdout, "", 0), "test", net.ParseIP("127.0.0.1"), logging.ErrorLogLevel)
	rndMsg := NewGossipRoundMessenger(l, pub, priv)
	target := discovery.NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, pub2)
	err := rndMsg.SendAck(target, []discovery.Peer{
		discovery.NewDiscoveredPeer(
			discovery.NewPeerIdentity(net.ParseIP("10.0.0.1"), 3000, pub3),
			discovery.NewPeerHeartbeatState(time.Now(), 1000),
			discovery.NewPeerAppState("1.0", discovery.OkPeerStatus, 10.0, 20.0, "", 300, 1, 1000),
		),
	})
	assert.Equal(t, err, discovery.ErrUnreachablePeer)
}

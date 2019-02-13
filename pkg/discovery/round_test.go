package discovery

import (
	"errors"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

/*
Scenario: Spread a gossip round and discover peers
	Given a initiator peer, a receiver peer and list of known peers
	When we start a gossip round we spread what we know
	Then we get the new peers discovered
*/
func TestRoundRunning(t *testing.T) {

	target := NewPeerIdentity(net.ParseIP("20.100.4.120"), 3000, "key2")

	p1 := NewDiscoveredPeer(
		NewPeerIdentity(net.ParseIP("50.20.100.2"), 3000, "key3"),
		NewPeerHeartbeatState(time.Now(), 0),
		NewPeerAppState("1.0", OkPeerStatus, 10.0, 20.0, "", 0, 1, 0),
	)

	p2 := NewDiscoveredPeer(
		NewPeerIdentity(net.ParseIP("50.10.30.2"), 3000, "uKey1"),
		NewPeerHeartbeatState(time.Now(), 0),
		NewPeerAppState("1.0", OkPeerStatus, 20.0, 19.4, "", 0, 1, 0),
	)

	r := round{
		peers:  []Peer{p1, p2},
		target: target,
	}

	discoveries, err := r.run(mockClient{})
	assert.Nil(t, err)
	assert.Len(t, discoveries, 1)
	assert.Equal(t, "dKey1", discoveries[0].Identity().PublicKey())
}

/*
Scenario: Spread gossip but unreach the target peer during the SYN request
	Given a initiator peer, a receiver peer and list of known peers
	When are sending the SYN, the target cannot be reached
	Then we get the unreached peer
*/
func TestRoundRunningWithUnreachWhenSYN(t *testing.T) {

	target := NewPeerIdentity(net.ParseIP("20.100.4.120"), 3000, "key2")

	p1 := NewDiscoveredPeer(
		NewPeerIdentity(net.ParseIP("50.20.100.2"), 3000, "key3"),
		NewPeerHeartbeatState(time.Now(), 0),
		NewPeerAppState("1.0", OkPeerStatus, 30.0, 10.0, "", 0, 1, 0),
	)

	p2 := NewDiscoveredPeer(
		NewPeerIdentity(net.ParseIP("50.10.30.2"), 3000, "uKey1"),
		NewPeerHeartbeatState(time.Now(), 0),
		NewPeerAppState("1.0", OkPeerStatus, 30.0, 10.0, "", 0, 1, 0),
	)

	r := round{
		peers:  []Peer{p1, p2},
		target: target,
	}

	_, err := r.run(mockClientWithSynFailure{})
	assert.Equal(t, err, ErrUnreachablePeer)
}

/*
Scenario: Spread gossip but unreach the target peer during the SYN request
	Given a initiator peer, a receiver peer and list of known peers
	When are sending the SYN, the target cannot be reached
	Then we get the unreached peer
*/
func TestRoundRunningWithUnreachWhenACK(t *testing.T) {

	target := NewPeerIdentity(net.ParseIP("20.100.4.120"), 3000, "key2")

	p1 := NewDiscoveredPeer(
		NewPeerIdentity(net.ParseIP("50.20.100.2"), 3000, "key3"),
		NewPeerHeartbeatState(time.Now(), 0),
		NewPeerAppState("1.0", OkPeerStatus, 30.0, 10.0, "", 0, 1, 0),
	)

	p2 := NewDiscoveredPeer(
		NewPeerIdentity(net.ParseIP("50.10.30.2"), 3000, "uKey1"),
		NewPeerHeartbeatState(time.Now(), 0),
		NewPeerAppState("1.0", OkPeerStatus, 30.0, 10.0, "", 0, 1, 0),
	)

	r := round{
		peers:  []Peer{p1, p2},
		target: target,
	}

	_, err := r.run(mockClientWithAckFailure{})
	assert.Equal(t, err, ErrUnreachablePeer)
}

type mockClientWithSynFailure struct {
}

func (m mockClientWithSynFailure) SendSyn(target PeerIdentity, known []Peer) (unknown []Peer, new []Peer, err error) {
	return nil, nil, ErrUnreachablePeer
}

func (m mockClientWithSynFailure) SendAck(target PeerIdentity, requested []Peer) error {
	return nil
}

type mockClientWithAckFailure struct {
}

func (m mockClientWithAckFailure) SendSyn(target PeerIdentity, known []Peer) (unknown []Peer, new []Peer, err error) {
	tar := NewSelfPeer("uKey1", net.ParseIP("200.18.186.39"), 3000, "1.1", 30.0, 10.0)

	hb := NewPeerHeartbeatState(time.Now(), 0)
	as := NewPeerAppState("1.0", OkPeerStatus, 30.0, 10.0, "", 0, 1, 0)

	np1 := NewDiscoveredPeer(
		NewPeerIdentity(net.ParseIP("35.200.100.2"), 3000, "dKey1"),
		hb, as,
	)

	newPeers := []Peer{np1}

	unknownPeers := []Peer{tar}

	return unknownPeers, newPeers, nil
}

func (m mockClientWithAckFailure) SendAck(target PeerIdentity, requested []Peer) error {
	return ErrUnreachablePeer
}

type mockClientUnexpectedFailure struct {
}

func (m mockClientUnexpectedFailure) SendSyn(target PeerIdentity, known []Peer) (unknown []Peer, new []Peer, err error) {
	return nil, nil, errors.New("Unexpected")
}

func (m mockClientUnexpectedFailure) SendAck(target PeerIdentity, requested []Peer) error {
	return nil
}

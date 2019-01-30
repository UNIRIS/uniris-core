package discovery

import (
	"errors"
	"net"
	"sync"
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
func TestSpreadDiscoveries(t *testing.T) {

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

	r := round{target, mockClient{}}

	kp := []Peer{p1, p2}

	discoveries := make(chan Peer)
	reaches := make(chan PeerIdentity)

	var wg sync.WaitGroup
	wg.Add(2)

	go r.run(kp, discoveries, reaches, nil)

	pp := make([]Peer, 0)
	go func() {
		for p := range discoveries {
			pp = append(pp, p)
			wg.Done()
			close(discoveries)
		}
	}()

	reachP := make([]PeerIdentity, 0)
	go func() {
		for p := range reaches {
			reachP = append(reachP, p)
			wg.Done()
			close(reaches)
		}
	}()

	wg.Wait()

	assert.NotEmpty(t, pp)
	assert.Equal(t, 1, len(pp))
	assert.Equal(t, "dKey1", pp[0].Identity().PublicKey())

	assert.NotEmpty(t, reachP)
}

/*
Scenario: Spread gossip but unreach the target peer during the SYN request
	Given a initiator peer, a receiver peer and list of known peers
	When are sending the SYN, the target cannot be reached
	Then we get the unreached peer
*/
func TestSYNSpreadUnreachables(t *testing.T) {

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

	r := round{target, mockClientWithSynFailure{}}

	kp := []Peer{p1, p2}

	unreaches := make(chan PeerIdentity)

	go func() {
		err := r.run(kp, nil, nil, unreaches)
		assert.Nil(t, err)
		close(unreaches)
	}()

	pp := make([]PeerIdentity, 0)
	for p := range unreaches {
		pp = append(pp, p)
	}

	assert.NotEmpty(t, pp)
	assert.Equal(t, 1, len(pp))
	assert.Equal(t, target.PublicKey(), pp[0].PublicKey())
}

/*
Scenario: Spread gossip but unreach the target peer during the SYN request
	Given a initiator peer, a receiver peer and list of known peers
	When are sending the SYN, the target cannot be reached
	Then we get the unreached peer
*/
func TestACKSpreadUnreachables(t *testing.T) {

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

	r := round{target, mockClientWithAckFailure{}}

	kp := []Peer{p1, p2}

	discoveries := make(chan Peer)
	reaches := make(chan PeerIdentity)
	unreaches := make(chan PeerIdentity)

	go func() {
		err := r.run(kp, discoveries, reaches, unreaches)
		assert.Nil(t, err)
		close(unreaches)
		close(discoveries)
		close(reaches)
	}()

	go func() {
		for range reaches {
		}
	}()

	unreachP := make([]PeerIdentity, 0)
	for p := range unreaches {
		unreachP = append(unreachP, p)
	}

	discovP := make([]Peer, 0)
	for p := range discoveries {
		discovP = append(discovP, p)
	}

	assert.NotEmpty(t, unreachP)
	assert.Equal(t, 1, len(unreachP))
	assert.Equal(t, target.PublicKey(), unreachP[0].PublicKey())

	assert.Empty(t, discovP)
}

/*
Scenario: Gossip but gets an unexpected error
	Given a initiator peer, a receiver peer and list of known peers
	When are gossiping, we get an unexpected error
	Then we get retrieve the error
*/
func TestSpreadUnexpectedError(t *testing.T) {

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

	r := round{target, mockClientUnexpectedFailure{}}

	kp := []Peer{p1, p2}

	err := r.run(kp, nil, nil, nil)
	assert.NotNil(t, err)
	assert.Error(t, err, "Unexpected")
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
	tar := NewLocalPeer("uKey1", net.ParseIP("200.18.186.39"), 3000, "1.1", 30.0, 10.0)

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

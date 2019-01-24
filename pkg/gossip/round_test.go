package gossip

import (
	"errors"
	"net"
	"sync"
	"testing"
	"time"

	uniris "github.com/uniris/uniris-core/pkg"

	"github.com/stretchr/testify/assert"
)

/*
Scenario: Spread a gossip round and discover peers
	Given a initiator peer, a receiver peer and list of known peers
	When we start a gossip round we spread what we know
	Then we get the new peers discovered
*/
func TestSpreadDiscoveries(t *testing.T) {
	initP := uniris.NewLocalPeer("key", net.ParseIP("127.0.0.1"), 3000, "1.0", 30.0, 10.0)

	target := uniris.NewPeerDigest(
		uniris.NewPeerIdentity(net.ParseIP("20.100.4.120"), 3000, "key2"),
		uniris.NewPeerHeartbeatState(time.Now(), 0))

	p1 := uniris.NewDiscoveredPeer(
		uniris.NewPeerIdentity(net.ParseIP("50.20.100.2"), 3000, "key3"),
		uniris.NewPeerHeartbeatState(time.Now(), 0),
		uniris.NewPeerAppState("1.0", uniris.OkPeerStatus, 10.0, 20.0, "", 0, 1, 0),
	)

	p2 := uniris.NewDiscoveredPeer(
		uniris.NewPeerIdentity(net.ParseIP("50.10.30.2"), 3000, "uKey1"),
		uniris.NewPeerHeartbeatState(time.Now(), 0),
		uniris.NewPeerAppState("1.0", uniris.OkPeerStatus, 20.0, 19.4, "", 0, 1, 0),
	)

	r := round{initP, target, mockMessenger{}}

	kp := []uniris.Peer{p1, p2}

	discoveries := make(chan uniris.Peer)
	reaches := make(chan uniris.Peer)

	var wg sync.WaitGroup
	wg.Add(2)

	go r.run(kp, discoveries, reaches, nil)

	pp := make([]uniris.Peer, 0)
	go func() {
		for p := range discoveries {
			pp = append(pp, p)
			wg.Done()
			close(discoveries)
		}
	}()

	reachP := make([]uniris.Peer, 0)
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
	initP := uniris.NewLocalPeer("key", net.ParseIP("127.0.0.1"), 3000, "1.0", 30.0, 10.0)

	target := uniris.NewPeerDigest(
		uniris.NewPeerIdentity(net.ParseIP("20.100.4.120"), 3000, "key2"),
		uniris.NewPeerHeartbeatState(time.Now(), 0))

	p1 := uniris.NewDiscoveredPeer(
		uniris.NewPeerIdentity(net.ParseIP("50.20.100.2"), 3000, "key3"),
		uniris.NewPeerHeartbeatState(time.Now(), 0),
		uniris.NewPeerAppState("1.0", uniris.OkPeerStatus, 30.0, 10.0, "", 0, 1, 0),
	)

	p2 := uniris.NewDiscoveredPeer(
		uniris.NewPeerIdentity(net.ParseIP("50.10.30.2"), 3000, "uKey1"),
		uniris.NewPeerHeartbeatState(time.Now(), 0),
		uniris.NewPeerAppState("1.0", uniris.OkPeerStatus, 30.0, 10.0, "", 0, 1, 0),
	)

	r := round{initP, target, mockMessengerWithSynFailure{}}

	kp := []uniris.Peer{p1, p2}

	unreaches := make(chan uniris.Peer)

	go func() {
		err := r.run(kp, nil, nil, unreaches)
		assert.Nil(t, err)
		close(unreaches)
	}()

	pp := make([]uniris.Peer, 0)
	for p := range unreaches {
		pp = append(pp, p)
	}

	assert.NotEmpty(t, pp)
	assert.Equal(t, 1, len(pp))
	assert.Equal(t, target.Identity().PublicKey(), pp[0].Identity().PublicKey())
}

/*
Scenario: Spread gossip but unreach the target peer during the SYN request
	Given a initiator peer, a receiver peer and list of known peers
	When are sending the SYN, the target cannot be reached
	Then we get the unreached peer
*/
func TestACKSpreadUnreachables(t *testing.T) {
	initP := uniris.NewLocalPeer("key", net.ParseIP("127.0.0.1"), 3000, "1.0", 30.0, 10.0)

	target := uniris.NewPeerDigest(
		uniris.NewPeerIdentity(net.ParseIP("20.100.4.120"), 3000, "key2"),
		uniris.NewPeerHeartbeatState(time.Now(), 0))

	p1 := uniris.NewDiscoveredPeer(
		uniris.NewPeerIdentity(net.ParseIP("50.20.100.2"), 3000, "key3"),
		uniris.NewPeerHeartbeatState(time.Now(), 0),
		uniris.NewPeerAppState("1.0", uniris.OkPeerStatus, 30.0, 10.0, "", 0, 1, 0),
	)

	p2 := uniris.NewDiscoveredPeer(
		uniris.NewPeerIdentity(net.ParseIP("50.10.30.2"), 3000, "uKey1"),
		uniris.NewPeerHeartbeatState(time.Now(), 0),
		uniris.NewPeerAppState("1.0", uniris.OkPeerStatus, 30.0, 10.0, "", 0, 1, 0),
	)

	r := round{initP, target, mockMessengerWithAckFailure{}}

	kp := []uniris.Peer{p1, p2}

	discoveries := make(chan uniris.Peer)
	reaches := make(chan uniris.Peer)
	unreaches := make(chan uniris.Peer)

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

	unreachP := make([]uniris.Peer, 0)
	for p := range unreaches {
		unreachP = append(unreachP, p)
	}

	discovP := make([]uniris.Peer, 0)
	for p := range discoveries {
		discovP = append(discovP, p)
	}

	assert.NotEmpty(t, unreachP)
	assert.Equal(t, 1, len(unreachP))
	assert.Equal(t, target.Identity().PublicKey(), unreachP[0].Identity().PublicKey())

	assert.Empty(t, discovP)
}

/*
Scenario: Gossip but gets an unexpected error
	Given a initiator peer, a receiver peer and list of known peers
	When are gossiping, we get an unexpected error
	Then we get retrieve the error
*/
func TestSpreadUnexpectedError(t *testing.T) {
	initP := uniris.NewLocalPeer("key", net.ParseIP("127.0.0.1"), 3000, "1.0", 30.0, 10.0)

	target := uniris.NewPeerDigest(
		uniris.NewPeerIdentity(net.ParseIP("20.100.4.120"), 3000, "key2"),
		uniris.NewPeerHeartbeatState(time.Now(), 0))

	p1 := uniris.NewDiscoveredPeer(
		uniris.NewPeerIdentity(net.ParseIP("50.20.100.2"), 3000, "key3"),
		uniris.NewPeerHeartbeatState(time.Now(), 0),
		uniris.NewPeerAppState("1.0", uniris.OkPeerStatus, 30.0, 10.0, "", 0, 1, 0),
	)

	p2 := uniris.NewDiscoveredPeer(
		uniris.NewPeerIdentity(net.ParseIP("50.10.30.2"), 3000, "uKey1"),
		uniris.NewPeerHeartbeatState(time.Now(), 0),
		uniris.NewPeerAppState("1.0", uniris.OkPeerStatus, 30.0, 10.0, "", 0, 1, 0),
	)

	r := round{initP, target, mockMessengerUnexpectedFailure{}}

	kp := []uniris.Peer{p1, p2}

	err := r.run(kp, nil, nil, nil)
	assert.NotNil(t, err)
	assert.Error(t, err, "Unexpected")
}

type mockMessengerWithSynFailure struct {
}

func (m mockMessengerWithSynFailure) SendSyn(source uniris.Peer, target uniris.Peer, known []uniris.Peer) (unknown []uniris.Peer, new []uniris.Peer, err error) {
	return nil, nil, ErrUnreachablePeer
}

func (m mockMessengerWithSynFailure) SendAck(source uniris.Peer, target uniris.Peer, requested []uniris.Peer) error {
	return nil
}

type mockMessengerWithAckFailure struct {
}

func (m mockMessengerWithAckFailure) SendSyn(source uniris.Peer, target uniris.Peer, known []uniris.Peer) (unknown []uniris.Peer, new []uniris.Peer, err error) {
	tar := uniris.NewLocalPeer("uKey1", net.ParseIP("200.18.186.39"), 3000, "1.1", 30.0, 10.0)

	hb := uniris.NewPeerHeartbeatState(time.Now(), 0)
	as := uniris.NewPeerAppState("1.0", uniris.OkPeerStatus, 30.0, 10.0, "", 0, 1, 0)

	np1 := uniris.NewDiscoveredPeer(
		uniris.NewPeerIdentity(net.ParseIP("35.200.100.2"), 3000, "dKey1"),
		hb, as,
	)

	newPeers := []uniris.Peer{np1}

	unknownPeers := []uniris.Peer{tar}

	return unknownPeers, newPeers, nil
}

func (m mockMessengerWithAckFailure) SendAck(source uniris.Peer, target uniris.Peer, requested []uniris.Peer) error {
	return ErrUnreachablePeer
}

type mockMessengerUnexpectedFailure struct {
}

func (m mockMessengerUnexpectedFailure) SendSyn(source uniris.Peer, target uniris.Peer, known []uniris.Peer) (unknown []uniris.Peer, new []uniris.Peer, err error) {
	return nil, nil, errors.New("Unexpected")
}

func (m mockMessengerUnexpectedFailure) SendAck(source uniris.Peer, target uniris.Peer, requested []uniris.Peer) error {
	return nil
}

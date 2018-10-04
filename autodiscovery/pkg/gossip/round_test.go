package gossip

import (
	"errors"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	discovery "github.com/uniris/uniris-core/autodiscovery/pkg"
)

/*
Scenario: Spread a gossip round and discover peers
	Given a initiator peer, a receiver peer and list of known peers
	When we start a gossip round we spread what we know
	Then we get the new peers discovered
*/
func TestSpreadDiscoveries(t *testing.T) {
	initP := discovery.NewStartupPeer([]byte("key"), net.ParseIP("127.0.0.1"), 3000, "1.0", discovery.PeerPosition{})

	target := discovery.NewPeerDigest(
		discovery.NewPeerIdentity(net.ParseIP("20.100.4.120"), 3000, []byte("key2")),
		discovery.NewPeerHeartbeatState(time.Now(), 0))

	p1 := discovery.NewDiscoveredPeer(
		discovery.NewPeerIdentity(net.ParseIP("50.20.100.2"), 3000, []byte("key3")),
		discovery.NewPeerHeartbeatState(time.Now(), 0),
		discovery.NewPeerAppState("1.0", discovery.OkStatus, discovery.PeerPosition{}, "", 0, 1, 0),
	)

	p2 := discovery.NewDiscoveredPeer(
		discovery.NewPeerIdentity(net.ParseIP("50.10.30.2"), 3000, []byte("uKey1")),
		discovery.NewPeerHeartbeatState(time.Now(), 0),
		discovery.NewPeerAppState("1.0", discovery.OkStatus, discovery.PeerPosition{}, "", 0, 1, 0),
	)

	g := NewGossipRound(initP, target, mockMessenger{})

	kp := []discovery.Peer{p1, p2}

	discoveries := make(chan discovery.Peer)
	reaches := make(chan discovery.Peer)

	var wg sync.WaitGroup
	wg.Add(2)

	go g.Spread(kp, discoveries, reaches, nil)

	pp := make([]discovery.Peer, 0)
	go func() {
		for p := range discoveries {
			pp = append(pp, p)
			wg.Done()
			close(discoveries)
		}
	}()

	reachP := make([]discovery.Peer, 0)
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
	assert.Equal(t, "dKey1", pp[0].Identity().PublicKey().String())

	assert.NotEmpty(t, reachP)
}

/*
Scenario: Spread gossip but unreach the target peer during the SYN request
	Given a initiator peer, a receiver peer and list of known peers
	When are sending the SYN, the target cannot be reached
	Then we get the unreached peer
*/
func TestSYNSpreadUnreachables(t *testing.T) {
	initP := discovery.NewStartupPeer([]byte("key"), net.ParseIP("127.0.0.1"), 3000, "1.0", discovery.PeerPosition{})

	target := discovery.NewPeerDigest(
		discovery.NewPeerIdentity(net.ParseIP("20.100.4.120"), 3000, []byte("key2")),
		discovery.NewPeerHeartbeatState(time.Now(), 0))

	p1 := discovery.NewDiscoveredPeer(
		discovery.NewPeerIdentity(net.ParseIP("50.20.100.2"), 3000, []byte("key3")),
		discovery.NewPeerHeartbeatState(time.Now(), 0),
		discovery.NewPeerAppState("1.0", discovery.OkStatus, discovery.PeerPosition{}, "", 0, 1, 0),
	)

	p2 := discovery.NewDiscoveredPeer(
		discovery.NewPeerIdentity(net.ParseIP("50.10.30.2"), 3000, []byte("uKey1")),
		discovery.NewPeerHeartbeatState(time.Now(), 0),
		discovery.NewPeerAppState("1.0", discovery.OkStatus, discovery.PeerPosition{}, "", 0, 1, 0),
	)

	g := NewGossipRound(initP, target, mockMessengerWithSynFailure{})

	kp := []discovery.Peer{p1, p2}

	unreaches := make(chan discovery.Peer)

	go func() {
		err := g.Spread(kp, nil, nil, unreaches)
		assert.Nil(t, err)
		close(unreaches)
	}()

	pp := make([]discovery.Peer, 0)
	for p := range unreaches {
		pp = append(pp, p)
	}

	assert.NotEmpty(t, pp)
	assert.Equal(t, 1, len(pp))
	assert.Equal(t, target.Identity().PublicKey().String(), pp[0].Identity().PublicKey().String())
}

/*
Scenario: Spread gossip but unreach the target peer during the SYN request
	Given a initiator peer, a receiver peer and list of known peers
	When are sending the SYN, the target cannot be reached
	Then we get the unreached peer
*/
func TestACKSpreadUnreachables(t *testing.T) {
	initP := discovery.NewStartupPeer([]byte("key"), net.ParseIP("127.0.0.1"), 3000, "1.0", discovery.PeerPosition{})

	target := discovery.NewPeerDigest(
		discovery.NewPeerIdentity(net.ParseIP("20.100.4.120"), 3000, []byte("key2")),
		discovery.NewPeerHeartbeatState(time.Now(), 0))

	p1 := discovery.NewDiscoveredPeer(
		discovery.NewPeerIdentity(net.ParseIP("50.20.100.2"), 3000, []byte("key3")),
		discovery.NewPeerHeartbeatState(time.Now(), 0),
		discovery.NewPeerAppState("1.0", discovery.OkStatus, discovery.PeerPosition{}, "", 0, 1, 0),
	)

	p2 := discovery.NewDiscoveredPeer(
		discovery.NewPeerIdentity(net.ParseIP("50.10.30.2"), 3000, []byte("uKey1")),
		discovery.NewPeerHeartbeatState(time.Now(), 0),
		discovery.NewPeerAppState("1.0", discovery.OkStatus, discovery.PeerPosition{}, "", 0, 1, 0),
	)

	g := NewGossipRound(initP, target, mockMessengerWithAckFailure{})

	kp := []discovery.Peer{p1, p2}

	discoveries := make(chan discovery.Peer)
	reaches := make(chan discovery.Peer)
	unreaches := make(chan discovery.Peer)

	go func() {
		err := g.Spread(kp, discoveries, reaches, unreaches)
		assert.Nil(t, err)
		close(unreaches)
		close(discoveries)
		close(reaches)
	}()

	go func() {
		for range reaches {
		}
	}()

	unreachP := make([]discovery.Peer, 0)
	for p := range unreaches {
		unreachP = append(unreachP, p)
	}

	discovP := make([]discovery.Peer, 0)
	for p := range discoveries {
		discovP = append(discovP, p)
	}

	assert.NotEmpty(t, unreachP)
	assert.Equal(t, 1, len(unreachP))
	assert.Equal(t, target.Identity().PublicKey().String(), unreachP[0].Identity().PublicKey().String())

	assert.Empty(t, discovP)
}

/*
Scenario: Gossip but gets an unexpected error
	Given a initiator peer, a receiver peer and list of known peers
	When are gossiping, we get an unexpected error
	Then we get retrieve the error
*/
func TestSpreadUnexpectedError(t *testing.T) {
	initP := discovery.NewStartupPeer([]byte("key"), net.ParseIP("127.0.0.1"), 3000, "1.0", discovery.PeerPosition{})

	target := discovery.NewPeerDigest(
		discovery.NewPeerIdentity(net.ParseIP("20.100.4.120"), 3000, []byte("key2")),
		discovery.NewPeerHeartbeatState(time.Now(), 0))

	p1 := discovery.NewDiscoveredPeer(
		discovery.NewPeerIdentity(net.ParseIP("50.20.100.2"), 3000, []byte("key3")),
		discovery.NewPeerHeartbeatState(time.Now(), 0),
		discovery.NewPeerAppState("1.0", discovery.OkStatus, discovery.PeerPosition{}, "", 0, 1, 0),
	)

	p2 := discovery.NewDiscoveredPeer(
		discovery.NewPeerIdentity(net.ParseIP("50.10.30.2"), 3000, []byte("uKey1")),
		discovery.NewPeerHeartbeatState(time.Now(), 0),
		discovery.NewPeerAppState("1.0", discovery.OkStatus, discovery.PeerPosition{}, "", 0, 1, 0),
	)

	g := NewGossipRound(initP, target, mockMessengerUnexpectedFailure{})

	kp := []discovery.Peer{p1, p2}

	err := g.Spread(kp, nil, nil, nil)
	assert.NotNil(t, err)
	assert.Error(t, err, "Unexpected")
}

type mockMessengerWithSynFailure struct {
}

func (m mockMessengerWithSynFailure) SendSyn(req SynRequest) (*SynAck, error) {
	return nil, ErrUnreachablePeer
}

func (m mockMessengerWithSynFailure) SendAck(req AckRequest) error {
	return nil
}

type mockMessengerWithAckFailure struct {
}

func (m mockMessengerWithAckFailure) SendSyn(req SynRequest) (*SynAck, error) {
	init := discovery.NewStartupPeer([]byte("key"), net.ParseIP("127.0.0.1"), 3000, "1.0", discovery.PeerPosition{})
	tar := discovery.NewStartupPeer([]byte("uKey1"), net.ParseIP("200.18.186.39"), 3000, "1.1", discovery.PeerPosition{})

	hb := discovery.NewPeerHeartbeatState(time.Now(), 0)
	as := discovery.NewPeerAppState("1.0", discovery.OkStatus, discovery.PeerPosition{}, "", 0, 1, 0)

	np1 := discovery.NewDiscoveredPeer(
		discovery.NewPeerIdentity(net.ParseIP("35.200.100.2"), 3000, []byte("dKey1")),
		hb, as,
	)

	newPeers := []discovery.Peer{np1}

	unknownPeers := []discovery.Peer{tar}

	return &SynAck{
		Initiator:    init,
		Target:       tar,
		NewPeers:     newPeers,
		UnknownPeers: unknownPeers,
	}, nil
}

func (m mockMessengerWithAckFailure) SendAck(req AckRequest) error {
	return ErrUnreachablePeer
}

type mockMessengerUnexpectedFailure struct {
}

func (m mockMessengerUnexpectedFailure) SendSyn(req SynRequest) (*SynAck, error) {
	return nil, errors.New("Unexpected")
}

func (m mockMessengerUnexpectedFailure) SendAck(req AckRequest) error {
	return nil
}

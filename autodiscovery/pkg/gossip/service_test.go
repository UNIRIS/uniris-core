package gossip

import (
	"errors"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	discovery "github.com/uniris/uniris-core/autodiscovery/pkg"
	"github.com/uniris/uniris-core/autodiscovery/pkg/mock"
)

/*
Scenario: Run cycle
	Given a initiator peer, a receiver peer and list of known peers
	When we start a gossip round, we run a gossip cycle to spread
	Then we get the new peers discovered
*/
func TestRunCycle(t *testing.T) {
	initP := discovery.NewStartupPeer([]byte("key"), net.ParseIP("127.0.0.1"), 3000, "1.0", discovery.PeerPosition{})

	repo := new(mock.Repository)
	repo.AddPeer(initP)

	id1 := discovery.NewPeerIdentity(net.ParseIP("20.100.4.120"), 3000, []byte("key2"))
	hb := discovery.NewPeerHeartbeatState(time.Now(), 0)
	as := discovery.NewPeerAppState("1.0", discovery.OkStatus, discovery.PeerPosition{}, "", 0, 1, 0)
	recP := discovery.NewDiscoveredPeer(id1, hb, as)

	p1 := discovery.NewDiscoveredPeer(
		discovery.NewPeerIdentity(net.ParseIP("50.20.100.2"), 3000, []byte("key3")),
		hb, as,
	)

	p2 := discovery.NewDiscoveredPeer(
		discovery.NewPeerIdentity(net.ParseIP("50.10.30.2"), 3000, []byte("uKey1")),
		hb, as,
	)

	g := NewService(repo, mockMessenger{}, new(mock.Notifier), mockMonitor{})

	newPeers, err := g.RunCycle(initP, recP, []discovery.Peer{p1, p2})
	assert.Nil(t, err)
	assert.NotEmpty(t, newPeers)

	assert.Equal(t, 1, len(newPeers))
	assert.Equal(t, "dKey1", newPeers[0].Identity().PublicKey().String())
}

/*
Scenario: Gossip across a selection of peers
	Given a initiator peer, seeds and known peers stored locally
	When we gossip
	Then the new peers are stored and notified
*/
func TestGossip(t *testing.T) {
	init := discovery.NewStartupPeer([]byte("key"), net.ParseIP("127.0.0.1"), 3000, "1.0", discovery.PeerPosition{})

	repo := new(mock.Repository)
	notif := new(mock.Notifier)

	repo.AddPeer(init)

	s := discovery.Seed{IP: net.ParseIP("10.0.0.1"), Port: 3000, PublicKey: "keyss"}
	repo.AddSeed(s)

	srv := NewService(repo, mockMessenger{}, notif, new(mockMonitor))
	err := srv.Spread(init)
	assert.Nil(t, err)

	peers, _ := repo.ListKnownPeers()
	assert.Equal(t, 2, len(peers))
	assert.Equal(t, "key", string(peers[0].Identity().PublicKey()))
	assert.Equal(t, "dKey1", string(peers[1].Identity().PublicKey()))

	assert.NotEmpty(t, notif.NotifiedPeers())
	assert.Equal(t, 1, len(notif.NotifiedPeers()))
}

/*
Scenario: Gossip across a selection of peers
	Given a initiator peer, seeds and known peers stored locally
	When we gossip spread get unreacheable error
	Then unreacheable peer is stored on the repo
*/
func TestAddUnreachable(t *testing.T) {
	init := discovery.NewStartupPeer([]byte("key"), net.ParseIP("127.0.0.1"), 3000, "1.0", discovery.PeerPosition{})
	repo := new(mock.Repository)
	notif := new(mock.Notifier)
	repo.AddPeer(init)

	s := discovery.Seed{IP: net.ParseIP("127.0.0.1"), Port: 3000, PublicKey: "key"}
	repo.AddSeed(s)

	hb := discovery.NewPeerHeartbeatState(time.Now(), 0)
	as := discovery.NewPeerAppState("1.0", discovery.OkStatus, discovery.PeerPosition{}, "", 0, 1, 0)
	p1 := discovery.NewDiscoveredPeer(
		discovery.NewPeerIdentity(net.ParseIP("50.10.30.2"), 3000, []byte("ukey")),
		hb, as,
	)

	repo.AddPeer(p1)
	assert.Equal(t, 2, len(repo.Peers))
	assert.Equal(t, 0, len(repo.UnreachablePeers))
	srv := NewService(repo, mockMessengerunreacheable{}, notif, new(mockMonitor))
	err := srv.Spread(init)
	assert.NotNil(t, err)
	fmt.Printf(err.Error())
	fmt.Print(repo.UnreachablePeers)
	assert.Equal(t, 1, len(repo.UnreachablePeers))
	assert.Equal(t, "ukey", repo.UnreachablePeers[0].String())
}

/*
Scenario: Gossip across a selection of peers
	Given a initiator peer, seeds and known peers stored locally and one unreacheable peer
	When we gossip
	Then unreacheable peer is removed from the repo
*/
func TestDelUnreachable(t *testing.T) {
	init := discovery.NewStartupPeer([]byte("key"), net.ParseIP("127.0.0.1"), 3000, "1.0", discovery.PeerPosition{})
	repo := new(mock.Repository)
	notif := new(mock.Notifier)
	repo.AddPeer(init)

	s := discovery.Seed{IP: net.ParseIP("127.0.0.1"), Port: 3000, PublicKey: "key"}
	repo.AddSeed(s)

	hb := discovery.NewPeerHeartbeatState(time.Now(), 0)
	as := discovery.NewPeerAppState("1.0", discovery.OkStatus, discovery.PeerPosition{}, "", 0, 1, 0)
	p1 := discovery.NewDiscoveredPeer(
		discovery.NewPeerIdentity(net.ParseIP("50.10.30.2"), 3000, []byte("ukey")),
		hb, as,
	)
	repo.AddPeer(p1)
	repo.AddUnreachablePeer(p1.Identity().PublicKey())
	rp, _ := repo.ListReachablePeers()
	assert.Equal(t, 1, len(rp))
	assert.Equal(t, 2, len(repo.Peers))
	assert.Equal(t, 1, len(repo.UnreachablePeers))
	srv := NewService(repo, mockMessenger{}, notif, new(mockMonitor))
	err := srv.Spread(init)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(repo.UnreachablePeers))
	assert.Equal(t, 3, len(repo.Peers))
}

type mockMessenger struct {
}

func (m mockMessenger) SendSyn(req SynRequest) (*SynAck, error) {
	init := discovery.NewStartupPeer([]byte("key"), net.ParseIP("127.0.0.1"), 3000, "1.0", discovery.PeerPosition{})
	rec := discovery.NewStartupPeer([]byte("uKey1"), net.ParseIP("200.18.186.39"), 3000, "1.1", discovery.PeerPosition{})

	hb := discovery.NewPeerHeartbeatState(time.Now(), 0)
	as := discovery.NewPeerAppState("1.0", discovery.OkStatus, discovery.PeerPosition{}, "", 0, 1, 0)

	np1 := discovery.NewDiscoveredPeer(
		discovery.NewPeerIdentity(net.ParseIP("35.200.100.2"), 3000, []byte("dKey1")),
		hb, as,
	)

	newPeers := []discovery.Peer{np1}

	unknownPeers := []discovery.Peer{rec}

	return &SynAck{
		Initiator:    init,
		Receiver:     rec,
		NewPeers:     newPeers,
		UnknownPeers: unknownPeers,
	}, nil
}

func (m mockMessenger) SendAck(req AckRequest) error {
	return nil
}

type mockMonitor struct{}

//RefreshPeer updates the peer's metrics retrieved from the peer monitor
func (s mockMonitor) RefreshPeer(p discovery.Peer) error {
	return nil
}

func (s mockMonitor) PeerStatus(p discovery.Peer) (discovery.PeerStatus, error) {
	return discovery.OkStatus, nil
}

type mockMessengerunreacheable struct {
}

func (m mockMessengerunreacheable) SendSyn(req SynRequest) (*SynAck, error) {
	return nil, errors.New("Unreachable Peer")
}

func (m mockMessengerunreacheable) SendAck(req AckRequest) error {
	return nil
}

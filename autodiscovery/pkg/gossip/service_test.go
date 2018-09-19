package gossip

import (
	"encoding/hex"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	discovery "github.com/uniris/uniris-core/autodiscovery/pkg"
)

/*
Scenario: Converts a list of peer into a map of peers
	Given a list of peers
	When we want to create a map of it
	Then we get a map of peer identified by their public key
*/
func TestMapPeers(t *testing.T) {
	p1 := discovery.NewStartupPeer([]byte("key"), net.ParseIP("127.0.0.1"), 3000, "1.0", discovery.PeerPosition{}, 1)
	p2 := discovery.NewPeerDigest([]byte("key2"), net.ParseIP("10.0.0.1"), 3000)

	g := service{}
	mPeers := g.mapPeers([]discovery.Peer{p1, p2})
	assert.NotNil(t, mPeers)
	assert.NotEmpty(t, mPeers)
	assert.Equal(t, 2, len(mPeers))

	assert.NotNil(t, mPeers[hex.EncodeToString([]byte("key"))])
	assert.NotNil(t, mPeers[hex.EncodeToString([]byte("key2"))])
	assert.Equal(t, "127.0.0.1", mPeers[hex.EncodeToString([]byte("key"))].IP().String())
}

/*
Scenario: Run cycle
	Given a initiator peer, a receiver peer and list of known peers
	When we start a gossip round, we run a gossip cycle to spread
	Then we get the new peers discovered
*/
func TestRunCycle(t *testing.T) {
	initP := discovery.NewStartupPeer([]byte("key"), net.ParseIP("127.0.0.1"), 3000, "1.0", discovery.PeerPosition{}, 1)
	recP := discovery.NewPeerDigest([]byte("key2"), net.ParseIP("20.100.4.120"), 3000)

	p1 := discovery.NewPeerDigest([]byte("key3"), net.ParseIP("50.20.100.2"), 3000)
	p2 := discovery.NewPeerDigest([]byte("uKey1"), net.ParseIP("50.10.30.2"), 3000)

	g := service{
		msg: mockMessenger{},
	}
	newPeers, err := g.RunCycle(initP, recP, []discovery.Peer{p1, p2})
	assert.Nil(t, err)
	assert.NotEmpty(t, newPeers)

	assert.Equal(t, 1, len(newPeers))
	assert.Equal(t, "dKey1", string(newPeers[0].PublicKey()))
}

/*
Scenario: Gossip across a selection of peers
	Given a initiator peer, seeds and known peers stored locally
	When we gossip
	Then the new peers are stored and notified
*/
func TestGossip(t *testing.T) {
	init := discovery.NewStartupPeer([]byte("key"), net.ParseIP("127.0.0.1"), 3000, "1.0", discovery.PeerPosition{}, 1)

	repo := new(mockPeerRepository)
	notif := new(mockNotifier)

	repo.AddPeer(init)

	s := discovery.Seed{IP: net.ParseIP("10.0.0.1"), Port: 3000}
	repo.AddSeed(s)

	srv := NewService(repo, mockMessenger{}, notif, new(mockInspect))
	err := srv.Spread(init)
	assert.Nil(t, err)

	peers, _ := repo.ListKnownPeers()
	assert.Equal(t, 2, len(peers))
	assert.Equal(t, "key", string(peers[0].PublicKey()))
	assert.Equal(t, "dKey1", string(peers[1].PublicKey()))

	assert.NotEmpty(t, notif.NotifiedPeers())
	assert.Equal(t, 1, len(notif.NotifiedPeers()))
}

/*
Scenario: Gets diff between our known peers and a list of peers with peers unknown from the both sides
	Given a unknown list of peers and a list known peers unknown from the sender
	When we want get to the diff between
	Then we retrieve a list of peers not include inside the list, and a peer unknows from us
*/
func TestDiffPeersWithDifferentPeers(t *testing.T) {
	repo := new(mockPeerRepository)

	kp := discovery.NewStartupPeer([]byte("key"), net.ParseIP("127.0.0.1"), 3000, "1.0", discovery.PeerPosition{}, 1)
	kp2 := discovery.NewStartupPeer([]byte("key2"), net.ParseIP("80.200.100.2"), 3000, "1.0", discovery.PeerPosition{}, 1)

	repo.AddPeer(kp)
	repo.AddPeer(kp2)

	srv := NewService(repo, new(mockMessenger), new(mockNotifier), new(mockInspect))

	np1 := discovery.NewPeerDigest([]byte("key3"), net.ParseIP("10.0.0.1"), 3000)
	np2 := discovery.NewPeerDigest([]byte("key4"), net.ParseIP("50.0.0.1"), 3000)

	diff, err := srv.DiffPeers([]discovery.Peer{np1, np2})
	assert.Nil(t, err)
	assert.NotEmpty(t, diff.UnknownLocally)
	assert.Equal(t, 2, len(diff.UnknownLocally))
	assert.Equal(t, "key3", string(diff.UnknownLocally[0].PublicKey()))
	assert.Equal(t, "key4", string(diff.UnknownLocally[1].PublicKey()))

	assert.NotEmpty(t, diff.UnknownRemotly)
	assert.Equal(t, 2, len(diff.UnknownRemotly))
	assert.Equal(t, "key", string(diff.UnknownRemotly[0].PublicKey()))
	assert.Equal(t, "key2", string(diff.UnknownRemotly[1].PublicKey()))
}

/*
Scenario: Gets diff between our known peers and a list of peers which include one of our peer
	Given a list of peers including one of our peer and a list of known peer
	When we want to get the diff
	Then we get the only the peer that the list don't know and we don' know
*/
func TestDiffPeerWithSomeKnownPeers(t *testing.T) {
	repo := new(mockPeerRepository)

	kp := discovery.NewStartupPeer([]byte("key"), net.ParseIP("127.0.0.1"), 3000, "1.0", discovery.PeerPosition{}, 1)
	kp2 := discovery.NewStartupPeer([]byte("key2"), net.ParseIP("80.200.100.2"), 3000, "1.0", discovery.PeerPosition{}, 1)

	repo.AddPeer(kp)
	repo.AddPeer(kp2)

	srv := NewService(repo, new(mockMessenger), new(mockNotifier), new(mockInspect))

	np1 := discovery.NewPeerDigest([]byte("key"), net.ParseIP("127.0.0.1"), 3000)
	np2 := discovery.NewPeerDigest([]byte("key4"), net.ParseIP("50.0.0.1"), 3000)

	diff, err := srv.DiffPeers([]discovery.Peer{np1, np2})
	assert.Nil(t, err)
	assert.NotEmpty(t, diff.UnknownLocally)
	assert.Equal(t, 1, len(diff.UnknownLocally))
	assert.Equal(t, "key4", string(diff.UnknownLocally[0].PublicKey()))
	assert.NotEmpty(t, diff.UnknownRemotly)
	assert.Equal(t, 1, len(diff.UnknownRemotly))
	assert.Equal(t, "key2", string(diff.UnknownRemotly[0].PublicKey()))
}

/*
Scenario: Gets diff between an empty list of peers and known peers
	Given a empty list of peers and a known list of peers
	When we want to get the diff
	Then we provide only our known peers
*/
func TestDiffWithEmptyPeers(t *testing.T) {
	repo := new(mockPeerRepository)

	kp := discovery.NewStartupPeer([]byte("key"), net.ParseIP("127.0.0.1"), 3000, "1.0", discovery.PeerPosition{}, 1)
	kp2 := discovery.NewPeerDigest([]byte("key2"), net.ParseIP("80.200.100.2"), 3000)

	repo.AddPeer(kp)
	repo.AddPeer(kp2)

	srv := NewService(repo, new(mockMessenger), new(mockNotifier), new(mockInspect))

	diff, err := srv.DiffPeers([]discovery.Peer{})
	assert.Nil(t, err)
	assert.Empty(t, diff.UnknownLocally)
	assert.NotEmpty(t, diff.UnknownRemotly)
	assert.Equal(t, 2, len(diff.UnknownRemotly))
	assert.Equal(t, "key", string(diff.UnknownRemotly[0].PublicKey()))
	assert.Equal(t, "key2", string(diff.UnknownRemotly[1].PublicKey()))
}

/*
Scenario: Gets diff between identically list of peers
	Given a list of peers identical to our known list of peers
	When we want to get the diff
	Then we provide empty lists
*/
func TestDiffPeerWithSamePeers(t *testing.T) {
	repo := new(mockPeerRepository)

	kp := discovery.NewStartupPeer([]byte("key"), net.ParseIP("127.0.0.1"), 3000, "1.0", discovery.PeerPosition{}, 1)
	kp2 := discovery.NewStartupPeer([]byte("key2"), net.ParseIP("80.200.100.2"), 3000, "1.0", discovery.PeerPosition{}, 1)

	repo.AddPeer(kp)
	repo.AddPeer(kp2)

	srv := NewService(repo, new(mockMessenger), new(mockNotifier), new(mockInspect))

	np1 := discovery.NewPeerDetailed([]byte("key"), net.ParseIP("127.0.0.1"), 3000, time.Now(), nil)
	np2 := discovery.NewPeerDetailed([]byte("key2"), net.ParseIP("80.200.100.2"), 3000, time.Now(), nil)

	diff, err := srv.DiffPeers([]discovery.Peer{np1, np2})
	assert.Nil(t, err)
	assert.Empty(t, diff.UnknownLocally)
	assert.Empty(t, diff.UnknownRemotly)
}

type mockMessenger struct {
}

func (m mockMessenger) SendSyn(req SynRequest) (*SynAck, error) {
	init := discovery.NewStartupPeer([]byte("key"), net.ParseIP("127.0.0.1"), 3000, "1.0", discovery.PeerPosition{}, 1)
	rec := discovery.NewStartupPeer([]byte("uKey1"), net.ParseIP("200.18.186.39"), 3000, "1.1", discovery.PeerPosition{}, 1)

	t, _ := time.Parse(
		time.RFC3339,
		"2012-11-01T22:08:41+00:00")
	np1 := discovery.NewPeerDetailed([]byte("dKey1"), net.ParseIP("35.200.100.2"), 3000, t, nil)

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

type mockNotifier struct {
	notifiedPeers []discovery.Peer
}

func (n mockNotifier) NotifiedPeers() []discovery.Peer {
	return n.notifiedPeers
}

func (n *mockNotifier) Notify(p discovery.Peer) {
	n.notifiedPeers = append(n.notifiedPeers, p)
}

type mockPeerRepository struct {
	peers []discovery.Peer
	seeds []discovery.Seed
}

func (r *mockPeerRepository) AddPeer(p discovery.Peer) error {
	if r.containsPeer(p) {
		return r.UpdatePeer(p)
	}
	r.peers = append(r.peers, p)
	return nil
}

func (r *mockPeerRepository) AddSeed(s discovery.Seed) error {
	r.seeds = append(r.seeds, s)
	return nil
}

func (r *mockPeerRepository) ListKnownPeers() ([]discovery.Peer, error) {
	return r.peers, nil
}

func (r *mockPeerRepository) ListSeedPeers() ([]discovery.Seed, error) {
	return r.seeds, nil
}

func (r *mockPeerRepository) UpdatePeer(peer discovery.Peer) error {
	for _, p := range r.peers {
		if string(p.PublicKey()) == string(peer.PublicKey()) {
			p = peer
			break
		}
	}
	return nil
}

func (r *mockPeerRepository) containsPeer(p discovery.Peer) bool {
	mPeers := make(map[string]discovery.Peer, 0)
	for _, p := range r.peers {
		mPeers[hex.EncodeToString(p.PublicKey())] = p
	}

	_, exist := mPeers[hex.EncodeToString(p.PublicKey())]
	return exist
}

type mockInspect struct{}

//RefreshPeer updates the peer's metrics retrieved from the peer monitor
func (s mockInspect) RefreshPeer(p *discovery.Peer) error {
	return nil
}

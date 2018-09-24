package gossip

import (
	"encoding/hex"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	discovery "github.com/uniris/uniris-core/autodiscovery/pkg"
	"github.com/uniris/uniris-core/autodiscovery/pkg/monitoring"
)

/*
Scenario: Run a gossip cycle
	Given a initator peer, a list of seeds and a list known peer
	When we gossip
	Then we get new peers are stored and notified
*/
func TestRunGossip(t *testing.T) {
	repo := new(mockPeerRepository)

	init := discovery.NewStartupPeer([]byte("key"), net.ParseIP("127.0.0.1"), 3000, "1.0", discovery.PeerPosition{}, 1)
	repo.SetPeer(init)

	kp := discovery.NewPeerDigest([]byte("key2"), net.ParseIP("10.0.0.2"), 3000)
	kp2 := discovery.NewPeerDigest([]byte("key3"), net.ParseIP("80.200.100.2"), 3000)
	repo.SetPeer(kp)
	repo.SetPeer(kp2)

	notif := new(mockNotifier)
	spr := new(mockSpreader)

	monSrv := monitoring.NewService(repo, new(mockMonitor))

	s := service{
		repo:  repo,
		mon:   monSrv,
		notif: notif,
		spr:   spr,
	}
	err := s.runGossip(init, []discovery.Seed{discovery.Seed{IP: net.ParseIP("30.0.0.1"), Port: 3000}})
	assert.Nil(t, err)

	pp, _ := repo.ListKnownPeers()
	assert.Equal(t, 4, len(pp))
	assert.Equal(t, "dKey1", string(pp[3].PublicKey()))

	npp := notif.NotifiedPeers()
	assert.NotEmpty(t, npp)
	assert.Equal(t, "dKey1", string(npp[0].PublicKey()))
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

	repo.SetPeer(kp)
	repo.SetPeer(kp2)

	srv := service{repo: repo}

	np1 := discovery.NewPeerDigest([]byte("key3"), net.ParseIP("10.0.0.1"), 3000)
	np2 := discovery.NewPeerDigest([]byte("key4"), net.ParseIP("50.0.0.1"), 3000)

	diff, err := srv.ComparePeers([]discovery.Peer{np1, np2})
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

	repo.SetPeer(kp)
	repo.SetPeer(kp2)

	srv := service{repo: repo}
	np1 := discovery.NewPeerDigest([]byte("key"), net.ParseIP("127.0.0.1"), 3000)
	np2 := discovery.NewPeerDigest([]byte("key4"), net.ParseIP("50.0.0.1"), 3000)

	diff, err := srv.ComparePeers([]discovery.Peer{np1, np2})
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

	repo.SetPeer(kp)
	repo.SetPeer(kp2)

	srv := service{repo: repo}
	diff, err := srv.ComparePeers([]discovery.Peer{})
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

	repo.SetPeer(kp)
	repo.SetPeer(kp2)

	srv := service{repo: repo}
	np1 := discovery.NewPeerDetailed([]byte("key"), net.ParseIP("127.0.0.1"), 3000, time.Now(), false, nil)
	np2 := discovery.NewPeerDetailed([]byte("key2"), net.ParseIP("80.200.100.2"), 3000, time.Now(), false, nil)

	diff, err := srv.ComparePeers([]discovery.Peer{np1, np2})
	assert.Nil(t, err)
	assert.Empty(t, diff.UnknownLocally)
	assert.Empty(t, diff.UnknownRemotly)
}

type mockPeerRepository struct {
	peers []discovery.Peer
	seeds []discovery.Seed
}

func (r *mockPeerRepository) GetOwnedPeer() (p discovery.Peer, err error) {
	for _, p := range r.peers {
		if p.IsOwned() {
			return p, nil
		}
	}
	return
}

func (r *mockPeerRepository) ListSeedPeers() ([]discovery.Seed, error) {
	return r.seeds, nil
}

func (r *mockPeerRepository) ListKnownPeers() ([]discovery.Peer, error) {
	return r.peers, nil
}

func (r *mockPeerRepository) SetPeer(peer discovery.Peer) error {
	if r.containsPeer(peer) {
		for _, p := range r.peers {
			if string(p.PublicKey()) == string(peer.PublicKey()) {
				p = peer
				break
			}
		}
	} else {
		r.peers = append(r.peers, peer)
	}
	return nil
}

func (r *mockPeerRepository) SetSeed(s discovery.Seed) error {
	r.seeds = append(r.seeds, s)
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

type mockSpreader struct {
}

func (m mockSpreader) SendSyn(req discovery.SynRequest) (*discovery.SynAck, error) {
	init := discovery.NewStartupPeer([]byte("key"), net.ParseIP("127.0.0.1"), 3000, "1.0", discovery.PeerPosition{}, 1)
	tar := discovery.NewStartupPeer([]byte("uKey1"), net.ParseIP("200.18.186.39"), 3000, "1.1", discovery.PeerPosition{}, 1)

	t, _ := time.Parse(
		time.RFC3339,
		"2012-11-01T22:08:41+00:00")
	np1 := discovery.NewPeerDetailed([]byte("dKey1"), net.ParseIP("35.200.100.2"), 3000, t, false, nil)

	newPeers := []discovery.Peer{np1}

	unknownPeers := []discovery.Peer{tar}

	return &discovery.SynAck{
		Initiator:    init,
		Target:       tar,
		NewPeers:     newPeers,
		UnknownPeers: unknownPeers,
	}, nil
}

func (m mockSpreader) SendAck(req discovery.AckRequest) error {
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

type mockMonitor struct{}

func (m mockMonitor) Status() (discovery.PeerStatus, error) {
	return discovery.OkStatus, nil
}

func (m mockMonitor) CPULoad() (string, error) {
	return "100.0.0", nil
}

func (m mockMonitor) FreeDiskSpace() (float64, error) {
	return 300.50, nil
}

func (m mockMonitor) IOWaitRate() (float64, error) {
	return 500, nil
}

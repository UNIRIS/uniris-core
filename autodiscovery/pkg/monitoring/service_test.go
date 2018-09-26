package monitoring

import (
	"encoding/hex"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"

	discovery "github.com/uniris/uniris-core/autodiscovery/pkg"
)

/*
Scenario: check refresh
	Given an initial seed
	When refresh
	Then status, CPUload, FreeDiskSpace and IOWaitRate are updated
*/

func TestRefresh(t *testing.T) {
	repo := new(mockPeerRepository)
	watch := new(mockWatcher)
	p1 := discovery.NewStartupPeer([]byte("key"), net.ParseIP("127.0.0.1"), 3000, "1.0", discovery.PeerPosition{})
	srv := NewService(repo, watch)
	err := srv.RefreshPeer(p1)
	assert.Nil(t, err)
	assert.Equal(t, "0.62 0.77 0.71 4/972 26361", p1.CPULoad())
	assert.Equal(t, discovery.OkStatus, p1.Status())
	assert.Equal(t, float64(212383852), p1.FreeDiskSpace())
	assert.Equal(t, 5, p1.DiscoveredPeers())
	assert.Equal(t, 1, p1.P2PFactor())

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

func (r *mockPeerRepository) GetPeerByIP(ip net.IP) (p discovery.Peer, err error) {
	for i := 0; i < len(r.peers); i++ {
		if string(ip) == string(r.peers[i].IP()) {
			return r.peers[i], nil
		}
	}
	return
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

func (r *mockPeerRepository) containsPeer(peer discovery.Peer) bool {
	mPeers := make(map[string]discovery.Peer, 0)
	for _, p := range r.peers {
		mPeers[hex.EncodeToString(p.PublicKey())] = peer
	}

	_, exist := mPeers[hex.EncodeToString(peer.PublicKey())]
	return exist
}

type mockWatcher struct{}

func (w mockWatcher) Status(p discovery.Peer, repo discovery.Repository) (discovery.PeerStatus, error) {
	return discovery.OkStatus, nil
}

func (w mockWatcher) CPULoad() (string, error) {
	return "0.62 0.77 0.71 4/972 26361", nil
}

func (w mockWatcher) FreeDiskSpace() (float64, error) {
	return 212383852, nil
}

func (w mockWatcher) DiscoveredPeer(repo discovery.Repository) (int, error) {
	return 5, nil
}

func (w mockWatcher) P2PFactor() (int, error) {
	return 1, nil
}

package monitoring

import (
	"encoding/hex"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	discovery "github.com/uniris/uniris-core/autodiscovery/pkg"
)

/*
Scenario: Refreshes own peer
	Given an owned peer
	When we refresh it
	Then we get monitoring updates and update the peer status
*/
func TestRefreshOwnPeer(t *testing.T) {
	repo := new(mockPeerRepository)
	p := discovery.NewStartupPeer([]byte("key"), net.ParseIP("127.0.0.1"), 3000, "1.0", discovery.PeerPosition{}, 1)
	repo.SetPeer(p)

	srv := NewService(repo, mockMonitor{})
	err := srv.RefreshOwnedPeer()
	assert.Nil(t, err)

	op, _ := repo.GetOwnedPeer()
	assert.Equal(t, discovery.OkStatus, op.Status())
	assert.Equal(t, "100.0.0", op.CPULoad())
	assert.Equal(t, float64(300.50), op.FreeDiskSpace())
	assert.Equal(t, float64(500), p.IOWaitRate())
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

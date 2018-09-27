package discovery

import (
	"encoding/hex"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

/*
Scenario: check GetSeedDiscoveredPeer
	Given a repo with 3 seed, seed1 discoveredPeersNumber = 5,seed2 discoveredPeersNumber = 6, seed3 discoveredPeersNumber = 7
	When GetSeedDiscoveredPeer call
	Then SeedDiscoveredPeer value is 6
*/

func TestGetSeedDiscoveredPeer(t *testing.T) {
	repo := new(mockPeerRepository)
	sdc := NewSeedDiscoveryCounter(repo)
	seed1 := Seed{IP: net.ParseIP("10.1.1.1"), Port: 3000}
	seed2 := Seed{IP: net.ParseIP("10.1.1.2"), Port: 3001}
	seed3 := Seed{IP: net.ParseIP("10.1.1.3"), Port: 3002}
	repo.SetSeed(seed1)
	repo.SetSeed(seed2)
	repo.SetSeed(seed3)
	assert.Equal(t, 3, len(repo.seeds))
	st1 := NewState("0.0", OkStatus, PeerPosition{}, "0.0.0", 0.0, 0, 5)
	st2 := NewState("0.0", OkStatus, PeerPosition{}, "0.0.0", 0.0, 0, 6)
	st3 := NewState("0.0", OkStatus, PeerPosition{}, "0.0.0", 0.0, 0, 7)
	p1 := NewPeerDetailed([]byte("key1"), seed1.IP, seed1.Port, time.Now(), st1)
	p2 := NewPeerDetailed([]byte("key2"), seed2.IP, seed1.Port, time.Now(), st2)
	p3 := NewPeerDetailed([]byte("key3"), seed3.IP, seed1.Port, time.Now(), st3)
	repo.SetPeer(p1)
	repo.SetPeer(p2)
	repo.SetPeer(p3)
	assert.Equal(t, 3, len(repo.peers))
	avg, _ := sdc.Average()
	assert.Equal(t, 6, avg)

}

/*
Scenario: check DiscoveredPeer
	Given a peer with 3 seed  / 5 peers on the repo
	When DiscoveredPeer
	Then return 5
*/

func TestDiscoveredPeer(t *testing.T) {
	repo := new(mockPeerRepository)
	initP := NewStartupPeer([]byte("key"), net.ParseIP("127.0.0.1"), 3000, "0.0", PeerPosition{})
	repo.SetPeer(initP)
	seed1 := Seed{IP: net.ParseIP("10.1.1.1"), Port: 3000}
	seed2 := Seed{IP: net.ParseIP("10.1.1.2"), Port: 3001}
	seed3 := Seed{IP: net.ParseIP("10.1.1.3"), Port: 3002}
	repo.SetSeed(seed1)
	repo.SetSeed(seed2)
	repo.SetSeed(seed3)
	assert.Equal(t, 3, len(repo.seeds))
	st1 := NewState("0.0", OkStatus, PeerPosition{}, "0.0.0", 0.0, 0, 5)
	st2 := NewState("0.0", OkStatus, PeerPosition{}, "0.0.0", 0.0, 0, 5)
	st3 := NewState("0.0", OkStatus, PeerPosition{}, "0.0.0", 0.0, 0, 5)
	p1 := NewPeerDetailed([]byte("key1"), seed1.IP, seed1.Port, time.Now(), st1)
	p2 := NewPeerDetailed([]byte("key2"), seed2.IP, seed1.Port, time.Now(), st2)
	p3 := NewPeerDetailed([]byte("key3"), seed3.IP, seed1.Port, time.Now(), st3)
	repo.SetPeer(p1)
	repo.SetPeer(p2)
	repo.SetPeer(p3)
	p4 := NewPeerDetailed([]byte("key4"), net.ParseIP("185.123.4.9"), 4000, time.Now(), st1)
	repo.SetPeer(p4)
	assert.Equal(t, 5, len(repo.peers))

	sdc := NewSeedDiscoveryCounter(repo)
	dn, err := sdc.Average()
	assert.Equal(t, nil, err)
	assert.Equal(t, 5, dn)

}

type mockPeerRepository struct {
	peers []Peer
	seeds []Seed
}

func (r *mockPeerRepository) CountKnownPeers() (int, error) {
	return len(r.peers), nil
}

func (r *mockPeerRepository) GetOwnedPeer() (p Peer, err error) {
	for _, p := range r.peers {
		if p.IsOwned() {
			return p, nil
		}
	}
	return
}

func (r *mockPeerRepository) SetPeer(p Peer) error {
	if r.containsPeer(p) {
		return r.UpdatePeer(p)
	}
	r.peers = append(r.peers, p)
	return nil
}

func (r *mockPeerRepository) SetSeed(s Seed) error {
	r.seeds = append(r.seeds, s)
	return nil
}

func (r *mockPeerRepository) ListKnownPeers() ([]Peer, error) {
	return r.peers, nil
}

func (r *mockPeerRepository) ListSeedPeers() ([]Seed, error) {
	return r.seeds, nil
}

func (r *mockPeerRepository) GetPeerByIP(ip net.IP) (p Peer, err error) {
	for i := 0; i < len(r.peers); i++ {
		if string(ip) == string(r.peers[i].IP()) {
			return r.peers[i], nil
		}
	}
	return
}

func (r *mockPeerRepository) UpdatePeer(peer Peer) error {
	for _, p := range r.peers {
		if string(p.PublicKey()) == string(peer.PublicKey()) {
			p = peer
			break
		}
	}
	return nil
}

func (r *mockPeerRepository) containsPeer(peer Peer) bool {
	mPeers := make(map[string]Peer, 0)
	for _, p := range r.peers {
		mPeers[hex.EncodeToString(p.PublicKey())] = peer
	}

	_, exist := mPeers[hex.EncodeToString(peer.PublicKey())]
	return exist
}

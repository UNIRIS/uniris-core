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
	repo.AddSeed(seed1)
	repo.AddSeed(seed2)
	repo.AddSeed(seed3)
	assert.Equal(t, 3, len(repo.seeds))

	st1 := NewPeerAppState("0.0", OkStatus, PeerPosition{}, "0.0.0", 0.0, 0, 5)
	st2 := NewPeerAppState("0.0", OkStatus, PeerPosition{}, "0.0.0", 0.0, 0, 6)
	st3 := NewPeerAppState("0.0", OkStatus, PeerPosition{}, "0.0.0", 0.0, 0, 7)

	p1 := NewDiscoveredPeer(
		NewPeerIdentity(seed1.IP, seed1.Port, []byte("key1")),
		NewPeerHeartbeatState(time.Now(), 0),
		st1,
	)

	p2 := NewDiscoveredPeer(
		NewPeerIdentity(seed2.IP, seed2.Port, []byte("key2")),
		NewPeerHeartbeatState(time.Now(), 0),
		st2,
	)

	p3 := NewDiscoveredPeer(
		NewPeerIdentity(seed3.IP, seed3.Port, []byte("key3")),
		NewPeerHeartbeatState(time.Now(), 0),
		st3,
	)

	repo.AddPeer(p1)
	repo.AddPeer(p2)
	repo.AddPeer(p3)
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
	repo.AddPeer(initP)
	seed1 := Seed{IP: net.ParseIP("10.1.1.1"), Port: 3000}
	seed2 := Seed{IP: net.ParseIP("10.1.1.2"), Port: 3001}
	seed3 := Seed{IP: net.ParseIP("10.1.1.3"), Port: 3002}
	repo.AddSeed(seed1)
	repo.AddSeed(seed2)
	repo.AddSeed(seed3)
	assert.Equal(t, 3, len(repo.seeds))

	st1 := NewPeerAppState("0.0", OkStatus, PeerPosition{}, "0.0.0", 0.0, 0, 5)
	st2 := NewPeerAppState("0.0", OkStatus, PeerPosition{}, "0.0.0", 0.0, 0, 5)
	st3 := NewPeerAppState("0.0", OkStatus, PeerPosition{}, "0.0.0", 0.0, 0, 5)

	p1 := NewDiscoveredPeer(
		NewPeerIdentity(seed1.IP, seed1.Port, []byte("key1")),
		NewPeerHeartbeatState(time.Now(), 0),
		st1,
	)

	p2 := NewDiscoveredPeer(
		NewPeerIdentity(seed2.IP, seed2.Port, []byte("key2")),
		NewPeerHeartbeatState(time.Now(), 0),
		st2,
	)

	p3 := NewDiscoveredPeer(
		NewPeerIdentity(seed3.IP, seed3.Port, []byte("key3")),
		NewPeerHeartbeatState(time.Now(), 0),
		st3,
	)

	repo.AddPeer(p1)
	repo.AddPeer(p2)
	repo.AddPeer(p3)

	p4 := NewDiscoveredPeer(
		NewPeerIdentity(net.ParseIP("185.123.4.9"), 4000, []byte("key4")),
		NewPeerHeartbeatState(time.Now(), 0),
		st1)

	repo.AddPeer(p4)
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
		if p.Owned() {
			return p, nil
		}
	}
	return
}

func (r *mockPeerRepository) AddPeer(p Peer) error {
	if r.containsPeer(p) {
		return r.UpdatePeer(p)
	}
	r.peers = append(r.peers, p)
	return nil
}

func (r *mockPeerRepository) AddSeed(s Seed) error {
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
		if string(ip) == string(r.peers[i].Identity().IP()) {
			return r.peers[i], nil
		}
	}
	return
}

func (r *mockPeerRepository) UpdatePeer(peer Peer) error {
	for _, p := range r.peers {
		if p.Identity().PublicKey().Equals(p.Identity().PublicKey()) {
			p = peer
			break
		}
	}
	return nil
}

func (r *mockPeerRepository) containsPeer(peer Peer) bool {
	mPeers := make(map[string]Peer, 0)
	for _, p := range r.peers {
		mPeers[hex.EncodeToString(p.Identity().PublicKey())] = peer
	}

	_, exist := mPeers[hex.EncodeToString(peer.Identity().PublicKey())]
	return exist
}

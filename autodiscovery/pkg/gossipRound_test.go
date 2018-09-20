package discovery

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
Scenario: Picks a random peer
	Given a list of peer
	When we want to pick a random peer
	Then we get a random peer
*/
func TestRandomPeer(t *testing.T) {
	p1 := NewPeerDigest([]byte("key"), net.ParseIP("127.0.0.1"), 3000)
	p2 := NewPeerDigest([]byte("key2"), net.ParseIP("10.0.0.1"), 3000)
	peers := []Peer{p1, p2}

	g := &GossipRound{
		initiator:  Peer{},
		knownPeers: peers,
		seedPeers:  []Seed{},
	}
	p := g.randomPeer()
	assert.NotNil(t, p)
}

/*
Scenario: Picks a random seed
	Given a list of seeds
	When we want to pick a random seed
	Then we get a random seed
*/
func TestRandomSeed(t *testing.T) {
	s1 := Seed{IP: net.ParseIP("127.0.0.1"), Port: 3000}
	s2 := Seed{IP: net.ParseIP("30.0.0.0"), Port: 3000}

	g := &GossipRound{
		initiator:  Peer{},
		knownPeers: []Peer{},
		seedPeers:  []Seed{s1, s2},
	}
	s := g.randomSeed()
	assert.NotNil(t, s)
}

/*
Scenario: Starts a gossip round without seeds
	Given a initiator peer, a empty list of seeds
	When we starts a gossip round
	Then an error is returned
*/
func TestGossipWithoutSeeds(t *testing.T) {
	_, err := NewGossipRound(Peer{}, []Peer{}, []Seed{})
	assert.Error(t, err, ErrEmptySeed)
}

/*
Scenario: Selects peers from seed and known peers
	Given a list of peers and seeds
	When we want select peers to gossip
	Then we get a random seed and a random peer (exluding ourself)
*/
func TestSelectPeers(t *testing.T) {

	s1 := Seed{IP: net.ParseIP("30.0.50.100"), Port: 3000}

	p1 := NewStartupPeer([]byte("key"), net.ParseIP("127.0.0.1"), 3000, "1.0", PeerPosition{}, 1)
	p2 := NewPeerDigest([]byte("key2"), net.ParseIP("10.0.0.1"), 3000)

	r, _ := NewGossipRound(Peer{}, []Peer{p1, p2}, []Seed{s1})

	peers, err := r.SelectPeers()
	assert.Nil(t, err)
	assert.NotNil(t, peers)
	assert.NotEmpty(t, peers)
	assert.Equal(t, 2, len(peers))
	assert.Equal(t, "30.0.50.100", peers[0].IP().String())
	assert.Equal(t, "10.0.0.1", peers[1].IP().String())
}

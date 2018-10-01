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
	p1 := &peer{identity: NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, PublicKey("key"))}
	p2 := &peer{identity: NewPeerIdentity(net.ParseIP("127.0.0.2"), 3000, PublicKey("key2"))}
	peers := []Peer{p1, p2}

	g := &GossipRound{
		initiator:         &peer{},
		reacheablePeers:   peers,
		seedPeers:         []Seed{},
		unreacheablePeers: []Peer{},
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
		initiator:         &peer{},
		reacheablePeers:   []Peer{},
		seedPeers:         []Seed{s1, s2},
		unreacheablePeers: []Peer{},
	}
	s := g.randomSeed()
	assert.NotNil(t, s)
}

/*
Scenario: Picks a random unreacheablepeer
	Given a list of unreacheablePeer
	When we want to pick a random unreacheablePeer
	Then we get a random unreacheablePeer
*/
func TestRandomUnreacheablePeer(t *testing.T) {
	p1 := &peer{identity: NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, PublicKey("key"))}
	p2 := &peer{identity: NewPeerIdentity(net.ParseIP("127.0.0.2"), 3000, PublicKey("key2"))}
	peers := []Peer{p1, p2}

	g := &GossipRound{
		initiator:         &peer{},
		reacheablePeers:   []Peer{},
		seedPeers:         []Seed{},
		unreacheablePeers: peers,
	}
	p := g.randomUnreacheablePeer()
	assert.NotNil(t, p)
}

/*
Scenario: Starts a gossip round without seeds
	Given a initiator peer, a empty list of seeds
	When we starts a gossip round
	Then an error is returned
*/
func TestGossipWithoutSeeds(t *testing.T) {
	_, err := NewGossipRound(&peer{}, []Peer{}, []Seed{}, []Peer{})
	assert.Error(t, err, ErrEmptySeed)
}

/*
Scenario: Selects peers from seed , reacheable peers and unreacheable peers
	Given a list of reacheable peers, seeds and unreacheable peers
	When we want select peers to gossip
	Then we get a random seed and a random peer (exluding ourself) and a random unreacheable peer
*/
func TestSelectPeers(t *testing.T) {

	s1 := Seed{IP: net.ParseIP("30.0.50.100"), Port: 3000}

	p1 := NewStartupPeer(PublicKey("key"), net.ParseIP("127.0.0.1"), 3000, "1.0", PeerPosition{})
	p2 := &peer{identity: NewPeerIdentity(net.ParseIP("10.0.0.1"), 3000, PublicKey("key2"))}
	p3 := &peer{identity: NewPeerIdentity(net.ParseIP("10.0.0.2"), 3000, PublicKey("key3"))}

	r, _ := NewGossipRound(&peer{}, []Peer{p1, p2}, []Seed{s1}, []Peer{p3})

	peers, err := r.SelectPeers()
	assert.Nil(t, err)
	assert.NotNil(t, peers)
	assert.NotEmpty(t, peers)
	assert.Equal(t, 3, len(peers))
	assert.Equal(t, "30.0.50.100", peers[0].Identity().IP().String())
	assert.Equal(t, "10.0.0.1", peers[1].Identity().IP().String())
	assert.Equal(t, "10.0.0.2", peers[2].Identity().IP().String())
}

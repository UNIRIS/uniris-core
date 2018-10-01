package gossip

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	discovery "github.com/uniris/uniris-core/autodiscovery/pkg"
)

/*
Scenario: Picks a random peer
	Given a list of peer
	When we want to pick a random peer
	Then we get a random peer
*/
func TestRandomPeer(t *testing.T) {
	p1 := discovery.NewPeerDigest(
		discovery.NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, []byte("key")),
		discovery.NewPeerHeartbeatState(time.Now(), 0))

	p2 := discovery.NewPeerDigest(
		discovery.NewPeerIdentity(net.ParseIP("10.0.0.1"), 3000, []byte("key2")),
		discovery.NewPeerHeartbeatState(time.Now(), 0))

	peers := []discovery.Peer{p1, p2}

	c := Cycle{}
	p := c.randomPeer(peers)
	assert.NotNil(t, p)
}

/*
Scenario: Picks a random seed
	Given a list of seeds
	When we want to pick a random seed
	Then we get a random seed
*/
func TestRandomSeed(t *testing.T) {
	s1 := discovery.Seed{IP: net.ParseIP("127.0.0.1"), Port: 3000}
	s2 := discovery.Seed{IP: net.ParseIP("30.0.0.0"), Port: 3000}

	g := Cycle{seedPeers: []discovery.Seed{s1, s2}}
	s := g.randomSeed()
	assert.NotNil(t, s)
}

/*
Scenario: Starts a gossip round without seeds
	Given a initiator peer, a empty list of seeds
	When we starts a gossip round
	Then an error is returned
*/
func TestCycleWithoutSeeds(t *testing.T) {
	p1 := discovery.NewPeerDigest(
		discovery.NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, []byte("key")),
		discovery.NewPeerHeartbeatState(time.Now(), 0))

	_, err := NewGossipCycle(p1, []discovery.Peer{}, []discovery.Seed{}, nil)
	assert.Error(t, err, ErrEmptySeed)
}

/*
Scenario: Selects peers from seed and known peers
	Given a list of peers and seeds
	When we want select peers to gossip
	Then we get a random seed and a random peer (exluding ourself)
*/
func TestSelectPeers(t *testing.T) {

	s1 := discovery.Seed{IP: net.ParseIP("30.0.50.100"), Port: 3000}

	p1 := discovery.NewStartupPeer([]byte("key"), net.ParseIP("127.0.0.1"), 3000, "1.0", discovery.PeerPosition{})

	p2 := discovery.NewPeerDigest(
		discovery.NewPeerIdentity(net.ParseIP("10.0.0.1"), 3000, []byte("key")),
		discovery.NewPeerHeartbeatState(time.Now(), 0))

	c, _ := NewGossipCycle(p1, []discovery.Peer{p1, p2}, []discovery.Seed{s1}, nil)

	peers := c.selectPeers()
	assert.NotNil(t, peers)
	assert.NotEmpty(t, peers)
	assert.Equal(t, 2, len(peers))
	assert.Equal(t, "30.0.50.100", peers[0].Identity().IP().String())
	assert.True(t, peers[1].Identity().IP().String() == "127.0.0.1" || peers[1].Identity().IP().String() == "10.0.0.1")
}

/*
Scenario: Run a gossip cycle
	Given a initator and a target
	When we create a round associated to a cycle
	Then we run it and get some discovered peers
*/
func TestRun(t *testing.T) {
	init := discovery.NewStartupPeer([]byte("key"), net.ParseIP("10.0.0.1"), 3000, "1.0", discovery.PeerPosition{})

	kp1 := discovery.NewPeerDigest(
		discovery.NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, []byte("key2")),
		discovery.NewPeerHeartbeatState(time.Now(), 0),
	)

	seeds := []discovery.Seed{discovery.Seed{IP: net.ParseIP("20.0.0.1"), Port: 3000, PublicKey: []byte("key3")}}

	c, err := NewGossipCycle(init, []discovery.Peer{kp1}, seeds, mockMessenger{})
	assert.Nil(t, err)

	go c.Run()

	newP := make([]discovery.Peer, 0)
	for p := range c.result.discoveries {
		newP = append(newP, p)
	}

	assert.NotEmpty(t, newP)
	assert.Equal(t, 2, len(newP))

	//Peer retrieved from the kp1
	assert.Equal(t, "dKey1", newP[0].Identity().PublicKey().String())

	//Peer retreived from the seed1
	assert.Equal(t, "dKey1", newP[1].Identity().PublicKey().String())

	assert.Empty(t, c.result.errors)
	assert.Empty(t, c.result.unreachables)
}

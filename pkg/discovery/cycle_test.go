package discovery

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

/*
Scenario: Picks a random peer from an only list of seeds
	Given a list of seeds
	When we want to pick a random seed
	Then we get a random seed
*/
func TestSelectPeersWithOnlySeeds(t *testing.T) {
	seeds := []PeerIdentity{
		NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "key"),
	}

	c := Cycle{}
	self := NewSelfPeer("key", net.ParseIP("10.0.0.1"), 3000, "1.0", 30.0, 12.0)
	pp, err := c.selectRandomPeers(self, seeds, []PeerIdentity{}, []PeerIdentity{})
	assert.Nil(t, err)
	assert.NotNil(t, pp)
	assert.Equal(t, 1, len(pp))
	assert.Equal(t, "127.0.0.1", pp[0].IP().String())
}

/*
Scenario: Pick two random peers (seed and a reachable peer)
	Given a list of seeds and a list of reachable peers
	When we want select peers
	Then we get a random seed and a random reachable peer
*/
func TestSelectPeerWithSomeReachablePeers(t *testing.T) {
	seeds := []PeerIdentity{
		NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "key"),
	}

	self := NewSelfPeer("key", net.ParseIP("10.0.0.1"), 3000, "1.0", 30.0, 12.0)

	c := Cycle{}

	reachP := []PeerIdentity{NewPeerIdentity(net.ParseIP("20.0.0.1"), 3000, "key")}

	pp, err := c.selectRandomPeers(self, seeds, reachP, []PeerIdentity{})
	assert.Nil(t, err)
	assert.NotNil(t, pp)
	assert.Equal(t, 2, len(pp))
	assert.Equal(t, "127.0.0.1", pp[0].IP().String())
	assert.Equal(t, "20.0.0.1", pp[1].IP().String())
}

/*
Scenario: Pick two random peers (seed and a reachable peer) but ourself is in the reachable peer
	Given a list of seeds and a list of reachable peers (including ourself)
	When we want select peers
	Then we get a random seed and no reachable peers
*/
func TestSelectPeerWithOurselfAsReachable(t *testing.T) {
	seeds := []PeerIdentity{
		NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "key"),
	}

	me := NewSelfPeer("key", net.ParseIP("10.0.0.1"), 3000, "1.0", 30.0, 12.0)

	c := Cycle{}

	reachP := []PeerIdentity{me.Identity()}

	pp, err := c.selectRandomPeers(me, seeds, reachP, []PeerIdentity{})
	assert.Nil(t, err)
	assert.NotNil(t, pp)
	assert.Equal(t, 1, len(pp))
	assert.Equal(t, "127.0.0.1", pp[0].IP().String())
}

/*
Scenario: Pick two random peers (seed and a unreachable peer)
	Given a list of seeds and a list of unreachable peers
	When we want select peers
	Then we get a random seed and a random unreachable peer
*/
func TestSelectPeerWithSomeUnReachablePeers(t *testing.T) {
	seeds := []PeerIdentity{
		NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "key"),
	}

	p1 := NewPeerDigest(
		NewPeerIdentity(net.ParseIP("20.0.0.1"), 3000, "key"),
		NewPeerHeartbeatState(time.Now(), 0))

	selfator := NewSelfPeer("key", net.ParseIP("10.0.0.1"), 3000, "1.0", 30.0, 12.0)

	c := Cycle{}

	unreachP := []PeerIdentity{p1.Identity()}

	pp, err := c.selectRandomPeers(selfator, seeds, []PeerIdentity{}, unreachP)
	assert.Nil(t, err)
	assert.NotNil(t, pp)
	assert.Equal(t, 2, len(pp))
	assert.Equal(t, "127.0.0.1", pp[0].IP().String())
	assert.Equal(t, "20.0.0.1", pp[1].IP().String())
}

/*
Scenario: Pick random peers (seed, reachable and unreachable)
	Given a list of seeds, a list of reachables peers and unreachables peers
	When we select randomly peers
	Then we get a seed, a reachable and an unreachable peer
*/
func TestSelectPeersFully(t *testing.T) {
	seeds := []PeerIdentity{
		NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "key"),
	}

	p1 := NewPeerIdentity(net.ParseIP("20.0.0.1"), 3000, "key")
	p2 := NewPeerIdentity(net.ParseIP("21.0.0.1"), 3000, "key")
	p3 := NewPeerIdentity(net.ParseIP("22.0.0.1"), 3000, "key")
	p4 := NewPeerIdentity(net.ParseIP("23.0.0.1"), 3000, "key")

	selfator := NewSelfPeer("key", net.ParseIP("10.0.0.1"), 3000, "1.0", 30.0, 12.0)
	c := Cycle{}

	reachP := []PeerIdentity{p1, p2}
	unreachP := []PeerIdentity{p3, p4}

	pp, err := c.selectRandomPeers(selfator, seeds, reachP, unreachP)
	assert.Nil(t, err)
	assert.NotNil(t, pp)
	assert.Equal(t, 3, len(pp))
	assert.True(t, pp[0].IP().String() == "127.0.0.1" || pp[0].IP().String() == "30.0.0.1")
	assert.True(t, pp[1].IP().String() == "20.0.0.1" || pp[1].IP().String() == "21.0.0.1")
	assert.True(t, pp[2].IP().String() == "22.0.0.1" || pp[2].IP().String() == "23.0.0.1")

}

/*
Scenario: Run a gossip cycle
	Given a selfator and a target
	When we create a round associated to a cycle
	Then we run it and get some discovered peers
*/
func TestRun(t *testing.T) {
	self := NewSelfPeer("key", net.ParseIP("10.0.0.1"), 3000, "1.0", 30.0, 12.0)

	kp1 := NewPeerDigest(
		NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "key2"),
		NewPeerHeartbeatState(time.Now(), 0))

	seeds := []PeerIdentity{
		NewPeerIdentity(net.ParseIP("20.0.0.1"), 3000, "key3"),
	}

	c := Cycle{}

	assert.Nil(t, c.run(self, mockClient{}, seeds, []Peer{kp1}, []PeerIdentity{kp1.Identity()}, []PeerIdentity{}))

	assert.Len(t, c.Discoveries, 2)
	assert.Len(t, c.Reaches, 2)

	//Peer retrieved from the kp1
	assert.Equal(t, "dKey1", c.Discoveries[0].Identity().PublicKey())

	//Peer retreived from the seed1
	assert.Equal(t, "dKey1", c.Discoveries[1].Identity().PublicKey())
}

type mockClient struct{}

func (m mockClient) SendSyn(target PeerIdentity, known []Peer) (unknown []Peer, new []Peer, err error) {
	tar := NewSelfPeer("uKey1", net.ParseIP("200.18.186.39"), 3000, "1.1", 40.4, 2.50)

	hb := NewPeerHeartbeatState(time.Now(), 0)
	as := NewPeerAppState("1.0", OkPeerStatus, 50.1, 22.1, "", 0, 1, 0)

	np1 := NewDiscoveredPeer(
		NewPeerIdentity(net.ParseIP("35.200.100.2"), 3000, "dKey1"),
		hb, as,
	)

	newPeers := []Peer{np1}
	unknownPeers := []Peer{tar}
	return unknownPeers, newPeers, nil
}

func (m mockClient) SendAck(target PeerIdentity, requested []Peer) error {
	return nil
}

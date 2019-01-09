package gossip

import (
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	discovery "github.com/uniris/uniris-core/autodiscovery/pkg"
)

/*
Scenario: Select peers without seeds
	Given no seeds
	When we i want select peers
	Then we get an error
*/
func TestSelectPeersWithoutSeeds(t *testing.T) {
	c := Cycle{}
	_, err := c.SelectPeers(nil, nil, nil)
	assert.Error(t, err, ErrEmptySeed)
}

/*
Scenario: Picks a random peer from an only list of seeds
	Given a list of seeds
	When we want to pick a random seed
	Then we get a random seed
*/
func TestSelectPeersWithOnlyPeers(t *testing.T) {
	seeds := []discovery.Seed{
		discovery.Seed{IP: net.ParseIP("127.0.0.1"), Port: 3000},
	}

	c := Cycle{
		initator: discovery.NewStartupPeer("key", net.ParseIP("10.0.0.1"), 3000, "1.0", discovery.PeerPosition{}),
	}
	pp, err := c.SelectPeers(seeds, []discovery.Peer{}, []discovery.Peer{})
	assert.Nil(t, err)
	assert.NotNil(t, pp)
	assert.Equal(t, 1, len(pp))
	assert.Equal(t, "127.0.0.1", pp[0].Identity().IP().String())
}

/*
Scenario: Picks a random peer from an only list of seeds excluding ourself
	Given a list of seeds including ourself
	When we want to pick a random seed
	Then we get a random seed
*/
func TestSelectPeersWithOnlyPeersExcludingOurself(t *testing.T) {
	seeds := []discovery.Seed{
		discovery.Seed{IP: net.ParseIP("127.0.0.1"), Port: 3000},
	}

	c := Cycle{
		initator: discovery.NewStartupPeer("key", net.ParseIP("127.0.0.1"), 3000, "1.0", discovery.PeerPosition{}),
	}
	pp, err := c.SelectPeers(seeds, []discovery.Peer{}, []discovery.Peer{})
	assert.Nil(t, err)
	assert.NotNil(t, pp)
	assert.Empty(t, pp)
}

/*
Scenario: Pick two random peers (seed and a reachable peer)
	Given a list of seeds and a list of reachable peers
	When we want select peers
	Then we get a random seed and a random reachable peer
*/
func TestSelectPeerWithSomeReachablePeers(t *testing.T) {
	seeds := []discovery.Seed{
		discovery.Seed{IP: net.ParseIP("127.0.0.1"), Port: 3000},
	}

	p1 := discovery.NewPeerDigest(
		discovery.NewPeerIdentity(net.ParseIP("20.0.0.1"), 3000, "key"),
		discovery.NewPeerHeartbeatState(time.Now(), 0))

	c := Cycle{
		initator: discovery.NewStartupPeer("key", net.ParseIP("10.0.0.1"), 3000, "1.0", discovery.PeerPosition{}),
	}

	reachP := []discovery.Peer{p1}

	pp, err := c.SelectPeers(seeds, reachP, []discovery.Peer{})
	assert.Nil(t, err)
	assert.NotNil(t, pp)
	assert.Equal(t, 2, len(pp))
	assert.Equal(t, "127.0.0.1", pp[0].Identity().IP().String())
	assert.Equal(t, "20.0.0.1", pp[1].Identity().IP().String())
}

/*
Scenario: Pick two random peers (seed and a reachable peer) but ourself is in the reachable peer
	Given a list of seeds and a list of reachable peers (including ourself)
	When we want select peers
	Then we get a random seed and no reachable peers
*/
func TestSelectPeerWithOurselfAsReachable(t *testing.T) {
	seeds := []discovery.Seed{
		discovery.Seed{IP: net.ParseIP("127.0.0.1"), Port: 3000},
	}

	me := discovery.NewStartupPeer("key", net.ParseIP("10.0.0.1"), 3000, "1.0", discovery.PeerPosition{})

	c := Cycle{
		initator: me,
	}

	reachP := []discovery.Peer{me}

	pp, err := c.SelectPeers(seeds, reachP, []discovery.Peer{})
	assert.Nil(t, err)
	assert.NotNil(t, pp)
	assert.Equal(t, 1, len(pp))
	assert.Equal(t, "127.0.0.1", pp[0].Identity().IP().String())
}

/*
Scenario: Pick two random peers (seed and a unreachable peer)
	Given a list of seeds and a list of unreachable peers
	When we want select peers
	Then we get a random seed and a random unreachable peer
*/
func TestSelectPeerWithSomeUnReachablePeers(t *testing.T) {
	seeds := []discovery.Seed{
		discovery.Seed{IP: net.ParseIP("127.0.0.1"), Port: 3000},
	}

	p1 := discovery.NewPeerDigest(
		discovery.NewPeerIdentity(net.ParseIP("20.0.0.1"), 3000, "key"),
		discovery.NewPeerHeartbeatState(time.Now(), 0))

	c := Cycle{
		initator: discovery.NewStartupPeer("key", net.ParseIP("10.0.0.1"), 3000, "1.0", discovery.PeerPosition{}),
	}

	unreachP := []discovery.Peer{p1}

	pp, err := c.SelectPeers(seeds, []discovery.Peer{}, unreachP)
	assert.Nil(t, err)
	assert.NotNil(t, pp)
	assert.Equal(t, 2, len(pp))
	assert.Equal(t, "127.0.0.1", pp[0].Identity().IP().String())
	assert.Equal(t, "20.0.0.1", pp[1].Identity().IP().String())
}

/*
Scenario: Pick random peers (seed, reachable and unreachable)
	Given a list of seeds, a list of reachables peers and unreachables peers
	When we select randomly peers
	Then we get a seed, a reachable and an unreachable peer
*/
func TestSelectPeersFully(t *testing.T) {
	seeds := []discovery.Seed{
		discovery.Seed{IP: net.ParseIP("127.0.0.1"), Port: 3000},
		discovery.Seed{IP: net.ParseIP("30.0.0.1"), Port: 3000},
	}

	p1 := discovery.NewPeerDigest(
		discovery.NewPeerIdentity(net.ParseIP("20.0.0.1"), 3000, "key"),
		discovery.NewPeerHeartbeatState(time.Now(), 0))
	p2 := discovery.NewPeerDigest(
		discovery.NewPeerIdentity(net.ParseIP("21.0.0.1"), 3000, "key"),
		discovery.NewPeerHeartbeatState(time.Now(), 0))
	p3 := discovery.NewPeerDigest(
		discovery.NewPeerIdentity(net.ParseIP("22.0.0.1"), 3000, "key"),
		discovery.NewPeerHeartbeatState(time.Now(), 0))
	p4 := discovery.NewPeerDigest(
		discovery.NewPeerIdentity(net.ParseIP("23.0.0.1"), 3000, "key"),
		discovery.NewPeerHeartbeatState(time.Now(), 0))

	c := Cycle{
		initator: discovery.NewStartupPeer("key", net.ParseIP("10.0.0.1"), 3000, "1.0", discovery.PeerPosition{}),
	}

	reachP := []discovery.Peer{p1, p2}
	unreachP := []discovery.Peer{p3, p4}

	pp, err := c.SelectPeers(seeds, reachP, unreachP)
	assert.Nil(t, err)
	assert.NotNil(t, pp)
	assert.Equal(t, 3, len(pp))
	assert.True(t, pp[0].Identity().IP().String() == "127.0.0.1" || pp[0].Identity().IP().String() == "30.0.0.1")
	assert.True(t, pp[1].Identity().IP().String() == "20.0.0.1" || pp[1].Identity().IP().String() == "21.0.0.1")
	assert.True(t, pp[2].Identity().IP().String() == "22.0.0.1" || pp[2].Identity().IP().String() == "23.0.0.1")

}

/*
Scenario: Run a gossip cycle
	Given a initator and a target
	When we create a round associated to a cycle
	Then we run it and get some discovered peers
*/
func TestRun(t *testing.T) {
	init := discovery.NewStartupPeer("key", net.ParseIP("10.0.0.1"), 3000, "1.0", discovery.PeerPosition{})

	kp1 := discovery.NewPeerDigest(
		discovery.NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "key2"),
		discovery.NewPeerHeartbeatState(time.Now(), 0),
	)

	seeds := []discovery.Seed{discovery.Seed{IP: net.ParseIP("20.0.0.1"), Port: 3000, PublicKey: "key3"}}

	c := NewGossipCycle(init, mockMessenger{})

	pp, err := c.SelectPeers(seeds, []discovery.Peer{kp1}, []discovery.Peer{})
	assert.Nil(t, err)

	var wg sync.WaitGroup
	wg.Add(2)

	go c.Run(init, pp, []discovery.Peer{kp1})

	newP := make([]discovery.Peer, 0)
	go func() {
		for p := range c.result.discoveries {
			newP = append(newP, p)
			wg.Done()
		}
	}()

	go func() {
		for range c.result.reaches {
		}
	}()

	wg.Wait()
	close(c.result.discoveries)
	close(c.result.unreachables)
	close(c.result.errors)
	close(c.result.reaches)

	assert.NotEmpty(t, newP)
	assert.Equal(t, 2, len(newP))

	//Peer retrieved from the kp1
	assert.Equal(t, "dKey1", newP[0].Identity().PublicKey())

	//Peer retreived from the seed1
	assert.Equal(t, "dKey1", newP[1].Identity().PublicKey())

	assert.Empty(t, c.result.errors)
	assert.Empty(t, c.result.unreachables)
}

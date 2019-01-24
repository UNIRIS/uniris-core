package gossip

import (
	"log"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	uniris "github.com/uniris/uniris-core/pkg"
)

/*
Scenario: Select peers without seeds
	Given no seeds
	When we i want select peers
	Then we get an error
*/
func TestSelectPeersWithoutSeeds(t *testing.T) {
	c := cycle{}
	_, err := c.selectRandomPeers(nil, nil, nil)
	assert.Error(t, err, ErrEmptySeed)
}

/*
Scenario: Picks a random peer from an only list of seeds
	Given a list of seeds
	When we want to pick a random seed
	Then we get a random seed
*/
func TestSelectPeersWithOnlyPeers(t *testing.T) {
	seeds := []uniris.Seed{
		uniris.Seed{
			PeerIdentity: uniris.NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "key"),
		},
	}

	c := cycle{
		initator: uniris.NewLocalPeer("key", net.ParseIP("10.0.0.1"), 3000, "1.0", 30.0, 12.0),
	}
	pp, err := c.selectRandomPeers(seeds, []uniris.Peer{}, []uniris.Peer{})
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
	seeds := []uniris.Seed{
		uniris.Seed{
			PeerIdentity: uniris.NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "key"),
		},
	}

	c := cycle{
		initator: uniris.NewLocalPeer("key", net.ParseIP("127.0.0.1"), 3000, "1.0", 30.0, 12.0),
	}
	pp, err := c.selectRandomPeers(seeds, []uniris.Peer{}, []uniris.Peer{})
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
	seeds := []uniris.Seed{
		uniris.Seed{
			PeerIdentity: uniris.NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "key"),
		},
	}

	p1 := uniris.NewPeerDigest(
		uniris.NewPeerIdentity(net.ParseIP("20.0.0.1"), 3000, "key"),
		uniris.NewPeerHeartbeatState(time.Now(), 0))

	c := cycle{
		initator: uniris.NewLocalPeer("key", net.ParseIP("10.0.0.1"), 3000, "1.0", 30.0, 12.0),
	}

	reachP := []uniris.Peer{p1}

	pp, err := c.selectRandomPeers(seeds, reachP, []uniris.Peer{})
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
	seeds := []uniris.Seed{
		uniris.Seed{
			PeerIdentity: uniris.NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "key"),
		},
	}

	me := uniris.NewLocalPeer("key", net.ParseIP("10.0.0.1"), 3000, "1.0", 30.0, 12.0)

	c := cycle{
		initator: me,
	}

	reachP := []uniris.Peer{me}

	pp, err := c.selectRandomPeers(seeds, reachP, []uniris.Peer{})
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
	seeds := []uniris.Seed{
		uniris.Seed{
			PeerIdentity: uniris.NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "key"),
		},
	}

	p1 := uniris.NewPeerDigest(
		uniris.NewPeerIdentity(net.ParseIP("20.0.0.1"), 3000, "key"),
		uniris.NewPeerHeartbeatState(time.Now(), 0))

	c := cycle{
		initator: uniris.NewLocalPeer("key", net.ParseIP("10.0.0.1"), 3000, "1.0", 30.0, 12.0),
	}

	unreachP := []uniris.Peer{p1}

	pp, err := c.selectRandomPeers(seeds, []uniris.Peer{}, unreachP)
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
	seeds := []uniris.Seed{
		uniris.Seed{
			PeerIdentity: uniris.NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "key"),
		},
	}

	p1 := uniris.NewPeerDigest(
		uniris.NewPeerIdentity(net.ParseIP("20.0.0.1"), 3000, "key"),
		uniris.NewPeerHeartbeatState(time.Now(), 0))
	p2 := uniris.NewPeerDigest(
		uniris.NewPeerIdentity(net.ParseIP("21.0.0.1"), 3000, "key"),
		uniris.NewPeerHeartbeatState(time.Now(), 0))
	p3 := uniris.NewPeerDigest(
		uniris.NewPeerIdentity(net.ParseIP("22.0.0.1"), 3000, "key"),
		uniris.NewPeerHeartbeatState(time.Now(), 0))
	p4 := uniris.NewPeerDigest(
		uniris.NewPeerIdentity(net.ParseIP("23.0.0.1"), 3000, "key"),
		uniris.NewPeerHeartbeatState(time.Now(), 0))

	c := cycle{
		initator: uniris.NewLocalPeer("key", net.ParseIP("10.0.0.1"), 3000, "1.0", 30.0, 12.0),
	}

	reachP := []uniris.Peer{p1, p2}
	unreachP := []uniris.Peer{p3, p4}

	pp, err := c.selectRandomPeers(seeds, reachP, unreachP)
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
	init := uniris.NewLocalPeer("key", net.ParseIP("10.0.0.1"), 3000, "1.0", 30.0, 12.0)

	kp1 := uniris.NewPeerDigest(
		uniris.NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "key2"),
		uniris.NewPeerHeartbeatState(time.Now(), 0),
	)

	seeds := []uniris.Seed{uniris.Seed{
		PeerIdentity: uniris.NewPeerIdentity(net.ParseIP("20.0.0.1"), 3000, "key3")}}

	c := newCycle(init, mockMessenger{}, []uniris.Peer{kp1}, []uniris.Peer{})

	var wg sync.WaitGroup
	wg.Add(4)

	go c.run(init, seeds, []uniris.Peer{kp1})

	newP := make([]uniris.Peer, 0)
	go func() {
		for p := range c.discoveryChan {
			log.Print(p.Identity().PublicKey())
			newP = append(newP, p)
			wg.Done()
		}
	}()

	go func() {
		for range c.reachChan {
			wg.Done()
		}
	}()

	wg.Wait()
	close(c.discoveryChan)
	close(c.unreachChan)
	close(c.errChan)
	close(c.reachChan)

	assert.NotEmpty(t, newP)
	assert.Equal(t, 2, len(newP))

	//Peer retrieved from the kp1
	assert.Equal(t, "dKey1", newP[0].Identity().PublicKey())

	//Peer retreived from the seed1
	assert.Equal(t, "dKey1", newP[1].Identity().PublicKey())
}

type mockMessenger struct{}

func (m mockMessenger) SendSyn(source uniris.Peer, target uniris.Peer, known []uniris.Peer) (unknown []uniris.Peer, new []uniris.Peer, err error) {
	tar := uniris.NewLocalPeer("uKey1", net.ParseIP("200.18.186.39"), 3000, "1.1", 40.4, 2.50)

	hb := uniris.NewPeerHeartbeatState(time.Now(), 0)
	as := uniris.NewPeerAppState("1.0", uniris.OkPeerStatus, 50.1, 22.1, "", 0, 1, 0)

	np1 := uniris.NewDiscoveredPeer(
		uniris.NewPeerIdentity(net.ParseIP("35.200.100.2"), 3000, "dKey1"),
		hb, as,
	)

	newPeers := []uniris.Peer{np1}
	unknownPeers := []uniris.Peer{tar}
	return unknownPeers, newPeers, nil
}

func (m mockMessenger) SendAck(source uniris.Peer, target uniris.Peer, requested []uniris.Peer) error {
	return nil
}

type mockPeerInfo struct{}

func (i mockPeerInfo) GeoPosition() (lon float64, lat float64, err error) {
	return 10.0, 30.0, nil
}

func (i mockPeerInfo) FreeDiskSpace() (float64, error) {
	return 200, nil
}

func (i mockPeerInfo) CPULoad() (string, error) {
	return "", nil
}

func (i mockPeerInfo) IP() (net.IP, error) {
	return net.ParseIP("127.0.0.1"), nil
}

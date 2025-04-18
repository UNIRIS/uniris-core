package mem

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	discovery "github.com/uniris/uniris-core/autodiscovery/pkg"
)

/*
Scenario: test add unreacheablePeer on the repo
	Given 2 unreacheable Peer
	When we add them to the repo
	Then unreacheable peer on the repo is 2
*/
func TestAddUnreacheablePeer(t *testing.T) {
	repo := NewRepository()
	st1 := discovery.NewPeerAppState("0.0", discovery.OkStatus, discovery.PeerPosition{}, "0.0.0", 0.0, 0, 5)
	st2 := discovery.NewPeerAppState("0.0", discovery.OkStatus, discovery.PeerPosition{}, "0.0.0", 0.0, 0, 5)
	st3 := discovery.NewPeerAppState("0.0", discovery.OkStatus, discovery.PeerPosition{}, "0.0.0", 0.0, 0, 5)
	p1 := discovery.NewDiscoveredPeer(
		discovery.NewPeerIdentity(net.ParseIP("10.10.0.1"), 3000, "key"),
		discovery.NewPeerHeartbeatState(time.Now(), 0),
		st1,
	)
	p2 := discovery.NewDiscoveredPeer(
		discovery.NewPeerIdentity(net.ParseIP("10.10.0.2"), 3000, "key2"),
		discovery.NewPeerHeartbeatState(time.Now(), 0),
		st2,
	)
	p3 := discovery.NewDiscoveredPeer(
		discovery.NewPeerIdentity(net.ParseIP("10.10.0.3"), 3000, "key3"),
		discovery.NewPeerHeartbeatState(time.Now(), 0),
		st3,
	)
	repo.SetKnownPeer(p1)
	repo.SetKnownPeer(p2)
	repo.SetKnownPeer(p3)
	unp1, _ := repo.GetKnownPeerByIP(net.ParseIP("10.10.0.1"))
	repo.SetUnreachablePeer(unp1.Identity().PublicKey())
	unp2, _ := repo.GetKnownPeerByIP(net.ParseIP("10.10.0.2"))
	repo.SetUnreachablePeer(unp2.Identity().PublicKey())
	l, _ := repo.ListUnreachablePeers()
	assert.Equal(t, 2, len(l))
}

/*
Scenario: test del unreacheablePeer on the repo
	Given 2 unreacheable Peer
	When we add them to the repo + delete one of them
	Then unreacheable peer on the repo is 1
*/
func TestDelUnreacheablePeer(t *testing.T) {
	repo := NewRepository()
	st1 := discovery.NewPeerAppState("0.0", discovery.OkStatus, discovery.PeerPosition{}, "0.0.0", 0.0, 0, 5)
	st2 := discovery.NewPeerAppState("0.0", discovery.OkStatus, discovery.PeerPosition{}, "0.0.0", 0.0, 0, 5)
	st3 := discovery.NewPeerAppState("0.0", discovery.OkStatus, discovery.PeerPosition{}, "0.0.0", 0.0, 0, 5)
	p1 := discovery.NewDiscoveredPeer(
		discovery.NewPeerIdentity(net.ParseIP("10.10.0.1"), 3000, "key"),
		discovery.NewPeerHeartbeatState(time.Now(), 0),
		st1,
	)
	p2 := discovery.NewDiscoveredPeer(
		discovery.NewPeerIdentity(net.ParseIP("10.10.0.2"), 3000, "key2"),
		discovery.NewPeerHeartbeatState(time.Now(), 0),
		st2,
	)
	p3 := discovery.NewDiscoveredPeer(
		discovery.NewPeerIdentity(net.ParseIP("10.10.0.3"), 3000, "key3"),
		discovery.NewPeerHeartbeatState(time.Now(), 0),
		st3,
	)
	repo.SetKnownPeer(p1)
	repo.SetKnownPeer(p2)
	repo.SetKnownPeer(p3)
	unp1, _ := repo.GetKnownPeerByIP(net.ParseIP("10.10.0.1"))
	repo.SetUnreachablePeer(unp1.Identity().PublicKey())
	unp2, _ := repo.GetKnownPeerByIP(net.ParseIP("10.10.0.2"))
	repo.SetUnreachablePeer(unp2.Identity().PublicKey())
	l, _ := repo.ListUnreachablePeers()
	assert.Equal(t, 2, len(l))
	repo.RemoveUnreachablePeer(unp2.Identity().PublicKey())
	l, _ = repo.ListUnreachablePeers()
	assert.Equal(t, 1, len(l))
	lu, _ := repo.ListUnreachablePeers()
	assert.Equal(t, "10.10.0.1", lu[0].Identity().IP().String())
}

/*
Scenario: test reachable peer method
	Given 3 Peer
	When we add one as an unreacheable peer
	Then reacheable peer list on the repo is 2
*/

func TestReacheablePeer(t *testing.T) {
	repo := NewRepository()
	st1 := discovery.NewPeerAppState("0.0", discovery.OkStatus, discovery.PeerPosition{}, "0.0.0", 0.0, 0, 5)
	st2 := discovery.NewPeerAppState("0.0", discovery.OkStatus, discovery.PeerPosition{}, "0.0.0", 0.0, 0, 5)
	st3 := discovery.NewPeerAppState("0.0", discovery.OkStatus, discovery.PeerPosition{}, "0.0.0", 0.0, 0, 5)
	p1 := discovery.NewDiscoveredPeer(
		discovery.NewPeerIdentity(net.ParseIP("10.10.0.1"), 3000, "key"),
		discovery.NewPeerHeartbeatState(time.Now(), 0),
		st1,
	)
	p2 := discovery.NewDiscoveredPeer(
		discovery.NewPeerIdentity(net.ParseIP("10.10.0.2"), 3000, "key2"),
		discovery.NewPeerHeartbeatState(time.Now(), 0),
		st2,
	)
	p3 := discovery.NewDiscoveredPeer(
		discovery.NewPeerIdentity(net.ParseIP("10.10.0.3"), 3000, "key3"),
		discovery.NewPeerHeartbeatState(time.Now(), 0),
		st3,
	)
	repo.SetKnownPeer(p1)
	repo.SetKnownPeer(p2)
	repo.SetKnownPeer(p3)
	unp1, _ := repo.GetKnownPeerByIP(net.ParseIP("10.10.0.1"))
	repo.SetUnreachablePeer(unp1.Identity().PublicKey())
	rp, _ := repo.ListReachablePeers()
	assert.Equal(t, 2, len(rp))
}

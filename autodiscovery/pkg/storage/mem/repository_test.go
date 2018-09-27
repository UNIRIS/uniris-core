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
	st1 := discovery.NewState("0.0", discovery.OkStatus, discovery.PeerPosition{}, "0.0.0", 0.0, 0, 5)
	st2 := discovery.NewState("0.0", discovery.OkStatus, discovery.PeerPosition{}, "0.0.0", 0.0, 0, 5)
	st3 := discovery.NewState("0.0", discovery.OkStatus, discovery.PeerPosition{}, "0.0.0", 0.0, 0, 5)
	p1 := discovery.NewPeerDetailed([]byte("key1"), net.ParseIP("10.10.0.1"), 3000, time.Now(), st1)
	p2 := discovery.NewPeerDetailed([]byte("key2"), net.ParseIP("10.10.0.2"), 3000, time.Now(), st2)
	p3 := discovery.NewPeerDetailed([]byte("key3"), net.ParseIP("10.10.0.3"), 3000, time.Now(), st3)
	repo.AddPeer(p1)
	repo.AddPeer(p2)
	repo.AddPeer(p3)
	unp1, _ := repo.GetPeerByIP(net.ParseIP("10.10.0.1"))
	repo.AddUnreacheablePeer(unp1)
	unp2, _ := repo.GetPeerByIP(net.ParseIP("10.10.0.2"))
	repo.AddUnreacheablePeer(unp2)
	l, _ := repo.ListUnrecheablePeers()
	assert.Equal(t, 2, len(l))
}

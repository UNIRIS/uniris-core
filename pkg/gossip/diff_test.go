package gossip

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	uniris "github.com/uniris/uniris-core/pkg"
)

/*
Scenario: Compare peers with different key and get the unknown
	Given a known peer and a different peer
	When I want to get the unknown peer
	Then I get the second peer
*/
func TestGetUnknownPeersWithDifferentKey(t *testing.T) {
	kp := uniris.NewPeerDigest(
		uniris.NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "key1"),
		uniris.NewPeerHeartbeatState(time.Now(), 1000),
	)

	comparee := uniris.NewPeerDigest(
		uniris.NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "key2"),
		uniris.NewPeerHeartbeatState(time.Now(), 1000),
	)

	unkwown := getUnknownPeers([]uniris.Peer{kp}, []uniris.Peer{comparee})
	assert.Len(t, unkwown, 1)
	assert.Equal(t, "key2", unkwown[0].Identity().PublicKey())
}

/*
Scenario: Compare 2 equal peers and get no one
	Given a known peer and another peer equal
	When I want to get the unknown peer
	Then I get no unknwown peers
*/
func TestGetUnknownPeersWithSameGenerationTime(t *testing.T) {
	kp := uniris.NewPeerDigest(
		uniris.NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "key1"),
		uniris.NewPeerHeartbeatState(time.Now(), 1000),
	)

	comparee := uniris.NewPeerDigest(
		uniris.NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "key1"),
		uniris.NewPeerHeartbeatState(time.Now(), 1000),
	)

	unkwown := getUnknownPeers([]uniris.Peer{kp}, []uniris.Peer{comparee})
	assert.Empty(t, unkwown, 1)
}

/*
Scenario: Compare 2 set of peers with different time and get the recent one
	Given known peers and received peers with different elapsed heartbeats
	When I want to get the unknown peers
	Then I get the peer with the highest elapsed heartbeats
*/
func TestGetUnknownPeersMoreRecent(t *testing.T) {
	kp1 := uniris.NewPeerDigest(
		uniris.NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "key1"),
		uniris.NewPeerHeartbeatState(time.Now(), 1000),
	)
	kp2 := uniris.NewPeerDigest(
		uniris.NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "key2"),
		uniris.NewPeerHeartbeatState(time.Now(), 1000),
	)

	comparee1 := uniris.NewPeerDigest(
		uniris.NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "key2"),
		uniris.NewPeerHeartbeatState(time.Now(), 1000),
	)
	comparee2 := uniris.NewPeerDigest(
		uniris.NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "key1"),
		uniris.NewPeerHeartbeatState(time.Now(), 1200),
	)

	unkwown := getUnknownPeers([]uniris.Peer{kp1, kp2}, []uniris.Peer{comparee1, comparee2})
	assert.Len(t, unkwown, 1)
	assert.Equal(t, "key1", unkwown[0].Identity().PublicKey())
	assert.Equal(t, int64(1200), unkwown[0].HeartbeatState().ElapsedHeartbeats())
}

/*
Scenario: Compare peers with different key and get the new one
	Given a known peer and a received peer
	When I want to get the new peer
	Then I get the first peer
*/
func TestGetNewPeersWithDifferentKey(t *testing.T) {
	kp := uniris.NewPeerDigest(
		uniris.NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "key1"),
		uniris.NewPeerHeartbeatState(time.Now(), 1000),
	)

	comparee := uniris.NewPeerDigest(
		uniris.NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "key2"),
		uniris.NewPeerHeartbeatState(time.Now(), 1000),
	)

	news := getNewPeers([]uniris.Peer{kp}, []uniris.Peer{comparee})
	assert.Len(t, news, 1)
	assert.Equal(t, "key1", news[0].Identity().PublicKey())

}

/*
Scenario: Compare 2 set of peers with different time and get the recent one
	Given known peers and received peers with different elapsed heartbeats
	When I want to get the news peer
	Then I get the peer with the highest elapsed heartbeats
*/
func TestGetNewsPeersMoreRecent(t *testing.T) {
	kp1 := uniris.NewPeerDigest(
		uniris.NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "key1"),
		uniris.NewPeerHeartbeatState(time.Now(), 1000),
	)
	kp2 := uniris.NewPeerDigest(
		uniris.NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "key2"),
		uniris.NewPeerHeartbeatState(time.Now(), 1200),
	)

	comparee1 := uniris.NewPeerDigest(
		uniris.NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "key2"),
		uniris.NewPeerHeartbeatState(time.Now(), 1000),
	)
	comparee2 := uniris.NewPeerDigest(
		uniris.NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "key1"),
		uniris.NewPeerHeartbeatState(time.Now(), 1000),
	)

	news := getNewPeers([]uniris.Peer{kp1, kp2}, []uniris.Peer{comparee1, comparee2})
	assert.Len(t, news, 1)
	assert.Equal(t, "key2", news[0].Identity().PublicKey())
	assert.Equal(t, int64(1200), news[0].HeartbeatState().ElapsedHeartbeats())
}

/*
Scenario: Compare 2 equal peers and get no one
	Given a known peer and another peer equal
	When I want to get the unknown peer
	Then I get no unknwown peers
*/
func TestGetNewPeersWithSameGenerationTime(t *testing.T) {
	kp := uniris.NewPeerDigest(
		uniris.NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "key1"),
		uniris.NewPeerHeartbeatState(time.Now(), 1000),
	)

	comparee := uniris.NewPeerDigest(
		uniris.NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "key1"),
		uniris.NewPeerHeartbeatState(time.Now(), 1000),
	)

	news := getNewPeers([]uniris.Peer{kp}, []uniris.Peer{comparee})
	assert.Empty(t, news, 1)
}

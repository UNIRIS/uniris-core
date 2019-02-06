package discovery

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

/*
Scenario: Compare peers with different key and get the unknown
	Given a known peer and a different peer
	When I want to get the unknown peer
	Then I get the second peer
*/
func TestExcludedOrRecentWithDifferentKey(t *testing.T) {
	kp := NewPeerDigest(
		NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "key1"),
		NewPeerHeartbeatState(time.Now(), 1000),
	)

	comparee := NewPeerDigest(
		NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "key2"),
		NewPeerHeartbeatState(time.Now(), 1000),
	)

	unknown := ExcludedOrRecent([]Peer{kp}, []Peer{comparee})
	assert.Len(t, unknown, 1)
	assert.Equal(t, "key2", unknown[0].Identity().PublicKey())
}

/*
Scenario: Compare 2 equal peers and get no one
	Given a known peer and another peer equal
	When I want to get the unknown peer
	Then I get no unknwown peers
*/
func TestExcludedOrRecentWithSameGenerationTime(t *testing.T) {
	kp := NewPeerDigest(
		NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "key1"),
		NewPeerHeartbeatState(time.Now(), 1000),
	)

	comparee := NewPeerDigest(
		NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "key1"),
		NewPeerHeartbeatState(time.Now(), 1000),
	)

	unknown := ExcludedOrRecent([]Peer{kp}, []Peer{comparee})
	assert.Empty(t, unknown, 1)
}

/*
Scenario: Compare 2 set of peers with different time and get the recent one
	Given known peers and received peers with different elapsed heartbeats
	When I want to get the unknown peers
	Then I get the peer with the highest elapsed heartbeats
*/
func TestExcludedOrRecentMoreRecent(t *testing.T) {
	kp1 := NewPeerDigest(
		NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "key1"),
		NewPeerHeartbeatState(time.Now(), 1000),
	)
	kp2 := NewPeerDigest(
		NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "key2"),
		NewPeerHeartbeatState(time.Now(), 1000),
	)

	comparee1 := NewPeerDigest(
		NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "key2"),
		NewPeerHeartbeatState(time.Now(), 1000),
	)
	comparee2 := NewPeerDigest(
		NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "key1"),
		NewPeerHeartbeatState(time.Now(), 1200),
	)

	unknown := ExcludedOrRecent([]Peer{kp1, kp2}, []Peer{comparee1, comparee2})
	assert.Len(t, unknown, 1)
	assert.Equal(t, "key1", unknown[0].Identity().PublicKey())
	assert.Equal(t, int64(1200), unknown[0].HeartbeatState().ElapsedHeartbeats())
}

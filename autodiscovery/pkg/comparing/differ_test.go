package comparing

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	discovery "github.com/uniris/uniris-core/autodiscovery/pkg"
)

/*
Scenario: Compare with the same peers
	Given a repo with one peer and the same peer to compare
	When the both list contains the same element with the same generation time
	Then we compare them, we got not unknown peer and no new peer to provide
*/
func TestDiffSameGenTime(t *testing.T) {

	kp1 := discovery.NewPeerDigest(
		discovery.NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, []byte("key1")),
		discovery.NewPeerHeartbeatState(time.Now(), 1000),
	)

	comparee := discovery.NewPeerDigest(
		discovery.NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, []byte("key1")),
		discovery.NewPeerHeartbeatState(time.Now(), 1000),
	)

	diff := NewPeerDiffer([]discovery.Peer{kp1})
	up := diff.UnknownPeers([]discovery.Peer{comparee})
	pp := diff.ProvidePeers([]discovery.Peer{comparee})
	assert.Empty(t, up)
	assert.Empty(t, pp)
}

/*
Scenario: Compare different peers
	Given a repo with one peer and a new peer to compare
	When we compare them
	Then we get the unknown peers and new peers to provide
*/
func TestDiffFull(t *testing.T) {

	kp1 := discovery.NewPeerDigest(
		discovery.NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, []byte("key1")),
		discovery.NewPeerHeartbeatState(time.Now(), 1000),
	)

	comparee := discovery.NewPeerDigest(
		discovery.NewPeerIdentity(net.ParseIP("127.0.0.1"), 4000, []byte("key2")),
		discovery.NewPeerHeartbeatState(time.Now(), 1000),
	)

	diff := NewPeerDiffer([]discovery.Peer{kp1})
	up := diff.UnknownPeers([]discovery.Peer{comparee})
	pp := diff.ProvidePeers([]discovery.Peer{comparee})
	assert.NotEmpty(t, up)
	assert.Equal(t, 1, len(up))
	assert.Equal(t, "key2", up[0].Identity().PublicKey().String())
	assert.NotEmpty(t, pp)
	assert.Equal(t, 1, len(pp))
	assert.Equal(t, "key1", pp[0].Identity().PublicKey().String())

}

/*
Scenario: Compare peers and gets the more recents
	Given a repo with two peer and a new list of peer to compare but with a more recent peer
	When we compare them
	Then we get the unknown peer more recent
*/
func TestDiffUnknownRecentPeers(t *testing.T) {

	kp1 := discovery.NewPeerDigest(
		discovery.NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, []byte("key1")),
		discovery.NewPeerHeartbeatState(time.Now(), 1000),
	)
	kp2 := discovery.NewPeerDigest(
		discovery.NewPeerIdentity(net.ParseIP("127.0.0.1"), 4000, []byte("key2")),
		discovery.NewPeerHeartbeatState(time.Now(), 1000),
	)

	c1 := discovery.NewPeerDigest(
		discovery.NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, []byte("key1")),
		discovery.NewPeerHeartbeatState(time.Now(), 1200),
	)
	c2 := discovery.NewPeerDigest(
		discovery.NewPeerIdentity(net.ParseIP("127.0.0.1"), 4000, []byte("key2")),
		discovery.NewPeerHeartbeatState(time.Now(), 1000),
	)

	diff := NewPeerDiffer([]discovery.Peer{kp1, kp2})
	up := diff.UnknownPeers([]discovery.Peer{c1, c2})
	pp := diff.ProvidePeers([]discovery.Peer{c1, c2})

	assert.Empty(t, pp)
	assert.NotEmpty(t, up)

	assert.Equal(t, 1, len(up))
	assert.Equal(t, int64(1200), up[0].HeartbeatState().ElapsedHeartbeats())
}

/*
Scenario: Compare peers and gets the more recents
	Given a repo with two peer and a new list of peer to compare but with a more older peer
	When we compare them
	Then we get the new peer more recent
*/
func TestDiffNewRecentPeers(t *testing.T) {

	kp1 := discovery.NewPeerDigest(
		discovery.NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, []byte("key1")),
		discovery.NewPeerHeartbeatState(time.Now(), 1000),
	)
	kp2 := discovery.NewPeerDigest(
		discovery.NewPeerIdentity(net.ParseIP("127.0.0.1"), 4000, []byte("key2")),
		discovery.NewPeerHeartbeatState(time.Now(), 1200),
	)

	c1 := discovery.NewPeerDigest(
		discovery.NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, []byte("key1")),
		discovery.NewPeerHeartbeatState(time.Now(), 1000),
	)
	c2 := discovery.NewPeerDigest(
		discovery.NewPeerIdentity(net.ParseIP("127.0.0.1"), 4000, []byte("key2")),
		discovery.NewPeerHeartbeatState(time.Now(), 1000),
	)

	diff := NewPeerDiffer([]discovery.Peer{kp1, kp2})
	up := diff.UnknownPeers([]discovery.Peer{c1, c2})
	pp := diff.ProvidePeers([]discovery.Peer{c1, c2})

	assert.Empty(t, up)
	assert.NotEmpty(t, pp)

	assert.Equal(t, 1, len(pp))
	assert.Equal(t, int64(1200), pp[0].HeartbeatState().ElapsedHeartbeats())
}

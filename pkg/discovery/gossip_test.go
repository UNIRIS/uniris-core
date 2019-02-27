package discovery

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

/*
Scenario: Store and notifies cycle discovered peers
	Given a gossip cycle with a discovered peer
	When I want to add them
	Then the store will included them and will be sent to the notifier
*/
func TestAddDiscoveries(t *testing.T) {

	store := &mockDatabase{}
	notif := &mockNotifier{}

	c := Cycle{
		Discoveries: []Peer{
			NewDiscoveredPeer(
				NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "key"),
				NewPeerHeartbeatState(time.Now(), 0),
				NewPeerAppState("1.0", OkPeerStatus, 30.0, 10.0, "", 200, 1, 0),
			),
		},
	}
	assert.Nil(t, addDiscoveries(c, store, notif))

	assert.Len(t, store.discoveredPeers, 1)
	assert.Equal(t, "key", store.discoveredPeers[0].identity.publicKey)
	assert.Equal(t, "key", notif.discoveries[0].identity.publicKey)
}

/*
Scenario: Store and notifies cycle unreachables peers
	Given a gossip cycle with a unreachable peer
	When I want to add them
	Then the store will included them and will be sent to the notifier
*/
func TestAddUnreachable(t *testing.T) {

	db := &mockDatabase{}
	notif := &mockNotifier{}

	c := Cycle{
		Unreaches: []PeerIdentity{
			NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "key"),
		},
	}
	assert.Nil(t, addUnreaches(c, db, notif))

	assert.Len(t, db.unreachablePeers, 1)
	assert.Equal(t, "key", db.unreachablePeers[0].publicKey)
	assert.Equal(t, "key", notif.unreaches[0].publicKey)

}

/*
Scenario: Store and notifies cycle reachable peers
	Given a gossip cycle with a reachable peer
	When I want to add them
	Then the database will included them and will be sent to the notifier
*/
func TestAddReachable(t *testing.T) {

	db := &mockDatabase{}
	notif := &mockNotifier{}

	err := db.WriteUnreachablePeer(NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "key"))
	assert.Nil(t, err)
	assert.Len(t, db.unreachablePeers, 1)

	c := Cycle{
		Reaches: []PeerIdentity{
			NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "key"),
		},
	}

	assert.Nil(t, addReaches(c, db, notif))
	assert.Equal(t, "key", notif.reaches[0].publicKey)
}

/*
Scenario: Store and notifies cycle reachables peers after unreach
	Given a database including a unreachble peer
	When the new cycle get as reachable this peer
	Then the peer is removed from the unreachable store and stored as reachable
*/
func TestAddReachbleAfterBeingUnreachale(t *testing.T) {
	db := &mockDatabase{}
	notif := &mockNotifier{}

	err := db.WriteUnreachablePeer(NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "key"))
	assert.Nil(t, err)
	assert.Len(t, db.unreachablePeers, 1)

	c := Cycle{
		Reaches: []PeerIdentity{
			NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "key"),
		},
	}

	assert.Nil(t, addReaches(c, db, notif))
	assert.Len(t, db.unreachablePeers, 0)
	assert.Equal(t, "key", notif.reaches[0].PublicKey())
}

/*
Scenario: Process a gossip cycle
	Given a list of seeds and a local peer
	When I want to gossip and run a cycle
	THen I get some peers stored and notified
*/
func TestGossip(t *testing.T) {
	db := &mockDatabase{}
	notif := &mockNotifier{}

	self := NewSelfPeer("key", net.ParseIP("127.0.0.1"), 3000, "", 30.0, 10.0)
	seeds := []PeerIdentity{
		NewPeerIdentity(net.ParseIP("10.0.0.1"), 3000, "key1"),
	}

	c, err := Gossip(self, seeds, db, mockNetworkChecker{}, mockSystemReader{}, mockClient{}, notif)
	assert.Nil(t, err)
	assert.Len(t, c.Discoveries, 1)
	assert.Len(t, c.Reaches, 1)
	assert.Len(t, c.Unreaches, 0)
	assert.Len(t, db.discoveredPeers, 1)

	assert.Equal(t, "dKey1", db.discoveredPeers[0].Identity().PublicKey())
}

type mockDatabase struct {
	unreachablePeers []PeerIdentity
	discoveredPeers  []Peer
}

func (db *mockDatabase) DiscoveredPeers() ([]Peer, error) {
	return db.discoveredPeers, nil
}

func (db *mockDatabase) WriteDiscoveredPeer(peer Peer) error {
	if db.containsPeer(peer) {
		for _, p := range db.discoveredPeers {
			if p.Identity().PublicKey() == peer.Identity().PublicKey() {
				p = peer
				break
			}
		}
	} else {
		db.discoveredPeers = append(db.discoveredPeers, peer)
	}
	return nil
}

func (db *mockDatabase) UnreachablePeers() ([]PeerIdentity, error) {
	pp := make([]PeerIdentity, 0)
	for i := 0; i < len(db.discoveredPeers); i++ {
		if exist, _ := db.ContainsUnreachablePeer(db.discoveredPeers[i].Identity()); exist {
			pp = append(pp, db.discoveredPeers[i].Identity())
		}
	}
	return pp, nil
}

func (db *mockDatabase) WriteUnreachablePeer(pk PeerIdentity) error {
	if exist, _ := db.ContainsUnreachablePeer(pk); !exist {
		db.unreachablePeers = append(db.unreachablePeers, pk)
	}
	return nil
}

func (db *mockDatabase) RemoveUnreachablePeer(pk PeerIdentity) error {
	if exist, _ := db.ContainsUnreachablePeer(pk); exist {
		for i := 0; i < len(db.unreachablePeers); i++ {
			if db.unreachablePeers[i].publicKey == pk.publicKey {
				db.unreachablePeers = db.unreachablePeers[:i+copy(db.unreachablePeers[i:], db.unreachablePeers[i+1:])]
			}
		}
	}
	return nil
}

func (db *mockDatabase) ContainsUnreachablePeer(pk PeerIdentity) (bool, error) {
	for _, up := range db.unreachablePeers {
		if up.publicKey == pk.PublicKey() {
			return true, nil
		}
	}
	return false, nil
}

func (db *mockDatabase) containsPeer(p Peer) bool {
	mdiscoveredPeers := make(map[string]Peer, 0)
	for _, p := range db.discoveredPeers {
		mdiscoveredPeers[p.Identity().PublicKey()] = p
	}

	_, exist := mdiscoveredPeers[p.Identity().PublicKey()]
	return exist
}

type mockNotifier struct {
	reaches     []PeerIdentity
	unreaches   []PeerIdentity
	discoveries []Peer
}

func (n *mockNotifier) NotifyReachable(p PeerIdentity) error {
	n.reaches = append(n.reaches, p)
	return nil
}
func (n *mockNotifier) NotifyUnreachable(p PeerIdentity) error {
	n.unreaches = append(n.unreaches, p)
	return nil
}

func (n *mockNotifier) NotifyDiscovery(p Peer) error {
	n.discoveries = append(n.discoveries, p)
	return nil
}

type mockNetworkChecker struct{}

func (nc mockNetworkChecker) CheckNtpState() error {
	return nil
}

func (nc mockNetworkChecker) CheckInternetState() error {
	return nil
}

func (nc mockNetworkChecker) CheckGRPCServer() error {
	return nil
}

type mockSystemReader struct{}

func (i mockSystemReader) GeoPosition() (lon float64, lat float64, err error) {
	return 10.0, 30.0, nil
}

func (i mockSystemReader) FreeDiskSpace() (float64, error) {
	return 200, nil
}

func (i mockSystemReader) CPULoad() (string, error) {
	return "", nil
}

func (i mockSystemReader) IP() (net.IP, error) {
	return net.ParseIP("127.0.0.1"), nil
}

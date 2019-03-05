package discovery

import (
	"errors"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

/*
Scenario: Spread a gossip round and discover peers
	Given a initiator peer, a receiver peer and list of known peers
	When we start a gossip round we spread what we know
	Then we get the new peers discovered
*/
func TestStartRound(t *testing.T) {

	target := NewPeerIdentity(net.ParseIP("20.100.4.120"), 3000, "key2")

	p1 := NewDiscoveredPeer(
		NewPeerIdentity(net.ParseIP("50.20.100.2"), 3000, "key3"),
		NewPeerHeartbeatState(time.Now(), 0),
		NewPeerAppState("1.0", OkPeerStatus, 10.0, 20.0, "", 0, 1, 0),
	)

	p2 := NewDiscoveredPeer(
		NewPeerIdentity(net.ParseIP("50.10.30.2"), 3000, "uKey1"),
		NewPeerHeartbeatState(time.Now(), 0),
		NewPeerAppState("1.0", OkPeerStatus, 20.0, 19.4, "", 0, 1, 0),
	)

	discoveries, err := startRound(target, []Peer{p1, p2}, mockMessenger{})
	assert.Nil(t, err)
	assert.Len(t, discoveries, 1)
	assert.Equal(t, "dKey1", discoveries[0].Identity().PublicKey())
}

/*
Scenario: Spread gossip but unreach the target peer during the SYN request
	Given a initiator peer, a receiver peer and list of known peers
	When are sending the SYN, the target cannot be reached
	Then we get the unreached peer
*/
func TestStartRoundWithUnreachWhenSYN(t *testing.T) {

	target := NewPeerIdentity(net.ParseIP("20.100.4.120"), 3000, "key2")

	p1 := NewDiscoveredPeer(
		NewPeerIdentity(net.ParseIP("50.20.100.2"), 3000, "key3"),
		NewPeerHeartbeatState(time.Now(), 0),
		NewPeerAppState("1.0", OkPeerStatus, 30.0, 10.0, "", 0, 1, 0),
	)

	p2 := NewDiscoveredPeer(
		NewPeerIdentity(net.ParseIP("50.10.30.2"), 3000, "uKey1"),
		NewPeerHeartbeatState(time.Now(), 0),
		NewPeerAppState("1.0", OkPeerStatus, 30.0, 10.0, "", 0, 1, 0),
	)

	_, err := startRound(target, []Peer{p1, p2}, mockMessengerWithSynFailure{})
	assert.Equal(t, err, ErrUnreachablePeer)
}

/*
Scenario: Spread gossip but unreach the target peer during the SYN request
	Given a initiator peer, a receiver peer and list of known peers
	When are sending the SYN, the target cannot be reached
	Then we get the unreached peer
*/
func TestStartRoundWithUnreachWhenACK(t *testing.T) {

	target := NewPeerIdentity(net.ParseIP("20.100.4.120"), 3000, "key2")

	p1 := NewDiscoveredPeer(
		NewPeerIdentity(net.ParseIP("50.20.100.2"), 3000, "key3"),
		NewPeerHeartbeatState(time.Now(), 0),
		NewPeerAppState("1.0", OkPeerStatus, 30.0, 10.0, "", 0, 1, 0),
	)

	p2 := NewDiscoveredPeer(
		NewPeerIdentity(net.ParseIP("50.10.30.2"), 3000, "uKey1"),
		NewPeerHeartbeatState(time.Now(), 0),
		NewPeerAppState("1.0", OkPeerStatus, 30.0, 10.0, "", 0, 1, 0),
	)

	_, err := startRound(target, []Peer{p1, p2}, mockMessengerWithAckFailure{})
	assert.Equal(t, err, ErrUnreachablePeer)
}

/*
Scenario: Get a random peers from a list
	Given a list of peers
	When I want to get a random peer
	Then I get one
*/
func TestRandomPeers(t *testing.T) {
	p1 := PeerIdentity{ip: net.ParseIP("127.0.0.1"), port: 3000, publicKey: "key1"}
	p2 := PeerIdentity{ip: net.ParseIP("127.0.0.2"), port: 4000, publicKey: "key2"}
	p3 := PeerIdentity{ip: net.ParseIP("127.0.0.3"), port: 5000, publicKey: "key3"}

	r1 := randomPeer([]PeerIdentity{p1, p2, p3})
	r2 := randomPeer([]PeerIdentity{p1, p2, p3})

	assert.NotEqual(t, r1, r2)

	r1 = randomPeer([]PeerIdentity{p1})
	assert.Equal(t, p1, r1)
}

/*
Scenario: Run a gossip cycle
	Given a selfator and a target
	When we create a round associated to a cycle
	Then we run it and get some discovered peers
*/
func TestRunCycle(t *testing.T) {
	self := NewSelfPeer("key", net.ParseIP("10.0.0.1"), 3000, "1.0", 30.0, 12.0)

	kp1 := NewPeerDigest(
		NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "key2"),
		NewPeerHeartbeatState(time.Now(), 0))

	seeds := []PeerIdentity{
		NewPeerIdentity(net.ParseIP("20.0.0.1"), 3000, "key3"),
	}

	discoveries, reachables, unreachables, err := startCycle(self, mockMessenger{}, seeds, []Peer{kp1}, []PeerIdentity{kp1.Identity()}, []PeerIdentity{})

	assert.Nil(t, err)
	assert.Len(t, discoveries, 2)
	assert.Len(t, reachables, 2)
	assert.Empty(t, unreachables)

	//Peer retrieved from the kp1
	assert.Equal(t, "dKey1", discoveries[0].Identity().PublicKey())

	//Peer retreived from the seed1
	assert.Equal(t, "dKey1", discoveries[1].Identity().PublicKey())
}

/*
Scenario: Run a gossip cycle with gossip failure
	Given a some peers and a target
	When we run cycle and cannot reach the target
	Then we run it and get some discovered peers
*/
func TestRunCycleGetUnreachable(t *testing.T) {
	self := NewSelfPeer("key", net.ParseIP("10.0.0.1"), 3000, "1.0", 30.0, 12.0)

	kp1 := NewPeerDigest(
		NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "key2"),
		NewPeerHeartbeatState(time.Now(), 0))

	seeds := []PeerIdentity{
		NewPeerIdentity(net.ParseIP("20.0.0.1"), 3000, "key3"),
	}

	_, _, unreachables, err := startCycle(self, mockMessengerWithSynFailure{}, seeds, []Peer{kp1}, []PeerIdentity{kp1.Identity()}, []PeerIdentity{})

	assert.Nil(t, err)
	assert.Len(t, unreachables, 2)
}

/*
Scenario: Run a gossip cycle with unreachable peers
	Given a some peers reachables and unreachables and a target
	When we run cycle
	Then we get an unreachable peer
*/
func TestRunCycleWithUnreachable(t *testing.T) {
	self := NewSelfPeer("key", net.ParseIP("10.0.0.1"), 3000, "1.0", 30.0, 12.0)

	kp1 := NewPeerDigest(
		NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "key2"),
		NewPeerHeartbeatState(time.Now(), 0))

	ur1 := NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "key4")

	seeds := []PeerIdentity{
		NewPeerIdentity(net.ParseIP("20.0.0.1"), 3000, "key3"),
	}

	_, reachables, unreachables, err := startCycle(self, mockMessenger{}, seeds, []Peer{kp1}, []PeerIdentity{kp1.Identity()}, []PeerIdentity{ur1})

	assert.Nil(t, err)
	assert.Len(t, unreachables, 0)
	assert.Len(t, reachables, 3)
}

/*
Scenario: Store and notifies cycle discovered peers
	Given a gossip cycle with a discovered peer
	When I want to add them
	Then the store will included them and will be sent to the notifier
*/
func TestAddDiscoveries(t *testing.T) {

	store := &mockDatabase{}
	notif := &mockNotifier{}

	err := addDiscoveries([]Peer{
		NewDiscoveredPeer(
			NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "key"),
			NewPeerHeartbeatState(time.Now(), 0),
			NewPeerAppState("1.0", OkPeerStatus, 30.0, 10.0, "", 200, 1, 0),
		),
	}, store, notif)
	assert.Nil(t, err)

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

	err := addUnreaches([]PeerIdentity{
		NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "key"),
	}, nil, db, notif)
	assert.Nil(t, err)

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

	err := addReaches([]PeerIdentity{
		NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "key"),
	}, []PeerIdentity{
		NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "key"),
	}, db, notif)
	assert.Nil(t, err)
	assert.Equal(t, "key", notif.reaches[0].publicKey)
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

	err := Gossip(self, seeds, db, mockNetworkChecker{}, mockSystemReader{}, mockMessenger{}, notif)
	assert.Nil(t, err)
	assert.Len(t, db.discoveredPeers, 2)
	assert.Len(t, db.unreachablePeers, 0)
	assert.Equal(t, "key", db.discoveredPeers[0].Identity().PublicKey())
	assert.Equal(t, "dKey1", db.discoveredPeers[1].Identity().PublicKey())
	assert.Len(t, notif.discoveries, 1)
	assert.Len(t, notif.reaches, 0)
	assert.Len(t, notif.unreaches, 0)
}

/*
Scenario: Compare peers with different key and get the unknown
	Given a known peer and a different peer
	When I want to get the unknown peer
	Then I get the second peer
*/
func TestCompareWithDifferentKey(t *testing.T) {
	kp := NewPeerDigest(
		NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "key1"),
		NewPeerHeartbeatState(time.Now(), 1000),
	)

	comparee := NewPeerDigest(
		NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "key2"),
		NewPeerHeartbeatState(time.Now(), 1000),
	)

	unknown := ComparePeers([]Peer{kp}, []Peer{comparee})
	assert.Len(t, unknown, 1)
	assert.Equal(t, "key2", unknown[0].Identity().PublicKey())
}

/*
Scenario: Compare 2 equal peers and get no one
	Given a known peer and another peer equal
	When I want to get the unknown peer
	Then I get no unknwown peers
*/
func TestCompareWithSameGenerationTime(t *testing.T) {
	kp := NewPeerDigest(
		NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "key1"),
		NewPeerHeartbeatState(time.Now(), 1000),
	)

	comparee := NewPeerDigest(
		NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "key1"),
		NewPeerHeartbeatState(time.Now(), 1000),
	)

	unknown := ComparePeers([]Peer{kp}, []Peer{comparee})
	assert.Empty(t, unknown)
}

/*
Scenario: Compare 2 set of peers with different time and get the recent one
	Given known peers and received peers with different elapsed heartbeats
	When I want to get the unknown peers
	Then I get the peer with the highest elapsed heartbeats
*/
func TestCompareMoreRecent(t *testing.T) {
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

	unknown := ComparePeers([]Peer{kp1, kp2}, []Peer{comparee1, comparee2})
	assert.Len(t, unknown, 1)
	assert.Equal(t, "key1", unknown[0].Identity().PublicKey())
	assert.Equal(t, int64(1200), unknown[0].HeartbeatState().ElapsedHeartbeats())
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

type mockMessengerWithSynFailure struct {
}

func (m mockMessengerWithSynFailure) SendSyn(target PeerIdentity, known []Peer) (reqPeers []PeerIdentity, discoveries []Peer, err error) {
	return nil, nil, ErrUnreachablePeer
}

func (m mockMessengerWithSynFailure) SendAck(target PeerIdentity, requested []Peer) error {
	return nil
}

type mockMessengerWithAckFailure struct {
}

func (m mockMessengerWithAckFailure) SendSyn(target PeerIdentity, known []Peer) (requested []PeerIdentity, discoveries []Peer, err error) {
	reqP := NewPeerIdentity(net.ParseIP("200.18.186.39"), 3000, "uKey1")

	hb := NewPeerHeartbeatState(time.Now(), 0)
	as := NewPeerAppState("1.0", OkPeerStatus, 30.0, 10.0, "", 0, 1, 0)

	np1 := NewDiscoveredPeer(
		NewPeerIdentity(net.ParseIP("35.200.100.2"), 3000, "dKey1"),
		hb, as,
	)

	return []PeerIdentity{reqP}, []Peer{np1}, nil
}

func (m mockMessengerWithAckFailure) SendAck(target PeerIdentity, requested []Peer) error {
	return ErrUnreachablePeer
}

type mockMessengerUnexpectedFailure struct {
}

func (m mockMessengerUnexpectedFailure) SendSyn(target PeerIdentity, known []Peer) (unknown []Peer, new []Peer, err error) {
	return nil, nil, errors.New("Unexpected")
}

func (m mockMessengerUnexpectedFailure) SendAck(target PeerIdentity, requested []Peer) error {
	return nil
}

type mockMessenger struct{}

func (m mockMessenger) SendSyn(target PeerIdentity, known []Peer) (requested []PeerIdentity, discoveries []Peer, err error) {
	reqP := NewPeerIdentity(net.ParseIP("200.18.186.39"), 3000, "uKey1")

	np1 := NewDiscoveredPeer(
		NewPeerIdentity(net.ParseIP("35.200.100.2"), 3000, "dKey1"),
		NewPeerHeartbeatState(time.Now(), 0),
		NewPeerAppState("1.0", OkPeerStatus, 50.1, 22.1, "", 0, 1, 0),
	)

	return []PeerIdentity{reqP}, []Peer{np1}, nil
}

func (m mockMessenger) SendAck(target PeerIdentity, requested []Peer) error {
	return nil
}

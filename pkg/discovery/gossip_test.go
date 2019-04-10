package discovery

import (
	"crypto/rand"
	"errors"
	"github.com/uniris/uniris-core/pkg/crypto"
	"github.com/uniris/uniris-core/pkg/logging"
	"log"
	"net"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	_, pubx, _ = crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
	_, puby, _ = crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
)

/*
Scenario: Spread a gossip round and discover peers
	Given a initiator peer, a receiver peer and list of known peers
	When we start a gossip round we spread what we know
	Then we get the new peers discovered
*/
func TestStartRound(t *testing.T) {
	_, pub1, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
	_, pub2, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
	_, pub3, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

	l := logging.NewLogger("stdout", log.New(os.Stdout, "", 0), "test", net.ParseIP("127.0.0.1"), logging.ErrorLogLevel)
	target := NewPeerIdentity(net.ParseIP("20.100.4.120"), 3000, pub1)

	p1 := NewDiscoveredPeer(
		NewPeerIdentity(net.ParseIP("50.20.100.2"), 3000, pub2),
		NewPeerHeartbeatState(time.Now(), 0),
		NewPeerAppState("1.0", OkPeerStatus, 10.0, 20.0, "", 0, 1, 0),
	)

	p2 := NewDiscoveredPeer(
		NewPeerIdentity(net.ParseIP("50.10.30.2"), 3000, pub3),
		NewPeerHeartbeatState(time.Now(), 0),
		NewPeerAppState("1.0", OkPeerStatus, 20.0, 19.4, "", 0, 1, 0),
	)

	discoveries, err := startRound(target, []Peer{p1, p2}, mockMessenger{l}, l)
	assert.Nil(t, err)
	assert.Len(t, discoveries, 1)
	assert.True(t, discoveries[0].Identity().PublicKey().Equals(puby))
}

/*
Scenario: Spread gossip but unreach the target peer during the SYN request
	Given a initiator peer, a receiver peer and list of known peers
	When are sending the SYN, the target cannot be reached
	Then we get the unreached peer error
*/
func TestStartRoundWithUnreachWhenSYN(t *testing.T) {

	l := logging.NewLogger("stdout", log.New(os.Stdout, "", 0), "test", net.ParseIP("127.0.0.1"), logging.ErrorLogLevel)
	_, pub1, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
	_, pub2, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
	_, pub3, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

	target := NewPeerIdentity(net.ParseIP("20.100.4.120"), 3000, pub1)

	p1 := NewDiscoveredPeer(
		NewPeerIdentity(net.ParseIP("50.20.100.2"), 3000, pub2),
		NewPeerHeartbeatState(time.Now(), 0),
		NewPeerAppState("1.0", OkPeerStatus, 30.0, 10.0, "", 0, 1, 0),
	)

	p2 := NewDiscoveredPeer(
		NewPeerIdentity(net.ParseIP("50.10.30.2"), 3000, pub3),
		NewPeerHeartbeatState(time.Now(), 0),
		NewPeerAppState("1.0", OkPeerStatus, 30.0, 10.0, "", 0, 1, 0),
	)

	_, err := startRound(target, []Peer{p1, p2}, mockMessengerWithSynFailure{l}, l)
	assert.Equal(t, err, ErrUnreachablePeer)
}

/*
Scenario: Spread gossip but unreach the target peer during the SYN request
	Given a initiator peer, a receiver peer and list of known peers
	When are sending the SYN, the target cannot be reached
	Then we get the unreached peer  error
*/
func TestStartRoundWithUnreachWhenACK(t *testing.T) {

	l := logging.NewLogger("stdout", log.New(os.Stdout, "", 0), "test", net.ParseIP("127.0.0.1"), logging.ErrorLogLevel)
	_, pub1, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
	_, pub2, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
	_, pub3, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

	target := NewPeerIdentity(net.ParseIP("20.100.4.120"), 3000, pub1)

	p1 := NewDiscoveredPeer(
		NewPeerIdentity(net.ParseIP("50.20.100.2"), 3000, pub2),
		NewPeerHeartbeatState(time.Now(), 0),
		NewPeerAppState("1.0", OkPeerStatus, 30.0, 10.0, "", 0, 1, 0),
	)

	p2 := NewDiscoveredPeer(
		NewPeerIdentity(net.ParseIP("50.10.30.2"), 3000, pub3),
		NewPeerHeartbeatState(time.Now(), 0),
		NewPeerAppState("1.0", OkPeerStatus, 30.0, 10.0, "", 0, 1, 0),
	)

	_, err := startRound(target, []Peer{p1, p2}, mockMessengerWithAckFailure{l}, l)
	assert.Equal(t, err, ErrUnreachablePeer)
}

/*
Scenario: Get a random peers from a list
	Given a list of peers
	When I want to get a random peer
	Then I get one
*/
func TestRandomPeers(t *testing.T) {

	_, pub1, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
	_, pub2, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
	_, pub3, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

	p1 := PeerIdentity{ip: net.ParseIP("127.0.0.1"), port: 3000, publicKey: pub1}
	p2 := PeerIdentity{ip: net.ParseIP("127.0.0.2"), port: 4000, publicKey: pub2}
	p3 := PeerIdentity{ip: net.ParseIP("127.0.0.3"), port: 5000, publicKey: pub3}

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
	_, pub1, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
	_, pub2, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
	_, pub3, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

	self := NewSelfPeer(pub1, net.ParseIP("10.0.0.1"), 3000, "1.0", 30.0, 12.0)

	kp1 := NewPeerDigest(
		NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, pub2),
		NewPeerHeartbeatState(time.Now(), 0))

	seeds := []PeerIdentity{
		NewPeerIdentity(net.ParseIP("20.0.0.1"), 3000, pub3),
	}

	l := logging.NewLogger("stdout", log.New(os.Stdout, "", 0), "test", net.ParseIP("127.0.0.1"), logging.ErrorLogLevel)

	discoveries, reachables, unreachables, err := startCycle(self, mockMessenger{l}, seeds, []Peer{kp1}, []PeerIdentity{kp1.Identity()}, []PeerIdentity{}, l)

	assert.Nil(t, err)
	assert.Len(t, discoveries, 2)
	assert.Len(t, reachables, 2)
	assert.Empty(t, unreachables)

	//Peer retrieved from the kp1
	assert.Equal(t, puby, discoveries[0].Identity().PublicKey())

	//Peer retreived from the seed1
	assert.Equal(t, puby, discoveries[1].Identity().PublicKey())
}

/*
Scenario: Run a gossip cycle with gossip failure
	Given a some peers and a target
	When we run cycle and cannot reach the target
	Then we run it and get some discovered peers
*/
func TestRunCycleGetUnreachable(t *testing.T) {
	_, pub1, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
	_, pub2, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
	_, pub3, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
	self := NewSelfPeer(pub1, net.ParseIP("10.0.0.1"), 3000, "1.0", 30.0, 12.0)

	kp1 := NewPeerDigest(
		NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, pub2),
		NewPeerHeartbeatState(time.Now(), 0))

	seeds := []PeerIdentity{
		NewPeerIdentity(net.ParseIP("20.0.0.1"), 3000, pub3),
	}

	l := logging.NewLogger("stdout", log.New(os.Stdout, "", 0), "test", net.ParseIP("127.0.0.1"), logging.ErrorLogLevel)

	_, _, unreachables, err := startCycle(self, mockMessengerWithSynFailure{l}, seeds, []Peer{kp1}, []PeerIdentity{kp1.Identity()}, []PeerIdentity{}, l)

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
	_, pub1, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
	_, pub2, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
	_, pub3, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
	_, pub4, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

	self := NewSelfPeer(pub1, net.ParseIP("10.0.0.1"), 3000, "1.0", 30.0, 12.0)

	kp1 := NewPeerDigest(
		NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, pub2),
		NewPeerHeartbeatState(time.Now(), 0))

	ur1 := NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, pub4)

	seeds := []PeerIdentity{
		NewPeerIdentity(net.ParseIP("20.0.0.1"), 3000, pub3),
	}

	l := logging.NewLogger("stdout", log.New(os.Stdout, "", 0), "test", net.ParseIP("127.0.0.1"), logging.ErrorLogLevel)

	_, reachables, unreachables, err := startCycle(self, mockMessenger{l}, seeds, []Peer{kp1}, []PeerIdentity{kp1.Identity()}, []PeerIdentity{ur1}, l)

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
	_, pub1, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

	store := &mockDatabase{}
	notif := &mockNotifier{}

	p1 := NewDiscoveredPeer(
		NewPeerIdentity(net.ParseIP("192.168.1.1"), 3000, pub1),
		NewPeerHeartbeatState(time.Now(), 0),
		NewPeerAppState("1.0", OkPeerStatus, 30.0, 10.0, "", 0, 1, 0),
	)

	err := addDiscoveries([]Peer{
		NewDiscoveredPeer(
			NewPeerIdentity(net.ParseIP("192.168.1.1"), 3000, pub1),
			NewPeerHeartbeatState(time.Now(), 0),
			NewPeerAppState("1.0", OkPeerStatus, 30.0, 10.0, "", 0, 1, 0),
		),
	}, []Peer{p1}, store, notif)
	assert.Nil(t, err)

	assert.Len(t, store.discoveredPeers, 1)
	assert.Equal(t, pub1, store.discoveredPeers[0].identity.publicKey)
	assert.Equal(t, 0, len(notif.discoveries))

	err = addDiscoveries([]Peer{
		NewDiscoveredPeer(
			NewPeerIdentity(net.ParseIP("192.168.1.1"), 3000, pub1),
			NewPeerHeartbeatState(time.Now(), 0),
			NewPeerAppState("1.0", OkPeerStatus, 30.0, 10.0, "", 10, 10, 10),
		),
	}, []Peer{p1}, store, notif)
	assert.Nil(t, err)
	assert.Len(t, store.discoveredPeers, 1)
	assert.Equal(t, pub1, store.discoveredPeers[0].identity.publicKey)
	assert.Equal(t, 1, len(notif.discoveries))
	assert.Equal(t, pub1, notif.discoveries[0].identity.publicKey)
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

	_, pub1, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

	err := addUnreaches([]PeerIdentity{
		NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, pub1),
	}, nil, db, notif)
	assert.Nil(t, err)

	assert.Len(t, db.unreachablePeers, 1)
	assert.Equal(t, pub1, db.unreachablePeers[0].publicKey)

	p, err := pub1.Marshal()
	assert.Nil(t, err)

	assert.Equal(t, string(p), notif.unreaches[0])

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

	_, pub1, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

	err := addReaches([]PeerIdentity{
		NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, pub1),
	}, []PeerIdentity{
		NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, pub1),
	}, db, notif)
	assert.Nil(t, err)

	p, err := pub1.Marshal()
	assert.Nil(t, err)

	assert.Equal(t, string(p), notif.reaches[0])
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

	_, pub1, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
	_, pub2, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

	self := NewSelfPeer(pub1, net.ParseIP("127.0.0.1"), 3000, "", 30.0, 10.0)
	seeds := []PeerIdentity{
		NewPeerIdentity(net.ParseIP("10.0.0.1"), 3000, pub2),
	}

	l := logging.NewLogger("stdout", log.New(os.Stdout, "", 0), "test", net.ParseIP("127.0.0.1"), logging.ErrorLogLevel)
	err := Gossip(self, seeds, db, mockNetworkChecker{}, mockSystemReader{}, mockMessenger{l}, notif, l)
	assert.Nil(t, err)
	assert.Len(t, db.discoveredPeers, 2)
	assert.Len(t, db.unreachablePeers, 0)
	assert.Equal(t, pub1, db.discoveredPeers[0].Identity().PublicKey())
	assert.Equal(t, puby, db.discoveredPeers[1].Identity().PublicKey())
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

	_, pub1, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
	_, pub2, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

	kp := NewPeerDigest(
		NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, pub1),
		NewPeerHeartbeatState(time.Now(), 1000),
	)

	comparee := NewPeerDigest(
		NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, pub2),
		NewPeerHeartbeatState(time.Now(), 1000),
	)

	unknown := ComparePeers([]Peer{kp}, []Peer{comparee})
	assert.Len(t, unknown, 1)
	assert.Equal(t, pub2, unknown[0].Identity().PublicKey())
}

/*
Scenario: Compare 2 equal peers and get no one
	Given a known peer and another peer equal
	When I want to get the unknown peer
	Then I get no unknwown peers
*/
func TestCompareWithSameGenerationTime(t *testing.T) {

	_, pub1, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

	kp := NewPeerDigest(
		NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, pub1),
		NewPeerHeartbeatState(time.Now(), 1000),
	)

	comparee := NewPeerDigest(
		NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, pub1),
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

	_, pub1, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
	_, pub2, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

	kp1 := NewPeerDigest(
		NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, pub1),
		NewPeerHeartbeatState(time.Now(), 1000),
	)
	kp2 := NewPeerDigest(
		NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, pub2),
		NewPeerHeartbeatState(time.Now(), 1000),
	)

	comparee1 := NewPeerDigest(
		NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, pub1),
		NewPeerHeartbeatState(time.Now(), 1000),
	)
	comparee2 := NewPeerDigest(
		NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, pub2),
		NewPeerHeartbeatState(time.Now(), 1200),
	)

	unknown := ComparePeers([]Peer{kp1, kp2}, []Peer{comparee1, comparee2})
	assert.Len(t, unknown, 1)
	assert.Equal(t, true, pub2.Equals(unknown[0].Identity().PublicKey()))
	assert.Equal(t, int64(1200), unknown[0].HeartbeatState().ElapsedHeartbeats())
}

/*
Scenario: Compare 2 version of a peer with a different appstate
	Given a source peer and comparee peers with different appstate
	When I want to compare appstate
	Then I get wanted result
*/

func TestComparePeerIDAndState(t *testing.T) {

	_, pub1, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

	p1 := NewDiscoveredPeer(
		NewPeerIdentity(net.ParseIP("50.20.100.2"), 3000, pub1),
		NewPeerHeartbeatState(time.Now(), 0),
		NewPeerAppState("1.0", OkPeerStatus, 30.0, 10.0, "", 0, 1, 0),
	)

	p2 := NewDiscoveredPeer(
		NewPeerIdentity(net.ParseIP("50.20.100.2"), 3000, pub1),
		NewPeerHeartbeatState(time.Now(), 0),
		NewPeerAppState("1.0", OkPeerStatus, 30.0, 10.0, "", 2, 2, 2),
	)

	p3 := NewDiscoveredPeer(
		NewPeerIdentity(net.ParseIP("50.20.100.2"), 3000, pub1),
		NewPeerHeartbeatState(time.Now(), 0),
		NewPeerAppState("1.0", OkPeerStatus, 30.0, 10.0, "", 0, 1, 0),
	)

	assert.False(t, comparePeerIDAndState(p1, p2))
	assert.True(t, comparePeerIDAndState(p1, p3))

	p4 := NewDiscoveredPeer(
		NewPeerIdentity(net.ParseIP("50.20.100.1"), 3000, pub1),
		NewPeerHeartbeatState(time.Now(), 0),
		NewPeerAppState("1.0", OkPeerStatus, 30.0, 10.0, "", 0, 1, 0),
	)

	p5 := NewDiscoveredPeer(
		NewPeerIdentity(net.ParseIP("50.20.100.4"), 3000, pub1),
		NewPeerHeartbeatState(time.Now(), 0),
		NewPeerAppState("1.0", OkPeerStatus, 30.0, 10.0, "", 0, 1, 0),
	)

	assert.False(t, comparePeerIDAndState(p4, p5))
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
			if p.Identity().PublicKey().Equals(peer.Identity().PublicKey()) {
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

func (db *mockDatabase) WriteUnreachablePeer(pi PeerIdentity) error {
	if exist, _ := db.ContainsUnreachablePeer(pi); !exist {
		db.unreachablePeers = append(db.unreachablePeers, pi)
	}
	return nil
}

func (db *mockDatabase) RemoveUnreachablePeer(pi PeerIdentity) error {
	if exist, _ := db.ContainsUnreachablePeer(pi); exist {
		for i := 0; i < len(db.unreachablePeers); i++ {
			if db.unreachablePeers[i].PublicKey().Equals(pi.PublicKey()) {
				db.unreachablePeers = db.unreachablePeers[:i+copy(db.unreachablePeers[i:], db.unreachablePeers[i+1:])]
			}
		}
	}
	return nil
}

func (db *mockDatabase) ContainsUnreachablePeer(pi PeerIdentity) (bool, error) {
	for _, up := range db.unreachablePeers {
		if up.PublicKey().Equals(pi.PublicKey()) {
			return true, nil
		}
	}
	return false, nil
}

func (db *mockDatabase) containsPeer(p Peer) bool {
	mdiscoveredPeers := make(map[string]Peer, 0)
	for _, p := range db.discoveredPeers {
		mdiscoveredPeers[string(p.Identity().PublicKey().Bytes())] = p
	}

	_, exist := mdiscoveredPeers[string(p.Identity().PublicKey().Bytes())]
	return exist
}

type mockNotifier struct {
	reaches     []string
	unreaches   []string
	discoveries []Peer
}

func (n *mockNotifier) NotifyReachable(pk crypto.PublicKey) error {
	p, err := pk.Marshal()
	if err != nil {
		return err
	}
	n.reaches = append(n.reaches, string(p))
	return nil
}
func (n *mockNotifier) NotifyUnreachable(pk crypto.PublicKey) error {
	p, err := pk.Marshal()
	if err != nil {
		return err
	}
	n.unreaches = append(n.unreaches, string(p))
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
	logger logging.Logger
}

func (m mockMessengerWithSynFailure) SendSyn(target PeerIdentity, known []Peer) (reqPeers []PeerIdentity, discoveries []Peer, err error) {
	return nil, nil, ErrUnreachablePeer
}

func (m mockMessengerWithSynFailure) SendAck(target PeerIdentity, requested []Peer) error {
	return nil
}

type mockMessengerWithAckFailure struct {
	logger logging.Logger
}

func (m mockMessengerWithAckFailure) SendSyn(target PeerIdentity, known []Peer) (requested []PeerIdentity, discoveries []Peer, err error) {
	reqP := NewPeerIdentity(net.ParseIP("200.18.186.39"), 3000, pubx)

	hb := NewPeerHeartbeatState(time.Now(), 0)
	as := NewPeerAppState("1.0", OkPeerStatus, 30.0, 10.0, "", 0, 1, 0)

	np1 := NewDiscoveredPeer(
		NewPeerIdentity(net.ParseIP("35.200.100.2"), 3000, puby),
		hb, as,
	)

	return []PeerIdentity{reqP}, []Peer{np1}, nil
}

func (m mockMessengerWithAckFailure) SendAck(target PeerIdentity, requested []Peer) error {
	return ErrUnreachablePeer
}

type mockMessengerUnexpectedFailure struct {
	logger logging.Logger
}

func (m mockMessengerUnexpectedFailure) SendSyn(target PeerIdentity, known []Peer) (unknown []Peer, new []Peer, err error) {
	return nil, nil, errors.New("Unexpected")
}

func (m mockMessengerUnexpectedFailure) SendAck(target PeerIdentity, requested []Peer) error {
	return nil
}

type mockMessenger struct {
	logger logging.Logger
}

func (m mockMessenger) SendSyn(target PeerIdentity, known []Peer) (requested []PeerIdentity, discoveries []Peer, err error) {
	reqP := NewPeerIdentity(net.ParseIP("200.18.186.39"), 3000, pubx)

	np1 := NewDiscoveredPeer(
		NewPeerIdentity(net.ParseIP("35.200.100.2"), 3000, puby),
		NewPeerHeartbeatState(time.Now(), 0),
		NewPeerAppState("1.0", OkPeerStatus, 50.1, 22.1, "", 0, 1, 0),
	)

	return []PeerIdentity{reqP}, []Peer{np1}, nil
}

func (m mockMessenger) SendAck(target PeerIdentity, requested []Peer) error {
	return nil
}

package discovery

import (
	"log"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

/*
Scenario: Store local peer
	Given a miner startup
	When I want to init the local peer
	Then I store it as owned peer
*/
func TestStoreLocalPeer(t *testing.T) {
	repo := &mockRepository{}
	s := Service{
		repo: repo,
		mon:  mockPeerMonitor{},
	}

	p, err := s.StoreLocalPeer("key", 3001, "1.0")
	assert.Nil(t, err)
	assert.Equal(t, "key", p.Identity().PublicKey())
	assert.Equal(t, 3001, p.Identity().Port())
	assert.True(t, p.IsLocal())

	assert.Len(t, repo.knownPeers, 1)
}

/*
Scenario: Store and notifies cycle discovered peers
	Given a gossip cycle
	When it has discovered new peers
	Then I store them and notify them
*/
func TestHandlingCycleDiscoveries(t *testing.T) {

	c := cycle{
		discoveryChan: make(chan Peer),
	}

	repo := &mockRepository{}
	notif := &mockNotifier{}

	srv := Service{
		repo:  repo,
		notif: notif,
	}
	discoveryChan := make(chan Peer)
	errChan := make(chan error)

	go srv.handleDiscoveries(c, discoveryChan, errChan)
	c.discoveryChan <- NewDiscoveredPeer(
		NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "key"),
		NewPeerHeartbeatState(time.Now(), 0),
		NewPeerAppState("1.0", OkPeerStatus, 30.0, 10.0, "", 200, 1, 0),
	)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for d := range discoveryChan {
			assert.Equal(t, "key", d.Identity().PublicKey())
			wg.Done()
			close(discoveryChan)
			close(errChan)
		}
	}()
	wg.Wait()

	assert.Len(t, repo.knownPeers, 1)
	assert.Equal(t, "key", repo.knownPeers[0].Identity().PublicKey())
	assert.Equal(t, "key", notif.discoveries[0].Identity().PublicKey())
}

/*
Scenario: Store and notifies cycle unreachables peers
	Given a gossip cycle
	When the cycle notifies a unreachable peer
	Then I store them and notify them
*/
func TestHandlingCycleUnreaches(t *testing.T) {
	c := cycle{
		unreachChan: make(chan PeerIdentity),
	}

	repo := &mockRepository{}
	notif := &mockNotifier{}

	srv := Service{
		repo:  repo,
		notif: notif,
	}
	unreachChan := make(chan PeerIdentity)
	errChan := make(chan error)

	go srv.handleUnreachables(c, unreachChan, errChan)
	c.unreachChan <- NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "key")

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for d := range unreachChan {
			assert.Equal(t, "key", d.PublicKey())
			wg.Done()

		}
	}()
	wg.Wait()
	close(unreachChan)
	close(errChan)

	assert.Len(t, repo.unreachablePeers, 1)
	assert.Equal(t, "key", repo.unreachablePeers[0])
	assert.Equal(t, "key", notif.unreaches[0].PublicKey())

}

/*
Scenario: Store and notifies cycle reachables peers
	Given a gossip cycle
	When the cycle notifies a areachable peer
	Then I store them and notify them
*/
func TestHandlingCycleReaches(t *testing.T) {
	c := cycle{
		reachChan: make(chan PeerIdentity),
	}

	repo := &mockRepository{}
	notif := &mockNotifier{}

	srv := Service{
		repo:  repo,
		notif: notif,
	}
	reachChan := make(chan PeerIdentity)
	errChan := make(chan error)

	var wg sync.WaitGroup
	wg.Add(1)

	go srv.handleReachables(c, reachChan, errChan)
	go func() {
		for d := range reachChan {
			assert.Equal(t, "key", d.PublicKey())
			wg.Done()
		}
	}()

	c.reachChan <- NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "key")

	wg.Wait()
	close(reachChan)
	close(errChan)

	assert.Len(t, repo.unreachablePeers, 0)
	assert.Equal(t, "key", notif.reaches[0].PublicKey())

}

/*
Scenario: Store and notifies cycle reachables peers after unreach
	Given a gossip cycle
	When the cycle notifies a unreach peer and after a reach peer
	Then I store them and notify them
*/
func TestHandlingCycleReachesAfterUnreach(t *testing.T) {
	c := cycle{
		reachChan:   make(chan PeerIdentity),
		unreachChan: make(chan PeerIdentity),
	}

	repo := &mockRepository{}
	notif := &mockNotifier{}

	srv := Service{
		repo:  repo,
		notif: notif,
	}
	unreachChan := make(chan PeerIdentity)
	reachChan := make(chan PeerIdentity)
	errChan := make(chan error)

	var wg sync.WaitGroup
	wg.Add(1)

	go srv.handleUnreachables(c, unreachChan, errChan)
	go func() {
		for d := range unreachChan {
			assert.Equal(t, "key", d.PublicKey())
			wg.Done()
		}
	}()

	c.unreachChan <- NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "key")

	wg.Wait()
	close(unreachChan)

	var wg2 sync.WaitGroup
	wg2.Add(1)
	go srv.handleReachables(c, reachChan, errChan)
	go func() {
		for d := range reachChan {
			assert.Equal(t, "key", d.PublicKey())
			wg2.Done()
		}
	}()

	c.reachChan <- NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "key")

	wg2.Wait()
	close(reachChan)
	close(errChan)

	assert.Len(t, repo.unreachablePeers, 0)
	assert.Equal(t, "key", notif.reaches[0].PublicKey())
}

func TestSpreadGossip(t *testing.T) {
	repo := &mockRepository{}
	notif := &mockNotifier{}
	s := NewService(repo, mockClient{}, notif, mockPeerNetworker{}, mockPeerMonitor{})

	ownPeer := NewLocalPeer("key", net.ParseIP("127.0.0.1"), 3000, "", 30.0, 10.0)
	seeds := []PeerIdentity{
		NewPeerIdentity(net.ParseIP("10.0.0.1"), 3000, "key1"),
	}

	dChan := make(chan Peer)
	rChan := make(chan PeerIdentity)
	uChan := make(chan PeerIdentity)
	errChan := make(chan error)
	s.startCycle(ownPeer, seeds, dChan, rChan, uChan, errChan)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for p := range dChan {
			log.Print(p.Identity().PublicKey())
			wg.Done()
		}
	}()

	wg.Wait()

	assert.Len(t, repo.knownPeers, 2)
	assert.Equal(t, "dKey1", repo.knownPeers[1].Identity().PublicKey())
}

/*
Scenario: Compare peers with different key and get the unknown
	Given a known peer and a different peer
	When I want to get the unknown peer
	Then I get the second peer
*/
func TestGetUnknownPeersWithDifferentKey(t *testing.T) {
	kp := NewPeerDigest(
		NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "key1"),
		NewPeerHeartbeatState(time.Now(), 1000),
	)

	comparee := NewPeerDigest(
		NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "key2"),
		NewPeerHeartbeatState(time.Now(), 1000),
	)

	unkwown := Service{}.getUnknownPeers([]Peer{kp}, []Peer{comparee})
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
	kp := NewPeerDigest(
		NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "key1"),
		NewPeerHeartbeatState(time.Now(), 1000),
	)

	comparee := NewPeerDigest(
		NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "key1"),
		NewPeerHeartbeatState(time.Now(), 1000),
	)

	unkwown := Service{}.getUnknownPeers([]Peer{kp}, []Peer{comparee})
	assert.Empty(t, unkwown, 1)
}

/*
Scenario: Compare 2 set of peers with different time and get the recent one
	Given known peers and received peers with different elapsed heartbeats
	When I want to get the unknown peers
	Then I get the peer with the highest elapsed heartbeats
*/
func TestGetUnknownPeersMoreRecent(t *testing.T) {
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

	unkwown := Service{}.getUnknownPeers([]Peer{kp1, kp2}, []Peer{comparee1, comparee2})
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
	kp := NewPeerDigest(
		NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "key1"),
		NewPeerHeartbeatState(time.Now(), 1000),
	)

	comparee := NewPeerDigest(
		NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "key2"),
		NewPeerHeartbeatState(time.Now(), 1000),
	)

	news := Service{}.getNewPeers([]Peer{kp}, []Peer{comparee})
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
	kp1 := NewPeerDigest(
		NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "key1"),
		NewPeerHeartbeatState(time.Now(), 1000),
	)
	kp2 := NewPeerDigest(
		NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "key2"),
		NewPeerHeartbeatState(time.Now(), 1200),
	)

	comparee1 := NewPeerDigest(
		NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "key2"),
		NewPeerHeartbeatState(time.Now(), 1000),
	)
	comparee2 := NewPeerDigest(
		NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "key1"),
		NewPeerHeartbeatState(time.Now(), 1000),
	)

	news := Service{}.getNewPeers([]Peer{kp1, kp2}, []Peer{comparee1, comparee2})
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
	kp := NewPeerDigest(
		NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "key1"),
		NewPeerHeartbeatState(time.Now(), 1000),
	)

	comparee := NewPeerDigest(
		NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "key1"),
		NewPeerHeartbeatState(time.Now(), 1000),
	)

	news := Service{}.getNewPeers([]Peer{kp}, []Peer{comparee})
	assert.Empty(t, news, 1)
}

type mockRepository struct {
	seedPeers        []PeerIdentity
	knownPeers       []Peer
	unreachablePeers []string
}

func (r *mockRepository) CountKnownPeers() (int, error) {
	return len(r.knownPeers), nil
}

func (r *mockRepository) ListSeedPeers() ([]PeerIdentity, error) {
	return r.seedPeers, nil
}

func (r *mockRepository) ListKnownPeers() ([]Peer, error) {
	return r.knownPeers, nil
}

func (r *mockRepository) StoreKnownPeer(peer Peer) error {
	if r.containsPeer(peer) {
		for _, p := range r.knownPeers {
			if p.Identity().PublicKey() == peer.Identity().PublicKey() {
				p = peer
				break
			}
		}
	} else {
		r.knownPeers = append(r.knownPeers, peer)
	}
	return nil
}

func (r *mockRepository) ListReachablePeers() ([]PeerIdentity, error) {
	pp := make([]PeerIdentity, 0)
	for i := 0; i < len(r.knownPeers); i++ {
		if !r.ContainsUnreachablePeer(r.knownPeers[i].Identity().PublicKey()) {
			pp = append(pp, r.knownPeers[i].Identity())
		}
	}
	return pp, nil
}

func (r *mockRepository) ListUnreachablePeers() ([]PeerIdentity, error) {
	pp := make([]PeerIdentity, 0)

	for i := 0; i < len(r.seedPeers); i++ {
		if r.ContainsUnreachablePeer(r.seedPeers[i].PublicKey()) {
			pp = append(pp, r.seedPeers[i])
		}
	}

	for i := 0; i < len(r.knownPeers); i++ {
		if r.ContainsUnreachablePeer(r.knownPeers[i].Identity().PublicKey()) {
			pp = append(pp, r.knownPeers[i].Identity())
		}
	}
	return pp, nil
}

func (r *mockRepository) StoreSeedPeer(s PeerIdentity) error {
	r.seedPeers = append(r.seedPeers, s)
	return nil
}

func (r *mockRepository) StoreUnreachablePeer(pk string) error {
	if !r.ContainsUnreachablePeer(pk) {
		r.unreachablePeers = append(r.unreachablePeers, pk)
	}
	return nil
}

func (r *mockRepository) RemoveUnreachablePeer(pk string) error {
	if r.ContainsUnreachablePeer(pk) {
		for i := 0; i < len(r.unreachablePeers); i++ {
			if r.unreachablePeers[i] == pk {
				r.unreachablePeers = r.unreachablePeers[:i+copy(r.unreachablePeers[i:], r.unreachablePeers[i+1:])]
			}
		}
	}
	return nil
}

func (r *mockRepository) ContainsUnreachablePeer(peerPubk string) bool {
	for _, up := range r.unreachablePeers {
		if up == peerPubk {
			return true
		}
	}
	return false
}

func (r *mockRepository) containsPeer(p Peer) bool {
	mdiscoveredPeers := make(map[string]Peer, 0)
	for _, p := range r.knownPeers {
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

type mockPeerNetworker struct{}

func (pn mockPeerNetworker) CheckNtpState() error {
	return nil
}

func (pn mockPeerNetworker) CheckInternetState() error {
	return nil
}

type mockPeerMonitor struct{}

func (i mockPeerMonitor) GeoPosition() (lon float64, lat float64, err error) {
	return 10.0, 30.0, nil
}

func (i mockPeerMonitor) FreeDiskSpace() (float64, error) {
	return 200, nil
}

func (i mockPeerMonitor) CPULoad() (string, error) {
	return "", nil
}

func (i mockPeerMonitor) IP() (net.IP, error) {
	return net.ParseIP("127.0.0.1"), nil
}

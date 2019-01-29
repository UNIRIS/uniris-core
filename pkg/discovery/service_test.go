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
		repo:  repo,
		pInfo: mockPeerInfo{},
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
		unreachChan: make(chan Peer),
	}

	repo := &mockRepository{}
	notif := &mockNotifier{}

	srv := Service{
		repo:  repo,
		notif: notif,
	}
	unreachChan := make(chan Peer)
	errChan := make(chan error)

	go srv.handleUnreachables(c, unreachChan, errChan)
	c.unreachChan <- NewDiscoveredPeer(
		NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "key"),
		NewPeerHeartbeatState(time.Now(), 0),
		NewPeerAppState("1.0", OkPeerStatus, 30.0, 10.0, "", 200, 1, 0),
	)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for d := range unreachChan {
			assert.Equal(t, "key", d.Identity().PublicKey())
			wg.Done()

		}
	}()
	wg.Wait()
	close(unreachChan)
	close(errChan)

	assert.Len(t, repo.unreachablePeers, 1)
	assert.Equal(t, "key", repo.unreachablePeers[0])
	assert.Equal(t, "key", notif.unreaches[0].Identity().PublicKey())

}

/*
Scenario: Store and notifies cycle reachables peers
	Given a gossip cycle
	When the cycle notifies a areachable peer
	Then I store them and notify them
*/
func TestHandlingCycleReaches(t *testing.T) {
	c := cycle{
		reachChan: make(chan Peer),
	}

	repo := &mockRepository{}
	notif := &mockNotifier{}

	srv := Service{
		repo:  repo,
		notif: notif,
	}
	reachChan := make(chan Peer)
	errChan := make(chan error)

	var wg sync.WaitGroup
	wg.Add(1)

	go srv.handleReachables(c, reachChan, errChan)
	go func() {
		for d := range reachChan {
			assert.Equal(t, "key", d.Identity().PublicKey())
			wg.Done()
		}
	}()

	c.reachChan <- NewDiscoveredPeer(
		NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "key"),
		NewPeerHeartbeatState(time.Now(), 0),
		NewPeerAppState("1.0", OkPeerStatus, 30.0, 10.0, "", 200, 1, 0),
	)

	wg.Wait()
	close(reachChan)
	close(errChan)

	assert.Len(t, repo.unreachablePeers, 0)
	assert.Equal(t, "key", notif.reaches[0].Identity().PublicKey())

}

/*
Scenario: Store and notifies cycle reachables peers after unreach
	Given a gossip cycle
	When the cycle notifies a unreach peer and after a reach peer
	Then I store them and notify them
*/
func TestHandlingCycleReachesAfterUnreach(t *testing.T) {
	c := cycle{
		reachChan:   make(chan Peer),
		unreachChan: make(chan Peer),
	}

	repo := &mockRepository{}
	notif := &mockNotifier{}

	srv := Service{
		repo:  repo,
		notif: notif,
	}
	unreachChan := make(chan Peer)
	reachChan := make(chan Peer)
	errChan := make(chan error)

	var wg sync.WaitGroup
	wg.Add(1)

	go srv.handleUnreachables(c, unreachChan, errChan)
	go func() {
		for d := range unreachChan {
			assert.Equal(t, "key", d.Identity().PublicKey())
			wg.Done()
		}
	}()

	c.unreachChan <- NewDiscoveredPeer(
		NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "key"),
		NewPeerHeartbeatState(time.Now(), 0),
		NewPeerAppState("1.0", OkPeerStatus, 30.0, 10.0, "", 200, 1, 0),
	)

	wg.Wait()
	close(unreachChan)

	var wg2 sync.WaitGroup
	wg2.Add(1)
	go srv.handleReachables(c, reachChan, errChan)
	go func() {
		for d := range reachChan {
			assert.Equal(t, "key", d.Identity().PublicKey())
			wg2.Done()
		}
	}()

	c.reachChan <- NewDiscoveredPeer(
		NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "key"),
		NewPeerHeartbeatState(time.Now(), 0),
		NewPeerAppState("1.0", OkPeerStatus, 30.0, 10.0, "", 200, 1, 0),
	)

	wg2.Wait()
	close(reachChan)
	close(errChan)

	assert.Len(t, repo.unreachablePeers, 0)
	assert.Equal(t, "key", notif.reaches[0].Identity().PublicKey())
}

func TestSpreadGossip(t *testing.T) {
	repo := &mockRepository{}
	notif := &mockNotifier{}
	s := NewService(repo, mockMessenger{}, notif, mockPeerNetworker{}, mockPeerInfo{})

	ownPeer := NewLocalPeer("key", net.ParseIP("127.0.0.1"), 3000, "", 30.0, 10.0)
	seeds := []Seed{
		Seed{
			PeerIdentity: NewPeerIdentity(net.ParseIP("10.0.0.1"), 3000, "key1"),
		},
	}

	dChan := make(chan Peer)
	rChan := make(chan Peer)
	uChan := make(chan Peer)
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

type mockRepository struct {
	seedPeers        []Seed
	knownPeers       []Peer
	unreachablePeers []string
}

func (r *mockRepository) CountKnownPeers() (int, error) {
	return len(r.knownPeers), nil
}

func (r *mockRepository) ListSeedPeers() ([]Seed, error) {
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

func (r *mockRepository) ListReachablePeers() ([]Peer, error) {
	pp := make([]Peer, 0)
	for i := 0; i < len(r.knownPeers); i++ {
		if !r.ContainsUnreachablePeer(r.knownPeers[i].Identity().PublicKey()) {
			pp = append(pp, r.knownPeers[i])
		}
	}
	return pp, nil
}

func (r *mockRepository) ListUnreachablePeers() ([]Peer, error) {
	pp := make([]Peer, 0)

	for i := 0; i < len(r.seedPeers); i++ {
		if r.ContainsUnreachablePeer(r.seedPeers[i].PublicKey()) {
			pp = append(pp, r.seedPeers[i].AsPeer())
		}
	}

	for i := 0; i < len(r.knownPeers); i++ {
		if r.ContainsUnreachablePeer(r.knownPeers[i].Identity().PublicKey()) {
			pp = append(pp, r.knownPeers[i])
		}
	}
	return pp, nil
}

func (r *mockRepository) StoreSeedPeer(s Seed) error {
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
	reaches     []Peer
	unreaches   []Peer
	discoveries []Peer
}

func (n *mockNotifier) NotifyReachable(p Peer) error {
	n.reaches = append(n.reaches, p)
	return nil
}
func (n *mockNotifier) NotifyUnreachable(p Peer) error {
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

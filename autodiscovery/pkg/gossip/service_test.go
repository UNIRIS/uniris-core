package gossip

import (
	"encoding/hex"
	"errors"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	discovery "github.com/uniris/uniris-core/autodiscovery/pkg"
	"github.com/uniris/uniris-core/autodiscovery/pkg/monitoring"
)

/*
Scenario: Gossip across a selection of peers
	Given a initiator peer, seeds and known peers stored locally
	When we gossip
	Then the new peers are stored and notified
*/
func TestGossip(t *testing.T) {

	repo := new(mockRepository)
	notif := new(notifier)
	msg := new(mockMessenger)
	mon := monitoring.NewService(repo, new(monitor), new(networker), new(robotWatcher))

	repo.SetSeed(discovery.Seed{IP: net.ParseIP("30.0.0.1"), Port: 3000})

	init := discovery.NewStartupPeer([]byte("key"), net.ParseIP("127.0.0.1"), 3000, "1.0", discovery.PeerPosition{})
	repo.SetPeer(init)

	s := service{
		msg:   msg,
		repo:  repo,
		notif: notif,
		mon:   mon,
	}

	seeds, _ := repo.ListSeedPeers()

	errs := make(chan error)
	newP := make(chan discovery.Peer)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		s.spread(init, seeds, errs, newP)
		for range newP {
			wg.Done()
		}
	}()

	wg.Wait()

	assert.Empty(t, errs)

	pp, _ := repo.ListDiscoveredPeers()
	assert.NotEmpty(t, pp)
	assert.Equal(t, "dKey1", pp[0].Identity().PublicKey().String())
}

/*
Scenario: Gossip with unexpected error
	Given a gossip round spread fails
	When the error is not about the target cannot be reach
	Then the error is catched
*/
func TestGossipFailureCatched(t *testing.T) {

	repo := new(mockRepository)
	notif := new(notifier)
	msg := new(mockMessengerUnexpectedFailure)
	mon := monitoring.NewService(repo, new(monitor), new(networker), new(robotWatcher))

	repo.SetSeed(discovery.Seed{IP: net.ParseIP("30.0.0.1"), Port: 3000})

	init := discovery.NewStartupPeer([]byte("key"), net.ParseIP("127.0.0.1"), 3000, "1.0", discovery.PeerPosition{})
	repo.SetPeer(init)

	s := service{
		msg:   msg,
		repo:  repo,
		notif: notif,
		mon:   mon,
	}

	seeds, _ := repo.ListSeedPeers()

	errs := make(chan error)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		s.spread(init, seeds, errs, nil)
		for err := range errs {
			assert.Error(t, err, "Unexpected failure")
			wg.Done()
		}
	}()

	wg.Wait()

}

//////////////////////////////////////////////////////////
// 						MOCKS
/////////////////////////////////////////////////////////

type mockMessenger struct {
}

func (m mockMessenger) SendSyn(req SynRequest) (*SynAck, error) {
	init := discovery.NewStartupPeer([]byte("key"), net.ParseIP("127.0.0.1"), 3000, "1.0", discovery.PeerPosition{})
	tar := discovery.NewStartupPeer([]byte("uKey1"), net.ParseIP("200.18.186.39"), 3000, "1.1", discovery.PeerPosition{})

	hb := discovery.NewPeerHeartbeatState(time.Now(), 0)
	as := discovery.NewPeerAppState("1.0", discovery.OkStatus, discovery.PeerPosition{}, "", 0, 1, 0)

	np1 := discovery.NewDiscoveredPeer(
		discovery.NewPeerIdentity(net.ParseIP("35.200.100.2"), 3000, []byte("dKey1")),
		hb, as,
	)

	newPeers := []discovery.Peer{np1}

	unknownPeers := []discovery.Peer{tar}

	return &SynAck{
		Initiator:    init,
		Target:       tar,
		NewPeers:     newPeers,
		UnknownPeers: unknownPeers,
	}, nil
}

func (m mockMessenger) SendAck(req AckRequest) error {
	return nil
}

type mockRepository struct {
	ownedPeer       discovery.Peer
	discoveredPeers []discovery.Peer
	seedPeers       []discovery.Seed
}

func (r *mockRepository) CountDiscoveredPeers() (int, error) {
	return len(r.discoveredPeers), nil
}

//GetOwnedPeer return the local peer
func (r *mockRepository) GetOwnedPeer() (discovery.Peer, error) {
	return r.ownedPeer, nil
}

//ListSeedPeers return all the seed on the mockRepository
func (r *mockRepository) ListSeedPeers() ([]discovery.Seed, error) {
	return r.seedPeers, nil
}

//ListDiscoveredPeers returns all the discoveredPeers on the mockRepository
func (r *mockRepository) ListDiscoveredPeers() ([]discovery.Peer, error) {
	return r.discoveredPeers, nil
}

func (r *mockRepository) SetPeer(peer discovery.Peer) error {
	if peer.Owned() {
		r.ownedPeer = peer
		return nil
	}
	if r.containsPeer(peer) {
		for _, p := range r.discoveredPeers {
			if p.Identity().PublicKey().Equals(peer.Identity().PublicKey()) {
				p = peer
				break
			}
		}
	} else {
		r.discoveredPeers = append(r.discoveredPeers, peer)
	}
	return nil
}

func (r *mockRepository) SetSeed(s discovery.Seed) error {
	r.seedPeers = append(r.seedPeers, s)
	return nil
}

//GetPeerByIP get a peer from the mockRepository using its ip
func (r *mockRepository) GetPeerByIP(ip net.IP) (p discovery.Peer, err error) {
	if r.ownedPeer.Identity().IP().Equal(ip) {
		return r.ownedPeer, nil
	}
	for i := 0; i < len(r.discoveredPeers); i++ {
		if r.discoveredPeers[i].Identity().IP().Equal(ip) {
			return r.discoveredPeers[i], nil
		}
	}
	return
}

func (r *mockRepository) containsPeer(p discovery.Peer) bool {
	mdiscoveredPeers := make(map[string]discovery.Peer, 0)
	for _, p := range r.discoveredPeers {
		mdiscoveredPeers[hex.EncodeToString(p.Identity().PublicKey())] = p
	}

	_, exist := mdiscoveredPeers[hex.EncodeToString(p.Identity().PublicKey())]
	return exist
}

type notifier struct {
	notifiedPeers []discovery.Peer
}

func (n notifier) NotifiedPeers() []discovery.Peer {
	return n.notifiedPeers
}

func (n *notifier) Notify(p discovery.Peer) {
	n.notifiedPeers = append(n.notifiedPeers, p)
}

type monitor struct{}

func (w monitor) CPULoad() (string, error) {
	return "0.62 0.77 0.71 4/972 26361", nil
}

func (w monitor) FreeDiskSpace() (float64, error) {
	return 212383852, nil
}

func (w monitor) P2PFactor() (int, error) {
	return 1, nil
}

type networker struct{}

func (n networker) IP() (net.IP, error) {
	return net.ParseIP("127.0.0.1"), nil
}

func (n networker) CheckInternetState() error {
	return nil
}

func (n networker) CheckNtpState() error {
	return nil
}

type networkerNTPFails struct{}

func (n networkerNTPFails) IP() (net.IP, error) {
	return net.ParseIP("127.0.0.1"), nil
}

func (n networkerNTPFails) CheckInternetState() error {
	return nil
}

func (n networkerNTPFails) CheckNtpState() error {
	return errors.New("System Clock have a big Offset check the ntp configuration of the system")
}

type networkerInternetFails struct{}

func (n networkerInternetFails) IP() (net.IP, error) {
	return net.ParseIP("127.0.0.1"), nil
}

func (n networkerInternetFails) CheckInternetState() error {
	return errors.New("required processes are not running")
}

func (n networkerInternetFails) CheckNtpState() error {
	return nil
}

type robotWatcher struct{}

func (r robotWatcher) CheckAutodiscoveryProcess(port int) error {
	return nil
}

func (r robotWatcher) CheckDataProcess() error {
	return nil
}

func (r robotWatcher) CheckMiningProcess() error {
	return nil
}

func (r robotWatcher) CheckAIProcess() error {
	return nil
}

func (r robotWatcher) CheckScyllaDbProcess() error {
	return nil
}

func (r robotWatcher) CheckRedisProcess() error {
	return nil
}

func (r robotWatcher) CheckRabbitmqProcess() error {
	return nil
}

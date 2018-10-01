package gossip

import (
	"net"
	"sync"
	"testing"
	"time"

	"github.com/uniris/uniris-core/autodiscovery/pkg/mock"

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

	repo := new(mock.Repository)
	notif := new(mock.Notifier)
	msg := new(mockMessenger)
	mon := monitoring.NewService(repo, new(mock.Monitor), new(mock.Networker), new(mock.RobotWatcher))

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
		for range s.spread(init, seeds, errs) {
			wg.Done()
		}
	}()

	go func() {
		for range errs {
			assert.Fail(t, "Cannot have errors")
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

	repo := new(mock.Repository)
	notif := new(mock.Notifier)
	msg := new(mockMessengerUnexpectedFailure)
	mon := monitoring.NewService(repo, new(mock.Monitor), new(mock.Networker), new(mock.RobotWatcher))

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
		s.spread(init, seeds, errs)

	}()

	go func() {
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

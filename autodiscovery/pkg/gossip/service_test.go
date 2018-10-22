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

	repo.SetSeedPeer(discovery.Seed{IP: net.ParseIP("30.0.0.1"), Port: 3000})

	init := discovery.NewStartupPeer("key", net.ParseIP("127.0.0.1"), 3000, "1.0", discovery.PeerPosition{})
	repo.SetKnownPeer(init)

	s := service{
		msg:   msg,
		repo:  repo,
		notif: notif,
		mon:   mon,
	}

	seeds, _ := repo.ListSeedPeers()

	var wg sync.WaitGroup
	wg.Add(1)

	res := NewSpreadResult()

	go s.spread(init, seeds, res.Discoveries, res.Unreaches, res.Errors)

	go func() {
		for range res.Discoveries {
			wg.Done()
		}
	}()

	go func() {
		for range res.Errors {
			assert.Fail(t, "Cannot have errors")
		}
	}()

	wg.Wait()

	res.CloseChannels()

	pp, _ := repo.ListKnownPeers()
	assert.NotEmpty(t, pp)
	assert.Len(t, pp, 2)
	assert.Equal(t, "key", pp[0].Identity().PublicKey())
	assert.Equal(t, "dKey1", pp[1].Identity().PublicKey())
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

	repo.SetSeedPeer(discovery.Seed{IP: net.ParseIP("30.0.0.1"), Port: 3000})

	init := discovery.NewStartupPeer("key", net.ParseIP("127.0.0.1"), 3000, "1.0", discovery.PeerPosition{})
	repo.SetKnownPeer(init)

	s := service{
		msg:   msg,
		repo:  repo,
		notif: notif,
		mon:   mon,
	}

	seeds, _ := repo.ListSeedPeers()

	var wg sync.WaitGroup
	wg.Add(1)

	res := NewSpreadResult()

	go s.spread(init, seeds, res.Discoveries, res.Unreaches, res.Errors)

	errs := make([]error, 0)
	go func() {
		for e := range res.Errors {
			errs = append(errs, e)
			wg.Done()
		}
	}()

	wg.Wait()

	res.CloseChannels()

	assert.NotEmpty(t, errs)
	assert.Equal(t, errs[0].Error(), "Unexpected")
}

/*
Scenario: Gossip across a selection of peers
	Given a initiator peer, seeds and known peers stored locally
	When we gossip spread get unreacheable error
	Then unreacheable peer is stored on the repo
*/
func TestAddUnreachable(t *testing.T) {
	repo := new(mock.Repository)
	notif := new(mock.Notifier)
	msg := new(mockMessengerWithSynFailure)
	mon := monitoring.NewService(repo, new(mock.Monitor), new(mock.Networker), new(mock.RobotWatcher))

	seed := discovery.Seed{IP: net.ParseIP("30.0.0.1"), Port: 3000, PublicKey: "key2"}
	repo.SetSeedPeer(seed)

	init := discovery.NewStartupPeer("key", net.ParseIP("127.0.0.1"), 3000, "1.0", discovery.PeerPosition{})
	repo.SetKnownPeer(init)

	s := service{
		msg:   msg,
		repo:  repo,
		notif: notif,
		mon:   mon,
	}

	seeds, _ := repo.ListSeedPeers()

	var wg sync.WaitGroup
	wg.Add(1)

	res := NewSpreadResult()

	go s.spread(init, seeds, res.Discoveries, res.Unreaches, res.Errors)

	go func() {
		for range res.Unreaches {
			wg.Done()
		}
	}()

	wg.Wait()

	res.CloseChannels()

	unreaches, _ := repo.ListUnreachablePeers()
	assert.NotEmpty(t, unreaches)
}

/*
Scenario: Gossip across a selection of peers
	Given a initiator peer, seeds and known peers stored locally and one unreacheable peer
	When we gossip
	Then unreacheable peer is removed from the repo
*/
func TestRemoveUnreachable(t *testing.T) {
	repo := new(mock.Repository)
	notif := new(mock.Notifier)
	msg := new(mockMessengerWithSynFailure)
	mon := monitoring.NewService(repo, new(mock.Monitor), new(mock.Networker), new(mock.RobotWatcher))

	seed := discovery.Seed{IP: net.ParseIP("30.0.0.1"), Port: 3000, PublicKey: "dKey1"}
	repo.SetSeedPeer(seed)

	init := discovery.NewStartupPeer("key", net.ParseIP("127.0.0.1"), 3000, "1.0", discovery.PeerPosition{})
	repo.SetKnownPeer(init)

	s := service{
		msg:   msg,
		repo:  repo,
		notif: notif,
		mon:   mon,
	}

	seeds, _ := repo.ListSeedPeers()

	var wg sync.WaitGroup
	wg.Add(1)

	res := NewSpreadResult()
	go s.spread(init, seeds, res.Discoveries, res.Unreaches, res.Errors)

	go func() {
		for range res.Unreaches {
			wg.Done()
		}
	}()

	wg.Wait()

	unreaches, _ := repo.ListUnreachablePeers()
	assert.NotEmpty(t, unreaches)

	msg2 := mockMessenger{}
	s = service{
		msg:   msg2,
		repo:  repo,
		notif: notif,
		mon:   mon,
	}

	var wg2 sync.WaitGroup
	wg2.Add(2)

	go s.spread(init, seeds, res.Discoveries, res.Unreaches, res.Errors)

	go func() {
		for range res.Discoveries {
			wg2.Done()
		}
	}()

	wg2.Wait()

	res.CloseChannels()

	unreaches, _ = repo.ListUnreachablePeers()
	assert.Empty(t, unreaches)

	reaches, _ := repo.ListReachablePeers()
	assert.NotEmpty(t, reaches)
	assert.Equal(t, "dKey1", reaches[1].Identity().PublicKey())
}

/*
Scenario: Stop the timer when a error is returned during the gossip spreading
	Given a gossip starting
	When an unexpected error occurred
	Then the gossip is stopped
*/
func TestStopTimerWhenGossipError(t *testing.T) {
	repo := new(mock.Repository)
	notif := new(mock.Notifier)
	msg := new(mockMessengerUnexpectedFailure)
	mon := monitoring.NewService(repo, new(mock.Monitor), new(mock.Networker), new(mock.RobotWatcher))

	seed := discovery.Seed{IP: net.ParseIP("30.0.0.1"), Port: 3000, PublicKey: "dKey1"}
	repo.SetSeedPeer(seed)

	init := discovery.NewStartupPeer("key", net.ParseIP("127.0.0.1"), 3000, "1.0", discovery.PeerPosition{})
	repo.SetKnownPeer(init)

	srv := NewService(repo, msg, notif, mon)

	var wg sync.WaitGroup
	wg.Add(1)

	ticker := time.NewTicker(1 * time.Second)
	res, _ := srv.Start(init, ticker)

	errs := make([]error, 0)
	go func() {
		for err := range res.Errors {
			errs = append(errs, err)
		}
	}()

	go func() {
		for range res.Finish {
			wg.Done()
		}
	}()

	wg.Wait()

	assert.Empty(t, res.Discoveries)
	assert.Empty(t, res.Unreaches)
	assert.NotEmpty(t, errs)
	assert.Equal(t, "Unexpected", errs[0].Error())
}

type mockMessenger struct {
}

func (m mockMessenger) SendSyn(req SynRequest) (*SynAck, error) {
	init := discovery.NewStartupPeer("key", net.ParseIP("127.0.0.1"), 3000, "1.0", discovery.PeerPosition{})
	tar := discovery.NewStartupPeer("uKey1", net.ParseIP("200.18.186.39"), 3000, "1.1", discovery.PeerPosition{})

	hb := discovery.NewPeerHeartbeatState(time.Now(), 0)
	as := discovery.NewPeerAppState("1.0", discovery.OkStatus, discovery.PeerPosition{}, "", 0, 1, 0)

	np1 := discovery.NewDiscoveredPeer(
		discovery.NewPeerIdentity(net.ParseIP("35.200.100.2"), 3000, "dKey1"),
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

package gossip

import (
	"log"
	"time"

	discovery "github.com/uniris/uniris-core/autodiscovery/pkg"
	"github.com/uniris/uniris-core/autodiscovery/pkg/monitoring"
)

//Service is the interface that provide gossip methods
type Service interface {
	Start(discovery.Peer) error
}

//Notifier is the interface that provides methods to notify gossip discovery
type Notifier interface {
	Notify(discovery.Peer)
}

type service struct {
	repo  discovery.Repository
	msg   Messenger
	notif Notifier
	mon   monitoring.Service
}

type result struct {
	err     chan error
	success chan bool
}

//Start initialize the gossip session every seconds
//It stores the discovered peers and unreachables peers
//It calls the notifier to dispatch the discovered peers to the AI service
func (s service) Start(init discovery.Peer) error {

	seeds, err := s.repo.ListSeedPeers()
	if err != nil {
		return err
	}

	ticker := time.NewTicker(1 * time.Second)
	for range ticker.C {

		errs := make(chan error)
		newPeers := make(chan discovery.Peer)

		go s.spread(init, seeds, errs, newPeers)

		if err := <-errs; err != nil {
			close(errs)
			ticker.Stop()
			return err
		}
	}

	return nil
}

func (s service) spread(init discovery.Peer, seeds []discovery.Seed, errs chan<- error, newPeers chan discovery.Peer) {

	//Refreshes owned peer state before sending any requests
	if err := s.mon.RefreshPeer(init); err != nil {
		errs <- err
	}

	//DEBUG OWNED PEER
	selfp, err := s.repo.GetOwnedPeer()
	if err != nil {
		errs <- err
	}
	log.Printf("DEBUG: cpu: %s, freedisk: %b, status: %d, discoveredPeersNumber: %d", selfp.AppState().CPULoad(), selfp.AppState().FreeDiskSpace(), selfp.AppState().Status(), selfp.AppState().DiscoveredPeersNumber())

	dp, err := s.repo.ListDiscoveredPeers()
	if err != nil {
		errs <- err
	}

	c, err := NewGossipCycle(init, dp, seeds, s.msg)
	if err != nil {
		errs <- err
	}

	go c.Run()

	go s.handleCycleErrors(c, errs)
	go s.handleCycleDiscoveries(c, errs, newPeers)
}

func (s service) handleCycleErrors(c *Cycle, errs chan<- error) {
	for err := range c.result.errors {
		errs <- err
	}
}

func (s service) handleCycleDiscoveries(c *Cycle, errs chan<- error, newPeers chan<- discovery.Peer) {
	for p := range c.result.discoveries {
		if err := s.repo.SetPeer(p); err != nil {
			errs <- err
		}
		s.notif.Notify(p)
		newPeers <- p
	}
}

//NewService creates a gossiping service its dependencies
func NewService(repo discovery.Repository, msg Messenger, notif Notifier) Service {
	return service{
		repo:  repo,
		msg:   msg,
		notif: notif,
	}
}

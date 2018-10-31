package gossip

import (
	"errors"
	"log"
	"time"

	discovery "github.com/uniris/uniris-core/autodiscovery/pkg"
	"github.com/uniris/uniris-core/autodiscovery/pkg/monitoring"
)

var (
	ErrNotFoundOnUnreachableList = errors.New("cannot found the peer in the unreachableKeys list")
)

//Service is the interface that provide gossip methods
type Service interface {
	Start(discovery.Peer, *time.Ticker) (*SpreadResult, error)
}

//Notifier is the interface that provides methods to notify gossip discovery
type Notifier interface {
	NotifyDiscoveries(discovery.Peer) error
	NotifyReachable(pubk string) error
	NotifyUnreachable(pubk string) error
}

type service struct {
	repo  discovery.Repository
	msg   Messenger
	notif Notifier
	mon   monitoring.Service
}

//SpreadResult represents the gossig results from the peer starting up
type SpreadResult struct {
	Errors      chan error
	Discoveries chan discovery.Peer
	Unreaches   chan discovery.Peer
	Reaches     chan discovery.Peer
	Finish      chan bool
}

//NewSpreadResult creates a new result for the gossip
func NewSpreadResult() *SpreadResult {
	return &SpreadResult{
		Discoveries: make(chan discovery.Peer),
		Unreaches:   make(chan discovery.Peer),
		Errors:      make(chan error, 1),
		Finish:      make(chan bool),
	}
}

//CloseChannels closes the opened channels during the gossip
func (r *SpreadResult) CloseChannels() {
	close(r.Errors)
	close(r.Discoveries)
	close(r.Unreaches)
	close(r.Finish)
}

//Start initialize the gossip session
//It stores the discovered peers and unreachables peers
//It calls the notifier to dispatch the discovered peers to the AI service
func (s service) Start(init discovery.Peer, ticker *time.Ticker) (*SpreadResult, error) {
	res := NewSpreadResult()

	seeds, err := s.repo.ListSeedPeers()
	if err != nil {
		return nil, err
	}

	go func() {

		for range ticker.C {
			go s.spread(init, seeds, res.Discoveries, res.Unreaches, res.Errors)

			//Stop the ticker when an unexpected error is returned
			//We also close all the channels to prevent any unexpected storage
			go func() {
				for range res.Errors {
					ticker.Stop()
					res.Finish <- true
					return
				}
			}()
		}
	}()

	return res, nil
}

func (s service) spread(init discovery.Peer, seeds []discovery.Seed, dChan chan<- discovery.Peer, uChan chan<- discovery.Peer, eChan chan<- error) {

	//Refreshes owned peer state before sending any requests
	if err := s.mon.RefreshPeer(init); err != nil {
		eChan <- err
		return
	}

	knownPeers, err := s.repo.ListKnownPeers()
	if err != nil {
		eChan <- err
		return
	}

	rp, err := s.repo.ListReachablePeers()
	if err != nil {
		eChan <- err
		return
	}

	up, err := s.repo.ListUnreachablePeers()
	if err != nil {
		eChan <- err
		return
	}

	//add reachables, seeds, unreachables
	c := NewGossipCycle(init, s.msg)
	if err != nil {
		eChan <- err
		return
	}

	pp, err := c.SelectPeers(seeds, rp, up)
	if err != nil {
		eChan <- err
		return
	}

	go c.Run(init, pp, knownPeers)

	//Handle gossip cycle returns
	go s.handleCycleErrors(c, eChan)
	go s.handleCycleDiscoveries(c, dChan, eChan)
	go s.handleCycleUnreachables(c, uChan, eChan)
	go s.handleCycleReachables(c, eChan)
}

func (s service) handleCycleReachables(c *Cycle, eChan chan<- error) {
	for p := range c.result.reaches {
		log.Printf("Gossip reached peer: %s", p.Endpoint())

		err := s.repo.ContainsUnreachableKey(p.Identity().PublicKey())
		if err != nil && err != ErrNotFoundOnUnreachableList {
			eChan <- err
			return
		}
		if err == nil {
			//Remove the target from the unreachable list if it is
			if err := s.repo.RemoveUnreachablePeer(p.Identity().PublicKey()); err != nil {
				eChan <- err
				return
			}
			//Notify for the reachable peer
			if err := s.notif.NotifyReachable(p.Identity().PublicKey()); err != nil {
				eChan <- err
				return
			}
		}

	}
}

func (s service) handleCycleUnreachables(c *Cycle, uChan chan<- discovery.Peer, eChan chan<- error) {
	for p := range c.result.unreachables {
		log.Printf("Gossip unreached peer: %s", p.Endpoint())
		err := s.repo.ContainsUnreachableKey(p.Identity().PublicKey())
		if err != nil {
			if err != ErrNotFoundOnUnreachableList {
				eChan <- err
				return
			}
			if err := s.repo.SetUnreachablePeer(p.Identity().PublicKey()); err != nil {
				eChan <- err
				return
			}
			if err := s.notif.NotifyUnreachable(p.Identity().PublicKey()); err != nil {
				eChan <- err
				return
			}
			uChan <- p
		}
	}
}

func (s service) handleCycleErrors(c *Cycle, eChan chan<- error) {
	for err := range c.result.errors {
		eChan <- err
		return
	}
}

func (s service) handleCycleDiscoveries(c *Cycle, dChan chan<- discovery.Peer, eChan chan<- error) {
	for p := range c.result.discoveries {
		log.Printf("Gossip discovered new peer: %s", p.String())

		//Add or update the discovered peer
		if err := s.repo.SetKnownPeer(p); err != nil {
			eChan <- err
			return
		}

		//Notify the new discovery
		if err := s.notif.NotifyDiscoveries(p); err != nil {
			eChan <- err
			return
		}

		dChan <- p
	}
}

//NewService creates a gossiping service its dependencies
func NewService(repo discovery.Repository, msg Messenger, notif Notifier, mon monitoring.Service) Service {
	return service{
		repo:  repo,
		msg:   msg,
		notif: notif,
		mon:   mon,
	}
}

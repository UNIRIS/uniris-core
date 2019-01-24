package gossip

import (
	"time"

	uniris "github.com/uniris/uniris-core/pkg"
)

//Repository provides access to the discovery persistence
type Repository interface {
	StoreSeedPeer(uniris.Seed) error
	StoreUnreachablePeer(pubKey string) error
	StoreKnownPeer(uniris.Peer) error
	RemoveUnreachablePeer(pubKey string) error

	ListSeedPeers() ([]uniris.Seed, error)
	CountKnownPeers() (int, error)
	ListKnownPeers() ([]uniris.Peer, error)
	ListUnreachablePeers() ([]uniris.Peer, error)
	ListReachablePeers() ([]uniris.Peer, error)
	ContainsUnreachablePeer(peerPublicKey string) bool
}

//Service handle the gossip spread and discovery of the network
type Service struct {
	repo  Repository
	pInfo PeerInformer
	pNet  PeerNetworker
	msg   RoundMessenger
	notif Notifier
}

//Notifier handle the notification of the gossip events
type Notifier interface {
	NotifyReachable(uniris.Peer) error
	NotifyUnreachable(uniris.Peer) error
	NotifyDiscovery(uniris.Peer) error
}

//NewService creates the gossip handler
func NewService(repo Repository, msg RoundMessenger, notif Notifier, pNet PeerNetworker, pInfo PeerInformer) Service {
	return Service{
		repo:  repo,
		msg:   msg,
		notif: notif,
		pInfo: pInfo,
		pNet:  pNet,
	}
}

//BootStrapingMinTime is the necessary minimum time on seconds to finish learning about the network
const BootStrapingMinTime = 1800

//Run initiate the gossip process by spreading local informations and updated the receiving
func (s Service) Run(localPeer uniris.Peer, seeds []uniris.Seed, ticker *time.Ticker) (abortChan chan error, err error) {
	for _, seed := range seeds {
		if err = s.repo.StoreSeedPeer(seed); err != nil {
			return
		}
	}

	discoveryChan := make(chan uniris.Peer)
	reachChan := make(chan uniris.Peer)
	unreachChan := make(chan uniris.Peer)
	abortChan = make(chan error, 1)
	errors := make(chan error)

	go func() {

		for range ticker.C {

			go s.startCycle(localPeer, seeds, discoveryChan, reachChan, unreachChan, errors)

			//Stop the ticker when an unexpected error is returned
			//We also close all the channels to prevent any unexpected storage
			go func() {
				for err := range errors {
					ticker.Stop()
					abortChan <- err
					return
				}
			}()
		}
	}()

	return
}

//StoreLocalPeer gets peer info and store it as owned peer
func (s Service) StoreLocalPeer(pbKey string, port int, ver string) (p uniris.Peer, err error) {
	lon, lat, ip, _, _, err := getPeerSystemInfo(s.pInfo)
	if err != nil {
		return
	}

	p = uniris.NewLocalPeer(pbKey, ip, port, ver, lon, lat)
	if err = s.repo.StoreKnownPeer(p); err != nil {
		return
	}

	return
}

func (s Service) startCycle(localP uniris.Peer, seeds []uniris.Seed, dChan chan<- uniris.Peer, rChan chan<- uniris.Peer, uChan chan<- uniris.Peer, eChan chan<- error) {
	if err := s.refreshLocalPeer(localP, seeds); err != nil {
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

	c := newCycle(localP, s.msg, rp, up)
	if err != nil {
		eChan <- err
		return
	}

	go s.handleDiscoveries(c, dChan, eChan)
	go s.handleUnreachables(c, uChan, eChan)
	go s.handleReachables(c, rChan, eChan)
	go c.run(localP, seeds, knownPeers)

	go func() {
		for err := range c.errChan {
			eChan <- err
			return
		}
	}()
}

func (s Service) refreshLocalPeer(p uniris.Peer, sp []uniris.Seed) error {
	kp, err := s.repo.ListKnownPeers()
	if err != nil {
		return err
	}
	status, err := getPeerStatus(p, getSeedDiscoveryAverage(sp, kp), s.pNet)
	if err != nil {
		return err
	}

	_, _, _, cpu, space, err := getPeerSystemInfo(s.pInfo)
	if err != nil {
		return err
	}

	p.Refresh(status, space, cpu, getP2PFactor(kp), len(kp))
	if err = s.repo.StoreKnownPeer(p); err != nil {
		return err
	}

	return nil
}

func (s Service) handleReachables(c cycle, rChan chan<- uniris.Peer, eChan chan<- error) {
	for p := range c.reachChan {
		if s.repo.ContainsUnreachablePeer(p.Identity().PublicKey()) {
			if err := s.repo.RemoveUnreachablePeer(p.Identity().PublicKey()); err != nil {
				eChan <- err
				return
			}
		}

		if err := s.notif.NotifyReachable(p); err != nil {
			eChan <- err
			return
		}
		rChan <- p
	}
}

func (s Service) handleUnreachables(c cycle, uChan chan<- uniris.Peer, eChan chan<- error) {
	for p := range c.unreachChan {
		if !s.repo.ContainsUnreachablePeer(p.Identity().PublicKey()) {
			if err := s.repo.StoreUnreachablePeer(p.Identity().PublicKey()); err != nil {
				eChan <- err
				return
			}
			if err := s.notif.NotifyUnreachable(p); err != nil {
				eChan <- err
				return
			}
		}
		uChan <- p
	}
}

func (s Service) handleDiscoveries(c cycle, dChan chan<- uniris.Peer, eChan chan<- error) {
	for p := range c.discoveryChan {

		if err := s.repo.StoreKnownPeer(p); err != nil {
			eChan <- err
			return
		}

		if err := s.notif.NotifyDiscovery(p); err != nil {
			eChan <- err
			return
		}

		dChan <- p
	}
}

//CompareSyncRequest compare the peers received from the SYN gossip request
// and returns the unknown peers and the peers the sender does not known
func (s Service) CompareSyncRequest(receivedPeers []uniris.Peer) ([]uniris.Peer, []uniris.Peer, error) {
	kp, err := s.repo.ListKnownPeers()
	if err != nil {
		return nil, nil, err
	}

	return getUnknownPeers(kp, receivedPeers), getNewPeers(kp, receivedPeers), nil
}

//StoreAcknowledgePeers store the receiving peers from the ACK gossip request
func (s Service) StoreAcknowledgePeers(receivedPeers []uniris.Peer) error {
	for _, p := range receivedPeers {
		if err := s.repo.StoreKnownPeer(p); err != nil {
			return err
		}
		if err := s.notif.NotifyDiscovery(p); err != nil {
			return err
		}
	}

	return nil
}

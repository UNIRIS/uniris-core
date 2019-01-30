package discovery

import (
	"time"
)

//Repository provides access to the discovery persistence
type Repository interface {
	StoreSeedPeer(PeerIdentity) error
	StoreUnreachablePeer(pubKey string) error
	StoreKnownPeer(Peer) error
	RemoveUnreachablePeer(pubKey string) error

	ListSeedPeers() ([]PeerIdentity, error)
	CountKnownPeers() (int, error)
	ListKnownPeers() ([]Peer, error)
	ListUnreachablePeers() ([]PeerIdentity, error)
	ListReachablePeers() ([]PeerIdentity, error)
	ContainsUnreachablePeer(peerPublicKey string) bool
}

//Service handle the gossip spread and discovery of the network
type Service struct {
	repo  Repository
	pInfo PeerInformer
	pNet  PeerNetworker
	cli   Client
	notif Notifier
}

//Notifier handle the notification of the gossip events
type Notifier interface {
	NotifyReachable(PeerIdentity) error
	NotifyUnreachable(PeerIdentity) error
	NotifyDiscovery(Peer) error
}

//NewService creates the gossip handler
func NewService(repo Repository, cli Client, notif Notifier, pNet PeerNetworker, pInfo PeerInformer) Service {
	return Service{
		repo:  repo,
		cli:   cli,
		notif: notif,
		pInfo: pInfo,
		pNet:  pNet,
	}
}

//BootStrapingMinTime is the necessary minimum time on seconds to finish learning about the network
const BootStrapingMinTime = 1800

//Gossip spreads local informations and updated the receiving
func (s Service) Gossip(localPeer Peer, seeds []PeerIdentity, ticker *time.Ticker) (abortChan chan error, err error) {
	for _, seed := range seeds {
		if err = s.repo.StoreSeedPeer(seed); err != nil {
			return
		}
	}

	discoveryChan := make(chan Peer)
	reachChan := make(chan PeerIdentity)
	unreachChan := make(chan PeerIdentity)
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
func (s Service) StoreLocalPeer(pbKey string, port int, ver string) (p Peer, err error) {
	lon, lat, ip, _, _, err := getPeerSystemInfo(s.pInfo)
	if err != nil {
		return
	}

	p = NewLocalPeer(pbKey, ip, port, ver, lon, lat)
	if err = s.repo.StoreKnownPeer(p); err != nil {
		return
	}

	return
}

func (s Service) startCycle(localP Peer, seeds []PeerIdentity, dChan chan<- Peer, rChan chan<- PeerIdentity, uChan chan<- PeerIdentity, eChan chan<- error) {
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

	c := newCycle(localP, s.cli, rp, up)
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

func (s Service) refreshLocalPeer(p Peer, sp []PeerIdentity) error {
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

func (s Service) handleReachables(c cycle, rChan chan<- PeerIdentity, eChan chan<- error) {
	for p := range c.reachChan {
		if s.repo.ContainsUnreachablePeer(p.PublicKey()) {
			if err := s.repo.RemoveUnreachablePeer(p.PublicKey()); err != nil {
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

func (s Service) handleUnreachables(c cycle, uChan chan<- PeerIdentity, eChan chan<- error) {
	for p := range c.unreachChan {
		if !s.repo.ContainsUnreachablePeer(p.PublicKey()) {
			if err := s.repo.StoreUnreachablePeer(p.PublicKey()); err != nil {
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

func (s Service) handleDiscoveries(c cycle, dChan chan<- Peer, eChan chan<- error) {
	for p := range c.discoveryChan {
		if err := s.storeDiscoveredPeer(p); err != nil {
			eChan <- err
			return
		}

		dChan <- p
	}
}

//ComparePeers compare the peers received from the SYN gossip request
// and returns the unknown peers and the peers the sender does not known
func (s Service) ComparePeers(receivedPeers []Peer) ([]Peer, []Peer, error) {
	kp, err := s.repo.ListKnownPeers()
	if err != nil {
		return nil, nil, err
	}

	unknownPeers := s.getUnknownPeers(kp, receivedPeers)
	newPeers := s.getNewPeers(kp, receivedPeers)

	return unknownPeers, newPeers, nil
}

func (s Service) getUnknownPeers(knownPeers []Peer, comparePP []Peer) []Peer {
	mapPeers := s.mapPeerSlice(knownPeers)

	diff := make([]Peer, 0)

	for _, p := range comparePP {

		//Checks if the compared peer is include inside the repository
		kp, exist := mapPeers[p.Identity().PublicKey()]

		if !exist {
			//Adds to the list if the peer is unknown
			diff = append(diff, p)
		} else if p.HeartbeatState().MoreRecentThan(kp.HeartbeatState()) {
			//Adds to the list if the peer is more recent
			diff = append(diff, p)
		}
	}

	return diff
}

func (s Service) getNewPeers(knownPeers []Peer, comparePP []Peer) []Peer {
	mapComparee := s.mapPeerSlice(comparePP)

	diff := make([]Peer, 0)

	for _, p := range knownPeers {

		//Checks if the known peer is include inside the list of compared peer
		c, exist := mapComparee[p.Identity().PublicKey()]

		if !exist {
			//Adds to the list if the peer is unknown
			diff = append(diff, p)
		} else if p.HeartbeatState().MoreRecentThan(c.HeartbeatState()) {
			//Adds to the list if the peer is more recent
			diff = append(diff, p)
		}
	}

	return diff
}

func (s Service) mapPeerSlice(pp []Peer) map[string]Peer {
	mPeers := make(map[string]Peer)
	for _, p := range pp {
		mPeers[p.Identity().PublicKey()] = p
	}
	return mPeers
}

//AcknowledgeNewPeers store the receiving peers from the ACK gossip request
func (s Service) AcknowledgeNewPeers(peers []Peer) error {
	for _, p := range peers {
		if err := s.storeDiscoveredPeer(p); err != nil {
			return err
		}
	}

	return nil
}

func (s Service) storeDiscoveredPeer(p Peer) error {
	if err := s.repo.StoreKnownPeer(p); err != nil {
		return err
	}
	if err := s.notif.NotifyDiscovery(p); err != nil {
		return err
	}

	return nil
}

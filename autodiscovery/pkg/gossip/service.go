package gossip

import (
	"encoding/hex"
	"time"

	discovery "github.com/uniris/uniris-core/autodiscovery/pkg"
	"github.com/uniris/uniris-core/autodiscovery/pkg/monitoring"
)

//PeerDiff describes a diff to identify the unknown peers from the initiator or the receiver a SYN request is received
type PeerDiff struct {

	//UnknownLocally describes the peer the SYN request receiver does not know
	UnknownLocally []discovery.Peer

	//UnknownRemotly describes the peer the SYN request initiator does not know
	UnknownRemotly []discovery.Peer
}

//Service is the interface that provide gossip methods
type Service interface {
	Gossip(discovery.Peer) error
	ScheduleGossip(discovery.Peer) error
	DiffPeers([]discovery.Peer) (*PeerDiff, error)
	NotifyDiscovery(p discovery.Peer)
}

type service struct {
	msg   discovery.GossipCycleMessenger
	repo  discovery.Repository
	notif discovery.GossipRoundNotifier
	mon   monitoring.Service
}

//Run start the gossip with peers
func (s service) ScheduleGossip(init discovery.Peer) error {
	ticker := time.NewTicker(1 * time.Second)
	for range ticker.C {
		if err := s.Gossip(init); err != nil {
			return err
		}
	}
	return nil
}

//Gossip initialize the gossip session by running a gossip round
func (s service) Gossip(i discovery.Peer) error {
	sp, err := s.repo.ListSeedPeers()
	if err != nil {
		return err
	}

	//Refresh own peer before to gossip and send new information
	if err := s.mon.RefreshOwnedPeer(); err != nil {
		return err
	}

	kp, err := s.repo.ListKnownPeers()
	if err != nil {
		return err
	}

	r, err := discovery.NewGossipRound(i, kp, sp, s.msg)
	if err != nil {
		return err
	}
	newPeers, err := r.Run(i, kp)
	if err != nil {
		return err
	}
	for _, p := range newPeers {
		if err := s.repo.SetPeer(p); err != nil {
			return err
		}
		s.NotifyDiscovery(p)
	}
	return nil
}

//DiffPeers returns the diff between known peers and given list of peer
func (s service) DiffPeers(given []discovery.Peer) (*PeerDiff, error) {
	diff := new(PeerDiff)

	kp, err := s.repo.ListKnownPeers()
	if err != nil {
		return nil, err
	}

	//Get the peers that the SYN initiator request does not known
	gpMap := s.mapPeers(given)
	for _, p := range kp {
		if _, exist := gpMap[hex.EncodeToString(p.PublicKey())]; exist == false {
			diff.UnknownRemotly = append(diff.UnknownRemotly, p)
		}
	}

	//Gets the peers unknown locally
	knMap := s.mapPeers(kp)
	for _, p := range given {
		if _, exist := knMap[hex.EncodeToString(p.PublicKey())]; exist == false {
			diff.UnknownLocally = append(diff.UnknownLocally, p)
		}
	}

	return diff, nil
}

func (s service) NotifyDiscovery(p discovery.Peer) {
	s.notif.DisptachNewPeer(p)
}

func (s service) mapPeers(pp []discovery.Peer) map[string]discovery.Peer {
	mPeers := make(map[string]discovery.Peer, 0)
	for _, p := range pp {
		mPeers[hex.EncodeToString(p.PublicKey())] = p
	}
	return mPeers
}

//NewService creates a gossiping service its dependencies
func NewService(repo discovery.Repository, msg discovery.GossipCycleMessenger, notif discovery.GossipRoundNotifier, mon monitoring.Service) Service {
	return service{
		repo:  repo,
		msg:   msg,
		notif: notif,
		mon:   mon,
	}
}

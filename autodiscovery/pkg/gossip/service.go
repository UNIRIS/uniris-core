package gossip

import (
	"encoding/hex"
	"errors"
	"time"

	discovery "github.com/uniris/uniris-core/autodiscovery/pkg"
	"github.com/uniris/uniris-core/autodiscovery/pkg/monitoring"
)

//Notifier is the interface that provide methods to notify gossip discoveries
type Notifier interface {
	Notify(discovery.Peer)
}

//ErrPeerUnreachable is returned when the gossip cycle cannot reach a peer
var ErrPeerUnreachable = errors.New("Cannot reach the peer %s")

//PeerDiff describes a diff to identify the unknown peers from the initiator or the receiver a SYN request is received
type PeerDiff struct {
	//UnknownLocally describes the peer the SYN request receiver does not know
	UnknownLocally []discovery.Peer
	//UnknownRemotly describes the peer the SYN request initiator does not know
	UnknownRemotly []discovery.Peer
}

//Service is the interface that provide gossip methods
type Service interface {
	Spread(discovery.Peer) error
	ComparePeers([]discovery.Peer) (*PeerDiff, error)
}

type service struct {
	spr   discovery.GossipSpreader
	repo  discovery.Repository
	notif Notifier
	mon   monitoring.Service
}

//Spread creates gossip cycles and spreads the known peers across the network
func (s service) Spread(init discovery.Peer) error {
	seeds, err := s.repo.ListSeedPeers()
	if err != nil {
		return err
	}
	ticker := time.NewTicker(1 * time.Second)
	for range ticker.C {
		if err := s.runGossip(init, seeds); err != nil {
			return err
		}
	}
	return nil
}

func (s service) runGossip(init discovery.Peer, seeds []discovery.Seed) error {
	//Refresh own peer before to gossip and send new information
	if err := s.mon.RefreshOwnedPeer(); err != nil {
		return err
	}
	kp, err := s.repo.ListKnownPeers()
	if err != nil {
		return err
	}

	c, err := discovery.NewGossipCycle(init, kp, seeds)
	if err != nil {
		return err
	}
	for _, p := range c.SelectPeers() {
		r := c.CreateRound(p)
		if err := r.Spread(kp, s.spr); err != nil {
			//We do not throw an error when the peer is unreachable
			//Gossip must continue
			if err == ErrPeerUnreachable {
				return nil
			}
			return err
		}
	}

	for _, p := range c.Discoveries() {
		if err := s.repo.SetPeer(p); err != nil {
			return err
		}
		s.notif.Notify(p)
	}
	return nil
}

//ComparePeers returns the diff between known peers and given list of peer
func (s service) ComparePeers(given []discovery.Peer) (*PeerDiff, error) {
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

func (s service) mapPeers(pp []discovery.Peer) map[string]discovery.Peer {
	mPeers := make(map[string]discovery.Peer, 0)
	for _, p := range pp {
		mPeers[hex.EncodeToString(p.PublicKey())] = p
	}
	return mPeers
}

//NewService creates a gossiping service its dependencies
func NewService(repo discovery.Repository, spr discovery.GossipSpreader, notif Notifier, mon monitoring.Service) Service {
	return service{
		repo:  repo,
		spr:   spr,
		notif: notif,
		mon:   mon,
	}
}

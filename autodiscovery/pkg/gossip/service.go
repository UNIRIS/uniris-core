package gossip

import (
	"encoding/hex"

	discovery "github.com/uniris/uniris-core/autodiscovery/pkg"
	"github.com/uniris/uniris-core/autodiscovery/pkg/monitoring"
)

//Service is the interface that provide gossip methods
type Service interface {
	Spread(discovery.Peer) error
	DiffPeers([]discovery.Peer) (*PeerDiff, error)
	RunCycle(initiator discovery.Peer, receiver discovery.Peer, knownPeers []discovery.Peer) ([]discovery.Peer, error)
}

//Messenger is the interface that provides methods to send gossip requests
type Messenger interface {

	//Sends a SYN request
	SendSyn(SynRequest) (*SynAck, error)

	//Sends a ACK request after receipt of the SYN request
	SendAck(AckRequest) error
}

//Notifier is the interface that provides methods to notify gossip discovery
type Notifier interface {

	//Notify a new peer has been discovered
	Notify(peer discovery.Peer)
}

type service struct {
	msg   Messenger
	repo  discovery.Repository
	notif Notifier
	monit monitoring.Service
}

//PeerDiff describes a diff to identify the unknown peers from the initiator or the receiver a SYN request is received
type PeerDiff struct {

	//UnknownLocally describes the peer the SYN request receiver does not know
	UnknownLocally []discovery.Peer

	//UnknownRemotly describes the peer the SYN request initiator does not know
	UnknownRemotly []discovery.Peer
}

//DiffPeers returns the diff between known peers and given list of peer
func (s service) DiffPeers(given []discovery.Peer) (*PeerDiff, error) {
	kp, err := s.repo.ListKnownPeers()
	if err != nil {
		return nil, err
	}

	diff := new(PeerDiff)

	//Get the peers that the SYN initiator request does not known
	gpMap := s.mapPeers(given)
	for _, p := range kp {
		if _, exist := gpMap[hex.EncodeToString(p.Identity().PublicKey())]; exist == false {
			diff.UnknownRemotly = append(diff.UnknownRemotly, p)
		}
	}

	//Gets the peers unknown locally
	knMap := s.mapPeers(kp)
	for _, p := range given {
		if _, exist := knMap[hex.EncodeToString(p.Identity().PublicKey())]; exist == false {
			diff.UnknownLocally = append(diff.UnknownLocally, p)
		}
	}

	return diff, nil
}

//Spread creates a gossip round, stores and notifies the discovered peers
func (s service) Spread(init discovery.Peer) error {
	sp, err := s.repo.ListSeedPeers()
	if err != nil {
		return err
	}

	kp, err := s.repo.ListKnownPeers()
	if err != nil {
		return err
	}

	r, err := discovery.NewGossipRound(init, kp, sp)
	if err != nil {
		return err
	}

	pSelected, err := r.SelectPeers()
	if err != nil {
		return err
	}
	for _, p := range pSelected {
		newPeers, err := s.RunCycle(init, p, kp)
		if err != nil {
			return err
		}
		for _, p := range newPeers {
			if err := s.repo.AddPeer(p); err != nil {
				return err
			}
			s.notif.Notify(p)
		}
	}
	return nil
}

func (s service) RunCycle(init discovery.Peer, recpt discovery.Peer, kp []discovery.Peer) ([]discovery.Peer, error) {

	owned, err := s.repo.GetOwnedPeer()
	if err != nil {
		return nil, err
	}

	//Refreshes owned peer state before sending any requests
	if err := s.monit.RefreshPeer(owned); err != nil {
		return nil, err
	}

	synAck, err := s.msg.SendSyn(NewSynRequest(init, recpt, kp))
	if err != nil {
		return nil, err
	}
	if len(synAck.UnknownPeers) > 0 {
		reqPeers := make([]discovery.Peer, 0)
		mapPeers := s.mapPeers(kp)
		for _, p := range synAck.UnknownPeers {
			if k, exist := mapPeers[hex.EncodeToString(p.Identity().PublicKey())]; exist {
				reqPeers = append(reqPeers, k)
			}
		}

		if err := s.msg.SendAck(NewAckRequest(init, recpt, reqPeers)); err != nil {
			return nil, err
		}
		return synAck.NewPeers, nil
	}
	return []discovery.Peer{}, nil
}

func (s service) mapPeers(pp []discovery.Peer) map[string]discovery.Peer {
	mPeers := make(map[string]discovery.Peer, 0)
	for _, p := range pp {
		mPeers[hex.EncodeToString(p.Identity().PublicKey())] = p
	}
	return mPeers
}

//NewService creates a gossiping service its dependencies
func NewService(repo discovery.Repository, msg Messenger, notif Notifier, monit monitoring.Service) Service {
	return service{
		repo:  repo,
		msg:   msg,
		notif: notif,
		monit: monit,
	}
}

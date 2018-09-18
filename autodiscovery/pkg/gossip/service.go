package gossip

import (
	"encoding/hex"

	discovery "github.com/uniris/uniris-core/autodiscovery/pkg"
)

//Service is the interface that provide gossip methods
type Service interface {
	Spread(discovery.Peer) error
	DiffPeers([]discovery.Peer) (*PeerDiff, error)
}

//Messenger is the interface that provide methods to send gossip requests
type Messenger interface {

	//Sends a SYN request
	SendSyn(SynRequest) (SynAck, error)

	//Sends a ACK request after receipt of the SYN request
	SendAck(AckRequest) error
}

//Notifier provides the gossip discovery notification
type Notifier interface {

	//Notify a new peer has been discovered
	Notify(peer discovery.Peer)
}

type service struct {
	msg   Messenger
	repo  discovery.Repository
	notif Notifier
}

type PeerDiff struct {
	UnknownLocally []discovery.Peer
	UnknownRemotly []discovery.Peer
}

//DiffPeers returns the diff between known peers and given list of peer
func (s service) DiffPeers(gp []discovery.Peer) (*PeerDiff, error) {
	kp, err := s.repo.ListKnownPeers()
	if err != nil {
		return nil, err
	}

	diff := new(PeerDiff)

	//Get the peers the given list does not know
	gpMap := s.mapPeers(gp)
	for _, p := range kp {
		if _, exist := gpMap[hex.EncodeToString(p.PublicKey())]; !exist {
			diff.UnknownRemotly = append(diff.UnknownRemotly, p)
		}
	}

	//Gets the peers that we don't known
	knMap := s.mapPeers(kp)
	for _, p := range gp {
		if _, exist := knMap[hex.EncodeToString(p.PublicKey())]; exist == false {
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

	for _, p := range r.SelectPeers() {
		newPeers, err := s.dial(init, p, kp)
		if err != nil {
			return nil
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

func (s service) dial(init discovery.Peer, recpt discovery.Peer, kp []discovery.Peer) ([]discovery.Peer, error) {
	synAck, err := s.msg.SendSyn(NewSynRequest(init, recpt, kp))
	if err != nil {
		return nil, err
	}
	if len(synAck.UnknownPeers) > 0 {
		reqPeers := make([]discovery.Peer, 0)
		mapPeers := s.mapPeers(kp)
		for _, p := range synAck.UnknownPeers {
			if k, exist := mapPeers[hex.EncodeToString(p.PublicKey())]; exist {
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
		mPeers[hex.EncodeToString(p.PublicKey())] = p
	}
	return mPeers
}

//NewService creates a gossiping service its dependencies
func NewService(repo discovery.Repository, msg Messenger, notif Notifier) Service {
	return service{
		repo:  repo,
		msg:   msg,
		notif: notif,
	}
}

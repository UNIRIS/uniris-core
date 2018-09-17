package gossip

import (
	"errors"
	"log"

	"github.com/uniris/uniris-core/autodiscovery/pkg/discovery"
)

type Repository interface {
	ListSeedPeers() ([]discovery.SeedPeer, error)
	ListKnownPeers() ([]discovery.Peer, error)
	AddPeer(discovery.Peer) error
}

//GossipNotifier wraps the gossip discovery notification
type GossipNotifier interface {
	Notify(peer discovery.Peer) error
}

type GossipMessenger interface {
	SendSyn(SynRequest) (SynAck, error)
	SendAck(AckRequest) error
}

type Service interface {
	Run(initiator discovery.Peer) error
}

type service struct {
	repo  Repository
	msg   GossipMessenger
	notif GossipNotifier
}

func NewService(repo Repository, msg GossipMessenger, notif GossipNotifier) Service {
	return service{
		repo:  repo,
		msg:   msg,
		notif: notif,
	}
}

func (s service) Run(initiator discovery.Peer) error {
	seeds, err := s.repo.ListSeedPeers()
	if err != nil {
		return err
	}
	if len(seeds) == 0 {
		return errors.New("Cannot start a gossip round without a list of peers")
	}

	knownPeers, err := s.repo.ListKnownPeers()
	if err != nil {
		return err
	}

	round := GossipRound{
		Initiator:  initiator,
		KnownPeers: knownPeers,
		Seeds:      seeds,
	}
	newPeers, err := s.gossip(round)
	if err != nil {
		return err
	}
	for _, p := range newPeers {
		if err := s.repo.AddPeer(p); err != nil {
			return err
		}
		//Notifies a peer discovery subscriber
		if err := s.notif.Notify(p); err != nil {
			return err
		}
	}
	return nil
}

func (s service) gossip(round GossipRound) ([]discovery.Peer, error) {
	newPeers := make([]discovery.Peer, 0)
	for _, p := range round.SelectPeers() {
		log.Printf("Gossiping with peer %s...", p.Endpoint())
		synAck, err := s.msg.SendSyn(NewSynRequest(round.Initiator, p, round.KnownPeers))
		if err != nil {
			return nil, err
		}

		if len(synAck.UnknownPeers) > 0 {
			req := round.RequestedPeers(synAck.UnknownPeers)
			if err := s.msg.SendAck(NewAckRequest(round.Initiator, p, req)); err != nil {
				return nil, err
			}
		}

		//Sets the new peers to insert or to update
		for _, p := range synAck.NewPeers {
			newPeers = append(newPeers, p)
		}
	}
	return newPeers, nil
}

package usecases

import (
	"errors"
	"log"

	"github.com/uniris/uniris-core/autodiscovery/core/domain"

	"github.com/uniris/uniris-core/autodiscovery/core/ports"
)

//StartGossipRound executes a gossip round
func StartGossipRound(repo ports.PeerRepository, broker ports.GossipBroker, notifier ports.DiscoveryNotifier) error {
	seeds, err := repo.ListSeeds()
	if err != nil {
		return err
	}
	if len(seeds) == 0 {
		return errors.New("Cannot start a gossip round without a list of peers")
	}

	knownPeers, err := repo.ListPeers()
	if err != nil {
		return err
	}

	ownedPeer, err := repo.GetOwnedPeer()
	if err != nil {
		return err
	}

	round := domain.NewGossipRound(seeds, knownPeers)

	for _, receiver := range round.SelectPeers() {
		log.Printf("Gossiping with peer %s...", receiver.GetDiscoveryEndpoint())
		newPeers, err := startGossip(gossip{
			Initiator:  ownedPeer,
			Receiver:   receiver,
			KnownPeers: knownPeers,
			Broker:     broker,
			Round:      round,
		})
		if err != nil {
			return err
		}
		if len(newPeers) > 0 {
			log.Printf("%d unknown peers discovered from %s", len(newPeers), receiver.GetDiscoveryEndpoint())
			for _, newPeer := range newPeers {
				exist, err := repo.ContainsPeer(newPeer)
				if err != nil {
					return err
				}
				if exist {
					if err := repo.UpdatePeer(newPeer); err != nil {
						return err
					}
				} else {
					if err := repo.InsertPeer(newPeer); err != nil {
						return err
					}
				}
				notifier.NotifyNewPeer(newPeer)
			}
		}
	}
	return nil
}

type gossip struct {
	Initiator  domain.Peer
	Receiver   domain.Peer
	KnownPeers []domain.Peer
	Broker     ports.GossipBroker
	Round      domain.GossipRound
}

func startGossip(g gossip) ([]domain.Peer, error) {
	newPeers := make([]domain.Peer, 0)
	synAck, err := g.Broker.SendSyn(domain.NewSynRequest(g.Initiator, g.Receiver, g.KnownPeers))
	log.Printf("%d", len(synAck.NewPeers))
	if err != nil {
		return nil, err
	}

	if len(synAck.UnknownPeers) > 0 {
		requestedPeers := g.Round.GetRequestedPeers(synAck.UnknownPeers)
		if err := g.Broker.SendAck(domain.NewAckRequest(g.Initiator, g.Receiver, requestedPeers)); err != nil {
			return nil, err
		}
	}

	//Sets the new peers to insert or to update
	for _, newPeer := range synAck.NewPeers {
		newPeers = append(newPeers, newPeer)
	}

	return newPeers, nil
}

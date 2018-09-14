package usecases

import (
	"errors"
	"math/rand"

	"github.com/uniris/uniris-core/autodiscovery/domain"
	"github.com/uniris/uniris-core/autodiscovery/usecases/ports"
	"github.com/uniris/uniris-core/autodiscovery/usecases/repositories"
)

//StartGossipRound initiates a gossip round
func StartGossipRound(repo repositories.PeerRepository, seedReader ports.SeedReader, messenger ports.GossipMessenger) error {
	seeds, err := seedReader.GetSeeds()
	if err != nil {
		return err
	}

	if len(seeds) == 0 {
		return errors.New("Cannot gossip without seed peers")
	}

	gossipTargets := make([]domain.Peer, 0)
	gossipTargets = append(gossipTargets, getRandomPeer(seeds))

	knownPeers, err := repo.ListPeers()
	if err != nil {
		return err
	}

	filterKnownPeers := excludeOwnPeer(knownPeers)
	if len(filterKnownPeers) > 0 {
		gossipTargets = append(gossipTargets, getRandomPeer(filterKnownPeers))
	}

	if err := Gossip(repo, messenger, gossipTargets); err != nil {
		return err
	}

	return nil
}

func excludeOwnPeer(peers []domain.Peer) []domain.Peer {
	noExclude := make([]domain.Peer, 0)
	for _, peer := range peers {
		if !peer.IsOwned {
			noExclude = append(noExclude, peer)
		}
	}
	return noExclude
}

func getRandomPeer(peers []domain.Peer) domain.Peer {
	if len(peers) > 1 {
		rnd := rand.Intn(len(peers) - 1)
		return peers[rnd]
	}
	return peers[0]
}

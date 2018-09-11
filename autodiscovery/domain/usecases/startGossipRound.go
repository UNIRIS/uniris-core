package usecases

import (
	"github.com/uniris/uniris-core/autodiscovery/domain/entities"
	"github.com/uniris/uniris-core/autodiscovery/domain/repositories"
	"github.com/uniris/uniris-core/autodiscovery/domain/services"
)

//StartGossipRound initiates the gossip round by calling the handshakes request to discover the network
func StartGossipRound(peerRepo repositories.PeerRepository, gossipService services.GossipService) error {
	seedPeers, err := peerRepo.ListSeedPeers()
	if err != nil {
		return err
	}

	discoveredPeers, err := peerRepo.ListDiscoveredPeers()
	if err != nil {
		return err
	}

	peersToCall := make([]*entities.Peer, 0)

	peersToCall = append(peersToCall, SelectRandomPeer(seedPeers))
	if len(discoveredPeers) > 0 {
		peersToCall = append(peersToCall, SelectRandomPeer(discoveredPeers))
	}

	knownPeers, err := peerRepo.ListDiscoveredPeers()
	if err != nil {
		return err
	}

	for _, peer := range peersToCall {
		ack, err := gossipService.Synchronize(peer, knownPeers)
		if err != nil {
			return err
		}

		peer.Category = entities.DiscoveredCategory
		err = SetNewPeers(peerRepo, ack.UnknownInitiatorPeers)
		if err != nil {
			return err
		}
	}

	return nil
}

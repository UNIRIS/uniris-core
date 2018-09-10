package usecases

import (
	"github.com/uniris/uniris-core/autodiscovery/domain/repositories"
	"github.com/uniris/uniris-core/autodiscovery/domain/services"
)

//DiscoverPeers initiates the gossip process by communicating known peers
func DiscoverPeers(peerRepo repositories.PeerRepository, gossip services.GossipService) error {
	selectedPeer, err := SelectRandomPeer(peerRepo)
	if err != nil {
		return err
	}

	knownPeers, err := ListKnownPeers(peerRepo)
	if err != nil {
		return err
	}

	newPeers, err := gossip.DiscoverPeers(*selectedPeer, knownPeers)
	if err != nil {
		return err
	}
	if err := SetNewPeers(peerRepo, newPeers); err != nil {
		return err
	}
	return nil
}

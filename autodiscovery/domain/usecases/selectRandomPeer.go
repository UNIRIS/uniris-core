package usecases

import (
	"fmt"
	"math/rand"

	"github.com/uniris/uniris-core/autodiscovery/domain/entities"
	"github.com/uniris/uniris-core/autodiscovery/domain/repositories"
)

//SelectRandomPeer pick a known random peer to initiate the gossip protocol
func SelectRandomPeer(peerRepo repositories.PeerRepository) (*entities.Peer, error) {
	knownPeers, err := ListKnownPeers(peerRepo)
	if err != nil {
		return nil, err
	}

	if len(knownPeers) == 0 {
		return nil, fmt.Errorf("Random peer selection cannot be done with no known peers")
	}

	if len(knownPeers) > 1 {
		rnd := rand.Intn(len(knownPeers) - 1)
		return knownPeers[rnd], nil
	}
	return knownPeers[0], nil
}

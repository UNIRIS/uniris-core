package usecases

import (
	"github.com/uniris/uniris-core/autodiscovery/domain"
	"github.com/uniris/uniris-core/autodiscovery/usecases/repositories"
)

//ProvideNewPeers retrieves the peers that a sender peer does not known
func ProvideNewPeers(repo repositories.PeerRepository, receivedPeers []domain.Peer) ([]domain.Peer, error) {
	newPeers := make([]domain.Peer, 0)
	if len(receivedPeers) == 0 {
		return newPeers, nil
	}

	knownPeers, err := repo.ListPeers()
	if err != nil {
		return nil, err
	}

	for _, peer := range knownPeers {
		if exist, _ := ContainsPeer(receivedPeers, peer); exist == false {
			newPeers = append(newPeers, peer)
		}
	}
	return newPeers, nil
}

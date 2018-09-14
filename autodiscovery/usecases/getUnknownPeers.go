package usecases

import (
	"github.com/uniris/uniris-core/autodiscovery/domain"
	"github.com/uniris/uniris-core/autodiscovery/usecases/repositories"
)

//GetUnknownPeers retrieves the peers that a receiver peer does not known
func GetUnknownPeers(repo repositories.PeerRepository, receivedPeers []domain.Peer) ([]domain.Peer, error) {
	unknownPeers := make([]domain.Peer, 0)
	if len(receivedPeers) == 0 {
		return unknownPeers, nil
	}

	knownPeers, err := repo.ListPeers()
	if err != nil {
		return nil, err
	}

	for _, peer := range receivedPeers {
		if exist, _ := ContainsPeer(knownPeers, peer); exist == false {
			unknownPeers = append(unknownPeers, peer)
		}
	}
	return unknownPeers, nil
}

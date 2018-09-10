package usecases

import (
	"github.com/uniris/uniris-core/autodiscovery/domain/entities"
	"github.com/uniris/uniris-core/autodiscovery/domain/repositories"
)

//ListKnownPeers returns list of known peers
func ListKnownPeers(peerRepo repositories.PeerRepository) ([]*entities.Peer, error) {
	mapPeer, err := peerRepo.GetPeers()
	if err != nil {
		return nil, err
	}

	//Convert map to slices
	peers := make([]*entities.Peer, 0)
	for _, peer := range mapPeer {
		peers = append(peers, peer)
	}
	return peers, nil
}

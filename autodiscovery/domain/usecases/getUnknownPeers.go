package usecases

import (
	"github.com/uniris/uniris-core/autodiscovery/domain/entities"
	"github.com/uniris/uniris-core/autodiscovery/domain/repositories"
)

//GetUnknownPeers compare known peers and new peers and returns the diff
func GetUnknownPeers(peerRepo repositories.PeerRepository, peersToCompare []*entities.Peer) ([]*entities.Peer, error) {
	knownPeers, err := ListKnownPeers(peerRepo)
	if err != nil {
		return nil, err
	}

	mapPeers := make(map[string]entities.Peer)
	for _, peer := range peersToCompare {
		mapPeers[string(peer.PublicKey)] = *peer
	}

	unknownPeers := make([]*entities.Peer, 0)
	for _, knownPeer := range knownPeers {
		if _, exist := mapPeers[string(knownPeer.PublicKey)]; exist == false {
			unknownPeers = append(unknownPeers, knownPeer)
		}
	}

	return unknownPeers, nil
}

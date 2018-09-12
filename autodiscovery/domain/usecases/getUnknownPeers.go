package usecases

import (
	"github.com/uniris/uniris-core/autodiscovery/domain/entities"
)

//GetUnknownPeers compare known peers and new peers and returns the diff
func GetUnknownPeers(knownPeers []*entities.Peer, peersToCompare []*entities.Peer) (unknownPeers []*entities.Peer) {
	if len(peersToCompare) == 0 {
		return knownPeers
	}

	mappedPeersForComparison := MapPeers(peersToCompare)
	for _, knownPeer := range knownPeers {
		if exist, _ := IsMapContainsPeer(mappedPeersForComparison, knownPeer); exist == false {
			unknownPeers = append(unknownPeers, knownPeer)
		}
	}

	return unknownPeers
}

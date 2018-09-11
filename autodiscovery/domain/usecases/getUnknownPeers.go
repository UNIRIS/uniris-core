package usecases

import (
	"encoding/hex"

	"github.com/uniris/uniris-core/autodiscovery/domain/entities"
)

//GetUnknownPeers compare known peers and new peers and returns the diff
func GetUnknownPeers(knownPeers []*entities.Peer, peersToCompare []*entities.Peer, mySelf *entities.Peer) []*entities.Peer {
	if len(peersToCompare) == 0 {
		return knownPeers
	}

	//Map the the peers to compare
	mapPeers := make(map[string]*entities.Peer)
	for _, peer := range peersToCompare {
		mapPeers[hex.EncodeToString(peer.PublicKey)] = peer
		mapPeers[hex.EncodeToString(mySelf.PublicKey)] = mySelf
	}

	unknownPeers := make([]*entities.Peer, 0)
	for _, knownPeer := range knownPeers {
		if _, exist := mapPeers[hex.EncodeToString(knownPeer.PublicKey)]; exist == false {
			unknownPeers = append(unknownPeers, knownPeer)
		}
	}

	return unknownPeers
}

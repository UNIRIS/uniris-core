package usecases

import (
	"encoding/hex"

	"github.com/uniris/uniris-core/autodiscovery/domain/entities"
)

//MapPeers converts a list of peers in map
func MapPeers(peers []*entities.Peer) map[string]*entities.Peer {
	mappedPeers := make(map[string]*entities.Peer, 0)
	for _, peer := range peers {
		mappedPeers[hex.EncodeToString(peer.PublicKey)] = peer
	}
	return mappedPeers
}

//IsMapContainsPeer checks if a map of peers contains a specific peer identified by public key
func IsMapContainsPeer(peers map[string]*entities.Peer, peer *entities.Peer) (bool, *entities.Peer) {
	knownPeer, exist := peers[hex.EncodeToString(peer.PublicKey)]
	return exist, knownPeer
}

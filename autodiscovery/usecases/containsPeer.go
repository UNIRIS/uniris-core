package usecases

import (
	"encoding/hex"

	"github.com/uniris/uniris-core/autodiscovery/domain"
)

//ContainsPeer checks if a peer exists in the repository
func ContainsPeer(source []domain.Peer, peer domain.Peer) (exist bool, knownPeer domain.Peer) {
	cMapped := make(map[string]domain.Peer, 0)
	for _, peer := range source {
		cMapped[hex.EncodeToString(peer.PublicKey)] = peer
	}
	knownPeer, exist = cMapped[hex.EncodeToString(peer.PublicKey)]
	return exist, knownPeer
}

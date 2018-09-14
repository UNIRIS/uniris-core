package usecases

import (
	"encoding/hex"

	"github.com/uniris/uniris-core/autodiscovery/core/domain"
	"github.com/uniris/uniris-core/autodiscovery/core/ports"
)

//GetUnknownPeers retrieves the peers that a receiver peer does not known
func GetUnknownPeers(repo ports.PeerRepository, receivedPeers []domain.Peer) ([]domain.Peer, error) {
	unknownPeers := make([]domain.Peer, 0)
	if len(receivedPeers) == 0 {
		return unknownPeers, nil
	}

	knownPeers, err := repo.ListPeers()
	if err != nil {
		return nil, err
	}

	for _, peer := range receivedPeers {
		if exist, _ := containsPeer(knownPeers, peer); exist == false {
			unknownPeers = append(unknownPeers, peer)
		}
	}
	return unknownPeers, nil
}

//ProvideNewPeers retrieves the peers that a sender peer does not known
func ProvideNewPeers(repo ports.PeerRepository, receivedPeers []domain.Peer) ([]domain.Peer, error) {
	newPeers := make([]domain.Peer, 0)
	if len(receivedPeers) == 0 {
		return newPeers, nil
	}

	knownPeers, err := repo.ListPeers()
	if err != nil {
		return nil, err
	}

	for _, peer := range knownPeers {
		if exist, _ := containsPeer(receivedPeers, peer); exist == false {
			newPeers = append(newPeers, peer)
		}
	}
	return newPeers, nil
}

//ContainsPeer checks if a peer exists from a source of peer
func containsPeer(source []domain.Peer, peer domain.Peer) (exist bool, knownPeer domain.Peer) {
	cMapped := make(map[string]domain.Peer, 0)
	for _, peer := range source {
		cMapped[hex.EncodeToString(peer.PublicKey)] = peer
	}
	knownPeer, exist = cMapped[hex.EncodeToString(peer.PublicKey)]
	return exist, knownPeer
}

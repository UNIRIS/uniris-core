package usecases

import (
	"github.com/uniris/uniris-core/autodiscovery/domain"
	"github.com/uniris/uniris-core/autodiscovery/usecases/repositories"
)

//StorePeer insert or update the peer in the repository
func StorePeer(repo repositories.PeerRepository, peer domain.Peer) error {
	knownPeers, err := repo.ListPeers()
	if err != nil {
		return err
	}
	isExist, existing := ContainsPeer(knownPeers, peer)
	if isExist {
		if peer.GetElapsedHeartbeats() < existing.GetElapsedHeartbeats() {
			existing.Refresh(peer.IP, peer.Port, peer.GenerationTime, peer.State)
			if err := repo.UpdatePeer(existing); err != nil {
				return err
			}
		}
	} else {
		if err := repo.InsertPeer(peer); err != nil {
			return err
		}
	}
	return nil
}

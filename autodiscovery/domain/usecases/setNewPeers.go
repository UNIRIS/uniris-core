package usecases

import (
	"fmt"
	"log"

	"github.com/uniris/uniris-core/autodiscovery/domain/entities"
	"github.com/uniris/uniris-core/autodiscovery/domain/repositories"
)

//SetNewPeers handles insert or update of new peers
func SetNewPeers(peerRepo repositories.PeerRepository, newPeers []*entities.Peer) error {
	peers, err := peerRepo.ListPeers()
	if err != nil {
		return err
	}

	peersToInsert, peersToUpdate := getNewOrUpdatePeers(peers, newPeers)

	if len(peersToInsert) > 0 {
		if err = insertNewPeers(peersToInsert, peerRepo); err != nil {
			return err
		}
	}

	if len(peersToUpdate) > 0 {
		if err = updateExistingPeers(peersToUpdate, peerRepo); err != nil {
			return err
		}
	}

	return nil
}

//Define if the new peer must be inserted ou updated
func getNewOrUpdatePeers(peers []*entities.Peer, newPeers []*entities.Peer) (peerToInsert []*entities.Peer, peerToUpdate []*entities.Peer) {
	peersMapped := MapPeers(peers)

	for _, newPeer := range newPeers {
		exist, knownPeer := IsMapContainsPeer(peersMapped, newPeer)
		if exist {
			if newPeer.GetElapsedHeartbeats() < knownPeer.GetElapsedHeartbeats() {
				peerToUpdate = append(peerToUpdate, newPeer)
			}
		} else {
			peerToInsert = append(peerToInsert, newPeer)
		}
	}

	return peerToInsert, peerToUpdate
}

//Add on the repository the new peers
func insertNewPeers(peers []*entities.Peer, repo repositories.PeerRepository) error {
	for _, peer := range peers {
		if err := repo.AddPeer(peer); err != nil {
			return err
		}
	}

	peers, err := repo.ListPeers()
	if err != nil {
		return err
	}
	for _, peer := range peers {
		fmt.Println("%s:%d", peer.IP.String(), peer.Port)
	}
	log.Printf("Peers in repositories: %d", len(peers))
	return nil
}

//Update the existing peers
func updateExistingPeers(peers []*entities.Peer, repo repositories.PeerRepository) error {
	for _, peer := range peers {
		if err := repo.UpdatePeer(peer); err != nil {
			return err
		}
	}
	return nil
}

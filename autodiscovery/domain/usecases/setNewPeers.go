package usecases

import (
	"encoding/hex"

	"github.com/uniris/uniris-core/autodiscovery/domain/entities"
	"github.com/uniris/uniris-core/autodiscovery/domain/repositories"
)

//SetNewPeers handles insert or update of new peers
func SetNewPeers(peerRepo repositories.PeerRepository, newPeers []*entities.Peer) error {
	peers, err := peerRepo.ListPeers()
	if err != nil {
		return err
	}

	//Map the the peers
	mapPeers := make(map[string]*entities.Peer)
	for _, peer := range peers {
		mapPeers[hex.EncodeToString(peer.PublicKey)] = peer
	}

	//Fetch the peer to update or to add
	peerToUpdate := []*entities.Peer{}
	peerToInsert := []*entities.Peer{}
	for _, newPeer := range newPeers {
		knownPeer, exist := mapPeers[hex.EncodeToString(newPeer.PublicKey)]
		if exist {
			if newPeer.GetElapsedHeartbeats() < knownPeer.GetElapsedHeartbeats() {
				peerToUpdate = append(peerToUpdate, newPeer)
			}
		} else {
			peerToInsert = append(peerToInsert, newPeer)
		}
	}

	//Add on the repository the new peers
	for _, peer := range peerToInsert {
		if err := peerRepo.AddPeer(peer); err != nil {
			return err
		}
	}

	//Update the existing peers
	for _, peer := range peerToUpdate {
		if err := peerRepo.UpdatePeer(peer); err != nil {
			return err
		}
	}

	return nil
}

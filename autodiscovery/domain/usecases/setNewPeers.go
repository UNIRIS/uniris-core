package usecases

import (
	"github.com/uniris/uniris-core/autodiscovery/domain/entities"
	"github.com/uniris/uniris-core/autodiscovery/domain/repositories"
)

//SetNewPeers handles insert or update of new peers
func SetNewPeers(peerRepo repositories.PeerRepository, newPeers []*entities.Peer) error {

	//Retrieve the stored peers
	knownPeers, err := peerRepo.GetPeers()
	if err != nil {
		return err
	}

	//Fetch the peer to update or to add
	peerToUpdate := []*entities.Peer{}
	peerToInsert := []*entities.Peer{}
	for _, newPeer := range newPeers {
		knownPeer, exist := knownPeers[string(newPeer.PublicKey)]
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

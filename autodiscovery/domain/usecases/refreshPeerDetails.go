package usecases

import (
	"github.com/uniris/uniris-core/autodiscovery/domain/entities"
)

//RefreshPeerDetails refreshs the peer details
func RefreshPeerDetails(peer *entities.Peer) error {

	peer.UpdateElapsedHeartbeats()

	return nil
	//TODO: define others details properties
}

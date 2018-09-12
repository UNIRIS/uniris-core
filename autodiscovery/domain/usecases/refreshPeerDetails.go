package usecases

import "github.com/uniris/uniris-core/autodiscovery/domain/repositories"

//RefreshSelfPeer refreshs the peer details
func RefreshSelfPeer(repo repositories.PeerRepository) error {
	peer, err := repo.GetOwnPeer()
	if err != nil {
		return err
	}

	peer.UpdateElapsedHeartbeats()
	//TODO: define others details properties

	if err = repo.UpdatePeer(peer); err != nil {
		return err
	}

	return nil
}

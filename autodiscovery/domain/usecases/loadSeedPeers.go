package usecases

import (
	"github.com/uniris/uniris-core/autodiscovery/domain/services"

	"github.com/uniris/uniris-core/autodiscovery/domain/repositories"
)

//LoadSeedPeers read the seed listing file and load into the peer database
func LoadSeedPeers(seedLoader services.SeedLoader, peerRepo repositories.PeerRepository) error {

	seedPeerList, err := seedLoader.GetSeedPeers()
	if err != nil {
		return err
	}

	//Store the seed peers
	for _, peer := range seedPeerList {
		if err := peerRepo.AddPeer(peer); err != nil {
			return err
		}
	}
	return nil
}

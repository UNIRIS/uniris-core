package usecases

import (
	"fmt"
	"log"

	"github.com/uniris/uniris-core/autodiscovery/domain/entities"

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
		peer.Category = entities.SeedCategory
		if err := peerRepo.AddPeer(peer); err != nil {
			return err
		}
	}

	log.Println(fmt.Sprintf("%d peers are loaded as seed peers", len(seedPeerList)))

	return nil
}

package usecases

import (
	"github.com/uniris/uniris-core/autodiscovery/domain/entities"
	"github.com/uniris/uniris-core/autodiscovery/domain/repositories"
	"github.com/uniris/uniris-core/autodiscovery/domain/services"
)

//DiscoverPeers call peers to discover their known peers by sharing our known peers
func DiscoverPeers(peersToCall []*entities.Peer, peerRepo repositories.PeerRepository, gossipService services.GossipService) error {
	knownPeers, err := peerRepo.ListPeers()
	if err != nil {
		return err
	}
	for _, peer := range peersToCall {
		synReq := &entities.SynchronizationRequest{
			PeerReceiver:     peer,
			KnownSenderPeers: knownPeers,
		}

		ackRes, err := gossipService.Synchronize(synReq)
		if err != nil {
			return err
		}

		if err = SetNewPeers(peerRepo, ackRes.UnknownSenderPeers); err != nil {
			return err
		}

		//Send ACK2
	}
	return nil
}

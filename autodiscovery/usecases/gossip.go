package usecases

import (
	"errors"
	"log"

	"github.com/uniris/uniris-core/autodiscovery/domain"
	"github.com/uniris/uniris-core/autodiscovery/usecases/repositories"
)

//Gossip etablishes connection with the peers to reach and send them discovery requests
func Gossip(repo repositories.PeerRepository, conf GossipConfiguration, peersToReach []domain.Peer) error {
	if len(peersToReach) == 0 {
		return errors.New("Cannot gossip without peers to reach")
	}

	ownedPeer, err := repo.GetOwnedPeer()
	if err != nil {
		return err
	}

	if err := RefreshPeer(repo, &ownedPeer, conf.Geolocalizer); err != nil {
		return err
	}

	knownPeers, err := repo.ListPeers()
	if err != nil {
		return err
	}

	for _, peer := range peersToReach {
		log.Printf("Reaching peer %s...", peer.GetDiscoveryEndpoint())

		synReq := domain.NewSynRequest(ownedPeer, peer.IP, peer.Port, knownPeers)
		ackRes, err := conf.Messenger.SendSynchro(synReq)
		if err != nil {
			return err
		}

		if len(ackRes.NewPeers) > 0 {
			log.Printf("%d new peers retrieved from %s", len(ackRes.NewPeers), peer.GetDiscoveryEndpoint())
			for _, newPeer := range ackRes.NewPeers {
				if err := StorePeer(repo, newPeer); err != nil {
					return err
				}
				log.Printf("New peer stored: %s", newPeer.GetDiscoveryEndpoint())
			}

			knownPeers, err := repo.ListPeers()
			if err != nil {
				return err
			}

			detailedPeers := make([]domain.Peer, 0)
			for _, detailedPeer := range ackRes.DetailPeersRequested {
				if exist, existing := ContainsPeer(knownPeers, detailedPeer); exist && existing.IsDiscovered() {
					detailedPeers = append(detailedPeers, existing)
				}
			}

			if len(detailedPeers) > 0 {
				ackReq := domain.NewAckRequest(detailedPeers)
				if err := conf.Messenger.SendAcknowledgement(ackReq); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

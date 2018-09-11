package usecases

import (
	"encoding/hex"
	"math/rand"

	"github.com/uniris/uniris-core/autodiscovery/domain/entities"
	"github.com/uniris/uniris-core/autodiscovery/domain/repositories"
	"github.com/uniris/uniris-core/autodiscovery/domain/services"
)

//StartGossipRound initiates the gossip round by calling the handshakes request to discover the network
func StartGossipRound(seedLoader services.SeedLoader, peerRepo repositories.PeerRepository, gossipService services.GossipService) error {
	seedPeers, err := seedLoader.GetSeedPeers()
	if err != nil {
		return err
	}

	ownPeer, err := peerRepo.GetOwnPeer()
	if err != nil {
		return err
	}
	knownPeers, err := peerRepo.ListPeers()
	if err != nil {
		return err
	}

	peersToCall := make([]*entities.Peer, 0)
	peersToCall = append(peersToCall, selectRandomPeer(seedPeers, ownPeer))

	knownPeers = excludePeer(knownPeers, ownPeer)
	if len(knownPeers) > 0 {
		peersToCall = append(peersToCall, selectRandomPeer(knownPeers, ownPeer))
	}

	if err := DiscoverPeers(peersToCall, peerRepo, gossipService); err != nil {
		return err
	}

	return nil
}

func excludePeer(list []*entities.Peer, exclude *entities.Peer) []*entities.Peer {
	excluded := make([]*entities.Peer, 0)
	for _, peer := range list {
		if hex.EncodeToString(peer.PublicKey) != hex.EncodeToString(exclude.PublicKey) {
			excluded = append(excluded, peer)
		}
	}
	return excluded
}

//SelectRandomPeer picks a known random peer from a list of peers (exclude own peer)
func selectRandomPeer(peers []*entities.Peer, ownPeer *entities.Peer) *entities.Peer {
	if len(peers) > 1 {
		rnd := rand.Intn(len(peers) - 1)
		return peers[rnd]
	}
	return peers[0]
}

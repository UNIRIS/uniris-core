package services

import (
	"github.com/uniris/uniris-core/autodiscovery/domain/entities"
)

//GossipService represents the autodiscovery requests
type GossipService interface {
	DiscoverPeers(destPeer entities.Peer, knownPeers []*entities.Peer) ([]*entities.Peer, error)
}

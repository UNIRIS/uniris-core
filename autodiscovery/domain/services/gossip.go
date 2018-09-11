package services

import (
	"github.com/uniris/uniris-core/autodiscovery/domain/entities"
)

//GossipService represents the autodiscovery requests
type GossipService interface {
	Synchronize(destPeer *entities.Peer, knownPeers []*entities.Peer) (*entities.Acknowledge, error)
}

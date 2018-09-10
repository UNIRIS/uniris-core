package services

import "github.com/uniris/uniris-core/autodiscovery/domain/entities"

//SeedLoader define the seed loader interface to retrieve the peers for bootstraping
type SeedLoader interface {
	GetSeedPeers() ([]*entities.Peer, error)
}

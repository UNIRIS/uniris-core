package ports

import "github.com/uniris/uniris-core/autodiscovery/domain"

//SeedReader wraps the configured seeds loading
type SeedReader interface {
	GetSeeds() ([]domain.Peer, error)
}

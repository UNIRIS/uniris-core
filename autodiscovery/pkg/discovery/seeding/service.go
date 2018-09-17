package seeding

import "github.com/uniris/uniris-core/autodiscovery/pkg/discovery"

type SeedReader interface {
	GetSeeds() ([]discovery.SeedPeer, error)
}

type Service interface {
	LoadSeeds() error
}

type Repository interface {
	AddSeed(discovery.SeedPeer) error
}

type service struct {
	read SeedReader
	repo Repository
}

func NewService(read SeedReader, repo Repository) Service {
	return service{
		read: read,
		repo: repo,
	}
}

func (s service) LoadSeeds() error {

	seeds, err := s.read.GetSeeds()
	if err != nil {
		return err
	}

	for _, sd := range seeds {
		if err := s.repo.AddSeed(sd); err != nil {
			return err
		}
	}

	return nil
}

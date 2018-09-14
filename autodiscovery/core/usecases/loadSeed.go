package usecases

import (
	"errors"

	"github.com/uniris/uniris-core/autodiscovery/core/ports"
)

//LoadSeeds reads the seed configuration file and save them
func LoadSeeds(repo ports.PeerRepository, conf ports.ConfigurationReader) error {
	seeds, err := conf.GetSeeds()
	if err != nil {
		return err
	}

	if len(seeds) == 0 {
		return errors.New("Cannot load empty seed list")
	}

	for _, seed := range seeds {
		if err := repo.StoreSeed(seed); err != nil {
			return err
		}
	}
	return nil
}

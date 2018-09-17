package file

import (
	"encoding/json"
	"io/ioutil"

	"github.com/uniris/uniris-core/autodiscovery/pkg/discovery"
)

type SeedReader struct{}

func (r SeedReader) GetSeeds() ([]discovery.SeedPeer, error) {
	bytes, err := ioutil.ReadFile("conf/seeds.json")
	if err != nil {
		return nil, err
	}

	seeds := make([]discovery.SeedPeer, 0)
	json.Unmarshal(bytes, &seeds)
	return seeds, nil
}

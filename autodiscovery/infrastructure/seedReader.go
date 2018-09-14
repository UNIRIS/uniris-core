package infrastructure

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/uniris/uniris-core/autodiscovery/domain"
)

//SeedReader implements the seed loader interface to retrieve the peers for bootstraping
type SeedReader struct {
}

//GetSeeds loads the seed peers from the configuration file
func (s SeedReader) GetSeeds() ([]domain.Peer, error) {
	path, err := filepath.Abs("./seed.json")
	if err != nil {
		return nil, err
	}
	jsonFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	seedPeerList := make([]domain.Peer, 0)

	//Deserialize the seed peers
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(byteValue, &seedPeerList)
	return seedPeerList, nil
}

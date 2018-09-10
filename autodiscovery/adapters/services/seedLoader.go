package services

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/uniris/uniris-core/autodiscovery/domain/entities"
)

//SeedLoader implements the seed loader interface to retrieve the peers for bootstraping
type SeedLoader struct {
}

//GetSeedPeers loads the seed peers from the configuration file
func (s SeedLoader) GetSeedPeers() ([]*entities.Peer, error) {
	path, err := filepath.Abs("./seed.json")
	if err != nil {
		return nil, err
	}
	jsonFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	seedPeerList := make([]entities.Peer, 0)

	//Deserialize the seed peers
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(byteValue, &seedPeerList)
	return seedPeerList, nil
}

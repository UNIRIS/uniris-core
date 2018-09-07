package usecases

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/uniris/uniris-core/autodiscovery/domain/entities"
	"github.com/uniris/uniris-core/autodiscovery/domain/repositories"
)

//LoadSeedPeers read the seed listing file and load into the peer database
func LoadSeedPeers(peerRepo repositories.PeerRepository, seedURI string) error {
	jsonFile, err := os.Open(seedURI)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	seedPeerList := make([]entities.Peer, 0)

	//Deserialize the seed peers
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return err
	}

	json.Unmarshal(byteValue, &seedPeerList)

	//Store the seed peers
	for _, peer := range seedPeerList {
		if err := peerRepo.AddPeer(peer); err != nil {
			return err
		}
	}
	return nil
}

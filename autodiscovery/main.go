package main

import (
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/uniris/uniris-core/autodiscovery/adapters/repositories"
	"github.com/uniris/uniris-core/autodiscovery/adapters/services"
	"github.com/uniris/uniris-core/autodiscovery/domain/usecases"
	"github.com/uniris/uniris-core/autodiscovery/infrastructure"
)

func main() {

	port := flag.Int("port", 3545, "GRPC port")
	isGossip := flag.Bool("gossip", true, "Is the peer will gossip")
	flag.Parse()

	peerRepo := &repositories.InMemoryPeerRepository{}
	geoService := services.GeoService{}

	pubKey, err := loadPubKey()
	if err != nil {
		panic(err)
	}

	err = usecases.StartupPeer(peerRepo, geoService, pubKey, *port)
	if err != nil {
		panic(err)
	}

	if *isGossip {
		err = usecases.LoadSeedPeers(services.SeedLoader{}, peerRepo)
		if err != nil {
			panic(err)
		}

		//TODO: loop over every seconds
		err = usecases.StartGossipRound(peerRepo, services.GossipService{})
		if err != nil {
			panic(err)
		}
	}

	if err := infrastructure.StartServer(peerRepo, *port); err != nil {
		panic(err)
	}

}

func loadPubKey() ([]byte, error) {
	path, err := filepath.Abs("./id.pub")
	if err != nil {
		return nil, err
	}
	pubKeyFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer pubKeyFile.Close()

	pubKey, err := ioutil.ReadAll(pubKeyFile)
	if err != nil {
		return nil, err
	}
	return pubKey, nil
}

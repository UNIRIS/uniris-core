package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/uniris/uniris-core/autodiscovery/domain"
	"github.com/uniris/uniris-core/autodiscovery/usecases"
	"github.com/uniris/uniris-core/autodiscovery/usecases/repositories"

	"github.com/uniris/uniris-core/autodiscovery/adapters"
	"github.com/uniris/uniris-core/autodiscovery/infrastructure"
)

func main() {

	log.Print("Autodiscovery starting...")

	peerConf := loadConfiguration()

	log.Printf("GRPC port = %d", peerConf.Port)
	log.Printf("P2P replication factor = %d", peerConf.P2PFactor)
	log.Printf("Public key = %s", peerConf.PublicKey)
	log.Printf("Version = %s", peerConf.Version)

	peerRepo := new(adapters.InMemoryPeerRepository)
	geolocalizer := new(infrastructure.Geolocalizer)

	if err := usecases.StartPeer(peerRepo, geolocalizer, peerConf); err != nil {
		log.Panicln(err)
	}

	go func() {
		if err := infrastructure.StartServer(peerRepo, peerConf.Port); err != nil {
			log.Panicln(err)
		}
	}()

	time.Sleep(2 * time.Second)
	startGossip(peerRepo)
}

func loadConfiguration() domain.PeerConfiguration {
	port := flag.Int("port", 3545, "GRPC port")
	pubKeyFile := flag.String("pub-key-file", "id.pub", "Public key file")
	p2pFactor := flag.Int("p2p-factor", 1, "P2P replication factor")

	flag.Parse()

	pubKey, err := loadPubKey(*pubKeyFile)
	if err != nil {
		log.Panicln(err)
	}

	version, err := infrastructure.GetVersion()
	if err != nil {
		log.Panicln(err)
	}

	return domain.NewPeerConfiguration(version, pubKey, *port, *p2pFactor)
}

func loadPubKey(pubKeyFile string) ([]byte, error) {
	path, err := filepath.Abs(pubKeyFile)
	if err != nil {
		return nil, err
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	pubKey, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return pubKey, nil
}

func startGossip(peerRepo repositories.PeerRepository) {
	seedReader := new(infrastructure.SeedReader)
	messenger := new(infrastructure.GrpcClient)

	log.Print("Gossip starting...")
	ticker := time.NewTicker(1 * time.Second)
	for range ticker.C {
		go func() {
			if err := usecases.StartGossipRound(peerRepo, seedReader, messenger); err != nil {
				log.Printf("Gossip failure %s", err.Error())
			}
		}()
	}
}

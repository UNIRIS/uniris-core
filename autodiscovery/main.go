package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/uniris/uniris-core/autodiscovery/adapters/repositories"
	"github.com/uniris/uniris-core/autodiscovery/adapters/services"
	"github.com/uniris/uniris-core/autodiscovery/domain/usecases"
	"github.com/uniris/uniris-core/autodiscovery/infrastructure"
)

func main() {

	log.Println("Autodiscovery starting...")

	port := flag.Int("port", 3545, "GRPC port")
	initGossip := flag.Bool("init-gossip", true, "Is the node must init gossip")
	pubKeyFile := flag.String("pub-key-file", "id.pub", "Public key file")
	flag.Parse()

	log.Printf("GRPC port = %d\n", *port)
	log.Printf("Initialize gossip = %v\n", *initGossip)

	peerRepo := &repositories.InMemoryPeerRepository{}
	geoService := services.GeoService{}

	pubKey, err := loadPubKey(*pubKeyFile)
	if err != nil {
		log.Panicln(err)
	}

	log.Printf("Public key = %s\n\n", pubKey)
	if err = usecases.StartupPeer(peerRepo, geoService, pubKey, *port); err != nil {
		log.Panicln(err)
	}

	if *initGossip {
		ticker := time.NewTicker(1 * time.Second)
		for range ticker.C {
			go func() {
				if err = usecases.StartGossipRound(&services.SeedLoader{}, peerRepo, &services.GossipService{}); err != nil {
					log.Fatalln(fmt.Sprintf("Gossip failure %s", err.Error()))
				}
			}()
		}
	}

	if err := infrastructure.StartServer(peerRepo, *port); err != nil {
		log.Panicln(err)
	}

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

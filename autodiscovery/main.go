package main

import (
	"log"
	"time"

	"github.com/uniris/uniris-core/autodiscovery/core/ports"

	"github.com/uniris/uniris-core/autodiscovery/adapters"
	"github.com/uniris/uniris-core/autodiscovery/core/usecases"
	"github.com/uniris/uniris-core/autodiscovery/infrastructure"
)

func main() {

	log.Print("Autodiscovery starting...")

	httpClient := new(infrastructure.HTTPClient)
	fileReader := new(infrastructure.FileReader)
	repo := new(adapters.InMemoryPeerRepository)
	metric := new(adapters.MetricReader)
	conf := adapters.NewConfigurationReader(*fileReader, *httpClient)

	if err := usecases.StartPeer(repo, conf); err != nil {
		log.Fatal(err)
	}

	go func() {
		if err := startServer(repo, conf, metric); err != nil {
			log.Fatal(err)
		}
	}()

	time.Sleep(2 * time.Second)

	broker := new(adapters.GrpcGossipBroker)
	notifier := new(adapters.InMemoryDiscoveryNotifier)

	log.Print("Seed loading...")
	if err := usecases.LoadSeeds(repo, conf); err != nil {
		log.Fatal(err)
	}
	log.Print("Gossip starting...")

	ticker := time.NewTicker(1 * time.Second)
	for range ticker.C {
		go func() {
			if err := usecases.StartGossipRound(repo, broker, notifier); err != nil {
				log.Printf("Gossip failure %s", err.Error())
			}
		}()
	}
}

func startServer(r ports.PeerRepository, c ports.ConfigurationReader, m ports.MetricReader) error {
	port, err := c.GetPort()
	if err != nil {
		return err
	}

	lis, err := infrastructure.NewNetListener(port)
	if err != nil {
		return err
	}

	grpcServer := adapters.NewGRPC(r, c, m)
	log.Printf("Server listening on port %d", port)
	if err := grpcServer.Serve(lis); err != nil {
		return err
	}
	return nil
}

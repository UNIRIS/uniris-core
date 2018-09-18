package main

import (
	"log"
	"time"

	"github.com/uniris/uniris-core/autodiscovery/pkg/boostraping"
	"github.com/uniris/uniris-core/autodiscovery/pkg/gossip"
	"github.com/uniris/uniris-core/autodiscovery/pkg/mock"
)

func main() {

	pbKey, port, p2pfactor := loadConfiguration()
	repo := mock.NewRepository()
	loc := mock.NewPeerLocalizer()

	go func() {
		if err := startServer(port); err != nil {
			log.Fatal(err)
		}
	}()

	boot := boostraping.NewService(repo, loc)
	startPeer, err := boot.Startup(pbKey, port, p2pfactor)
	if err != nil {
		log.Fatal(err)
	}

	seedReader := file.SeedReader{}
	if err := seeding.NewService(seedReader, repo).LoadSeeds(); err != nil {
		log.Fatal(err)
	}

	msg := mock.NewGossipMessenger()
	notif := mock.NewGossipNotifier()

	gs := gossip.NewService(repo, msg, notif)

	ticker := time.NewTicker(1 * time.Second)
	for range ticker.C {
		if err := gs.Run(startPeer); err != nil {
			log.Printf("Gossip failure: %s", err)
		}
	}
}

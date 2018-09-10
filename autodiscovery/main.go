package main

import (
	"github.com/uniris/uniris-core/autodiscovery/domain/repositories"
	"github.com/uniris/uniris-core/autodiscovery/domain/services"
	"github.com/uniris/uniris-core/autodiscovery/domain/usecases"
	"github.com/uniris/uniris-core/autodiscovery/infrastructure"
)

func main() {

	peerRepo := &repositories.PeerRepository{}

	if err := infrastructure.StartServer(peerRepo); err != nil {
		panic(err)
	}

	usecases.LoadSeedPeers(&services.SeedLoader{}, peerRepo)
	go usecases.DiscoverPeers(peerRepo)
}

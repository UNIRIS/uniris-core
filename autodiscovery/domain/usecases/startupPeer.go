package usecases

import (
	"time"

	"github.com/uniris/uniris-core/autodiscovery/domain/entities"
	"github.com/uniris/uniris-core/autodiscovery/domain/repositories"
	"github.com/uniris/uniris-core/autodiscovery/domain/services"
)

//StartupPeer initializes a new peer
func StartupPeer(repo repositories.PeerRepository, geo services.GeolocService, publicKey []byte, port int) error {
	geoLoc, err := geo.Lookup()
	if err != nil {
		return err
	}

	peer := &entities.Peer{
		PublicKey: publicKey,
		IP:        geoLoc.IP,
		Port:      port,
		Heartbeat: entities.PeerHeartbeat{
			GenerationTime: time.Now(),
		},
		Category: entities.DiscoveredCategory,
		AppState: entities.PeerAppState{
			GeoCoordinates: entities.Coordinates{
				Lat: geoLoc.Lat,
				Lon: geoLoc.Lon,
			},
			State: entities.BootstrapingState,
		},
	}
	if err := repo.SetLocalPeer(peer); err != nil {
		return err
	}

	return nil
}

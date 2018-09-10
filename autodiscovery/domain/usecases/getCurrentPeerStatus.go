package usecases

import (
	"github.com/uniris/uniris-core/autodiscovery/domain/entities"
	"github.com/uniris/uniris-core/autodiscovery/domain/services"
)

//GetCurrentPeerStatus retrieves the current peer information and state
func GetCurrentPeerStatus(geo services.GeolocService) (*entities.Peer, error) {
	//Get the geolocalization information
	geoLoc, err := geo.Lookup()
	if err != nil {
		return nil, err
	}

	peer := &entities.Peer{
		IP: geoLoc.IP,
		Details: entities.PeerDetails{
			GeoCoordinates: entities.Coordinates{
				Lat: geoLoc.Lat,
				Lon: geoLoc.Lon,
			},
		},
		//TODO: define others details properties
	}
	peer.SetElapsedHeartbeats()
	return peer, nil
}

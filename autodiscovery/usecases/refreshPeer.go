package usecases

import (
	"github.com/uniris/uniris-core/autodiscovery/domain"
	"github.com/uniris/uniris-core/autodiscovery/usecases/ports"
	"github.com/uniris/uniris-core/autodiscovery/usecases/repositories"
)

//RefreshPeer recalculate and change
func RefreshPeer(repo repositories.PeerRepository, peer *domain.Peer, geo ports.Geolocalizer) error {
	geoPos, err := geo.Lookup()
	if err != nil {
		return err
	}

	//TODO: computes system statistics to generate new state
	version := ""
	status := domain.Ok
	p2pfactor := 1

	newState := peer.State
	if newState == nil {
		newState = &domain.PeerState{
			Status:    status,
			Version:   version,
			P2PFactor: p2pfactor,
		}
	}
	newState.GeoPosition = domain.GeoPosition{Lon: geoPos.Lon, Lat: geoPos.Lon}
	peer.Refresh(geoPos.IP, peer.Port, peer.GenerationTime, newState)

	if err := StorePeer(repo, *peer); err != nil {
		return err
	}

	return nil
}

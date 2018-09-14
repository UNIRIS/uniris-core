package usecases

import (
	"github.com/uniris/uniris-core/autodiscovery/domain"
	"github.com/uniris/uniris-core/autodiscovery/usecases/ports"
	"github.com/uniris/uniris-core/autodiscovery/usecases/repositories"
)

//StartPeer creates a new peer on startup and defined as owned on the repository
func StartPeer(repo repositories.PeerRepository, geo ports.Geolocalizer, conf domain.PeerConfiguration) error {
	geoPos, err := geo.Lookup()
	if err != nil {
		return err
	}
	peer := domain.NewPeer(conf.PublicKey, geoPos.IP, conf.Port, true)
	peer.State = domain.NewPeerState(
		domain.Bootstraping,
		conf.Version,
		geoPos,
		conf.P2PFactor,
	)
	if err = StorePeer(repo, peer); err != nil {
		return err
	}
	return nil
}

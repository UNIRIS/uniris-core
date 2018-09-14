package usecases

import (
	"log"
	"time"

	"github.com/uniris/uniris-core/autodiscovery/core/domain"
	"github.com/uniris/uniris-core/autodiscovery/core/ports"
)

//StartPeer creates a new peer on startup and defined as owned on the repository
func StartPeer(repo ports.PeerRepository, conf ports.ConfigurationReader) error {

	ip, err := conf.GetIP()
	if err != nil {
		return err
	}

	geoPos, err := conf.GetGeoPosition()
	if err != nil {
		return err
	}

	ver, err := conf.GetVersion()
	if err != nil {
		return err
	}

	pbKey, err := conf.GetPublicKey()
	if err != nil {
		return err
	}

	port, err := conf.GetPort()
	if err != nil {
		return err
	}

	p2pFactor, err := conf.GetP2PFactor()
	if err != nil {
		return err
	}

	peer := domain.Peer{
		GenerationTime: time.Now(),
		IP:             ip,
		IsOwned:        true,
		Port:           port,
		PublicKey:      pbKey,
		State: &domain.PeerState{
			Status:      domain.Bootstraping,
			Version:     ver,
			GeoPosition: geoPos,
			P2PFactor:   p2pFactor,
		},
	}
	if err := repo.InsertPeer(peer); err != nil {
		return err
	}

	log.Println("Owned peer initialized")

	return nil
}

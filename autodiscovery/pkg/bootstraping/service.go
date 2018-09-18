package bootstraping

import (
	"net"

	discovery "github.com/uniris/uniris-core/autodiscovery/pkg"
)

//PeerLocalizer is the interface that provides methods for the peer localization
type PeerLocalizer interface {
	GetIP() (net.IP, error)
	GetGeoPosition() (discovery.PeerPosition, error)
}

//Service is the interface that provide methods for the peer's bootstraping
type Service interface {
	Startup(pbKey []byte, port int, p2pFactor int, ver string) (discovery.Peer, error)
	LoadSeeds(ss []discovery.Seed) error
}

type service struct {
	repo discovery.Repository
	loc  PeerLocalizer
}

//Startup creates a new peer initiator, locates and stores it
func (s service) Startup(pbKey []byte, port int, p2pFactor int, ver string) (p discovery.Peer, err error) {
	pos, err := s.loc.GetGeoPosition()
	if err != nil {
		return
	}

	ip, err := s.loc.GetIP()
	if err != nil {
		return
	}

	p = discovery.NewStartupPeer(pbKey, ip, port, ver, pos, p2pFactor)
	if err = s.repo.AddPeer(p); err != nil {
		return
	}
	return
}

//LoadSeeds stores the provided seed peers
func (s service) LoadSeeds(ss []discovery.Seed) error {
	for _, sd := range ss {
		if err := s.repo.AddSeed(sd); err != nil {
			return err
		}
	}

	return nil
}

//NewService creates a bootstraping service its dependencies
func NewService(repo discovery.Repository, loc PeerLocalizer) Service {
	return &service{
		repo: repo,
		loc:  loc,
	}
}

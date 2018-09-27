package bootstraping

import (
	discovery "github.com/uniris/uniris-core/autodiscovery/pkg"
	"github.com/uniris/uniris-core/autodiscovery/pkg/monitoring"
)

//Service is the interface that provide methods for the peer's bootstraping
type Service interface {
	Startup(pbKey []byte, port int, ver string) (discovery.Peer, error)
	LoadSeeds(seeds []discovery.Seed) error
}

type service struct {
	repo discovery.Repository
	pp   monitoring.PeerPositionner
	pn   monitoring.PeerNetworker
}

//Startup creates a new peer initiator, locates and stores it
func (s service) Startup(pbKey []byte, port int, ver string) (p discovery.Peer, err error) {
	pos, err := s.pp.Position()
	if err != nil {
		return
	}

	ip, err := s.pn.IP()
	if err != nil {
		return
	}

	p = discovery.NewStartupPeer(pbKey, ip, port, ver, pos)
	if err = s.repo.SetPeer(p); err != nil {
		return
	}

	return
}

//LoadSeeds stores the provided seed peers
func (s service) LoadSeeds(ss []discovery.Seed) error {
	for _, sd := range ss {
		if err := s.repo.SetSeed(sd); err != nil {
			return err
		}
	}

	return nil
}

//NewService creates a bootstraping service its dependencies
func NewService(repo discovery.Repository, pp monitoring.PeerPositionner, pn monitoring.PeerNetworker) Service {
	return &service{
		repo: repo,
		pp:   pp,
		pn:   pn,
	}
}

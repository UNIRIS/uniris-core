package bootstraping

import (
	"net"

	discovery "github.com/uniris/uniris-core/autodiscovery/pkg"
)

//PeerPositionner is the interface that provide methods to identity the peer geo position
type PeerPositionner interface {
	//Position lookups the peer's geographic position
	Position() (discovery.PeerPosition, error)
}

//PeerNetworker is the interface that provides methods to get the peer network information
type PeerNetworker interface {
	//IP lookups the peer's IP
	IP() (net.IP, error)
}

//Service is the interface that provide methods for the peer's bootstraping
type Service interface {
	Startup(pbKey []byte, port uint16, ver string) (discovery.Peer, error)
	LoadSeeds(seeds []discovery.Seed) error
}

type service struct {
	repo discovery.Repository
	pp   PeerPositionner
	pn   PeerNetworker
}

//Startup creates a new peer initiator, locates and stores it
func (s service) Startup(pbKey []byte, port uint16, ver string) (p discovery.Peer, err error) {
	pos, err := s.pp.Position()
	if err != nil {
		return
	}

	ip, err := s.pn.IP()
	if err != nil {
		return
	}

	p = discovery.NewStartupPeer(pbKey, ip, port, ver, pos)
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
func NewService(repo discovery.Repository, pp PeerPositionner, pn PeerNetworker) Service {
	return &service{
		repo: repo,
		pp:   pp,
		pn:   pn,
	}
}

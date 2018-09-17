package boostraping

import (
	"net"

	"github.com/uniris/uniris-core/autodiscovery/pkg/discovery"
)

type Repository interface {
	AddPeer(discovery.Peer) error
}

type Service interface {
	Startup(pbKey []byte, port int, version string, p2pFactor int) error
}

type PeerLocalizer interface {
	GetIP() (net.IP, error)
	GetGeoPosition() (discovery.PeerPosition, error)
}

type service struct {
	repo Repository
	loc  PeerLocalizer
}

func NewService(repo Repository, loc PeerLocalizer) Service {
	return &service{
		repo: repo,
		loc:  loc,
	}
}

func (s service) Startup(pbKey []byte, port int, version string, p2pFactor int) error {

	pos, err := s.loc.GetGeoPosition()
	if err != nil {
		return err
	}

	ip, err := s.loc.GetIP()
	if err != nil {
		return err
	}

	p := discovery.StartPeer(pbKey, ip, port, version, pos, p2pFactor)
	if err := s.repo.AddPeer(p); err != nil {
		return err
	}
	return nil
}

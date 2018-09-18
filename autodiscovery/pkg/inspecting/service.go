package inspecting

import (
	discovery "github.com/uniris/uniris-core/autodiscovery/pkg"
)

//PeerMonitor is the interface that provides methods for the peer monitoring
type PeerMonitor interface {
	GetStatus() (discovery.PeerStatus, error)
	GetCPULoad() (string, error)
	GetFreeDiskSpace() (float64, error)
	GetIOWaitRate() (float64, error)
}

//Service defines the interface for the peer inpsection
type Service interface {
	RefreshPeer(*discovery.Peer) error
}

type service struct {
	pr   PeerMonitor
	repo discovery.Repository
}

//RefreshPeer updates the peer's metrics retrieved from the peer monitor
func (s service) RefreshPeer(p *discovery.Peer) error {
	status, err := s.pr.GetStatus()
	if err != nil {
		return err
	}

	cpu, err := s.pr.GetCPULoad()
	if err != nil {
		return err
	}

	disk, err := s.pr.GetFreeDiskSpace()
	if err != nil {
		return err
	}

	io, err := s.pr.GetIOWaitRate()
	if err != nil {
		return err
	}

	p.Refresh(status, disk, cpu, io)
	if err := s.repo.UpdatePeer(*p); err != nil {
		return err
	}
	return nil
}

//NewService creates a new inspection service
func NewService(repo discovery.Repository, pr PeerMonitor) Service {
	return service{
		repo: repo,
		pr:   pr,
	}
}

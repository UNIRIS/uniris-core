package inspecting

import (
	discovery "github.com/uniris/uniris-core/autodiscovery/pkg"
)

//PeerMonitor is the interface that provides methods for the peer monitoring
type PeerMonitor interface {

	//GetStatus computes the peer's status according to the health state of the system
	GetStatus() (discovery.PeerStatus, error)

	//GetCPULoad retrieves the load on the peer's CPU
	GetCPULoad() (string, error)

	//GetFreeDiskSpace retrieves the available free disk space of the peer
	GetFreeDiskSpace() (float64, error)

	//GetIOWaitRate computes the rate of the I/O operations of the peer
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

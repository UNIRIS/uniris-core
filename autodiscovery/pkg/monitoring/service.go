package monitoring

import (
	discovery "github.com/uniris/uniris-core/autodiscovery/pkg"
)

//Service defines the interface for the peer monitoring
type Service interface {
	RefreshOwnedPeer() error
}

type service struct {
	repo discovery.Repository
	mon  discovery.PeerMonitor
}

//RefreshOwnedPeer updates the owned peer's metrics retrieved from the peer monitor
func (s service) RefreshOwnedPeer() error {
	p, err := s.repo.GetOwnedPeer()
	if err != nil {
		return err
	}

	status, err := s.mon.Status()
	if err != nil {
		return err
	}

	cpu, err := s.mon.CPULoad()
	if err != nil {
		return err
	}

	disk, err := s.mon.FreeDiskSpace()
	if err != nil {
		return err
	}

	io, err := s.mon.IOWaitRate()
	if err != nil {
		return err
	}

	p.Refresh(status, disk, cpu, io)
	if err := s.repo.SetPeer(p); err != nil {
		return err
	}
	return nil
}

//NewService creates a new monitoring service
func NewService(repo discovery.Repository, mon discovery.PeerMonitor) Service {
	return service{
		repo: repo,
		mon:  mon,
	}
}

package monitoring

import (
	discovery "github.com/uniris/uniris-core/autodiscovery/pkg"
)

//Watcher is the interface that provides methods for the peer monitoring
type Watcher interface {

	//Status computes the peer's status according to the health state of the system
	Status() (discovery.PeerStatus, error)

	//CPULoad retrieves the load on the peer's CPU
	CPULoad() (string, error)

	//FreeDiskSpace retrieves the available free disk space of the peer
	FreeDiskSpace() (float64, error)

	//IOWaitRate computes the rate of the I/O operations of the peer
	IOWaitRate() (float64, error)
}

//PeerWatcher define the interface to retrieve the different state of the process running on a peer
type PeerWatcher interface {
	CheckProcessStates() (bool, error)
	CheckInternetState() (bool, error)
	CheckNtpState() (bool, error)
}

//SeedDiscoverdNodeWatcher define the interface to check the number of discovered node by a seed
type SeedDiscoverdNodeWatcher interface {
	GetSeedDiscoveredNode() (int, error)
}

//Service defines the interface for the peer inpsection
type Service interface {
	RefreshPeer(discovery.Peer) error
}

type service struct {
	w    Watcher
	repo discovery.Repository
}

//RefreshPeer updates the peer's metrics retrieved from the peer monitor
func (s service) RefreshPeer(p discovery.Peer) error {
	status, err := s.w.Status()
	if err != nil {
		return err
	}

	cpu, err := s.w.CPULoad()
	if err != nil {
		return err
	}

	disk, err := s.w.FreeDiskSpace()
	if err != nil {
		return err
	}

	io, err := s.w.IOWaitRate()
	if err != nil {
		return err
	}

	p.Refresh(status, disk, cpu, io)
	if err := s.repo.UpdatePeer(p); err != nil {
		return err
	}
	return nil
}

//NewService creates a new inspection service
func NewService(repo discovery.Repository, w Watcher) Service {
	return service{
		repo: repo,
		w:    w,
	}
}

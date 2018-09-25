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

	//P2PFactor retrieves the replication factor from the AI service
	P2PFactor() (int32, error)
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

	p2pFactor, err := s.w.P2PFactor()
	if err != nil {
		return err
	}

	if err := p.Refresh(status, disk, cpu, p2pFactor); err != nil {
		return err
	}
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

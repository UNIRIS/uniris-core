package monitoring

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
	IP() (net.IP, error)
	CheckNtpState() error
	CheckInternetState() error
}

//RobotWatcher is the interface that provides methods for the uniris robot process monitoring
type RobotWatcher interface {
	CheckAutodiscoveryProcess(discoveryPort int) error
	CheckMiningProcess() error
	CheckDataProcess() error
	CheckAIProcess() error
	CheckRedisProcess() error
	CheckScyllaDbProcess() error
	CheckRabbitmqProcess() error
}

//PeerMonitor is the interface that provides methods for the peer monitoring
type PeerMonitor interface {

	//CPULoad retrieves the load on the peer's CPU
	CPULoad() (string, error)

	//FreeDiskSpace retrieves the available free disk space of the peer
	FreeDiskSpace() (float64, error)

	//P2PFactor retrieves the replication factor from the AI service
	//and defines the number of robots that should validate a transaction
	P2PFactor() (int, error)
}

//Service defines the interface for the peer inpsection
type Service interface {
	PeerStatus(p discovery.Peer) (discovery.PeerStatus, error)
	RefreshPeer(discovery.Peer) error
}

type service struct {
	mon  PeerMonitor
	repo discovery.Repository
	pn   PeerNetworker
	sdc  discovery.SeedDiscoveryCounter
	rbtW RobotWatcher
}

//Status computes the peer's status according to the health state of the system
func (s service) PeerStatus(p discovery.Peer) (discovery.PeerStatus, error) {
	err := s.checkProcesses(p.Identity().Port())
	if err != nil {
		return discovery.FaultStatus, err
	}

	err = s.pn.CheckInternetState()
	if err != nil {
		return discovery.FaultStatus, err
	}
	err = s.pn.CheckNtpState()
	if err != nil {
		return discovery.StorageOnlyStatus, err
	}

	seedDiscoveries, err := s.sdc.Average()
	if err != nil {
		return discovery.FaultStatus, err
	}
	if seedDiscoveries == 0 {
		return discovery.BootstrapingStatus, nil
	}

	if t := p.HeartbeatState().ElapsedHeartbeats(); t < discovery.BootStrapingMinTime && seedDiscoveries > p.AppState().DiscoveredPeersNumber() {
		return discovery.BootstrapingStatus, nil
	} else if t < discovery.BootStrapingMinTime && seedDiscoveries <= p.AppState().DiscoveredPeersNumber() {
		return discovery.OkStatus, nil
	} else {
		return discovery.OkStatus, nil
	}
}

//RefreshPeer updates the peer's metrics retrieved from the peer monitor
func (s service) RefreshPeer(p discovery.Peer) error {
	status, err := s.PeerStatus(p)
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

	dp, err := s.repo.CountKnownPeers()
	if err != nil {
		return err
	}

	p2p, err := s.mon.P2PFactor()
	if err != nil {
		return err
	}

	if err := p.Refresh(status, disk, cpu, p2p, dp); err != nil {
		return err
	}
	if err := s.repo.UpdatePeer(p); err != nil {
		return err
	}
	return nil
}

//checkProcesses check the different state of the differents necessary services running on the peer
func (s service) checkProcesses(discoveryPort int) error {
	err := s.rbtW.CheckAutodiscoveryProcess(discoveryPort)
	if err != nil {
		return err
	}
	err = s.rbtW.CheckDataProcess()
	if err != nil {
		return err
	}
	err = s.rbtW.CheckMiningProcess()
	if err != nil {
		return err
	}
	err = s.rbtW.CheckAIProcess()
	if err != nil {
		return err
	}
	err = s.rbtW.CheckRedisProcess()
	if err != nil {
		return err
	}
	err = s.rbtW.CheckScyllaDbProcess()
	if err != nil {
		return err
	}
	err = s.rbtW.CheckRabbitmqProcess()
	if err != nil {
		return err
	}

	return nil
}

//NewService creates a new inspection service
func NewService(repo discovery.Repository, mon PeerMonitor, pn PeerNetworker, rbtW RobotWatcher) Service {
	return service{
		repo: repo,
		mon:  mon,
		pn:   pn,
		sdc:  discovery.NewSeedDiscoveryCounter(repo),
		rbtW: rbtW,
	}
}

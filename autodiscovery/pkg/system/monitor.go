package system

import (
	"errors"
	"net"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/beevik/ntp"
	discovery "github.com/uniris/uniris-core/autodiscovery/pkg"
	"github.com/uniris/uniris-core/autodiscovery/pkg/monitoring"
)

const (
	cdns          = "uniris.io"
	upmaxoffset   = 300
	downmaxoffset = -300
)

var (
	cntp = [...]string{"1.pool.ntp.org", "2.pool.ntp.org", "3.pool.ntp.org", "4.pool.ntp.org"}
)

//PeerWatcher define the interface to retrieve the different state of the process running on a peer
type PeerWatcher interface {
	CheckProcessStates(p discovery.Peer) (bool, error)
	CheckInternetState() (bool, error)
	CheckNtpState() (bool, error)
}

type peerWatcher struct {
}

//GetProcessStates check the different state of the differents necessary services running on the peer
func (Pwatcher *peerWatcher) CheckProcessStates(p discovery.Peer) (bool, error) {
	_, err := CheckAutodiscoveryProcess(p)
	if err != nil {
		return false, err
	}
	_, err = CheckDataProcess()
	if err != nil {
		return false, err
	}
	_, err = CheckMiningProcess()
	if err != nil {
		return false, err
	}
	_, err = CheckAIProcess()
	if err != nil {
		return false, err
	}
	_, err = CheckScyllaProcess()
	if err != nil {
		return false, err
	}
	_, err = CheckRedisProcess()
	if err != nil {
		return false, err
	}
	_, err = CheckRabitmqProcess()
	if err != nil {
		return false, err
	}

	return true, nil
}

//CheckInternetConfig check internet configuration on the node
func (Pwatcher *peerWatcher) CheckInternetState() (bool, error) {
	_, err := net.LookupIP(cdns)
	if err != nil {
		return false, err
	}
	return true, nil
}

//CheckNtp check time synchonization on the node
func (Pwatcher *peerWatcher) CheckNtpState() (bool, error) {
	for _, ntps := range cntp {
		r, err := ntp.QueryWithOptions(ntps, ntp.QueryOptions{Version: 4})
		if err == nil {
			if (int64(r.ClockOffset/time.Second) < downmaxoffset) || (int64(r.ClockOffset/time.Second) > upmaxoffset) {
				return false, nil
			}
			return true, nil
		}
		//log.Print(err)
	}
	return false, errors.New("Could not get reply from ntp servers")
}

//SeedDiscoverdPeerWatcher define the interface to check the number of discovered node by a seed
type SeedDiscoverdPeerWatcher interface {
	GetSeedDiscoveredPeer(rep discovery.Repository) (int, error)
}

type seedDiscoverdPeerWatcher struct {
}

//GetSeedDiscoveredPeer report the average of node detected by the differents known seeds
func (SdnWatcher *seedDiscoverdPeerWatcher) GetSeedDiscoveredPeer(rep discovery.Repository) (int, error) {
	listseed, err := rep.ListSeedPeers()
	if err != nil {
		return 0, err
	}
	avg := 0
	for i := 0; i < len(listseed); i++ {
		ipseed := listseed[i].IP
		p, err := rep.GetPeerByIP(ipseed)
		if err == nil {
			avg += p.DiscoveredPeers()
		}
	}
	avg = avg / len(listseed)
	return avg, nil
}

type watcher struct {
	Pwatcher   peerWatcher
	SdnWatcher seedDiscoverdPeerWatcher
	rep        discovery.Repository
}

//Status computes the peer's status according to the health state of the system
func (w watcher) Status(p discovery.Peer, repo discovery.Repository) (discovery.PeerStatus, error) {

	procState, err := w.Pwatcher.CheckProcessStates(p)
	if err != nil {
		return discovery.FaultStatus, err
	}
	if !procState {
		return discovery.FaultStatus, nil
	}
	internetState, err := w.Pwatcher.CheckInternetState()
	if err != nil {
		return discovery.FaultStatus, err
	}
	if !internetState {
		return discovery.FaultStatus, nil
	}
	ntpState, err := w.Pwatcher.CheckNtpState()
	if err != nil {
		return discovery.FaultStatus, err
	}
	if !ntpState {
		return discovery.StorageOnlyStatus, nil
	}
	seedDn, err := w.SdnWatcher.GetSeedDiscoveredPeer(repo)
	if err != nil {
		return discovery.FaultStatus, err
	}
	if seedDn == 0 {
		return discovery.BootstrapingStatus, nil
	}
	if t := p.GetElapsedHeartbeats(); t < discovery.BootStrapingMinTime && seedDn > p.DiscoveredPeers() {
		return discovery.BootstrapingStatus, nil
	} else if t < discovery.BootStrapingMinTime && seedDn <= p.DiscoveredPeers() {
		return discovery.OkStatus, nil
	} else {
		return discovery.OkStatus, nil
	}
}

//CPULoad retrieves the load on the peer's CPU
func (w watcher) CPULoad() (string, error) {
	cmd := exec.Command("cat", "/proc/loadavg")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "--", err
	}
	return string(out), nil
}

//FreeDiskSpace retrieves the available free disk (k bytes) space of the peer
func (w watcher) FreeDiskSpace() (float64, error) {
	var stat syscall.Statfs_t
	wd, err := os.Getwd()
	if err != nil {
		return 0.0, err
	}
	syscall.Statfs(wd, &stat)
	return float64((stat.Bavail * uint64(stat.Bsize)) / 1024), nil
}

//DiscoverdPeer computes the number of peer discovered by the local peer
func (w watcher) DiscoveredPeer(rep discovery.Repository) (int, error) {
	l, err := rep.ListKnownPeers()
	if err != nil {
		return 0, err
	}
	return len(l), nil
}

//P2PFactor request the update P2PFactor from the AI Daemon
func (w watcher) P2PFactor() (int, error) {
	return 0, nil
}

//NewSystemWatcher creates an instance which implements monitoring.Watcher
func NewSystemWatcher(rep discovery.Repository) monitoring.Watcher {
	return watcher{
		rep:        rep,
		Pwatcher:   peerWatcher{},
		SdnWatcher: seedDiscoverdPeerWatcher{},
	}
}

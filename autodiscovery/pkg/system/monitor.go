package system

import (
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
	cntp          = "pool.ntp.org"
	upmaxoffset   = 300
	downmaxoffset = -300
)

type peerWatcher struct{}

//GetProcessStates check the different state of the differents necessary services running on the peer
func (Pwatcher *peerWatcher) CheckProcessStates() (bool, error) {
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
	r, err := ntp.QueryWithOptions(cntp, ntp.QueryOptions{Version: 4})

	if err != nil {
		return false, err
	}

	if (int64(r.ClockOffset/time.Second) < downmaxoffset) || (int64(r.ClockOffset/time.Second) > upmaxoffset) {
		return false, nil
	}
	return true, nil
}

type seedDiscoverdNodeWatcher struct{}

//GetSeedDiscoveredNode report the average of node detected by the different known seed
func (SdnWatcher *seedDiscoverdNodeWatcher) GetSeedDiscoveredNode() (int, error) {
	return 5, nil
}

type watcher struct {
	Pwatcher   peerWatcher
	SdnWatcher seedDiscoverdNodeWatcher
	rep        discovery.Repository
}

//Status computes the peer's status according to the health state of the system
func (w watcher) Status() (discovery.PeerStatus, error) {

	selfpeer, err := w.rep.GetOwnedPeer()
	if err != nil {
		return discovery.FaultStatus, err
	}

	procState, err := w.Pwatcher.CheckProcessStates()
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

	seedDn, err := w.SdnWatcher.GetSeedDiscoveredNode()

	if t := selfpeer.GetElapsedHeartbeats(); t < discovery.BootStrapingMinTime && seedDn > selfpeer.DiscoveredNodes() {
		return discovery.BootstrapingStatus, nil
	} else if t < discovery.BootStrapingMinTime && seedDn <= selfpeer.DiscoveredNodes() {
		return discovery.OkStatus, nil
	} else if t > discovery.BootStrapingMinTime && seedDn > selfpeer.DiscoveredNodes() {
		return discovery.BootstrapingStatus, nil
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

//IOWaitRate computes the rate of the I/O operations of the peer
func (w watcher) IOWaitRate() (float64, error) {
	return 0.0, nil
}

//NewSystemWatcher creates an instance which implements monitoring.Watcher
func NewSystemWatcher(rep discovery.Repository) monitoring.Watcher {
	return watcher{
		rep: rep,
	}
}

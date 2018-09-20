package system

import (
	discovery "github.com/uniris/uniris-core/autodiscovery/pkg"
	"github.com/uniris/uniris-core/autodiscovery/pkg/monitoring"
)

type watcher struct{}

//Status computes the peer's status according to the health state of the system
func (w watcher) Status() (discovery.PeerStatus, error) {
	return discovery.OkStatus, nil
}

//CPULoad retrieves the load on the peer's CPU
func (w watcher) CPULoad() (string, error) {
	return "", nil
}

//FreeDiskSpace retrieves the available free disk space of the peer
func (w watcher) FreeDiskSpace() (float64, error) {
	return 0.0, nil
}

//IOWaitRate computes the rate of the I/O operations of the peer
func (w watcher) IOWaitRate() (float64, error) {
	return 0.0, nil
}

//NewSystemWatcher creates an instance which implements monitoring.Watcher
func NewSystemWatcher() monitoring.Watcher {
	return watcher{}
}

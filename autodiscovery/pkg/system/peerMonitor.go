package system

import (
	discovery "github.com/uniris/uniris-core/autodiscovery/pkg"
)

type monitor struct{}

//Status computes the peer's status according to the health state of the system
func (m monitor) Status() (discovery.PeerStatus, error) {
	return discovery.OkStatus, nil
}

//CPULoad retrieves the load on the peer's CPU
func (m monitor) CPULoad() (string, error) {
	return "", nil
}

//FreeDiskSpace retrieves the available free disk space of the peer
func (m monitor) FreeDiskSpace() (float64, error) {
	return 0.0, nil
}

//IOWaitRate computes the rate of the I/O operations of the peer
func (m monitor) IOWaitRate() (float64, error) {
	return 0.0, nil
}

//NewPeerMonitor implements peer monitor using system metrics
func NewPeerMonitor() discovery.PeerMonitor {
	return monitor{}
}

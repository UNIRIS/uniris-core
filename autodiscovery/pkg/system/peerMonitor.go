package system

import discovery "github.com/uniris/uniris-core/autodiscovery/pkg"

//PeerMonitor implements the PeerMonitor interface to provide peer's metrics monitoring
type PeerMonitor struct{}

//GetStatus computes the peer's status according to the health state of the system
func (pm PeerMonitor) GetStatus() (discovery.PeerStatus, error) {
	return discovery.OkStatus, nil
}

//GetCPULoad retrieves the load on the peer's CPU
func (pm PeerMonitor) GetCPULoad() (string, error) {
	return "", nil
}

//GetFreeDiskSpace retrieves the available free disk space of the peer
func (pm PeerMonitor) GetFreeDiskSpace() (float64, error) {
	return 0.0, nil
}

//GetIOWaitRate computes the rate of the I/O operations of the peer
func (pm PeerMonitor) GetIOWaitRate() (float64, error) {
	return 0.0, nil
}

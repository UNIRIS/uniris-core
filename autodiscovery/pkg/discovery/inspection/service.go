package inspection

import "github.com/uniris/uniris-core/autodiscovery/pkg/discovery"

//Monitor wraps the peer's metrics operations
type Monitor interface {
	Status() (discovery.PeerStatus, error)
	CPULoad() (string, error)
	FreeDiskSpace() (float64, error)
	IOWaitRate() (float64, error)
}

//Service wraps the peer's inspector/metrics operations
type Service struct {
	Monitor
}

//Status retrieve the peer's status
func (i Service) Status() (discovery.PeerStatus, error) {
	return discovery.OkStatus, nil
}

//CPULoad retrieve the peer's cpu load
func (i Service) CPULoad() (string, error) {
	return "0.0.0", nil
}

//FreeDiskSpace retrieve the peer's free disk splace
func (i Service) FreeDiskSpace() (float64, error) {
	return 0, nil

}

//IOWaitRate retrieve the peer's io wait rate
func (i Service) IOWaitRate() (float64, error) {
	return 0, nil
}

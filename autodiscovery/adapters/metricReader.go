package adapters

import "github.com/uniris/uniris-core/autodiscovery/core/domain"

//MetricReader wraps the system perfomance and health checks
type MetricReader struct{}

//GetStatus checks if the peer is healthy or not
func (r MetricReader) GetStatus() (domain.PeerStatus, error) {
	//TODO
	return domain.Ok, nil
}

//GetCPULoad gets the peer CPU load
func (r MetricReader) GetCPULoad() (string, error) {
	//TODO
	return "", nil

}

//GetFreeDiskSpace checks the free space on the peer's disk
func (r MetricReader) GetFreeDiskSpace() (float64, error) {
	//TODO
	return 200000, nil

}

//GetIOWaitRate gets a rate of the number of IO waits
func (r MetricReader) GetIOWaitRate() (float64, error) {
	//TODO
	return 0, nil
}

package ports

import "github.com/uniris/uniris-core/autodiscovery/core/domain"

//MetricReader wraps the peer state metrics analyzing
type MetricReader interface {
	GetStatus() (domain.PeerStatus, error)
	GetCPULoad() (string, error)
	GetFreeDiskSpace() (float64, error)
	GetIOWaitRate() (float64, error)
}

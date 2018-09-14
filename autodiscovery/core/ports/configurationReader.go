package ports

import (
	"net"

	"github.com/uniris/uniris-core/autodiscovery/core/domain"
)

type ConfigurationReader interface {
	GetPublicKey() ([]byte, error)
	GetPort() (int, error)
	GetGeoPosition() (domain.GeoPosition, error)
	GetIP() (net.IP, error)
	GetSeeds() ([]domain.Peer, error)
	GetVersion() (string, error)
	GetP2PFactor() (int, error)
}

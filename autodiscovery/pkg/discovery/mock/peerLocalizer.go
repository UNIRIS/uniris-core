package mock

import (
	"net"

	"github.com/uniris/uniris-core/autodiscovery/pkg/discovery"
	"github.com/uniris/uniris-core/autodiscovery/pkg/discovery/boostraping"
)

type MockPeerLocalizer struct{}

func NewPeerLocalizer() boostraping.PeerLocalizer {
	return MockPeerLocalizer{}
}

func (l MockPeerLocalizer) GetIP() (net.IP, error) {
	return net.ParseIP("127.0.0.1"), nil
}

func (l MockPeerLocalizer) GetGeoPosition() (discovery.PeerPosition, error) {
	return discovery.PeerPosition{
		Lon: 3.5,
		Lat: 65.2,
	}, nil
}

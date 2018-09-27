package mock

import discovery "github.com/uniris/uniris-core/autodiscovery/pkg"

type Positioner struct{}

func (l Positioner) Position() (discovery.PeerPosition, error) {
	return discovery.PeerPosition{
		Lon: 3.5,
		Lat: 65.2,
	}, nil
}

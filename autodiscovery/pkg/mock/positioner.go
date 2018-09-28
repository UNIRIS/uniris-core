package mock

import discovery "github.com/uniris/uniris-core/autodiscovery/pkg"

//Positioner is the struct providing the peer's geo informations
type Positioner struct{}

//Position returns the peer's geo position
func (l Positioner) Position() (discovery.PeerPosition, error) {
	return discovery.PeerPosition{
		Lon: 3.5,
		Lat: 65.2,
	}, nil
}

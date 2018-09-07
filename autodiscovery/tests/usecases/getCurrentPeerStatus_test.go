package tests

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/uniris/uniris-core/autodiscovery/domain/services"
	"github.com/uniris/uniris-core/autodiscovery/domain/usecases"
)

func TestGetCurrentPeerStatus(t *testing.T) {
	peer, err := usecases.GetCurrentPeerStatus(&GeolocService{})
	assert.Nil(t, err)
	assert.Equal(t, "127.0.0.1", peer.IP.String())
	assert.Equal(t, 2.33, peer.Details.GeoCoordinates[0])
	assert.Equal(t, 64.20, peer.Details.GeoCoordinates[1])

	//TODO: test other peer's details
}

type GeolocService struct{}

func (geo *GeolocService) Lookup() (services.GeoLoc, error) {
	return services.GeoLoc{
		IP:  net.ParseIP("127.0.0.1"),
		Lat: 2.33,
		Lon: 64.20,
	}, nil
}

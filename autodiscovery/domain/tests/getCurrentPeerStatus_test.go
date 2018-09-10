package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/uniris/uniris-core/autodiscovery/domain/usecases"
)

func TestGetCurrentPeerStatus(t *testing.T) {
	peer, err := usecases.GetCurrentPeerStatus(&GeolocService{})
	assert.Nil(t, err)
	assert.Equal(t, "127.0.0.1", peer.IP.String())
	assert.Equal(t, 2.33, peer.Details.GeoCoordinates.Lat)
	assert.Equal(t, 64.20, peer.Details.GeoCoordinates.Lon)

	//TODO: test other peer's details
}

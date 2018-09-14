package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
Scenario: Create a peer state
	Given a peer refresh order
	When we create a peer stage
	Then we retrieve the new information
*/
func TestNewPeerState(t *testing.T) {
	state := NewPeerState(Ok, "1.0.1", GeoPosition{Lat: 30, Lon: 2}, 2)
	assert.NotNil(t, state)
	assert.Equal(t, Ok, state.Status)
	assert.Equal(t, float64(30), state.GeoPosition.Lat)
	assert.Equal(t, float64(2), state.GeoPosition.Lon)
	assert.Equal(t, 2, state.P2PFactor)
}

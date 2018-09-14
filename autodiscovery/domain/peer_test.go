package domain

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

/*
Scenario: Create a new peer
	Given a public key, an ip and a port
	When we create a peer
	Then we retrieve a peer with a generation time setted
*/
func TestNewPeer(t *testing.T) {
	peer := NewPeer([]byte("my public key"), net.ParseIP("127.0.0.1"), 3545, true)
	assert.NotNil(t, peer)
	assert.Equal(t, "my public key", string(peer.PublicKey))
	assert.Equal(t, 3545, peer.Port)
	assert.Equal(t, "127.0.0.1", peer.IP.String())
	assert.Equal(t, peer.GenerationTime.Unix(), time.Now().Unix())
	assert.Nil(t, peer.State)
	assert.False(t, peer.IsDiscovered())
}

/*
Scenario: Refresh a created peer
	Given a created peer
	When we refresh a peer with new information
	Then the new information are loaded
*/
func TestRefreshPeer(t *testing.T) {
	peer := NewPeer([]byte("my public key"), net.ParseIP("127.0.0.1"), 3545, true)

	newState := &PeerState{
		Status:      Ok,
		Version:     "1.0.1",
		GeoPosition: GeoPosition{Lat: 50, Lon: 3},
		P2PFactor:   1,
	}
	oldGenerationTime := peer.GenerationTime

	time.Sleep(10 * time.Microsecond)
	peer.Refresh(net.ParseIP("10.0.0.1"), 4000, time.Now(), newState)

	assert.Equal(t, "10.0.0.1", peer.IP.String())
	assert.Equal(t, 4000, peer.Port)
	assert.NotEqual(t, oldGenerationTime, peer.GenerationTime)

	assert.NotNil(t, peer.State)
	assert.Equal(t, Ok, peer.State.Status)
	assert.Equal(t, "1.0.1", peer.State.Version)
	assert.Equal(t, float64(50), peer.State.GeoPosition.Lat)
	assert.Equal(t, float64(3), peer.State.GeoPosition.Lon)
	assert.Equal(t, 1, peer.State.P2PFactor)

	assert.True(t, peer.IsDiscovered())
}

/*
Scenario: Checks the discovering of a peer
	Given a created peer
	When the peer is discovered
	Then the peer's state is filled
*/
func TestIsPeerDiscovered(t *testing.T) {
	peer := NewPeer([]byte("my public key"), net.ParseIP("127.0.0.1"), 3545, true)
	assert.False(t, peer.IsDiscovered())
	newState := new(PeerState)
	peer.State = newState
	assert.False(t, peer.IsDiscovered())
	peer.State.Status = Ok
	assert.True(t, peer.IsDiscovered())
}

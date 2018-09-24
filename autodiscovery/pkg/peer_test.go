package discovery

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

/*
Scenario: Create a startup peer
	Given some inputs parameters
	When we create a peer that startup
	Then we get a new peer with the status bootstraping and specified as owned
*/
func TestNewPeer(t *testing.T) {
	p := NewStartupPeer([]byte("key"), net.ParseIP("127.0.0.1"), 3000, "1.0", PeerPosition{Lat: 50.0, Lon: 3.0}, 1)
	assert.NotNil(t, p)
	assert.Equal(t, "key", string(p.PublicKey()))
	assert.Equal(t, "127.0.0.1", p.IP().String())
	assert.Equal(t, 3000, p.Port())
	assert.Equal(t, "1.0", p.Version())
	assert.Equal(t, 50.0, p.GeoPosition().Lat)
	assert.Equal(t, 3.0, p.GeoPosition().Lon)
	assert.Equal(t, 1, p.P2PFactor())
	assert.True(t, p.IsOwned())
	assert.Equal(t, BootstrapingStatus, p.Status())
}

/*
Scenario: Gets the peer's endpoint
	Given a peer
	When we want the peer's endpoint
	Then we gets the IP followed by the port
*/
func TestEndpoint(t *testing.T) {
	p := NewStartupPeer([]byte("key"), net.ParseIP("127.0.0.1"), 3000, "1.0", PeerPosition{Lat: 50.0, Lon: 3.0}, 1)
	assert.Equal(t, "127.0.0.1:3000", p.Endpoint())
}

/*
Scenario: Gets the peer's elasped hearbeats
	Given a created peer
	When we wait 2 seconds and we want the elapsed heartbeats
	Then we get the 2 heartbeats
*/
func TestElapsedHeartbeats(t *testing.T) {
	p := NewStartupPeer([]byte("key"), net.ParseIP("127.0.0.1"), 3000, "1.0", PeerPosition{Lat: 50.0, Lon: 3.0}, 1)
	time.Sleep(2 * time.Second)
	assert.Equal(t, int64(2), p.ElapsedHeartbeats())
}

/*
Scenario: Refreshes a peer
	Given a created peer
	When we refresh the peer
	Then the new info are stored
*/
func TestRefreshPeer(t *testing.T) {
	p := NewStartupPeer([]byte("key"), net.ParseIP("127.0.0.1"), 3000, "1.0", PeerPosition{Lat: 50.0, Lon: 3.0}, 1)
	p.Refresh(OkStatus, 600.10, "300.200.100", 50.0)
	assert.Equal(t, OkStatus, p.Status())
	assert.True(t, p.IsOk())
	assert.Equal(t, 600.10, p.FreeDiskSpace())
	assert.Equal(t, "300.200.100", p.CPULoad())
	assert.Equal(t, 50.0, p.IOWaitRate())
}

/*
Scenario: Returns default values
	Given a peer without some information
	When we want theses information
	Then we get the default values
*/
func TestDefaultGetters(t *testing.T) {
	p := NewPeerDetailed([]byte("key"), net.ParseIP("127.0.0.1"), 3000, time.Now(), false, nil)
	assert.Equal(t, 1, p.P2PFactor())
	assert.Equal(t, "0.0.0", p.CPULoad())
	assert.Equal(t, 0.0, p.IOWaitRate())
	assert.Equal(t, 0.0, p.FreeDiskSpace())
	assert.Equal(t, BootstrapingStatus, p.Status())
	assert.Equal(t, "1.0.0", p.Version())
}

/*
Scenario: Convers a seed into a peer
	Given a seed
	When we want to convert it to a peer
	Then we get a peer with the IP and the Port defined
*/
func TestSeedToPeer(t *testing.T) {
	s := Seed{IP: net.ParseIP("127.0.0.1"), Port: 3000}
	p := s.ToPeer()
	assert.NotNil(t, p)
	assert.Equal(t, "127.0.0.1", p.IP().String())
	assert.Equal(t, 3000, p.Port())
}

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
	p := NewStartupPeer(PublicKey("key"), net.ParseIP("127.0.0.1"), 3000, "1.0", PeerPosition{Lat: 50.0, Lon: 3.0})
	assert.NotNil(t, p)
	assert.Equal(t, "key", p.Identity().PublicKey().String())
	assert.Equal(t, "127.0.0.1", p.Identity().IP().String())
	assert.Equal(t, uint16(3000), p.Identity().Port())
	assert.Equal(t, "1.0", p.AppState().Version())
	assert.Equal(t, 50.0, p.AppState().GeoPosition().Lat)
	assert.Equal(t, 3.0, p.AppState().GeoPosition().Lon)
	assert.Equal(t, uint8(1), p.AppState().P2PFactor())
	assert.True(t, p.Owned())
	assert.Equal(t, BootstrapingStatus, p.AppState().Status())
}

/*
Scenario: Gets the peer's endpoint
	Given a peer
	When we want the peer's endpoint
	Then we gets the IP followed by the port
*/
func TestEndpoint(t *testing.T) {
	p := NewStartupPeer(PublicKey("key"), net.ParseIP("127.0.0.1"), 3000, "1.0", PeerPosition{Lat: 50.0, Lon: 3.0})
	assert.Equal(t, "127.0.0.1:3000", p.Endpoint())
}

/*
Scenario: Refreshes a peer
	Given a owned  peer
	When we refresh the peer
	Then the new info are stored
*/
func TestRefreshPeer(t *testing.T) {
	p := NewStartupPeer(PublicKey("key"), net.ParseIP("127.0.0.1"), 3000, "1.0", PeerPosition{Lat: 50.0, Lon: 3.0})
	err := p.Refresh(OkStatus, 600.10, "300.200.100", 50.0)
	assert.Nil(t, err)
	assert.Equal(t, OkStatus, p.AppState().Status())
	assert.Equal(t, 600.10, p.AppState().FreeDiskSpace())
	assert.Equal(t, "300.200.100", p.AppState().CPULoad())
}

/*
Scenario: Refreshes a peer
	Given a not owned  peer
	When we refresh the peer
	Then the new info are stored
*/
func TestRefreshNotOwnedPeer(t *testing.T) {
	p := peer{}
	err := p.Refresh(OkStatus, 600.10, "300.200.100", 50.0)
	assert.Error(t, err, ErrChangeNotOwnedPeer)
}

/*
Scenario: Create a discovered peer
	Given all information related to a peer (identity, heartbeat, app state)
	When we want theses information
	Then all the fields are setted and owned flag is false
*/
func TestCreateDiscoveredPeer(t *testing.T) {

	identity := peerIdentity{
		publicKey: PublicKey("key"),
		ip:        net.ParseIP("127.0.0.1"),
		port:      3000,
	}

	hbState := heartbeatState{
		generationTime:    time.Now(),
		elapsedHeartbeats: 1000,
	}

	appState := appState{

		freeDiskSpace: 200.10,
		status:        BootstrapingStatus,
		cpuLoad:       "300.10.200",
		version:       "1.0.0",
		p2pFactor:     2,
	}

	p := NewDiscoveredPeer(identity, hbState, appState)
	assert.Equal(t, uint64(1000), p.HeartbeatState().ElapsedHeartbeats())
	assert.Equal(t, "key", p.Identity().PublicKey().String())
	assert.Equal(t, "127.0.0.1", p.Identity().IP().String())
	assert.Equal(t, uint16(3000), p.Identity().Port())
	assert.Equal(t, uint8(2), p.AppState().P2PFactor())
	assert.Equal(t, "300.10.200", p.AppState().CPULoad())
	assert.Equal(t, 200.10, p.AppState().FreeDiskSpace())
	assert.Equal(t, BootstrapingStatus, p.AppState().Status())
	assert.Equal(t, "1.0.0", p.AppState().Version())
	assert.False(t, p.Owned())
}

/*
Scenario: Convers a seed into a peer
	Given a seed
	When we want to convert it to a peer
	Then we get a peer with the IP and the Port defined
*/
func TestSeedToPeer(t *testing.T) {
	s := Seed{IP: net.ParseIP("127.0.0.1"), Port: 3000}
	p := s.AsPeer()
	assert.NotNil(t, p)
	assert.Equal(t, "127.0.0.1", p.Identity().IP().String())
	assert.Equal(t, uint16(3000), p.Identity().Port())
}

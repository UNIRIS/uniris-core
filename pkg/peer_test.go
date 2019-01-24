package uniris

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

/*
Scenario: Create a local peer
	Given some inputs parameters
	When we create a peer that startup
	Then we get a new peer with the status bootstraping and specified as owned
*/
func TestNewPeer(t *testing.T) {
	p := NewLocalPeer("key", net.ParseIP("127.0.0.1"), 3000, "1.0", 3.0, 50.0)
	assert.NotNil(t, p)
	assert.Equal(t, "key", p.Identity().PublicKey())
	assert.Equal(t, "127.0.0.1", p.Identity().IP().String())
	assert.Equal(t, 3000, p.Identity().Port())
	assert.Equal(t, "1.0", p.AppState().Version())
	assert.Equal(t, 50.0, p.AppState().GeoPosition().Latitude())
	assert.Equal(t, 3.0, p.AppState().GeoPosition().Longitude())
	assert.Equal(t, 0, p.AppState().P2PFactor())
	assert.True(t, p.IsLocal())
	assert.Equal(t, BootstrapingPeer, p.AppState().Status())
}

/*
Scenario: Refreshes a peer
	Given a local  peer
	When we refresh the peer
	Then the new info are stored
*/
func TestRefreshPeer(t *testing.T) {
	p := NewLocalPeer("key", net.ParseIP("127.0.0.1"), 3000, "1.0", 3.0, 50.0)
	p.Refresh(OkPeerStatus, 600.10, "300.200.100", 50.0, 1)
	assert.Equal(t, OkPeerStatus, p.AppState().Status())
	assert.Equal(t, 600.10, p.AppState().FreeDiskSpace())
	assert.Equal(t, "300.200.100", p.AppState().CPULoad())
	assert.Equal(t, 1, p.AppState().DiscoveredPeersNumber())
}

/*
Scenario: Create a discovered peer
	Given all information related to a peer (identity, heartbeat, app state)
	When we want theses information
	Then all the fields are setted and owned flag is false
*/
func TestCreateDiscoveredPeer(t *testing.T) {

	identity := PeerIdentity{
		publicKey: "key",
		ip:        net.ParseIP("127.0.0.1"),
		port:      3000,
	}

	hbState := PeerHeartbeatState{
		generationTime:    time.Now(),
		elapsedHeartbeats: 1000,
	}

	appState := PeerAppState{

		freeDiskSpace:         200.10,
		status:                BootstrapingPeer,
		cpuLoad:               "300.10.200",
		version:               "1.0.0",
		p2pFactor:             2,
		discoveredPeersNumber: 10,
	}

	p := NewDiscoveredPeer(identity, hbState, appState)
	assert.Equal(t, int64(1000), p.HeartbeatState().ElapsedHeartbeats())
	assert.Equal(t, "key", p.Identity().PublicKey())
	assert.Equal(t, "127.0.0.1", p.Identity().IP().String())
	assert.Equal(t, 3000, p.Identity().Port())
	assert.Equal(t, 2, p.AppState().P2PFactor())
	assert.Equal(t, "300.10.200", p.AppState().CPULoad())
	assert.Equal(t, 200.10, p.AppState().FreeDiskSpace())
	assert.Equal(t, BootstrapingPeer, p.AppState().Status())
	assert.Equal(t, "1.0.0", p.AppState().Version())
	assert.Equal(t, 10, p.AppState().DiscoveredPeersNumber())
	assert.False(t, p.IsLocal())
}

/*
Scenario: Convers a seed into a peer
	Given a seed
	When we want to convert it to a peer
	Then we get a peer with the IP and the Port defined
*/
func TestSeedToPeer(t *testing.T) {
	s := Seed{
		PeerIdentity: PeerIdentity{
			ip: net.ParseIP("127.0.0.1"), port: 3000, publicKey: "key",
		},
	}
	p := s.AsPeer()
	assert.NotNil(t, p)
	assert.Equal(t, "127.0.0.1", p.Identity().IP().String())
	assert.Equal(t, 3000, p.Identity().Port())
	assert.Equal(t, "key", p.Identity().PublicKey())
}

/*
Scenario: Refreshes elapsed hearbeats
	Given an heartbeat state
	When we want to refresh the elapsed heartbeats
	Then we get the new elapsed heartbeats based on the current time
*/
func TestRefreshElapsedHeartbeats(t *testing.T) {
	hb := PeerHeartbeatState{
		generationTime: time.Now(),
	}
	time.Sleep(2 * time.Second)
	hb.refreshElapsedHeartbeats()
	assert.Equal(t, int64(2), hb.ElapsedHeartbeats())
}

/*
Scenario: Gets the elapsed heartbeats when no previous refresh
	Given a fresh heartbeat state
	When we get the elaspsed hearbeats
	Then we refresh the elapsed hearbeats and returns it
*/
func TestGetElapsedBeatsWithoutPrevRefresh(t *testing.T) {
	hb := PeerHeartbeatState{
		generationTime: time.Now(),
	}
	time.Sleep(2 * time.Second)
	assert.Equal(t, int64(2), hb.ElapsedHeartbeats())
}

/*
Scenario: Checks if an heartbeat state is more recent based on the upper generation time
	Given an heartbeat state with a generation time set as (now + 2 seconds)
	When we compare with another state with generation time set as now
	Then the first heartbeat is more recent
*/
func TestMoreRecentUpperGenTime(t *testing.T) {
	hb := PeerHeartbeatState{generationTime: time.Now().Add(2 * time.Second)}
	hb2 := PeerHeartbeatState{generationTime: time.Now()}
	assert.True(t, hb.MoreRecentThan(hb2))
}

/*
Scenario: Checks if an heartbeat state is more recent based on the same generation time
	Given a heartbeat state with a generation time as now
	When we compare with another with the same generation time
	Then the first heartbeat is not more recent
*/
func TestMoreRecentSameGenTimeSameElapsedBeats(t *testing.T) {
	hb := PeerHeartbeatState{generationTime: time.Now()}
	hb2 := PeerHeartbeatState{generationTime: time.Now()}
	assert.False(t, hb.MoreRecentThan(hb2))
}

/*
Scenario: Checks if an heartbeat state is more recent based on  the same generation and upper elapsed beats
	Given an heartbeat state with a generation time set as now and 500 elapsed beats
	When we compare with another with the same generation time and 300 elapsed beats
	Then the first heartbeat is more recent
*/
func TestMoreRecentSameGenTimeUpperElapsedBeats(t *testing.T) {
	hb := PeerHeartbeatState{generationTime: time.Now(), elapsedHeartbeats: 500}
	hb2 := PeerHeartbeatState{generationTime: time.Now(), elapsedHeartbeats: 300}
	assert.True(t, hb.MoreRecentThan(hb2))
}

/*
Scenario: Checks if an heartbeat state is more recent based on the same generation time and lower elapsed beats
	Given an heartbeat state with a generation time set as now and 300 elapsed beats
	When we compare with another with the same generation time and 500 elapsed beats
	Then the first heartbeat is not more recent
*/
func TestMoreRecentSameGenTimeLowerElapsedBeats(t *testing.T) {
	hb := PeerHeartbeatState{generationTime: time.Now(), elapsedHeartbeats: 300}
	hb2 := PeerHeartbeatState{generationTime: time.Now(), elapsedHeartbeats: 500}
	assert.False(t, hb.MoreRecentThan(hb2))
}

/*
Scenario: Checks if an heartbeat state is more recent based on a lower generation time
	Given an hearbeat with a generation time set as now
	When we compare with another with generation set as now + 2 seconds
	Then the first heartbeat is not more recent
*/
func TestMoreRecentLowerGenTime(t *testing.T) {
	hb := PeerHeartbeatState{generationTime: time.Now()}
	hb2 := PeerHeartbeatState{generationTime: time.Now().Add(2 * time.Second)}
	assert.False(t, hb.MoreRecentThan(hb2))
}

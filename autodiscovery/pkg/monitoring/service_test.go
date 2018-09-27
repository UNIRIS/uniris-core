package monitoring

import (
	"net"
	"testing"
	"time"

	"github.com/uniris/uniris-core/autodiscovery/pkg/mock"

	"github.com/stretchr/testify/assert"

	discovery "github.com/uniris/uniris-core/autodiscovery/pkg"
)

/*
Scenario: Gets peer status when no access to internet
	Given a peer without internet connection
	When we checks its status
	Then we get a faulty status
*/
func TestPeerStatusFaulty(t *testing.T) {
	srv := NewService(new(mock.Repository), new(mock.Monitor), new(mock.NetworkerInternetFails), new(mock.RobotWatcher))
	p := discovery.NewPeerDigest(
		discovery.NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, []byte("test")),
		discovery.NewPeerHeartbeatState(time.Now(), 0))
	status, _ := srv.PeerStatus(p)
	assert.Equal(t, discovery.FaultStatus, status)
}

/*
Scenario: Gets peer status when no NTP synchro
	Given a peer without NTP synchro
	When we checks its status
	Then we get a storage only status
*/
func TestPeerStatusStorageOnly(t *testing.T) {
	srv := NewService(new(mock.Repository), new(mock.Monitor), new(mock.NetworkerNTPFails), new(mock.RobotWatcher))
	p := discovery.NewPeerDigest(
		discovery.NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, []byte("test")),
		discovery.NewPeerHeartbeatState(time.Now(), 0))
	status, _ := srv.PeerStatus(p)
	assert.Equal(t, discovery.StorageOnlyStatus, status)
}

/*
Scenario: check refresh
	Given an initial seed
	When refresh
	Then status, CPUload, FreeDiskSpace and IOWaitRate are updated
*/
func TestRefresh(t *testing.T) {
	repo := new(mock.Repository)

	seed1 := discovery.Seed{IP: net.ParseIP("10.1.1.1"), Port: 3000}
	seed2 := discovery.Seed{IP: net.ParseIP("10.1.1.2"), Port: 3001}
	seed3 := discovery.Seed{IP: net.ParseIP("10.1.1.3"), Port: 3002}
	repo.AddSeed(seed1)
	repo.AddSeed(seed2)
	repo.AddSeed(seed3)

	st2 := discovery.NewPeerAppState("0.0", discovery.OkStatus, discovery.PeerPosition{}, "0.0.0", 0.0, 0, 5)
	st3 := discovery.NewPeerAppState("0.0", discovery.OkStatus, discovery.PeerPosition{}, "0.0.0", 0.0, 0, 5)

	p1 := discovery.NewStartupPeer([]byte("key1"), seed1.IP, seed1.Port, "1.0", discovery.PeerPosition{})
	p1.Refresh(discovery.BootstrapingStatus, 0.0, "", 1, 5)

	p2 := discovery.NewDiscoveredPeer(
		discovery.NewPeerIdentity(seed2.IP, seed2.Port, []byte("key2")),
		discovery.NewPeerHeartbeatState(time.Now(), 0),
		st2,
	)

	p3 := discovery.NewDiscoveredPeer(
		discovery.NewPeerIdentity(seed3.IP, seed3.Port, []byte("key3")),
		discovery.NewPeerHeartbeatState(time.Now(), 0),
		st3,
	)

	repo.AddPeer(p1)
	repo.AddPeer(p2)
	repo.AddPeer(p3)

	srv := NewService(repo, new(mock.Monitor), new(mock.Networker), new(mock.RobotWatcher))
	err := srv.RefreshPeer(p1)
	assert.Nil(t, err)

	p, _ := repo.GetPeerByIP(seed1.IP)
	assert.Equal(t, "0.62 0.77 0.71 4/972 26361", p.AppState().CPULoad())
	assert.Equal(t, discovery.OkStatus, p.AppState().Status())
	assert.Equal(t, float64(212383852), p.AppState().FreeDiskSpace())
	assert.Equal(t, 3, p.AppState().DiscoveredPeersNumber())
	assert.Equal(t, 1, p.AppState().P2PFactor())
}

/*
Scenario: Gets peer status
	Given a peer with 3 seed (discoveredPeersNumber=5 for all seed) / 5 peers on the repo
	When DiscoveredPeer=5 and elapsedheartbeats < Bootstrapingmintime
	Then state is OkStatus
*/
func TestPeerStatusOkStatus(t *testing.T) {
	repo := new(mock.Repository)
	srv := NewService(repo, new(mock.Monitor), new(mock.Networker), new(mock.RobotWatcher))

	initP := discovery.NewStartupPeer([]byte("key"), net.ParseIP("127.0.0.1"), 3000, "0.0", discovery.PeerPosition{})
	repo.AddPeer(initP)

	seed1 := discovery.Seed{IP: net.ParseIP("10.1.1.1"), Port: 3000}
	seed2 := discovery.Seed{IP: net.ParseIP("10.1.1.2"), Port: 3001}
	seed3 := discovery.Seed{IP: net.ParseIP("10.1.1.3"), Port: 3002}

	repo.AddSeed(seed1)
	repo.AddSeed(seed2)
	repo.AddSeed(seed3)

	assert.Equal(t, 3, len(repo.Seeds))

	st1 := discovery.NewPeerAppState("0.0", discovery.OkStatus, discovery.PeerPosition{}, "0.0.0", 0.0, 0, 5)
	st2 := discovery.NewPeerAppState("0.0", discovery.OkStatus, discovery.PeerPosition{}, "0.0.0", 0.0, 0, 5)
	st3 := discovery.NewPeerAppState("0.0", discovery.OkStatus, discovery.PeerPosition{}, "0.0.0", 0.0, 0, 5)

	p1 := discovery.NewDiscoveredPeer(
		discovery.NewPeerIdentity(seed1.IP, seed1.Port, []byte("key1")),
		discovery.NewPeerHeartbeatState(time.Now(), 0),
		st1,
	)

	p2 := discovery.NewDiscoveredPeer(
		discovery.NewPeerIdentity(seed2.IP, seed2.Port, []byte("key2")),
		discovery.NewPeerHeartbeatState(time.Now(), 0),
		st2,
	)

	p3 := discovery.NewDiscoveredPeer(
		discovery.NewPeerIdentity(seed3.IP, seed3.Port, []byte("key3")),
		discovery.NewPeerHeartbeatState(time.Now(), 0),
		st3,
	)
	repo.AddPeer(p1)
	repo.AddPeer(p2)
	repo.AddPeer(p3)

	p4 := discovery.NewDiscoveredPeer(
		discovery.NewPeerIdentity(net.ParseIP("185.123.4.9"), 4000, []byte("key4")),
		discovery.NewPeerHeartbeatState(time.Now(), 0),
		st1)

	repo.AddPeer(p4)

	assert.Equal(t, 5, len(repo.Peers))

	selfpeer, err := repo.GetOwnedPeer()
	selfpeer.Refresh(discovery.BootstrapingStatus, 0.0, "0.0.0", 5, 5)

	s, err := srv.PeerStatus(selfpeer)
	assert.Equal(t, nil, err)
	assert.Equal(t, discovery.OkStatus, s)
}

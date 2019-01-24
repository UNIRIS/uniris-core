package gossip

import (
	"errors"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	uniris "github.com/uniris/uniris-core/pkg"
)

/*
Scenario: Get peer status when no peer has been discovered yet by it
	Given a peer without discoveries
	When I want get its status
	Then I get a boostraping status
*/
func TestPeerStatusWithNoDiscoveries(t *testing.T) {
	p := uniris.NewPeerDigest(
		uniris.NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "test"),
		uniris.NewPeerHeartbeatState(time.Now(), 0))

	status, err := getPeerStatus(p, 0, mockPeerNetworker{})
	assert.Nil(t, err)
	assert.Equal(t, uniris.BootstrapingPeer, status)
}

/*
Scenario: Gets peer status when no access to internet
	Given a peer without internet connection
	When we checks its status
	Then we get a faulty status
*/
func TestPeerStatusWithNotInternet(t *testing.T) {
	p := uniris.NewPeerDigest(
		uniris.NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "test"),
		uniris.NewPeerHeartbeatState(time.Now(), 0))
	status, _ := getPeerStatus(p, 1, mockFailInternetNetworker{})
	assert.Equal(t, uniris.FaultyPeer, status)
}

/*
Scenario: Gets peer status when no NTP synchro
	Given a peer without NTP synchro
	When we checks its status
	Then we get a storage only status
*/
func TestPeerStatusWithBadNTP(t *testing.T) {
	p := uniris.NewPeerDigest(
		uniris.NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "test"),
		uniris.NewPeerHeartbeatState(time.Now(), 0))
	status, _ := getPeerStatus(p, 1, mockFailNTPNetworker{})
	assert.Equal(t, uniris.StorageOnlyPeer, status)
}

/*
Scenario: Gets peer status when the elapsed time lower than the bootstraping time
	Given a peer just starting
	When I want get its status
	Then I get a bootstraping status
*/
func TestPeerStatusWithElapsedTimeLowerBootstrapingTime(t *testing.T) {
	p := uniris.NewPeerDigest(
		uniris.NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "test"),
		uniris.NewPeerHeartbeatState(time.Now(), 0))

	status, _ := getPeerStatus(p, 1, mockPeerNetworker{})
	assert.Equal(t, uniris.BootstrapingPeer, status)
}

/*
Scenario: Gets peer status when the average of discoveries is greater than the peer discovery
	Given a peer just starting with a avergage of discoveries greater than the peer discovery
	When I want get its status
	Then I get a bootstraping status
*/
func TestPeerStatusWithAvgDiscoveriesGreaterThanPeerDiscovery(t *testing.T) {
	p := uniris.NewDiscoveredPeer(
		uniris.NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "test"),
		uniris.NewPeerHeartbeatState(time.Now(), 0),
		uniris.NewPeerAppState("1.0", uniris.BootstrapingPeer, 30.0, 10.0, "", 100, 1, 1))

	status, _ := getPeerStatus(p, 5, mockPeerNetworker{})
	assert.Equal(t, uniris.BootstrapingPeer, status)
}

/*
Scenario: Gets peer status when the peer time it less than the boostraping and avg of discoveries is less the discovery number
	Given a peer just starting with a avergage of discoveries less its discoveries
	When I want get its status
	Then I get a OK status
*/
func TestPeerStatusWithAvgDiscoveriesLessThanPeerDiscovery(t *testing.T) {
	p := uniris.NewDiscoveredPeer(
		uniris.NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "test"),
		uniris.NewPeerHeartbeatState(time.Now(), 0),
		uniris.NewPeerAppState("1.0", uniris.BootstrapingPeer, 30.0, 10.0, "", 100, 1, 5))

	status, _ := getPeerStatus(p, 1, mockPeerNetworker{})
	assert.Equal(t, uniris.OkPeerStatus, status)
}

/*
Scenario: Gets peer status a peer live longer than the bootstraping time
	Given a peer started for a while
	When I want get its status
	Then I get a OK status
*/
func TestPeerStatusWithLongTTL(t *testing.T) {
	p := uniris.NewDiscoveredPeer(
		uniris.NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "test"),
		uniris.NewPeerHeartbeatState(time.Now(), 5000),
		uniris.NewPeerAppState("1.0", uniris.BootstrapingPeer, 30.0, 10.0, "", 100, 1, 5))

	status, _ := getPeerStatus(p, 1, mockPeerNetworker{})
	assert.Equal(t, uniris.OkPeerStatus, status)
}

func TestGetP2PFactor(t *testing.T) {

	peers := []uniris.Peer{}

	assert.Equal(t, 1, getP2PFactor(peers))
}

/*
Scenario: Avergage number of discoveries with 3 seeds as the only known peers
	Given 3 seeds, seed1 discoveredPeersNumber = 5,seed2 discoveredPeersNumber = 6, seed3 discoveredPeersNumber = 7
	When I want to the retrieve the avergage of number discovered peers
	Then I get 6
*/
func TestAvergageDiscoveriesWithOnlySeeds(t *testing.T) {

	seed1 := uniris.Seed{
		PeerIdentity: uniris.NewPeerIdentity(net.ParseIP("10.0.0.1"), 3001, "key1"),
	}
	seed2 := uniris.Seed{
		PeerIdentity: uniris.NewPeerIdentity(net.ParseIP("10.0.0.2"), 3002, "key2"),
	}
	seed3 := uniris.Seed{
		PeerIdentity: uniris.NewPeerIdentity(net.ParseIP("10.0.0.3"), 3003, "key3"),
	}

	p1 := uniris.NewDiscoveredPeer(
		seed1.PeerIdentity,
		uniris.NewPeerHeartbeatState(time.Now(), 0),
		uniris.NewPeerAppState("0.0", uniris.OkPeerStatus, 30.0, 12.0, "0.0.0", 0.0, 0, 5),
	)

	p2 := uniris.NewDiscoveredPeer(
		seed2.PeerIdentity,
		uniris.NewPeerHeartbeatState(time.Now(), 0),
		uniris.NewPeerAppState("0.0", uniris.OkPeerStatus, 20.0, 5.0, "0.0.0", 0.0, 0, 6),
	)

	p3 := uniris.NewDiscoveredPeer(
		seed3.PeerIdentity,
		uniris.NewPeerHeartbeatState(time.Now(), 0),
		uniris.NewPeerAppState("0.0", uniris.OkPeerStatus, 10.0, 3.0, "0.0.0", 0.0, 0, 7),
	)

	avg := getSeedDiscoveryAverage([]uniris.Seed{seed1, seed2, seed3}, []uniris.Peer{p1, p2, p3})
	assert.Equal(t, 6, avg)
}

/*
Scenario: Avergage discoveries peers including seeds and discovered peers
	Given 3 seeds and a discovered peer including for each a number of discovered peers equal to 5
	When I want to the retrieve the avergage of number discovered peers
	Then I get 5
*/
func TestAvergageDiscoveriesWithSeedAndDiscoveries(t *testing.T) {

	seed1 := uniris.Seed{
		PeerIdentity: uniris.NewPeerIdentity(net.ParseIP("10.0.0.1"), 3001, "key1"),
	}
	seed2 := uniris.Seed{
		PeerIdentity: uniris.NewPeerIdentity(net.ParseIP("10.0.0.2"), 3002, "key2"),
	}
	seed3 := uniris.Seed{
		PeerIdentity: uniris.NewPeerIdentity(net.ParseIP("10.0.0.3"), 3003, "key3"),
	}

	p1 := uniris.NewDiscoveredPeer(
		seed1.PeerIdentity,
		uniris.NewPeerHeartbeatState(time.Now(), 0),
		uniris.NewPeerAppState("0.0", uniris.OkPeerStatus, 30.0, 12.0, "0.0.0", 0.0, 0, 5),
	)

	p2 := uniris.NewDiscoveredPeer(
		seed2.PeerIdentity,
		uniris.NewPeerHeartbeatState(time.Now(), 0),
		uniris.NewPeerAppState("0.0", uniris.OkPeerStatus, 20.0, 5.0, "0.0.0", 0.0, 0, 5),
	)

	p3 := uniris.NewDiscoveredPeer(
		seed3.PeerIdentity,
		uniris.NewPeerHeartbeatState(time.Now(), 0),
		uniris.NewPeerAppState("0.0", uniris.OkPeerStatus, 10.0, 3.0, "0.0.0", 0.0, 0, 5),
	)

	p4 := uniris.NewDiscoveredPeer(
		uniris.NewPeerIdentity(net.ParseIP("185.123.4.9"), 4000, "key4"),
		uniris.NewPeerHeartbeatState(time.Now(), 0),
		uniris.NewPeerAppState("0.0", uniris.OkPeerStatus, 40.0, 3.0, "0.0.0", 0.0, 0, 5),
	)

	avg := getSeedDiscoveryAverage([]uniris.Seed{seed1, seed2, seed3}, []uniris.Peer{p1, p2, p3, p4})
	assert.Equal(t, 5, avg)
}

type mockFailInternetNetworker struct{}

func (pn mockFailInternetNetworker) CheckNtpState() error {
	return nil
}

func (pn mockFailInternetNetworker) CheckInternetState() error {
	return errors.New("Unexpected")
}

type mockFailNTPNetworker struct{}

func (pn mockFailNTPNetworker) CheckNtpState() error {
	return ErrNTPFailure
}

func (pn mockFailNTPNetworker) CheckInternetState() error {
	return nil
}

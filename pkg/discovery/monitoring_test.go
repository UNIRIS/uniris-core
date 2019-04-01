package discovery

import (
	"errors"
	"log"
	"net"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/uniris/uniris-core/pkg/logging"
)

/*
Scenario: Get peer status when no peer has been discovered yet by it
	Given a peer without discoveries
	When I want get its status
	Then I get a boostraping status
*/
func TestPeerStatusWithNoDiscoveries(t *testing.T) {
	resettimerState()
	p := NewPeerDigest(
		NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "test"),
		NewPeerHeartbeatState(time.Now(), 0))
	l := logging.NewLogger("stdout", log.New(os.Stdout, "", 0), "test", net.ParseIP("127.0.0.1"), logging.ErrorLogLevel)
	status, err := localStatus(p, 0, mockNetworkChecker{}, l)
	assert.Nil(t, err)
	assert.Equal(t, BootstrapingPeer, status)
}

/*
Scenario: Gets peer status when no access to internet
	Given a peer without internet connection
	When we checks its status
	Then we get a faulty status
*/
func TestPeerStatusWithNotInternet(t *testing.T) {
	resettimerState()
	p := NewPeerDigest(
		NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "test"),
		NewPeerHeartbeatState(time.Now(), 0))
	l := logging.NewLogger("stdout", log.New(os.Stdout, "", 0), "test", net.ParseIP("127.0.0.1"), logging.ErrorLogLevel)
	status, _ := localStatus(p, 1, mockFailInternetNetworker{}, l)
	assert.Equal(t, FaultyPeer, status)
}

/*
Scenario: Gets peer status when no access to internet and get the state after the internet problem is resolved
	Given a peer without internet connection
	When we checks its status
	Then we get the good state depending on the timerState
*/
func TestPeerStatusWithNotInternetWithTimerStateEnabled(t *testing.T) {
	BootstrapingMinTime = 5
	p := NewPeerDigest(
		NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "test"),
		NewPeerHeartbeatState(time.Now(), 0))

	l := logging.NewLogger("stdout", log.New(os.Stdout, "", 0), "test", net.ParseIP("127.0.0.1"), logging.ErrorLogLevel)
	status, _ := localStatus(p, 1, mockFailInternetNetworker{}, l)
	assert.Equal(t, FaultyPeer, status)
	ts := checktimerState()
	assert.Equal(t, false, ts)
	status2, _ := localStatus(p, 1, mockNetworkChecker{}, l)
	assert.Equal(t, BootstrapingPeer, status2)
	time.Sleep(5 * time.Second)
	status3, _ := localStatus(p, 1, mockNetworkChecker{}, l)
	assert.Equal(t, OkPeerStatus, status3)

}

/*
Scenario: Gets peer status when no NTP synchro
	Given a peer without NTP synchro
	When we checks its status
	Then we get a faulty status
*/
func TestPeerStatusWithBadNTP(t *testing.T) {
	resettimerState()
	p := NewPeerDigest(
		NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "test"),
		NewPeerHeartbeatState(time.Now(), 0))

	l := logging.NewLogger("stdout", log.New(os.Stdout, "", 0), "test", net.ParseIP("127.0.0.1"), logging.ErrorLogLevel)
	status, _ := localStatus(p, 1, mockFailNTPNetworker{}, l)
	assert.Equal(t, FaultyPeer, status)
}

/*
Scenario: Gets peer status when no GRPC are reached
	Given a peer without GRPC server running
	When we checks its status
	Then we get a faulty status
*/
func TestPeerStatusWithNoGRPC(t *testing.T) {
	resettimerState()
	p := NewPeerDigest(
		NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "test"),
		NewPeerHeartbeatState(time.Now(), 0))

	l := logging.NewLogger("stdout", log.New(os.Stdout, "", 0), "test", net.ParseIP("127.0.0.1"), logging.ErrorLogLevel)
	status, _ := localStatus(p, 1, mockFailGRPCServersNetworker{}, l)
	assert.Equal(t, FaultyPeer, status)
}

/*
Scenario: Gets peer status when the elapsed time lower than the bootstraping time
	Given a peer just starting
	When I want get its status
	Then I get a bootstraping status
*/
func TestPeerStatusWithElapsedTimeLowerBootstrapingTime(t *testing.T) {
	resettimerState()
	p := NewPeerDigest(
		NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "test"),
		NewPeerHeartbeatState(time.Now(), 0))

	l := logging.NewLogger("stdout", log.New(os.Stdout, "", 0), "test", net.ParseIP("127.0.0.1"), logging.ErrorLogLevel)
	status, _ := localStatus(p, 1, mockNetworkChecker{}, l)
	assert.Equal(t, BootstrapingPeer, status)
}

/*
Scenario: Gets peer status when the average of discoveries is greater than the peer discovery
	Given a peer just starting with a avergage of discoveries greater than the peer discovery
	When I want get its status
	Then I get a bootstraping status
*/
func TestPeerStatusWithAvgDiscoveriesGreaterThanPeerDiscovery(t *testing.T) {
	resettimerState()
	p := NewDiscoveredPeer(
		NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "test"),
		NewPeerHeartbeatState(time.Now(), 0),
		NewPeerAppState("1.0", BootstrapingPeer, 30.0, 10.0, "", 100, 1, 1))

	l := logging.NewLogger("stdout", log.New(os.Stdout, "", 0), "test", net.ParseIP("127.0.0.1"), logging.ErrorLogLevel)
	status, _ := localStatus(p, 5, mockNetworkChecker{}, l)
	assert.Equal(t, BootstrapingPeer, status)
}

/*
Scenario: Gets peer status when the peer time it less than the boostraping and avg of discoveries is less the discovery number
	Given a peer just starting with a avergage of discoveries less its discoveries
	When I want get its status
	Then I get a OK status
*/
func TestPeerStatusWithAvgDiscoveriesLessThanPeerDiscovery(t *testing.T) {
	resettimerState()
	p := NewDiscoveredPeer(
		NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "test"),
		NewPeerHeartbeatState(time.Now(), 0),
		NewPeerAppState("1.0", BootstrapingPeer, 30.0, 10.0, "", 100, 1, 5))

	l := logging.NewLogger("stdout", log.New(os.Stdout, "", 0), "test", net.ParseIP("127.0.0.1"), logging.ErrorLogLevel)
	status, _ := localStatus(p, 1, mockNetworkChecker{}, l)
	assert.Equal(t, OkPeerStatus, status)
}

/*
Scenario: Gets peer status a peer live longer than the bootstraping time
	Given a peer started for a while
	When I want get its status
	Then I get a OK status
*/
func TestPeerStatusWithLongTTL(t *testing.T) {
	resettimerState()
	p := NewDiscoveredPeer(
		NewPeerIdentity(net.ParseIP("127.0.0.1"), 3000, "test"),
		NewPeerHeartbeatState(time.Now(), 5000),
		NewPeerAppState("1.0", BootstrapingPeer, 30.0, 10.0, "", 100, 1, 5))

	l := logging.NewLogger("stdout", log.New(os.Stdout, "", 0), "test", net.ParseIP("127.0.0.1"), logging.ErrorLogLevel)
	status, _ := localStatus(p, 1, mockNetworkChecker{}, l)
	assert.Equal(t, OkPeerStatus, status)
}

func TestGetP2PFactor(t *testing.T) {

	peers := []Peer{}

	assert.Equal(t, 1, p2pFactor(peers))
}

/*
Scenario: Avergage number of discoveries with 3 seeds as the only known peers
	Given 3 seeds, seed1 discoveredPeersNumber = 5,seed2 discoveredPeersNumber = 6, seed3 discoveredPeersNumber = 7
	When I want to the retrieve the avergage of number discovered peers
	Then I get 6
*/
func TestAvergageReachableWithOnlySeeds(t *testing.T) {

	seed1 := NewPeerIdentity(net.ParseIP("10.0.0.1"), 3001, "key1")
	seed2 := NewPeerIdentity(net.ParseIP("10.0.0.2"), 3002, "key2")
	seed3 := NewPeerIdentity(net.ParseIP("10.0.0.3"), 3003, "key3")

	p1 := NewDiscoveredPeer(
		seed1,
		NewPeerHeartbeatState(time.Now(), 0),
		NewPeerAppState("0.0", OkPeerStatus, 30.0, 12.0, "0.0.0", 0.0, 0, 5),
	)

	p2 := NewDiscoveredPeer(
		seed2,
		NewPeerHeartbeatState(time.Now(), 0),
		NewPeerAppState("0.0", OkPeerStatus, 20.0, 5.0, "0.0.0", 0.0, 0, 6),
	)

	p3 := NewDiscoveredPeer(
		seed3,
		NewPeerHeartbeatState(time.Now(), 0),
		NewPeerAppState("0.0", OkPeerStatus, 10.0, 3.0, "0.0.0", 0.0, 0, 7),
	)

	avg := seedReachableAverage([]PeerIdentity{seed1, seed2, seed3}, []Peer{p1, p2, p3})
	assert.Equal(t, 6, avg)
}

/*
Scenario: Avergage discoveries peers including seeds and discovered peers
	Given 3 seeds and a discovered peer including for each a number of discovered peers equal to 5
	When I want to the retrieve the avergage of number discovered peers
	Then I get 5
*/
func TestAvergageDiscoveriesWithSeedAndDiscoveries(t *testing.T) {

	seed1 := NewPeerIdentity(net.ParseIP("10.0.0.1"), 3001, "key1")
	seed2 := NewPeerIdentity(net.ParseIP("10.0.0.2"), 3002, "key2")
	seed3 := NewPeerIdentity(net.ParseIP("10.0.0.3"), 3003, "key3")

	p1 := NewDiscoveredPeer(
		seed1,
		NewPeerHeartbeatState(time.Now(), 0),
		NewPeerAppState("0.0", OkPeerStatus, 30.0, 12.0, "0.0.0", 0.0, 0, 5),
	)

	p2 := NewDiscoveredPeer(
		seed2,
		NewPeerHeartbeatState(time.Now(), 0),
		NewPeerAppState("0.0", OkPeerStatus, 20.0, 5.0, "0.0.0", 0.0, 0, 5),
	)

	p3 := NewDiscoveredPeer(
		seed3,
		NewPeerHeartbeatState(time.Now(), 0),
		NewPeerAppState("0.0", OkPeerStatus, 10.0, 3.0, "0.0.0", 0.0, 0, 5),
	)

	p4 := NewDiscoveredPeer(
		NewPeerIdentity(net.ParseIP("185.123.4.9"), 4000, "key4"),
		NewPeerHeartbeatState(time.Now(), 0),
		NewPeerAppState("0.0", OkPeerStatus, 40.0, 3.0, "0.0.0", 0.0, 0, 5),
	)

	avg := seedReachableAverage([]PeerIdentity{seed1, seed2, seed3}, []Peer{p1, p2, p3, p4})
	assert.Equal(t, 5, avg)
}

type mockFailInternetNetworker struct{}

func (n mockFailInternetNetworker) CheckNtpState() error {
	return nil
}

func (n mockFailInternetNetworker) CheckInternetState() error {
	return errors.New("Unexpected")
}

func (n mockFailInternetNetworker) CheckGRPCServer() error {
	return nil
}

type mockFailNTPNetworker struct{}

func (n mockFailNTPNetworker) CheckNtpState() error {
	return ErrNTPFailure
}

func (n mockFailNTPNetworker) CheckInternetState() error {
	return nil
}

func (n mockFailNTPNetworker) CheckGRPCServer() error {
	return nil
}

type mockFailGRPCServersNetworker struct{}

func (n mockFailGRPCServersNetworker) CheckNtpState() error {
	return nil
}

func (n mockFailGRPCServersNetworker) CheckInternetState() error {
	return nil
}

func (n mockFailGRPCServersNetworker) CheckGRPCServer() error {
	return ErrGRPCServer
}

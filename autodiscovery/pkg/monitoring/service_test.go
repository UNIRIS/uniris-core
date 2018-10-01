package monitoring

import (
	"encoding/hex"
	"errors"
	"net"
	"testing"
	"time"

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
	srv := NewService(new(mockRepository), new(monitor), new(NetworkerInternetFails), new(robotWatcher))
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
	srv := NewService(new(mockRepository), new(monitor), new(NetworkerNTPFails), new(robotWatcher))
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
	repo := new(mockRepository)

	seed1 := discovery.Seed{IP: net.ParseIP("10.1.1.1"), Port: 3000}
	seed2 := discovery.Seed{IP: net.ParseIP("10.1.1.2"), Port: 3001}
	seed3 := discovery.Seed{IP: net.ParseIP("10.1.1.3"), Port: 3002}
	repo.SetSeed(seed1)
	repo.SetSeed(seed2)
	repo.SetSeed(seed3)

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

	repo.SetPeer(p1)
	repo.SetPeer(p2)
	repo.SetPeer(p3)

	srv := NewService(repo, new(monitor), new(Networker), new(robotWatcher))
	err := srv.RefreshPeer(p1)
	assert.Nil(t, err)

	p, _ := repo.GetPeerByIP(seed1.IP)
	assert.Equal(t, "0.62 0.77 0.71 4/972 26361", p.AppState().CPULoad())
	assert.Equal(t, discovery.OkStatus, p.AppState().Status())
	assert.Equal(t, float64(212383852), p.AppState().FreeDiskSpace())
	assert.Equal(t, 2, p.AppState().DiscoveredPeersNumber())
	assert.Equal(t, 1, p.AppState().P2PFactor())
}

/*
Scenario: Gets peer status
	Given a peer with 3 seed (discoveredPeersNumber=5 for all seed) / 5 peers on the repo
	When DiscoveredPeer=5 and elapsedheartbeats < Bootstrapingmintime
	Then state is OkStatus
*/
func TestPeerStatusOkStatus(t *testing.T) {
	repo := new(mockRepository)
	srv := NewService(repo, new(monitor), new(Networker), new(robotWatcher))

	initP := discovery.NewStartupPeer([]byte("key"), net.ParseIP("127.0.0.1"), 3000, "0.0", discovery.PeerPosition{})
	repo.SetPeer(initP)

	seed1 := discovery.Seed{IP: net.ParseIP("10.1.1.1"), Port: 3000}
	seed2 := discovery.Seed{IP: net.ParseIP("10.1.1.2"), Port: 3001}
	seed3 := discovery.Seed{IP: net.ParseIP("10.1.1.3"), Port: 3002}

	repo.SetSeed(seed1)
	repo.SetSeed(seed2)
	repo.SetSeed(seed3)

	seeds, _ := repo.ListSeedPeers()
	assert.Equal(t, 3, len(seeds))

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
	repo.SetPeer(p1)
	repo.SetPeer(p2)
	repo.SetPeer(p3)

	p4 := discovery.NewDiscoveredPeer(
		discovery.NewPeerIdentity(net.ParseIP("185.123.4.9"), 4000, []byte("key4")),
		discovery.NewPeerHeartbeatState(time.Now(), 0),
		st1)

	repo.SetPeer(p4)

	peers, _ := repo.ListDiscoveredPeers()
	assert.Equal(t, 4, len(peers))

	selfpeer, err := repo.GetOwnedPeer()
	selfpeer.Refresh(discovery.BootstrapingStatus, 0.0, "0.0.0", 5, 5)

	s, err := srv.PeerStatus(selfpeer)
	assert.Equal(t, nil, err)
	assert.Equal(t, discovery.OkStatus, s)
}

//////////////////////////////////////////////////////////
// 						MOCKS
/////////////////////////////////////////////////////////

type monitor struct{}

func (w monitor) CPULoad() (string, error) {
	return "0.62 0.77 0.71 4/972 26361", nil
}

func (w monitor) FreeDiskSpace() (float64, error) {
	return 212383852, nil
}

func (w monitor) P2PFactor() (int, error) {
	return 1, nil
}

type mockRepository struct {
	ownedPeer       discovery.Peer
	discoveredPeers []discovery.Peer
	seedPeers       []discovery.Seed
}

func (r *mockRepository) CountDiscoveredPeers() (int, error) {
	return len(r.discoveredPeers), nil
}

//GetOwnedPeer return the local peer
func (r *mockRepository) GetOwnedPeer() (discovery.Peer, error) {
	return r.ownedPeer, nil
}

//ListSeedPeers return all the seed on the mockRepository
func (r *mockRepository) ListSeedPeers() ([]discovery.Seed, error) {
	return r.seedPeers, nil
}

//ListDiscoveredPeers returns all the discoveredPeers on the mockRepository
func (r *mockRepository) ListDiscoveredPeers() ([]discovery.Peer, error) {
	return r.discoveredPeers, nil
}

func (r *mockRepository) SetPeer(peer discovery.Peer) error {
	if peer.Owned() {
		r.ownedPeer = peer
		return nil
	}
	if r.containsPeer(peer) {
		for _, p := range r.discoveredPeers {
			if p.Identity().PublicKey().Equals(peer.Identity().PublicKey()) {
				p = peer
				break
			}
		}
	} else {
		r.discoveredPeers = append(r.discoveredPeers, peer)
	}
	return nil
}

func (r *mockRepository) SetSeed(s discovery.Seed) error {
	r.seedPeers = append(r.seedPeers, s)
	return nil
}

//GetPeerByIP get a peer from the mockRepository using its ip
func (r *mockRepository) GetPeerByIP(ip net.IP) (p discovery.Peer, err error) {
	if r.ownedPeer.Identity().IP().Equal(ip) {
		return r.ownedPeer, nil
	}
	for i := 0; i < len(r.discoveredPeers); i++ {
		if r.discoveredPeers[i].Identity().IP().Equal(ip) {
			return r.discoveredPeers[i], nil
		}
	}
	return
}

func (r *mockRepository) containsPeer(p discovery.Peer) bool {
	mdiscoveredPeers := make(map[string]discovery.Peer, 0)
	for _, p := range r.discoveredPeers {
		mdiscoveredPeers[hex.EncodeToString(p.Identity().PublicKey())] = p
	}

	_, exist := mdiscoveredPeers[hex.EncodeToString(p.Identity().PublicKey())]
	return exist
}

type Networker struct{}

func (n Networker) IP() (net.IP, error) {
	return net.ParseIP("127.0.0.1"), nil
}

func (n Networker) CheckInternetState() error {
	return nil
}

func (n Networker) CheckNtpState() error {
	return nil
}

type NetworkerNTPFails struct{}

func (n NetworkerNTPFails) IP() (net.IP, error) {
	return net.ParseIP("127.0.0.1"), nil
}

func (n NetworkerNTPFails) CheckInternetState() error {
	return nil
}

func (n NetworkerNTPFails) CheckNtpState() error {
	return errors.New("System Clock have a big Offset check the ntp configuration of the system")
}

type NetworkerInternetFails struct{}

func (n NetworkerInternetFails) IP() (net.IP, error) {
	return net.ParseIP("127.0.0.1"), nil
}

func (n NetworkerInternetFails) CheckInternetState() error {
	return errors.New("required processes are not running")
}

func (n NetworkerInternetFails) CheckNtpState() error {
	return nil
}

type robotWatcher struct{}

func (r robotWatcher) CheckAutodiscoveryProcess(port int) error {
	return nil
}

func (r robotWatcher) CheckDataProcess() error {
	return nil
}

func (r robotWatcher) CheckMiningProcess() error {
	return nil
}

func (r robotWatcher) CheckAIProcess() error {
	return nil
}

func (r robotWatcher) CheckScyllaDbProcess() error {
	return nil
}

func (r robotWatcher) CheckRedisProcess() error {
	return nil
}

func (r robotWatcher) CheckRabbitmqProcess() error {
	return nil
}

type mockSystemNetworker2 struct{}

func (n mockSystemNetworker2) IP() (net.IP, error) {
	return net.ParseIP("127.0.0.1"), nil
}

func (n mockSystemNetworker2) CheckInternetState() error {
	return nil
}

func (n mockSystemNetworker2) CheckNtpState() error {
	return errors.New("System Clock have a big Offset check the ntp configuration of the system")
}

type mockSystemNetworker3 struct{}

func (n mockSystemNetworker3) IP() (net.IP, error) {
	return net.ParseIP("127.0.0.1"), nil
}

func (n mockSystemNetworker3) CheckInternetState() error {
	return errors.New("required processes are not running")
}

func (n mockSystemNetworker3) CheckNtpState() error {
	return nil
}

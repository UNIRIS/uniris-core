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
Scenario: check refresh
	Given an initial seed
	When refresh
	Then status, CPUload, FreeDiskSpace and IOWaitRate are updated
*/

func TestRefresh(t *testing.T) {
	repo := new(mockPeerRepository)

	seed1 := discovery.Seed{IP: net.ParseIP("10.1.1.1"), Port: 3000}
	seed2 := discovery.Seed{IP: net.ParseIP("10.1.1.2"), Port: 3001}
	seed3 := discovery.Seed{IP: net.ParseIP("10.1.1.3"), Port: 3002}
	repo.AddSeed(seed1)
	repo.AddSeed(seed2)
	repo.AddSeed(seed3)

	st1 := discovery.NewState("0.0", discovery.OkStatus, discovery.PeerPosition{}, "0.0.0", 0.0, 0, 5)
	st2 := discovery.NewState("0.0", discovery.OkStatus, discovery.PeerPosition{}, "0.0.0", 0.0, 0, 5)
	st3 := discovery.NewState("0.0", discovery.OkStatus, discovery.PeerPosition{}, "0.0.0", 0.0, 0, 5)
	p1 := discovery.NewPeerDetailed([]byte("key1"), seed1.IP, seed1.Port, time.Now(), st1)
	p2 := discovery.NewPeerDetailed([]byte("key2"), seed2.IP, seed1.Port, time.Now(), st2)
	p3 := discovery.NewPeerDetailed([]byte("key3"), seed3.IP, seed1.Port, time.Now(), st3)

	repo.AddPeer(p1)
	repo.AddPeer(p2)
	repo.AddPeer(p3)

	srv := NewService(repo, new(mockPeerMonitor), new(mockPeerNetworker), new(mockRobotWatcher))
	err := srv.RefreshPeer(p1)
	assert.Nil(t, err)
	assert.Equal(t, "0.62 0.77 0.71 4/972 26361", p1.CPULoad())
	assert.Equal(t, discovery.OkStatus, p1.Status())
	assert.Equal(t, float64(212383852), p1.FreeDiskSpace())
	assert.Equal(t, 3, p1.DiscoveredPeersNumber())
	assert.Equal(t, 1, p1.P2PFactor())
}

/*
Scenario: check state1
	Given a peer with 3 seed (discoveredPeersNumber=5 for all seed) / 5 peers on the repo
	When DiscoveredPeer=5 and elapsedheartbeats < Bootstrapingmintime
	Then state is OkStatus
*/

func TestState1(t *testing.T) {
	repo := new(mockPeerRepository)
	srv := NewService(repo, new(mockPeerMonitor), new(mockPeerNetworker), new(mockRobotWatcher))

	initP := discovery.NewStartupPeer([]byte("key"), net.ParseIP("127.0.0.1"), 3000, "0.0", discovery.PeerPosition{})
	repo.AddPeer(initP)
	seed1 := discovery.Seed{IP: net.ParseIP("10.1.1.1"), Port: 3000}
	seed2 := discovery.Seed{IP: net.ParseIP("10.1.1.2"), Port: 3001}
	seed3 := discovery.Seed{IP: net.ParseIP("10.1.1.3"), Port: 3002}
	repo.AddSeed(seed1)
	repo.AddSeed(seed2)
	repo.AddSeed(seed3)
	assert.Equal(t, 3, len(repo.seeds))
	st1 := discovery.NewState("0.0", discovery.OkStatus, discovery.PeerPosition{}, "0.0.0", 0.0, 0, 5)
	st2 := discovery.NewState("0.0", discovery.OkStatus, discovery.PeerPosition{}, "0.0.0", 0.0, 0, 5)
	st3 := discovery.NewState("0.0", discovery.OkStatus, discovery.PeerPosition{}, "0.0.0", 0.0, 0, 5)
	p1 := discovery.NewPeerDetailed([]byte("key1"), seed1.IP, seed1.Port, time.Now(), st1)
	p2 := discovery.NewPeerDetailed([]byte("key2"), seed2.IP, seed1.Port, time.Now(), st2)
	p3 := discovery.NewPeerDetailed([]byte("key3"), seed3.IP, seed1.Port, time.Now(), st3)
	repo.AddPeer(p1)
	repo.AddPeer(p2)
	repo.AddPeer(p3)
	p4 := discovery.NewPeerDetailed([]byte("key4"), net.ParseIP("185.123.4.9"), 4000, time.Now(), st1)
	repo.AddPeer(p4)
	assert.Equal(t, 5, len(repo.peers))
	selfpeer, err := repo.GetOwnedPeer()
	selfpeer.Refresh(discovery.BootstrapingStatus, 0.0, "0.0.0", 5, 0)
	s, err := srv.PeerStatus(selfpeer)
	assert.Equal(t, nil, err)
	assert.Equal(t, discovery.OkStatus, s)
}

/*
Scenario: check state2
	Given a peer with 3 seed (discoveredPeersNumber=5 for all seeds) / 5 peers on the repo / ntp offset is not fine
	When check state
	Then state is StorageOnlystate
*/

func TestState2(t *testing.T) {
	repo := new(mockPeerRepository)
	srv := NewService(repo, new(mockPeerMonitor), new(mockSystemNetworker2), new(mockRobotWatcher))

	initP := discovery.NewStartupPeer([]byte("key"), net.ParseIP("127.0.0.1"), 3000, "0.0", discovery.PeerPosition{})
	repo.AddPeer(initP)
	seed1 := discovery.Seed{IP: net.ParseIP("10.1.1.1"), Port: 3000}
	seed2 := discovery.Seed{IP: net.ParseIP("10.1.1.2"), Port: 3001}
	seed3 := discovery.Seed{IP: net.ParseIP("10.1.1.3"), Port: 3002}
	repo.AddSeed(seed1)
	repo.AddSeed(seed2)
	repo.AddSeed(seed3)
	assert.Equal(t, 3, len(repo.seeds))
	st1 := discovery.NewState("0.0", discovery.OkStatus, discovery.PeerPosition{}, "0.0.0", 0.0, 0, 5)
	st2 := discovery.NewState("0.0", discovery.OkStatus, discovery.PeerPosition{}, "0.0.0", 0.0, 0, 5)
	st3 := discovery.NewState("0.0", discovery.OkStatus, discovery.PeerPosition{}, "0.0.0", 0.0, 0, 5)
	p1 := discovery.NewPeerDetailed([]byte("key1"), seed1.IP, seed1.Port, time.Now(), st1)
	p2 := discovery.NewPeerDetailed([]byte("key2"), seed2.IP, seed1.Port, time.Now(), st2)
	p3 := discovery.NewPeerDetailed([]byte("key3"), seed3.IP, seed1.Port, time.Now(), st3)
	repo.AddPeer(p1)
	repo.AddPeer(p2)
	repo.AddPeer(p3)
	p4 := discovery.NewPeerDetailed([]byte("key4"), net.ParseIP("185.123.4.9"), 4000, time.Now(), st1)
	repo.AddPeer(p4)
	assert.Equal(t, 5, len(repo.peers))
	selfpeer, _ := repo.GetOwnedPeer()
	selfpeer.Refresh(discovery.BootstrapingStatus, 0.0, "0.0.0", 5, 0)
	s, _ := srv.PeerStatus(selfpeer)
	assert.Equal(t, discovery.StorageOnlyStatus, s)
}

/*
Scenario: check state3
	Given a peer with 3 seed (discoveredPeersNumber=5 for all seeds) / 5 peers on the repo / processstate is KO
	When check state
	Then state is FaultyState
*/

func TestState3(t *testing.T) {
	repo := new(mockPeerRepository)
	srv := NewService(repo, new(mockPeerMonitor), new(mockSystemNetworker3), new(mockRobotWatcher))

	initP := discovery.NewStartupPeer([]byte("key"), net.ParseIP("127.0.0.1"), 3000, "0.0", discovery.PeerPosition{})
	repo.AddPeer(initP)
	seed1 := discovery.Seed{IP: net.ParseIP("10.1.1.1"), Port: 3000}
	seed2 := discovery.Seed{IP: net.ParseIP("10.1.1.2"), Port: 3001}
	seed3 := discovery.Seed{IP: net.ParseIP("10.1.1.3"), Port: 3002}
	repo.AddSeed(seed1)
	repo.AddSeed(seed2)
	repo.AddSeed(seed3)
	assert.Equal(t, 3, len(repo.seeds))
	st1 := discovery.NewState("0.0", discovery.OkStatus, discovery.PeerPosition{}, "0.0.0", 0.0, 0, 5)
	st2 := discovery.NewState("0.0", discovery.OkStatus, discovery.PeerPosition{}, "0.0.0", 0.0, 0, 5)
	st3 := discovery.NewState("0.0", discovery.OkStatus, discovery.PeerPosition{}, "0.0.0", 0.0, 0, 5)
	p1 := discovery.NewPeerDetailed([]byte("key1"), seed1.IP, seed1.Port, time.Now(), st1)
	p2 := discovery.NewPeerDetailed([]byte("key2"), seed2.IP, seed1.Port, time.Now(), st2)
	p3 := discovery.NewPeerDetailed([]byte("key3"), seed3.IP, seed1.Port, time.Now(), st3)
	repo.AddPeer(p1)
	repo.AddPeer(p2)
	repo.AddPeer(p3)
	p4 := discovery.NewPeerDetailed([]byte("key4"), net.ParseIP("185.123.4.9"), 4000, time.Now(), st1)
	repo.AddPeer(p4)
	assert.Equal(t, 5, len(repo.peers))
	selfpeer, _ := repo.GetOwnedPeer()
	selfpeer.Refresh(discovery.BootstrapingStatus, 0.0, "0.0.0", 5, 0)
	s, _ := srv.PeerStatus(selfpeer)
	assert.Equal(t, discovery.FaultStatus, s)
}

type mockPeerRepository struct {
	peers []discovery.Peer
	seeds []discovery.Seed
}

func (r *mockPeerRepository) CountKnownPeers() (int, error) {
	return len(r.peers), nil
}

func (r *mockPeerRepository) GetOwnedPeer() (p discovery.Peer, err error) {
	for _, p := range r.peers {
		if p.IsOwned() {
			return p, nil
		}
	}
	return
}

func (r *mockPeerRepository) AddPeer(p discovery.Peer) error {
	if r.containsPeer(p) {
		return r.UpdatePeer(p)
	}
	r.peers = append(r.peers, p)
	return nil
}

func (r *mockPeerRepository) AddSeed(s discovery.Seed) error {
	r.seeds = append(r.seeds, s)
	return nil
}

func (r *mockPeerRepository) ListKnownPeers() ([]discovery.Peer, error) {
	return r.peers, nil
}

func (r *mockPeerRepository) ListSeedPeers() ([]discovery.Seed, error) {
	return r.seeds, nil
}

func (r *mockPeerRepository) GetPeerByIP(ip net.IP) (p discovery.Peer, err error) {
	for i := 0; i < len(r.peers); i++ {
		if string(ip) == string(r.peers[i].IP()) {
			return r.peers[i], nil
		}
	}
	return
}

func (r *mockPeerRepository) UpdatePeer(peer discovery.Peer) error {
	for _, p := range r.peers {
		if string(p.PublicKey()) == string(peer.PublicKey()) {
			p = peer
			break
		}
	}
	return nil
}

func (r *mockPeerRepository) containsPeer(peer discovery.Peer) bool {
	mPeers := make(map[string]discovery.Peer, 0)
	for _, p := range r.peers {
		mPeers[hex.EncodeToString(p.PublicKey())] = peer
	}

	_, exist := mPeers[hex.EncodeToString(peer.PublicKey())]
	return exist
}

type mockPeerMonitor struct{}

func (w mockPeerMonitor) CPULoad() (string, error) {
	return "0.62 0.77 0.71 4/972 26361", nil
}

func (w mockPeerMonitor) FreeDiskSpace() (float64, error) {
	return 212383852, nil
}

func (w mockPeerMonitor) P2PFactor() (int, error) {
	return 1, nil
}

func (w mockPeerMonitor) CheckAutodiscoveryProcess(discoveryPort int) error {
	return nil
}
func (w mockPeerMonitor) CheckMiningProcess() error {
	return nil
}
func (w mockPeerMonitor) CheckDataProcess() error {
	return nil
}
func (w mockPeerMonitor) CheckAIProcess() error {
	return nil
}
func (w mockPeerMonitor) CheckRedisProcess() error {
	return nil
}
func (w mockPeerMonitor) CheckScyllaDbProcess() error {
	return nil
}
func (w mockPeerMonitor) CheckRabbitmqProcess() error {
	return nil
}

type mockPeerNetworker struct{}

func (n mockPeerNetworker) IP() (net.IP, error) {
	return net.ParseIP("127.0.0.1"), nil
}

func (n mockPeerNetworker) CheckInternetState() error {
	return nil
}

func (n mockPeerNetworker) CheckNtpState() error {
	return nil
}

type mockRobotWatcher struct{}

func (r mockRobotWatcher) CheckAutodiscoveryProcess(port int) error {
	return nil
}

func (r mockRobotWatcher) CheckDataProcess() error {
	return nil
}

func (r mockRobotWatcher) CheckMiningProcess() error {
	return nil
}

func (r mockRobotWatcher) CheckAIProcess() error {
	return nil
}

func (r mockRobotWatcher) CheckScyllaDbProcess() error {
	return nil
}

func (r mockRobotWatcher) CheckRedisProcess() error {
	return nil
}

func (r mockRobotWatcher) CheckRabbitmqProcess() error {
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

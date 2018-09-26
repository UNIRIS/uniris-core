package system

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
Scenario: check GetSeedDiscoveredPeer
	Given a repo with 3 seed, seed1 discoveredPeersNumber = 5,seed2 discoveredPeersNumber = 6, seed3 discoveredPeersNumber = 7
	When GetSeedDiscoveredPeer call
	Then SeedDiscoveredPeer value is 6
*/

func TestGetSeedDiscoveredPeer(t *testing.T) {
	repo := new(mockPeerRepository)
	sdnw := mockseedDiscoverdPeerWatcher{}
	seed1 := discovery.Seed{IP: net.ParseIP("10.1.1.1"), Port: 3000}
	seed2 := discovery.Seed{IP: net.ParseIP("10.1.1.2"), Port: 3001}
	seed3 := discovery.Seed{IP: net.ParseIP("10.1.1.3"), Port: 3002}
	repo.AddSeed(seed1)
	repo.AddSeed(seed2)
	repo.AddSeed(seed3)
	assert.Equal(t, 3, len(repo.seeds))
	st1 := discovery.NewState("0.0", discovery.OkStatus, discovery.PeerPosition{}, "0.0.0", 0.0, 0, 5)
	st2 := discovery.NewState("0.0", discovery.OkStatus, discovery.PeerPosition{}, "0.0.0", 0.0, 0, 6)
	st3 := discovery.NewState("0.0", discovery.OkStatus, discovery.PeerPosition{}, "0.0.0", 0.0, 0, 7)
	p1 := discovery.NewPeerDetailed([]byte("key1"), seed1.IP, seed1.Port, time.Now(), st1)
	p2 := discovery.NewPeerDetailed([]byte("key2"), seed2.IP, seed1.Port, time.Now(), st2)
	p3 := discovery.NewPeerDetailed([]byte("key3"), seed3.IP, seed1.Port, time.Now(), st3)
	repo.AddPeer(p1)
	repo.AddPeer(p2)
	repo.AddPeer(p3)
	assert.Equal(t, 3, len(repo.peers))
	sdn, _ := sdnw.CountSeedDiscoveredPeer(repo)
	assert.Equal(t, 6, sdn)

}

/*
Scenario: check DiscoveredPeer
	Given a peer with 3 seed  / 5 peers on the repo
	When DiscoveredPeer
	Then return 5
*/

func TestDiscoveredPeer(t *testing.T) {
	repo := new(mockPeerRepository)
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
	sw := NewSystemWatcher(repo)
	dn, err := sw.CountDiscoveredPeer()
	assert.Equal(t, nil, err)
	assert.Equal(t, 5, dn)

}

/*
Scenario: check state1
	Given a peer with 3 seed (discoveredPeersNumber=5 for all seed) / 5 peers on the repo
	When DiscoveredPeer=5 and elapsedheartbeats < Bootstrapingmintime
	Then state is OkStatus
*/

func TestState1(t *testing.T) {
	repo := new(mockPeerRepository)
	sdnw := mockseedDiscoverdPeerWatcher{}
	pw := mockpeerWatcher{}
	w := mockwatcher{
		Pwatcher:   pw,
		SdnWatcher: sdnw,
		rep:        repo,
	}
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
	s, err := w.Status(selfpeer)
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
	sdnw := mockseedDiscoverdPeerWatcher{}
	pw := mockpeerWatcher2{}
	w := mockwatcher2{
		Pwatcher:   pw,
		SdnWatcher: sdnw,
		rep:        repo,
	}
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
	s, _ := w.Status(selfpeer)
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
	sdnw := mockseedDiscoverdPeerWatcher{}
	pw := mockpeerWatcher3{rep: repo}
	w := mockwatcher3{
		Pwatcher:   pw,
		SdnWatcher: sdnw,
		rep:        repo,
	}
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
	s, _ := w.Status(selfpeer)
	assert.Equal(t, discovery.FaultStatus, s)
}

type mockPeerRepository struct {
	peers []discovery.Peer
	seeds []discovery.Seed
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

type mockpeerWatcher struct {
}

func (pw mockpeerWatcher) CheckProcessStates(peer discovery.Peer) error {
	return nil
}

func (pw mockpeerWatcher) CheckInternetState() error {
	return nil
}

func (pw mockpeerWatcher) CheckNtpState() error {
	return nil
}

type mockpeerWatcher2 struct {
}

func (pw mockpeerWatcher2) CheckProcessStates(peer discovery.Peer) error {
	return nil
}

func (pw mockpeerWatcher2) CheckInternetState() error {
	return nil
}

func (pw mockpeerWatcher2) CheckNtpState() error {
	return errors.New("System Clock have a big Offset check the ntp configuration of the system")
}

type mockpeerWatcher3 struct {
	rep discovery.Repository
}

func (pw mockpeerWatcher3) CheckProcessStates(peer discovery.Peer) error {
	return errors.New("required processes are not running")
}

func (pw mockpeerWatcher3) CheckInternetState() error {
	return nil
}

func (pw mockpeerWatcher3) CheckNtpState() error {
	return nil
}

type mockseedDiscoverdPeerWatcher struct {
}

func (sdnw mockseedDiscoverdPeerWatcher) CountSeedDiscoveredPeer(rep discovery.Repository) (int, error) {
	listseed, err := rep.ListSeedPeers()
	if err != nil {
		return 0, err
	}
	avg := 0
	for i := 0; i < len(listseed); i++ {
		ipseed := listseed[i].IP
		p, err := rep.GetPeerByIP(ipseed)
		if err == nil {
			avg += p.DiscoveredPeersNumber()
		}
	}
	avg = avg / len(listseed)
	return avg, nil
}

type mockwatcher struct {
	Pwatcher   mockpeerWatcher
	SdnWatcher mockseedDiscoverdPeerWatcher
	rep        discovery.Repository
}

func (w mockwatcher) Status(p discovery.Peer) (discovery.PeerStatus, error) {
	err := w.Pwatcher.CheckProcessStates(p)
	if err != nil {
		return discovery.FaultStatus, err
	}
	err = w.Pwatcher.CheckInternetState()
	if err != nil {
		return discovery.FaultStatus, err
	}
	err = w.Pwatcher.CheckNtpState()
	if err != nil {
		return discovery.StorageOnlyStatus, err
	}
	seedDn, err := w.SdnWatcher.CountSeedDiscoveredPeer(w.rep)
	if err != nil {
		return discovery.FaultStatus, err
	}
	if seedDn == 0 {
		return discovery.BootstrapingStatus, nil
	}
	if t := p.GetElapsedHeartbeats(); t < discovery.BootStrapingMinTime && seedDn > p.DiscoveredPeersNumber() {
		return discovery.BootstrapingStatus, nil
	} else if t < discovery.BootStrapingMinTime && seedDn <= p.DiscoveredPeersNumber() {
		return discovery.OkStatus, nil
	} else {
		return discovery.OkStatus, nil
	}
}

type mockwatcher2 struct {
	Pwatcher   mockpeerWatcher2
	SdnWatcher mockseedDiscoverdPeerWatcher
	rep        discovery.Repository
}

func (w mockwatcher2) Status(p discovery.Peer) (discovery.PeerStatus, error) {
	err := w.Pwatcher.CheckProcessStates(p)
	if err != nil {
		return discovery.FaultStatus, err
	}
	err = w.Pwatcher.CheckInternetState()
	if err != nil {
		return discovery.FaultStatus, err
	}
	err = w.Pwatcher.CheckNtpState()
	if err != nil {
		return discovery.StorageOnlyStatus, err
	}
	seedDn, err := w.SdnWatcher.CountSeedDiscoveredPeer(w.rep)
	if err != nil {
		return discovery.FaultStatus, err
	}
	if seedDn == 0 {
		return discovery.BootstrapingStatus, nil
	}
	if t := p.GetElapsedHeartbeats(); t < discovery.BootStrapingMinTime && seedDn > p.DiscoveredPeersNumber() {
		return discovery.BootstrapingStatus, nil
	} else if t < discovery.BootStrapingMinTime && seedDn <= p.DiscoveredPeersNumber() {
		return discovery.OkStatus, nil
	} else {
		return discovery.OkStatus, nil
	}
}

type mockwatcher3 struct {
	Pwatcher   mockpeerWatcher3
	SdnWatcher mockseedDiscoverdPeerWatcher
	rep        discovery.Repository
}

func (w mockwatcher3) Status(p discovery.Peer) (discovery.PeerStatus, error) {
	err := w.Pwatcher.CheckProcessStates(p)
	if err != nil {
		return discovery.FaultStatus, err
	}
	err = w.Pwatcher.CheckInternetState()
	if err != nil {
		return discovery.FaultStatus, err
	}
	err = w.Pwatcher.CheckNtpState()
	if err != nil {
		return discovery.StorageOnlyStatus, err
	}
	seedDn, err := w.SdnWatcher.CountSeedDiscoveredPeer(w.rep)
	if err != nil {
		return discovery.FaultStatus, err
	}
	if seedDn == 0 {
		return discovery.BootstrapingStatus, nil
	}
	if t := p.GetElapsedHeartbeats(); t < discovery.BootStrapingMinTime && seedDn > p.DiscoveredPeersNumber() {
		return discovery.BootstrapingStatus, nil
	} else if t < discovery.BootStrapingMinTime && seedDn <= p.DiscoveredPeersNumber() {
		return discovery.OkStatus, nil
	} else {
		return discovery.OkStatus, nil
	}
}

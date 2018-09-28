package rpc

import (
	"encoding/hex"
	"errors"
	"net"
	"testing"
	"time"

	"github.com/uniris/uniris-core/autodiscovery/pkg/gossip"
	"github.com/uniris/uniris-core/autodiscovery/pkg/monitoring"

	"github.com/stretchr/testify/assert"

	"golang.org/x/net/context"

	api "github.com/uniris/uniris-core/autodiscovery/api/protobuf-spec"
	discovery "github.com/uniris/uniris-core/autodiscovery/pkg"
)

/*
Scenario: Process a SYN request by returning new peers and unknown peers
	Given a GRPC server and a peer in our repository
	When we receive a SYN request we compute the diff between our peers and the sent peers
	Then we returns the unknown peers from the sender and the unknown peers locally
*/
func TestHandleSynRequest(t *testing.T) {

	repo := new(mockRepository)
	repo.SetPeer(discovery.NewPeerDigest(
		discovery.NewPeerIdentity(net.ParseIP("20.10.0.1"), 3000, []byte("key1")),
		discovery.NewPeerHeartbeatState(time.Now(), 1000),
	))

	notif := new(mockNotifier)
	msg := new(mockMessenger)
	mon := monitoring.NewService(repo, new(mockMonitor), new(mockNetworker), new(mockRobotWatcher))
	h := NewHandler(repo, gossip.NewService(repo, msg, notif), mon, notif)

	req := &api.SynRequest{
		Initiator: &api.PeerDigest{},
		Target:    &api.PeerDigest{},
		KnownPeers: []*api.PeerDigest{
			&api.PeerDigest{
				Identity: &api.PeerIdentity{
					IP:        "30.0.0.1",
					Port:      3000,
					PublicKey: []byte("key2"),
				},
				HeartbeatState: &api.PeerHeartbeatState{
					GenerationTime:    int64(time.Now().Unix()),
					ElapsedHeartbeats: 1000,
				},
			},
		},
	}
	res, err := h.Synchronize(context.TODO(), req)
	assert.Nil(t, err)
	assert.NotNil(t, res)

	assert.NotEmpty(t, res.NewPeers)
	assert.NotEmpty(t, res.UnknownPeers)
	assert.Equal(t, "key1", string(res.NewPeers[0].Identity.PublicKey))
	assert.Equal(t, "key2", string(res.UnknownPeers[0].Identity.PublicKey))
}

/*
Scenario: Process a SYN request by returning only the new peers
	Given a GRPC server and a peer in our repository
	When we receive a SYN request we compute the diff between our peers and the sent peers
	Then we returns the unknown peers from the sender
*/
func TestHandleSynRequestNewPeers(t *testing.T) {
	repo := new(mockRepository)
	repo.SetPeer(discovery.NewPeerDigest(
		discovery.NewPeerIdentity(net.ParseIP("20.10.0.1"), 3000, []byte("key1")),
		discovery.NewPeerHeartbeatState(time.Now(), 1000),
	))
	repo.SetPeer(discovery.NewPeerDigest(
		discovery.NewPeerIdentity(net.ParseIP("30.0.0.1"), 3000, []byte("key2")),
		discovery.NewPeerHeartbeatState(time.Now(), 1000),
	))

	notif := new(mockNotifier)
	msg := new(mockMessenger)
	mon := monitoring.NewService(repo, new(mockMonitor), new(mockNetworker), new(mockRobotWatcher))
	h := NewHandler(repo, gossip.NewService(repo, msg, notif), mon, notif)

	req := &api.SynRequest{
		Initiator: &api.PeerDigest{},
		Target:    &api.PeerDigest{},
		KnownPeers: []*api.PeerDigest{
			&api.PeerDigest{
				Identity: &api.PeerIdentity{
					IP:        "30.0.0.1",
					Port:      3000,
					PublicKey: []byte("key2"),
				},
				HeartbeatState: &api.PeerHeartbeatState{
					GenerationTime:    int64(time.Now().Unix()),
					ElapsedHeartbeats: 1000,
				},
			},
		},
	}
	res, err := h.Synchronize(context.TODO(), req)
	assert.Nil(t, err)
	assert.NotNil(t, res)

	assert.NotEmpty(t, res.NewPeers)
	assert.Empty(t, res.UnknownPeers)
	assert.Equal(t, "key1", string(res.NewPeers[0].Identity.PublicKey))
}

/*
Scenario: Process a SYN request by returning only the new peers recent
	Given a GRPC server and a peer in our repository
	When we receive a SYN request we compute the diff between our peers and the sent peers
	Then we returns the unknown peers more recent from the sender
*/
func TestHandleSynRequestNewPeersRecentOnly(t *testing.T) {
	repo := new(mockRepository)
	repo.SetPeer(discovery.NewPeerDigest(
		discovery.NewPeerIdentity(net.ParseIP("20.10.0.1"), 3000, []byte("key1")),
		discovery.NewPeerHeartbeatState(time.Now(), 1000),
	))
	repo.SetPeer(discovery.NewPeerDigest(
		discovery.NewPeerIdentity(net.ParseIP("30.0.0.1"), 3000, []byte("key2")),
		discovery.NewPeerHeartbeatState(time.Now(), 1000),
	))

	notif := new(mockNotifier)
	msg := new(mockMessenger)
	mon := monitoring.NewService(repo, new(mockMonitor), new(mockNetworker), new(mockRobotWatcher))
	h := NewHandler(repo, gossip.NewService(repo, msg, notif), mon, notif)

	req := &api.SynRequest{
		Initiator: &api.PeerDigest{},
		Target:    &api.PeerDigest{},
		KnownPeers: []*api.PeerDigest{
			&api.PeerDigest{
				Identity: &api.PeerIdentity{
					IP:        "30.0.0.1",
					Port:      3000,
					PublicKey: []byte("key2"),
				},
				HeartbeatState: &api.PeerHeartbeatState{
					GenerationTime:    int64(time.Now().Unix()),
					ElapsedHeartbeats: 1000,
				},
			},
			&api.PeerDigest{
				Identity: &api.PeerIdentity{
					IP:        "20.10.0.1",
					Port:      3000,
					PublicKey: []byte("key2"),
				},
				HeartbeatState: &api.PeerHeartbeatState{
					GenerationTime:    int64(time.Now().Unix()),
					ElapsedHeartbeats: 1500,
				},
			},
		},
	}
	res, err := h.Synchronize(context.TODO(), req)
	assert.Nil(t, err)
	assert.NotNil(t, res)

	assert.NotEmpty(t, res.NewPeers)
	assert.NotEmpty(t, res.UnknownPeers)

	assert.Equal(t, int64(1500), res.UnknownPeers[0].HeartbeatState.ElapsedHeartbeats)
}

/*
Scenario: Process a SYN request by returning only the unknown peers
	Given a GRPC server and a peer in our repository
	When we receive a SYN request we compute the diff between our peers and the sent peers
	Then we returns the unknown peers locally
*/
func TestHandleSynRequestUnknownPeers(t *testing.T) {
	repo := new(mockRepository)
	repo.SetPeer(discovery.NewPeerDigest(
		discovery.NewPeerIdentity(net.ParseIP("20.10.0.1"), 3000, []byte("key1")),
		discovery.NewPeerHeartbeatState(time.Now(), 1000),
	))

	notif := new(mockNotifier)
	msg := new(mockMessenger)
	mon := monitoring.NewService(repo, new(mockMonitor), new(mockNetworker), new(mockRobotWatcher))
	h := NewHandler(repo, gossip.NewService(repo, msg, notif), mon, notif)

	req := &api.SynRequest{
		Initiator: &api.PeerDigest{},
		Target:    &api.PeerDigest{},
		KnownPeers: []*api.PeerDigest{
			&api.PeerDigest{
				Identity: &api.PeerIdentity{
					IP:        "30.0.0.1",
					Port:      3000,
					PublicKey: []byte("key2"),
				},
				HeartbeatState: &api.PeerHeartbeatState{
					GenerationTime:    int64(time.Now().Unix()),
					ElapsedHeartbeats: 1000,
				},
			},
			&api.PeerDigest{
				Identity: &api.PeerIdentity{
					IP:        "20.10.0.1",
					Port:      3000,
					PublicKey: []byte("key1"),
				},
				HeartbeatState: &api.PeerHeartbeatState{
					GenerationTime:    int64(time.Now().Unix()),
					ElapsedHeartbeats: 1000,
				},
			},
		},
	}
	res, err := h.Synchronize(context.TODO(), req)
	assert.Nil(t, err)
	assert.NotNil(t, res)

	assert.Empty(t, res.NewPeers)
	assert.NotEmpty(t, res.UnknownPeers)
	assert.Equal(t, "key2", string(res.UnknownPeers[0].Identity.PublicKey))
}

/*
Scenario: Process a SYN request by returning only the unknown peers more recent
	Given a GRPC server and a peer in our repository
	When we receive a SYN request we compute the diff between our peers and the sent peers
	Then we returns the unknown peers locally more recent
*/
func TestHandleSynRequestUnknownPeersRecentOnly(t *testing.T) {
	repo := new(mockRepository)
	repo.SetPeer(discovery.NewPeerDigest(
		discovery.NewPeerIdentity(net.ParseIP("20.10.0.1"), 3000, []byte("key1")),
		discovery.NewPeerHeartbeatState(time.Now(), 1000),
	))
	repo.SetPeer(discovery.NewPeerDigest(
		discovery.NewPeerIdentity(net.ParseIP("30.0.0.1"), 3000, []byte("key2")),
		discovery.NewPeerHeartbeatState(time.Now(), 1000),
	))

	notif := new(mockNotifier)
	msg := new(mockMessenger)
	mon := monitoring.NewService(repo, new(mockMonitor), new(mockNetworker), new(mockRobotWatcher))
	h := NewHandler(repo, gossip.NewService(repo, msg, notif), mon, notif)

	req := &api.SynRequest{
		Initiator: &api.PeerDigest{},
		Target:    &api.PeerDigest{},
		KnownPeers: []*api.PeerDigest{
			&api.PeerDigest{
				Identity: &api.PeerIdentity{
					IP:        "30.0.0.1",
					Port:      3000,
					PublicKey: []byte("key2"),
				},
				HeartbeatState: &api.PeerHeartbeatState{
					GenerationTime:    int64(time.Now().Unix()),
					ElapsedHeartbeats: 1200,
				},
			},
			&api.PeerDigest{
				Identity: &api.PeerIdentity{
					IP:        "20.10.0.1",
					Port:      3000,
					PublicKey: []byte("key1"),
				},
				HeartbeatState: &api.PeerHeartbeatState{
					GenerationTime:    int64(time.Now().Unix()),
					ElapsedHeartbeats: 1000,
				},
			},
		},
	}
	res, err := h.Synchronize(context.TODO(), req)
	assert.Nil(t, err)
	assert.NotNil(t, res)

	assert.Empty(t, res.NewPeers)
	assert.NotEmpty(t, res.UnknownPeers)
	assert.Equal(t, "key2", string(res.UnknownPeers[0].Identity.PublicKey))
	assert.Equal(t, int64(1200), res.UnknownPeers[0].HeartbeatState.ElapsedHeartbeats)
}

/*
Scenario: Process a ACK request by saving and notifying the request detailed peers
	Given a GRPC server
	When we receive a ACK request
	Then we store and notified the new peers
*/
func TestHandlAckRequest(t *testing.T) {
	notif := new(mockNotifier)
	msg := new(mockMessenger)
	repo := new(mockRepository)
	mon := monitoring.NewService(repo, new(mockMonitor), new(mockNetworker), new(mockRobotWatcher))
	h := NewHandler(repo, gossip.NewService(repo, msg, notif), mon, notif)

	req := &api.AckRequest{
		Initiator: &api.PeerDigest{},
		Target:    &api.PeerDigest{},
		RequestedPeers: []*api.PeerDiscovered{
			&api.PeerDiscovered{
				Identity: &api.PeerIdentity{
					IP:        "20.10.0.1",
					Port:      3000,
					PublicKey: []byte("key1"),
				},
				HeartbeatState: &api.PeerHeartbeatState{
					GenerationTime:    int64(time.Now().Unix()),
					ElapsedHeartbeats: 1000,
				},
				AppState: &api.PeerAppState{
					CPULoad:       "00-00-000",
					FreeDiskSpace: 3000,
					GeoPosition: &api.PeerAppState_GeoCoordinates{
						Lat: 30.0,
						Lon: 50.0,
					},
					P2PFactor: 2,
					Status:    api.PeerAppState_Ok,
					Version:   "1.0",
				},
			},
		},
	}
	_, err := h.Acknowledge(context.TODO(), req)
	assert.Nil(t, err)

	kp, _ := repo.ListDiscoveredPeers()
	assert.Equal(t, 1, len(kp))
	assert.Equal(t, "key1", kp[0].Identity().PublicKey().String())
	assert.Equal(t, "20.10.0.1", kp[0].Identity().IP().String())
	assert.Equal(t, 3000, kp[0].Identity().Port())
	assert.Equal(t, int64(1000), kp[0].HeartbeatState().ElapsedHeartbeats())
	assert.Equal(t, "00-00-000", kp[0].AppState().CPULoad())
	assert.Equal(t, float64(3000), kp[0].AppState().FreeDiskSpace())
	assert.Equal(t, 30.0, kp[0].AppState().GeoPosition().Lat)
	assert.Equal(t, 50.0, kp[0].AppState().GeoPosition().Lon)
	assert.Equal(t, 2, kp[0].AppState().P2PFactor())
	assert.Equal(t, "1.0", kp[0].AppState().Version())
	assert.Equal(t, discovery.OkStatus, kp[0].AppState().Status())

	assert.NotEmpty(t, notif.NotifiedPeers)
	assert.Equal(t, "key1", notif.NotifiedPeers()[0].Identity().PublicKey().String())
}

//////////////////////////////////////////////////////////
// 						MOCKS
/////////////////////////////////////////////////////////

type mockMessenger struct {
}

func (m mockMessenger) SendSyn(req gossip.SynRequest) (*gossip.SynAck, error) {
	init := discovery.NewStartupPeer([]byte("key"), net.ParseIP("127.0.0.1"), 3000, "1.0", discovery.PeerPosition{})
	tar := discovery.NewStartupPeer([]byte("uKey1"), net.ParseIP("200.18.186.39"), 3000, "1.1", discovery.PeerPosition{})

	hb := discovery.NewPeerHeartbeatState(time.Now(), 0)
	as := discovery.NewPeerAppState("1.0", discovery.OkStatus, discovery.PeerPosition{}, "", 0, 1, 0)

	np1 := discovery.NewDiscoveredPeer(
		discovery.NewPeerIdentity(net.ParseIP("35.200.100.2"), 3000, []byte("dKey1")),
		hb, as,
	)

	newPeers := []discovery.Peer{np1}

	unknownPeers := []discovery.Peer{tar}

	return &gossip.SynAck{
		Initiator:    init,
		Target:       tar,
		NewPeers:     newPeers,
		UnknownPeers: unknownPeers,
	}, nil
}

func (m mockMessenger) SendAck(req gossip.AckRequest) error {
	return nil
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

type mockNotifier struct {
	notifiedPeers []discovery.Peer
}

func (n mockNotifier) NotifiedPeers() []discovery.Peer {
	return n.notifiedPeers
}

func (n *mockNotifier) Notify(p discovery.Peer) {
	n.notifiedPeers = append(n.notifiedPeers, p)
}

type mockMonitor struct{}

func (w mockMonitor) CPULoad() (string, error) {
	return "0.62 0.77 0.71 4/972 26361", nil
}

func (w mockMonitor) FreeDiskSpace() (float64, error) {
	return 212383852, nil
}

func (w mockMonitor) P2PFactor() (int, error) {
	return 1, nil
}

type mockNetworker struct{}

func (n mockNetworker) IP() (net.IP, error) {
	return net.ParseIP("127.0.0.1"), nil
}

func (n mockNetworker) CheckInternetState() error {
	return nil
}

func (n mockNetworker) CheckNtpState() error {
	return nil
}

type mockNetworkerNTPFails struct{}

func (n mockNetworkerNTPFails) IP() (net.IP, error) {
	return net.ParseIP("127.0.0.1"), nil
}

func (n mockNetworkerNTPFails) CheckInternetState() error {
	return nil
}

func (n mockNetworkerNTPFails) CheckNtpState() error {
	return errors.New("System Clock have a big Offset check the ntp configuration of the system")
}

type mockNetworkerInternetFails struct{}

func (n mockNetworkerInternetFails) IP() (net.IP, error) {
	return net.ParseIP("127.0.0.1"), nil
}

func (n mockNetworkerInternetFails) CheckInternetState() error {
	return errors.New("required processes are not running")
}

func (n mockNetworkerInternetFails) CheckNtpState() error {
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

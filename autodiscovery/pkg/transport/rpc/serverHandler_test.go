package rpc

import (
	"net"
	"testing"
	"time"

	"github.com/uniris/uniris-core/autodiscovery/pkg/mock"

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

	repo := new(mock.Repository)
	repo.SetKnownPeer(discovery.NewPeerDigest(
		discovery.NewPeerIdentity(net.ParseIP("20.10.0.1"), 3000, "key1"),
		discovery.NewPeerHeartbeatState(time.Now(), 1000),
	))

	repo.SetKnownPeer(
		discovery.NewStartupPeer("key2", net.ParseIP("127.0.0.1"), 3000, "1.0", discovery.PeerPosition{}))

	notif := new(mock.Notifier)

	h := NewServerHandler(repo, notif)

	req := &api.SynRequest{
		Initiator: &api.PeerDigest{},
		Target:    &api.PeerDigest{},
		KnownPeers: []*api.PeerDigest{
			&api.PeerDigest{
				Identity: &api.PeerIdentity{
					IP:        "30.0.0.1",
					Port:      3000,
					PublicKey: "key2",
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
	repo := new(mock.Repository)
	repo.SetKnownPeer(discovery.NewPeerDigest(
		discovery.NewPeerIdentity(net.ParseIP("20.10.0.1"), 3000, "key1"),
		discovery.NewPeerHeartbeatState(time.Now(), 1000),
	))
	repo.SetKnownPeer(discovery.NewPeerDigest(
		discovery.NewPeerIdentity(net.ParseIP("30.0.0.1"), 3000, "key2"),
		discovery.NewPeerHeartbeatState(time.Now(), 1000),
	))

	repo.SetKnownPeer(
		discovery.NewStartupPeer("key3", net.ParseIP("127.0.0.1"), 3000, "1.0", discovery.PeerPosition{}))

	notif := new(mock.Notifier)

	h := NewServerHandler(repo, notif)

	req := &api.SynRequest{
		Initiator: &api.PeerDigest{},
		Target:    &api.PeerDigest{},
		KnownPeers: []*api.PeerDigest{
			&api.PeerDigest{
				Identity: &api.PeerIdentity{
					IP:        "30.0.0.1",
					Port:      3000,
					PublicKey: "key2",
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
	repo := new(mock.Repository)
	repo.SetKnownPeer(discovery.NewPeerDigest(
		discovery.NewPeerIdentity(net.ParseIP("20.10.0.1"), 3000, "key1"),
		discovery.NewPeerHeartbeatState(time.Now(), 1000),
	))
	repo.SetKnownPeer(discovery.NewPeerDigest(
		discovery.NewPeerIdentity(net.ParseIP("30.0.0.1"), 3000, "key2"),
		discovery.NewPeerHeartbeatState(time.Now(), 1000),
	))

	repo.SetKnownPeer(
		discovery.NewStartupPeer("key3", net.ParseIP("127.0.0.1"), 3000, "1.0", discovery.PeerPosition{}))

	notif := new(mock.Notifier)

	h := NewServerHandler(repo, notif)

	req := &api.SynRequest{
		Initiator: &api.PeerDigest{},
		Target:    &api.PeerDigest{},
		KnownPeers: []*api.PeerDigest{
			&api.PeerDigest{
				Identity: &api.PeerIdentity{
					IP:        "30.0.0.1",
					Port:      3000,
					PublicKey: "key2",
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
					PublicKey: "key2",
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
	repo := new(mock.Repository)
	repo.SetKnownPeer(discovery.NewPeerDigest(
		discovery.NewPeerIdentity(net.ParseIP("20.10.0.1"), 3000, "key1"),
		discovery.NewPeerHeartbeatState(time.Now(), 1000),
	))

	repo.SetKnownPeer(
		discovery.NewStartupPeer("key2", net.ParseIP("127.0.0.1"), 3000, "1.0", discovery.PeerPosition{}))

	notif := new(mock.Notifier)
	h := NewServerHandler(repo, notif)

	req := &api.SynRequest{
		Initiator: &api.PeerDigest{},
		Target:    &api.PeerDigest{},
		KnownPeers: []*api.PeerDigest{
			&api.PeerDigest{
				Identity: &api.PeerIdentity{
					IP:        "30.0.0.1",
					Port:      3000,
					PublicKey: "key2",
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
					PublicKey: "key1",
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
	repo := new(mock.Repository)
	repo.SetKnownPeer(discovery.NewPeerDigest(
		discovery.NewPeerIdentity(net.ParseIP("30.0.0.1"), 3000, "key2"),
		discovery.NewPeerHeartbeatState(time.Now(), 1000),
	))

	repo.SetKnownPeer(
		discovery.NewStartupPeer("key1", net.ParseIP("20.10.0.1"), 3000, "1.0", discovery.PeerPosition{}))

	notif := new(mock.Notifier)

	h := NewServerHandler(repo, notif)

	req := &api.SynRequest{
		Initiator: &api.PeerDigest{},
		Target:    &api.PeerDigest{},
		KnownPeers: []*api.PeerDigest{
			&api.PeerDigest{
				Identity: &api.PeerIdentity{
					IP:        "30.0.0.1",
					Port:      3000,
					PublicKey: "key2",
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
					PublicKey: "key1",
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
	notif := new(mock.Notifier)
	repo := new(mock.Repository)
	h := NewServerHandler(repo, notif)

	req := &api.AckRequest{
		Initiator: &api.PeerDigest{},
		Target:    &api.PeerDigest{},
		RequestedPeers: []*api.PeerDiscovered{
			&api.PeerDiscovered{
				Identity: &api.PeerIdentity{
					IP:        "20.10.0.1",
					Port:      3000,
					PublicKey: "key1",
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

	kp, _ := repo.ListKnownPeers()
	assert.Equal(t, 1, len(kp))
	assert.Equal(t, "key1", kp[0].Identity().PublicKey())
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
	assert.Equal(t, "key1", notif.NotifiedPeers()[0].Identity().PublicKey())
}

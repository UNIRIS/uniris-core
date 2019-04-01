package rpc

import (
	"context"
	"github.com/uniris/uniris-core/pkg/logging"
	"log"
	"net"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	api "github.com/uniris/uniris-core/api/protobuf-spec"
	"github.com/uniris/uniris-core/pkg/discovery"
)

/*
Scenario: Receive a synchronize request without knowing any peers
	Given not peers discovered
	When I receive a syn request including a peer
	Then I when I make diff , I return the sended peer as unknown
*/
func TestHandleSynchronizeRequestWithoutKnownPeers(t *testing.T) {
	db := &mockDiscoveryDB{}
	l := logging.NewLogger("stdout", log.New(os.Stdout, "", 0), "test", net.ParseIP("127.0.0.1"), logging.ErrorLogLevel)

	discoverySrv := NewDiscoveryServer(db, &mockDiscoveryNotifier{}, l)

	res, err := discoverySrv.Synchronize(context.TODO(), &api.SynRequest{
		KnownPeers: []*api.PeerDigest{
			&api.PeerDigest{
				HeartbeatState: &api.PeerHeartbeatState{
					ElapsedHeartbeats: 0,
					GenerationTime:    time.Now().Unix(),
				},
				Identity: &api.PeerIdentity{
					Ip:        "127.0.0.1",
					Port:      3000,
					PublicKey: "pubkey",
				},
			},
		},
	})
	assert.Nil(t, err)
	assert.NotEmpty(t, res.RequestedPeers)
	assert.Equal(t, "pubkey", res.RequestedPeers[0].PublicKey)
}

/*
Scenario: Receive a synchronize request with knowing any peers
	Given a peer discovered
	When I receive a syn request
	Then I when I make diff , I return the a peer the sender does not known
*/
func TestHandleSynchronizeRequestByKnowingPeer(t *testing.T) {
	db := &mockDiscoveryDB{
		discoveredPeers: []discovery.Peer{
			discovery.NewDiscoveredPeer(
				discovery.NewPeerIdentity(net.ParseIP("10.0.0.1"), 3000, "pubkey000"),
				discovery.NewPeerHeartbeatState(time.Now(), 1000),
				discovery.NewPeerAppState("1.0.1", discovery.OkPeerStatus, 10.0, 20.0, "", 300, 1, 100),
			),
		},
	}

	l := logging.NewLogger("stdout", log.New(os.Stdout, "", 0), "test", net.ParseIP("127.0.0.1"), logging.ErrorLogLevel)
	discoverySrv := NewDiscoveryServer(db, &mockDiscoveryNotifier{}, l)

	res, err := discoverySrv.Synchronize(context.TODO(), &api.SynRequest{
		KnownPeers: []*api.PeerDigest{
			&api.PeerDigest{
				HeartbeatState: &api.PeerHeartbeatState{
					ElapsedHeartbeats: 0,
					GenerationTime:    time.Now().Unix(),
				},
				Identity: &api.PeerIdentity{
					Ip:        "127.0.0.1",
					Port:      3000,
					PublicKey: "pubkey",
				},
			},
		},
	})
	assert.Nil(t, err)
	assert.NotEmpty(t, res.RequestedPeers)
	assert.Equal(t, "pubkey", res.RequestedPeers[0].PublicKey)
	assert.NotEmpty(t, res.DiscoveredPeers)
	assert.Equal(t, "pubkey000", res.DiscoveredPeers[0].Identity.PublicKey)
}

/*
Scenario: Receive an acknowledgement request with the details for the requested peers
	Given a requested peers details
	When I want to acknowledge them
	Then I store inside the db the discovered peers
*/
func TestHandleAcknowledgeRequest(t *testing.T) {
	db := &mockDiscoveryDB{}
	l := logging.NewLogger("stdout", log.New(os.Stdout, "", 0), "test", net.ParseIP("127.0.0.1"), logging.ErrorLogLevel)
	discoverySrv := NewDiscoveryServer(db, &mockDiscoveryNotifier{}, l)

	_, err := discoverySrv.Acknowledge(context.TODO(), &api.AckRequest{
		RequestedPeers: []*api.PeerDiscovered{
			&api.PeerDiscovered{
				Identity: &api.PeerIdentity{
					Ip:        "127.0.0.1",
					Port:      3000,
					PublicKey: "pubkey",
				},
				HeartbeatState: &api.PeerHeartbeatState{
					ElapsedHeartbeats: 1000,
					GenerationTime:    time.Now().Unix(),
				},
				AppState: &api.PeerAppState{
					CpuLoad:              "",
					ReachablePeersNumber: 100,
					FreeDiskSpace:        300,
					GeoPosition: &api.PeerAppState_GeoCoordinates{
						Latitude:  30.0,
						Longitude: 20.0,
					},
					P2PFactor: 1,
					Status:    api.PeerAppState_OK,
					Version:   "1.0",
				},
			},
		},
	})

	assert.Nil(t, err)
	assert.Len(t, db.discoveredPeers, 1)
	assert.Equal(t, "pubkey", db.discoveredPeers[0].Identity().PublicKey())
}

type mockDiscoveryDB struct {
	discoveredPeers  []discovery.Peer
	unreachablePeers []discovery.PeerIdentity
}

func (db mockDiscoveryDB) DiscoveredPeers() ([]discovery.Peer, error) {
	return db.discoveredPeers, nil
}

func (db *mockDiscoveryDB) WriteDiscoveredPeer(peer discovery.Peer) error {
	for i, p := range db.discoveredPeers {
		if p.Identity().PublicKey() == peer.Identity().PublicKey() {
			db.discoveredPeers[i] = peer
			return nil
		}
	}
	db.discoveredPeers = append(db.discoveredPeers, peer)
	return nil
}

func (db mockDiscoveryDB) UnreachablePeers() ([]discovery.PeerIdentity, error) {
	pp := make([]discovery.PeerIdentity, 0)
	for i := 0; i < len(db.discoveredPeers); i++ {
		if exist, _ := db.ContainsUnreachablePeer(db.discoveredPeers[i].Identity()); exist {
			pp = append(pp, db.discoveredPeers[i].Identity())
		}
	}
	return pp, nil
}

func (db *mockDiscoveryDB) WriteUnreachablePeer(p discovery.PeerIdentity) error {
	if exist, _ := db.ContainsUnreachablePeer(p); !exist {
		db.unreachablePeers = append(db.unreachablePeers, p)
	}
	return nil
}

func (db *mockDiscoveryDB) RemoveUnreachablePeer(p discovery.PeerIdentity) error {
	if exist, _ := db.ContainsUnreachablePeer(p); exist {
		for i := 0; i < len(db.unreachablePeers); i++ {
			if db.unreachablePeers[i].PublicKey() == p.PublicKey() {
				db.unreachablePeers = db.unreachablePeers[:i+copy(db.unreachablePeers[i:], db.unreachablePeers[i+1:])]
			}
		}
	}
	return nil
}

func (db mockDiscoveryDB) ContainsUnreachablePeer(peerPubK discovery.PeerIdentity) (bool, error) {
	for _, up := range db.unreachablePeers {
		if up.PublicKey() == peerPubK.PublicKey() {
			return true, nil
		}
	}
	return false, nil
}

type mockDiscoveryNotifier struct {
	reaches     []string
	unreaches   []string
	discoveries []discovery.Peer
}

func (n *mockDiscoveryNotifier) NotifyReachable(pk string) error {
	n.reaches = append(n.reaches, pk)
	return nil
}
func (n *mockDiscoveryNotifier) NotifyUnreachable(pk string) error {
	n.unreaches = append(n.unreaches, pk)
	return nil
}

func (n *mockDiscoveryNotifier) NotifyDiscovery(p discovery.Peer) error {
	n.discoveries = append(n.discoveries, p)
	return nil
}

type mockPeerNetworker struct{}

func (pn mockPeerNetworker) CheckNtpState() error {
	return nil
}

func (pn mockPeerNetworker) CheckInternetState() error {
	return nil
}

type mockPeerInfo struct{}

func (i mockPeerInfo) GeoPosition() (lon float64, lat float64, err error) {
	return 10.0, 30.0, nil
}

func (i mockPeerInfo) FreeDiskSpace() (float64, error) {
	return 200, nil
}

func (i mockPeerInfo) CPULoad() (string, error) {
	return "", nil
}

func (i mockPeerInfo) IP() (net.IP, error) {
	return net.ParseIP("127.0.0.1"), nil
}

package rpc

import (
	"context"
	"net"
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
	repo := &mockDiscoveryRepo{}

	service := discovery.NewService(repo, nil, &mockDiscoveryNotifier{}, mockPeerNetworker{}, mockPeerInfo{})
	srv := NewDiscoveryServer(service)

	res, err := srv.Synchronize(context.TODO(), &api.SynRequest{
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
	assert.NotEmpty(t, res.UnknownPeers)
	assert.Equal(t, "pubkey", res.UnknownPeers[0].Identity.PublicKey)
}

/*
Scenario: Receive a synchronize request with knowing any peers
	Given a peer discovered
	When I receive a syn request
	Then I when I make diff , I return the a peer the sender does not known
*/
func TestHandleSynchronizeRequestByKnowingPeer(t *testing.T) {
	repo := &mockDiscoveryRepo{
		knownPeers: []discovery.Peer{
			discovery.NewDiscoveredPeer(
				discovery.NewPeerIdentity(net.ParseIP("10.0.0.1"), 3000, "pubkey000"),
				discovery.NewPeerHeartbeatState(time.Now(), 1000),
				discovery.NewPeerAppState("1.0.1", discovery.OkPeerStatus, 10.0, 20.0, "", 300, 1, 100),
			),
		},
	}

	service := discovery.NewService(repo, nil, &mockDiscoveryNotifier{}, mockPeerNetworker{}, mockPeerInfo{})
	srv := NewDiscoveryServer(service)

	res, err := srv.Synchronize(context.TODO(), &api.SynRequest{
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
	assert.NotEmpty(t, res.UnknownPeers)
	assert.Equal(t, "pubkey", res.UnknownPeers[0].Identity.PublicKey)
	assert.NotEmpty(t, res.NewPeers)
	assert.Equal(t, "pubkey000", res.NewPeers[0].Identity.PublicKey)
}

/*
Scenario: Receive an acknowledgement request with the details for the requested peers
	Given a requested peers details
	When I want to acknowledge them
	Then I store inside the db the discovered peers
*/
func TestHandleAcknowledgeRequest(t *testing.T) {
	repo := &mockDiscoveryRepo{}

	service := discovery.NewService(repo, nil, &mockDiscoveryNotifier{}, mockPeerNetworker{}, mockPeerInfo{})
	srv := NewDiscoveryServer(service)

	_, err := srv.Acknowledge(context.TODO(), &api.AckRequest{
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
					CpuLoad:               "",
					DiscoveredPeersNumber: 100,
					FreeDiskSpace:         300,
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
	assert.Len(t, repo.knownPeers, 1)
	assert.Equal(t, "pubkey", repo.knownPeers[0].Identity().PublicKey())
}

type mockDiscoveryRepo struct {
	seedPeers        []discovery.PeerIdentity
	knownPeers       []discovery.Peer
	unreachablePeers []string
}

func (r *mockDiscoveryRepo) CountKnownPeers() (int, error) {
	return len(r.knownPeers), nil
}

func (r *mockDiscoveryRepo) ListSeedPeers() ([]discovery.PeerIdentity, error) {
	return r.seedPeers, nil
}

func (r *mockDiscoveryRepo) ListKnownPeers() ([]discovery.Peer, error) {
	return r.knownPeers, nil
}

func (r *mockDiscoveryRepo) StoreKnownPeer(peer discovery.Peer) error {
	if r.containsPeer(peer) {
		for _, p := range r.knownPeers {
			if p.Identity().PublicKey() == peer.Identity().PublicKey() {
				p = peer
				break
			}
		}
	} else {
		r.knownPeers = append(r.knownPeers, peer)
	}
	return nil
}

func (r *mockDiscoveryRepo) ListReachablePeers() ([]discovery.PeerIdentity, error) {
	pp := make([]discovery.PeerIdentity, 0)
	for i := 0; i < len(r.knownPeers); i++ {
		if !r.ContainsUnreachablePeer(r.knownPeers[i].Identity().PublicKey()) {
			pp = append(pp, r.knownPeers[i].Identity())
		}
	}
	return pp, nil
}

func (r *mockDiscoveryRepo) ListUnreachablePeers() ([]discovery.PeerIdentity, error) {
	pp := make([]discovery.PeerIdentity, 0)

	for i := 0; i < len(r.seedPeers); i++ {
		if r.ContainsUnreachablePeer(r.seedPeers[i].PublicKey()) {
			pp = append(pp, r.seedPeers[i])
		}
	}

	for i := 0; i < len(r.knownPeers); i++ {
		if r.ContainsUnreachablePeer(r.knownPeers[i].Identity().PublicKey()) {
			pp = append(pp, r.knownPeers[i].Identity())
		}
	}
	return pp, nil
}

func (r *mockDiscoveryRepo) StoreSeedPeer(s discovery.PeerIdentity) error {
	r.seedPeers = append(r.seedPeers, s)
	return nil
}

func (r *mockDiscoveryRepo) StoreUnreachablePeer(pk string) error {
	if !r.ContainsUnreachablePeer(pk) {
		r.unreachablePeers = append(r.unreachablePeers, pk)
	}
	return nil
}

func (r *mockDiscoveryRepo) RemoveUnreachablePeer(pk string) error {
	if r.ContainsUnreachablePeer(pk) {
		for i := 0; i < len(r.unreachablePeers); i++ {
			if r.unreachablePeers[i] == pk {
				r.unreachablePeers = r.unreachablePeers[:i+copy(r.unreachablePeers[i:], r.unreachablePeers[i+1:])]
			}
		}
	}
	return nil
}

func (r *mockDiscoveryRepo) ContainsUnreachablePeer(peerPubk string) bool {
	for _, up := range r.unreachablePeers {
		if up == peerPubk {
			return true
		}
	}
	return false
}

func (r *mockDiscoveryRepo) containsPeer(p discovery.Peer) bool {
	mdiscoveredPeers := make(map[string]discovery.Peer, 0)
	for _, p := range r.knownPeers {
		mdiscoveredPeers[p.Identity().PublicKey()] = p
	}

	_, exist := mdiscoveredPeers[p.Identity().PublicKey()]
	return exist
}

type mockDiscoveryNotifier struct {
	reaches     []discovery.PeerIdentity
	unreaches   []discovery.PeerIdentity
	discoveries []discovery.Peer
}

func (n *mockDiscoveryNotifier) NotifyReachable(p discovery.PeerIdentity) error {
	n.reaches = append(n.reaches, p)
	return nil
}
func (n *mockDiscoveryNotifier) NotifyUnreachable(p discovery.PeerIdentity) error {
	n.unreaches = append(n.unreaches, p)
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

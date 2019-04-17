package rpc

import (
	"context"
	"crypto/rand"
	"log"
	"net"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	api "github.com/uniris/uniris-core/api/protobuf-spec"
	"github.com/uniris/uniris-core/pkg/crypto"
	"github.com/uniris/uniris-core/pkg/discovery"
	"github.com/uniris/uniris-core/pkg/logging"
	"github.com/uniris/uniris-core/pkg/shared"
)

/*
Scenario: Receive a synchronize request without knowing any peers
	Given not peers discovered
	When I receive a syn request including a peer
	Then I when I make diff , I return the sended peer as unknown
*/
func TestHandleSynchronizeRequestWithoutKnownPeers(t *testing.T) {
	db := &mockDiscoveryDB{}
	sharedDB := mochSharedDB{}

	priv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
	_, pub2, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

	sharedDB.WriteAuthorizedNode(pub)
	sharedDB.WriteAuthorizedNode(pub2)

	l := logging.NewLogger("stdout", log.New(os.Stdout, "", 0), "test", net.ParseIP("127.0.0.1"), logging.ErrorLogLevel)

	discoverySrv := NewDiscoveryServer(db, &mockDiscoveryNotifier{}, l, pub, priv, sharedDB)

	p, err := pub.Marshal()
	assert.Nil(t, err)

	sig, err := priv.Sign(p)
	assert.Nil(t, err)

	p2, err := pub2.Marshal()
	assert.Nil(t, err)

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
					PublicKey: p2,
				},
			},
		},
		PublicKey: p,
		Signature: sig,
	})

	assert.Nil(t, err)
	assert.NotEmpty(t, res.RequestedPeers)
	assert.Equal(t, p2, crypto.VersionnedKey(res.RequestedPeers[0].PublicKey))
}

/*
Scenario: Receive a synchronize request with knowing any peers
	Given a peer discovered
	When I receive a syn request
	Then I when I make diff , I return the a peer the sender does not known
*/
func TestHandleSynchronizeRequestByKnowingPeer(t *testing.T) {
	priv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
	_, pub2, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
	_, pub3, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

	p2, err := pub2.Marshal()
	assert.Nil(t, err)

	p3, err := pub3.Marshal()
	assert.Nil(t, err)

	sharedDB := mochSharedDB{}

	sharedDB.WriteAuthorizedNode(pub)
	sharedDB.WriteAuthorizedNode(pub2)
	sharedDB.WriteAuthorizedNode(pub3)

	db := &mockDiscoveryDB{
		discoveredPeers: []discovery.Peer{
			discovery.NewDiscoveredPeer(
				discovery.NewPeerIdentity(net.ParseIP("10.0.0.1"), 3000, pub2),
				discovery.NewPeerHeartbeatState(time.Now(), 1000),
				discovery.NewPeerAppState("1.0.1", discovery.OkPeerStatus, 10.0, 20.0, "", 300, 1, 100),
			),
		},
	}

	l := logging.NewLogger("stdout", log.New(os.Stdout, "", 0), "test", net.ParseIP("127.0.0.1"), logging.ErrorLogLevel)
	discoverySrv := NewDiscoveryServer(db, &mockDiscoveryNotifier{}, l, pub, priv, sharedDB)

	p, err := pub.Marshal()
	assert.Nil(t, err)

	sig, err := priv.Sign(p)
	assert.Nil(t, err)

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
					PublicKey: p3,
				},
			},
		},
		PublicKey: p,
		Signature: sig,
	})
	assert.Nil(t, err)
	assert.NotEmpty(t, res.RequestedPeers)
	assert.Equal(t, p3, crypto.VersionnedKey(res.RequestedPeers[0].PublicKey))
	assert.NotEmpty(t, res.DiscoveredPeers)
	assert.Equal(t, p2, crypto.VersionnedKey(res.DiscoveredPeers[0].Identity.PublicKey))
}

/*
Scenario: Receive an acknowledgement request with the details for the requested peers
	Given a requested peers details
	When I want to acknowledge them
	Then I store inside the db the discovered peers
*/
func TestHandleAcknowledgeRequest(t *testing.T) {
	priv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
	_, pub2, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

	p2, err := pub2.Marshal()
	assert.Nil(t, err)

	sharedDB := mochSharedDB{}
	sharedDB.WriteAuthorizedNode(pub)
	sharedDB.WriteAuthorizedNode(pub2)

	db := &mockDiscoveryDB{}

	l := logging.NewLogger("stdout", log.New(os.Stdout, "", 0), "test", net.ParseIP("127.0.0.1"), logging.ErrorLogLevel)

	discoverySrv := NewDiscoveryServer(db, &mockDiscoveryNotifier{}, l, pub, priv, sharedDB)

	p, err := pub.Marshal()
	assert.Nil(t, err)

	sig, err := priv.Sign(p)
	assert.Nil(t, err)

	_, err = discoverySrv.Acknowledge(context.TODO(), &api.AckRequest{
		RequestedPeers: []*api.PeerDiscovered{
			&api.PeerDiscovered{
				Identity: &api.PeerIdentity{
					Ip:        "127.0.0.1",
					Port:      3000,
					PublicKey: p2,
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
		PublicKey: p,
		Signature: sig,
	})

	assert.Nil(t, err)
	assert.Len(t, db.discoveredPeers, 1)
	assert.Equal(t, pub2, db.discoveredPeers[0].Identity().PublicKey())
}

type mochSharedDB struct {
	nodeCrossKeys      []shared.NodeCrossKeyPair
	emitterCrossKeys   []shared.EmitterCrossKeyPair
	authNodePublicKeys []crypto.PublicKey

	shared.KeyReadWriter
}

//EmitterCrossKeypairs retrieve the list of the cross emitter keys
func (db mochSharedDB) EmitterCrossKeypairs() ([]shared.EmitterCrossKeyPair, error) {
	return db.emitterCrossKeys, nil
}

//FirstEmitterCrossKeypair retrieves the first public key
func (db mochSharedDB) FirstEmitterCrossKeypair() (shared.EmitterCrossKeyPair, error) {
	return db.emitterCrossKeys[0], nil
}

//FirstNodeCrossKeypair retrieve the first shared crosskeys for the nodes
func (db mochSharedDB) FirstNodeCrossKeypair() (shared.NodeCrossKeyPair, error) {
	return db.nodeCrossKeys[0], nil
}

//LastNodeCrossKeypair retrieve the last shared crosskeys for the nodes
func (db mochSharedDB) LastNodeCrossKeypair() (shared.NodeCrossKeyPair, error) {
	return db.nodeCrossKeys[len(db.nodeCrossKeys)-1], nil
}

//AuthorizedNodesPublicKeys retrieves the list of public keys of the authorized nodes
func (db mochSharedDB) AuthorizedNodesPublicKeys() ([]crypto.PublicKey, error) {
	return db.authNodePublicKeys, nil
}

//WriteAuthorizedNode inserts a new node public key as an authorized node
func (db *mochSharedDB) WriteAuthorizedNode(pub crypto.PublicKey) error {
	var found bool
	for _, k := range db.authNodePublicKeys {
		if k.Equals(pub) {
			found = true
			break
		}
	}

	if !found {
		db.authNodePublicKeys = append(db.authNodePublicKeys, pub)
	}

	return nil
}

func (db mochSharedDB) IsAuthorizedNode(pub crypto.PublicKey) bool {
	found := false
	for _, k := range db.authNodePublicKeys {
		if k.Equals(pub) {
			found = true
			break
		}
	}
	return found
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

func (n *mockDiscoveryNotifier) NotifyReachable(pk crypto.PublicKey) error {
	p, err := pk.Marshal()
	if err != nil {
		return err
	}
	n.reaches = append(n.reaches, string(p))
	return nil
}
func (n *mockDiscoveryNotifier) NotifyUnreachable(pk crypto.PublicKey) error {
	p, err := pk.Marshal()
	if err != nil {
		return err
	}
	n.unreaches = append(n.unreaches, string(p))
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

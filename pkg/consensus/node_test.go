package consensus

// import (
// 	"crypto/rand"
// 	"net"
// 	"testing"

// 	"github.com/uniris/uniris-core/pkg/crypto"

// 	"github.com/stretchr/testify/assert"
// )

// /*
// Scenario: test WriteDiscoveredNode
// 	Given a node
// 	When I want write it on the nodeDB
// 	Then I get the wanted result.
// */
// func TestWriteDiscoveredNode(t *testing.T) {
// 	store := &mockNodeDatabase{}

// 	_, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
// 	_, pub2, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

// 	node := NewNode(net.ParseIP("192.168.1.1"), 5000, pub, 1, "0, 0, 0", 0, "v1.0", 1, 1, 0.000, 0.000, GeoPatch{}, false)
// 	store.WriteDiscoveredNode(node)
// 	assert.Equal(t, 1, len(store.nodes))
// 	node2 := NewNode(net.ParseIP("192.168.1.1"), 5000, pub, 1, "0, 0, 1", 0, "v1.0", 1, 1, 0.000, 0.000, GeoPatch{}, false)
// 	store.WriteDiscoveredNode(node2)
// 	assert.Equal(t, 1, len(store.nodes))
// 	node3 := NewNode(net.ParseIP("192.168.1.1"), 5000, pub2, 1, "0, 0, 0", 0, "v1.0", 1, 1, 0.000, 0.000, GeoPatch{}, false)
// 	store.WriteDiscoveredNode(node3)
// 	assert.Equal(t, 2, len(store.nodes))

// }

// type mockNodeDatabase struct {
// 	nodes []Node
// }

// func (db *mockNodeDatabase) WriteDiscoveredNode(node Node) error {
// 	for i, n := range db.nodes {
// 		if n.publicKey.Equals(node.PublicKey()) {
// 			db.nodes[i] = node
// 			return nil
// 		}
// 	}
// 	db.nodes = append(db.nodes, node)
// 	return nil
// }

// func (db *mockNodeDatabase) WriteReachableNode(publicKey crypto.PublicKey) error {
// 	for i, n := range db.nodes {
// 		if n.publicKey.Equals(publicKey) {
// 			node := NewNode(n.IP(), n.Port(), n.PublicKey(), n.Status(), n.CPULoad(), n.FreeDiskSpace(), n.Version(), n.P2PFactor(), n.ReachablePeersNumber(), n.Latitude(), n.Longitude(), n.Patch(), true)
// 			db.nodes[i] = node
// 			break
// 		}
// 	}
// 	return nil
// }

// func (db *mockNodeDatabase) WriteUnreachableNode(publicKey crypto.PublicKey) error {
// 	for i, n := range db.nodes {
// 		if n.publicKey.Equals(publicKey) {
// 			node := NewNode(n.IP(), n.Port(), n.PublicKey(), n.Status(), n.CPULoad(), n.FreeDiskSpace(), n.Version(), n.P2PFactor(), n.ReachablePeersNumber(), n.Latitude(), n.Longitude(), n.Patch(), false)
// 			db.nodes[i] = node
// 			break
// 		}
// 	}
// 	return nil
// }

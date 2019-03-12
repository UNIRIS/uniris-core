package consensus

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
Scenario: Find PatchID
	Given a latitude and logitude
	When I want to get the patch id
	Then I get the wanted result.
*/
func TestPatchId(t *testing.T) {
	lat1 := 0.0
	lon1 := 0.0
	p := ComputeGeoPatch(lat1, lon1)
	assert.NotEqual(t, 0, p.patchid)

	//position of Eiffel Tower, Paris
	lat2 := 48.8583728827653310
	lon2 := 2.2944796085357666
	//position of triumphal arch, Paris
	lat3 := 48.873804445573874
	lon3 := 2.2950267791748047

	p2 := ComputeGeoPatch(lat2, lon2)
	p3 := ComputeGeoPatch(lat3, lon3)
	assert.Equal(t, p2.patchid, p3.patchid)

	//position of statue of liberty, New york
	lat4 := 40.689039
	lon4 := -74.044396
	//position of Clock Habib Bourguiba, Tunis
	lat5 := 36.800236
	lon5 := 10.186422

	p4 := ComputeGeoPatch(lat4, lon4)
	p5 := ComputeGeoPatch(lat5, lon5)
	assert.NotEqual(t, p4.patchid, p5.patchid)
}

/*
Scenario: test WriteDiscoveredNode
	Given a node
	When I want write it on the nodeDB
	Then I get the wanted result.
*/
func TestWriteDiscoveredNode(t *testing.T) {
	store := &mockNodeDatabase{}
	node := NewNode(net.ParseIP("192.168.1.1"), 5000, "key1", 1, "0, 0, 0", 0, "v1.0", 1, 1, 0.000, 0.000, GeoPatch{}, false)
	store.WriteDiscoveredNode(node)
	assert.Equal(t, 1, len(store.nodes))
	node2 := NewNode(net.ParseIP("192.168.1.1"), 5000, "key1", 1, "0, 0, 1", 0, "v1.0", 1, 1, 0.000, 0.000, GeoPatch{}, false)
	store.WriteDiscoveredNode(node2)
	assert.Equal(t, 1, len(store.nodes))
	node3 := NewNode(net.ParseIP("192.168.1.1"), 5000, "key3", 1, "0, 0, 0", 0, "v1.0", 1, 1, 0.000, 0.000, GeoPatch{}, false)
	store.WriteDiscoveredNode(node3)
	assert.Equal(t, 2, len(store.nodes))

}

type mockNodeDatabase struct {
	nodes []Node
}

func (db *mockNodeDatabase) WriteDiscoveredNode(node Node) error {
	for i, n := range db.nodes {
		if n.PublicKey() == node.PublicKey() {
			db.nodes[i] = node
			return nil
		}
	}
	db.nodes = append(db.nodes, node)
	return nil
}

func (db *mockNodeDatabase) WriteReachableNode(publicKey string) error {
	for i, n := range db.nodes {
		if n.PublicKey() == publicKey {
			node := NewNode(n.IP(), n.Port(), n.PublicKey(), n.Status(), n.CPULoad(), n.FreeDiskSpace(), n.Version(), n.P2PFactor(), n.ReachablePeersNumber(), n.Latitude(), n.Longitude(), n.Patch(), true)
			db.nodes[i] = node
			break
		}
	}
	return nil
}

func (db *mockNodeDatabase) WriteUnreachableNode(publicKey string) error {
	for i, n := range db.nodes {
		if n.PublicKey() == publicKey {
			node := NewNode(n.IP(), n.Port(), n.PublicKey(), n.Status(), n.CPULoad(), n.FreeDiskSpace(), n.Version(), n.P2PFactor(), n.ReachablePeersNumber(), n.Latitude(), n.Longitude(), n.Patch(), false)
			db.nodes[i] = node
			break
		}
	}
	return nil
}

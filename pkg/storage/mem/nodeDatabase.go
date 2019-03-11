package memstorage

import (
	"reflect"

	"github.com/uniris/uniris-core/pkg/consensus"
)

//NodeDatabase is a network nodes database in memory
type NodeDatabase struct {
	nodes []consensus.Node

	consensus.NodeReader
	consensus.NodeWriter
}

//WriteDiscoveredNode stores a new discovered node
func (db *NodeDatabase) WriteDiscoveredNode(node consensus.Node) error {
	for i, n := range db.nodes {
		if n.PublicKey() == node.PublicKey() && !reflect.DeepEqual(n, node) {
			db.nodes[i] = node
			return nil
		}
	}
	db.nodes = append(db.nodes, node)
	return nil
}

//WriteReachableNode defines a node by its public key as reachable
func (db *NodeDatabase) WriteReachableNode(publicKey string) error {
	for i, n := range db.nodes {
		if n.PublicKey() == publicKey {
			node := consensus.NewNode(n.IP(), n.Port(), n.PublicKey(), n.Status(), n.CPULoad(), n.FreeDiskSpace(), n.Version(), n.P2PFactor(), n.ReachablePeersNumber(), n.Latitude(), n.Longitude(), n.Patch(), true)
			db.nodes[i] = node
			break
		}
	}
	return nil
}

//WriteUnreachableNode defines a node by its public key as unreachable
func (db *NodeDatabase) WriteUnreachableNode(publicKey string) error {
	for i, n := range db.nodes {
		if n.PublicKey() == publicKey {
			node := consensus.NewNode(n.IP(), n.Port(), n.PublicKey(), n.Status(), n.CPULoad(), n.FreeDiskSpace(), n.Version(), n.P2PFactor(), n.ReachablePeersNumber(), n.Latitude(), n.Longitude(), n.Patch(), false)
			db.nodes[i] = node
			break
		}
	}
	return nil
}

package memstorage

import (
	"github.com/uniris/uniris-core/pkg/consensus"
)

//NodeDatabase is a network nodes database in memory
type NodeDatabase struct {
	nodes []consensus.Node

	consensus.NodeReader
	consensus.NodeWriter
}

//WriteDiscoveredNode stores a new discovered node
func (db *NodeDatabase) WriteDiscoveredNode(n consensus.Node) error {
	db.nodes = append(db.nodes, n)
	return nil
}

//WriteReachableNode defines a node by its public key as reachable
func (db *NodeDatabase) WriteReachableNode(publicKey string) error {
	for i, n := range db.nodes {
		if n.PublicKey() == publicKey {
			db.nodes[i] = consensus.NewNode(n.IP(), n.Port(), n.PublicKey(), n.Status(), n.CPULoad(), n.FreeDiskSpace(), n.Version(), n.P2PFactor(), n.ReachablePeersNumber(), n.Patch(), true)
			break
		}
	}
	return nil
}

//WriteUnreachableNode defines a node by its public key as unreachable
func (db *NodeDatabase) WriteUnreachableNode(publicKey string) error {
	for i, n := range db.nodes {
		if n.PublicKey() == publicKey {
			db.nodes[i] = consensus.NewNode(n.IP(), n.Port(), n.PublicKey(), n.Status(), n.CPULoad(), n.FreeDiskSpace(), n.Version(), n.P2PFactor(), n.ReachablePeersNumber(), n.Patch(), false)
			break
		}
	}
	return nil
}

package memstorage

import (
	"github.com/uniris/uniris-core/pkg/crypto"

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
		if n.PublicKey() == node.PublicKey() {
			db.nodes[i] = node
			return nil
		}
	}
	db.nodes = append(db.nodes, node)
	return nil
}

//WriteReachableNode defines a node by its public key as reachable
func (db *NodeDatabase) WriteReachableNode(publicKey crypto.PublicKey) error {
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
func (db *NodeDatabase) WriteUnreachableNode(publicKey crypto.PublicKey) error {
	for i, n := range db.nodes {
		if n.PublicKey() == publicKey {
			node := consensus.NewNode(n.IP(), n.Port(), n.PublicKey(), n.Status(), n.CPULoad(), n.FreeDiskSpace(), n.Version(), n.P2PFactor(), n.ReachablePeersNumber(), n.Latitude(), n.Longitude(), n.Patch(), false)
			db.nodes[i] = node
			break
		}
	}
	return nil
}

//Reachables retrieves the nodes flagged as reachable
func (db NodeDatabase) Reachables() (reachables []consensus.Node, err error) {
	for _, n := range db.nodes {
		if n.IsReachable() {
			reachables = append(reachables, n)
		}
	}
	return
}

//Unreachables retrieves the nodes flagged as unreachable
func (db NodeDatabase) Unreachables() (unreachables []consensus.Node, err error) {
	for _, n := range db.nodes {
		if !n.IsReachable() {
			unreachables = append(unreachables, n)
		}
	}
	return
}

//CountReachables retrieves the number of reachable nodes
func (db NodeDatabase) CountReachables() (nb int, err error) {
	for _, n := range db.nodes {
		if n.IsReachable() {
			nb++
		}
	}
	return
}

//FindByPublicKey retrieves a node from a public key
func (db *NodeDatabase) FindByPublicKey(publicKey string) (found consensus.Node, err error) {
	for _, n := range db.nodes {
		if n.PublicKey() == publicKey {
			return n, nil
		}
	}
	return
}

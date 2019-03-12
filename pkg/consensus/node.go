package consensus

import (
	"net"
)

//NodeWriter persists network node
type NodeWriter interface {

	//WriteDiscoveredNode stores a new discovered node
	WriteDiscoveredNode(n Node) error

	//WriteReachableNode defines a node by its public key as reachable
	WriteReachableNode(publicKey string) error

	//WriteUnreachableNode defines a node by its public key as unreachable
	WriteUnreachableNode(publicKey string) error
}

//NodeReader provides queries to fetch network nodes
type NodeReader interface {
	//CountReachables retrieves the number of reachable nodes
	CountReachables() (int, error)

	//Reachables retrieves the nodes flagged as reachable
	Reachables() ([]Node, error)

	//Unreachables retrieves the nodes flagged as unreachable
	Unreachables() ([]Node, error)

	//FindByPublicKey retrieves a node from a public key
	FindByPublicKey(publicKey string) (Node, error)
}

//Node represents a discovered peers with some additional computed data
type Node struct {
	ip                   net.IP
	port                 int
	publicKey            string
	status               NodeStatus
	cpuLoad              string
	freeDiskSpace        float64
	version              string
	p2pFactor            int
	reachablePeersNumber int
	latitude             float64
	longitude            float64
	patch                GeoPatch
	isReachable          bool
}

//NodeStatus identifies the status of a node
type NodeStatus int

const (
	//NodeBootstraping identifies a node which is not fully synchronized and ready
	NodeBootstraping NodeStatus = iota

	//NodeOK identifies the peer as stable
	NodeOK

	//NodeFaulty identifies a peer with some errors (NTP drifts, DNS unfound, GeoPosition failed, GRPC not running)
	NodeFaulty
)

//NewNode creates a new enhanced discovered peers with geo patch
func NewNode(ip net.IP, port int, pubK string, status NodeStatus, cpu string, disk float64, ver string, p2pFactor int, reachNumbers int, lat float64, lon float64, patch GeoPatch, isReachable bool) Node {
	return Node{
		ip:                   ip,
		port:                 port,
		publicKey:            pubK,
		status:               status,
		cpuLoad:              cpu,
		freeDiskSpace:        disk,
		version:              ver,
		p2pFactor:            p2pFactor,
		reachablePeersNumber: reachNumbers,
		latitude:             lat,
		longitude:            lon,
		patch:                patch,
		isReachable:          isReachable,
	}
}

//IP returns the node's ip
func (n Node) IP() net.IP {
	return n.ip
}

//Port returns the node's port
func (n Node) Port() int {
	return n.port
}

//PublicKey returns the node's public key
func (n Node) PublicKey() string {
	return n.publicKey
}

//Status returns the node's status
func (n Node) Status() NodeStatus {
	return n.status
}

//CPULoad returns the node's cpu load
func (n Node) CPULoad() string {
	return n.cpuLoad
}

//FreeDiskSpace returns the node's free disk space
func (n Node) FreeDiskSpace() float64 {
	return n.freeDiskSpace
}

//Version returns the node's version
func (n Node) Version() string {
	return n.version
}

//ReachablePeersNumber returns the number of peers has been discovered and are reachable
func (n Node) ReachablePeersNumber() int {
	return n.reachablePeersNumber
}

//P2PFactor returns the node's p2p factor
func (n Node) P2PFactor() int {
	return n.p2pFactor
}

//Longitude returns the node's longitude coordinates
func (n Node) Longitude() float64 {
	return n.longitude
}

//Latitude returns the node's latitude coordinates
func (n Node) Latitude() float64 {
	return n.latitude
}

//Patch returns the geo patch of the peer
func (n Node) Patch() GeoPatch {
	return n.patch
}

//IsReachable returns true if the node has been reached recently
func (n Node) IsReachable() bool {
	return n.isReachable
}

//GeoPatch represents a geographic section on the earth based on latitude and longitude
type GeoPatch struct {
	patchid int
	left    float64
	right   float64
	top     float64
	bottom  float64
}

//ID returns the geo patch ID
func (p GeoPatch) ID() int {
	return p.patchid
}

func createMapPatches(xDegree float64, yDegree float64) []GeoPatch {

	geoPatches := make([]GeoPatch, 0)
	i := 0
	for x := -180.0; x < 180.0; x += xDegree {
		for y := -90.0; y < 90.0; y += yDegree {
			geoPatches = append(geoPatches, GeoPatch{
				patchid: i,
				left:    x,
				right:   x + 10,
				bottom:  y,
				top:     y + 10,
			})
			i++
		}
	}
	return geoPatches
}

//ComputeGeoPatch identifies a geographic patch from a given geographic peer position
func ComputeGeoPatch(lat float64, lon float64) (p GeoPatch) {
	geoPatches := createMapPatches(10, 10)

	for _, patch := range geoPatches {

		if lon >= patch.left && lon <= patch.right && lat >= patch.bottom && lat <= patch.top {
			p = patch
			break
		}
	}

	return p
}

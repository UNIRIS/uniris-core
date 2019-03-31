package redis

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net"
	"strconv"

	"github.com/go-redis/redis"
	"github.com/uniris/uniris-core/pkg/consensus"
	"github.com/uniris/uniris-core/pkg/crypto"
)

const (
	nodesKey = "nodes"
)

type nodeDB struct {
	client *redis.Client
}

//NewNodeDatabase creates a new node database handler using redis as storage
//An error is returned if the redis instance is not reached
func NewNodeDatabase(hostname string, port int, pwd string) (consensus.NodeReadWriter, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", hostname, port),
		Password: pwd,
		DB:       1,
	})

	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}

	return &nodeDB{
		client: client,
	}, nil
}

func (db nodeDB) CountReachables() (int, error) {
	r, err := db.Reachables()
	if err != nil {
		return 0, err
	}
	return len(r), nil
}

func (db nodeDB) FindByPublicKey(pubKey crypto.PublicKey) (n consensus.Node, err error) {

	pubB, err := pubKey.Marshal()
	if err != nil {
		return
	}
	pubKeyEnc := base64.StdEncoding.EncodeToString(pubB)

	list, err := db.fetchList(nodesKey)
	if err != nil {
		return
	}
	for _, r := range list {
		if r["publicKey"] == pubKeyEnc {
			n, err := db.formatHashToNode(r)
			if err != nil {
				return consensus.Node{}, err
			}
			return n, nil
		}
	}

	return consensus.Node{}, errors.New("node not found")
}

func (db nodeDB) Reachables() (nodes []consensus.Node, err error) {
	list, err := db.fetchList(nodesKey)
	if err != nil {
		return nil, err
	}
	for _, r := range list {
		if r["isReachable"] == "1" {
			n, err := db.formatHashToNode(r)
			if err != nil {
				return nil, err
			}
			nodes = append(nodes, n)
		}
	}
	return
}

func (db nodeDB) Unreachables() (nodes []consensus.Node, err error) {
	list, err := db.fetchList(nodesKey)
	if err != nil {
		return nil, err
	}
	for _, r := range list {
		if r["isReachable"] == "0" {
			n, err := db.formatHashToNode(r)
			if err != nil {
				return nil, err
			}
			nodes = append(nodes, n)
		}
	}
	return
}

func (db nodeDB) WriteDiscoveredNode(n consensus.Node) error {
	pubB, err := n.PublicKey().Marshal()
	if err != nil {
		return err
	}
	pubKeyEnc := base64.StdEncoding.EncodeToString(pubB)
	id := fmt.Sprintf("%s:%s", nodesKey, pubKeyEnc)
	nHash, err := db.formatNodeToHash(n)
	if err != nil {
		return err
	}
	cmd := db.client.HMSet(id, nHash)
	if cmd.Err() != nil {
		return cmd.Err()
	}
	return nil
}

func (db nodeDB) WriteReachableNode(publicKey crypto.PublicKey) error {
	pubB, err := publicKey.Marshal()
	if err != nil {
		return err
	}
	pubKeyEnc := base64.StdEncoding.EncodeToString(pubB)
	id := fmt.Sprintf("%s:%s", nodesKey, pubKeyEnc)

	cmd := db.client.HExists(id, "publicKey")
	if cmd.Err() != nil {
		return cmd.Err()
	}
	exist, err := cmd.Result()
	if err != nil {
		return err
	}
	if exist {
		cmd := db.client.HMSet(id, map[string]interface{}{
			"isReachable": 1,
		})
		if cmd.Err() != nil {
			return cmd.Err()
		}
	}

	return nil
}

func (db nodeDB) WriteUnreachableNode(publicKey crypto.PublicKey) error {
	pubB, err := publicKey.Marshal()
	if err != nil {
		return err
	}
	pubKeyEnc := base64.StdEncoding.EncodeToString(pubB)
	id := fmt.Sprintf("%s:%s", nodesKey, pubKeyEnc)

	cmd := db.client.HExists(id, "publicKey")
	if cmd.Err() != nil {
		return cmd.Err()
	}
	exist, err := cmd.Result()
	if err != nil {
		return err
	}
	if exist {
		cmd := db.client.HMSet(id, map[string]interface{}{
			"isReachable": 0,
		})
		if cmd.Err() != nil {
			return cmd.Err()
		}
	}

	return nil
}

func (db nodeDB) fetchList(key string) ([]map[string]string, error) {
	cmdKeys := db.client.Keys(key)
	if cmdKeys.Err() != nil {
		return nil, cmdKeys.Err()
	}
	keys, err := cmdKeys.Result()
	if err != nil {
		return nil, err
	}

	list := make([]map[string]string, 0)
	for _, k := range keys {
		cmdGet := db.client.HGetAll(k)
		if cmdGet.Err() != nil {
			return nil, cmdGet.Err()
		}
		res, err := cmdGet.Result()
		if err != nil {
			return nil, err
		}
		list = append(list, res)
	}
	return list, nil
}

func (db nodeDB) formatNodeToHash(n consensus.Node) (map[string]interface{}, error) {

	pubBytes, err := n.PublicKey().Marshal()
	if err != nil {
		return nil, err
	}

	var isReachable int
	if n.IsReachable() {
		isReachable = 1
	}

	return map[string]interface{}{
		"publicKey":            base64.StdEncoding.EncodeToString(pubBytes),
		"ip":                   n.IP().String(),
		"port":                 strconv.Itoa(n.Port()),
		"status":               strconv.Itoa(int(n.Status())),
		"cpuLoad":              n.CPULoad(),
		"freeDiskSpace":        fmt.Sprintf("%f", n.FreeDiskSpace()),
		"version":              n.Version(),
		"p2pFactor":            strconv.Itoa(n.P2PFactor()),
		"reachablePeersNumber": strconv.Itoa(n.ReachablePeersNumber()),
		"latitude":             fmt.Sprintf("%f", n.Latitude()),
		"longitude":            fmt.Sprintf("%f", n.Longitude()),
		"patchID":              strconv.Itoa(n.Patch().ID()),
		"isReachable":          strconv.Itoa(isReachable),
	}, nil
}

func (db nodeDB) formatHashToNode(h map[string]string) (n consensus.Node, err error) {
	port, err := strconv.Atoi(h["port"])
	if err != nil {
		return
	}
	ip := net.ParseIP(h["ip"])

	pbKeyEncoded := h["publicKey"]
	pubBytes, err := base64.StdEncoding.DecodeString(pbKeyEncoded)
	if err != nil {
		return
	}

	pubK, err := crypto.ParsePublicKey(pubBytes)
	if err != nil {
		return
	}

	status, err := strconv.Atoi(h["status"])
	if err != nil {
		return
	}

	cpuLoad := h["cpuLoad"]
	freeDiskSpace, err := strconv.ParseFloat(h["freeDiskSpace"], 64)
	if err != nil {
		return
	}
	version := h["version"]
	p2pFactor, err := strconv.Atoi(h["p2pFactor"])
	if err != nil {
		return
	}
	reachablesNumbers, err := strconv.Atoi(h["reachablePeersNumber"])
	if err != nil {
		return
	}
	lat, err := strconv.ParseFloat(h["latitude"], 64)
	if err != nil {
		return
	}
	lon, err := strconv.ParseFloat(h["longitude"], 64)
	if err != nil {
		return
	}

	patchID, err := strconv.Atoi(h["patchID"])
	if err != nil {
		return
	}
	patch, err := consensus.NewGeoPatch(patchID)
	if err != nil {
		return
	}

	isReachable, err := strconv.Atoi(h["isReachable"])
	if err != nil {
		return
	}

	return consensus.NewNode(ip, port, pubK, consensus.NodeStatus(status), cpuLoad,
		freeDiskSpace, version, p2pFactor, reachablesNumbers, lat, lon, patch, isReachable != 0), nil
}

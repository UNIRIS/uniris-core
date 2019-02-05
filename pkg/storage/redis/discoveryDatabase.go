package redis

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/uniris/uniris-core/pkg/discovery"

	"github.com/go-redis/redis"
)

const (
	unreachablesKey = "unreachabled-peers"
	discoveriesKey  = "discovered-peer"
)

type discoveryDb struct {
	client *redis.Client
}

//NewDiscoveryDatabase creates a new repository using redis as storage
//An error is returned if the redis instance is not reached
func NewDiscoveryDatabase(hostname string, port int, pwd string) (discovery.Database, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", hostname, port),
		Password: pwd,
		DB:       0,
	})

	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}

	return &discoveryDb{
		client: client,
	}, nil
}

func (r discoveryDb) DiscoveredPeers() ([]discovery.Peer, error) {
	list, err := r.fetchList(fmt.Sprintf("%s:*", discoveriesKey))
	if err != nil {
		return nil, err
	}

	peers := make([]discovery.Peer, 0)
	for _, p := range list {
		peer := formatHashToPeer(p)
		peers = append(peers, peer)
	}

	return peers, nil
}

func (r discoveryDb) WriteDiscoveredPeer(p discovery.Peer) error {
	id := fmt.Sprintf("%s:%s", discoveriesKey, p.Identity().PublicKey())
	cmd := r.client.HMSet(id, formatPeerToHash(p))
	if cmd.Err() != nil {
		return cmd.Err()
	}
	return nil
}

func (r discoveryDb) WriteUnreachablePeer(pID discovery.PeerIdentity) error {
	id := fmt.Sprintf("%s:%s", unreachablesKey, pID.PublicKey())
	cmd := r.client.HMSet(id, formatPeerIdentityToHash(pID))
	if cmd.Err() != nil {
		return cmd.Err()
	}

	return nil
}

func (r discoveryDb) UnreachablePeers() ([]discovery.PeerIdentity, error) {
	list, err := r.fetchList(fmt.Sprintf("%s:*", unreachablesKey))
	if err != nil {
		return nil, err
	}

	peers := make([]discovery.PeerIdentity, 0)
	for _, unr := range list {
		pID := formatHashToPeerIdentity(unr)
		peers = append(peers, pID)
	}
	return peers, nil
}

func (r discoveryDb) ContainsUnreachablePeer(pID discovery.PeerIdentity) (bool, error) {
	id := fmt.Sprintf("%s:%s", unreachablesKey, pID.PublicKey())
	cmd := r.client.HKeys(id)
	if cmd.Err() != nil {
		return false, cmd.Err()
	}
	res, err := cmd.Result()
	if err != nil {
		return false, err
	}

	return len(res) > 0, nil
}

func (r discoveryDb) RemoveUnreachablePeer(pID discovery.PeerIdentity) error {
	id := fmt.Sprintf("%s:%s", unreachablesKey, pID.PublicKey())

	exist, err := r.ContainsUnreachablePeer(pID)
	if err != nil {
		return err
	}
	if exist {
		cmd := r.client.Del(id)
		if cmd.Err() != nil {
			return cmd.Err()
		}
		return nil
	}
	return nil
}

func (r discoveryDb) fetchList(key string) ([]map[string]string, error) {
	cmdKeys := r.client.Keys(key)
	if cmdKeys.Err() != nil {
		return nil, cmdKeys.Err()
	}
	keys, err := cmdKeys.Result()
	if err != nil {
		return nil, err
	}

	list := make([]map[string]string, 0)
	for _, k := range keys {
		cmdGet := r.client.HGetAll(k)
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

func formatPeerToHash(p discovery.Peer) map[string]interface{} {
	return map[string]interface{}{
		"publicKey":             p.Identity().PublicKey(),
		"port":                  strconv.Itoa(p.Identity().Port()),
		"ip":                    p.Identity().IP().String(),
		"generationTime":        strconv.Itoa(int(p.HeartbeatState().GenerationTime().Unix())),
		"elapsedHeartbeats":     strconv.Itoa(int(p.HeartbeatState().ElapsedHeartbeats())),
		"status":                string(p.AppState().Status()),
		"cpuLoad":               p.AppState().CPULoad(),
		"freeDiskSpace":         fmt.Sprintf("%f", p.AppState().FreeDiskSpace()),
		"version":               p.AppState().Version(),
		"geoPosition":           fmt.Sprintf("%f;%f", p.AppState().GeoPosition().Latitude(), p.AppState().GeoPosition().Longitude()),
		"p2pFactor":             string(p.AppState().P2PFactor()),
		"discoveredPeersNumber": fmt.Sprintf("%d", p.AppState().DiscoveredPeersNumber()),
	}
}

func formatPeerIdentityToHash(pID discovery.PeerIdentity) map[string]interface{} {
	return map[string]interface{}{
		"publicKey": pID.PublicKey(),
		"port":      strconv.Itoa(pID.Port()),
		"ip":        pID.IP().String(),
	}
}

func formatHashToPeerIdentity(hash map[string]string) discovery.PeerIdentity {
	pbKey := hash["publicKey"]
	port, _ := strconv.Atoi(hash["port"])
	ip := net.ParseIP(hash["ip"])

	return discovery.NewPeerIdentity(ip, port, pbKey)
}

func formatHashToPeer(hash map[string]string) discovery.Peer {

	pbKey := hash["publicKey"]
	port, _ := strconv.Atoi(hash["port"])
	ip := net.ParseIP(hash["ip"])

	gen, _ := strconv.Atoi(hash["generationTime"])
	generationTime := time.Unix(int64(gen), 0)

	elapsedHeartbeats, _ := strconv.Atoi(hash["elapsedHeartbeats"])
	elpased := int64(elapsedHeartbeats)

	s, _ := strconv.Atoi(hash["status"])
	status := discovery.PeerStatus(s)

	cpuLoad := hash["cpuLoad"]
	freeDiskSpace, _ := strconv.ParseFloat(hash["freeDiskSpace"], 64)
	version := hash["version"]
	p2pFactor, _ := strconv.Atoi(hash["p2pFactor"])
	posArr := strings.Split(hash["geoPosition"], ";")

	lat, _ := strconv.ParseFloat(posArr[0], 64)
	lon, _ := strconv.ParseFloat(posArr[1], 64)

	dpn, _ := strconv.Atoi(hash["discoveredPeersNumber"])

	p := discovery.NewDiscoveredPeer(
		discovery.NewPeerIdentity(ip, port, pbKey),
		discovery.NewPeerHeartbeatState(generationTime, elpased),
		discovery.NewPeerAppState(version, status, lat, lon, cpuLoad, freeDiskSpace, p2pFactor, dpn),
	)
	return p
}

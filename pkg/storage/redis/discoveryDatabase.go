package redis

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/uniris/uniris-core/pkg/crypto"
	"github.com/uniris/uniris-core/pkg/discovery"

	"github.com/go-redis/redis"
)

const (
	unreachablesKey = "unreachabled-peers"
	discoveriesKey  = "discovered-peer"
)

var (
	errorPubKeyParsePeerToHash             = errors.New("Unable to format peer to a hash, due to a problem on the public key")
	errorPubKeyParsePeerIDentityToHash     = errors.New("Unable to format peer Identity to a hash, due to a problem on the public key")
	errorPubKeyParseHashToPeer             = errors.New("Unable to format hash to a peer, due to a problem on the public key")
	errorPubKeyParseHashToPeerIdentity     = errors.New("Unable to format hash to a peer Identity, due to a problem on the public key")
	errorPubKeyWriteDiscoveredPeerToDB     = errors.New("Unable to write peer on db, due to a problem on the public key")
	errorPubKeyWriteUnreachablePeerToDB    = errors.New("Unable to write unreacheable peer on db, due to a problem on the public key")
	errorPubKeyRemoveUnreachablePeerFromDB = errors.New("Unable to Remove unreacheable peer from db, due to a problem on the public key")
)

type discoveryDb struct {
	client *redis.Client
}

//NewDiscoveryDatabase creates a new discovery database handler using redis as storage
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
		fpeer, err := formatHashToPeer(p)
		if err != nil {
			return nil, err
		}
		peers = append(peers, fpeer)
	}

	return peers, nil
}

func (r discoveryDb) WriteDiscoveredPeer(p discovery.Peer) error {
	pk, err := p.Identity().PublicKey().Marshal()
	if err != nil {
		return errorPubKeyWriteDiscoveredPeerToDB
	}
	id := fmt.Sprintf("%s:%s", discoveriesKey, string(pk))

	fp, err := formatPeerToHash(p)
	if err != nil {
		return err
	}

	cmd := r.client.HMSet(id, fp)
	if cmd.Err() != nil {
		return cmd.Err()
	}
	return nil
}

func (r discoveryDb) WriteUnreachablePeer(pID discovery.PeerIdentity) error {
	pk, err := pID.PublicKey().Marshal()
	if err != nil {
		return errorPubKeyWriteUnreachablePeerToDB
	}
	id := fmt.Sprintf("%s:%s", unreachablesKey, string(pk))

	fpid, err := formatPeerIdentityToHash(pID)
	if err != nil {
		return err
	}
	cmd := r.client.HMSet(id, fpid)
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
		funr, err := formatHashToPeerIdentity(unr)
		if err != nil {
			return nil, err
		}
		peers = append(peers, funr)
	}
	return peers, nil
}

func (r discoveryDb) RemoveUnreachablePeer(pID discovery.PeerIdentity) error {
	pk, err := pID.PublicKey().Marshal()
	if err != nil {
		return errorPubKeyRemoveUnreachablePeerFromDB
	}
	id := fmt.Sprintf("%s:%s", unreachablesKey, string(pk))

	cmd := r.client.Del(id)
	if cmd.Err() != nil {
		return cmd.Err()
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

func formatPeerToHash(p discovery.Peer) (map[string]interface{}, error) {
	pk, err := p.Identity().PublicKey().Marshal()
	if err != nil {
		return nil, errorPubKeyParsePeerToHash
	}

	return map[string]interface{}{
		"publicKey":            string(pk),
		"port":                 strconv.Itoa(p.Identity().Port()),
		"ip":                   p.Identity().IP().String(),
		"generationTime":       strconv.Itoa(int(p.HeartbeatState().GenerationTime().Unix())),
		"elapsedHeartbeats":    strconv.Itoa(int(p.HeartbeatState().ElapsedHeartbeats())),
		"status":               string(p.AppState().Status()),
		"cpuLoad":              p.AppState().CPULoad(),
		"freeDiskSpace":        fmt.Sprintf("%f", p.AppState().FreeDiskSpace()),
		"version":              p.AppState().Version(),
		"geoPosition":          fmt.Sprintf("%f;%f", p.AppState().GeoPosition().Latitude(), p.AppState().GeoPosition().Longitude()),
		"p2pFactor":            string(p.AppState().P2PFactor()),
		"reachablePeersNumber": fmt.Sprintf("%d", p.AppState().ReachablePeersNumber()),
	}, nil
}

func formatPeerIdentityToHash(pID discovery.PeerIdentity) (map[string]interface{}, error) {
	pk, err := pID.PublicKey().Marshal()
	if err != nil {
		return nil, errorPubKeyParsePeerIDentityToHash
	}
	return map[string]interface{}{
		"publicKey": string(pk),
		"port":      strconv.Itoa(pID.Port()),
		"ip":        pID.IP().String(),
	}, nil
}

func formatHashToPeerIdentity(hash map[string]string) (pid discovery.PeerIdentity, err error) {
	pbKey := hash["publicKey"]
	port, _ := strconv.Atoi(hash["port"])
	ip := net.ParseIP(hash["ip"])

	pk, err := crypto.ParsePublicKey([]byte(pbKey))
	if err != nil {
		return pid, errorPubKeyParseHashToPeerIdentity
	}

	return discovery.NewPeerIdentity(ip, port, pk), nil
}

func formatHashToPeer(hash map[string]string) (peer discovery.Peer, err error) {

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

	rpn, _ := strconv.Atoi(hash["reachablePeersNumber"])

	pk, err := crypto.ParsePublicKey([]byte(pbKey))
	if err != nil {
		return peer, errorPubKeyParseHashToPeer
	}

	peer = discovery.NewDiscoveredPeer(
		discovery.NewPeerIdentity(ip, port, pk),
		discovery.NewPeerHeartbeatState(generationTime, elpased),
		discovery.NewPeerAppState(version, status, lat, lon, cpuLoad, freeDiskSpace, p2pFactor, rpn),
	)
	return peer, nil
}

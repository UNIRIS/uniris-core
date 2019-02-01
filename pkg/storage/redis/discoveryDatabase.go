package redis

import (
	"fmt"
	"sort"

	"github.com/uniris/uniris-core/pkg/discovery"

	"github.com/go-redis/redis"
)

type discoveryDb struct {
	client *redis.Client
}

const (
	unreachablesKey = "unreachabled-peers"
	peerKey         = "discovered-peer"
	seedKey         = "seed-peer"
)

//ListSeedPeers return all the seed on the repository
func (r discoveryDb) ListSeedPeers() ([]discovery.PeerIdentity, error) {
	id := fmt.Sprintf("%s:*", seedKey)
	list, err := r.fetchList(id)
	if err != nil {
		return nil, err
	}

	seeds := make([]discovery.PeerIdentity, 0)
	for _, s := range list {
		seed := FormatHashToSeed(s)
		seeds = append(seeds, seed)
	}

	return seeds, nil
}

//ListKnownPeers returns all the discoveredPeers on the repository
func (r discoveryDb) ListKnownPeers() ([]discovery.Peer, error) {
	list, err := r.fetchList(fmt.Sprintf("%s:*", peerKey))
	if err != nil {
		return nil, err
	}

	peers := make([]discovery.Peer, 0)
	for _, p := range list {
		peer := FormatHashToPeer(p)
		peers = append(peers, peer)
	}

	return peers, nil
}

//SetKnownPeer add a peer to the repository
func (r discoveryDb) StoreKnownPeer(p discovery.Peer) error {
	id := fmt.Sprintf("%s:%s", peerKey, p.Identity().PublicKey())
	cmd := r.client.HMSet(id, FormatPeerToHash(p))
	if cmd.Err() != nil {
		return cmd.Err()
	}
	return nil
}

//StoreSeedPeer add a seed to the repository
func (r discoveryDb) StoreSeedPeer(s discovery.PeerIdentity) error {
	id := fmt.Sprintf("%s:%s", seedKey, s.PublicKey())
	cmd := r.client.HMSet(id, FormatSeedToHash(s))
	if cmd.Err() != nil {
		return cmd.Err()
	}
	return nil
}

//CountKnownPeers return the number of Known peers
func (r discoveryDb) CountKnownPeers() (int, error) {
	cmdKeys := r.client.Keys(fmt.Sprintf("%s:*", peerKey))
	if cmdKeys.Err() != nil {
		return 0, cmdKeys.Err()
	}
	keys, err := cmdKeys.Result()
	if err != nil {
		return 0, err
	}
	return len(keys), nil
}

//StoreUnreachablePeer add an unreachable peer to the repository
func (r discoveryDb) StoreUnreachablePeer(pbKey string) error {
	cmd := r.client.SAdd(unreachablesKey, pbKey)
	if cmd.Err() != nil {
		return cmd.Err()
	}
	return nil
}

//RemoveUnreachablePeer remove an unreachable peer to the repository
func (r discoveryDb) RemoveUnreachablePeer(pbKey string) error {
	boolCmd := r.client.SIsMember(unreachablesKey, pbKey)
	if boolCmd.Err() != nil {
		return boolCmd.Err()
	}

	exists, err := boolCmd.Result()
	if err != nil {
		return err
	}

	if exists {
		cmd := r.client.SRem(unreachablesKey, pbKey, 0)
		if cmd.Err() != nil {
			return cmd.Err()
		}
	}
	return nil
}

//ListReachablePeers returns all the reachable peers on the repository
func (r discoveryDb) ListReachablePeers() ([]discovery.PeerIdentity, error) {
	peers, err := r.ListKnownPeers()
	if err != nil {
		return nil, err
	}

	unreachableKeys, err := r.listUnreachableKeys()
	if err != nil {
		return nil, err
	}

	peerIds := make([]discovery.PeerIdentity, 0)
	for _, p := range peers {
		peerIds = append(peerIds, p.Identity())
	}

	//When there is no unreachables returns the list of known peers
	if len(unreachableKeys) == 0 {
		return peerIds, nil
	}

	//We want to get the peers no include inside the list of unreachables
	pp := make([]discovery.PeerIdentity, 0)
	sort.Strings(unreachableKeys)
	for _, p := range peerIds {
		idx := sort.SearchStrings(unreachableKeys, p.PublicKey())
		if idx >= len(unreachableKeys) || p.PublicKey() != unreachableKeys[idx] {
			pp = append(pp, p)
		}
	}

	return pp, nil
}

//ListUnreacheablePeers returns all unreachable peers on the repository
func (r discoveryDb) ListUnreachablePeers() ([]discovery.PeerIdentity, error) {
	unreachableKeys, err := r.listUnreachableKeys()
	if err != nil {
		return nil, err
	}
	peers, err := r.ListKnownPeers()
	if err != nil {
		return nil, err
	}

	pp := make([]discovery.PeerIdentity, 0)

	//Avoid looping if there is not unreachable keys
	if len(unreachableKeys) == 0 {
		return pp, nil
	}

	sort.Strings(unreachableKeys)

	for i, p := range peers {
		idx := sort.SearchStrings(unreachableKeys, p.Identity().PublicKey())
		if idx < len(unreachableKeys) && peers[i].Identity().PublicKey() == unreachableKeys[idx] {
			pp = append(pp, peers[i].Identity())
		}
	}

	return pp, nil
}

//ContainsUnreachableKey check if the pubk is in the list of unreacheable keys
func (r discoveryDb) ContainsUnreachablePeer(pubk string) bool {
	unreachableKeys, err := r.listUnreachableKeys()
	if err != nil {
		return false
	}
	sort.Strings(unreachableKeys)
	idx := sort.SearchStrings(unreachableKeys, pubk)
	if idx < len(unreachableKeys) && unreachableKeys[idx] == pubk {
		return true
	}
	return false
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

func (r discoveryDb) listUnreachableKeys() ([]string, error) {
	cmd := r.client.SMembers(unreachablesKey)
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}

	res, err := cmd.Result()
	if err != nil {
		return nil, err
	}
	return res, nil
}

//NewDiscoveryDatabase creates a new repository using redis as storage
//
//An error is returned if the redis instance is not reached
func NewDiscoveryDatabase(hostname string, port int, pwd string) (discovery.Repository, error) {
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

package redis

import (
	"fmt"
	"net"
	"sort"

	"github.com/go-redis/redis"

	discovery "github.com/uniris/uniris-core/autodiscovery/pkg"
	gossip "github.com/uniris/uniris-core/autodiscovery/pkg/gossip"
	"github.com/uniris/uniris-core/autodiscovery/pkg/system"
)

type redisRepository struct {
	client *redis.Client
}

const (
	unreachablesKey = "unreachabled-peers"
	peerKey         = "discovered-peer"
	seedKey         = "seed-peer"
)

//GetOwnedPeer return the local peer
func (r redisRepository) GetOwnedPeer() (discovery.Peer, error) {
	peers, err := r.ListKnownPeers()
	if err != nil {
		return nil, err
	}
	for _, p := range peers {
		if p.Owned() {
			return p, nil
		}
	}
	return nil, nil
}

//ListSeedPeers return all the seed on the repository
func (r redisRepository) ListSeedPeers() ([]discovery.Seed, error) {
	id := fmt.Sprintf("%s:*", seedKey)
	list, err := r.fetchList(id)
	if err != nil {
		return nil, err
	}

	seeds := make([]discovery.Seed, 0)
	for _, s := range list {
		seed := FormatHashToSeed(s)
		seeds = append(seeds, seed)
	}

	return seeds, nil
}

//ListKnownPeers returns all the discoveredPeers on the repository
func (r redisRepository) ListKnownPeers() ([]discovery.Peer, error) {
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
func (r redisRepository) SetKnownPeer(p discovery.Peer) error {
	id := fmt.Sprintf("%s:%s", peerKey, p.Identity().PublicKey())
	cmd := r.client.HMSet(id, FormatPeerToHash(p))
	if cmd.Err() != nil {
		return cmd.Err()
	}
	return nil
}

//SetSeedPeer add a seed to the repository
func (r redisRepository) SetSeedPeer(s discovery.Seed) error {
	id := fmt.Sprintf("%s:%s", seedKey, s.PublicKey)
	cmd := r.client.HMSet(id, FormatSeedToHash(s))
	if cmd.Err() != nil {
		return cmd.Err()
	}
	return nil
}

//CountKnownPeers return the number of Known peers
func (r redisRepository) CountKnownPeers() (int, error) {
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

//GetPeerByIP get a peer from the repository using its ip
func (r redisRepository) GetKnownPeerByIP(ip net.IP) (discovery.Peer, error) {
	peers, err := r.ListKnownPeers()
	if err != nil {
		return nil, err
	}
	for _, p := range peers {
		if p.Identity().IP().Equal(ip) {
			return p, nil
		}
	}
	return nil, nil
}

//SetUnreachablePeer add an unreachable peer to the repository
func (r redisRepository) SetUnreachablePeer(pbKey string) error {
	cmd := r.client.SAdd(unreachablesKey, pbKey)
	if cmd.Err() != nil {
		return cmd.Err()
	}
	return nil
}

//RemoveUnreachablePeer remove an unreachable peer to the repository
func (r redisRepository) RemoveUnreachablePeer(pbKey string) error {
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
func (r redisRepository) ListReachablePeers() ([]discovery.Peer, error) {
	peers, err := r.ListKnownPeers()
	if err != nil {
		return nil, err
	}

	unreachableKeys, err := r.listUnreachableKeys()
	if err != nil {
		return nil, err
	}

	//When there is no unreachables returns the list of known peers
	if len(unreachableKeys) == 0 {
		return peers, nil
	}

	//We want to get the peers no include inside the list of unreachables
	pp := make([]discovery.Peer, 0)

	sort.Strings(unreachableKeys)

	for i, p := range peers {
		idx := sort.SearchStrings(unreachableKeys, p.Identity().PublicKey())
		if idx >= len(unreachableKeys) || peers[i].Identity().PublicKey() != unreachableKeys[idx] {
			pp = append(pp, peers[i])
		}
	}

	return pp, nil
}

//ListUnreacheablePeers returns all unreachable peers on the repository
func (r redisRepository) ListUnreachablePeers() ([]discovery.Peer, error) {
	unreachableKeys, err := r.listUnreachableKeys()
	if err != nil {
		return nil, err
	}
	peers, err := r.ListKnownPeers()
	if err != nil {
		return nil, err
	}

	pp := make([]discovery.Peer, 0)

	//Avoid looping if there is not unreachable keys
	if len(unreachableKeys) == 0 {
		return pp, nil
	}

	sort.Strings(unreachableKeys)

	for i, p := range peers {
		idx := sort.SearchStrings(unreachableKeys, p.Identity().PublicKey())
		if idx < len(unreachableKeys) && peers[i].Identity().PublicKey() == unreachableKeys[idx] {
			pp = append(pp, peers[i])
		}
	}

	return pp, nil
}

//ContainsUnreachableKey check if the pubk is in the list of unreacheable keys
func (r redisRepository) ContainsUnreachableKey(pubk string) error {
	unreachableKeys, err := r.listUnreachableKeys()
	if err != nil {
		return err
	}
	sort.Strings(unreachableKeys)
	idx := sort.SearchStrings(unreachableKeys, pubk)
	if idx < len(unreachableKeys) {
		return nil
	}
	return gossip.ErrNotFoundOnUnreachableList
}

func (r redisRepository) fetchList(key string) ([]map[string]string, error) {
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

func (r redisRepository) listUnreachableKeys() ([]string, error) {
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

//NewRepository creates a new repository using redis as storage
//
//An error is returned if the redis instance is not reached
func NewRepository(conf system.RedisConfig) (discovery.Repository, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		Password: conf.Pwd,
		DB:       0,
	})

	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}

	return redisRepository{
		client: client,
	}, nil
}

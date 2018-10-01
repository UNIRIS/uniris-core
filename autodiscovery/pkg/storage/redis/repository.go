package redis

import (
	"encoding/hex"
	"fmt"
	"net"

	"github.com/uniris/uniris-core/autodiscovery/pkg/system"

	"github.com/go-redis/redis"
	discovery "github.com/uniris/uniris-core/autodiscovery/pkg"
)

type redisRepository struct {
	client *redis.Client
}

const (
	ownedKey = "owned-peer"
	peerKey  = "discovered-peer"
	seedKey  = "seed-peer"
)

func (r *redisRepository) GetOwnedPeer() (p discovery.Peer, err error) {
	cmd := r.client.HGetAll(ownedKey)
	if cmd.Err() != nil {
		err = cmd.Err()
		return
	}

	res, err := cmd.Result()
	if err != nil {
		return
	}

	p = FormatHashToPeer(res)
	return
}

func (r *redisRepository) ListSeedPeers() ([]discovery.Seed, error) {
	cmdKeys := r.client.Keys(fmt.Sprintf("%s:*", seedKey))
	if cmdKeys.Err() != nil {
		return nil, cmdKeys.Err()
	}
	keys, err := cmdKeys.Result()
	if err != nil {
		return nil, err
	}

	ss := make([]discovery.Seed, 0)
	for _, k := range keys {
		cmdGet := r.client.HGetAll(k)
		if cmdGet.Err() != nil {
			return nil, cmdGet.Err()
		}
		res, err := cmdGet.Result()
		if err != nil {
			return nil, err
		}

		s := FormatHashToSeed(res)
		ss = append(ss, s)
	}

	return ss, nil
}

func (r *redisRepository) ListDiscoveredPeers() ([]discovery.Peer, error) {
	cmdKeys := r.client.Keys(fmt.Sprintf("%s:*", peerKey))
	if cmdKeys.Err() != nil {
		return nil, cmdKeys.Err()
	}
	keys, err := cmdKeys.Result()
	if err != nil {
		return nil, err
	}

	pp := make([]discovery.Peer, 0)
	for _, k := range keys {
		cmdGet := r.client.HGetAll(k)
		if cmdGet.Err() != nil {
			return nil, cmdGet.Err()
		}
		res, err := cmdGet.Result()
		if err != nil {
			return nil, err
		}

		p := FormatHashToPeer(res)
		pp = append(pp, p)
	}

	return pp, nil
}

func (r *redisRepository) SetPeer(p discovery.Peer) error {
	var id string
	if p.Owned() {
		id = ownedKey
	} else {
		id = fmt.Sprintf("%s:%s", peerKey, hex.EncodeToString(p.Identity().PublicKey()))
	}

	cmd := r.client.HMSet(id, FormatPeerToHash(p))
	if cmd.Err() != nil {
		return cmd.Err()
	}
	return nil
}

func (r *redisRepository) SetSeed(s discovery.Seed) error {
	id := fmt.Sprintf("%s:%s", seedKey, hex.EncodeToString(s.PublicKey))
	cmd := r.client.HMSet(id, FormatSeedToHash(s))
	if cmd.Err() != nil {
		return cmd.Err()
	}
	return nil
}

func (r *redisRepository) CountDiscoveredPeers() (int, error) {
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

func (r *redisRepository) GetPeerByIP(ip net.IP) (p discovery.Peer, err error) {
	owned, err := r.GetOwnedPeer()
	if err != nil {
		return
	}
	if owned.Identity().IP().Equal(ip) {
		return owned, nil
	}

	list, err := r.ListDiscoveredPeers()
	if err != nil {
		return
	}

	for _, p := range list {
		if p.Identity().IP().Equal(ip) {
			return p, nil
		}
	}

	return
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

	return &redisRepository{
		client: client,
	}, nil
}

package redis

import (
	"encoding/hex"
	"fmt"

	"github.com/go-redis/redis"
	discovery "github.com/uniris/uniris-core/autodiscovery/pkg"
)

type redisRepository struct {
	client *redis.Client

	peers []discovery.Peer
	seeds []discovery.Seed
}

func (r *redisRepository) GetOwnedPeer() (peer discovery.Peer, err error) {
	pp, err := r.ListKnownPeers()
	if err != nil {
		return
	}

	for _, p := range pp {
		if p.IsOwned() {
			peer = p
			break
		}
	}

	return
}

func (r *redisRepository) ListSeedPeers() ([]discovery.Seed, error) {
	cmdKeys := r.client.Keys("seed:*")
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

func (r *redisRepository) ListKnownPeers() ([]discovery.Peer, error) {
	cmdKeys := r.client.Keys("peer:*")
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
	id := fmt.Sprintf("peer:%s", hex.EncodeToString(p.PublicKey()))
	cmd := r.client.HMSet(id, FormatPeerToHash(p))
	if cmd.Err() != nil {
		return cmd.Err()
	}
	return nil
}

func (r *redisRepository) SetSeed(s discovery.Seed) error {
	id := fmt.Sprintf("seed:%s", s.IP.String())
	cmd := r.client.HMSet(id, FormatSeedToHash(s))
	if cmd.Err() != nil {
		return cmd.Err()
	}
	return nil
}

//NewRepository creates a new repository using redis as storage
//
//An error is returned if the redis instance is not reached
func NewRepository(host string, port int, pwd string) (discovery.Repository, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: pwd,
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

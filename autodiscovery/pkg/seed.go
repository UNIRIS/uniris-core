package discovery

import (
	"net"
)

//Seed is initial peer need to startup the discovery process
type Seed struct {
	IP        net.IP
	Port      int
	PublicKey PublicKey
}

//AsPeer converts a seed into a peer
func (s Seed) AsPeer() Peer {
	return &peer{
		identity: peerIdentity{
			ip:        s.IP,
			port:      s.Port,
			publicKey: s.PublicKey,
		},
	}
}

//SeedDiscoveryCounter define the interface to check the number of discovered node by a seed
type SeedDiscoveryCounter interface {
	CountDiscoveries() (int, error)
}

type seedDiscoveryCounter struct {
	repo Repository
}

//NewSeedDiscoveryCounter creates a new Counter for the seed discoveries
func NewSeedDiscoveryCounter(repo Repository) SeedDiscoveryCounter {
	return seedDiscoveryCounter{repo}
}

//CountDiscoveries report the average of node detected by the differents known seeds
func (sdc seedDiscoveryCounter) CountDiscoveries() (int, error) {
	listseed, err := sdc.repo.ListSeedPeers()
	if err != nil {
		return 0, err
	}
	avg := 0
	for i := 0; i < len(listseed); i++ {
		ipseed := listseed[i].IP
		p, err := sdc.repo.GetPeerByIP(ipseed)
		if p == nil {
			continue
		}
		if err == nil {
			avg += p.AppState().DiscoveredPeersNumber()
		}
	}
	avg = avg / len(listseed)
	return avg, nil
}

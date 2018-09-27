package discovery

import (
	"net"
)

//Seed is initial peer need to startup the discovery process
type Seed struct {
	IP   net.IP
	Port int
}

//ToPeer converts a seed into a peer
func (s Seed) ToPeer() Peer {
	return Peer{
		ip:   s.IP,
		port: s.Port,
	}
}

//SeedDiscoveryCounter define the interface to check the number of discovered node by a seed
type SeedDiscoveryCounter interface {
	Average() (int, error)
}

type seedDiscoveryCounter struct {
	repo Repository
}

//NewSeedDiscoveryCounter creates a new Counter for the seed discoveries
func NewSeedDiscoveryCounter(repo Repository) SeedDiscoveryCounter {
	return seedDiscoveryCounter{repo}
}

//Average report the average of node detected by the differents known seeds
func (sdc seedDiscoveryCounter) Average() (int, error) {
	listseed, err := sdc.repo.ListSeedPeers()
	if err != nil {
		return 0, err
	}
	avg := 0
	for i := 0; i < len(listseed); i++ {
		ipseed := listseed[i].IP
		p, err := sdc.repo.GetPeerByIP(ipseed)
		if err == nil {
			avg += p.DiscoveredPeersNumber()
		}
	}
	avg = avg / len(listseed)
	return avg, nil
}

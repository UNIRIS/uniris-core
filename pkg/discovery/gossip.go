package discovery

import (
	"errors"
	"log"
)

//Notifier handle the notification of the gossip events
type Notifier interface {
	NotifyReachable(PeerIdentity) error
	NotifyUnreachable(PeerIdentity) error
	NotifyDiscovery(Peer) error
}

//BootStrapingMinTime is the necessary minimum time on seconds to finish learning about the network
const BootStrapingMinTime = 1800

//Gossip initialize a cycle to spread the local view by updating its local peer and store the results
func Gossip(self Peer, seeds []PeerIdentity, db Database, netCheck NetworkChecker, sysR SystemReader, msg RoundMessenger, n Notifier) (Cycle, error) {
	if len(seeds) == 0 {
		return Cycle{}, errors.New("Cannot start a gossip round without a list seeds")
	}

	peers, err := db.DiscoveredPeers()
	if err != nil {
		return Cycle{}, err
	}

	reaches, err := reachablePeers(db)
	if err != nil {
		return Cycle{}, err
	}

	if err != nil {
		return Cycle{}, err
	}

	//Refresh ourself and append to the list of discovered peers
	self, err = updateSelf(self, reaches, seeds, db, netCheck, sysR)
	if err != nil {
		return Cycle{}, err
	}
	peers = append(peers, self)

	unreaches, err := db.UnreachablePeers()
	if err != nil {
		return Cycle{}, err
	}

	//Start the gossip Cycle
	c := Cycle{}
	if err := c.run(self, msg, seeds, peers, reaches.Identities(), unreaches); err != nil {
		return Cycle{}, err
	}

	//GossipStores and notifies the gossip Cycle result
	if err := addDiscoveries(c, db, n); err != nil {
		return c, err
	}
	if err := addReaches(c, db, n); err != nil {
		return c, err
	}
	if err := addUnreaches(c, db, n); err != nil {
		return c, err
	}

	return c, nil
}

func updateSelf(self Peer, reachables []Peer, seeds []PeerIdentity, db DatabaseWriter, netCheck NetworkChecker, sysR SystemReader) (Peer, error) {
	status, err := localStatus(self, seedReachableAverage(seeds, reachables), netCheck)
	if err != nil {
		return self, err
	}

	_, _, _, cpu, space, err := systemInfo(sysR)
	if err != nil {
		if err == ErrGeoPosition {
			status = FaultyPeer
			log.Println(ErrGeoPosition)
		} else {
			return self, err
		}
	}

	self.SelfRefresh(status, space, cpu, p2pFactor(reachables), len(reachables))
	return self, nil
}

func addDiscoveries(c Cycle, db DatabaseWriter, n Notifier) error {
	for _, p := range c.Discoveries {
		if err := db.WriteDiscoveredPeer(p); err != nil {
			return err
		}
		if err := n.NotifyDiscovery(p); err != nil {
			return err
		}
	}
	return nil
}

func addReaches(c Cycle, db Database, n Notifier) error {
	for _, p := range c.Reaches {
		exist, err := db.ContainsUnreachablePeer(p)
		if err != nil {
			return err
		}
		if exist {
			if err := db.RemoveUnreachablePeer(p); err != nil {
				return err
			}
			if err := n.NotifyReachable(p); err != nil {
				return err
			}
		}
	}
	return nil
}

func addUnreaches(c Cycle, db Database, n Notifier) error {
	for _, p := range c.Unreaches {
		exist, err := db.ContainsUnreachablePeer(p)
		if err != nil {
			return err
		}
		if !exist {
			if err := db.WriteUnreachablePeer(p); err != nil {
				return err
			}
			if err := n.NotifyUnreachable(p); err != nil {
				return err
			}
		}
	}
	return nil
}

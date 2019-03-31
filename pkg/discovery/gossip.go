package discovery

import (
	"errors"
	"github.com/uniris/uniris-core/pkg/logging"
	"math/rand"
	"reflect"
	"sync"
)

//ErrUnreachablePeer is returns when no owned peers has been stored
var ErrUnreachablePeer = errors.New("unreachable peer")

//Database wrap discovery database queries, persistence and removals
type Database interface {
	dbReader
	dbWriter
	dbRemover
}

type dbReader interface {
	//DiscoveredPeers retrieves the discovered peers from the discovery database
	DiscoveredPeers() ([]Peer, error)

	//UnreachablePeers retrieves the unreachables peers from the unreachables database
	UnreachablePeers() ([]PeerIdentity, error)
}

type dbWriter interface {

	//WriteDiscoveredPeer inserts or updates the peer in the discovery database
	WriteDiscoveredPeer(p Peer) error

	//WriteUnreachablePeer inserts the peer in the unreachable database
	WriteUnreachablePeer(p PeerIdentity) error
}

type dbRemover interface {
	//RemoveUnreachablePeer deletes the peer from the unreachable database
	RemoveUnreachablePeer(p PeerIdentity) error
}

//Notifier handle the notification of the gossip events
type Notifier interface {

	//NotifyReachable notifies the peer's public key which became reachable
	NotifyReachable(publicKey string) error

	//NotifyUnreachable notifies the peer's public key which became unreachable
	NotifyUnreachable(publicKey string) error

	//NotifyDiscovery notifies the peer's which has been discovered
	NotifyDiscovery(Peer) error
}

//Messenger represents a gossip client
type Messenger interface {
	SendSyn(target PeerIdentity, known []Peer, l logging.Logger) (requested []PeerIdentity, discovered []Peer, err error)
	SendAck(target PeerIdentity, requested []Peer, l logging.Logger) error
}

//Gossip initialize a cycle to spread the local view by updating its local peer and store the results
func Gossip(self Peer, seeds []PeerIdentity, db Database, netCheck NetworkChecker, sysR SystemReader, msg Messenger, n Notifier, l logging.Logger) error {
	if len(seeds) == 0 {
		return errors.New("cannot start a gossip round without a list seeds")
	}

	peers, err := db.DiscoveredPeers()
	if err != nil {
		return err
	}

	unreaches, err := db.UnreachablePeers()
	if err != nil {
		return err
	}
	reaches := reachablePeers(unreaches, peers)

	if err := updateSelf(self, reaches, seeds, db, netCheck, sysR, l); err != nil {
		return err
	}

	cDiscoveries, cReachables, cUnreachables, err := startCycle(self, msg, seeds, peers, reaches.identities(), unreaches, l)
	if err != nil {
		return err
	}

	if err := addDiscoveries(cDiscoveries, peers, db, n); err != nil {
		return err
	}
	if err := addReaches(cReachables, unreaches, db, n); err != nil {
		return err
	}
	if err := addUnreaches(cUnreachables, unreaches, db, n); err != nil {
		return err
	}

	return nil
}

//startCycle initiate a gossip cycle by creating rounds from a peer selection to spread the known peers and discover new peers
func startCycle(self Peer, msg Messenger, seeds []PeerIdentity, peers []Peer, reaches []PeerIdentity, unreaches []PeerIdentity, l logging.Logger) (discoveries []Peer, reachables []PeerIdentity, unreachables []PeerIdentity, err error) {

	//Pick the peers to gossip with
	selectedPeers := make([]PeerIdentity, 0)

	//We always pick a seed peer (as boostraping peer)
	selectedPeers = append(selectedPeers, randomPeer(seeds))

	//Because the self peer is included inside the database, we need to filter it out to not gossip with ourself
	reachFiltered := make([]PeerIdentity, 0)
	for _, r := range reaches {
		if r.publicKey != self.identity.publicKey {
			reachFiltered = append(reachFiltered, r)
		}
	}
	if len(reachFiltered) > 0 {
		selectedPeers = append(selectedPeers, randomPeer(reachFiltered))
	}

	//We include also a random unreachable peer to try to set it as reachable
	if len(unreaches) > 0 {
		selectedPeers = append(selectedPeers, randomPeer(unreaches))
	}

	//Need to wait the gossip with the selected peers to complete the cycle
	var wg sync.WaitGroup
	wg.Add(len(selectedPeers))

	//Start gossip for every selected peers
	for _, p := range selectedPeers {
		go func(target PeerIdentity) {
			defer wg.Done()

			//We initiate a gossip round and stores as discovery and reachable when the peers answers.
			//Otherwise it the peer cannot be reached, it will stored as unreachable for a later retry
			pp, err := startRound(target, peers, msg, l)
			if err != nil {
				if err == ErrUnreachablePeer {
					unreachables = append(unreachables, target)
					return
				}
				l.Error("unexpected error during round execution: " + err.Error())
				return
			}
			discoveries = append(discoveries, pp...)
			reachables = append(reachables, target)
		}(p)
	}

	wg.Wait()

	return
}

func randomPeer(items []PeerIdentity) PeerIdentity {
	if len(items) > 1 {
		rnd := rand.Intn(len(items))
		return items[rnd]
	}
	return items[0]
}

//startRound initiate the gossip round by messenging with the target peer
func startRound(target PeerIdentity, peers []Peer, msg Messenger, l logging.Logger) ([]Peer, error) {
	reqPeers, discoveries, err := msg.SendSyn(target, peers, l)
	if err != nil {
		return nil, err
	}
	//if some peers are requested, we send back the details of these peers
	if len(reqPeers) > 0 {
		reqDetailed := make([]Peer, 0)

		//Find details of the requested peers from the known peers
		for i := 0; i < len(reqPeers); i++ {
			var found bool
			var j int
			for !found && j < len(peers) {
				if peers[j].identity.publicKey == reqPeers[i].publicKey {
					reqDetailed = append(reqDetailed, peers[j])
					found = true
				}
				j++
			}
		}

		//Send to the SYN receiver an ACK with the peer detailed requested
		if err := msg.SendAck(target, reqDetailed, l); err != nil {
			return nil, err
		}
	}
	return discoveries, nil
}

//addDiscoveries persists the cycle discoveries in the database and send notifcation
func addDiscoveries(cDiscoveries []Peer, peers []Peer, db dbWriter, n Notifier) error {

	var oldFound bool
	var comparee Peer

	for _, dp := range cDiscoveries {

		oldFound = false

		if err := db.WriteDiscoveredPeer(dp); err != nil {
			return err
		}

		for _, p := range peers {
			if p.identity.publicKey == dp.identity.publicKey {
				comparee = p
				oldFound = true
				if !comparePeerIDAndState(dp, comparee) {
					if err := n.NotifyDiscovery(dp); err != nil {
						return err
					}
				}
				break
			}
		}

		if !oldFound {
			if err := n.NotifyDiscovery(dp); err != nil {
				return err
			}
		}
	}
	return nil
}

//addReaches removes the reachables peers from the unreachables in the database and send notification
func addReaches(cReaches []PeerIdentity, unreaches []PeerIdentity, db dbRemover, n Notifier) error {
	for _, p := range cReaches {
		if isUnreachable(p, unreaches) {
			if err := db.RemoveUnreachablePeer(p); err != nil {
				return err
			}
			if err := n.NotifyReachable(p.publicKey); err != nil {
				return err
			}
		}
	}
	return nil
}

//addUnreaches persist the unreachables if is not present in the database and  send notification
func addUnreaches(cUnreachables []PeerIdentity, unreaches []PeerIdentity, db dbWriter, n Notifier) error {
	for _, p := range cUnreachables {
		if !isUnreachable(p, unreaches) {
			if err := db.WriteUnreachablePeer(p); err != nil {
				return err
			}
			if err := n.NotifyUnreachable(p.publicKey); err != nil {
				return err
			}
		}
	}
	return nil
}

//ComparePeers compares a source of peers with an other list of peers
//and returns the peers that are not included inside the source or
func ComparePeers(source []Peer, comparees []Peer) []Peer {

	diff := make([]Peer, 0)

	for i := 0; i < len(comparees); i++ {
		var found bool
		var j int

		for !found && j < len(source) {
			//Add it to the list if the compared peer is include inside the source and if it's more recent
			if source[j].identity.publicKey == comparees[i].identity.publicKey {
				if comparees[i].hbState.MoreRecentThan(source[j].hbState) {
					diff = append(diff, comparees[i])
					found = true
				}
				found = true
			}
			j++
		}

		if !found {
			diff = append(diff, comparees[i])
		}
	}

	return diff
}

//isUnreachable determinates if a peer is contained inside the unreachable peer list
func isUnreachable(p PeerIdentity, unreaches []PeerIdentity) (found bool) {
	var i int
	for !found && i < len(unreaches) {
		if unreaches[i].publicKey == p.publicKey {
			found = true
		}
		i++
	}
	return
}

//reachablePeers filters known peers based on the unreachables list
//if there is not unreachables, the reachables are the known peers
func reachablePeers(unreachables []PeerIdentity, knownPeers []Peer) peerList {

	if len(unreachables) == 0 {
		return knownPeers
	}

	reachables := make([]Peer, 0)

	for _, p := range knownPeers {
		var found bool

		//detect if the peers is unreachable
		for _, u := range unreachables {
			if u.publicKey == p.identity.publicKey {
				found = true
				break
			}
		}

		//if not we add it to the list of reachables
		if !found {
			reachables = append(reachables, p)
		}
	}

	return reachables
}

//comparePeerIDAndState compares a peer with an other peer
//and returns false if at least identity and app state is different between the source and the comparee
func comparePeerIDAndState(source Peer, comparee Peer) bool {
	return source.identity.ip.Equal(comparee.identity.ip) &&
		source.identity.port == source.identity.port &&
		reflect.DeepEqual(source.appState, comparee.appState)
}

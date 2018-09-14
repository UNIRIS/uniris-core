package domain

import (
	"encoding/hex"
	"math/rand"
)

//GossipRound wraps the gossip round process
type GossipRound struct {
	Initiator  *Peer
	Seeds      []Peer
	KnownPeers []Peer
}

//NewGossipRound creates a gossip round
func NewGossipRound(seeds []Peer, knownPeers []Peer) GossipRound {
	return GossipRound{
		Seeds:      seeds,
		KnownPeers: knownPeers,
	}
}

//ExcludeOwnPeer filter to get only the peers that we don't own
func (r GossipRound) excludeOwnPeer() []Peer {
	noExclude := make([]Peer, 0)
	for _, peer := range r.KnownPeers {
		if !peer.IsOwned {
			noExclude = append(noExclude, peer)
		}
	}
	return noExclude
}

//SelectPeers perform the gossip peer selection
func (r GossipRound) SelectPeers() []Peer {
	receivers := make([]Peer, 0)
	randomSeed := r.getRandomPeer(r.Seeds)
	receivers = append(receivers, randomSeed)
	if len(r.excludeOwnPeer()) > 0 {
		randomPeer := r.getRandomPeer(r.KnownPeers)
		receivers = append(receivers, randomPeer)
	}

	return receivers
}

//GetRequestedPeers returns the known peers from a list of requested peers
func (r GossipRound) GetRequestedPeers(reqPeers []Peer) []Peer {
	knownPeersMap := make(map[string]Peer, 0)
	for _, peer := range r.KnownPeers {
		knownPeersMap[hex.EncodeToString(peer.PublicKey)] = peer
	}

	peers := make([]Peer, 0)
	for _, reqPeer := range reqPeers {
		if known, exist := knownPeersMap[hex.EncodeToString(reqPeer.PublicKey)]; exist && known.IsDiscovered() {
			peers = append(peers, known)
		}
	}
	return peers
}

//GetRandomPeer returns a random peer from a list of peers
func (r GossipRound) getRandomPeer(peers []Peer) Peer {
	if len(peers) > 1 {
		rnd := rand.Intn(len(peers) - 1)
		return peers[rnd]
	}
	return peers[0]
}

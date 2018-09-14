package usecases

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uniris/uniris-core/autodiscovery/domain"
)

/*
Scenario: Gossip with an empty targets
	Given an empty list of peer to reach
	When the peer gossip
	Then an error is returned
*/
func TestGossipWithoutPeers(t *testing.T) {
	repo := new(GossipTestPeerRepository)
	messenger := new(GosssipTestMessenger)
	geo := new(GosssipTestGeo)

	conf := GossipConfiguration{Messenger: messenger, Geolocalizer: geo}

	peersToReach := make([]domain.Peer, 0)
	err := Gossip(repo, conf, peersToReach)
	assert.Error(t, err, "Cannot gossip without peers to reach")
}

/*
Scenario: Execute the gossip with gossip targets
	Given a list of targets
	When the gossip is executed
	Then an error is returned
*/
func TestExecuteGossip(t *testing.T) {
	repo := new(GossipTestPeerRepository)
	messenger := new(GosssipTestMessenger)
	geo := new(GosssipTestGeo)

	conf := GossipConfiguration{Messenger: messenger, Geolocalizer: geo}

	gossipTargets := make([]domain.Peer, 0)
	err := Gossip(repo, conf, gossipTargets)
	assert.Error(t, err, "Cannot gossip without peers to reach")
}

//=========================
//INTERFACE IMPLEMENTATIONS
//=========================

type GosssipTestGeo struct{}

func (geo GosssipTestGeo) Lookup() (domain.GeoPosition, error) {
	return domain.GeoPosition{Lat: 10, Lon: 50, IP: net.ParseIP("127.0.0.1")}, nil
}

type GosssipTestMessenger struct{}

func (ms GosssipTestMessenger) SendSynchro(req domain.SynRequest) (*domain.SynAck, error) {
	newPeers := make([]domain.Peer, 0)
	newPeers = append(newPeers, domain.NewPeer([]byte("public key"), net.ParseIP("127.0.0.1"), 3545, false))
	return &domain.SynAck{
		NewPeers: newPeers,
	}, nil
}

func (ms GosssipTestMessenger) SendAcknowledgement(req domain.AckRequest) error {
	return nil
}

type GossipTestPeerRepository struct {
	peers []domain.Peer
}

func (r GossipTestPeerRepository) ListPeers() ([]domain.Peer, error) {
	return r.peers, nil
}

func (r *GossipTestPeerRepository) InsertPeer(p domain.Peer) error {
	r.peers = append(r.peers, p)
	return nil
}

func (r *GossipTestPeerRepository) UpdatePeer(p domain.Peer) error {
	newPeers := make([]domain.Peer, 0)
	for _, peer := range r.peers {
		if peer.Equals(p) {
			newPeers = append(newPeers, p)
		} else {
			newPeers = append(newPeers, peer)
		}
	}
	r.peers = newPeers
	return nil
}

func (r GossipTestPeerRepository) GetOwnedPeer() (owned domain.Peer, err error) {
	for _, peer := range r.peers {
		if peer.IsOwned {
			owned = peer
		}
	}
	return owned, nil
}

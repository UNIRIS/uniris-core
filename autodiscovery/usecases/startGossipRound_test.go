package usecases

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/uniris/uniris-core/autodiscovery/domain"
)

/*
Scenario: Start a gossip round without loaded seed
	Given a peer is starting up without seeds
	When the gossip round start
	Then a error is returned
*/
func TestStartGossipRoundWithoutSeeds(t *testing.T) {
	repo := new(StartGossipTestRepository)
	seedLoader := new(StartGossipTestSeedReader)
	messenger := new(StartGossipTestMessenger)

	err := StartGossipRound(repo, seedLoader, messenger)
	assert.Error(t, err, "Cannot gossip without seed peers")
}

/*
Scenario: Starts a gossip round without known peers
	Given a peer is starting up
	When the gossip start
	Then itself is registered on the repository
*/
func TestExecuteGossipRounder(t *testing.T) {
	repo := new(StartGossipTestRepository)
	seedLoader := new(StartGossipTestSeedReader)
	messenger := new(StartGossipTestMessenger)

	seedLoader.InitSeed()

	err := StartGossipRound(repo, seedLoader, messenger)
	assert.Nil(t, err)

	peers, _ := repo.ListPeers()
	assert.NotEmpty(t, peers)
}

//=========================
//INTERFACE IMPLEMENTATIONS
//=========================

type StartGossipTestSeedReader struct {
	seeds []domain.Peer
}

func (s StartGossipTestSeedReader) GetSeeds() ([]domain.Peer, error) {
	return s.seeds, nil
}

func (s *StartGossipTestSeedReader) InitSeed() {
	s.seeds = append(s.seeds, domain.NewPeer([]byte("public key"), net.ParseIP("127.0.0.1"), 3545, false))
}

type StartGossipTestMessenger struct{}

func (m StartGossipTestMessenger) SendSynchro(req domain.SynRequest) (*domain.SynAck, error) {
	newPeers := make([]domain.Peer, 0)
	newPeers = append(newPeers, domain.NewPeer([]byte("public key"), net.ParseIP("127.0.0.1"), 3545, false))
	return &domain.SynAck{
		NewPeers: newPeers,
	}, nil
}

func (m StartGossipTestMessenger) SendAcknowledgement(req domain.AckRequest) error {
	return nil
}

type StartGossipTestRepository struct {
	peers []domain.Peer
}

func (r StartGossipTestRepository) ListPeers() ([]domain.Peer, error) {
	return r.peers, nil
}

func (r *StartGossipTestRepository) InsertPeer(p domain.Peer) error {
	r.peers = append(r.peers, p)
	return nil
}

func (r *StartGossipTestRepository) UpdatePeer(p domain.Peer) error {
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

func (r StartGossipTestRepository) GetOwnedPeer() (owned domain.Peer, err error) {
	for _, peer := range r.peers {
		if peer.IsOwned {
			owned = peer
		}
	}
	return owned, nil
}

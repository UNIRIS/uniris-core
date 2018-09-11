package tests

import (
	"net"
	"testing"

	"github.com/uniris/uniris-core/autodiscovery/domain/entities"

	"github.com/stretchr/testify/assert"
	"github.com/uniris/uniris-core/autodiscovery/domain/usecases"
)

/*
Scenario: Get the unkwnown peers by comparing peer's list with the known peer's list
	Given a empty peer's list
	When we compare with a known peer's list
	Then we get the known peer's list
*/
func TestGetUnknownPeerFromEmpty(t *testing.T) {
	repo := GetRepo()
	usecases.LoadSeedPeers(&SeedLoader{}, repo)

	mySelf := &entities.Peer{
		PublicKey: GetValidPublicKey(),
	}

	knownPeers, _ := repo.ListPeers()
	unknownPeers := usecases.GetUnknownPeers(knownPeers, []*entities.Peer{}, mySelf)
	assert.NotEmpty(t, unknownPeers)
	assert.Equal(t, len(knownPeers), len(unknownPeers))
	assert.Equal(t, "127.0.0.1", unknownPeers[0].IP.String())
}

/*
Scenario: Get the unknown peers by comparing a peer's list with the known peer's list
	Given a peer's list
	When we compare with a known peer's list
	Then we get the only the unknown peers
*/
func TestGetUnknownPeers(t *testing.T) {
	repo := GetRepo()
	usecases.LoadSeedPeers(&SeedLoader{}, repo)

	mySelf := &entities.Peer{
		PublicKey: GetSecondValidPublicKey(),
	}

	peers := []*entities.Peer{
		&entities.Peer{
			IP:        net.ParseIP("20.10.200.10"),
			PublicKey: GetThirdValidPublicKey(),
		},
	}

	knownPeers, _ := repo.ListPeers()

	unknownPeers := usecases.GetUnknownPeers(knownPeers, peers, mySelf)
	assert.NotEmpty(t, unknownPeers)
	assert.Equal(t, GetValidPublicKey(), unknownPeers[0].PublicKey)
	assert.Equal(t, "127.0.0.1", unknownPeers[0].IP.String())
}

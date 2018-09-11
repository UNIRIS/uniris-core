package tests

import (
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
	knownPeers := []*entities.Peer{
		&entities.Peer{
			PublicKey: GetValidPublicKey(),
		},
	}
	unknownPeers := usecases.GetUnknownPeers(knownPeers, []*entities.Peer{})
	assert.NotEmpty(t, unknownPeers)
	assert.Equal(t, len(knownPeers), len(unknownPeers))
	assert.Equal(t, GetValidPublicKey(), unknownPeers[0].PublicKey, "Public key is not %s", string(GetValidPublicKey()))
}

/*
Scenario: Get the unknown peers by comparing a peer's list with the known peer's list
	Given a peer's list
	When we compare with a known peer's list
	Then we get the only the unknown peers
*/
func TestGetUnknownPeers(t *testing.T) {
	peers := []*entities.Peer{
		&entities.Peer{
			PublicKey: GetThirdValidPublicKey(),
		},
	}

	knownPeers := []*entities.Peer{
		&entities.Peer{
			PublicKey: GetValidPublicKey(),
		},
	}

	unknownPeers := usecases.GetUnknownPeers(knownPeers, peers)
	assert.NotEmpty(t, unknownPeers)
	assert.Equal(t, GetValidPublicKey(), unknownPeers[0].PublicKey)
	assert.Equal(t, GetValidPublicKey(), unknownPeers[0].PublicKey, "Public key is not %s", GetValidPublicKey())
}

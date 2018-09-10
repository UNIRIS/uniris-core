package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/uniris/uniris-core/autodiscovery/domain/usecases"
)

/*
Scenario: List known peers
	Given no peers knowns in the database
	When we ask the known peers
	Then an empty list is returned
*/
func TestListKnownPeers_Empty(t *testing.T) {
	repoPeer := new(PeerRepository)
	peers, err := usecases.ListKnownPeers(repoPeer)
	assert.Nil(t, err)
	assert.Empty(t, peers, "Peers should be empty")
}

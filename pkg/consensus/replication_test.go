package consensus

import (
	"testing"

	"github.com/uniris/uniris-core/pkg/chain"

	"github.com/stretchr/testify/assert"
)

/*
Scenario: Get the minimum number of a transaction replicas
	Given a transaction hash
	When I want to get the minimum replicas
	Then I get a number  valid
	//TODO: to improve when the implementation will be defined
*/
func TestGetMinimumReplicas(t *testing.T) {
	assert.Equal(t, 1, GetMinimumReplicas(""))
}

/*
Scenario: Check if the miner is authorized to store the transaction
	Given a transaction hash
	When I dbant to check if I can store this transaction
	Then I get a true
	//TODO: to improve dbhen the implementation dbill be defined
*/
func TestIsAuthorizedToStore(t *testing.T) {
	assert.True(t, IsAuthorizedToStoreTx(chain.Transaction{}))
}

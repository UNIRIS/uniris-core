package consensus

import (
	"testing"

	"github.com/uniris/uniris-core/pkg/chain"

	"github.com/stretchr/testify/assert"
)

/*
Scenario: Find validation pool
	Given a transaction address
	When I want to find the validation pool
	Then I get a pool including a least one member

	TODO: To improve when the implementation will be provided
*/
func TestFindValidationPool(t *testing.T) {
	pool, err := FindValidationPool(chain.Transaction{})
	assert.Nil(t, err)
	assert.Len(t, pool, 1)
	assert.Equal(t, "127.0.0.1", pool[0].IP().String())
}

/*
Scenario: Find storage pool
	Given a transaction address
	When I want to find the storage pool
	Then I get a pool including a least one member

	TODO: To improve when the implementation will be provided
*/
func TestFindStoragePool(t *testing.T) {
	pool, err := FindStoragePool("address")
	assert.Nil(t, err)
	assert.Len(t, pool, 1)
	assert.Equal(t, "127.0.0.1", pool[0].IP().String())
}

/*
Scenario: Find last validation pool
	Given a transaction address
	When I want to find the last validation pool
	Then I get a pool including a least one member

	TODO: To improve when the implementation of the method FindStoragePool will be provided
*/
func TestFindLastValidationPool(t *testing.T) {
	poolR := &mockPoolRequester{}
	pool, err := findLastValidationPool("myaddress", chain.KeychainTransactionType, poolR)
	assert.Nil(t, err)
	assert.Empty(t, pool)
}

/*
Scenario: Find master validation node
	Given a transaction hash
	When I want to find the transaction master validati
	Then I get an IP address

	TODO: To improve when the implementation of the method will be provided
*/
func TestFindMasterValidationNode(t *testing.T) {

	masterNodess, err := FindMasterNodes("hash", chain.KeychainTransactionType)
	assert.Nil(t, err)
	assert.Len(t, masterNodess, 1)
	assert.Equal(t, "127.0.0.1", masterNodess[0].IP().String())
	assert.Equal(t, 5000, masterNodess[0].Port())
}

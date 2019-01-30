package transaction

import (
	"net"
	"testing"

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
	s := PoolFindingService{}

	pool, err := s.FindValidationPool("myaddress")
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
	s := PoolFindingService{}

	pool, err := s.FindStoragePool("myaddress")
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
	s := PoolFindingService{
		pRetr: mockPoolRetriever{},
	}

	pool, err := s.FindLastValidationPool("myaddress", KeychainType)
	assert.Nil(t, err)
	assert.Len(t, pool, 1)
	assert.Equal(t, "127.0.0.1", pool[0].IP().String())
}

/*
Scenario: Find master validation peer
	Given a transaction hash
	When I want to find the transaction master validation
	Then I get an IP address

	TODO: To improve when the implementation of the method will be provided
*/
func TestFindMasterValidationPeer(t *testing.T) {
	s := PoolFindingService{}

	masterPeerIP, masterPeerPort := s.FindTransactionMasterPeer("hash")
	assert.Equal(t, "127.0.0.1", masterPeerIP)
	assert.Equal(t, 3545, masterPeerPort)
}

type mockPoolRetriever struct{}

func (pr mockPoolRetriever) RequestLastTransaction(pool Pool, txAddr string, txType Type) (*Transaction, error) {
	return &Transaction{
		masterV: MasterValidation{
			prevMiners: Pool{
				PoolMember{ip: net.ParseIP("127.0.0.1"), port: 3545},
			},
		},
	}, nil
}

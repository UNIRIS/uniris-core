package checks

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uniris/uniris-core/datamining/pkg/account"
	"github.com/uniris/uniris-core/datamining/pkg/mining"
)

/*
Scenario: Check data integrity
	Given a data and a transaction hash related to the data
	When we check the data integrity
	Then I get not error
*/
func TestIntegrityCheckData(t *testing.T) {
	data := &account.KeyChainData{}
	c := NewIntegrityChecker(mockHasher{})
	assert.Nil(t, c.CheckData(data, "hash"))
}

/*
Scenario: Check data integrity
	Given a data and a transaction hash related to the data
	When we check the data integrity
	Then I get not error
*/
func TestIntegrityCheckDataFails(t *testing.T) {
	c := NewIntegrityChecker(mockHasher{})
	data := &account.KeyChainData{}
	err := c.CheckData(data, "bad hash")
	assert.NotNil(t, err)
	assert.Equal(t, mining.ErrInvalidTransaction, err)
}

type mockHasher struct{}

func (h mockHasher) HashTransactionData(data interface{}) (string, error) {
	return "hash", nil
}

type mockBadHasher struct{}

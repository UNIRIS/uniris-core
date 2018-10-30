package checks

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
Scenario: Check data integrity
	Given a data and a transaction hash related to the data
	When we check the data integrity
	Then I get not error
*/
func TestIntegrityCheckData(t *testing.T) {
	c := NewIntegrityChecker(mockHasher{})
	assert.Nil(t, c.CheckData("data", "hash"))
}

/*
Scenario: Check data integrity
	Given a data and a transaction hash related to the data
	When we check the data integrity
	Then I get not error
*/
func TestIntegrityCheckDataFails(t *testing.T) {
	c := NewIntegrityChecker(mockHasher{})
	err := c.CheckData("data", "hash000")
	assert.NotNil(t, err)
	assert.Equal(t, errors.New("Invalid transaction"), err)
}

type mockHasher struct{}

func (h mockHasher) HashTransactionData(data interface{}) (string, error) {
	return "hash", nil
}

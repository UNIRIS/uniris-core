package file

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

/*
Scenario: Read keys fro files
	Given a keyreader
	When I read keys from files
	Then return values are not empty
*/
func TestKeyReader(t *testing.T) {

	r, err := NewReader()
	assert.Nil(t, err)
	assert.NotNil(t, r)

	k1, err := r.SharedBiodPublicKey()
	assert.NotNil(t, k1)
	assert.Nil(t, err)

	k2, err := r.SharedRobotPrivateKey()
	assert.Nil(t, err)
	assert.NotNil(t, k2)

	k3, err := r.SharedRobotPublicKey()
	assert.Nil(t, err)
	assert.NotNil(t, k3)

}

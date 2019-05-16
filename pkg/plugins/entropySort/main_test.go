package main

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
Scenario: Sort by entropy a list of authorized keys using the starting point characters
	Given a starting point (1d62567ec763002c9f88728a480629412cd33c673156a227bcd79b7adc8ac877) and a list of 3 keys (where hashes are: 1BD2B169A9E74A32133550E72E053AECD00500161BF87EB33D921A0DC63D1A71, BEF57EC7F53A6D40BEB640A780A639C83BC29AC8A9816F1FC6C5C6DCD93C4721,
	31C666C96118537BE81216E3A232DC7601779CD8D0D633980F0143FFC9B75FE6)
	When I want to sort the list by entropy
	Then I get the list sorted:
		- BEF57EC7F53A6D40BEB640A780A639C83BC29AC8A9816F1FC6C5C6DCD93C4721
		- 1BD2B169A9E74A32133550E72E053AECD00500161BF87EB33D921A0DC63D1A71
		- 31C666C96118537BE81216E3A232DC7601779CD8D0D633980F0143FFC9B75FE6
*/
func TestEntropySortWithStartingPointCharacter(t *testing.T) {

	pub1Hex := "0044657dab453d34f9adc2100a2cb8f38f644ef48e34b1d99d7c4d9371068e9438"
	pub2Hex := "00a8e0f20d4da185d0bf8bd0a45995dfc7926d545e5bbff0194fe34c42bf5e221b"
	pub3Hex := "00ee7a047a226e08ea14fe60ec4f6d328e56ebdb2ee2b9f5b1120e231e05c956a3"

	pub1, _ := hex.DecodeString(pub1Hex)
	pub2, _ := hex.DecodeString(pub2Hex)
	pub3, _ := hex.DecodeString(pub3Hex)

	pv, _ := hex.DecodeString("000c3bb61141f052e1936823a4a56224f2aae04084265655ff4c83d885295b570344657dab453d34f9adc2100a2cb8f38f644ef48e34b1d99d7c4d9371068e9438")

	sortedKeys, err := EntropySort([]byte("myhash"), [][]byte{pub1, pub2, pub3}, pv)
	assert.Nil(t, err)
	assert.Len(t, sortedKeys, 3)

	assert.Equal(t, pub2Hex, hex.EncodeToString(sortedKeys[0]))
	assert.Equal(t, pub1Hex, hex.EncodeToString(sortedKeys[1]))
	assert.Equal(t, pub3Hex, hex.EncodeToString(sortedKeys[2]))
}

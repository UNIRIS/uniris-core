package main

import (
	"bytes"
	"crypto/rand"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	ed25519 "golang.org/x/crypto/ed25519"
)

func TestMain(m *testing.M) {
	dir, _ := os.Getwd()
	os.Setenv("PLUGINS_DIR", filepath.Join(dir, "../"))
	m.Run()
}

/*
Scenario: Get the required number of coordinator
	Given 2 nodes in the network
	When I want to get the required number of coordinator
	Then I get 1
*/
func TestRequiredNumberOfCoordinatorsWith2Nodes(t *testing.T) {
	assert.Equal(t, RequiredNumberOfCoordinators(2, 1), 1)
}

/*
Scenario: Get the required number of coordinator
	Given 5 nodes in the network and 2 reachables
	When I want to get the required number of coordinator
	Then I get 1
*/
func TestRequiredNumberOfCoordinatorsWith5NodesAnd2Reachables(t *testing.T) {
	assert.Equal(t, RequiredNumberOfCoordinators(5, 2), 1)
}

/*
Scenario: Get the required number of coordinator
	Given 6 nodes in the network and 6 reachables
	When I want to get the required number of coordinator
	Then I get 5
*/
func TestRequiredNumberOfCoordinatorsWith6NodesAnd6Reachables(t *testing.T) {
	assert.Equal(t, RequiredNumberOfCoordinators(6, 6), 5)
}

/*
Scenario: Find coordinator pool
	Given a transaction hash, a list of 8 nodes in the networks with 5 reachables
	When I want to find node elected to the coordinator validation
	Then I get a list coordinator nodes with 5 reachables
*/
func TestFindcoordinatorPool(t *testing.T) {

	key1, _, _ := ed25519.GenerateKey(rand.Reader)
	key2, _, _ := ed25519.GenerateKey(rand.Reader)
	key3, _, _ := ed25519.GenerateKey(rand.Reader)
	key4, _, _ := ed25519.GenerateKey(rand.Reader)
	key5, _, _ := ed25519.GenerateKey(rand.Reader)
	key6, _, _ := ed25519.GenerateKey(rand.Reader)
	key7, _, _ := ed25519.GenerateKey(rand.Reader)
	key8, _, _ := ed25519.GenerateKey(rand.Reader)

	nodeReader := &mockNodeReader{
		nodes: []mockNode{
			mockNode{isReachable: false, key: mockPublicKey{bytes: key1[:], curve: 0}, patch: mockGeoPatch{id: 1}},
			mockNode{isReachable: true, key: mockPublicKey{bytes: key2[:], curve: 0}, patch: mockGeoPatch{id: 2}},
			mockNode{isReachable: true, key: mockPublicKey{bytes: key3[:], curve: 0}, patch: mockGeoPatch{id: 3}},
			mockNode{isReachable: false, key: mockPublicKey{bytes: key4[:], curve: 0}, patch: mockGeoPatch{id: 2}},
			mockNode{isReachable: true, key: mockPublicKey{bytes: key5[:], curve: 0}, patch: mockGeoPatch{id: 4}},
			mockNode{isReachable: true, key: mockPublicKey{bytes: key6[:], curve: 0}, patch: mockGeoPatch{id: 1}},
			mockNode{isReachable: true, key: mockPublicKey{bytes: key7[:], curve: 0}, patch: mockGeoPatch{id: 5}},
			mockNode{isReachable: true, key: mockPublicKey{bytes: key8[:], curve: 0}, patch: mockGeoPatch{id: 3}},
		},
	}

	nodePub, nodePv, _ := ed25519.GenerateKey(rand.Reader)
	nodePubKey := mockPublicKey{bytes: nodePub[:], curve: 0}
	nodePvKey := mockPrivateKey{bytes: nodePv[:]}

	_, crossPv, _ := ed25519.GenerateKey(rand.Reader)

	coordinatorNodes, err := FindCoordinatorPool([]byte("hash"), [][]byte{key1[:], key2[:], key3[:], key4[:], key5[:], key6[:], key7[:], key8[:]}, crossPv[:], nodePvKey, nodePubKey, nodeReader)
	assert.Nil(t, err)

	var nbReachables int
	for _, n := range coordinatorNodes.(ElectedNodeList).Nodes() {
		if !n.(electedNode).IsUnreachable() {
			nbReachables++
		}
	}
	assert.Equal(t, nbReachables, 5)
}

/*
Scenario: Find validation pool
    Given a transaction required 5 validations and 12 nodes with 3 unreachables located into 5 patches
    When I want to find the validation pool
    Then I get at least 7 nodes in the pool (5 + 5/2)
*/
func TestFindValidationPool(t *testing.T) {

	coordinatorPub, _, _ := ed25519.GenerateKey(rand.Reader)

	key1, _, _ := ed25519.GenerateKey(rand.Reader)
	key2, _, _ := ed25519.GenerateKey(rand.Reader)
	key3, _, _ := ed25519.GenerateKey(rand.Reader)
	key4, _, _ := ed25519.GenerateKey(rand.Reader)
	key5, _, _ := ed25519.GenerateKey(rand.Reader)
	key6, _, _ := ed25519.GenerateKey(rand.Reader)
	key7, _, _ := ed25519.GenerateKey(rand.Reader)
	key8, _, _ := ed25519.GenerateKey(rand.Reader)
	key9, _, _ := ed25519.GenerateKey(rand.Reader)
	key10, _, _ := ed25519.GenerateKey(rand.Reader)
	key11, _, _ := ed25519.GenerateKey(rand.Reader)
	key12, _, _ := ed25519.GenerateKey(rand.Reader)

	nodeReader := &mockNodeReader{
		nodes: []mockNode{
			mockNode{key: mockPublicKey{bytes: key1[:], curve: 0}, isReachable: false, patch: mockGeoPatch{id: 1}},
			mockNode{key: mockPublicKey{bytes: key2[:], curve: 0}, isReachable: true, patch: mockGeoPatch{id: 2}},
			mockNode{key: mockPublicKey{bytes: key3[:], curve: 0}, isReachable: true, patch: mockGeoPatch{id: 5}},
			mockNode{key: mockPublicKey{bytes: key4[:], curve: 0}, isReachable: true, patch: mockGeoPatch{id: 2}},
			mockNode{key: mockPublicKey{bytes: key5[:], curve: 0}, isReachable: false, patch: mockGeoPatch{id: 4}},
			mockNode{key: mockPublicKey{bytes: key6[:], curve: 0}, isReachable: true, patch: mockGeoPatch{id: 3}},
			mockNode{key: mockPublicKey{bytes: key7[:], curve: 0}, isReachable: true, patch: mockGeoPatch{id: 2}},
			mockNode{key: mockPublicKey{bytes: key8[:], curve: 0}, isReachable: false, patch: mockGeoPatch{id: 1}},
			mockNode{key: mockPublicKey{bytes: key9[:], curve: 0}, isReachable: true, patch: mockGeoPatch{id: 1}},
			mockNode{key: mockPublicKey{bytes: key10[:], curve: 0}, isReachable: true, patch: mockGeoPatch{id: 2}},
			mockNode{key: mockPublicKey{bytes: key11[:], curve: 0}, isReachable: true, patch: mockGeoPatch{id: 3}},
			mockNode{key: mockPublicKey{bytes: key12[:], curve: 0}, isReachable: true, patch: mockGeoPatch{id: 4}},
		},
	}

	_, crossPv, _ := ed25519.GenerateKey(rand.Reader)

	nodePub, nodePv, _ := ed25519.GenerateKey(rand.Reader)
	nodePubKey := mockPublicKey{bytes: nodePub[:], curve: 0}
	nodePvKey := mockPrivateKey{bytes: nodePv[:]}

	authKeys := [][]byte{key1[:], key2[:], key3[:], key4[:], key5[:], key6[:], key7[:], key8[:], key9[:], key10[:], key11[:], key12[:]}

	pool, err := FindValidationPool([]byte("address"), 5, coordinatorPub[:], authKeys, crossPv[:], nodePvKey, nodePubKey, nodeReader)
	assert.Nil(t, err)
	assert.True(t, len(pool.(ElectedNodeList).Nodes()) >= 7)

	distinctPatches := make([]int, 0)
	for _, h := range pool.(ElectedNodeList).Nodes() {
		var found bool
		for _, p := range distinctPatches {
			if p == h.(electedNode).PatchNumber() {
				found = true
				break
			}
		}
		if !found {
			distinctPatches = append(distinctPatches, h.(electedNode).PatchNumber())
		}
	}

	assert.True(t, len(distinctPatches) >= 5)
}

/*
Scenario: Get the minimum validation number for a system transaction with tiny network
	Given a system transaction with 1 nodes and 1 reachable
	When I want to get the validation required number
	Then I get 1
*/
func TestTestValidationNumberNetworkBasedWithTinyNetwork(t *testing.T) {
	nbValidations, err := requiredValidationNumberForSysTX(1, 1)
	assert.Nil(t, err)
	assert.Equal(t, 1, nbValidations)
}

/*
Scenario: Get the minimum validation number for a system transaction with small network
	Given a system transaction with less than 5 nodes and 2 reachables
	When I want to get the validation required number
	Then I get 2
*/
func TestTestValidationNumberNetworkBasedWithSmallNetwork(t *testing.T) {

	nbValidations, err := requiredValidationNumberForSysTX(5, 2)
	assert.Nil(t, err)
	assert.Equal(t, 2, nbValidations)
}

/*
Scenario: Get the minimum validation number for a system transaction with normal network
	Given a system transaction with 10 nodes and 10 reachables
	When I want to get the validation required number
	Then I get 2
*/
func TestValidationNumberNetworkBasedWithNormalNetwork(t *testing.T) {
	nbValidations, err := requiredValidationNumberForSysTX(10, 10)
	assert.Nil(t, err)
	assert.Equal(t, 5, nbValidations)
}

/*
Scenario: Get the minimum validation number for a system transaction with too less nodes
	Given a system transaction with 5 nodes and 1 reachable
	When I want to get the validation required number
	Then I get an error
*/
func TestValidationNumberNetworkBasedWithUnsufficientNetwork(t *testing.T) {
	_, err := requiredValidationNumberForSysTX(6, 1)
	assert.EqualError(t, err, "no enough nodes in the network to validate this transaction")
}

/*
SCenario: Get the minimum validation number for transaction with fees
	Given a transaction fees as 1 UCO and 10 reachables nodes
	When I want to get the validation required number
	Then I get 9 validations neeed
*/
func TestValidationNumberFeesBasedFor1UCOFeesWith10Nodes(t *testing.T) {
	nb, err := requiredValidationNumberWithFees(1, 10)
	assert.Nil(t, err)
	assert.Equal(t, 9, nb)
}

/*
SCenario: Get the minimum validation number for transaction with fees
	Given a transaction fees as 1 UCO and 8 reachables nodes
	When I want to get the validation required number
	Then I get 9 validations neeed
*/
func TestValidationNumberFeesBasedFor1UCOFeesWith8Nodes(t *testing.T) {
	nb, err := requiredValidationNumberWithFees(1, 8)
	assert.Nil(t, err)
	assert.Equal(t, 8, nb)
}

/*
Scenario: Get the minimum validation for normal transaction with less 3 nodes
	Given a system transaction and less 3 nodes
	When I want to get the validation required number
	Then I get an error
*/
func TestRequiredValidationNumberNormalTxWithLess3Nodes(t *testing.T) {
	authKeys := make([][]byte, 2)
	_, err := RequiredValidationNumber(0, 0.001, 2, authKeys)
	assert.EqualError(t, err, "no enough nodes in the network to validate this transaction")
}

/*
Scenario: Get the minimum validation for system transaction
	Given a system transaction and 5 nodes and 5 reachables
	When I want to get the validation required number
	Then I get 5 validation
*/
func TestRequiredValidationNumberSystemTxWith5Nodes(t *testing.T) {
	authKeys := make([][]byte, 5)
	nbValidations, err := RequiredValidationNumber(4, 0, 5, authKeys)
	assert.Nil(t, err)
	assert.Equal(t, 5, nbValidations)
}

/*
Scenario: Get the minimum validation  for normal transaction
	Given a normal transaction (minium fees: 0.001 => 3 validations)
	When I want to get the validation required number
	Then I get 3 validation
*/
func TestRequiredValidationNumberWith5UCO(t *testing.T) {
	authKeys := make([][]byte, 5)
	nbValidations, err := RequiredValidationNumber(2, 0.001, 5, authKeys)
	assert.Nil(t, err)
	assert.Equal(t, 3, nbValidations)
}

type mockNodeReader struct {
	nodes []mockNode
}

func (r mockNodeReader) CountReachables() (int, error) {
	count := 0
	for _, n := range r.nodes {
		if n.IsReachable() {
			count++
		}
	}
	return count, nil
}

func (r mockNodeReader) Reachables() (nodes []node, err error) {
	for _, n := range r.nodes {
		if n.IsReachable() {
			nodes = append(nodes, n)
		}
	}
	return
}

func (r mockNodeReader) FindByPublicKey(key []byte) (node, error) {
	for _, n := range r.nodes {
		if bytes.Equal(key, n.key.bytes) {
			return n, nil
		}
	}
	return nil, errors.New("node not found")
}

type mockNode struct {
	key         mockPublicKey
	isReachable bool
	patch       geoPatch
	status      int
}

func (n mockNode) PublicKey() interface{} {
	return n.key
}

func (n mockNode) IsReachable() bool {
	return n.isReachable
}

func (n mockNode) Patch() geoPatch {
	return n.patch
}

func (n mockNode) Status() int {
	return n.status
}

type mockGeoPatch struct {
	id int
}

func (g mockGeoPatch) ID() int {
	return g.id
}

type mockPublicKey struct {
	bytes []byte
	curve int
}

func (pb mockPublicKey) Marshal() []byte {
	out := make([]byte, 1+len(pb.bytes))
	out[0] = byte(int(pb.curve))
	copy(out[1:], pb.bytes)
	return out
}

func (pb mockPublicKey) Verify(data []byte, sig []byte) (bool, error) {
	return ed25519.Verify(pb.bytes, data, sig), nil
}

type mockPrivateKey struct {
	bytes []byte
}

func (pv mockPrivateKey) Sign(data []byte) ([]byte, error) {
	return ed25519.Sign(pv.bytes, data), nil
}

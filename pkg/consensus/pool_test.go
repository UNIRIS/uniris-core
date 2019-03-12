package consensus

import (
	"log"
	"testing"

	"github.com/uniris/uniris-core/pkg/chain"
	"github.com/uniris/uniris-core/pkg/crypto"

	"github.com/stretchr/testify/assert"

	"github.com/uniris/uniris-core/pkg/shared"
)

/*
Scenario: Get the required number of master
	Given 2 nodes in the network
	When I want to get the required number of master
	Then I get 1
*/
func TestRequiredNumberOfMasterWith2Nodes(t *testing.T) {
	assert.Equal(t, requiredNumberOfMaster(2, 1), 1)
}

/*
Scenario: Get the required number of master
	Given 5 nodes in the network and 2 reachables
	When I want to get the required number of master
	Then I get 1
*/
func TestRequiredNumberOfMasterWith5NodesAnd2Reachables(t *testing.T) {
	assert.Equal(t, requiredNumberOfMaster(5, 2), 1)
}

/*
Scenario: Get the required number of master
	Given 6 nodes in the network and 6 reachables
	When I want to get the required number of master
	Then I get 5
*/
func TestRequiredNumberOfMasterWith6NodesAnd6Reachables(t *testing.T) {
	assert.Equal(t, requiredNumberOfMaster(6, 6), 5)
}

/*
Scenario: Create a starting point for entropy sorting
	Given a transaction hash and a node private key
	When I want to create the starting point
	Then I get an HMAC and I can reproduce the same output
*/
func TestBuildStartingPoint(t *testing.T) {
	hmac := buildStartingPoint("myhash", "mykey")
	assert.NotEmpty(t, hmac)
	assert.Equal(t, hmac, buildStartingPoint("myhash", "mykey"))
}

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

	authKeys := []string{
		"000000000000000000", //1BD2B169A9E74A32133550E72E053AECD00500161BF87EB33D921A0DC63D1A71
		"abcdef",             //BEF57EC7F53A6D40BEB640A780A639C83BC29AC8A9816F1FC6C5C6DCD93C4721
		"afdsfsdf",           //31C666C96118537BE81216E3A232DC7601779CD8D0D633980F0143FFC9B75FE6
		"abc",                //BA7816BF8F01CFEA414140DE5DAE2223B00361A396177A9CB410FF61F20015AD
	}

	sortedKeys := entropySort("myhash", authKeys, "mykey")
	log.Print(sortedKeys)
	assert.Len(t, sortedKeys, 4)
	assert.Equal(t, "abcdef", sortedKeys[0])
	assert.Equal(t, "000000000000000000", sortedKeys[1])
	assert.Equal(t, "abc", sortedKeys[2])
	assert.Equal(t, "afdsfsdf", sortedKeys[3])
}

/*
Scenario: Find master validation node
	Given a transaction hash, a list of 8 nodes in the networks with 5 reachables
	When I want to find node elected to the master validation
	Then I get a list master nodes with 5 reachables
*/
func TestFindMasterValidationNode(t *testing.T) {
	nodeDB := &mockNodeDatabase{
		nodes: []Node{
			Node{publicKey: "pub1", isReachable: false},
			Node{publicKey: "pub2", isReachable: true},
			Node{publicKey: "pub3", isReachable: true},
			Node{publicKey: "pub4", isReachable: true},
			Node{publicKey: "pub5", isReachable: false},
			Node{publicKey: "pub6", isReachable: true},
			Node{publicKey: "pub7", isReachable: true},
			Node{publicKey: "pub8", isReachable: true},
		},
	}

	masterNodes, err := FindMasterNodes("hash", nodeDB, &mockSharedNodeReader{})
	assert.Nil(t, err)

	var nbReachables int
	for _, n := range masterNodes {
		if n.isReachable {
			nbReachables++
		}
	}
	assert.Equal(t, nbReachables, 5)
}

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

type mockSharedNodeReader struct{}

func (r mockSharedNodeReader) NodeFirstKeys() (shared.KeyPair, error) {
	pub, pv := crypto.GenerateKeys()
	return shared.NewKeyPair(pub, pv)
}

func (r mockSharedNodeReader) NodeLastKeys() (shared.KeyPair, error) {
	return shared.KeyPair{}, nil
}

func (r mockSharedNodeReader) AuthorizedPublicKeys() ([]string, error) {
	return []string{
		"pub1",
		"pub2",
		"pub3",
		"pub4",
		"pub5",
		"pub6",
		"pub7",
		"pub8",
	}, nil
}

type mockNodeReader struct {
	nodes []Node
}

func (db mockNodeDatabase) Reachables() (reachables []Node, err error) {
	for _, n := range db.nodes {
		if n.isReachable {
			reachables = append(reachables, n)
		}
	}
	return
}

func (db mockNodeDatabase) Unreachables() (unreachables []Node, err error) {
	for _, n := range db.nodes {
		if !n.isReachable {
			unreachables = append(unreachables, n)
		}
	}
	return
}

func (db mockNodeDatabase) CountReachables() (nb int, err error) {
	for _, n := range db.nodes {
		if n.isReachable {
			nb++
		}
	}
	return
}

func (db *mockNodeDatabase) FindByPublicKey(publicKey string) (found Node, err error) {
	for _, n := range db.nodes {
		if n.publicKey == publicKey {
			return n, nil
		}
	}
	return
}

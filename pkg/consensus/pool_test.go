package consensus

import (
	"crypto/rand"
	"encoding/hex"
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
	pv, _, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
	sp, err := buildStartingPoint(crypto.VersionnedHash("myhash"), pv)
	assert.Nil(t, err)
	assert.NotEmpty(t, sp)

	sp2, err := buildStartingPoint(crypto.VersionnedHash("myhash"), pv)
	assert.Nil(t, err)

	assert.Equal(t, sp, sp2)
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

	pub1Hex := "0044657dab453d34f9adc2100a2cb8f38f644ef48e34b1d99d7c4d9371068e9438"
	pub2Hex := "00a8e0f20d4da185d0bf8bd0a45995dfc7926d545e5bbff0194fe34c42bf5e221b"
	pub3Hex := "00ee7a047a226e08ea14fe60ec4f6d328e56ebdb2ee2b9f5b1120e231e05c956a3"

	pb1B, _ := hex.DecodeString(pub1Hex)
	pb2B, _ := hex.DecodeString(pub2Hex)
	pb3B, _ := hex.DecodeString(pub3Hex)

	pub1, _ := crypto.ParsePublicKey(pb1B)
	pub2, _ := crypto.ParsePublicKey(pb2B)
	pub3, _ := crypto.ParsePublicKey(pb3B)

	authKeys := []crypto.PublicKey{pub1, pub2, pub3}

	pvBytes, _ := hex.DecodeString("000c3bb61141f052e1936823a4a56224f2aae04084265655ff4c83d885295b570344657dab453d34f9adc2100a2cb8f38f644ef48e34b1d99d7c4d9371068e9438")
	pv, _ := crypto.ParsePrivateKey(pvBytes)

	sortedKeys, err := entropySort([]byte("myhash"), authKeys, pv)
	assert.Nil(t, err)
	assert.Len(t, sortedKeys, 3)

	sorted1Bytes, _ := sortedKeys[0].Marshal()
	sorted2Bytes, _ := sortedKeys[1].Marshal()
	sorted3Bytes, _ := sortedKeys[2].Marshal()

	assert.Equal(t, pub3Hex, hex.EncodeToString(sorted1Bytes))
	assert.Equal(t, pub1Hex, hex.EncodeToString(sorted2Bytes))
	assert.Equal(t, pub2Hex, hex.EncodeToString(sorted3Bytes))
}

/*
Scenario: Find master validation node
	Given a transaction hash, a list of 8 nodes in the networks with 5 reachables
	When I want to find node elected to the master validation
	Then I get a list master nodes with 5 reachables
*/
func TestFindMasterValidationNode(t *testing.T) {

	_, pub1, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
	_, pub2, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
	_, pub3, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
	_, pub4, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
	_, pub5, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
	_, pub6, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
	_, pub7, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
	_, pub8, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)

	nodeDB := &mockNodeDatabase{
		nodes: []Node{
			Node{publicKey: pub1, isReachable: false},
			Node{publicKey: pub2, isReachable: true},
			Node{publicKey: pub3, isReachable: true},
			Node{publicKey: pub4, isReachable: true},
			Node{publicKey: pub5, isReachable: false},
			Node{publicKey: pub6, isReachable: true},
			Node{publicKey: pub7, isReachable: true},
			Node{publicKey: pub8, isReachable: true},
		},
	}

	pv, pub, _ := crypto.GenerateECKeyPair(crypto.Ed25519Curve, rand.Reader)
	crossNodeKeys, _ := shared.NewNodeCrossKeyPair(pub, pv)

	masterNodes, err := FindMasterNodes([]byte("hash"), nodeDB, &mockSharedKeyReader{
		authKeys:      []crypto.PublicKey{pub1, pub2, pub3, pub4, pub5, pub6, pub7, pub8},
		crossNodeKeys: []shared.NodeCrossKeyPair{crossNodeKeys},
	})
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
	pool, err := FindStoragePool([]byte("address"))
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
	pool, err := findLastValidationPool([]byte("address"), chain.KeychainTransactionType, poolR)
	assert.Nil(t, err)
	assert.Empty(t, pool)
}

type mockSharedKeyReader struct {
	crossNodeKeys    []shared.NodeCrossKeyPair
	crossEmitterKeys []shared.EmitterCrossKeyPair
	authKeys         []crypto.PublicKey
}

func (r mockSharedKeyReader) EmitterCrossKeypairs() ([]shared.EmitterCrossKeyPair, error) {
	return r.crossEmitterKeys, nil
}

func (r mockSharedKeyReader) FirstNodeCrossKeypair() (shared.NodeCrossKeyPair, error) {
	return r.crossNodeKeys[0], nil
}

func (r mockSharedKeyReader) LastNodeCrossKeypair() (shared.NodeCrossKeyPair, error) {
	return r.crossNodeKeys[len(r.crossNodeKeys)-1], nil
}

func (r mockSharedKeyReader) AuthorizedNodesPublicKeys() ([]crypto.PublicKey, error) {
	return r.authKeys, nil
}

func (r mockSharedKeyReader) CrossEmitterPublicKeys() (pubKeys []crypto.PublicKey, err error) {
	for _, kp := range r.crossEmitterKeys {
		pubKeys = append(pubKeys, kp.PublicKey())
	}
	return
}

func (r mockSharedKeyReader) FirstEmitterCrossKeypair() (shared.EmitterCrossKeyPair, error) {
	return r.crossEmitterKeys[0], nil
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

func (db *mockNodeDatabase) FindByPublicKey(publicKey crypto.PublicKey) (found Node, err error) {
	for _, n := range db.nodes {
		if n.publicKey.Equals(publicKey) {
			return n, nil
		}
	}
	return
}

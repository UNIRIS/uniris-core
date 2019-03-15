package consensus

import (
	"crypto/hmac"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"

	"github.com/uniris/uniris-core/pkg/shared"

	"github.com/uniris/uniris-core/pkg/chain"
	"github.com/uniris/uniris-core/pkg/crypto"
)

//Pool represent a pool either for sharding or validation
type Pool struct {
	nodes   []Node
	headers []chain.NodeHeader
}

//Nodes returns the nodes of the pool
func (p Pool) Nodes() []Node {
	return p.nodes
}

//PoolRequester handles the request to perform on a pool during the mining
type PoolRequester interface {
	//RequestTransactionTimeLock asks a pool to timelock a transaction using the address related
	RequestTransactionTimeLock(pool Pool, txHash crypto.VersionnedHash, txAddr crypto.VersionnedHash, masterPublicKey crypto.PublicKey) error

	//RequestTransactionValidations asks a pool to validation a transaction
	RequestTransactionValidations(pool Pool, tx chain.Transaction, minValid int, masterValid chain.MasterValidation) ([]chain.Validation, error)

	//RequestTransactionStorage asks a pool to store a transaction
	RequestTransactionStorage(pool Pool, minStorage int, tx chain.Transaction) error

	//RequestLastTransaction asks a pool to retrieve the last transaction of an address
	RequestLastTransaction(pool Pool, txAddr crypto.VersionnedHash, txType chain.TransactionType) (*chain.Transaction, error)
}

//FindMasterNodes finds a list of master nodes by using an entropy sKeysing based on the transaction and minimum number of master
func FindMasterNodes(txHash crypto.VersionnedHash, nodeReader NodeReader, sharedKeyReader shared.KeyReader) (mPool Pool, err error) {
	authKeys, err := sharedKeyReader.AuthorizedNodesPublicKeys()
	if err != nil {
		return
	}

	firstKeys, err := sharedKeyReader.FirstNodeCrossKeypair()
	if err != nil {
		return
	}

	sKeys, err := entropySort(txHash, authKeys, firstKeys.PrivateKey())
	if err != nil {
		return
	}
	nbReachables, err := nodeReader.CountReachables()
	if err != nil {
		return
	}
	nbMasters := requiredNumberOfMaster(len(authKeys), nbReachables)

	nbReachableMasters := 0

	for i := 0; nbReachableMasters < nbMasters && i < len(sKeys); i++ {
		n, err := nodeReader.FindByPublicKey(sKeys[i])
		if err != nil {
			return Pool{}, err
		}

		//check if the node exists, happens only when there is some networking issues
		//or if the node has not been discovered yet by the gossip service
		if n.publicKey == nil {
			continue
		}

		mPool.nodes = append(mPool.nodes, n)
		if n.isReachable {
			nbReachableMasters++
		}
	}
	if nbReachableMasters != nbMasters {
		return Pool{}, fmt.Errorf("cannot proceed transaction with an invalid number of reachables master nodes (%d)", nbReachableMasters)
	}
	return
}

//requiredNumberOfMaster returns the number of master based on the network capacity
func requiredNumberOfMaster(nbNodes int, nbReachables int) int {
	if nbNodes < 5 && nbReachables >= 1 {
		return 1
	} else if nbNodes >= 5 && nbReachables <= 5 {
		return 1
	}
	return 5
}

//buildStartingPoint creates a starting point by using an HMAC of the transaction hash and the first node shared private key
func buildStartingPoint(txHash crypto.VersionnedHash, nodeMinerPrivateKey crypto.PrivateKey) (string, error) {
	pvBytes, err := nodeMinerPrivateKey.Marshal()
	if err != nil {
		return "", err
	}
	h := hmac.New(crypto.DefaultHashAlgo.New, pvBytes)
	h.Write([]byte(txHash))
	return hex.EncodeToString(h.Sum(nil)), nil
}

//entropySort sKeyss a list of nodes public keys using a "starting point" (HMAC of the transaction hash with the first node shared private key) and the hashes of the node public keys
func entropySort(txHash crypto.VersionnedHash, authKeys []crypto.PublicKey, nodeFirstKey crypto.PrivateKey) (sortedKeys []crypto.PublicKey, err error) {
	startingPoint, err := buildStartingPoint(txHash, nodeFirstKey)
	if err != nil {
		return nil, err
	}

	//Building list of public keys and map of hash-›key
	hashKeys := make([]string, len(authKeys))
	mHashKeys := make(map[string]crypto.PublicKey, 0)
	for i, k := range authKeys {
		keyBytes, err := k.Marshal()
		if err != nil {
			return nil, err
		}

		h := crypto.Hash(keyBytes)
		mHashKeys[hex.EncodeToString(h)] = k
		hashKeys[i] = hex.EncodeToString(h)
	}

	hashKeys = append(hashKeys, startingPoint)
	sort.Strings(hashKeys)
	var startPointIndex int
	for i, k := range hashKeys {
		if startingPoint == k {
			startPointIndex = i
			break
		}
	}

	end := crypto.DefaultHashAlgo.Size()

	//Sort keys by comparing the last character of the key with a starting point character
	for p := 0; len(sortedKeys) < len(hashKeys)-1 && p < end; p++ {

		//iterating from the starting point to the end of the list
		//add add the key if the latest character matchew the start point position
		for i := startPointIndex + 1; i < len(hashKeys); i++ {
			if []rune(hashKeys[i])[end-1] == []rune(startingPoint)[p] {
				var contains bool
				for _, k := range sortedKeys {
					if k.Equals(mHashKeys[hashKeys[i]]) {
						contains = true
						break
					}
				}
				if !contains {
					sortedKeys = append(sortedKeys, mHashKeys[hashKeys[i]])
				}
			}
		}

		//iterating from the 0 to the starting point
		//and add the key if the latest character matches the start point position
		for i := 0; i < startPointIndex; i++ {
			if []rune(hashKeys[i])[end-1] == []rune(startingPoint)[p] {
				var contains bool
				for _, k := range sortedKeys {
					if k.Equals(mHashKeys[hashKeys[i]]) {
						contains = true
						break
					}
				}
				if !contains {
					sortedKeys = append(sortedKeys, mHashKeys[hashKeys[i]])
				}
			}
		}
	}

	//We have tested all the characters of the staring point and not yet finished the sKeysing operation, we will loop on all the hex characters to finish sKeysing
	if len(sortedKeys) < len(hashKeys)-1 {
		hexChar := []rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f'}
		for p := 0; len(sortedKeys) < len(hashKeys)-1 && p < len(hexChar); p++ {

			//iterating from the starting point to the end of the list
			//add add the key if the latest character matchew the start point position
			for i := startPointIndex + 1; i < len(hashKeys); i++ {
				if []rune(hashKeys[i])[end-1] == hexChar[p] {
					var contains bool
					for _, k := range sortedKeys {
						if k.Equals(mHashKeys[hashKeys[i]]) {
							contains = true
							break
						}
					}
					if !contains {
						sortedKeys = append(sortedKeys, mHashKeys[hashKeys[i]])
					}
				}
			}

			//iterating from the 0 to the starting point
			//and add the key if the latest character matches the start point position
			for i := 0; i < startPointIndex; i++ {
				if []rune(hashKeys[i])[end-1] == hexChar[p] {
					var contains bool
					for _, k := range sortedKeys {
						if k.Equals(mHashKeys[hashKeys[i]]) {
							contains = true
							break
						}
					}
					if !contains {
						sortedKeys = append(sortedKeys, mHashKeys[hashKeys[i]])
					}
				}
			}
		}
	}

	return
}

//FindStoragePool searches a storage pool for the given address
func FindStoragePool(address crypto.VersionnedHash, r NodeReader) (Pool, error) {
	//Because of the entropy of the master election and without the sharding implementation
	//We cannot be sure to retrieve the data
	//So in waiting the sharding implementation, we need to select one of the master peers elected to retrieve data
	//TODO: implement storage pool election
	nodes, err := r.Reachables()
	if err != nil {
		return Pool{}, err
	}

	h := make([]chain.NodeHeader, 0)
	for _, n := range nodes {
		h = append(h, chain.NewNodeHeader(n.publicKey, !n.isReachable, false, n.patch.patchid, n.status == NodeOK))
	}

	return Pool{
		nodes:   nodes,
		headers: h,
	}, err
}

//FindValidationPool lookups a validation pool from a transaction hash and a required number using the entropy sKeys
func FindValidationPool(txHash crypto.VersionnedHash, minValidations int, masterNodeKey crypto.PublicKey, nodeReader NodeReader, sharedKeyReader shared.KeyReader) (vPool Pool, err error) {
	authKeys, err := sharedKeyReader.AuthorizedNodesPublicKeys()
	if err != nil {
		return
	}

	firstKeys, err := sharedKeyReader.FirstNodeCrossKeypair()
	if err != nil {
		return
	}

	sKeys, err := entropySort(txHash, authKeys, firstKeys.PrivateKey())
	if err != nil {
		return
	}

	requiredPatchNb, err := validationRequiredPatchNumber(minValidations, nodeReader)
	if err != nil {
		return
	}

	nbReachables := 0
	patchIds := make([]int, 0)

	var i int

	//challenge the validations nodes by providing more nodes validations
	maxNbValidations := minValidations + (minValidations / 2)

	//adding the master node to the validation headers
	masterNode, err := nodeReader.FindByPublicKey(masterNodeKey)
	if err != nil {
		return
	}
	vPool.headers = append(vPool.headers, chain.NewNodeHeader(masterNodeKey, false, true, masterNode.patch.patchid, masterNode.status == NodeOK))

	for nbReachables < maxNbValidations && len(patchIds) < requiredPatchNb && i < len(sKeys) {
		n, err := nodeReader.FindByPublicKey(sKeys[i])
		if err != nil {
			return Pool{}, err
		}

		//Add a validation headers
		vPool.headers = append(vPool.headers,
			chain.NewNodeHeader(sKeys[i],
				!n.isReachable,
				false,
				n.patch.patchid,
				n.status == NodeOK))

		//Add the node to the pool
		vPool.nodes = append(vPool.nodes, n)

		//Need a view of the reachable and unreachables for a better validation
		if n.isReachable {

			//Reference the patch of the node if it's not already insert by helping to determinate
			//the number of distinct patches retrieved for the check of the required number of patches
			var existingPatch bool
			for _, id := range patchIds {
				if id == n.patch.patchid {
					existingPatch = true
					break
				}
			}
			if !existingPatch {
				patchIds = append(patchIds, n.patch.patchid)
			}

			nbReachables++
		}

		i++
	}

	if nbReachables != maxNbValidations {
		return Pool{}, errors.New("cannot proceed transaction with an invalid number of reachables validation nodes")
	}

	if len(patchIds) != requiredPatchNb {
		return Pool{}, errors.New("cannot proceed transaction with missing patches validation nodes")
	}

	return vPool, nil
}

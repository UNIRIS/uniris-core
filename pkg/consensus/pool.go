package consensus

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net"
	"sort"

	"github.com/uniris/uniris-core/pkg/shared"

	"github.com/uniris/uniris-core/pkg/chain"
)

//Pool represent a pool either for sharding or validation
type Pool []Node

//PoolRequester handles the request to perform on a pool during the mining
type PoolRequester interface {
	//RequestTransactionTimeLock asks a pool to timelock a transaction using the address related
	RequestTransactionTimeLock(pool Pool, txHash string, txAddr string, masterPublicKey string) error

	//RequestTransactionValidations asks a pool to validation a transaction
	RequestTransactionValidations(pool Pool, tx chain.Transaction, minValid int, masterValid chain.MasterValidation) ([]chain.Validation, error)

	//RequestTransactionStorage asks a pool to store a transaction
	RequestTransactionStorage(pool Pool, minStorage int, tx chain.Transaction) error

	//RequestLastTransaction asks a pool to retrieve the last transaction of an address
	RequestLastTransaction(pool Pool, txAddr string, txType chain.TransactionType) (*chain.Transaction, error)
}

//FindMasterNodes finds a list of master nodes by using an entropy sorting based on the transaction and minimum number of master
func FindMasterNodes(txHash string, r NodeReader, sharedReader shared.NodeReader) (Pool, error) {
	authKeys, err := sharedReader.AuthorizedPublicKeys()
	if err != nil {
		return nil, err
	}

	nodeFirstKeys, err := sharedReader.NodeFirstKeys()
	if err != nil {
		return nil, err
	}

	sort := entropySort(txHash, authKeys, nodeFirstKeys.PrivateKey())
	nbReachables, err := r.CountReachables()
	if err != nil {
		return nil, err
	}
	nbMasters := requiredNumberOfMaster(len(authKeys), nbReachables)

	masters := make([]Node, 0)
	nbReachableMasters := 0

	var i int
	for nbReachableMasters < nbMasters && i < len(sort) {
		n, err := r.FindByPublicKey(sort[i])
		if err != nil {
			return nil, err
		}
		masters = append(masters, n)
		if n.isReachable {
			nbReachableMasters++
		}
		i++
	}
	if nbReachableMasters != nbMasters {
		return nil, errors.New("cannot proceed transaction with an invalid number of reachables master nodes")
	}
	return masters, nil
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
func buildStartingPoint(txHash string, nodeMinerPrivateKey string) string {
	h := hmac.New(sha256.New, []byte(nodeMinerPrivateKey))
	h.Write([]byte(txHash))
	return hex.EncodeToString(h.Sum(nil))
}

//entropySort sorts a list of nodes public keys using a "starting point" (HMAC of the transaction hash with the first node shared private key) and the hashes of the node public keys
func entropySort(txHash string, authKeys []string, nodeFirstKey string) (sortedKeys []string) {
	startingPoint := buildStartingPoint(txHash, nodeFirstKey)

	//Building list of public keys and map of hash-›key
	hashKeys := make([]string, len(authKeys))
	mHashKeys := make(map[string]string, 0)
	for i, k := range authKeys {
		h := sha256.New()
		h.Write([]byte(k))
		hash := hex.EncodeToString(h.Sum(nil))
		mHashKeys[hash] = k
		hashKeys[i] = hash
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

	maxpos := 64 //64 is the size of a sha256 hash

	var p int
	//Sort keys by comparing the last character of the key with a starting point character
	for len(sortedKeys) < len(hashKeys)-1 && p < maxpos {

		//iterating from the starting point to the end of the list
		//add add the key if the latest character matchew the start point position
		for i := startPointIndex + 1; i < len(hashKeys); i++ {
			if []rune(hashKeys[i])[maxpos-1] == []rune(startingPoint)[p] {
				var contains bool
				var j int
				for !contains && j < len(sortedKeys) {
					if sortedKeys[j] == mHashKeys[hashKeys[i]] {
						contains = true
						break
					}
					j++
				}
				if !contains {
					sortedKeys = append(sortedKeys, mHashKeys[hashKeys[i]])
				}
			}
		}

		//iterating from the 0 to the starting point
		//and add the key if the latest character matches the start point position
		for i := 0; i < startPointIndex; i++ {
			if []rune(hashKeys[i])[maxpos-1] == []rune(startingPoint)[p] {
				var contains bool
				var j int
				for !contains && j < len(sortedKeys) {
					if sortedKeys[j] == mHashKeys[hashKeys[i]] {
						contains = true
					}
					j++
				}
				if !contains {
					sortedKeys = append(sortedKeys, mHashKeys[hashKeys[i]])
				}
			}
		}

		//We advance on the starting point character
		p++
	}

	//We have tested all the characters of the staring point and not yet finished the sorting operation, we will loop on all the hex characters to finish sorting
	if len(sortedKeys) < len(hashKeys)-1 {
		hexChar := []rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f'}
		var p int
		for len(sortedKeys) < len(hashKeys)-1 && p < len(hexChar) {

			//iterating from the starting point to the end of the list
			//add add the key if the latest character matchew the start point position
			for i := startPointIndex + 1; i < len(hashKeys); i++ {
				if []rune(hashKeys[i])[maxpos-1] == hexChar[p] {
					var contains bool
					var j int
					for !contains && j < len(sortedKeys) {
						if sortedKeys[j] == mHashKeys[hashKeys[i]] {
							contains = true
							break
						}
						j++
					}
					if !contains {
						sortedKeys = append(sortedKeys, mHashKeys[hashKeys[i]])
					}
				}
			}

			//iterating from the 0 to the starting point
			//and add the key if the latest character matches the start point position
			for i := 0; i < startPointIndex; i++ {
				if []rune(hashKeys[i])[maxpos-1] == hexChar[p] {
					var contains bool
					var j int
					for !contains && j < len(sortedKeys) {
						if sortedKeys[j] == mHashKeys[hashKeys[i]] {
							contains = true
						}
						j++
					}
					if !contains {
						sortedKeys = append(sortedKeys, mHashKeys[hashKeys[i]])
					}
				}
			}

			//We advance on the hexadecimal characters
			p++
		}
	}

	return
}

//FindStoragePool searches a storage pool for the given address
//TODO: Implements AI lookups to identify the right storage pool
func FindStoragePool(address string) (Pool, error) {
	return Pool{
		Node{
			ip:        net.ParseIP("127.0.0.1"),
			port:      5000,
			publicKey: "3059301306072a8648ce3d020106082a8648ce3d0301070342000408f4b4026d2560aaa552244bdf8ec421bb41378b56487d9d4ca5a57fd6e64ef7ae2f2c6530f18bd0f359342b4fa7fdaeaa60c45a1197260eb1c267cc996bec81",
		},
	}, nil
}

//FindValidationPool searches a validation pool from a transaction hash
//TODO: Implements AI lookups to identify the right validation pool
func FindValidationPool(tx chain.Transaction) (Pool, error) {
	return Pool{
		Node{
			ip:        net.ParseIP("127.0.0.1"),
			port:      5000,
			publicKey: "3059301306072a8648ce3d020106082a8648ce3d0301070342000408f4b4026d2560aaa552244bdf8ec421bb41378b56487d9d4ca5a57fd6e64ef7ae2f2c6530f18bd0f359342b4fa7fdaeaa60c45a1197260eb1c267cc996bec81",
		},
	}, nil
}

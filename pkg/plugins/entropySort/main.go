package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"sort"
)

//buildStartingPoint creates a starting point by using an HMAC of the transaction hash and the first node shared private key
func buildStartingPoint(txHash []byte, nodeCrossFirstPvKey []byte) (string, error) {
	h := hmac.New(sha256.New, nodeCrossFirstPvKey)
	h.Write(txHash)
	return hex.EncodeToString(h.Sum(nil)), nil
}

func hash(data []byte) []byte {
	h := sha256.New()
	h.Write(data)
	return h.Sum(nil)
}

//EntropySort sorts a list of nodes public keys using a "starting point" (HMAC of the transaction hash with the first node shared private key) and the hashes of the node public keys
func EntropySort(txHash []byte, authKeys [][]byte, nodeCrossFirstPvKey []byte) (sortedKeys [][]byte, err error) {
	startingPoint, err := buildStartingPoint(txHash, nodeCrossFirstPvKey)
	if err != nil {
		return nil, err
	}

	//Building list of public keys and map of hash-›key
	hashKeys := make([]string, len(authKeys))
	mHashKeys := make(map[string][]byte, 0)
	for i, k := range authKeys {
		h := hash(k)
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

	end := sha256.Size

	//Sort keys by comparing the last character of the key with a starting point character
	for p := 0; len(sortedKeys) < len(hashKeys)-1 && p < end; p++ {

		//iterating from the starting point to the end of the list
		//add add the key if the latest character matchew the start point position
		for i := startPointIndex + 1; i < len(hashKeys); i++ {
			if []rune(hashKeys[i])[end-1] == []rune(startingPoint)[p] {
				var contains bool
				for _, k := range sortedKeys {
					if bytes.Equal(k, mHashKeys[hashKeys[i]]) {
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
					if bytes.Equal(k, mHashKeys[hashKeys[i]]) {
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

	//We have tested all the characters of the staring point and not yet finished the sort operation, we will loop on all the hex characters to finish the sort
	if len(sortedKeys) < len(hashKeys)-1 {
		hexChar := []rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f'}
		for p := 0; len(sortedKeys) < len(hashKeys)-1 && p < len(hexChar); p++ {

			//iterating from the starting point to the end of the list
			//add add the key if the latest character matchew the start point position
			for i := startPointIndex + 1; i < len(hashKeys); i++ {
				if []rune(hashKeys[i])[end-1] == hexChar[p] {
					var contains bool
					for _, k := range sortedKeys {
						if bytes.Equal(k, mHashKeys[hashKeys[i]]) {
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
						if bytes.Equal(k, mHashKeys[hashKeys[i]]) {
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

package main

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"
	"time"
)

var timeLockCountdown = 60 * time.Second

type timeLocker struct {
	txHash         []byte
	txAddress      []byte
	coordPublicKey []byte
	end            time.Time
	ticker         *time.Ticker
}

var timeLockers []timeLocker
var timerLockMux = &sync.Mutex{}

//TimeLockTransaction write a timelock.
//if a lock exists already an error is returned
//the timelock will be removed when the countdown is reached
func TimeLockTransaction(txHash []byte, txAddr []byte, coordPublicKey []byte) error {
	if _, found, _ := findTimelock(txHash, txAddr); found {
		return errors.New("a lock already exist for this transaction")
	}

	ticker := time.NewTicker(1 * time.Second)
	end := time.Now().Add(timeLockCountdown)
	tLocker := timeLocker{
		txAddress:      txAddr,
		txHash:         txHash,
		coordPublicKey: coordPublicKey,
		end:            end,
		ticker:         ticker,
	}
	timerLockMux.Lock()
	timeLockers = append(timeLockers, tLocker)
	timerLockMux.Unlock()

	//Remove the timelock when the countdown is reached
	go func() {
		for range ticker.C {
			if time.Now().Unix() == end.Unix() {
				removeTimeLock(txHash, txAddr)
			}
		}
	}()

	return nil
}

//ContainsTransactionTimeLock checks if a transaction timelock exists
func ContainsTransactionTimeLock(txHash []byte, txAddr []byte) bool {
	_, found, _ := findTimelock(txHash, txAddr)
	return found
}

func findTimelock(txHash, txAddr []byte) (tl timeLocker, found bool, index int) {
	for index, l := range timeLockers {
		if bytes.Equal(l.txHash, txHash) && bytes.Equal(l.txAddress, txAddr) {
			return l, true, index
		}
	}
	return
}

func transactionHashTimeLocked(txHash []byte) bool {
	for _, l := range timeLockers {
		if bytes.Equal(l.txHash, txHash) {
			return true
		}
	}
	return false
}

func removeTimeLock(txHash []byte, txAddr []byte) {
	tl, found, index := findTimelock(txHash, txAddr)
	if !found {
		return
	}

	timerLockMux.Lock()
	tl.ticker.Stop()
	timeLockers[index] = timeLockers[len(timeLockers)-1]
	timeLockers = timeLockers[:len(timeLockers)-1]
	timerLockMux.Unlock()

	fmt.Printf("transaction %s is unlocked\n", hex.EncodeToString(txHash))
}

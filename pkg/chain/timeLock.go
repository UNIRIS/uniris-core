package chain

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/uniris/uniris-core/pkg/crypto"
)

var timeLockCountdown time.Duration = 60 * time.Second

type timeLocker struct {
	txHash          crypto.VersionnedHash
	txAddress       crypto.VersionnedHash
	masterPublicKey crypto.PublicKey
	end             time.Time
	ticker          *time.Ticker
}

var timeLockers []timeLocker
var timerLockMux = &sync.Mutex{}

//TimeLockTransaction write a timelock.
//if a lock exists already an error is returned
//the timelock will be removed when the countdown is reached
func TimeLockTransaction(txHash crypto.VersionnedHash, txAddr crypto.VersionnedHash, masterPubk crypto.PublicKey) error {
	if _, found, _ := findTimelock(txHash, txAddr); found {
		return errors.New("a lock already exist for this transaction")
	}

	ticker := time.NewTicker(1 * time.Second)
	end := time.Now().Add(timeLockCountdown)
	tLocker := timeLocker{
		txAddress:       txAddr,
		txHash:          txHash,
		masterPublicKey: masterPubk,
		end:             end,
		ticker:          ticker,
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

//ContainsTimeLock checks if a transaction timelock exists
func ContainsTimeLock(txHash crypto.VersionnedHash, txAddr crypto.VersionnedHash) bool {
	_, found, _ := findTimelock(txHash, txAddr)
	return found
}

func findTimelock(txHash, txAddr crypto.VersionnedHash) (tl timeLocker, found bool, index int) {
	for index, l := range timeLockers {
		if bytes.Equal(l.txHash, txHash) && bytes.Equal(l.txAddress, txAddr) {
			return l, true, index
		}
	}
	return
}

func transactionHashTimeLocked(txHash crypto.VersionnedHash) bool {
	for _, l := range timeLockers {
		if bytes.Equal(l.txHash, txHash) {
			return true
		}
	}
	return false
}

func removeTimeLock(txHash crypto.VersionnedHash, txAddr crypto.VersionnedHash) {
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

package chain

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/uniris/uniris-core/pkg/crypto"
)

var timeLockCountdown time.Duration = 60 * time.Second

type timeLocker struct {
	txHash          string
	txAddress       string
	masterPublicKey string
	end             time.Time
	ticker          *time.Ticker
}

var timeLockers []timeLocker
var timerLockMux = &sync.Mutex{}

//TimeLockTransaction write a timelock.
//if a lock exists already an error is returned
//the timelock will be removed when the countdown is reached
func TimeLockTransaction(txHash string, txAddr string, masterPubk string) error {
	if _, err := crypto.IsHash(txHash); err != nil {
		return fmt.Errorf("lock transaction hash: %s", err.Error())
	}

	if _, err := crypto.IsHash(txAddr); err != nil {
		return fmt.Errorf("lock transaction address: %s", err.Error())
	}

	if _, err := crypto.IsPublicKey(masterPubk); err != nil {
		return fmt.Errorf("lock transaction public key: %s", err.Error())
	}

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
func ContainsTimeLock(txHash string, txAddr string) bool {
	_, found, _ := findTimelock(txHash, txAddr)
	return found
}

func findTimelock(txHash, txAddr string) (tl timeLocker, found bool, index int) {
	for index, l := range timeLockers {
		if l.txHash == txHash && l.txAddress == txAddr {
			return l, true, index
		}
	}
	return
}

func transactionHashTimeLocked(txHash string) bool {
	for _, l := range timeLockers {
		if l.txHash == txHash {
			return true
		}
	}
	return false
}

func removeTimeLock(txHash string, txAddr string) {
	tl, found, index := findTimelock(txHash, txAddr)
	if !found {
		return
	}

	timerLockMux.Lock()
	tl.ticker.Stop()
	timeLockers[index] = timeLockers[len(timeLockers)-1]
	timeLockers = timeLockers[:len(timeLockers)-1]
	timerLockMux.Unlock()

	fmt.Printf("transaction %s is unlocked\n", txHash)
}

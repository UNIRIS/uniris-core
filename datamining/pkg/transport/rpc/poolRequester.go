package rpc

import (
	"errors"
	"log"
	"sync"
	"sync/atomic"

	datamining "github.com/uniris/uniris-core/datamining/pkg"
	"github.com/uniris/uniris-core/datamining/pkg/account"
	"github.com/uniris/uniris-core/datamining/pkg/lock"
	"github.com/uniris/uniris-core/datamining/pkg/mining"
	"github.com/uniris/uniris-core/datamining/pkg/system"
)

//PoolRequester define methods for pool requesting
type PoolRequester interface {
	mining.PoolRequester
	account.PoolRequester
}

type poolR struct {
	cli    ExternalClient
	conf   system.UnirisConfig
	crypto Crypto
	api    apiBuilder
	data   dataBuilder
}

//NewPoolRequester creates a new pool requester using GRPC
func NewPoolRequester(cli ExternalClient, conf system.UnirisConfig, crypto Crypto) PoolRequester {
	return poolR{
		cli:    cli,
		conf:   conf,
		crypto: crypto,
		api:    apiBuilder{},
		data:   dataBuilder{},
	}
}

func (pR poolR) RequestID(sPool datamining.Pool, idHash string) (account.EndorsedID, error) {
	ids := make([]account.EndorsedID, 0)

	var wg sync.WaitGroup
	wg.Add(len(sPool.Peers()))

	idChan := make(chan account.EndorsedID)

	for _, p := range sPool.Peers() {
		go func(p datamining.Peer) {
			defer wg.Done()

			id, err := pR.cli.RequestID(p.IP.String(), idHash)
			if err != nil {
				log.Printf("Unexpected error during ID requesting for the peer %s\n", p.IP.String())
				log.Printf("Details: %s\n", err.Error())
				return
			}

			idChan <- id
		}(p)
	}

	go func() {
		wg.Wait()
		close(idChan)
	}()

	for id := range idChan {
		ids = append(ids, id)
	}

	if len(ids) == 0 {
		return nil, errors.New(pR.conf.Services.Datamining.Errors.AccountNotExist)
	}

	//Checks the consistency of the retrieved results
	firstHash, err := pR.crypto.hasher.HashEndorsedID(ids[0])
	if err != nil {
		return nil, err
	}
	for i := 1; i < len(ids); i++ {
		hash, err := pR.crypto.hasher.HashEndorsedID(ids[i])
		if err != nil {
			return nil, err
		}
		if hash != firstHash {
			return nil, errors.New("Inconsistent data")
		}
	}

	return ids[0], nil
}

func (pR poolR) RequestKeychain(sPool datamining.Pool, encAddress string) (account.EndorsedKeychain, error) {

	keychains := make([]account.EndorsedKeychain, 0)

	var wg sync.WaitGroup
	wg.Add(len(sPool.Peers()))

	kcChan := make(chan account.EndorsedKeychain)

	for _, p := range sPool.Peers() {
		go func(p datamining.Peer) {
			defer wg.Done()
			kc, err := pR.cli.RequestKeychain(p.IP.String(), encAddress)
			if err != nil {
				log.Printf("Unexpected error during keychain requesting for the peer %s\n", p.IP.String())
				log.Printf("Details: %s\n", err.Error())
				return
			}

			kcChan <- kc
		}(p)
	}

	go func() {
		wg.Wait()
		close(kcChan)
	}()

	for kc := range kcChan {
		keychains = append(keychains, kc)
	}

	if len(keychains) == 0 {
		return nil, errors.New(pR.conf.Services.Datamining.Errors.AccountNotExist)
	}

	//Checks the consistency of the retrieved results
	firstHash, err := pR.crypto.hasher.HashEndorsedKeychain(keychains[0])
	if err != nil {
		return nil, err
	}
	for i := 1; i < len(keychains); i++ {
		hash, err := pR.crypto.hasher.HashEndorsedKeychain(keychains[i])
		if err != nil {
			return nil, err
		}
		if hash != firstHash {
			return nil, errors.New("Inconsistent data")
		}
	}

	return keychains[0], nil
}

func (pR poolR) RequestLock(lastValidPool datamining.Pool, txLock lock.TransactionLock) error {

	var wg sync.WaitGroup
	wg.Add(len(lastValidPool.Peers()))

	var ackLock int32
	lockChan := make(chan bool)

	for _, p := range lastValidPool.Peers() {
		go func(p datamining.Peer) {
			defer wg.Done()
			if err := pR.cli.RequestLock(p.IP.String(), txLock); err != nil {
				log.Printf("Unexpected error during lock requesting for the peer %s\n", p.IP.String())
				log.Printf("Details: %s\n", err.Error())
				return
			}
			atomic.AddInt32(&ackLock, 1)
		}(p)
	}

	wg.Wait()
	close(lockChan)

	//TODO: specify minium required locks
	minLocks := 1
	lockFinal := atomic.LoadInt32(&ackLock)
	if int(lockFinal) < minLocks {
		return errors.New("Transaction locking failed")
	}

	return nil
}
func (pR poolR) RequestUnlock(lastValidPool datamining.Pool, txLock lock.TransactionLock) error {

	var wg sync.WaitGroup
	wg.Add(len(lastValidPool.Peers()))

	var ackUnLock int32

	for _, p := range lastValidPool.Peers() {
		go func(p datamining.Peer) {
			defer wg.Done()
			if err := pR.cli.RequestUnlock(p.IP.String(), txLock); err != nil {
				log.Printf("Unexpected error during unlock requesting for the peer %s\n", p.IP.String())
				log.Printf("Details: %s\n", err.Error())
				return
			}
			atomic.AddInt32(&ackUnLock, 1)
		}(p)
	}

	wg.Wait()

	//TODO: specify minium required locks
	minUnlocks := 1
	unlockFinal := atomic.LoadInt32(&ackUnLock)
	if int(unlockFinal) < minUnlocks {
		return errors.New("Transaction unlocking failed")
	}

	return nil
}

func (pR poolR) RequestValidations(minValids int, validPool datamining.Pool, txHash string, data interface{}, txType mining.TransactionType) ([]mining.Validation, error) {
	valids := make([]mining.Validation, 0)

	var wg sync.WaitGroup
	wg.Add(len(validPool.Peers()))

	vChan := make(chan mining.Validation)

	for _, p := range validPool.Peers() {
		go func(p datamining.Peer) {
			defer wg.Done()
			v, err := pR.cli.RequestValidation(p.IP.String(), txType, txHash, data)
			if err != nil {
				log.Printf("Unexpected error during validation requesting for the peer %s\n", p.IP.String())
				log.Printf("Details: %s\n", err.Error())
				return
			}

			vChan <- v
		}(p)
	}

	go func() {
		wg.Wait()
		close(vChan)
	}()

	for v := range vChan {
		valids = append(valids, v)
	}

	if len(valids) < minValids {
		return nil, errors.New("Minimum validations are not reached")
	}

	return valids, nil
}

func (pR poolR) RequestStorage(minReplicas int, sPool datamining.Pool, data interface{}, end mining.Endorsement, txType mining.TransactionType) error {

	var wg sync.WaitGroup
	wg.Add(len(sPool.Peers()))

	var ackStore int32

	for _, p := range sPool.Peers() {
		go func(p datamining.Peer) {
			defer wg.Done()
			if err := pR.cli.RequestStorage(p.IP.String(), txType, data, end); err != nil {
				log.Printf("Unexpected error during storage requesting for the peer %s\n", p.IP.String())
				log.Printf("Details: %s\n", err.Error())
				return
			}
			atomic.AddInt32(&ackStore, 1)
		}(p)
	}

	wg.Wait()

	ackStoreFinal := atomic.LoadInt32(&ackStore)
	if int(ackStoreFinal) < minReplicas {
		return errors.New("Transaction storage failed")
	}

	return nil
}

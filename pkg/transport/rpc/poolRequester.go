package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/uniris/uniris-core/pkg/transaction"

	api "github.com/uniris/uniris-core/api/protobuf-spec"
	"github.com/uniris/uniris-core/pkg/crypto"
	"google.golang.org/grpc"
)

type poolR struct {
	sharedPubk string
	sharedPvk  string
}

//NewPoolRequester creates a new pool requester as a GRPC client
func NewPoolRequester(sharedPubk string, sharedPvk string) transaction.PoolRequester {
	return poolR{
		sharedPubk: sharedPubk,
		sharedPvk:  sharedPvk,
	}
}

func (pr poolR) RequestTransactionLock(pool transaction.Pool, lock transaction.Lock) error {

	var wg sync.WaitGroup
	wg.Add(len(pool))

	var ackUnlock int32

	req := &api.LockRequest{
		Address:             lock.Address(),
		TransactionHash:     lock.TransactionHash(),
		MasterPeerPublicKey: lock.MasterRobotKey(),
		Timestamp:           time.Now().Unix(),
	}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return err
	}
	sig, err := crypto.Sign(string(reqBytes), pr.sharedPvk)
	if err != nil {
		return err
	}
	req.SignatureRequest = sig

	for _, p := range pool {
		go func(p transaction.PoolMember) {
			defer wg.Done()

			serverAddr := fmt.Sprintf("%s:%d", p.IP(), p.Port())
			conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
			if err != nil {
				fmt.Printf("LOCK TRANSACTION - ERROR: %s", err.Error())
			}
			defer conn.Close()
			cli := api.NewTransactionServiceClient(conn)
			res, err := cli.LockTransaction(context.Background(), req)
			if err != nil {
				fmt.Printf("LOCK TRANSACTION RESPONSE - ERROR: %s", err.Error())
				return
			}

			fmt.Printf("LOCK TRANSACTION RESPONSE - %s", time.Unix(res.Timestamp, 0).String())

			resBytes, err := json.Marshal(&api.LockResponse{
				Timestamp: req.Timestamp,
			})
			if err != nil {
				fmt.Printf("LOCK TRANSACTION RESPONSE - ERROR: %s", err.Error())
				return
			}
			if err := crypto.VerifySignature(string(resBytes), pr.sharedPubk, res.SignatureResponse); err != nil {
				fmt.Printf("LOCK TRANSACTION RESPONSE - ERROR: %s", err.Error())
				return
			}

			atomic.AddInt32(&ackUnlock, 1)
		}(p)
	}

	wg.Wait()

	//TODO: specify minium required locks
	minLocks := 1
	lockFinal := atomic.LoadInt32(&ackUnlock)
	if int(lockFinal) < minLocks {
		return errors.New("number of locks are not reached")
	}

	return nil
}

func (pr poolR) RequestTransactionUnlock(pool transaction.Pool, lock transaction.Lock) error {
	var wg sync.WaitGroup
	wg.Add(len(pool))

	var ackUnlock int32

	req := &api.LockRequest{
		Address:             lock.Address(),
		TransactionHash:     lock.TransactionHash(),
		MasterPeerPublicKey: lock.MasterRobotKey(),
		Timestamp:           time.Now().Unix(),
	}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return err
	}
	sig, err := crypto.Sign(string(reqBytes), pr.sharedPvk)
	if err != nil {
		return err
	}
	req.SignatureRequest = sig

	for _, p := range pool {
		go func(p transaction.PoolMember) {
			defer wg.Done()

			serverAddr := fmt.Sprintf("%s:%d", p.IP(), p.Port())
			conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
			if err != nil {
				fmt.Printf("UNLOCK TRANSACTION - ERROR: %s", err.Error())
			}
			defer conn.Close()
			cli := api.NewTransactionServiceClient(conn)
			res, err := cli.UnlockTransaction(context.Background(), req)
			if err != nil {
				fmt.Printf("UNLOCK TRANSACTION RESPONSE - ERROR: %s", err.Error())
				return
			}

			fmt.Printf("UNLOCK TRANSACTION RESPONSE - %s", time.Unix(res.Timestamp, 0).String())

			resBytes, err := json.Marshal(&api.LockResponse{
				Timestamp: req.Timestamp,
			})
			if err != nil {
				fmt.Printf("UNLOCK TRANSACTION RESPONSE - ERROR: %s", err.Error())
				return
			}
			if err := crypto.VerifySignature(string(resBytes), pr.sharedPubk, res.SignatureResponse); err != nil {
				fmt.Printf("UNLOCK TRANSACTION RESPONSE - ERROR: %s", err.Error())
				return
			}

			atomic.AddInt32(&ackUnlock, 1)
		}(p)
	}

	wg.Wait()

	//TODO: specify minium required unlocks
	minUnlocks := 1
	unlockFinal := atomic.LoadInt32(&ackUnlock)
	if int(unlockFinal) < minUnlocks {
		return errors.New("number of unlocks are not reached")
	}

	return nil
}

func (pr poolR) RequestTransactionValidations(pool transaction.Pool, tx transaction.Transaction, masterValid transaction.MasterValidation, validChan chan<- transaction.MinerValidation) {

	req := &api.ConfirmValidationRequest{
		MasterValidation: formatAPIMasterValidation(masterValid),
		Transaction:      formatAPITransaction(tx),
		Timestamp:        time.Now().Unix(),
	}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		fmt.Printf("CONFIRM VALIDATION TRANSACTION REQUEST - ERROR: %s", err.Error())
		return
	}
	sig, err := crypto.Sign(string(reqBytes), pr.sharedPvk)
	if err != nil {
		fmt.Printf("CONFIRM VALIDATION TRANSACTION REQUEST - ERROR: %s", err.Error())
		return
	}
	req.SignatureRequest = sig

	for _, p := range pool {

		serverAddr := fmt.Sprintf("%s:%d", p.IP().String(), p.Port())
		conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
		if err != nil {

		}
		defer conn.Close()
		cli := api.NewTransactionServiceClient(conn)
		res, err := cli.ConfirmTransactionValidation(context.Background(), req)
		if err != nil {
			fmt.Printf("CONFIRM VALIDATION TRANSACTION RESPONSE - ERROR: %s", err.Error())
			return
		}

		fmt.Printf("CONFIRM VALIDATION TRANSACTION RESPONSE - %s", time.Unix(res.Timestamp, 0).String())

		resBytes, err := json.Marshal(&api.ConfirmValidationResponse{
			Timestamp:  res.Timestamp,
			Validation: res.Validation,
		})
		if err != nil {
			fmt.Printf("CONFIRM VALIDATION TRANSACTION RESPONSE - ERROR: %s", err.Error())
			return
		}
		if err := crypto.VerifySignature(string(resBytes), pr.sharedPubk, res.SignatureResponse); err != nil {
			fmt.Printf("CONFIRM VALIDATION TRANSACTION RESPONSE - ERROR: %s", err.Error())
			return
		}

		v, err := transaction.NewMinerValidation(transaction.ValidationStatus(res.Validation.Status), time.Unix(res.Timestamp, 0), res.Validation.PublicKey, res.Validation.Signature)
		if err != nil {
			return
		}
		validChan <- v
	}

}

func (pr poolR) RequestTransactionStorage(pool transaction.Pool, tx transaction.Transaction, ackChan chan<- bool) {
	confValids := make([]*api.MinerValidation, 0)
	for _, v := range tx.ConfirmationsValidations() {
		confValids = append(confValids, formatAPIValidation(v))
	}

	req := &api.StoreRequest{
		MinedTransaction: &api.MinedTransaction{
			MasterValidation:   formatAPIMasterValidation(tx.MasterValidation()),
			ConfirmValidations: confValids,
			Transaction:        formatAPITransaction(tx),
		},
		Timestamp: time.Now().Unix(),
	}

	reqBytes, err := json.Marshal(req)
	if err != nil {
		fmt.Printf("STORE TRANSACTION REQUEST - ERROR: %s", err.Error())
		return
	}
	sig, err := crypto.Sign(string(reqBytes), pr.sharedPvk)
	if err != nil {
		fmt.Printf("STORE TRANSACTION REQUEST - ERROR: %s", err.Error())
		return
	}

	req.SignatureRequest = sig

	for _, p := range pool {
		serverAddr := fmt.Sprintf("%s:%d", p.IP().String(), p.Port())
		conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
		if err != nil {

		}
		defer conn.Close()
		cli := api.NewTransactionServiceClient(conn)
		res, err := cli.StoreTransaction(context.Background(), req)
		if err != nil {
			fmt.Printf("STORE TRANSACTION RESPONSE - ERROR: %s", err.Error())
			return
		}

		fmt.Printf("STORE TRANSACTION RESPONSE - %s", time.Unix(res.Timestamp, 0).String())

		resBytes, err := json.Marshal(&api.StoreResponse{
			Timestamp: res.Timestamp,
		})
		if err != nil {
			fmt.Printf("STORE TRANSACTION RESPONSE - ERROR: %s", err.Error())
			return
		}
		if err := crypto.VerifySignature(string(resBytes), pr.sharedPubk, res.SignatureResponse); err != nil {
			fmt.Printf("STORE TRANSACTION RESPONSE - ERROR: %s", err.Error())
			return
		}

		ackChan <- true
	}
}

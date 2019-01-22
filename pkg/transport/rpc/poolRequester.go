package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	api "github.com/uniris/uniris-core/api/protobuf-spec"
	uniris "github.com/uniris/uniris-core/pkg"
	"github.com/uniris/uniris-core/pkg/crypto"
	"github.com/uniris/uniris-core/pkg/pooling"
	"google.golang.org/grpc"
)

type poolR struct {
	sharedPubk string
	sharedPvk  string
}

//NewPoolRequester creates a new pool requester as a GRPC client
func NewPoolRequester(sharedPubk string, sharedPvk string) pooling.PoolRequester {
	return poolR{
		sharedPubk: sharedPubk,
		sharedPvk:  sharedPvk,
	}
}

func (pr poolR) RequestLastTransaction(pool pooling.Pool, addr string) (tx uniris.Transaction, err error) {
	return
}

func (pr poolR) RequestTransactionLock(pool pooling.Pool, txHash string, address string, masterPeerIP string) error {

	var wg sync.WaitGroup
	wg.Add(len(pool.Peers()))

	var ackLock int32
	lockChan := make(chan bool)

	req := &api.LockRequest{
		Address:         address,
		TransactionHash: txHash,
		MasterPeerIp:    masterPeerIP,
		Timestamp:       time.Now().Unix(),
	}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return err
	}
	sig, err := crypto.Sign(string(reqBytes), pr.sharedPubk)
	if err != nil {
		return err
	}
	req.SignatureRequest = sig

	for _, p := range pool.Peers() {
		go func(p string) {
			defer wg.Done()

			serverAddr := fmt.Sprintf("%s:1717", p)
			conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
			if err != nil {

			}
			defer conn.Close()
			cli := api.NewTransactionServiceClient(conn)
			res, err := cli.LockTransaction(context.Background(), req)
			if err != nil {
				fmt.Printf("LOCK TRANSACTION RESPONSE - ERROR: %s", err.Error())
				return
			}

			fmt.Printf("LOCK TRANSACTION RESPONSE - %s", time.Unix(res.Timestamp, 0).String())

			resBytes, err := json.Marshal(res)
			if err != nil {
				fmt.Printf("LOCK TRANSACTION RESPONSE UNMARSHALING - ERROR: %s", err.Error())
				return
			}
			if err := crypto.VerifySignature(string(resBytes), pr.sharedPubk, res.SignatureResponse); err != nil {
				fmt.Printf("LOCK TRANSACTION RESPONSE VERIFICATION - ERROR: %s", err.Error())
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

func (pr poolR) RequestTransactionUnlock(pool pooling.Pool, txHash string, address string) error {
	return nil
}

func (pr poolR) RequestTransactionValidations(pool pooling.Pool, tx uniris.Transaction, masterValid uniris.MasterValidation, validChan chan<- uniris.MinerValidation, replyChan chan<- bool) {

	req := &api.ConfirmValidationRequest{
		MasterValidation: formatAPIMasterValidationAPI(masterValid),
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

	for _, p := range pool.Peers() {

		serverAddr := fmt.Sprintf("%s:1717", p)
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
		replyChan <- true

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

		validChan <- uniris.NewMinerValidation(uniris.ValidationStatus(res.Validation.Status), time.Unix(res.Timestamp, 0), res.Validation.PublicKey, res.Validation.Signature)
	}

}

func (pr poolR) RequestTransactionStorage(pool pooling.Pool, tx uniris.Transaction, ackChan chan<- bool) {
	confValids := make([]*api.MinerValidation, 0)
	for _, v := range tx.ConfirmationsValidations() {
		confValids = append(confValids, formatAPIValidation(v))
	}

	req := &api.StoreRequest{
		MinedTransaction: &api.MinedTransaction{
			MasterValidation:   formatAPIMasterValidationAPI(tx.MasterValidation()),
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

	for _, p := range pool.Peers() {
		serverAddr := fmt.Sprintf("%s:1717", p)
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

package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/uniris/uniris-core/pkg/shared"

	"google.golang.org/grpc/status"

	"github.com/uniris/uniris-core/pkg/transaction"

	api "github.com/uniris/uniris-core/api/protobuf-spec"
	"github.com/uniris/uniris-core/pkg/crypto"
	"google.golang.org/grpc"
)

type poolR struct {
	sharedService shared.Service
}

//NewPoolRequester creates a new pool requester as a GRPC client
func NewPoolRequester(sharedService shared.Service) transaction.PoolRequester {
	return poolR{
		sharedService: sharedService,
	}
}

func (pr poolR) RequestTransactionLock(pool transaction.Pool, lock transaction.Lock) error {

	lastMinerKeys, err := pr.sharedService.GetSharedMinerKeys()
	if err != nil {
		return err
	}

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
	sig, err := crypto.Sign(string(reqBytes), lastMinerKeys.PrivateKey())
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
				grpcErr, _ := status.FromError(err)
				fmt.Printf("LOCK TRANSACTION REQUEST - ERROR: %s\n", grpcErr.Message())
				return
			}
			defer conn.Close()
			cli := api.NewTransactionServiceClient(conn)
			res, err := cli.LockTransaction(context.Background(), req)
			if err != nil {
				grpcErr, _ := status.FromError(err)
				fmt.Printf("LOCK TRANSACTION RESPONSE - ERROR: %s\n", grpcErr.Message())
				return
			}

			fmt.Printf("LOCK TRANSACTION RESPONSE - %s\n", time.Unix(res.Timestamp, 0).String())

			resBytes, err := json.Marshal(&api.LockResponse{
				Timestamp: req.Timestamp,
			})
			if err != nil {
				fmt.Printf("LOCK TRANSACTION RESPONSE - ERROR: %s", err.Error())
				return
			}
			if err := crypto.VerifySignature(string(resBytes), lastMinerKeys.PublicKey(), res.SignatureResponse); err != nil {
				fmt.Printf("LOCK TRANSACTION RESPONSE - ERROR: %s\n", err.Error())
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

	lastMinerKeys, err := pr.sharedService.GetSharedMinerKeys()
	if err != nil {
		return err
	}

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
	sig, err := crypto.Sign(string(reqBytes), lastMinerKeys.PrivateKey())
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
				grpcErr, _ := status.FromError(err)
				fmt.Printf("UNLOCK TRANSACTION REQUEST - ERROR: %s\n", grpcErr.Message())
				return
			}
			defer conn.Close()
			cli := api.NewTransactionServiceClient(conn)
			res, err := cli.UnlockTransaction(context.Background(), req)
			if err != nil {
				grpcErr, _ := status.FromError(err)
				fmt.Printf("UNLOCK TRANSACTION RESPONSE - ERROR: %s\n", grpcErr.Message())
				return
			}

			fmt.Printf("UNLOCK TRANSACTION RESPONSE - %s\n", time.Unix(res.Timestamp, 0).String())

			resBytes, err := json.Marshal(&api.LockResponse{
				Timestamp: req.Timestamp,
			})
			if err != nil {
				fmt.Printf("UNLOCK TRANSACTION RESPONSE - ERROR: %s\n", err.Error())
				return
			}
			if err := crypto.VerifySignature(string(resBytes), lastMinerKeys.PublicKey(), res.SignatureResponse); err != nil {
				fmt.Printf("UNLOCK TRANSACTION RESPONSE - ERROR: %s\n", err.Error())
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

	lastMinerKeys, err := pr.sharedService.GetSharedMinerKeys()
	if err != nil {
		log.Printf("CONFIRM VALIDATION TRANSACTION REQUEST - ERROR: %s\n", err.Error())
		return
	}

	req := &api.ConfirmValidationRequest{
		MasterValidation: formatAPIMasterValidation(masterValid),
		Transaction:      formatAPITransaction(tx),
		Timestamp:        time.Now().Unix(),
	}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		fmt.Printf("CONFIRM VALIDATION TRANSACTION REQUEST - ERROR: %s\n", err.Error())
		return
	}
	sig, err := crypto.Sign(string(reqBytes), lastMinerKeys.PrivateKey())
	if err != nil {
		fmt.Printf("CONFIRM VALIDATION TRANSACTION REQUEST - ERROR: %s\n", err.Error())
		return
	}
	req.SignatureRequest = sig

	for _, p := range pool {

		serverAddr := fmt.Sprintf("%s:%d", p.IP().String(), p.Port())
		conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
		if err != nil {
			grpcErr, _ := status.FromError(err)
			fmt.Printf("CONFIRM VALIDATION TRANSACTION REQUEST - ERROR: %s\n", grpcErr.Message())
			return
		}
		defer conn.Close()
		cli := api.NewTransactionServiceClient(conn)
		res, err := cli.ConfirmTransactionValidation(context.Background(), req)
		if err != nil {
			grpcErr, _ := status.FromError(err)
			fmt.Printf("CONFIRM VALIDATION TRANSACTION RESPONSE - ERROR: %s\n", grpcErr.Message())
			return
		}

		fmt.Printf("CONFIRM VALIDATION TRANSACTION RESPONSE - %s\n", time.Unix(res.Timestamp, 0).String())

		resBytes, err := json.Marshal(&api.ConfirmValidationResponse{
			Timestamp:  res.Timestamp,
			Validation: res.Validation,
		})
		if err != nil {
			fmt.Printf("CONFIRM VALIDATION TRANSACTION RESPONSE - ERROR: %s\n", err.Error())
			return
		}
		if err := crypto.VerifySignature(string(resBytes), lastMinerKeys.PublicKey(), res.SignatureResponse); err != nil {
			fmt.Printf("CONFIRM VALIDATION TRANSACTION RESPONSE - ERROR: %s\n", err.Error())
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

	lastMinerKeys, err := pr.sharedService.GetSharedMinerKeys()
	if err != nil {
		log.Printf("STORE TRANSACTION REQUEST - ERROR: %s\n", err.Error())
		return
	}

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
		fmt.Printf("STORE TRANSACTION REQUEST - ERROR: %s\n", err.Error())
		return
	}
	sig, err := crypto.Sign(string(reqBytes), lastMinerKeys.PrivateKey())
	if err != nil {
		fmt.Printf("STORE TRANSACTION REQUEST - ERROR: %s\n", err.Error())
		return
	}

	req.SignatureRequest = sig

	for _, p := range pool {
		serverAddr := fmt.Sprintf("%s:%d", p.IP().String(), p.Port())
		conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
		if err != nil {
			grpcErr, _ := status.FromError(err)
			fmt.Printf("STORE TRANSACTION REQUEST - ERROR: %s\n", grpcErr.Message())
			return
		}
		defer conn.Close()
		cli := api.NewTransactionServiceClient(conn)
		res, err := cli.StoreTransaction(context.Background(), req)
		if err != nil {
			grpcErr, _ := status.FromError(err)
			fmt.Printf("STORE TRANSACTION RESPONSE - ERROR: %s\n", grpcErr.Message())
			return
		}

		fmt.Printf("STORE TRANSACTION RESPONSE - %s\n", time.Unix(res.Timestamp, 0).String())

		resBytes, err := json.Marshal(&api.StoreResponse{
			Timestamp: res.Timestamp,
		})
		if err != nil {
			fmt.Printf("STORE TRANSACTION RESPONSE - ERROR: %s\n", err.Error())
			return
		}
		if err := crypto.VerifySignature(string(resBytes), lastMinerKeys.PublicKey(), res.SignatureResponse); err != nil {
			fmt.Printf("STORE TRANSACTION RESPONSE - ERROR: %s\n", err.Error())
			return
		}

		ackChan <- true
	}
}

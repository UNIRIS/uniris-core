package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"google.golang.org/grpc/codes"

	"github.com/uniris/uniris-core/pkg/shared"

	"google.golang.org/grpc/status"

	api "github.com/uniris/uniris-core/api/protobuf-spec"
	"github.com/uniris/uniris-core/pkg/chain"
	"github.com/uniris/uniris-core/pkg/consensus"
	"github.com/uniris/uniris-core/pkg/crypto"
	"google.golang.org/grpc"
)

type poolRequester struct {
	sharedKeyReader shared.KeyReader
}

//NewPoolRequester creates a new pool requester as GRPC client
func NewPoolRequester(skr shared.KeyReader) consensus.PoolRequester {
	return poolRequester{
		sharedKeyReader: skr,
	}
}

func (pr poolRequester) RequestLastTransaction(pool consensus.Pool, txAddr crypto.VersionnedHash, txType chain.TransactionType) (*chain.Transaction, error) {

	nodeLastKeys, err := pr.sharedKeyReader.LastNodeCrossKeypair()
	if err != nil {
		return nil, err
	}

	ackChan := make(chan bool)
	var nbReqAcks int32
	var nbReqFailures int32
	var failed bool
	minReplies := 1 //TODO: specify minium required transaction replies

	req := &api.GetLastTransactionRequest{
		TransactionAddress: txAddr,
		Type:               api.TransactionType(txType),
		Timestamp:          time.Now().Unix(),
	}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	sig, err := nodeLastKeys.PrivateKey().Sign(reqBytes)
	if err != nil {
		return nil, err
	}
	req.SignatureRequest = sig

	txRes := make([]chain.Transaction, 0)

	for _, p := range pool {
		go func(p consensus.Node) {

			serverAddr := fmt.Sprintf("%s:%d", p.IP(), p.Port())
			conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
			if err != nil {
				fmt.Printf("GET LAST TRANSACTION - ERROR: %s\n", err.Error())
				ackChan <- false
				return
			}
			defer conn.Close()
			cli := api.NewTransactionServiceClient(conn)
			res, err := cli.GetLastTransaction(context.Background(), req)
			if err != nil {
				if status.Code(err) == codes.NotFound {
					ackChan <- true
					return
				}
				grpcErr, _ := status.FromError(err)
				fmt.Printf("GET LAST TRANSACTION RESPONSE - ERROR: %s\n", grpcErr.Message())
				ackChan <- false
				return
			}

			fmt.Printf("GET LAST TRANSACTION RESPONSE - %s\n", time.Unix(res.Timestamp, 0).String())

			resBytes, err := json.Marshal(&api.GetLastTransactionResponse{
				Timestamp:   res.Timestamp,
				Transaction: res.Transaction,
			})
			if err != nil {
				fmt.Printf("GET LAST TRANSACTION RESPONSE - ERROR: %s\n", err.Error())
				ackChan <- false
				return
			}

			if !nodeLastKeys.PublicKey().Verify(resBytes, res.SignatureResponse) {
				fmt.Println("GET LAST TRANSACTION RESPONSE - ERROR: invalid signature")
				ackChan <- false
				return
			}

			if res.Transaction != nil {
				tx, err := formatTransaction(res.Transaction)
				if err != nil {
					fmt.Printf("GET LAST TRANSACTION RESPONSE - ERROR: %s\n", err.Error())
					ackChan <- false
					return
				}
				txRes = append(txRes, tx)
			}
			ackChan <- true
		}(p)
	}

	for ack := range ackChan {
		if ack {
			atomic.AddInt32(&nbReqAcks, 1)
		} else {
			atomic.AddInt32(&nbReqFailures, 1)
		}
		if atomic.LoadInt32(&nbReqFailures) == int32(len(pool)) {
			failed = true
			break
		}
		if atomic.LoadInt32(&nbReqAcks) == int32(minReplies) {
			break
		}
	}

	if failed {
		return nil, errors.New("transaction request failed")
	}

	if len(txRes) == 0 {
		return nil, nil
	}

	//TODO: consensus to implement to get the right result
	return &txRes[0], nil
}

func (pr poolRequester) RequestTransactionTimeLock(pool consensus.Pool, txHash crypto.VersionnedHash, txAddress crypto.VersionnedHash, masterPublicKey crypto.PublicKey) error {

	lastKeys, err := pr.sharedKeyReader.LastNodeCrossKeypair()
	if err != nil {
		return err
	}

	ackChan := make(chan bool)
	var nbLockAcks int32
	var nbLockFailures int32
	var failed bool

	minLocks := 1 //TODO: specify minium required timelocks

	masterKey, err := masterPublicKey.Marshal()
	if err != nil {
		return err
	}
	req := &api.TimeLockTransactionRequest{
		Address:             txAddress,
		TransactionHash:     txHash,
		MasterNodePublicKey: masterKey,
		Timestamp:           time.Now().Unix(),
	}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return err
	}
	sig, err := lastKeys.PrivateKey().Sign(reqBytes)
	if err != nil {
		return err
	}
	req.SignatureRequest = sig

	for _, p := range pool {
		go func(p consensus.Node) {

			serverAddr := fmt.Sprintf("%s:%d", p.IP(), p.Port())
			conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
			if err != nil {
				grpcErr, _ := status.FromError(err)
				fmt.Printf("TIMELOCK TRANSACTION REQUEST - ERROR: %s\n", grpcErr.Message())
				ackChan <- false
				return
			}
			defer conn.Close()
			cli := api.NewTransactionServiceClient(conn)
			res, err := cli.TimeLockTransaction(context.Background(), req)
			if err != nil {
				grpcErr, _ := status.FromError(err)
				fmt.Printf("TIMELOCK TRANSACTION RESPONSE - ERROR: %s\n", grpcErr.Message())
				ackChan <- false
				return
			}

			fmt.Printf("TIMELOCK TRANSACTION RESPONSE - %s\n", time.Unix(res.Timestamp, 0).String())

			resBytes, err := json.Marshal(&api.TimeLockTransactionResponse{
				Timestamp: req.Timestamp,
			})
			if err != nil {
				fmt.Printf("TIMELOCK TRANSACTION RESPONSE - ERROR: %s", err.Error())
				ackChan <- false
				return
			}
			if !lastKeys.PublicKey().Verify(resBytes, res.SignatureResponse) {
				fmt.Println("LOCK TRANSACTION RESPONSE - ERROR: invalid signature")
				ackChan <- false
				return
			}

			ackChan <- true
		}(p)
	}

	for ack := range ackChan {
		if ack {
			atomic.AddInt32(&nbLockAcks, 1)
		} else {
			atomic.AddInt32(&nbLockFailures, 1)
		}
		if atomic.LoadInt32(&nbLockFailures) == int32(len(pool)) {
			failed = true
			break
		}
		if atomic.LoadInt32(&nbLockAcks) == int32(minLocks) {
			break
		}
	}

	if failed {
		return errors.New("transaction timelock failed")
	}

	return nil
}

func (pr poolRequester) RequestTransactionValidations(pool consensus.Pool, tx chain.Transaction, minValids int, masterValid chain.MasterValidation) ([]chain.Validation, error) {

	lastKeys, err := pr.sharedKeyReader.LastNodeCrossKeypair()
	if err != nil {
		return nil, fmt.Errorf("confirm validation request error: %s", err.Error())
	}

	txf, err := formatAPITransaction(tx)
	if err != nil {
		return nil, err
	}
	mvf, err := formatAPIMasterValidation(masterValid)
	if err != nil {
		return nil, err
	}

	req := &api.ConfirmTransactionValidationRequest{
		MasterValidation: mvf,
		Transaction:      txf,
		Timestamp:        time.Now().Unix(),
	}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("confirm validation request error: %s", err.Error())
	}
	sig, err := lastKeys.PrivateKey().Sign(reqBytes)
	if err != nil {
		return nil, fmt.Errorf("confirm validation request error: %s", err.Error())
	}
	req.SignatureRequest = sig

	validations := make([]chain.Validation, 0)

	ackChan := make(chan bool)
	var nbValidationAck int32
	var nbValidationFailures int32
	var failed bool

	for _, p := range pool {

		go func(p consensus.Node) {

			serverAddr := fmt.Sprintf("%s:%d", p.IP().String(), p.Port())
			conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
			if err != nil {
				grpcErr, _ := status.FromError(err)
				fmt.Printf("CONFIRM VALIDATION TRANSACTION REQUEST - ERROR: %s\n", grpcErr.Message())
				ackChan <- false
				return
			}
			defer conn.Close()
			cli := api.NewTransactionServiceClient(conn)
			res, err := cli.ConfirmTransactionValidation(context.Background(), req)
			if err != nil {
				grpcErr, _ := status.FromError(err)
				fmt.Printf("CONFIRM VALIDATION TRANSACTION RESPONSE - ERROR: %s\n", grpcErr.Message())
				ackChan <- false
				return
			}

			fmt.Printf("CONFIRM VALIDATION TRANSACTION RESPONSE - %s\n", time.Unix(res.Timestamp, 0).String())

			resBytes, err := json.Marshal(&api.ConfirmTransactionValidationResponse{
				Timestamp:  res.Timestamp,
				Validation: res.Validation,
			})
			if err != nil {
				fmt.Printf("CONFIRM VALIDATION TRANSACTION RESPONSE - ERROR: %s\n", err.Error())
				ackChan <- false
				return
			}
			if !lastKeys.PublicKey().Verify(resBytes, res.SignatureResponse) {
				fmt.Println("CONFIRM VALIDATION TRANSACTION RESPONSE - ERROR: invalid signature")
				ackChan <- false
				return
			}

			vKey, err := crypto.ParsePublicKey(res.Validation.PublicKey)
			if err != nil {
				fmt.Println("CONFIRM VALIDATION TRANSACTION RESPONSE - ERROR: invalid validator public key")
				ackChan <- false
				return
			}
			v, err := chain.NewValidation(chain.ValidationStatus(res.Validation.Status), time.Unix(res.Timestamp, 0), vKey, res.Validation.Signature)
			if err != nil {
				ackChan <- false
				return
			}
			validations = append(validations, v)
			ackChan <- true
		}(p)
	}

	for ack := range ackChan {
		if ack {
			atomic.AddInt32(&nbValidationAck, 1)
		} else {
			atomic.AddInt32(&nbValidationFailures, 1)
		}
		if atomic.LoadInt32(&nbValidationFailures) == int32(len(pool)) {
			failed = true
			break
		}
		if atomic.LoadInt32(&nbValidationAck) == int32(minValids) {
			break
		}
	}

	if failed {
		return nil, errors.New("transaction storage failed")
	}

	return validations, nil
}

func (pr poolRequester) RequestTransactionStorage(pool consensus.Pool, minStorage int, tx chain.Transaction) error {

	lastKeys, err := pr.sharedKeyReader.LastNodeCrossKeypair()
	if err != nil {
		return fmt.Errorf("store transaction request error: %s", err.Error())
	}

	confValids := make([]*api.Validation, 0)
	for _, v := range tx.ConfirmationsValidations() {
		vf, err := formatAPIValidation(v)
		if err != nil {
			return err
		}
		confValids = append(confValids, vf)
	}

	txf, err := formatAPITransaction(tx)
	if err != nil {
		return err
	}

	mvf, err := formatAPIMasterValidation(tx.MasterValidation())
	if err != nil {
		return err
	}

	req := &api.StoreTransactionRequest{
		MinedTransaction: &api.MinedTransaction{
			MasterValidation:   mvf,
			ConfirmValidations: confValids,
			Transaction:        txf,
		},
		Timestamp: time.Now().Unix(),
	}

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("store transaction request error: %s", err.Error())
	}
	sig, err := lastKeys.PrivateKey().Sign(reqBytes)
	if err != nil {
		return fmt.Errorf("store transaction request error: %s", err.Error())
	}

	req.SignatureRequest = sig

	ackChan := make(chan bool)
	var nbStorageAck int32
	var nbStorageFailures int32
	var failed bool

	for _, p := range pool {
		go func(p consensus.Node) {

			serverAddr := fmt.Sprintf("%s:%d", p.IP().String(), p.Port())
			conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
			if err != nil {
				grpcErr, _ := status.FromError(err)
				fmt.Printf("STORE TRANSACTION REQUEST - ERROR: %s\n", grpcErr.Message())
				ackChan <- false
				return
			}
			defer conn.Close()
			cli := api.NewTransactionServiceClient(conn)
			res, err := cli.StoreTransaction(context.Background(), req)
			if err != nil {
				grpcErr, _ := status.FromError(err)
				fmt.Printf("STORE TRANSACTION RESPONSE - ERROR: %s\n", grpcErr.Message())
				ackChan <- false
				return
			}

			fmt.Printf("STORE TRANSACTION RESPONSE - %s\n", time.Unix(res.Timestamp, 0).String())

			resBytes, err := json.Marshal(&api.StoreTransactionResponse{
				Timestamp: res.Timestamp,
			})
			if err != nil {
				fmt.Printf("STORE TRANSACTION RESPONSE - ERROR: %s\n", err.Error())
				ackChan <- false
				return
			}
			if !lastKeys.PublicKey().Verify(resBytes, res.SignatureResponse) {
				fmt.Println("STORE TRANSACTION RESPONSE - ERROR: invalid signature")
				ackChan <- false
				return
			}

			ackChan <- true
		}(p)
	}

	for ack := range ackChan {
		if ack {
			atomic.AddInt32(&nbStorageAck, 1)
		} else {
			atomic.AddInt32(&nbStorageFailures, 1)
		}
		if atomic.LoadInt32(&nbStorageFailures) == int32(len(pool)) {
			failed = true
			break
		}
		if atomic.LoadInt32(&nbStorageAck) == int32(minStorage) {
			break
		}
	}

	if failed {
		return errors.New("transaction storage failed")
	}

	return nil
}

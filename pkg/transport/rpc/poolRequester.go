package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/uniris/uniris-core/pkg/shared"

	"google.golang.org/grpc/status"

	api "github.com/uniris/uniris-core/api/protobuf-spec"
	"github.com/uniris/uniris-core/pkg/chain"
	"github.com/uniris/uniris-core/pkg/consensus"
	"github.com/uniris/uniris-core/pkg/crypto"
	"google.golang.org/grpc"
)

type poolRequester struct {
	techDB shared.TechDatabaseReader
}

//NewPoolRequester creates a new pool requester as GRPC client
func NewPoolRequester(tDB shared.TechDatabaseReader) consensus.PoolRequester {
	return poolRequester{
		techDB: tDB,
	}
}

func (pr poolRequester) RequestLastTransaction(pool consensus.Pool, txAddr string, txType chain.TransactionType) (*chain.Transaction, error) {

	nodeLastKeys, err := pr.techDB.NodeLastKeys()
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	wg.Add(len(pool))

	req := &api.LastTransactionRequest{
		TransactionAddress: txAddr,
		Type:               api.TransactionType(txType),
		Timestamp:          time.Now().Unix(),
	}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	sig, err := crypto.Sign(string(reqBytes), nodeLastKeys.PrivateKey())
	if err != nil {
		return nil, err
	}
	req.SignatureRequest = sig

	txRes := make([]chain.Transaction, 0)

	for _, p := range pool {
		go func(p consensus.PoolMember) {
			defer wg.Done()

			serverAddr := fmt.Sprintf("%s:%d", p.IP(), p.Port())
			conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
			if err != nil {
				fmt.Printf("GET LAST TRANSACTION - ERROR: %s\n", err.Error())
				return
			}
			defer conn.Close()
			cli := api.NewStorageServiceClient(conn)
			res, err := cli.GetLastTransaction(context.Background(), req)
			if err != nil {
				grpcErr, _ := status.FromError(err)
				fmt.Printf("GET LAST TRANSACTION RESPONSE - ERROR: %s\n", grpcErr.Message())
				return
			}

			fmt.Printf("GET LAST TRANSACTION RESPONSE - %s\n", time.Unix(res.Timestamp, 0).String())

			resBytes, err := json.Marshal(&api.LastTransactionResponse{
				Timestamp:   res.Timestamp,
				Transaction: res.Transaction,
			})
			if err != nil {
				fmt.Printf("GET LAST TRANSACTION RESPONSE - ERROR: %s\n", err.Error())
				return
			}
			if err := crypto.VerifySignature(string(resBytes), nodeLastKeys.PublicKey(), res.SignatureResponse); err != nil {
				fmt.Printf("GET LAST TRANSACTION RESPONSE - ERROR: %s\n", err.Error())
				return
			}

			if res.Transaction != nil {
				tx, err := formatTransaction(res.Transaction)
				if err != nil {
					fmt.Printf("GET LAST TRANSACTION RESPONSE - ERROR: %s\n", err.Error())
					return
				}
				txRes = append(txRes, tx)
			}
		}(p)
	}

	wg.Wait()

	if len(txRes) == 0 {
		return nil, nil
	}

	//TODO: consensus to implement to get the right result
	return &txRes[0], nil
}

func (pr poolRequester) RequestTransactionLock(pool consensus.Pool, txHash string, txAddress string, masterPublicKey string) error {

	lastKeys, err := pr.techDB.NodeLastKeys()
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(len(pool))

	var ackUnlock int32

	req := &api.LockRequest{
		Address:             txAddress,
		TransactionHash:     txHash,
		MasterNodePublicKey: masterPublicKey,
		Timestamp:           time.Now().Unix(),
	}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return err
	}
	sig, err := crypto.Sign(string(reqBytes), lastKeys.PrivateKey())
	if err != nil {
		return err
	}
	req.SignatureRequest = sig

	for _, p := range pool {
		go func(p consensus.PoolMember) {
			defer wg.Done()

			serverAddr := fmt.Sprintf("%s:%d", p.IP(), p.Port())
			conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
			if err != nil {
				grpcErr, _ := status.FromError(err)
				fmt.Printf("LOCK TRANSACTION REQUEST - ERROR: %s\n", grpcErr.Message())
				return
			}
			defer conn.Close()
			cli := api.NewStorageServiceClient(conn)
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
			if err := crypto.VerifySignature(string(resBytes), lastKeys.PublicKey(), res.SignatureResponse); err != nil {
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

func (pr poolRequester) RequestTransactionValidations(pool consensus.Pool, tx chain.Transaction, minValids int, masterValid chain.MasterValidation) ([]chain.Validation, error) {

	lastKeys, err := pr.techDB.NodeLastKeys()
	if err != nil {
		return nil, fmt.Errorf("confirm validation request error: %s", err.Error())
	}

	req := &api.ConfirmValidationRequest{
		MasterValidation: formatAPIMasterValidation(masterValid),
		Transaction:      formatAPITransaction(tx),
		Timestamp:        time.Now().Unix(),
	}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("confirm validation request error: %s", err.Error())
	}
	sig, err := crypto.Sign(string(reqBytes), lastKeys.PrivateKey())
	if err != nil {
		return nil, fmt.Errorf("confirm validation request error: %s", err.Error())
	}
	req.SignatureRequest = sig

	var wg sync.WaitGroup
	wg.Add(minValids)

	validations := make([]chain.Validation, 0)

	for _, p := range pool {

		go func(p consensus.PoolMember) {
			defer wg.Done()

			serverAddr := fmt.Sprintf("%s:%d", p.IP().String(), p.Port())
			conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
			if err != nil {
				grpcErr, _ := status.FromError(err)
				fmt.Printf("CONFIRM VALIDATION TRANSACTION REQUEST - ERROR: %s\n", grpcErr.Message())
				return
			}
			defer conn.Close()
			cli := api.NewMiningServiceClient(conn)
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
			if err := crypto.VerifySignature(string(resBytes), lastKeys.PublicKey(), res.SignatureResponse); err != nil {
				fmt.Printf("CONFIRM VALIDATION TRANSACTION RESPONSE - ERROR: %s\n", err.Error())
				return
			}

			v, err := chain.NewValidation(chain.ValidationStatus(res.Validation.Status), time.Unix(res.Timestamp, 0), res.Validation.PublicKey, res.Validation.Signature)
			if err != nil {
				return
			}
			validations = append(validations, v)
		}(p)
	}

	wg.Wait()

	return validations, nil
}

func (pr poolRequester) RequestTransactionStorage(pool consensus.Pool, minStorage int, tx chain.Transaction) error {

	lastKeys, err := pr.techDB.NodeLastKeys()
	if err != nil {
		return fmt.Errorf("store transaction request error: %s", err.Error())
	}

	confValids := make([]*api.Validation, 0)
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
		return fmt.Errorf("store transaction request error: %s", err.Error())
	}
	sig, err := crypto.Sign(string(reqBytes), lastKeys.PrivateKey())
	if err != nil {
		return fmt.Errorf("store transaction request error: %s", err.Error())
	}

	req.SignatureRequest = sig

	var wg sync.WaitGroup
	wg.Add(minStorage)

	for _, p := range pool {
		go func(p consensus.PoolMember) {

			defer wg.Done()

			serverAddr := fmt.Sprintf("%s:%d", p.IP().String(), p.Port())
			conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
			if err != nil {
				grpcErr, _ := status.FromError(err)
				fmt.Printf("STORE TRANSACTION REQUEST - ERROR: %s\n", grpcErr.Message())
				return
			}
			defer conn.Close()
			cli := api.NewStorageServiceClient(conn)
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
			if err := crypto.VerifySignature(string(resBytes), lastKeys.PublicKey(), res.SignatureResponse); err != nil {
				fmt.Printf("STORE TRANSACTION RESPONSE - ERROR: %s\n", err.Error())
				return
			}

		}(p)
	}

	wg.Wait()

	return nil
}

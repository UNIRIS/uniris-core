package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	api "github.com/uniris/uniris-core/api/protobuf-spec"
	"github.com/uniris/uniris-core/pkg/crypto"
	"github.com/uniris/uniris-core/pkg/transaction"
	"google.golang.org/grpc"
)

type poolRtrv struct {
	sharedPubk string
	sharedPvk  string
}

//NewPoolRetriever creates a new pool retriever as a GRPC client
func NewPoolRetriever(sharedPubk string, sharedPvk string) transaction.PoolRetriever {
	return poolRtrv{
		sharedPubk: sharedPubk,
		sharedPvk:  sharedPvk,
	}
}

func (pr poolRtrv) RequestLastTransaction(pool transaction.Pool, txAddr string, txType transaction.Type) (*transaction.Transaction, error) {

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
	sig, err := crypto.Sign(string(reqBytes), pr.sharedPubk)
	if err != nil {
		return nil, err
	}
	req.SignatureRequest = sig

	txRes := make([]transaction.Transaction, 0)

	for _, p := range pool {
		go func(p transaction.PoolMember) {
			defer wg.Done()

			serverAddr := fmt.Sprintf("%s:%d", p.IP(), p.Port())
			conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
			if err != nil {
				fmt.Printf("GET LAST TRANSACTION - ERROR: %s", err.Error())
			}
			defer conn.Close()
			cli := api.NewTransactionServiceClient(conn)
			res, err := cli.GetLastTransaction(context.Background(), req)
			if err != nil {
				fmt.Printf("GET LAST TRANSACTION RESPONSE - ERROR: %s", err.Error())
				return
			}

			fmt.Printf("GET LAST TRANSACTION RESPONSE - %s", time.Unix(res.Timestamp, 0).String())

			resBytes, err := json.Marshal(res)
			if err != nil {
				fmt.Printf("GET LAST TRANSACTION RESPONSE - ERROR: %s", err.Error())
				return
			}
			if err := crypto.VerifySignature(string(resBytes), pr.sharedPubk, res.SignatureResponse); err != nil {
				fmt.Printf("GET LAST TRANSACTION RESPONSE - ERROR: %s", err.Error())
				return
			}

			tx, err := formatTransaction(res.Transaction)
			if err != nil {
				fmt.Printf("GET LAST TRANSACTION RESPONSE - ERROR: %s", err.Error())
				return
			}
			txRes = append(txRes, tx)
		}(p)
	}

	wg.Wait()

	if len(txRes) == 0 {
		return nil, nil
	}

	//TODO: consensus to implement to get the right result
	return &txRes[0], nil
}

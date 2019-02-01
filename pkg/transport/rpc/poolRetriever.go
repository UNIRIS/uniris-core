package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/uniris/uniris-core/pkg/shared"

	api "github.com/uniris/uniris-core/api/protobuf-spec"
	"github.com/uniris/uniris-core/pkg/crypto"
	"github.com/uniris/uniris-core/pkg/transaction"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type poolRtrv struct {
	sharedSrv shared.Service
}

//NewPoolRetriever creates a new pool retriever as a GRPC client
func NewPoolRetriever(sharedSrv shared.Service) transaction.PoolRetriever {
	return poolRtrv{
		sharedSrv: sharedSrv,
	}
}

func (pr poolRtrv) RequestLastTransaction(pool transaction.Pool, txAddr string, txType transaction.Type) (*transaction.Transaction, error) {

	lastMinersKeys, err := pr.sharedSrv.GetSharedMinerKeys()
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
	sig, err := crypto.Sign(string(reqBytes), lastMinersKeys.PrivateKey())
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
				fmt.Printf("GET LAST TRANSACTION - ERROR: %s\n", err.Error())
				return
			}
			defer conn.Close()
			cli := api.NewTransactionServiceClient(conn)
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
			if err := crypto.VerifySignature(string(resBytes), lastMinersKeys.PublicKey(), res.SignatureResponse); err != nil {
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

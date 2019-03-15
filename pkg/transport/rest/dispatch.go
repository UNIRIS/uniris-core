package rest

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/uniris/uniris-core/pkg/chain"

	api "github.com/uniris/uniris-core/api/protobuf-spec"
	"github.com/uniris/uniris-core/pkg/consensus"
	"github.com/uniris/uniris-core/pkg/crypto"
	"github.com/uniris/uniris-core/pkg/shared"
	"google.golang.org/grpc"
)

func requestTransactionMining(tx *api.Transaction, pvKey crypto.PrivateKey, pubKey crypto.PublicKey, nodeReader consensus.NodeReader, sharedKeyReader shared.KeyReader) (transactionResponse, *httpError) {

	masterNodes, err := consensus.FindMasterNodes(tx.TransactionHash, nodeReader, sharedKeyReader)
	if err != nil {
		return transactionResponse{}, &httpError{
			code:      http.StatusInternalServerError,
			Error:     err.Error(),
			Timestamp: time.Now().Unix(),
			Status:    http.StatusText(http.StatusInternalServerError),
		}
	}

	//Get the minimum number of validations based on the type of the transaction and its fees
	txFees := consensus.TransactionFees(chain.TransactionType(tx.Type), tx.Data)
	minValidations, err := consensus.RequiredValidationNumber(chain.TransactionType(tx.Type), txFees, nodeReader, sharedKeyReader)
	if err != nil {
		return transactionResponse{}, &httpError{
			code:      http.StatusInternalServerError,
			Error:     err.Error(),
			Timestamp: time.Now().Unix(),
			Status:    http.StatusText(http.StatusInternalServerError),
		}
	}

	//Building the welcome node headers
	wHeaders := make([]*api.NodeHeader, 0)
	for _, n := range masterNodes.Nodes() {

		pubKey, err := n.PublicKey().Marshal()
		if err != nil {
			return transactionResponse{}, &httpError{
				code:      http.StatusInternalServerError,
				Error:     err.Error(),
				Timestamp: time.Now().Unix(),
				Status:    http.StatusText(http.StatusInternalServerError),
			}
		}

		wHeaders = append(wHeaders, &api.NodeHeader{
			IsMaster:      true,
			IsUnreachable: !n.IsReachable(),
			PublicKey:     pubKey,
			PatchNumber:   int32(n.Patch().ID()),
			IsOK:          n.Status() == consensus.NodeOK,
		})
	}

	req := &api.LeadTransactionMiningRequest{
		Transaction:        tx,
		MinimumValidations: int32(minValidations),
		Timestamp:          time.Now().Unix(),
		WelcomeHeaders:     wHeaders,
	}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return transactionResponse{}, &httpError{
			code:      http.StatusInternalServerError,
			Error:     err.Error(),
			Timestamp: time.Now().Unix(),
			Status:    http.StatusText(http.StatusInternalServerError),
		}
	}
	reqSig, err := pvKey.Sign(reqBytes)
	if err != nil {
		return transactionResponse{}, &httpError{
			code:      http.StatusInternalServerError,
			Error:     err.Error(),
			Timestamp: time.Now().Unix(),
			Status:    http.StatusText(http.StatusInternalServerError),
		}
	}
	req.SignatureRequest = reqSig

	minAckMaster := 1 //TODO: define how many master must ack before to send response to the client

	var nbMasterAck int32
	var nbMasterFailures int32
	var failed bool
	ackChan := make(chan bool)

	//Send concurrently mining request to several masters
	for _, n := range masterNodes.Nodes() {
		go func(n consensus.Node) {

			masterAddr := fmt.Sprintf("%s:%d", n.IP().String(), n.Port())
			conn, err := grpc.Dial(masterAddr, grpc.WithInsecure())
			defer conn.Close()
			if err != nil {
				fmt.Printf("error - master unreachable: %s\n", err.Error())
				ackChan <- false
				return
			}

			cli := api.NewTransactionServiceClient(conn)
			res, err := cli.LeadTransactionMining(context.Background(), req)
			if err != nil {
				_, message := parseGrpcError(err)
				fmt.Printf("error - master dispatch: transaction failed - cause: %s\n", message)
				ackChan <- false
				return
			}

			fmt.Printf("LEAD TRANSACTION MINING RESPONSE - %s\n", time.Unix(res.Timestamp, 0).String())
			resBytes, err := json.Marshal(&api.LeadTransactionMiningResponse{
				Timestamp: res.Timestamp,
			})
			if err != nil {
				fmt.Printf("error - master dispatch: transaction response bad format - cause: %s\n", err.Error())
				ackChan <- false
				return
			}
			if !pubKey.Verify(resBytes, res.SignatureResponse) {
				fmt.Printf("error - master dispatch: invalid signature response\n")
				ackChan <- false
				return
			}

			ackChan <- true

		}(n)
	}

	//Waiting the min ack of master
	for ack := range ackChan {
		if ack {
			atomic.AddInt32(&nbMasterAck, 1)
		} else {
			atomic.AddInt32(&nbMasterFailures, 1)
		}
		if atomic.LoadInt32(&nbMasterFailures) == int32(len(masterNodes.Nodes())) {
			failed = true
			break
		}
		if atomic.LoadInt32(&nbMasterAck) == int32(minAckMaster) {
			break
		}
	}

	//When no master could reply to the transaction mining request
	//Then we send an error response back to the client
	if failed {
		return transactionResponse{}, &httpError{
			code:      http.StatusInternalServerError,
			Error:     "transaction failed", //TODO: provide more details
			Timestamp: time.Now().Unix(),
			Status:    http.StatusText(http.StatusInternalServerError),
		}
	}

	//When the minimum ack from the master has been reached
	//Then we send a success response back to the client
	txRes := transactionResponse{
		Timestamp:          time.Now().Unix(),
		TransactionReceipt: encodeTxReceipt(tx),
	}
	txResBytes, err := json.Marshal(txRes)
	if err != nil {
		return transactionResponse{}, &httpError{
			code:      http.StatusInternalServerError,
			Error:     err.Error(),
			Timestamp: time.Now().Unix(),
			Status:    http.StatusText(http.StatusInternalServerError),
		}
	}
	sig, err := pvKey.Sign(txResBytes)
	if err != nil {
		return transactionResponse{}, &httpError{
			code:      http.StatusInternalServerError,
			Error:     err.Error(),
			Timestamp: time.Now().Unix(),
			Status:    http.StatusText(http.StatusInternalServerError),
		}
	}
	txRes.Signature = hex.EncodeToString(sig)

	return txRes, nil
}

func findLastTransaction(txAddr crypto.VersionnedHash, txType api.TransactionType, pvKey crypto.PrivateKey, nodeReader consensus.NodeReader) (*api.Transaction, *httpError) {
	storagePool, err := consensus.FindStoragePool(txAddr, nodeReader)
	if err != nil {
		return nil, &httpError{
			code:      http.StatusInternalServerError,
			Error:     err.Error(),
			Timestamp: time.Now().Unix(),
			Status:    http.StatusText(http.StatusInternalServerError),
		}
	}
	storageMasterNode := fmt.Sprintf("%s:%d", storagePool.Nodes()[0].IP().String(), storagePool.Nodes()[0].Port())
	conn, err := grpc.Dial(storageMasterNode, grpc.WithInsecure())
	defer conn.Close()
	if err != nil {
		return nil, &httpError{
			code:      http.StatusServiceUnavailable,
			Error:     err.Error(),
			Timestamp: time.Now().Unix(),
			Status:    http.StatusText(http.StatusServiceUnavailable),
		}
	}

	req := &api.GetLastTransactionRequest{
		TransactionAddress: txAddr,
		Type:               txType,
		Timestamp:          time.Now().Unix(),
	}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, &httpError{
			code:      http.StatusInternalServerError,
			Error:     err.Error(),
			Timestamp: time.Now().Unix(),
			Status:    http.StatusText(http.StatusInternalServerError),
		}
	}

	sig, err := pvKey.Sign(reqBytes)
	if err != nil {
		return nil, &httpError{
			code:      http.StatusInternalServerError,
			Error:     err.Error(),
			Timestamp: time.Now().Unix(),
			Status:    http.StatusText(http.StatusInternalServerError),
		}
	}
	req.SignatureRequest = sig

	cli := api.NewTransactionServiceClient(conn)
	res, err := cli.GetLastTransaction(context.Background(), req)
	if err != nil {
		code, message := parseGrpcError(err)
		return nil, &httpError{
			code:      code,
			Error:     message,
			Timestamp: time.Now().Unix(),
			Status:    http.StatusText(code),
		}
	}

	return res.Transaction, nil
}

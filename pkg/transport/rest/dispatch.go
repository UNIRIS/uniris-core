package rest

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	api "github.com/uniris/uniris-core/api/protobuf-spec"
	"github.com/uniris/uniris-core/pkg/chain"
	"github.com/uniris/uniris-core/pkg/consensus"
	"github.com/uniris/uniris-core/pkg/crypto"
	"google.golang.org/grpc"
)

func requestTransactionMining(tx *api.Transaction, pvKey crypto.PrivateKey, pubKey crypto.PublicKey) (transactionResponse, *httpError) {

	masterNodes, err := consensus.FindMasterNodes(tx.TransactionHash, chain.TransactionType(tx.Type))
	if err != nil {
		return transactionResponse{}, &httpError{
			code:      http.StatusInternalServerError,
			Error:     err.Error(),
			Timestamp: time.Now().Unix(),
			Status:    http.StatusText(http.StatusInternalServerError),
		}
	}

	minValidations := consensus.GetMinimumValidation(tx.TransactionHash)

	wHeaders := make([]*api.NodeHeader, 0)
	for _, n := range masterNodes {

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
			IsUnreachable: true, //TODO: ensures it
			PublicKey:     pubKey,
			PatchNumber:   0,    //TODO:
			IsOK:          true, //TODO:
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

	//TODO: handle multiple master node sending
	masterAddr := fmt.Sprintf("%s:%d", masterNodes[0].IP().String(), masterNodes[0].Port())
	conn, err := grpc.Dial(masterAddr, grpc.WithInsecure())
	defer conn.Close()
	if err != nil {
		return transactionResponse{}, &httpError{
			code:      http.StatusServiceUnavailable,
			Error:     err.Error(),
			Timestamp: time.Now().Unix(),
			Status:    http.StatusText(http.StatusServiceUnavailable),
		}
	}

	cli := api.NewTransactionServiceClient(conn)
	res, err := cli.LeadTransactionMining(context.Background(), req)
	if err != nil {
		code, message := parseGrpcError(err)
		return transactionResponse{}, &httpError{
			code:      code,
			Error:     message,
			Timestamp: time.Now().Unix(),
			Status:    http.StatusText(code),
		}
	}

	fmt.Printf("LEAD TRANSACTION MINING RESPONSE - %s\n", time.Unix(res.Timestamp, 0).String())

	resBytes, err := json.Marshal(&api.LeadTransactionMiningResponse{
		Timestamp: res.Timestamp,
	})
	if err != nil {
		return transactionResponse{}, &httpError{
			code:      http.StatusInternalServerError,
			Error:     err.Error(),
			Timestamp: time.Now().Unix(),
			Status:    http.StatusText(http.StatusInternalServerError),
		}
	}
	if !pubKey.Verify(resBytes, res.SignatureResponse) {
		return transactionResponse{}, &httpError{
			code:      http.StatusInternalServerError,
			Error:     "invalid signature",
			Timestamp: time.Now().Unix(),
			Status:    http.StatusText(http.StatusInternalServerError),
		}
	}

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

func findLastTransaction(txAddr []byte, txType api.TransactionType, pvKey crypto.PrivateKey) (*api.Transaction, *httpError) {
	storagePool, err := consensus.FindStoragePool(txAddr)
	if err != nil {
		return nil, &httpError{
			code:      http.StatusInternalServerError,
			Error:     err.Error(),
			Timestamp: time.Now().Unix(),
			Status:    http.StatusText(http.StatusInternalServerError),
		}
	}
	storageMasterNode := fmt.Sprintf("%s:%d", storagePool[0].IP().String(), storagePool[0].Port())
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

package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	api "github.com/uniris/uniris-core/api/protobuf-spec"
	"github.com/uniris/uniris-core/pkg/consensus"
	"github.com/uniris/uniris-core/pkg/crypto"
	"github.com/uniris/uniris-core/pkg/shared"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

//GetTransactionStatusHandler defines an HTTP handler to get the status of a transaction
func GetTransactionStatusHandler(techReader shared.TechDatabaseReader) func(c *gin.Context) {
	return func(c *gin.Context) {
		txReceipt := c.Param("txReceipt")
		txAddress, txHash, err := decodeTxReceipt(txReceipt)
		if err != nil {
			c.JSON(http.StatusBadRequest, httpError{
				Error:     fmt.Sprintf("tx receipt decoding: %s", err.Error()),
				Status:    http.StatusText(http.StatusBadRequest),
				Timestamp: time.Now().Unix(),
			})
			return
		}

		sPool, err := consensus.FindStoragePool(txAddress)
		if err != nil {
			c.JSON(http.StatusInternalServerError, httpError{
				Error:     err.Error(),
				Status:    http.StatusText(http.StatusInternalServerError),
				Timestamp: time.Now().Unix(),
			})
			return
		}

		//Send request to the storage master node
		serverAddr := fmt.Sprintf("%s:%d", sPool[0].IP().String(), sPool[0].Port())
		conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, httpError{
				Error:     err.Error(),
				Status:    http.StatusText(http.StatusServiceUnavailable),
				Timestamp: time.Now().Unix(),
			})
			return
		}
		defer conn.Close()

		cli := api.NewTransactionServiceClient(conn)
		reqStatus := &api.GetTransactionStatusRequest{
			TransactionHash: txHash,
			Timestamp:       time.Now().Unix(),
		}
		reqBytes, err := json.Marshal(reqStatus)
		if err != nil {
			c.JSON(http.StatusInternalServerError, httpError{
				Error:     err.Error(),
				Status:    http.StatusText(http.StatusInternalServerError),
				Timestamp: time.Now().Unix(),
			})
			return
		}

		nodeLastKeys, err := techReader.NodeLastKeys()
		if err != nil {
			c.JSON(http.StatusInternalServerError, httpError{
				Error:     err.Error(),
				Status:    http.StatusText(http.StatusInternalServerError),
				Timestamp: time.Now().Unix(),
			})
			return
		}
		sig, err := crypto.Sign(string(reqBytes), nodeLastKeys.PrivateKey())
		if err != nil {
			c.JSON(http.StatusInternalServerError, httpError{
				Error:     err.Error(),
				Status:    http.StatusText(http.StatusInternalServerError),
				Timestamp: time.Now().Unix(),
			})
			return
		}
		reqStatus.SignatureRequest = sig

		res, err := cli.GetTransactionStatus(context.Background(), reqStatus)
		if err != nil {
			grpcStatus, _ := status.FromError(err)
			code, message := parseGrpcError(grpcStatus.Err())
			c.JSON(code, httpError{
				Error:     message,
				Status:    http.StatusText(code),
				Timestamp: time.Now().Unix(),
			})
			return
		}

		fmt.Printf("GET TRANSACTION STATUS RESPONSE - %s\n", time.Unix(res.Timestamp, 0).String())
		resBytes, err := json.Marshal(&api.GetTransactionStatusResponse{
			Status:    res.Status,
			Timestamp: res.Timestamp,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, httpError{
				Error:     err.Error(),
				Status:    http.StatusText(http.StatusInternalServerError),
				Timestamp: time.Now().Unix(),
			})
			return
		}

		if err := crypto.VerifySignature(string(resBytes), nodeLastKeys.PublicKey(), res.SignatureResponse); err != nil {
			c.JSON(http.StatusInternalServerError, httpError{
				Error:     err.Error(),
				Status:    http.StatusText(http.StatusInternalServerError),
				Timestamp: time.Now().Unix(),
			})
			return
		}

		c.JSON(http.StatusOK, transactionStatusResponse{
			Status:    res.Status.String(),
			Timestamp: res.Timestamp,
			Signature: res.SignatureResponse,
		})
	}
}

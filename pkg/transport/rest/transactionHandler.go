package rest

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	api "github.com/uniris/uniris-core/api/protobuf-spec"
	"github.com/uniris/uniris-core/pkg/crypto"
	"google.golang.org/grpc"
)

//NewTransactionHandler creates a new HTTP handler for the transaction endpoints
func NewTransactionHandler(r *gin.RouterGroup, internalPort int) {
	r.GET("/transaction/:txReceipt/status", getTransactionStatus(internalPort))
}

func getTransactionStatus(internalPort int) func(c *gin.Context) {
	return func(c *gin.Context) {
		txReceipt := c.Param("txReceipt")
		txAddress, txHash, err := decodeTxReceipt(txReceipt)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("tx receipt decoding: %s", err.Error())})
			return
		}

		serverAddr := fmt.Sprintf("localhost:%d", internalPort)
		conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
		defer conn.Close()
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			return
		}
		cli := api.NewInternalServiceClient(conn)
		req := &api.InternalTransactionStatusRequest{
			Timestamp:          time.Now().Unix(),
			TransactionHash:    txHash,
			TransactionAddress: txAddress,
		}
		res, err := cli.GetTransactionStatus(context.Background(), req)
		if err != nil {
			c.JSON(parseGrpcError(err))
			return
		}

		c.JSON(http.StatusOK, map[string]interface{}{
			"status":    res.Status.String(),
			"timestamp": res.Timestamp,
			"signature": res.SignatureResponse,
		})
	}
}

func decodeTxReceipt(receipt string) (addr, hash string, err error) {
	if _, err = hex.DecodeString(receipt); err != nil {
		err = errors.New("must be hexadecimal")
		return
	}

	/*
		Length from sha256 hash is 64 bytes.
		a transaction receipt is a set of the hash of the address and the hash of the transaction
		So a transaction receipt is 128 bytes
	*/
	if len(receipt) != 128 {
		err = errors.New("invalid length")
		return
	}

	addr = receipt[:64]
	hash = receipt[64:]

	if _, err = crypto.IsHash(addr); err != nil {
		return
	}

	if _, err = crypto.IsHash(hash); err != nil {
		return
	}

	return
}

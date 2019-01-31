package rest

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	api "github.com/uniris/uniris-core/api/protobuf-spec"
	"github.com/uniris/uniris-core/pkg/crypto"
	"google.golang.org/grpc"
)

//NewTransactionHandler creates a new transaction HTTP handler
func NewTransactionHandler(r *gin.RouterGroup, internalPort int) {
	r.GET("/transaction/:addr/status/:hash", getTransactionStatus(internalPort))
}

func getTransactionStatus(internalPort int) func(c *gin.Context) {
	return func(c *gin.Context) {
		txAddress := c.Param("addr")
		if _, err := crypto.IsHash(txAddress); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("address: %s", err.Error())})
			return
		}

		txHash := c.Param("hash")
		if _, err := crypto.IsHash(txHash); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("hash: %s", err.Error())})
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

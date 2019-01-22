package http

import (
	"context"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	api "github.com/uniris/uniris-core/api/protobuf-spec"
	"google.golang.org/grpc"
)

//NewTransactionHandler creates a new transaction HTTP handler
func NewTransactionHandler(r *gin.Engine) {
	r.GET("/transaction/:hash/status", getTransactionStatus())
}

func getTransactionStatus() func(c *gin.Context) {
	return func(c *gin.Context) {
		txHash := c.Param("hash")
		_, err := hex.DecodeString(txHash)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
		}

		serverAddr := fmt.Sprintf("localhost:1717")
		conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
		defer conn.Close()

		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		cli := api.NewTransactionServiceClient(conn)
		res, err := cli.GetTransactionStatus(context.Background(), &api.TransactionStatusRequest{
			Timestamp:        time.Now().Unix(),
			TransactionHash:  txHash,
			SignatureRequest: "", //TODO
		})
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		}

		c.JSON(http.StatusOK, res)
	}
}

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

//NewAccountHandler creates a new account HTTP handler
func NewAccountHandler(r *gin.Engine) {
	r.HEAD("/account/:hash", isAccountExist())
	r.GET("/account/:hash")
	r.POST("/account", createAccount())
}

func isAccountExist() func(c *gin.Context) {
	return func(c *gin.Context) {

		address := c.Param("address")
		sig := c.Query("signature")

		_, err := hex.DecodeString(address)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
		}

		_, err = hex.DecodeString(sig)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
		}

		serverAddr := fmt.Sprintf("localhost:1717")
		conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
		defer conn.Close()

		if err != nil {
			c.AbortWithError(500, err)
			return
		}
		cli := api.NewInternalServiceClient(conn)
		account, err := cli.GetAccount(context.Background(), &api.GetAccountRequest{
			EncryptedIdAddress: address,
			SignatureRequest:   sig,
		})
		if account == nil {
			c.Header("Account-Exist", "false")
			return
		}

		c.Header("Account-Exist", "true")
	}
}

func getAccount() func(c *gin.Context) {
	return func(c *gin.Context) {

		address := c.Param("address")
		sig := c.Query("signature")
		_, err := hex.DecodeString(address)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
		}

		_, err = hex.DecodeString(sig)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
		}

		serverAddr := fmt.Sprintf("localhost:1717")
		conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
		defer conn.Close()

		if err != nil {
			c.AbortWithError(500, err)
			return
		}
		cli := api.NewInternalServiceClient(conn)
		account, err := cli.GetAccount(context.Background(), &api.GetAccountRequest{
			EncryptedIdAddress: address,
			SignatureRequest:   sig,
		})

		if err != nil {
			c.AbortWithError(500, err)
			return
		}

		c.JSON(http.StatusOK, account)
	}
}

func createAccount() func(c *gin.Context) {
	return func(c *gin.Context) {

		var form accountCreationForm
		if err := c.ShouldBindJSON(&form); err != nil {
			c.AbortWithError(500, err)
			return
		}

		serverAddr := fmt.Sprintf("localhost:1717")
		conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
		defer conn.Close()

		if err != nil {
			c.AbortWithError(500, err)
			return
		}
		cli := api.NewInternalServiceClient(conn)
		resID, err := cli.HandleTransaction(context.Background(), &api.IncomingTransaction{
			EncryptedTransaction: form.EncryptedID,
			Timestamp:            time.Now().Unix(),
			Type:                 api.TransactionType_ID,
		})
		if err != nil {
			c.AbortWithError(500, err)
		}

		resKeychain, err := cli.HandleTransaction(context.Background(), &api.IncomingTransaction{
			EncryptedTransaction: form.EncryptedID,
			Timestamp:            time.Now().Unix(),
			Type:                 api.TransactionType_KEYCHAIN,
		})
		if err != nil {
			c.AbortWithError(500, err)
		}

		c.JSON(http.StatusCreated, accountCreationResult{
			IDTransactionHash:       resID.TransactionHash,
			KeychainTransactionHash: resKeychain.TransactionHash,
		})
	}
}

type accountCreationForm struct {
	EncryptedKeychain string `json:"encrypted_id"`
	EncryptedID       string `json:"encrypted_keychain"`
	Signature         string `json:"signature"`
}

type accountCreationResult struct {
	IDTransactionHash       string
	KeychainTransactionHash string
}

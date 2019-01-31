package rest

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	api "github.com/uniris/uniris-core/api/protobuf-spec"
	"github.com/uniris/uniris-core/pkg/crypto"
	"google.golang.org/grpc"
)

//NewAccountHandler creates a new account HTTP handler
func NewAccountHandler(apiGroup *gin.RouterGroup, intServerPort int, sharedPubk string) {
	apiGroup.GET("/account/:hash", getAccount(intServerPort, sharedPubk))
	apiGroup.POST("/account", createAccount(intServerPort, sharedPubk))
}

func getAccount(intServerPort int, sharedPubk string) func(c *gin.Context) {
	return func(c *gin.Context) {
		hash := c.Param("hash")
		if _, err := hex.DecodeString(hash); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("id hash: must be hexadecimal")})
			return
		}

		sig := c.Query("signature")
		if _, err := crypto.IsSignature(sig); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("request signature: %s", err.Error())})
			return
		}

		if err := crypto.VerifySignature(hash, sharedPubk, sig); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("request signature: %s", err.Error())})
			return
		}

		serverAddr := fmt.Sprintf("localhost:%d", intServerPort)
		conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
		defer conn.Close()

		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			return
		}
		cli := api.NewInternalServiceClient(conn)
		account, err := cli.GetAccount(context.Background(), &api.GetAccountRequest{
			EncryptedIdAddress: hash,
		})

		if err != nil {
			c.JSON(parseGrpcError(err))
			return
		}

		c.JSON(http.StatusOK, map[string]string{
			"encrypted_aes_key":  account.EncryptedAesKey,
			"encrypted_wallet":   account.EncryptedWallet,
			"signature_response": account.SignatureResponse,
		})
	}
}

func createAccount(intServerPort int, sharedPubk string) func(c *gin.Context) {
	return func(c *gin.Context) {

		var form struct {
			EncryptedID       string `json:"encrypted_id" binding:"required"`
			EncryptedKeychain string `json:"encrypted_keychain" binding:"required"`
			Signature         string `json:"signature" binding:"required"`
		}

		if err := c.ShouldBindJSON(&form); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if _, err := hex.DecodeString(form.EncryptedID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "encrypted id: must be hexadecimal"})
			return
		}

		if _, err := hex.DecodeString(form.EncryptedKeychain); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "encrypted keychain: must be hexadecimal"})
			return
		}

		if _, err := crypto.IsSignature(form.Signature); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("signature request: %s", err.Error())})
			return
		}
		formBytes, _ := json.Marshal(map[string]string{
			"encrypted_id":       form.EncryptedID,
			"encrypted_keychain": form.EncryptedKeychain,
		})
		if err := crypto.VerifySignature(string(formBytes), sharedPubk, form.Signature); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("signature request: %s", err.Error())})
			return
		}

		serverAddr := fmt.Sprintf("localhost:%d", intServerPort)
		conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
		defer conn.Close()

		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			return
		}
		cli := api.NewInternalServiceClient(conn)
		resID, err := cli.HandleTransaction(context.Background(), &api.IncomingTransaction{
			EncryptedTransaction: form.EncryptedID,
			Timestamp:            time.Now().Unix(),
		})
		if err != nil {
			c.JSON(parseGrpcError(err))
			return
		}

		resKeychain, err := cli.HandleTransaction(context.Background(), &api.IncomingTransaction{
			EncryptedTransaction: form.EncryptedKeychain,
			Timestamp:            time.Now().Unix(),
		})
		if err != nil {
			c.JSON(parseGrpcError(err))
			return
		}

		c.JSON(http.StatusCreated, map[string]interface{}{
			"id_transaction": map[string]interface{}{
				"transaction_hash": resID.TransactionHash,
				"timestamp":        resID.Timestamp,
				"signature":        resID.Signature,
			},
			"keychain_transaction": map[string]interface{}{
				"transaction_hash": resKeychain.TransactionHash,
				"timestamp":        resKeychain.Timestamp,
				"signature":        resKeychain.Signature,
			},
		})
	}
}

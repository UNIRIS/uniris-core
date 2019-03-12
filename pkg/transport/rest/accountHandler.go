package rest

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/uniris/uniris-core/pkg/consensus"

	api "github.com/uniris/uniris-core/api/protobuf-spec"
	"github.com/uniris/uniris-core/pkg/shared"

	"github.com/gin-gonic/gin"
	"github.com/uniris/uniris-core/pkg/crypto"
)

//GetAccountHandler is an HTTP handler which retrieves an account from an ID public key hash
//It requests the storage pool from the id address, decrypts the encrypted keychain address and request the keychain from its dedicated pool
//Then it aggregates the ID and Keychain data
func GetAccountHandler(techReader shared.TechDatabaseReader) func(c *gin.Context) {
	return func(c *gin.Context) {

		encIDHash := c.Param("idHash")
		if _, err := hex.DecodeString(encIDHash); err != nil {
			c.JSON(http.StatusBadRequest, httpError{
				Error:     "id hash: must be hexadecimal",
				Status:    http.StatusText(http.StatusBadRequest),
				Timestamp: time.Now().Unix(),
			})
			return
		}

		sigReq := c.Query("signature")
		if _, err := crypto.IsSignature(sigReq); err != nil {
			c.JSON(http.StatusBadRequest, httpError{
				Error:     fmt.Sprintf("signature request: %s", err.Error()),
				Status:    http.StatusText(http.StatusBadRequest),
				Timestamp: time.Now().Unix(),
			})
			return
		}

		emKeys, err := techReader.EmitterKeys()
		if err != nil {
			c.JSON(http.StatusInternalServerError, httpError{
				Error:     err.Error(),
				Status:    http.StatusText(http.StatusInternalServerError),
				Timestamp: time.Now().Unix(),
			})
			return
		}

		if err := crypto.VerifySignature(encIDHash, emKeys.RequestKey(), sigReq); err != nil {
			c.JSON(http.StatusBadRequest, httpError{
				Error:     fmt.Sprintf("signature request: %s", err.Error()),
				Status:    http.StatusText(http.StatusBadRequest),
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

		idHash, err := crypto.Decrypt(encIDHash, nodeLastKeys.PrivateKey())
		if err != nil {
			c.JSON(http.StatusBadRequest, httpError{
				Error:     err.Error(),
				Status:    http.StatusText(http.StatusBadRequest),
				Timestamp: time.Now().Unix(),
			})
			return
		}
		idTx, httpErr := findLastTransaction(idHash, api.TransactionType_ID, nodeLastKeys.PrivateKey())
		if httpErr != nil {
			httpErr.Error = fmt.Sprintf("ID: %s", httpErr.Error)
			c.JSON(httpErr.code, httpErr)
			return
		}

		encKeychainAddr := idTx.Data["encrypted_address_by_node"]
		keychainAddr, err := crypto.Decrypt(encKeychainAddr, nodeLastKeys.PrivateKey())
		if err != nil {
			c.JSON(http.StatusBadRequest, httpError{
				Error:     err.Error(),
				Status:    http.StatusText(http.StatusBadRequest),
				Timestamp: time.Now().Unix(),
			})
			return
		}

		keychainTx, httpErr := findLastTransaction(keychainAddr, api.TransactionType_KEYCHAIN, nodeLastKeys.PrivateKey())
		if httpErr != nil {
			httpErr.Error = fmt.Sprintf("Keychain: %s", httpErr.Error)
			c.JSON(httpErr.code, httpErr)
			return
		}

		res := accountFindResponse{
			EncryptedAESKey: idTx.Data["encrypted_aes_key"],
			EncryptedWallet: keychainTx.Data["encrypted_wallet"],
			Timestamp:       time.Now().Unix(),
		}

		resBytes, err := json.Marshal(res)
		if err != nil {
			c.JSON(http.StatusInternalServerError, httpError{
				Error:     err.Error(),
				Status:    http.StatusText(http.StatusInternalServerError),
				Timestamp: time.Now().Unix(),
			})
			return
		}
		sigRes, err := crypto.Sign(string(resBytes), nodeLastKeys.PrivateKey())
		if err != nil {
			c.JSON(http.StatusInternalServerError, httpError{
				Error:     err.Error(),
				Status:    http.StatusText(http.StatusInternalServerError),
				Timestamp: time.Now().Unix(),
			})
			return
		}

		res.Signature = sigRes

		c.JSON(http.StatusOK, res)
	}
}

//CreateAccountHandler is an HTTP handler which forwards ID and Keychain transaction to master nodes and reply with the transaction receipts
func CreateAccountHandler(techReader shared.TechDatabaseReader, p2pNodeReader consensus.NodeReader) func(c *gin.Context) {
	return func(c *gin.Context) {

		var form accountCreationRequest

		if err := c.ShouldBindJSON(&form); err != nil {
			c.JSON(http.StatusBadRequest, httpError{
				Error:     err.Error(),
				Status:    http.StatusText(http.StatusBadRequest),
				Timestamp: time.Now().Unix(),
			})
		}

		if _, err := hex.DecodeString(form.EncryptedID); err != nil {
			c.JSON(http.StatusBadRequest, httpError{
				Error:     "encrypted id: must be hexadecimal",
				Status:    http.StatusText(http.StatusBadRequest),
				Timestamp: time.Now().Unix(),
			})
			return
		}

		if _, err := hex.DecodeString(form.EncryptedKeychain); err != nil {
			c.JSON(http.StatusBadRequest, httpError{
				Error:     "encrypted keychain: must be hexadecimal",
				Status:    http.StatusText(http.StatusBadRequest),
				Timestamp: time.Now().Unix(),
			})
			return
		}

		if _, err := crypto.IsSignature(form.Signature); err != nil {
			c.JSON(http.StatusBadRequest, httpError{
				Error:     fmt.Sprintf("signature request: %s", err.Error()),
				Status:    http.StatusText(http.StatusBadRequest),
				Timestamp: time.Now().Unix(),
			})
			return
		}

		emKeys, err := techReader.EmitterKeys()
		if err != nil {
			c.JSON(http.StatusInternalServerError, httpError{
				Error:     err.Error(),
				Status:    http.StatusText(http.StatusInternalServerError),
				Timestamp: time.Now().Unix(),
			})
			return
		}

		formBytes, _ := json.Marshal(accountCreationRequest{
			EncryptedID:       form.EncryptedID,
			EncryptedKeychain: form.EncryptedKeychain,
		})
		if err := crypto.VerifySignature(string(formBytes), emKeys.RequestKey(), form.Signature); err != nil {
			c.JSON(http.StatusBadRequest, httpError{
				Error:     fmt.Sprintf("signature request: %s", err.Error()),
				Status:    http.StatusText(http.StatusBadRequest),
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

		idTx, err := decodeTransactionRaw(form.EncryptedID, nodeLastKeys.PrivateKey())
		if err != nil {
			c.JSON(http.StatusInternalServerError, httpError{
				Error:     err.Error(),
				Status:    http.StatusText(http.StatusInternalServerError),
				Timestamp: time.Now().Unix(),
			})
		}
		idTxRes, httpErr := requestTransactionMining(idTx, nodeLastKeys.PrivateKey(), nodeLastKeys.PublicKey(), p2pNodeReader, techReader)
		if httpErr != nil {
			c.JSON(httpErr.code, httpErr)
			return
		}

		keychainTx, err := decodeTransactionRaw(form.EncryptedKeychain, nodeLastKeys.PrivateKey())
		if err != nil {
			c.JSON(http.StatusInternalServerError, httpError{
				Error:     err.Error(),
				Status:    http.StatusText(http.StatusInternalServerError),
				Timestamp: time.Now().Unix(),
			})
		}
		keychainTxRes, httpErr := requestTransactionMining(keychainTx, nodeLastKeys.PrivateKey(), nodeLastKeys.PublicKey(), p2pNodeReader, techReader)
		if httpErr != nil {
			c.JSON(httpErr.code, httpErr)
			return
		}

		c.JSON(http.StatusCreated, accountCreationResponse{
			IDTransaction:       idTxRes,
			KeychainTransaction: keychainTxRes,
		})
	}
}

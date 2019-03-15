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
)

//GetAccountHandler is an HTTP handler which retrieves an account from an ID public key hash
//It requests the storage pool from the id address, decrypts the encrypted keychain address and request the keychain from its dedicated pool
//Then it aggregates the ID and Keychain data
func GetAccountHandler(sharedKeyReader shared.KeyReader, nodeReader consensus.NodeReader) func(c *gin.Context) {
	return func(c *gin.Context) {

		encIDHash := c.Param("idHash")
		encIDBytes, err := hex.DecodeString(encIDHash)
		if err != nil {
			c.JSON(http.StatusBadRequest, httpError{
				Error:     "id hash is not in hexadecimal",
				Status:    http.StatusText(http.StatusBadRequest),
				Timestamp: time.Now().Unix(),
			})
			return
		}

		sigReq := c.Query("signature")
		if sigReq == "" {
			c.JSON(http.StatusBadRequest, httpError{
				Error:     "signature is missing",
				Status:    http.StatusText(http.StatusBadRequest),
				Timestamp: time.Now().Unix(),
			})
			return
		}

		firstEmKeys, err := sharedKeyReader.FirstEmitterCrossKeypair()
		if err != nil {
			c.JSON(http.StatusBadRequest, httpError{
				Error:     "signature is not in hexadecimal",
				Status:    http.StatusText(http.StatusBadRequest),
				Timestamp: time.Now().Unix(),
			})
			return
		}

		sigBytes, err := hex.DecodeString(sigReq)
		if err != nil {
			c.JSON(http.StatusBadRequest, httpError{
				Error:     "signature is not in hexadecimal",
				Status:    http.StatusText(http.StatusBadRequest),
				Timestamp: time.Now().Unix(),
			})
			return
		}
		if !firstEmKeys.PublicKey().Verify(encIDBytes, sigBytes) {
			c.JSON(http.StatusBadRequest, httpError{
				Error:     "signature is invalid",
				Status:    http.StatusText(http.StatusBadRequest),
				Timestamp: time.Now().Unix(),
			})
			return
		}

		nodeLastKeys, err := sharedKeyReader.LastNodeCrossKeypair()
		if err != nil {
			c.JSON(http.StatusInternalServerError, httpError{
				Error:     err.Error(),
				Status:    http.StatusText(http.StatusInternalServerError),
				Timestamp: time.Now().Unix(),
			})
			return
		}

		idHash, err := nodeLastKeys.PrivateKey().Decrypt(encIDBytes)
		if err != nil {
			c.JSON(http.StatusBadRequest, httpError{
				Error:     err.Error(),
				Status:    http.StatusText(http.StatusBadRequest),
				Timestamp: time.Now().Unix(),
			})
			return
		}
		idTx, httpErr := findLastTransaction(idHash, api.TransactionType_ID, nodeLastKeys.PrivateKey(), nodeReader)
		if httpErr != nil {
			httpErr.Error = fmt.Sprintf("ID: %s", httpErr.Error)
			c.JSON(httpErr.code, httpErr)
			return
		}

		encKeychainAddr := idTx.Data["encrypted_address_by_node"]
		keychainAddr, err := nodeLastKeys.PrivateKey().Decrypt(encKeychainAddr)
		if err != nil {
			c.JSON(http.StatusBadRequest, httpError{
				Error:     err.Error(),
				Status:    http.StatusText(http.StatusBadRequest),
				Timestamp: time.Now().Unix(),
			})
			return
		}

		keychainTx, httpErr := findLastTransaction(keychainAddr, api.TransactionType_KEYCHAIN, nodeLastKeys.PrivateKey(), nodeReader)
		if httpErr != nil {
			httpErr.Error = fmt.Sprintf("Keychain: %s", httpErr.Error)
			c.JSON(httpErr.code, httpErr)
			return
		}

		res := accountFindResponse{
			EncryptedAESKey: hex.EncodeToString(idTx.Data["encrypted_aes_key"]),
			EncryptedWallet: hex.EncodeToString(keychainTx.Data["encrypted_wallet"]),
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

		sigRes, err := nodeLastKeys.PrivateKey().Sign(resBytes)
		if err != nil {
			c.JSON(http.StatusInternalServerError, httpError{
				Error:     err.Error(),
				Status:    http.StatusText(http.StatusInternalServerError),
				Timestamp: time.Now().Unix(),
			})
			return
		}

		res.Signature = hex.EncodeToString(sigRes)

		c.JSON(http.StatusOK, res)
	}
}

//CreateAccountHandler is an HTTP handler which forwards ID and Keychain transaction to master nodes and reply with the transaction receipts
func CreateAccountHandler(sharedKeyReader shared.KeyReader, nodeReader consensus.NodeReader) func(c *gin.Context) {
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
		firstEmKeys, err := sharedKeyReader.FirstEmitterCrossKeypair()
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

		sig, _ := hex.DecodeString(form.Signature)
		if !firstEmKeys.PublicKey().Verify(formBytes, sig) {
			c.JSON(http.StatusBadRequest, httpError{
				Error:     "signature is invalid",
				Status:    http.StatusText(http.StatusBadRequest),
				Timestamp: time.Now().Unix(),
			})
			return
		}

		nodeLastKeys, err := sharedKeyReader.LastNodeCrossKeypair()
		if err != nil {
			c.JSON(http.StatusInternalServerError, httpError{
				Error:     err.Error(),
				Status:    http.StatusText(http.StatusInternalServerError),
				Timestamp: time.Now().Unix(),
			})
			return
		}

		encIDBytes, _ := hex.DecodeString(form.EncryptedID)
		idTx, httpErr := decodeTransactionRaw(encIDBytes, nodeLastKeys.PrivateKey())
		if httpErr != nil {
			c.JSON(httpErr.code, httpErr)
			return
		}
		idTxRes, httpErr := requestTransactionMining(idTx, nodeLastKeys.PrivateKey(), nodeLastKeys.PublicKey(), nodeReader, sharedKeyReader)
		if httpErr != nil {
			c.JSON(httpErr.code, httpErr)
			return
		}

		encKeychainBytes, _ := hex.DecodeString(form.EncryptedKeychain)
		keychainTx, httpErr := decodeTransactionRaw(encKeychainBytes, nodeLastKeys.PrivateKey())
		if httpErr != nil {
			c.JSON(httpErr.code, httpErr)
			return
		}
		keychainTxRes, httpErr := requestTransactionMining(keychainTx, nodeLastKeys.PrivateKey(), nodeLastKeys.PublicKey(), nodeReader, sharedKeyReader)
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

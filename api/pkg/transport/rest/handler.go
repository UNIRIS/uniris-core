package rest

import (
	"net/http"

	"github.com/uniris/uniris-core/api/pkg/crypto"

	"github.com/gin-gonic/gin"

	"github.com/uniris/uniris-core/api/pkg/adding"
	"github.com/uniris/uniris-core/api/pkg/listing"
)

//ErrorMessage define an HTTP error
type ErrorMessage struct {
	Message string `json:"error_message"`
	Code    int    `json:"error_code"`
}

//Handler manages http rest methods handling
func Handler(r *gin.Engine, l listing.Service, a adding.Service) {

	api := r.Group("/api")
	{
		api.GET("/transaction/:addr/status/:hash", getTransactionStatus(l))
		api.POST("/account", createAccount(a))
		api.HEAD("/account/:hash", checkAccount(l))
		api.GET("/account/:hash", getAccount(l))
		api.GET("/sharedkeys/:publicKey", getSharedKeys(l))

		api.POST("/contract", createContract(a))
		api.POST("/contract/:addr/message", createContractMessage(a))
		api.GET("/contract/:addr/call/:method", getContractState(l))
	}
}

func createAccount(a adding.Service) func(c *gin.Context) {
	return func(c *gin.Context) {

		var req *accountRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			e := createError(http.StatusBadRequest, err)
			c.JSON(e.Code, e)
			return
		}

		res, err := a.AddAccount(adding.NewAccountCreationRequest(req.EncryptedID, req.EncryptedKeychain, req.Signature))
		if err != nil {
			if err == crypto.ErrInvalidSignature {
				e := createError(http.StatusBadRequest, err)
				c.JSON(e.Code, e)
				return
			}
			e := createError(http.StatusInternalServerError, err)
			c.JSON(e.Code, e)
			return
		}

		c.JSON(http.StatusCreated, accountCreationResult{
			Signature: res.Signature(),
			Transactions: accountCreationTransactionsResult{
				ID: transactionResult{
					MasterPeerIP:    res.ResultTransactions().ID().MasterPeerIP(),
					Signature:       res.ResultTransactions().ID().Signature(),
					TransactionHash: res.ResultTransactions().ID().TransactionHash(),
				},
				Keychain: transactionResult{
					MasterPeerIP:    res.ResultTransactions().Keychain().MasterPeerIP(),
					Signature:       res.ResultTransactions().Keychain().Signature(),
					TransactionHash: res.ResultTransactions().Keychain().TransactionHash(),
				},
			},
		})
	}
}

func checkAccount(l listing.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		hash := c.Param("hash")
		sig := c.Query("signature")

		err := l.ExistAccount(hash, sig)
		if err != nil {
			if err == crypto.ErrInvalidSignature {
				c.Header("Error", err.Error())
				return
			}
			if err == listing.ErrAccountNotExist {
				c.Header("Account-Exist", "false")
				return
			}
			e := createError(http.StatusInternalServerError, err)
			c.JSON(e.Code, e)
			return
		}

		c.Header("Account-Exist", "true")
	}
}

func getAccount(l listing.Service) func(c *gin.Context) {
	return func(c *gin.Context) {

		hash := c.Param("hash")
		sig := c.Query("signature")

		res, err := l.GetAccount(hash, sig)
		if err != nil {
			if err == crypto.ErrInvalidSignature {
				e := createError(http.StatusBadRequest, err)
				c.JSON(e.Code, e)
				return
			}
			if err == listing.ErrAccountNotExist {
				e := createError(http.StatusNotFound, err)
				c.JSON(e.Code, e)
				return
			}
			e := createError(http.StatusInternalServerError, err)
			c.JSON(e.Code, e)
			return
		}

		c.JSON(http.StatusOK, accountResult{
			EncryptedAddress: res.EncryptedAddress(),
			EncryptedAESKey:  res.EncryptedAESKey(),
			EncryptedWallet:  res.EncryptedWallet(),
			Signature:        res.Signature(),
		})
	}
}

func getSharedKeys(l listing.Service) func(c *gin.Context) {
	return func(c *gin.Context) {

		emPublicKey := c.Param("publicKey")
		sig := c.Query("signature")

		keys, err := l.GetSharedKeys(emPublicKey, sig)
		if err != nil {
			if err == listing.ErrUnauthorized {
				e := createError(http.StatusUnauthorized, err)
				c.JSON(e.Code, e)
				return
			}
			e := createError(http.StatusInternalServerError, err)
			c.JSON(e.Code, e)
			return
		}

		sharedEms := make([]sharedEmitterKeys, 0)
		for _, kp := range keys.EmitterKeyPairs() {
			sharedEms = append(sharedEms, sharedEmitterKeys{
				PublicKey:           kp.PublicKey(),
				EncryptedPrivateKey: kp.EncryptedPrivateKey(),
			})
		}

		c.JSON(http.StatusOK, sharedKeys{
			SharedEmitterKeys:    sharedEms,
			SharedRobotPublicKey: keys.RobotPublicKey(),
		})
	}
}

func getTransactionStatus(l listing.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		addr := c.Param("addr")
		txHash := c.Param("hash")

		status, err := l.GetTransactionStatus(addr, txHash)
		if err != nil {
			e := createError(http.StatusInternalServerError, err)
			c.JSON(e.Code, e)
			return
		}

		c.JSON(http.StatusOK, struct {
			Status string `json:"status"`
		}{
			Status: status.String(),
		})
	}
}

func createContract(a adding.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		var req contractCreationRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			e := createError(http.StatusBadRequest, err)
			c.JSON(e.Code, e)
			return
		}

		res, err := a.AddContract(adding.NewContractCreationRequest(req.EncryptedContract, req.Signature))
		if err != nil {
			if err == crypto.ErrInvalidSignature {
				e := createError(http.StatusBadRequest, err)
				c.JSON(e.Code, e)
				return
			}
			e := createError(http.StatusInternalServerError, err)
			c.JSON(e.Code, e)
			return
		}

		c.JSON(http.StatusCreated, transactionResult{
			TransactionHash: res.TransactionHash(),
			MasterPeerIP:    res.MasterPeerIP(),
			Signature:       res.Signature(),
		})
	}
}

func createContractMessage(a adding.Service) func(c *gin.Context) {
	return func(c *gin.Context) {

		address := c.Param("addr")

		var req contractMessageRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			e := createError(http.StatusBadRequest, err)
			c.JSON(e.Code, e)
			return
		}

		res, err := a.AddContractMessage(adding.NewContractMessageCreationRequest(address, req.EncryptedMessage, req.Signature))
		if err != nil {
			if err == crypto.ErrInvalidSignature {
				e := createError(http.StatusBadRequest, err)
				c.JSON(e.Code, e)
				return
			}
			e := createError(http.StatusInternalServerError, err)
			c.JSON(e.Code, e)
			return
		}

		c.JSON(http.StatusCreated, transactionResult{
			TransactionHash: res.TransactionHash(),
			MasterPeerIP:    res.MasterPeerIP(),
			Signature:       res.Signature(),
		})
	}
}

func getContractState(l listing.Service) func(c *gin.Context) {
	return func(c *gin.Context) {

		address := c.Param("addr")

		state, err := l.GetContractState(address)
		if err != nil {
			if err == crypto.ErrInvalidSignature {
				e := createError(http.StatusBadRequest, err)
				c.JSON(e.Code, e)
				return
			}
			e := createError(http.StatusInternalServerError, err)
			c.JSON(e.Code, e)
			return
		}

		c.JSON(http.StatusOK, contractState{
			Data:      state.Data(),
			Signature: state.Signature(),
		})
	}
}

func createError(handleErrorCode int, handleErr error) ErrorMessage {
	return ErrorMessage{
		Message: handleErr.Error(),
		Code:    http.StatusInternalServerError,
	}
}

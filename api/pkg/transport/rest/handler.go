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
func Handler(r *gin.Engine, robotPvKey string, l listing.Service, a adding.Service) {

	api := r.Group("/api")
	{
		api.POST("/account", createAccount(a, robotPvKey))
		api.HEAD("/account/:hash", checkAccount(l, robotPvKey))
		api.GET("/account/:hash", getAccount(l, robotPvKey))
	}
}

func createAccount(a adding.Service, robotPvKey string) func(c *gin.Context) {
	return func(c *gin.Context) {

		var req adding.AccountCreationRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			e := createError(http.StatusBadRequest, err, robotPvKey)
			c.JSON(e.Code, e)
			return
		}

		res, err := a.AddAccount(req)
		if err != nil {
			if err == crypto.ErrInvalidSignature {
				e := createError(http.StatusBadRequest, err, robotPvKey)
				c.JSON(e.Code, e)
				return
			}
			e := createError(http.StatusInternalServerError, err, robotPvKey)
			c.JSON(e.Code, e)
			return
		}

		c.JSON(http.StatusCreated, res)
	}
}

func checkAccount(l listing.Service, robotPvKey string) func(c *gin.Context) {
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
			e := createError(http.StatusInternalServerError, err, robotPvKey)
			c.JSON(e.Code, e)
			return
		}

		c.Header("Account-Exist", "true")
	}
}

func getAccount(l listing.Service, robotPvKey string) func(c *gin.Context) {
	return func(c *gin.Context) {

		hash := c.Param("hash")
		sig := c.Query("signature")

		acc, err := l.GetAccount(hash, sig)
		if err != nil {
			if err == crypto.ErrInvalidSignature {
				e := createError(http.StatusBadRequest, err, robotPvKey)
				c.JSON(e.Code, e)
				return
			}
			if err == listing.ErrAccountNotExist {
				e := createError(http.StatusNotFound, err, robotPvKey)
				c.JSON(e.Code, e)
				return
			}
			e := createError(http.StatusInternalServerError, err, robotPvKey)
			c.JSON(e.Code, e)
			return
		}

		c.JSON(http.StatusOK, acc)
	}
}

func createError(handleErrorCode int, handleErr error, robotPvKey string) ErrorMessage {
	return ErrorMessage{
		Message: handleErr.Error(),
		Code:    http.StatusInternalServerError,
	}
}

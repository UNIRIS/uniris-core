package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/uniris/uniris-core/api/pkg/adding"
	"github.com/uniris/uniris-core/api/pkg/listing"
)

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
			if err == adding.ErrInvalidSignature {
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
			if err == listing.ErrInvalidSignature {
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
			if err == listing.ErrInvalidSignature {
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

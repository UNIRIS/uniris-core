package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/uniris/uniris-core/api/pkg/adding"
	"github.com/uniris/uniris-core/api/pkg/listing"
)

//Handler manages http rest methods handling
func Handler(r *gin.Engine, l listing.Service, a adding.Service) {
	api := r.Group("/api")
	{
		api.POST("/account", enrollAccount(a))
		api.GET("/account/:hash", getAccount(l))
	}
}

func enrollAccount(a adding.Service) func(c *gin.Context) {
	return func(c *gin.Context) {

		var req adding.EnrollmentRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		res, err := a.AddAccount(req)
		if err != nil {
			if err == adding.ErrInvalidSignature {
				c.JSON(http.StatusBadRequest, err)
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, res)
	}
}

func getAccount(l listing.Service) func(c *gin.Context) {
	return func(c *gin.Context) {

		hash := c.Param("hash")
		sig := c.Query("signature")

		acc, err := l.GetAccount(hash, sig)
		if err != nil {
			if err == listing.ErrInvalidSignature {
				c.JSON(http.StatusBadRequest, err)
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, acc)
	}
}

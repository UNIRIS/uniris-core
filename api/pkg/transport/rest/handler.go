package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/uniris/uniris-core/api/pkg/adding"
	"github.com/uniris/uniris-core/api/pkg/listing"
)

//BuildAPI create a API router
func BuildAPI(r *gin.Engine, robotPvKey string, l listing.Service, a adding.Service) {

	api := r.Group("/api")
	{
		api.POST("/account", createAccount(a, robotPvKey))
		api.HEAD("/account/:hash", checkAccount(l, robotPvKey))
		api.GET("/account/:hash", getAccount(l, robotPvKey))

		api.GET("/peer/master", getMasterPeer(l, robotPvKey))
	}
}

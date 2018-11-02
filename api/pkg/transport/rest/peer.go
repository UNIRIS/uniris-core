package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/uniris/uniris-core/api/pkg/listing"
)

func getMasterPeer(l listing.Service, robotPvKey string) func(c *gin.Context) {
	return func(c *gin.Context) {

		masterPeer, err := l.GetMasterPeer()
		if err != nil {
			e := createError(http.StatusInternalServerError, err, robotPvKey)
			c.JSON(e.Code, e)
			return
		}

		c.JSON(http.StatusOK, masterPeer)
	}
}

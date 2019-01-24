package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/uniris/uniris-core/pkg/system"
	"github.com/uniris/uniris-core/pkg/transport/http"
)

func startAPI(conf system.UnirisConfig) {
	r := gin.Default()

	http.NewAccountHandler(r)
	http.NewTransactionHandler(r)

	r.Run(fmt.Sprintf(":%d", conf.Services.API.Port))
}

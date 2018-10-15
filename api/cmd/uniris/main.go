package main

import (
	"flag"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/uniris/uniris-core/api/pkg/adding"
	"github.com/uniris/uniris-core/api/pkg/crypto"
	"github.com/uniris/uniris-core/api/pkg/listing"
	"github.com/uniris/uniris-core/api/pkg/mock"
	"github.com/uniris/uniris-core/api/pkg/transport/rest"
)

func main() {

	port := flag.Int("port", 8080, "API port")
	flag.Parse()

	r := gin.Default()

	staticDir, _ := filepath.Abs("../../web/static")
	r.StaticFS("/static/", http.Dir(staticDir))

	rootPage, _ := filepath.Abs("../../web/index.html")
	r.StaticFile("/", rootPage)
	swaggerFile, _ := filepath.Abs("../../api/swagger-spec/swagger.yaml")
	r.StaticFile("/swagger.yaml", swaggerFile)

	sharedBioPrivKey := []byte("")
	client := mock.NewClient()
	validator := new(crypto.RequestValidator)
	lister := listing.NewService(sharedBioPrivKey, client, validator)
	adder := adding.NewService(sharedBioPrivKey, client, validator)

	rest.Handler(r, lister, adder)

	r.Run(fmt.Sprintf(":%d", *port))
}

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/uniris/uniris-core/api/pkg/adding"
	"github.com/uniris/uniris-core/api/pkg/crypto"
	"github.com/uniris/uniris-core/api/pkg/listing"
	"github.com/uniris/uniris-core/api/pkg/system"
	"github.com/uniris/uniris-core/api/pkg/transport/rest"
	"github.com/uniris/uniris-core/api/pkg/transport/rpc"
)

const (
	defaultConfFile = "../../../conf.yaml"
)

func main() {
	config, err := loadConfiguration()
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()

	staticDir, _ := filepath.Abs("../../web/static")
	r.StaticFS("/static/", http.Dir(staticDir))

	rootPage, _ := filepath.Abs("../../web/index.html")
	r.StaticFile("/", rootPage)
	swaggerFile, _ := filepath.Abs("../../api/swagger-spec/swagger.yaml")
	r.StaticFile("/swagger.yaml", swaggerFile)

	signer := crypto.NewSigner()
	client := rpc.NewRobotClient(config, signer)
	lister := listing.NewService(config, client, signer)
	adder := adding.NewService(config, client, signer)

	rest.Handler(r, lister, adder)

	r.Run(fmt.Sprintf(":%d", config.Services.API.Port))
}

func loadConfiguration() (conf system.UnirisConfig, err error) {
	confFile := flag.String("config", defaultConfFile, "Configuration file")
	flag.Parse()

	confFilePath, err := filepath.Abs(*confFile)
	conf, err = system.BuildFromFile(confFilePath)
	if err != nil {
		return
	}
	return conf, nil
}

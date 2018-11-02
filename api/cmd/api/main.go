package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
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
	defaultConfFile = "../../../default-conf.yml"
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

	client := rpc.NewRobotClient(config.Datamining, config.SharedKeys.RobotPrivateKey)
	validator := new(crypto.RequestValidator)
	lister := listing.NewService(config.SharedKeys.BiodPublicKey, client, validator)
	adder := adding.NewService(config.SharedKeys.BiodPublicKey, client, validator)

	rest.BuildAPI(r, config.SharedKeys.RobotPrivateKey, lister, adder)

	r.Run(fmt.Sprintf(":%d", config.API.Port))
}

func loadConfiguration() (*system.UnirisConfig, error) {
	confFile := flag.String("config", defaultConfFile, "Configuration file")
	flag.Parse()

	confFilePath, err := filepath.Abs(*confFile)
	if _, err := os.Stat(confFilePath); os.IsNotExist(err) {
		conf, err := system.BuildFromEnv()
		if err != nil {
			return nil, err
		}
		return conf, nil
	}

	conf, err := system.BuildFromFile(confFilePath)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

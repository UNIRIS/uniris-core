package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/urfave/cli"

	"github.com/uniris/uniris-core/pkg/crypto"
	"github.com/uniris/uniris-core/pkg/system"
)

func main() {

	conf := system.UnirisConfig{}

	app := cli.NewApp()
	app.Name = "uniris-miner"
	app.Usage = "UNIRIS miner"
	app.Version = "0.0.1"
	app.Flags = getCliFlags(&conf)
	app.Action = func(c *cli.Context) error {
		if c.String("private-key") == "" {
			fmt.Printf("Error: missing private key\n\n")
			return cli.ShowAppHelp(c)
		}

		if c.String("discovery-seeds") == "" {
			fmt.Printf("Error: missing seeds\n\n")
			return cli.ShowAppHelp(c)
		}

		conf.Version = app.Version

		pub, err := crypto.GetPublicKeyFromPrivate(conf.PrivateKey)
		if err != nil {
			panic(err)
		}
		conf.PublicKey = pub

		fmt.Println("----------")
		fmt.Println("UNIRIS MINER")
		fmt.Println("----------")
		fmt.Printf("Version: %s\n", conf.Version)
		fmt.Printf("Public key: %s\n", pub)
		fmt.Printf("Network: %s\n", conf.Network.Type)
		fmt.Printf("Network interface: %s\n", conf.Network.Interface)

		go startAPI(conf)
		go startDatamining(conf)

		time.Sleep(2 * time.Second)

		startDiscovery(conf)

		return nil
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

	// startDiscovery(conf)
	// go startDatamining(*conf)
	// startAPI(*conf)
}

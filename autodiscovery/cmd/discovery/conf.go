package main

import (
	"flag"
	"io/ioutil"
	"log"
)

func loadConfiguration() ([]byte, int, string, int) {
	port := flag.Int("port", 3545, "Discovery port")
	p2pFactor := flag.Int("p2p-factor", 1, "P2P replication factor")
	pbKeyFile := flag.String("public-key-file", "id.pub", "Public key file")

	pbKey, err := getPublicKey(*pbKeyFile)
	if err != nil {
		log.Panic(err)
	}

	flag.Parse()

	return nil, pbKey, *port, *p2pFactor
}

func getPublicKey(file string) ([]byte, error) {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

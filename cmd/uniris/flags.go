package main

import (
	"github.com/uniris/uniris-core/pkg/system"
	"github.com/urfave/cli"
)

func getCliFlags(conf *system.UnirisConfig) []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:        "private-key",
			Usage:       "Miner private key in hexadecimal",
			EnvVar:      "UNIRIS_PRIVATE_KEY",
			Destination: &conf.PrivateKey,
		},
		cli.StringFlag{
			Name:        "network",
			EnvVar:      "UNIRIS_NETWORK_TYPE",
			Value:       "public",
			Usage:       "Type of the blockchain network (public or private) - Help to identify the IP address",
			Destination: &conf.Network.Type,
		},
		cli.StringFlag{
			Name:        "network-interface",
			EnvVar:      "UNIRIS_NETWORK_INTERFACE",
			Usage:       "Name of the network interface when type of network is private",
			Destination: &conf.Network.Interface,
		},
		cli.IntFlag{
			Name:        "discovery-port",
			EnvVar:      "UNIRIS_DISCOVERY_PORT",
			Value:       4000,
			Usage:       "Discovery service port",
			Destination: &conf.Services.Discovery.Port,
		},
		cli.StringFlag{
			Name:        "discovery-seeds",
			EnvVar:      "UNIRIS_DISCOVERY_SEEDS",
			Usage:       "List of the seeds peers to bootstrap the miner `IP:PORT:PUBLIC_KEY;IP:PORT:PUBLIC_KEY`",
			Destination: &conf.Services.Discovery.Seeds,
		},
		cli.StringFlag{
			Name:        "discovery-redis-host",
			EnvVar:      "UNIRIS_DISCOVERY_REDIS_PORT",
			Value:       "localhost",
			Usage:       "Redis instance hostname",
			Destination: &conf.Services.Discovery.Redis.Host,
		},
		cli.IntFlag{
			Name:        "discovery-redis-port",
			Value:       6379,
			EnvVar:      "UNIRIS_DISCOVERY_REDIS_PORT",
			Usage:       "Redis instance port",
			Destination: &conf.Services.Discovery.Redis.Port,
		},
		cli.StringFlag{
			Name:        "discovery-redis-password",
			EnvVar:      "UNIRIS_DISCOVERY_REDIS_PWD",
			Usage:       "Redis instance password",
			Destination: &conf.Services.Discovery.Redis.Pwd,
		},
		cli.StringFlag{
			Name:        "discovery-rabbitmq-host",
			EnvVar:      "UNIRIS_DISCOVERY_RABBITMQ_HOST",
			Value:       "localhost",
			Usage:       "RabbitMQ instance hostname",
			Destination: &conf.Services.Discovery.AMQP.Host,
		},
		cli.IntFlag{
			Name:        "discovery-rabbitmq-port",
			EnvVar:      "UNIRIS_DISCOVERY_RABBITMQ_PORT",
			Value:       5672,
			Usage:       "Rabbitmq instance port",
			Destination: &conf.Services.Discovery.AMQP.Port,
		},
		cli.StringFlag{
			Name:        "discovery-rabbitmq-user",
			EnvVar:      "UNIRIS_DISCOVERY_RABBITMQ_USER",
			Value:       "guest",
			Usage:       "Rabbitmq instance user",
			Destination: &conf.Services.Discovery.AMQP.Username,
		},
		cli.StringFlag{
			Name:        "discovery-rabbitmq-password",
			EnvVar:      "UNIRIS_DISCOVERY_RABBITMQ_PWD",
			Value:       "guest",
			Usage:       "Rabbitmq instance password",
			Destination: &conf.Services.Discovery.AMQP.Password,
		},
		cli.IntFlag{
			Name:        "datamining-port",
			EnvVar:      "UNIRIS_DATAMINING_PORT",
			Value:       5000,
			Usage:       "Datamining port",
			Destination: &conf.Services.Datamining.ExternalPort,
		},
		cli.IntFlag{
			Name:        "datamining-internal-port",
			Value:       3009,
			Usage:       "Datamining internal port",
			Hidden:      true,
			Destination: &conf.Services.Datamining.InternalPort,
		},
		cli.IntFlag{
			Name:        "api-port",
			Value:       8080,
			Usage:       "API port",
			Destination: &conf.Services.API.Port,
		},
	}
}

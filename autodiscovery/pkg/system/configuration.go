package system

import (
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

//UnirisConfig describes the uniris robot main configuration
type UnirisConfig struct {
	Network          string `yaml:"network"`
	NetworkInterface string `yaml:"networkInterface"`
	PublicKey        string `yaml:"publicKey"`
	Version          string `yaml:"version"`
	Discovery        DiscoveryConfig
}

//DiscoveryConfig describes the autodiscovery configuration
type DiscoveryConfig struct {
	Port      int          `yaml:"port"`
	P2PFactor int          `yaml:"p2pFactor"`
	Seeds     []SeedConfig `yaml:"seeds"`
	Redis     RedisConfig  `yaml:"redis"`
	AMQP      AMQPConfig   `yaml:"amqp"`
}

//SeedConfig describes the autodiscovery seed configuration
type SeedConfig struct {
	IP        string `yaml:"ip"`
	Port      int    `yaml:"port"`
	PublicKey string `yaml:"publicKey"`
}

//RedisConfig describes the Redis database configuration
type RedisConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	Pwd  string `yaml:"pwd"`
}

//AMQPConfig describes the AMQP notifier configuration
type AMQPConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

//BuildFromEnv creates configurtion from env variables
func BuildFromEnv() (*UnirisConfig, error) {
	ver := os.Getenv("UNIRIS_VERSION")
	pbKey := os.Getenv("UNIRIS_PUBLICKEY")
	network := os.Getenv("UNIRIS_NETWORK")
	netiface := os.Getenv("UNIRIS_NETWORK_INTERFACE")
	port := os.Getenv("UNIRIS_DISCOVERY_PORT")
	p2pFactor := os.Getenv("UNIRIS_DISCOVERY_P2PFACTOR")
	seeds := os.Getenv("UNIRIS_DISCOVERY_SEEDS")
	redisHost := os.Getenv("UNIRIS_DISCOVERY_REDIS_HOST")
	redisPort := os.Getenv("UNIRIS_DISCOVERY_REDIS_PORT")
	redisPwd := os.Getenv("UNIRIS_DISCOVERY_REDIS_PWD")
	amqpHost := os.Getenv("UNIRIS_DISCOVERY_AMQP_HOST")
	amqpPort := os.Getenv("UNIRIS_DISCOVERY_AMQP_PORT")
	amqpUsername := os.Getenv("UNIRIS_DISCOVERY_AMQP_USER")
	amqpPassword := os.Getenv("UNIRIS_DISCOVERY_AMQP_PWD")

	_seeds := make([]SeedConfig, 0)
	ss := strings.Split(seeds, ",")
	for _, s := range ss {
		addr := strings.Split(s, ":")
		sPort, err := strconv.Atoi(addr[1])
		if err != nil {
			return nil, err
		}

		ips, err := net.LookupIP(addr[0])
		if err != nil {
			return nil, err
		}

		_seeds = append(_seeds, SeedConfig{
			IP:        ips[0].String(),
			Port:      sPort,
			PublicKey: addr[2],
		})
	}

	_port, err := strconv.Atoi(port)
	if err != nil {
		return nil, err
	}
	_p2pFactor, err := strconv.Atoi(p2pFactor)
	if err != nil {
		return nil, err
	}
	_redisPort, err := strconv.Atoi(redisPort)
	if err != nil {
		return nil, err
	}
	_amqpPort, err := strconv.Atoi(amqpPort)
	if err != nil {
		return nil, err
	}

	return &UnirisConfig{
		Version:          ver,
		PublicKey:        pbKey,
		Network:          network,
		NetworkInterface: netiface,
		Discovery: DiscoveryConfig{
			Port:      _port,
			P2PFactor: _p2pFactor,
			Seeds:     _seeds,
			Redis: RedisConfig{
				Host: redisHost,
				Port: _redisPort,
				Pwd:  redisPwd,
			},
			AMQP: AMQPConfig{
				Host:     amqpHost,
				Port:     _amqpPort,
				Username: amqpUsername,
				Password: amqpPassword,
			},
		},
	}, nil
}

//BuildFromFile creates configuration from configuration file
func BuildFromFile(confFile string) (*UnirisConfig, error) {
	confFilePath, err := filepath.Abs(confFile)
	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadFile(confFilePath)
	if err != nil {
		return nil, err
	}

	var conf UnirisConfig
	err = yaml.Unmarshal(bytes, &conf)
	if err != nil {
		return nil, err
	}

	return &conf, nil
}

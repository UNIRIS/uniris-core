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

	_seeds := make([]SeedConfig, 0)
	ss := strings.Split(seeds, ",")
	for _, s := range ss {
		addr := strings.Split(s, ":")
		sPort, _ := strconv.Atoi(addr[1])

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

	_port, _ := strconv.Atoi(port)
	_p2pFactor, _ := strconv.Atoi(p2pFactor)
	_redisPort, _ := strconv.Atoi(redisPort)

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

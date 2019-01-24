package system

//UnirisConfig describes the uniris robot main configuration
type UnirisConfig struct {
	Network    UnirisNetwork         `yaml:"network"`
	PublicKey  string                `yaml:"publicKey"`
	PrivateKey string                `yaml:"privateKey"`
	Version    string                `yaml:"version"`
	Services   ServicesConfiguration `yaml:"services"`
}

//ServicesConfiguration describe the robot services configuration
type ServicesConfiguration struct {
	API        APIConfiguration        `yaml:"api"`
	Discovery  UnirisDiscoveryConfig   `yaml:"discovery"`
	Datamining DataMiningConfiguration `yaml:"datamining"`
}

//UnirisNetwork describe the robot network
type UnirisNetwork struct {
	Type      string `yaml:"type"`
	Interface string `yaml:"interface"`
}

//APIConfiguration describes the api service configuration
type APIConfiguration struct {
	Port int `yaml:"port"`
}

//UnirisDiscoveryConfig describes the autodiscovery configuration
type UnirisDiscoveryConfig struct {
	Port  int         `yaml:"port"`
	Seeds string      `yaml:"seeds"`
	Redis RedisConfig `yaml:"redis"`
	AMQP  AMQPConfig  `yaml:"amqp"`
}

//DataMiningConfiguration describes the datamining configuration
type DataMiningConfiguration struct {
	InternalPort int `yaml:"internalPort"`
	ExternalPort int `yaml:"externalPort"`
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

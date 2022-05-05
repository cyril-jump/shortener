package config

// config

type EnvVar struct {
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:":8080"`
	BaseURL       string `env:"BASE_URL" envDefault:"http://localhost:8080"`
}

//http://localhost:8080/f845599b098517893fc2712d32774f53
//http://localhost:8080/f845599b098517893fc2712d32774f53
type Config struct {
	serverAddress string
	baseURL       string
}

//getters

func (c Config) SrvAddr() string {
	return c.serverAddress
}

func (c Config) HostName() string {
	return c.baseURL
}

//constructor

func NewConfig(srvAddr, hostName string) *Config {
	return &Config{
		serverAddress: srvAddr,
		baseURL:       hostName,
	}
}

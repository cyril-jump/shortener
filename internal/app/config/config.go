package config

// config

type Config struct {
	serverAddress   string
	baseURL         string
	fileStoragePath string
}

//getters

func (c Config) SrvAddr() string {
	return c.serverAddress
}

func (c Config) HostName() string {
	return c.baseURL
}

func (c Config) FileStoragePath() string {
	return c.fileStoragePath
}

//constructor

func NewConfig(srvAddr, hostName, fileStoragePath string) *Config {
	return &Config{
		serverAddress:   srvAddr,
		baseURL:         hostName,
		fileStoragePath: fileStoragePath,
	}
}

package config

// config

type Config struct {
	srvAddr  string
	hostName string
}

//getters

func (c Config) SrvAddr() string {
	return c.srvAddr
}

func (c Config) HostName() string {
	return c.hostName
}

//constructor

func NewConfig(srvAddr, hostName string) *Config {
	return &Config{
		srvAddr:  srvAddr,
		hostName: hostName,
	}
}

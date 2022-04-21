package storage

//data base

type DB struct {
	StorageURL map[string]string
}

func NewDB() *DB {
	return &DB{
		StorageURL: make(map[string]string),
	}
}

//config

type Config struct {
	SrvAddr  string
	HostName string
}

func NewConfig(srvAddr, hostName string) *Config {
	return &Config{
		SrvAddr:  srvAddr,
		HostName: hostName,
	}
}

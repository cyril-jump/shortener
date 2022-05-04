package interfaces

//storage

type Storage interface {
	BaseURL(shortURL string) string
	SetURL(shortURL, baseURL string)
}

//config

type Config interface {
	SrvAddr() string
	HostName() string
}

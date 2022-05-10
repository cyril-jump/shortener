package interfaces

import "errors"

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
)

//storage

type Storage interface {
	GetBaseURL(shortURL string) (string, error)
	SetShortURL(shortURL, baseURL string) error
	Close() error
}

//config

type Config interface {
	SrvAddr() string
	HostName() string
}

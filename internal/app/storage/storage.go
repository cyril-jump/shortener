package storage

type DB interface {
	GetBaseURL(shortURL string) (string, error)
	SetShortURL(shortURL, baseURL string) error
	Close()
}

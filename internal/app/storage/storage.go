package storage

// DB

type DB struct {
	storageURL map[string]string
}

//getters

func (D *DB) BaseURL(shortURL string) string {
	return D.storageURL[shortURL]
}

//setters

func (D *DB) SetURL(shortURL, baseURL string) {
	D.storageURL[shortURL] = baseURL
}

//constructor

func NewDB() *DB {
	return &DB{storageURL: make(map[string]string)}
}

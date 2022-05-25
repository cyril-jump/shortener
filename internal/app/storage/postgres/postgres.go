package postgres

import (
	"database/sql"
	"github.com/cyril-jump/shortener/internal/app/storage"
	_ "github.com/lib/pq"
	"log"
)

type DB struct {
	db *sql.DB
}

func New(psqlConn string) *DB {
	db, err := sql.Open("postgres", psqlConn)
	if err != nil {
		log.Fatal(err)
	}
	// close database
	//defer db.Close()

	// check db
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to DB!")
	return &DB{
		db: db,
	}
}

func (D *DB) GetBaseURL(shortURL string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (D *DB) GetAllURLsByUserID(userID string) ([]storage.ModelURL, error) {
	//TODO implement me
	panic("implement me")
}

func (D *DB) SetShortURL(userID, shortURL, baseURL string) error {
	//TODO implement me
	panic("implement me")
}

func (D *DB) Ping() error {
	return D.db.Ping()
}

func (D *DB) Close() {
	D.Close()
}

package postgres

import (
	"database/sql"
	"github.com/cyril-jump/shortener/internal/app/storage"
	_ "github.com/jackc/pgx/v4/stdlib"
	"log"
)

type DB struct {
	db *sql.DB
}

func New(psqlConn string) *DB {
	db, err := sql.Open("pgx", psqlConn)
	if err != nil {
		log.Fatal(err)
	}

	// check db
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	if _, err = db.Exec(schema); err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to DB!")
	return &DB{
		db: db,
	}
}

func (D *DB) GetBaseURL(shortURL string) (string, error) {
	var URL string
	selectStmt, err := D.db.Prepare("SELECT url FROM urls WHERE short_url=$1;")
	if err != nil {
		return "", err
	}
	defer selectStmt.Close()

	if err = selectStmt.QueryRow(shortURL).Scan(&URL); err != nil {
		return "", err
	}
	return URL, nil

}

func (D *DB) GetAllURLsByUserID(userID string) ([]storage.ModelURL, error) {
	//TODO implement me
	panic("implement me")
}

func (D *DB) SetShortURL(userID, shortURL, baseURL string) error {
	var URL string
	selectStmt, err := D.db.Prepare("SELECT short_url FROM urls WHERE url=$1 and user_id=$2;")
	if err != nil {
		return err
	}
	defer selectStmt.Close()

	insertStmt, err := D.db.Prepare("INSERT INTO urls (user_id, url, short_url) VALUES ($1, $2, $3);")
	if err != nil {
		return err
	}
	defer insertStmt.Close()

	if err = selectStmt.QueryRow(baseURL, userID).Scan(&URL); err == sql.ErrNoRows {
		_, err = insertStmt.Exec(userID, baseURL, shortURL)
		if err != nil {
			return err
		}
	}

	return err
}

func (D *DB) Ping() error {
	return D.db.Ping()
}

func (D *DB) Close() error {
	return D.db.Close()
}

var schema = `CREATE TABLE IF NOT EXISTS urls (
		id bigserial not null,
		user_id text not null,
		url text not null,
		short_url text not null 
	);`

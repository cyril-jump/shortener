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
	//TODO implement me
	panic("implement me")
}

func (D *DB) GetAllURLsByUserID(userID string) ([]storage.ModelURL, error) {
	//TODO implement me
	panic("implement me")
}

func (D *DB) SetShortURL(userID, shortURL, baseURL string) error {
	var URL string
	var selectStmt, err = D.db.Prepare("SELECT short_url FROM urls WHERE url=$1 and user_id=$2;")
	if err != nil {
		return err
	}

	defer func(selectStmt *sql.Stmt) {
		err := selectStmt.Close()
		if err != nil {

		}
	}(selectStmt)
	err = selectStmt.QueryRow(baseURL, userID).Scan(&URL)
	if err != nil {
		log.Println(err)
		var insertStmt, err = D.db.Prepare("INSERT INTO urls (user_id, url, short_url) VALUES ($1, $2, $3);")
		if err != nil {
			return err
		}
		defer func(insertStmt *sql.Stmt) {
			err := insertStmt.Close()
			if err != nil {

			}
		}(insertStmt)
		_, err = insertStmt.Exec(userID, baseURL, shortURL)
		if err != nil {
			return err
		}

	}
	return nil
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

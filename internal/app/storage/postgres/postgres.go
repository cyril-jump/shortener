package postgres

import (
	"context"
	"database/sql"
	"github.com/cyril-jump/shortener/internal/app/storage"
	"github.com/cyril-jump/shortener/internal/app/utils/errs"
	_ "github.com/jackc/pgx/v4/stdlib"
	"log"
)

type DB struct {
	db  *sql.DB
	ctx context.Context
}

func (D *DB) GetBaseURL(shortURL string) (string, error) {
	var baseURL string
	selectStmt, err := D.db.Prepare("SELECT base_url FROM urls WHERE short_url=$1;")
	if err != nil {
		return "", err
	}
	defer selectStmt.Close()

	if err = selectStmt.QueryRow(shortURL).Scan(&baseURL); err != nil {
		return "", err
	}
	return baseURL, nil

}

func New(ctx context.Context, psqlConn string) *DB {
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
		db:  db,
		ctx: ctx,
	}
}

func (D *DB) GetAllURLsByUserID(userID string) ([]storage.ModelURL, error) {
	var modelURL []storage.ModelURL
	var model storage.ModelURL
	selectStmt, err := D.db.Prepare("SELECT short_url, base_url FROM users_url RIGHT JOIN urls u on users_url.url_id=u.id WHERE user_id=$1;")
	if err != nil {
		return nil, err
	}
	defer selectStmt.Close()

	row, err := selectStmt.Query(userID)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	if err = row.Err(); err != nil {
		log.Println(err)
	}

	for row.Next() {
		err := row.Scan(&model.ShortURL, &model.BaseURL)
		if err != nil {
			return nil, err
		}
		modelURL = append(modelURL, model)
	}

	return modelURL, nil
}

func (D *DB) SetShortURL(userID, shortURL, baseURL string) error {
	var id int
	var userURLID int

	insertStmt1, err := D.db.Prepare("INSERT INTO urls (base_url, short_url) VALUES ($1, $2) RETURNING id")
	if err != nil {
		return err
	}
	defer insertStmt1.Close()

	insertStmt2, err := D.db.Prepare("INSERT INTO users_url (user_id, url_id) VALUES ($1, $2);")
	if err != nil {
		return err
	}
	defer insertStmt2.Close()

	selectStmt, err := D.db.Prepare("SELECT id FROM urls WHERE base_url = $1;")
	if err != nil {
		return err
	}
	defer selectStmt.Close()

	insertStmt1.QueryRow(baseURL, shortURL).Scan(&id)
	if id != 0 {
		_, err = insertStmt2.Exec(userID, id)
		if err != nil {
			log.Println(errs.ErrAlreadyExists)
			return errs.ErrAlreadyExists
		}
	} else {
		selectStmt.QueryRow(baseURL).Scan(&userURLID)
		_, err := insertStmt2.Exec(userID, userURLID)
		if err != nil {
			return errs.ErrAlreadyExists
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

var schema = `
	CREATE TABLE IF NOT EXISTS urls (
		id serial primary key,
		base_url text not null unique,
		short_url text not null 
	);
	CREATE TABLE IF NOT EXISTS users_url(
	  user_id text not null,
	  url_id int not null references urls(id),
	  CONSTRAINT unique_url UNIQUE (user_id, url_id)
	);
	`

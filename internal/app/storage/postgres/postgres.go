package postgres

import (
	"context"
	"database/sql"
	"log"
	"sync"

	_ "github.com/jackc/pgx/v4/stdlib"

	"github.com/cyril-jump/shortener/internal/app/dto"
	"github.com/cyril-jump/shortener/internal/app/utils/errs"
)

//DB struct
type DB struct {
	mu  sync.Mutex
	db  *sql.DB
	ctx context.Context
}

//New DB constructor
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

//GetBaseURL Get base URL from DB
func (D *DB) GetBaseURL(shortURL string) (string, error) {
	D.mu.Lock()
	var baseURL string
	countURL := 0
	selectStmt, err := D.db.Prepare("SELECT base_url, count_url FROM urls WHERE short_url=$1;")
	if err != nil {
		return "", err
	}
	defer func() {
		selectStmt.Close()
		D.mu.Unlock()
	}()

	if err = selectStmt.QueryRow(shortURL).Scan(&baseURL, &countURL); err != nil {
		return "", err
	}
	if countURL == 0 {
		return "", errs.ErrWasDeleted
	}
	return baseURL, nil

}

//GetAllURLsByUserID Get all URLs by UserID from DB
func (D *DB) GetAllURLsByUserID(userID string) ([]dto.ModelURL, error) {
	D.mu.Lock()
	modelURL := make([]dto.ModelURL, 20000)
	model := dto.ModelURL{}
	selectStmt, err := D.db.Prepare("SELECT short_url, base_url FROM users_url RIGHT JOIN urls u on users_url.url_id=u.id WHERE user_id=$1 AND  count_url > 0;")
	if err != nil {
		return nil, err
	}
	defer func() {
		selectStmt.Close()
		D.mu.Unlock()
	}()

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

//SetShortURL Set short URL in DB
func (D *DB) SetShortURL(userID, shortURL, baseURL string) error {
	D.mu.Lock()
	var id, userURLID int

	insertStmt1, err := D.db.Prepare("INSERT INTO urls (base_url, short_url) VALUES ($1, $2) RETURNING (id)")
	if err != nil {
		return err
	}

	insertStmt2, err := D.db.Prepare("INSERT INTO users_url (user_id, url_id)  VALUES ($1, $2);")
	if err != nil {
		return err
	}

	selectStmt1, err := D.db.Prepare("SELECT id FROM urls WHERE base_url = $1;")
	if err != nil {
		return err
	}

	updateStmt1, err := D.db.Prepare("UPDATE urls SET count_url = count_url + 1  WHERE base_url = $1;")
	if err != nil {
		return err
	}

	tx, err := D.db.BeginTx(D.ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		insertStmt1.Close()
		insertStmt2.Close()
		selectStmt1.Close()
		updateStmt1.Close()
		tx.Rollback()
		D.mu.Unlock()
	}()

	insertStmt1.QueryRow(baseURL, shortURL).Scan(&id)
	if id != 0 {
		_, err = tx.StmtContext(D.ctx, insertStmt2).ExecContext(D.ctx, userID, id)
		if err != nil {
			return err
		}
	} else {
		selectStmt1.QueryRow(baseURL).Scan(&userURLID)
		_, err = tx.StmtContext(D.ctx, insertStmt2).ExecContext(D.ctx, userID, userURLID)
		if err != nil {
			return errs.ErrAlreadyExists
		}
	}
	_, err = tx.StmtContext(D.ctx, updateStmt1).ExecContext(D.ctx, baseURL)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

//DelBatchShortURLs Delete batch short URLs in DB
func (D *DB) DelBatchShortURLs(tasks []dto.Task) error {
	D.mu.Lock()
	id := 0
	updateStmt2, err := D.db.Prepare("UPDATE users_url SET is_deleted = true WHERE user_id = $1 AND url_id = $2")
	if err != nil {
		return err
	}
	updateStmt1, err := D.db.Prepare("UPDATE urls SET count_url = count_url - 1  WHERE short_url = $1 RETURNING id;")
	if err != nil {
		return err
	}
	tx, err := D.db.BeginTx(D.ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		updateStmt1.Close()
		updateStmt2.Close()
		D.mu.Unlock()
		tx.Rollback()
	}()

	for _, t := range tasks {
		_ = tx.StmtContext(D.ctx, updateStmt1).QueryRowContext(D.ctx, t.ShortURL).Scan(&id)
		if err != nil {
			return err
		}

		_, err = tx.StmtContext(D.ctx, updateStmt2).ExecContext(D.ctx, t.ID, id)
		if err != nil {
			return err
		}
		id = 0
	}
	tx.Commit()
	return nil
}

//Ping Ping DB
func (D *DB) Ping() error {
	return D.db.Ping()
}

//Close Close DB connection
func (D *DB) Close() error {
	return D.db.Close()
}

//schema DB schema
var schema = `
	CREATE TABLE IF NOT EXISTS urls (
		id serial primary key,
		base_url text not null unique,
		short_url text not null,
	    count_url integer not null DEFAULT 0
	);
	CREATE TABLE IF NOT EXISTS users_url(
	  user_id text not null,
	  url_id int not null references urls(id),
	  is_deleted boolean not null DEFAULT false, 
	  CONSTRAINT unique_url UNIQUE (user_id, url_id)
	);
	`

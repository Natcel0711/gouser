package database

import (
	"database/sql"
	"errors"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

func StartDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", os.Getenv("GOINONE_DB_DSN"))
	if err != nil {
		return nil, errors.New("error connecting to DB")
	}
	return db, nil
}

func CloseDB() error {
	err := db.Close()
	if err != nil {
		return errors.New("something happened while closing db")
	}
	return nil
}

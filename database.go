package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type database struct {
	db    *sql.DB
	prod  bool
	errCh chan error
}

func Database(source string, errCh chan error) (database, error) {

	db, err := sql.Open("sqlite3", source)
	if err != nil {
		return database{}, fmt.Errorf("Failed to Open database: %s", err)
	}

	prod := os.Getenv("production") == "true"

	return database{
		prod:  prod,
		db:    db,
		errCh: errCh,
	}, nil

}

func (d *database) Destroy() {
	d.db.Close()
}

func (d *database) InsertUser(u User) {
	q := `insert into users (id, name, url) values (%v, %v, %s)`
	d.db.Exec(q, u.id, u.name, u.urls)
}

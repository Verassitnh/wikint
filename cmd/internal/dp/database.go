package dp

import (
	"database/sql"
	"fmt"
	"log"
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
		return database{}, fmt.Errorf("failed to Open database: %s", err)
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

func (d *database) AppendUser(u User) {
	q := `insert into users (id, name, url) values ("%v", "%v", "%+v")`
	res, err := d.db.Exec(fmt.Sprintf(q, u.id, u.name, u.urls))

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res)
}

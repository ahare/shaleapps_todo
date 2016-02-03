package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var (
	db      *sql.DB
	dbPath  = "./todo.db"
	queries = map[string]string{}
	prep    = map[string]*sql.Stmt{}
)

func Load() {
	var err error
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("unable to open db: %q", err)
	}

	err = createSchema()
	if err != nil {
		log.Fatalf("unable to create schema: %q", err)
	}
}

func Close() {
	db.Close()
}

func createSchema() error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS todos (
						id INTEGER NOT NULL PRIMARY KEY,
						text TEXT, done INTEGER);`)
	return err
}

func getPrep(name string) (*sql.Stmt, error) {
	if p, ok := prep[name]; ok {
		return p, nil
	}

	if s, ok := queries[name]; ok {
		p, err := db.Prepare(s)
		if err != nil {
			return nil, fmt.Errorf("can't prepare %q: %q", s, err)
		}
		prep[name] = p
		return p, nil
	}

	return nil, fmt.Errorf("no query named %q is defined", name)
}

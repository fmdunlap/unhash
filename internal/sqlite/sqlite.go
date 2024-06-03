package sqlite

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type SqliteStore struct {
	sq3 *sql.DB
}

func NewSqliteStore(dataStore string, clearOnStartup bool) *SqliteStore {
	if clearOnStartup {
		err := os.Remove(dataStore)
		if err != nil {
			log.Println("Error clearing data persistence", err)
		}
	}
	if _, err := os.Stat(dataStore); os.IsNotExist(err) {
		_, err := os.Create(dataStore)
		if err != nil {
			panic(err)
		}
	}

	db, err := sql.Open("sqlite3", dataStore)
	if err != nil {
		panic(err)
	}

	if clearOnStartup {
		_, err := db.Exec("create table users (id text, data jsonb)")
		if err != nil {
			panic(err)
		}
		_, err = db.Exec("create table hashjobs (id text, data jsonb)")
		if err != nil {
			panic(err)
		}
	}

	return &SqliteStore{sq3: db}
}

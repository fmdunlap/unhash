//go:build testing

package sqlite

import "database/sql"

func CreateTestDb() *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("create table users (id text, data jsonb)")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("create table hashjobs (id text, data jsonb)")
	if err != nil {
		panic(err)
	}

	return db
}

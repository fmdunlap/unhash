package sqlite

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/fmdunlap/unhash/internal/hashjob"
	"github.com/fmdunlap/unhash/internal/uerr"
	"github.com/fmdunlap/unhash/internal/user"
	_ "github.com/mattn/go-sqlite3"
)

type SqliteStore struct {
	dataStore string
	db        *sql.DB
}

func NewSqliteStore(dataStore string, clearOnStartup bool) *SqliteStore {
	if clearOnStartup {
		err := os.Remove(dataStore)
		if err != nil {
			log.Println("Error clearing data store", err)
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

	return &SqliteStore{dataStore: dataStore, db: db}
}

func (s *SqliteStore) InsertUser(u user.User) error {
	if u.ID == "" {
		return errors.New("id is required")
	}

	statement, err := s.db.Prepare("insert into users (id, data) values (?, ?)")
	if err != nil {
		return err
	}
	defer statement.Close()

	rawData, err := json.Marshal(u)
	if err != nil {
		return err
	}

	_, err = statement.Exec(u.ID, rawData)
	if err != nil {
		return &uerr.ErrorCannotInsert{Err: err}
	}

	return nil
}

func (s *SqliteStore) GetUser(id string) (*user.User, error) {
	var data []byte
	err := s.db.QueryRow("select data from users where id = ?", id).Scan(&data)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &uerr.ErrorNotFound{Err: err}
		}

		return nil, err
	}

	u, err := user.Unmarshal(data)
	if err != nil {
		return nil, err
	}

	if u.ID == "" {
		return nil, &uerr.ErrorNotFound{Err: errors.New("user not found")}
	}

	return u, nil
}

func (s *SqliteStore) GetUserByEmail(email string) (*user.User, error) {
	var data []byte
	err := s.db.QueryRow("select data from users where data->>'email' = ?", email).Scan(&data)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &uerr.ErrorNotFound{Err: err}
		}

		return nil, err
	}

	u, err := user.Unmarshal(data)
	if err != nil {
		return nil, err
	}

	if u.ID == "" {
		return nil, &uerr.ErrorNotFound{Err: errors.New("user not found")}
	}

	return u, nil
}

func (s *SqliteStore) ListUsers() ([]user.User, error) {
	rows, err := s.db.Query("select data from users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]user.User, 0)
	for rows.Next() {
		var data []byte
		err = rows.Scan(&data)
		if err != nil {
			return nil, err
		}

		u, err := user.Unmarshal(data)
		if err != nil {
			return nil, err
		}

		if u.ID == "" {
			return nil, &uerr.ErrorNotFound{Err: errors.New("user not found")}
		}

		users = append(users, *u)
	}

	return users, nil
}

func (s *SqliteStore) DeleteUser(id string) error {
	statement, err := s.db.Prepare("delete from users where id = ?")
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(id)
	if err != nil {
		return &uerr.ErrorCannotDelete{Err: err}
	}

	return nil
}

func (s *SqliteStore) InsertHashJob(h hashjob.HashJob) error {
	if h.ID == "" {
		return errors.New("id is required")
	}

	statement, err := s.db.Prepare("insert into hashjobs (id, data) values (?, ?)")
	if err != nil {
		return err
	}
	defer statement.Close()

	rawData, err := json.Marshal(h)
	if err != nil {
		return err
	}

	_, err = statement.Exec(h.ID, rawData)
	if err != nil {
		return &uerr.ErrorCannotInsert{Err: err}
	}

	return nil
}

func (s *SqliteStore) GetHashJob(id string) (*hashjob.HashJob, error) {
	var data []byte
	err := s.db.QueryRow("select data from hashjobs where id = ?", id).Scan(&data)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &uerr.ErrorNotFound{Err: err}
		}
		return nil, err
	}

	h, err := hashjob.Unmarshal(data)
	if err != nil {
		return nil, err
	}

	if h.ID == "" {
		return nil, &uerr.ErrorNotFound{Err: errors.New("hashjob not found")}
	}

	return h, nil
}

func (s *SqliteStore) DeleteHashJob(id string) error {
	statement, err := s.db.Prepare("delete from hashjobs where id = ?")
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(id)
	if err != nil {
		return &uerr.ErrorCannotDelete{Err: err}
	}

	return nil
}

package sqlite

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/fmdunlap/unhash/internal/hashjob"
	"github.com/fmdunlap/unhash/internal/uerr"
)

func (s *SqliteStore) InsertHashJob(h hashjob.HashJob) error {
	if h.ID == "" {
		return errors.New("id is required")
	}
	if h.OwnerId == "" {
		return errors.New("owner id is required")
	}

	statement, err := s.sq3.Prepare("insert into hashjobs (id, data) values (?, ?)")
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
	err := s.sq3.QueryRow("select data from hashjobs where id = ?", id).Scan(&data)
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
	statement, err := s.sq3.Prepare("delete from hashjobs where id = ?")
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

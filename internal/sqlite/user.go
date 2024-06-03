package sqlite

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/fmdunlap/unhash/internal/uerr"
	"github.com/fmdunlap/unhash/internal/user"
)

func (s *SqliteStore) InsertUser(u *user.User) error {
	if u.ID == "" {
		return errors.New("id is required")
	}

	statement, err := s.sq3.Prepare("insert into users (id, data) values (?, ?)")
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

func (s *SqliteStore) UpdateUser(u *user.User) error {
	rawData, err := json.Marshal(u)
	if err != nil {
		return err
	}

	_, err = s.sq3.Query("UPDATE users SET data = ? WHERE id = ?", u.ID, rawData)
	if err != nil {
		return &uerr.ErrorCannotUpdate{Err: err}
	}

	return nil
}

func (s *SqliteStore) GetUser(id string) (*user.User, error) {
	var data []byte
	err := s.sq3.QueryRow("select data from users where id = ?", id).Scan(&data)
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
	err := s.sq3.QueryRow("select data from users where data->>'email' = ?", email).Scan(&data)
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
	rows, err := s.sq3.Query("select data from users")
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
	statement, err := s.sq3.Prepare("delete from users where id = ?")
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

func (s *SqliteStore) UserExists(u *user.User) bool {
	rows, err := s.sq3.Query("SELECT id FROM users WHERE id = ? OR data->>'email' = ? OR data->>'username' = ? LIMIT 1", u.ID, u.Email, u.Username)
	if err != nil {
		return false
	}
	defer rows.Close()

	return rows.Next()
}

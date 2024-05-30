package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/fmdunlap/unhash/internal/uerr"
)

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type UserStore interface {
	InsertUser(u User) error
	GetUser(id string) (*User, error)
	GetUserByEmail(email string) (*User, error)
	DeleteUser(id string) error
	ListUsers() ([]User, error)
}

type UserService struct {
	store UserStore
}

func NewUserService(s UserStore) *UserService {
	return &UserService{store: s}
}

func Decode(r io.ReadCloser) (*User, error) {
	var u User
	err := json.NewDecoder(r).Decode(&u)
	if err != nil {
		return nil, fmt.Errorf("error decoding user: %w", err)
	}

	return &u, nil
}

func Unmarshal(data []byte) (*User, error) {
	var u User
	err := json.Unmarshal(data, &u)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling user: %w", err)
	}

	return &u, nil
}

func (u *UserService) CreateUser(id, name, email string) error {
	existingUser, err := u.store.GetUser(id)
	if err != nil {
		if errors.Is(err, &uerr.ErrorCannotInsert{}) {
			return fmt.Errorf("user already exists: %w", err)
		}
	}

	if existingUser != nil {
		return fmt.Errorf("user already exists")
	}

	return u.store.InsertUser(User{ID: id, Username: name, Email: email})
}

func (u *UserService) GetUser(id string) (*User, error) {
	user, err := u.store.GetUser(id)
	if err != nil {
		if errors.Is(err, &uerr.ErrorNotFound{}) {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("error getting user: %w", err)
	}

	return user, nil
}

func (u *UserService) GetUserByEmail(email string) (*User, error) {
	user, err := u.store.GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, &uerr.ErrorNotFound{}) {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("error getting user: %w", err)
	}

	return user, nil
}

func (u *UserService) DeleteUser(id string) error {
	return u.store.DeleteUser(id)
}

func (u *UserService) ListAllUsers() ([]User, error) {
	users, err := u.store.ListUsers()
	if err != nil {
		return nil, fmt.Errorf("error listing users: %w", err)
	}

	return users, nil
}

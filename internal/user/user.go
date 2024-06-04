package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fmdunlap/unhash/internal/uerr"
	"io"
	"log"
)

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type UserCache interface {
	GetUser(id string) (*User, error)
	GetUserByEmail(email string) (*User, error)
	SetUser(u *User) error
	ClearUser(u *User) error
}

type UserStore interface {
	InsertUser(u *User) error
	GetUser(id string) (*User, error)
	GetUserByEmail(email string) (*User, error)
	DeleteUser(id string) error
	ListUsers() ([]User, error)
}

type UserService struct {
	store UserStore
	cache UserCache
}

func NewUserService(s UserStore, c UserCache) *UserService {

	return &UserService{store: s, cache: c}
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

func (user *User) Validate() error {
	if user.ID == "" {
		return fmt.Errorf("user ID cannot be empty")
	}
	if user.Username == "" {
		return fmt.Errorf("username cannot be empty")
	}
	if user.Email == "" {
		return fmt.Errorf("email cannot be empty")
	}
	return nil
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

	user := &User{ID: id, Username: name, Email: email}
	err = u.store.InsertUser(user)
	if err != nil {
		return err
	}
	err = u.cache.SetUser(user)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserService) GetUser(id string) (*User, error) {
	user, err := u.cache.GetUser(id)
	if err == nil {
		return user, nil
	}

	if errors.Is(err, &uerr.ErrorNotFound{}) {
		log.Print("User not found in cache, getting from DB")
	} else {
		return nil, fmt.Errorf("unknown error when getting user: %w", err)
	}

	user, err = u.store.GetUser(id)
	if err != nil {
		if errors.Is(err, &uerr.ErrorNotFound{}) {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("error getting user: %w", err)
	}

	return user, nil
}

func (u *UserService) GetUserByEmail(email string) (*User, error) {
	user, err := u.cache.GetUserByEmail(email)
	if err == nil {
		return user, nil
	}

	if errors.Is(err, &uerr.ErrorNotFound{}) {
		log.Print("User not found in cache, getting from DB")
	} else {
		return nil, fmt.Errorf("unknown error when getting user: %w", err)
	}

	user, err = u.store.GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, &uerr.ErrorNotFound{}) {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("error getting user: %w", err)
	}

	return user, nil
}

func (u *UserService) DeleteUser(user *User) error {
	err := u.cache.ClearUser(user)
	if err != nil {
		log.Printf("Encountered error when clearing user from cache: %v", err)
	}
	return u.store.DeleteUser(user.ID)
}

func (u *UserService) ListAllUsers() ([]User, error) {
	users, err := u.store.ListUsers()
	if err != nil {
		return nil, fmt.Errorf("error listing users: %w", err)
	}

	return users, nil
}

package rediscache

import (
	"encoding/json"
	"fmt"
	"github.com/fmdunlap/unhash/internal/uerr"
	"github.com/fmdunlap/unhash/internal/user"
	"time"
)

const (
	USER_CACHE_TTL = 1 * time.Hour
)

func (r *RedisCache) userIdKey(userId string) string {
	return fmt.Sprintf("userid#%v", userId)
}

func (r *RedisCache) userEmailKey(userEmail string) string {
	return fmt.Sprintf("useremail#%v", userEmail)
}

func (r *RedisCache) GetUser(id string) (*user.User, error) {
	val, err := r.Client.Get(r.Context, r.userIdKey(id)).Result()
	if err != nil {
		return nil, fmt.Errorf("key not found %v: %w", id, &uerr.ErrorCannotDelete{Err: err})
	}

	u, err := user.Unmarshal([]byte(val))
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (r *RedisCache) GetUserByEmail(email string) (*user.User, error) {
	val, err := r.Client.Get(r.Context, r.userEmailKey(email)).Result()
	if err != nil {
		return nil, fmt.Errorf("key not found %v: %w", email, &uerr.ErrorCannotDelete{Err: err})
	}

	u, err := user.Unmarshal([]byte(val))
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (r *RedisCache) SetUser(u *user.User) error {
	marshaledUser, err := json.Marshal(u)
	if err != nil {
		return fmt.Errorf("cannot insert user %v into cache: %w", u, &uerr.ErrorCannotInsert{Err: err})
	}

	err = r.Client.Set(r.Context, r.userIdKey(u.ID), marshaledUser, USER_CACHE_TTL).Err()
	if err != nil {
		return fmt.Errorf("cannot insert user %v into cache: %w", u, &uerr.ErrorCannotInsert{Err: err})
	}

	err = r.Client.Set(r.Context, r.userEmailKey(u.Email), marshaledUser, USER_CACHE_TTL).Err()
	if err != nil {
		return fmt.Errorf("cannot insert user %v into cache: %w", u, &uerr.ErrorCannotInsert{Err: err})
	}

	return nil
}

func (r *RedisCache) ClearUser(u *user.User) error {
	err := r.Client.Del(r.Context, r.userIdKey(u.ID), r.userEmailKey(u.Email)).Err()
	if err != nil {
		return fmt.Errorf("cannot delete user %v redis keys: %w", u, &uerr.ErrorCannotDelete{Err: err})
	}
	return nil
}

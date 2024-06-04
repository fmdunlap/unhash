package rediscache

import (
	"encoding/json"
	"fmt"
	"github.com/fmdunlap/unhash/internal/hashjob"
	"github.com/fmdunlap/unhash/internal/uerr"
)

func (r *RedisCache) jobIdKey(jobId string) string {
	return fmt.Sprintf("userid#%v", jobId)
}

func (r *RedisCache) GetHashJob(id string) (*hashjob.HashJob, error) {
	val, err := r.Client.Get(r.Context, r.jobIdKey(id)).Result()
	if err != nil {
		return nil, fmt.Errorf("key not found %v: %w", id, &uerr.ErrorCannotDelete{Err: err})
	}

	hj, err := hashjob.Unmarshal([]byte(val))
	if err != nil {
		return nil, err
	}

	return hj, nil
}

func (r *RedisCache) SetHashJob(h hashjob.HashJob) error {
	marshaledJob, err := json.Marshal(h)
	if err != nil {
		return fmt.Errorf("cannot insert job %v into cache: %w", h, &uerr.ErrorCannotInsert{Err: err})
	}

	err = r.Client.Set(r.Context, r.jobIdKey(h.ID), marshaledJob, USER_CACHE_TTL).Err()
	if err != nil {
		return fmt.Errorf("cannot insert job %v into cache: %w", h, &uerr.ErrorCannotInsert{Err: err})
	}

	return nil
}

func (r *RedisCache) ClearHashJob(h hashjob.HashJob) error {
	err := r.Client.Del(r.Context, r.jobIdKey(h.ID)).Err()
	if err != nil {
		return fmt.Errorf("cannot delete job %v from cache: %w", h, &uerr.ErrorCannotDelete{Err: err})
	}

	return nil
}

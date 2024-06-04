package hashjob

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/fmdunlap/unhash/internal/uerr"
	"github.com/fmdunlap/unhash/internal/user"
	"github.com/google/uuid"
)

type HashJobStatus string

const (
	HashJobStatusPending HashJobStatus = "pending"
	HashJobStatusRunning HashJobStatus = "running"
	HashJobStatusDone    HashJobStatus = "done"
	HashJobStatusError   HashJobStatus = "error"
)

type HashJob struct {
	ID      string        `json:"id"`
	OwnerId string        `json:"ownerId"`
	Status  HashJobStatus `json:"status"`
	Hashes  []string      `json:"hash"`
}

type HashJobStore interface {
	InsertHashJob(h HashJob) error
	GetHashJob(id string) (*HashJob, error)
	DeleteHashJob(id string) error
}

type HashJobCache interface {
	GetHashJob(id string) (*HashJob, error)
	SetHashJob(h HashJob) error
	ClearHashJob(h HashJob) error
}

type HashJobService struct {
	store HashJobStore
	cache HashJobCache
}

func NewHashJobService(s HashJobStore, c HashJobCache) *HashJobService {
	return &HashJobService{store: s, cache: c}
}

func Unmarshal(data []byte) (*HashJob, error) {
	var hj HashJob
	err := json.Unmarshal(data, &hj)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling hashjob: %w", err)
	}

	return &hj, nil
}

func (h *HashJobService) CreateHashJob(hashes []string, owner *user.User) (string, error) {
	if len(hashes) == 0 {
		return "", errors.New("hashes cannot be empty")
	}
	if owner == nil {
		return "", errors.New("owner cannot be nil")
	}
	if owner.Validate() != nil {
		return "", fmt.Errorf("invalid owner: %w", owner.Validate())
	}

	hj := HashJob{
		ID:      uuid.New().String(),
		OwnerId: owner.ID,
		Status:  HashJobStatusPending,
		Hashes:  hashes,
	}

	err := h.store.InsertHashJob(hj)
	if err != nil {
		if errors.Is(err, &uerr.ErrorCannotInsert{}) {
			return "", fmt.Errorf("hashjob already exists: %w", err)
		}
		return "", err
	}

	err = h.cache.SetHashJob(hj)
	if err != nil {
		return "", err
	}

	return hj.ID, nil
}

func (h *HashJobService) GetHashJob(id string) (*HashJob, error) {
	hj, err := h.cache.GetHashJob(id)
	if err == nil {
		return hj, nil
	}

	log.Printf("Cache miss for hashjob %v", id)

	hj, err = h.store.GetHashJob(id)
	if err != nil {
		if errors.Is(err, &uerr.ErrorNotFound{}) {
			return nil, fmt.Errorf("hashjob not found: %w", err)
		}
		return nil, err
	}

	err = h.cache.SetHashJob(*hj)
	if err != nil {
		return nil, err
	}

	return hj, nil
}

func (h *HashJobService) DeleteHashJob(id string) error {
	err := h.store.DeleteHashJob(id)
	if err != nil {
		if errors.Is(err, &uerr.ErrorCannotDelete{}) {
			return fmt.Errorf("hashjob not found: %w", err)
		}
		return nil
	}

	return nil
}

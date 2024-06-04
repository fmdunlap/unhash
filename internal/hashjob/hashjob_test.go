package hashjob

import (
	"encoding/json"
	"github.com/fmdunlap/unhash/internal/uerr"
	"github.com/fmdunlap/unhash/internal/user"
	"testing"
)

// Helpers

func checkJobInMap(t *testing.T, m map[string]HashJob, id string, hj HashJob) {
	_, ok := m[id]
	if !ok {
		t.Errorf("checkJobInMap() got = %v, want %v", ok, true)
	}
	if m[id].ID != hj.ID {
		t.Errorf("checkJobInMap() got = %v, want %v", m[id].ID, hj.ID)
	}
	if m[id].OwnerId != hj.OwnerId {
		t.Errorf("checkJobInMap() got = %v, want %v", m[id].OwnerId, hj.OwnerId)
	}
	if m[id].Status != hj.Status {
		t.Errorf("checkJobInMap() got = %v, want %v", m[id].Status, hj.Status)
	}
	if len(m[id].Hashes) != len(hj.Hashes) {
		t.Errorf("checkJobInMap() got = %v, want %v", len(m[id].Hashes), len(hj.Hashes))
	}
	for i, hash := range m[id].Hashes {
		if hash != hj.Hashes[i] {
			t.Errorf("checkJobInMap() got = %v, want %v", hash, hj.Hashes[i])
		}
	}
}

// MockHashJobStore implements HashJobStore

type MockHashJobStore struct {
	HashJobs map[string]HashJob
}

func (m *MockHashJobStore) InsertHashJob(h HashJob) error {
	m.HashJobs[h.ID] = h
	return nil
}

func (m *MockHashJobStore) GetHashJob(id string) (*HashJob, error) {
	h, ok := m.HashJobs[id]
	if !ok {
		return nil, &uerr.ErrorNotFound{}
	}
	return &h, nil
}

func (m *MockHashJobStore) DeleteHashJob(id string) error {
	_, ok := m.HashJobs[id]
	if !ok {
		return &uerr.ErrorCannotDelete{}
	}
	delete(m.HashJobs, id)
	return nil
}

// MockHashJobCache implements HashJobCache

type MockHashJobCache struct {
	HashJobs map[string]HashJob
}

func (m *MockHashJobCache) GetHashJob(id string) (*HashJob, error) {
	h, ok := m.HashJobs[id]
	if !ok {
		return nil, &uerr.ErrorNotFound{}
	}
	return &h, nil
}

func (m *MockHashJobCache) SetHashJob(h HashJob) error {
	m.HashJobs[h.ID] = h
	return nil
}

func (m *MockHashJobCache) ClearHashJob(h HashJob) error {
	_, ok := m.HashJobs[h.ID]
	if !ok {
		return &uerr.ErrorCannotDelete{}
	}
	delete(m.HashJobs, h.ID)
	return nil
}

// HashJob Tests

func TestUnmarshal(t *testing.T) {
	t.Run("Test Unmarshal", func(t *testing.T) {
		premarshaledHashJob := HashJob{
			ID:      "test",
			OwnerId: "test",
			Status:  HashJobStatusPending,
			Hashes:  []string{"test"},
		}

		data, err := json.Marshal(premarshaledHashJob)
		if err != nil {
			t.Errorf("Marshal() error = %v", err)
			return
		}

		hj, err := Unmarshal(data)
		if err != nil {
			t.Errorf("Unmarshal() error = %v", err)
			return
		}

		if hj.ID != "test" {
			t.Errorf("Unmarshal() ID = %v, want test", hj.ID)
		}

		if hj.OwnerId != "test" {
			t.Errorf("Unmarshal() OwnerId = %v, want test", hj.OwnerId)
		}

		if hj.Status != HashJobStatusPending {
			t.Errorf("Unmarshal() Status = %v, want pending", hj.Status)
		}

		if len(hj.Hashes) != 1 {
			t.Errorf("Unmarshal() Hashes = %v, want 1", len(hj.Hashes))
		}

		if hj.Hashes[0] != "test" {
			t.Errorf("Unmarshal() Hashes[0] = %v, want test", hj.Hashes[0])
		}
	})
}

// HashJobService Tests

func TestHashJobService_CreateHashJob(t *testing.T) {

	testOwner := &user.User{
		ID:       "test",
		Username: "test",
		Email:    "test",
	}

	type args struct {
		hashes []string
		owner  *user.User
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test CreateHashJob",
			args: args{
				hashes: []string{"test"},
				owner:  testOwner,
			},
			wantErr: false,
		},
		{
			name: "Test CreateHashJob with multiple hashes",
			args: args{
				hashes: []string{"test", "test2"},
				owner:  testOwner,
			},
			wantErr: false,
		},
		{
			name: "Test CreateHashJob with no hashes",
			args: args{
				hashes: []string{},
				owner:  testOwner,
			},
			wantErr: true,
		},
		{
			name: "Test CreateHashJob with no owner",
			args: args{
				hashes: []string{"test"},
				owner:  nil,
			},
			wantErr: true,
		},
		{
			name: "Test CreateHashJob with invalid owner",
			args: args{
				hashes: []string{"test"},
				owner: &user.User{
					ID:       "",
					Username: "test",
					Email:    "",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hashJobStoreMap := make(map[string]HashJob)
			hashJobCacheMap := make(map[string]HashJob)
			h := &HashJobService{
				store: &MockHashJobStore{HashJobs: hashJobStoreMap},
				cache: &MockHashJobCache{HashJobs: hashJobCacheMap},
			}
			got, err := h.CreateHashJob(tt.args.hashes, tt.args.owner)

			if tt.wantErr {
				if err == nil {
					t.Errorf("CreateHashJob() got = %v, want error", got)
				}
				return
			}

			if err != nil {
				t.Errorf("CreateHashJob() error = %v", err)
			}

			if got == "" && !tt.wantErr {
				t.Errorf("CreateHashJob() got = %v, want not empty", got)
			}

			// Check Store
			checkJobInMap(t, hashJobStoreMap, got, HashJob{
				ID:      got,
				OwnerId: tt.args.owner.ID,
				Status:  HashJobStatusPending,
				Hashes:  tt.args.hashes,
			})

			// Check Cache
			checkJobInMap(t, hashJobCacheMap, got, hashJobStoreMap[got])
		})
	}
}

func TestHashJobService_GetHashJob(t *testing.T) {
	testOwner := &user.User{
		ID:       "test",
		Username: "test",
		Email:    "test",
	}

	type args struct {
		id string
	}

	tests := []struct {
		name       string
		args       args
		want       *HashJob
		wantErr    bool
		before     func(*HashJobService)
		checkCache bool
	}{
		{
			name: "Test GetHashJob",
			args: args{
				id: "test",
			},
			want: &HashJob{
				ID:      "test",
				OwnerId: testOwner.ID,
				Status:  HashJobStatusPending,
				Hashes:  []string{"test"},
			},
			wantErr: false,
			before: func(h *HashJobService) {
				h.store.InsertHashJob(HashJob{
					ID:      "test",
					OwnerId: testOwner.ID,
					Status:  HashJobStatusPending,
					Hashes:  []string{"test"},
				})
			},
			checkCache: true,
		},
		{
			name: "Test GetHashJob with no job",
			args: args{
				id: "test",
			},
			want:       nil,
			wantErr:    true,
			before:     nil,
			checkCache: false,
		},
		{
			name: "Test GetHashJob with no job in cache",
			args: args{
				id: "test",
			},
			want: &HashJob{
				ID:      "test",
				OwnerId: testOwner.ID,
				Status:  HashJobStatusPending,
				Hashes:  []string{"test"},
			},
			wantErr: false,
			before: func(h *HashJobService) {
				hj := HashJob{
					ID:      "test",
					OwnerId: testOwner.ID,
					Status:  HashJobStatusPending,
					Hashes:  []string{"test"},
				}

				h.store.InsertHashJob(hj)
				h.cache.ClearHashJob(hj)
			},
			checkCache: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storeMap := make(map[string]HashJob)
			cacheMap := make(map[string]HashJob)
			h := &HashJobService{
				store: &MockHashJobStore{HashJobs: storeMap},
				cache: &MockHashJobCache{HashJobs: cacheMap},
			}

			if tt.before != nil {
				tt.before(h)
			}

			if tt.checkCache {
				if len(cacheMap) != 0 {
					t.Errorf("GetHashJob() cacheMap = %v, want 0", len(cacheMap))
				}
			}

			got, err := h.GetHashJob(tt.args.id)

			if tt.wantErr {
				if err == nil {
					t.Errorf("GetHashJob() got = %v, want error", got)
				}
				return
			}

			if err != nil {
				t.Errorf("GetHashJob() error = %v", err)
				return
			}

			checkJobInMap(t, storeMap, tt.args.id, *tt.want)

			if tt.checkCache {
				checkJobInMap(t, cacheMap, tt.args.id, *tt.want)
			}
		})
	}
}

func TestHashJobService_DeleteHashJob(t *testing.T) {
	testOwner := &user.User{
		ID:       "test",
		Username: "test",
		Email:    "test",
	}

	type args struct {
		id string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
		before  func(*HashJobService)
	}{
		{
			name: "Test DeleteHashJob",
			args: args{
				id: "test",
			},
			wantErr: false,
			before: func(h *HashJobService) {
				h.store.InsertHashJob(HashJob{
					ID:      "test",
					OwnerId: testOwner.ID,
					Status:  HashJobStatusPending,
					Hashes:  []string{"test"},
				})
			},
		},
		{
			name: "Test DeleteHashJob with no job",
			args: args{
				id: "test",
			},
			wantErr: true,
			before:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storeMap := make(map[string]HashJob)
			cacheMap := make(map[string]HashJob)

			h := &HashJobService{
				store: &MockHashJobStore{HashJobs: storeMap},
				cache: &MockHashJobCache{HashJobs: cacheMap},
			}

			if tt.before != nil {
				tt.before(h)
			}

			err := h.DeleteHashJob(tt.args.id)

			if tt.wantErr {
				if err == nil {
					t.Errorf("DeleteHashJob() got = %v, want error", err)
				}
				return
			}

			if err != nil {
				t.Errorf("DeleteHashJob() error = %v", err)
				return
			}

			if _, ok := storeMap[tt.args.id]; ok {
				t.Errorf("DeleteHashJob() storeMap = %v, want empty", storeMap)
			}
			if _, ok := cacheMap[tt.args.id]; ok {
				t.Errorf("DeleteHashJob() cacheMap = %v, want empty", cacheMap)
			}

			_, err = h.store.GetHashJob(tt.args.id)
			if err == nil {
				t.Errorf("DeleteHashJob() got = %v, want error", err)
			}

		})
	}
}

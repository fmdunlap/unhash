package sqlite

import (
	"database/sql"
	"github.com/fmdunlap/unhash/internal/hashjob"
	"testing"
)

func TestSqliteStore_InsertHashJob(t *testing.T) {
	type fields struct {
		sq3 *sql.DB
	}
	type args struct {
		h hashjob.HashJob
	}

	testHashJob := hashjob.HashJob{
		ID:      "test",
		OwnerId: "test",
		Status:  hashjob.HashJobStatusPending,
		Hashes:  []string{"test"},
	}

	testHashJobTwo := hashjob.HashJob{
		ID:      "test2",
		OwnerId: "test2",
		Status:  hashjob.HashJobStatusPending,
		Hashes:  []string{"test21", "test22"},
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		expect  []hashjob.HashJob
		before  func(s *SqliteStore)
	}{
		{
			name: "Test InsertHashJob",
			fields: fields{
				sq3: CreateTestDb(),
			},
			args: args{
				h: testHashJob,
			},
			wantErr: false,
			expect: []hashjob.HashJob{
				testHashJob,
			},
		},
		{
			name: "Test InsertHashJob with multiple",
			fields: fields{
				sq3: CreateTestDb(),
			},
			args: args{
				h: testHashJob,
			},
			wantErr: false,
			expect: []hashjob.HashJob{
				testHashJob,
				testHashJobTwo,
			},
			before: func(s *SqliteStore) {
				err := s.InsertHashJob(testHashJobTwo)
				if err != nil {
					t.Errorf("Error inserting hashjob: %v", err)
				}
			},
		},
		{
			name: "Test InsertHashJob with empty ID",
			fields: fields{
				sq3: CreateTestDb(),
			},
			args: args{
				h: hashjob.HashJob{
					ID:      "",
					OwnerId: "test",
					Status:  hashjob.HashJobStatusPending,
					Hashes:  []string{"test"},
				},
			},
			wantErr: true,
			expect:  []hashjob.HashJob{},
		},
		{
			name: "Test InsertHashJob with empty OwnerId",
			fields: fields{
				sq3: CreateTestDb(),
			},
			args: args{
				h: hashjob.HashJob{
					ID:      "test",
					OwnerId: "",
					Status:  hashjob.HashJobStatusPending,
					Hashes:  []string{"test"},
				},
			},
			wantErr: true,
			expect:  []hashjob.HashJob{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SqliteStore{
				sq3: tt.fields.sq3,
			}

			if tt.before != nil {
				tt.before(s)
			}

			if err := s.InsertHashJob(tt.args.h); (err != nil) != tt.wantErr {
				t.Errorf("InsertHashJob() error = %v, wantErr %v", err, tt.wantErr)
			}

			countResult, err := s.sq3.Query("select count(*) from hashjobs")
			if err != nil {
				t.Errorf("Error querying hashjobs: %v", err)
			}
			defer countResult.Close()

			var count int
			for countResult.Next() {
				err = countResult.Scan(&count)
				if err != nil {
					t.Errorf("Error scanning hashjobs: %v", err)
				}
			}
			if count != len(tt.expect) {
				t.Errorf("Expected %d hashjobs, got %d", len(tt.expect), count)
			}

			for _, h := range tt.expect {
				row, err := s.sq3.Query("select * from hashjobs where id = ?", h.ID)
				if err != nil {
					t.Errorf("Error querying hashjob: %v", err)
				}
				defer row.Close()

				var id string
				var data []byte
				for row.Next() {
					err = row.Scan(&id, &data)
					if err != nil {
						t.Errorf("Error scanning hashjob: %v", err)
					}

					h2, err := hashjob.Unmarshal(data)
					if err != nil {
						t.Errorf("Error unmarshaling hashjob: %v", err)
					}

					if h2.ID != h.ID {
						t.Errorf("Expected ID %s, got %s", h.ID, h2.ID)
					}
					if h2.OwnerId != h.OwnerId {
						t.Errorf("Expected OwnerId %s, got %s", h.OwnerId, h2.OwnerId)
					}
					if h2.Status != h.Status {
						t.Errorf("Expected Status %s, got %s", h.Status, h2.Status)
					}
					if len(h2.Hashes) != len(h.Hashes) {
						t.Errorf("Expected %d hashes, got %d", len(h.Hashes), len(h2.Hashes))
					}
					for i, hash := range h.Hashes {
						if h2.Hashes[i] != hash {
							t.Errorf("Expected hash %s, got %s", hash, h2.Hashes[i])
						}
					}
				}
			}
		})
	}
}

func TestSqliteStore_GetHashJob(t *testing.T) {
	type fields struct {
		sq3 *sql.DB
	}
	type args struct {
		id string
	}

	testHashJob := hashjob.HashJob{
		ID:      "test",
		OwnerId: "test",
		Status:  hashjob.HashJobStatusPending,
		Hashes:  []string{"test"},
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *hashjob.HashJob
		wantErr bool
		before  func(s *SqliteStore)
	}{
		{
			name: "Test GetHashJob",
			fields: fields{
				sq3: CreateTestDb(),
			},
			args: args{
				id: "test",
			},
			want:    &testHashJob,
			wantErr: false,
			before: func(s *SqliteStore) {
				err := s.InsertHashJob(testHashJob)
				if err != nil {
					t.Errorf("Error inserting hashjob: %v", err)
				}
			},
		},
		{
			name: "Test GetHashJob with non-existent ID",
			fields: fields{
				sq3: CreateTestDb(),
			},
			args: args{
				id: "test",
			},
			want:    nil,
			wantErr: true,
			before:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SqliteStore{
				sq3: tt.fields.sq3,
			}

			if tt.before != nil {
				tt.before(s)
			}

			got, err := s.GetHashJob(tt.args.id)
			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
				return
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("GetHashJob() error = %v, wantErr %v", err, tt.wantErr)
			}

			if got.ID != tt.want.ID {
				t.Errorf("GetHashJob() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSqliteStore_DeleteHashJob(t *testing.T) {
	type fields struct {
		sq3 *sql.DB
	}
	type args struct {
		id string
	}

	testHashJob := hashjob.HashJob{
		ID:      "test",
		OwnerId: "test",
		Status:  hashjob.HashJobStatusPending,
		Hashes:  []string{"test"},
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		before  func(s *SqliteStore)
		expect  []hashjob.HashJob
	}{
		{
			name: "Test DeleteHashJob",
			fields: fields{
				sq3: CreateTestDb(),
			},
			args: args{
				id: "test",
			},
			wantErr: false,
			before: func(s *SqliteStore) {
				err := s.InsertHashJob(testHashJob)
				if err != nil {
					t.Errorf("Error inserting hashjob: %v", err)
				}
			},
			expect: []hashjob.HashJob{},
		},
		{
			name: "Test DeleteHashJob with non-existent ID",
			fields: fields{
				sq3: CreateTestDb(),
			},
			args: args{
				id: "test",
			},
			wantErr: false,
			before:  nil,
			expect:  []hashjob.HashJob{},
		},
		{
			name: "Test DeleteHashJob with multiple",
			fields: fields{
				sq3: CreateTestDb(),
			},
			args: args{
				id: "test",
			},
			wantErr: false,
			before: func(s *SqliteStore) {
				err := s.InsertHashJob(testHashJob)
				if err != nil {
					t.Errorf("Error inserting hashjob: %v", err)
				}
				err = s.InsertHashJob(hashjob.HashJob{
					ID:      "test2",
					OwnerId: "test2",
					Status:  hashjob.HashJobStatusPending,
					Hashes:  []string{"test21", "test22"},
				})
				if err != nil {
					t.Errorf("Error inserting hashjob: %v", err)
				}
			},
			expect: []hashjob.HashJob{
				{
					ID:      "test2",
					OwnerId: "test2",
					Status:  hashjob.HashJobStatusPending,
					Hashes:  []string{"test21", "test22"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SqliteStore{
				sq3: tt.fields.sq3,
			}

			if tt.before != nil {
				tt.before(s)
			}

			if err := s.DeleteHashJob(tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("DeleteHashJob() error = %v, wantErr %v", err, tt.wantErr)
			}

			countResult, err := s.sq3.Query("select count(*) from hashjobs")
			if err != nil {
				t.Errorf("Error querying hashjobs: %v", err)
			}
			defer countResult.Close()

			var count int
			for countResult.Next() {
				err = countResult.Scan(&count)
				if err != nil {
					t.Errorf("Error scanning hashjobs: %v", err)
				}
			}
			if count != len(tt.expect) {
				t.Errorf("Expected 0 hashjobs, got %d", count)
			}
		})
	}
}

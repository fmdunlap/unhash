package rediscache

import (
	"context"
	"github.com/fmdunlap/unhash/internal/hashjob"
	"github.com/redis/go-redis/v9"
	"testing"
)

func setupJobsInCache(r *RedisCache, jobs []hashjob.HashJob) error {
	for _, j := range jobs {
		err := r.SetHashJob(j)
		if err != nil {
			return err
		}
	}
	return nil
}

func clearCache(r *RedisCache) {
	r.Client.FlushAll(r.Context)
}

func TestRedisCache_SetHashJob(t *testing.T) {
	type fields struct {
		Client  *redis.Client
		Context context.Context
	}
	type args struct {
		h hashjob.HashJob
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Test SetHashJob",
			fields: fields{
				Client:  redis.NewClient(&redis.Options{}),
				Context: context.Background(),
			},
			args: args{
				h: hashjob.HashJob{
					ID:      "test",
					OwnerId: "test",
					Status:  hashjob.HashJobStatusPending,
					Hashes:  []string{"test"},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RedisCache{
				Client:  tt.fields.Client,
				Context: tt.fields.Context,
			}
			if err := r.SetHashJob(tt.args.h); (err != nil) != tt.wantErr {
				t.Errorf("SetHashJob() error = %v, wantErr %v", err, tt.wantErr)
				if tt.wantErr {
					return
				}
			}

			gotJob, err := tt.fields.Client.Get(tt.fields.Context, r.jobIdKey(tt.args.h.ID)).Result()
			if err != nil {
				t.Errorf("SetHashJob() error = %v", err)
			}

			job, err := hashjob.Unmarshal([]byte(gotJob))
			if err != nil {
				t.Errorf("SetHashJob() error = %v", err)
			}

			if job.ID != tt.args.h.ID {
				t.Errorf("SetHashJob() got = %v, want %v", job, tt.args.h)
			}
			if job.Status != tt.args.h.Status {
				t.Errorf("SetHashJob() got = %v, want %v", job, tt.args.h)
			}
			if job.OwnerId != tt.args.h.OwnerId {
				t.Errorf("SetHashJob() got = %v, want %v", job, tt.args.h)
			}
			if job.Hashes[0] != tt.args.h.Hashes[0] {
				t.Errorf("SetHashJob() got = %v, want %v", job, tt.args.h)
			}

			// Clean up
			clearCache(r)
		})
	}
}

func TestRedisCache_GetHashJob(t *testing.T) {
	type fields struct {
		Client  *redis.Client
		Context context.Context
	}
	type args struct {
		id string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *hashjob.HashJob
		wantErr bool
	}{
		{
			name: "Test GetHashJob",
			fields: fields{
				Client:  redis.NewClient(&redis.Options{}),
				Context: context.Background(),
			},
			args: args{
				id: "test",
			},
			want: &hashjob.HashJob{
				ID:      "test",
				OwnerId: "test",
				Status:  hashjob.HashJobStatusPending,
				Hashes:  []string{"test"},
			},
			wantErr: false,
		},
		{
			name: "Test GetHashJob with no job",
			fields: fields{
				Client:  redis.NewClient(&redis.Options{}),
				Context: context.Background(),
			},
			args: args{
				id: "test",
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RedisCache{
				Client:  tt.fields.Client,
				Context: tt.fields.Context,
			}

			if tt.want != nil {
				err := setupJobsInCache(r, []hashjob.HashJob{*tt.want})
				if err != nil {
					t.Errorf("GetHashJob() error = %v", err)
				}
			}

			got, err := r.GetHashJob(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetHashJob() error = %v, wantErr %v, got = %v", err, tt.wantErr, got)
				return
			}

			if got == nil && tt.want == nil {
				return
			}

			if got.ID != tt.want.ID {
				t.Errorf("GetHashJob() got = %v, want %v", got, tt.want)
			}
			if got.Status != tt.want.Status {
				t.Errorf("GetHashJob() got = %v, want %v", got, tt.want)
			}
			if got.OwnerId != tt.want.OwnerId {
				t.Errorf("GetHashJob() got = %v, want %v", got, tt.want)
			}

			clearCache(r)
		})
	}
}

func TestRedisCache_ClearHashJob(t *testing.T) {
	type fields struct {
		Client  *redis.Client
		Context context.Context
	}
	type args struct {
		h hashjob.HashJob
	}

	tests := []struct {
		name           string
		fields         fields
		args           args
		createJobFirst bool
		wantErr        bool
	}{
		{
			name: "Test ClearHashJob",
			fields: fields{
				Client:  redis.NewClient(&redis.Options{}),
				Context: context.Background(),
			},
			args: args{
				h: hashjob.HashJob{
					ID:      "test",
					OwnerId: "test",
					Status:  hashjob.HashJobStatusPending,
					Hashes:  []string{"test"},
				},
			},
			createJobFirst: true,
			wantErr:        false,
		},
		{
			name: "Test ClearHashJob with no job",
			fields: fields{
				Client:  redis.NewClient(&redis.Options{}),
				Context: context.Background(),
			},
			args: args{
				h: hashjob.HashJob{
					ID:      "test",
					OwnerId: "test",
					Status:  hashjob.HashJobStatusPending,
					Hashes:  []string{"test"},
				},
			},
			createJobFirst: false,
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RedisCache{
				Client:  tt.fields.Client,
				Context: tt.fields.Context,
			}

			if tt.createJobFirst {
				err := r.SetHashJob(tt.args.h)
				if err != nil {
					t.Errorf("ClearHashJob() error = %v", err)
				}
			}

			if err := r.ClearHashJob(tt.args.h); (err != nil) != tt.wantErr {
				t.Errorf("ClearHashJob() error = %v, wantErr %v", err, tt.wantErr)
				if tt.wantErr {
					return
				}
			}

			_, err := tt.fields.Client.Get(tt.fields.Context, r.jobIdKey(tt.args.h.ID)).Result()
			if err == nil {
				t.Errorf("ClearHashJob() error = %v", err)
			}

			clearCache(r)
		})
	}
}

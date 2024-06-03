package rediscache

import "github.com/redis/go-redis/v9"

import "context"

type RedisCache struct {
	Client  *redis.Client
	Context context.Context
}

func NewRedisClient(address, password string, clearOnStartup bool) *RedisCache {
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       0,
	})

	if clearOnStartup {
		rdb.FlushAll(ctx)
	}

	return &RedisCache{
		Client: rdb,
		// Probably not ideal, but I don't have a deep enough understanding of Go's context just yet.
		Context: ctx,
	}
}

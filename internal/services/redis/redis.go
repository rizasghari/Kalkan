package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rizasghari/kalkan/internal/cfg"
)

type RedisService struct {
	client *redis.Client
}

func Initialize(cfg *cfg.Configuration) *RedisService {
	return &RedisService{
		client: redis.NewClient(&redis.Options{
			Addr:     cfg.Redis.Url,
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.DB,
		}),
	}
}

func (rs *RedisService) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	return rs.client.Set(ctx, key, value, expiration).Err()
}

func (rs *RedisService) Get(ctx context.Context, key string) (any, error) {
	return rs.client.Get(ctx, key).Result()
}

func (rs *RedisService) Delete(ctx context.Context, key string) error {
	return rs.client.Del(ctx, key).Err()
}

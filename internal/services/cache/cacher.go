package cache

import (
	"context"
	"time"
)

type Cacher interface {
	Set(ctx context.Context, string, value any, expiration time.Duration) error
	Get(ctx context.Context, key string) (any, error)
	Delete(ctx context.Context, key string) error
}

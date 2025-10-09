package interfaces

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type CacheManager interface {
	Set(
		ctx context.Context,
		key string,
		value any,
		expiration time.Duration,
	) *redis.StatusCmd
	Get(ctx context.Context, key string) *redis.StringCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
}

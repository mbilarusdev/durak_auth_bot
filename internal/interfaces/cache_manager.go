package interfaces

import (
	"context"
	"time"
)

type CacheManager interface {
	Set(
		ctx context.Context,
		key string,
		value any,
		expiration time.Duration,
	) CacheStatusCmd
	Get(ctx context.Context, key string) CacheStringCmd
	Del(ctx context.Context, keys ...string) CacheIntCmd
}

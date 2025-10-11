package adapter

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/mbilarusdev/durak_auth_bot/internal/interfaces"
)

type AdapterCacheManager struct {
	manager *redis.Client
}

func NewCacheManagerAdapter(manager *redis.Client) interfaces.CacheManager {
	return &AdapterCacheManager{manager: manager}
}

func (c *AdapterCacheManager) Set(
	ctx context.Context,
	key string,
	value interface{},
	expiration time.Duration,
) interfaces.CacheStatusCmd {
	cmd := c.manager.Set(ctx, key, value, expiration)

	return cmd
}

func (c *AdapterCacheManager) Get(
	ctx context.Context, key string,
) interfaces.CacheStringCmd {
	cmd := c.manager.Get(ctx, key)

	return cmd
}

func (c *AdapterCacheManager) Del(
	ctx context.Context, keys ...string,
) interfaces.CacheIntCmd {
	cmd := c.manager.Del(ctx, keys...)

	return cmd
}

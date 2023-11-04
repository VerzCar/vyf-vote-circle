package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type Client interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Get(ctx context.Context, key string) *redis.StringCmd
	FlushDB(ctx context.Context) *redis.StatusCmd
	Exists(ctx context.Context, keys ...string) *redis.IntCmd
	ZAdd(ctx context.Context, key string, members ...*redis.Z) *redis.IntCmd
	ZRevRangeWithScores(ctx context.Context, key string, start int64, stop int64) *redis.ZSliceCmd
	ZScore(ctx context.Context, key string, member string) *redis.FloatCmd
	ZCount(ctx context.Context, key string, min string, max string) *redis.IntCmd
	Pipelined(ctx context.Context, fn func(redis.Pipeliner) error) ([]redis.Cmder, error)
}

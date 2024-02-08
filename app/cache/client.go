package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type Client interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Get(ctx context.Context, key string) *redis.StringCmd
	HSet(ctx context.Context, key string, values ...interface{}) *redis.IntCmd
	HGet(ctx context.Context, key string, field string) *redis.StringCmd
	HGetAll(ctx context.Context, key string) *redis.StringStringMapCmd
	HDel(ctx context.Context, key string, fields ...string) *redis.IntCmd
	FlushDB(ctx context.Context) *redis.StatusCmd
	Exists(ctx context.Context, keys ...string) *redis.IntCmd
	ZAdd(ctx context.Context, key string, members ...*redis.Z) *redis.IntCmd
	ZRem(ctx context.Context, key string, members ...interface{}) *redis.IntCmd
	ZRevRangeWithScores(ctx context.Context, key string, start int64, stop int64) *redis.ZSliceCmd
	ZScore(ctx context.Context, key string, member string) *redis.FloatCmd
	ZRevRank(ctx context.Context, key string, member string) *redis.IntCmd
	Pipelined(ctx context.Context, fn func(redis.Pipeliner) error) ([]redis.Cmder, error)
}

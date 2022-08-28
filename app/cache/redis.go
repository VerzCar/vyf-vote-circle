package cache

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"gitlab.vecomentman.com/libs/logger"
	"gitlab.vecomentman.com/vote-your-face/service/vote_circle/api/model"
	"gitlab.vecomentman.com/vote-your-face/service/vote_circle/app/config"
	"time"
)

type RedisCache interface {
	UpdateRanking(
		ctx context.Context,
		circleId int64,
		identityId model.UserIdentityId,
		votes int64,
	) error
}

type redisCache struct {
	redis  *redis.Client
	config *config.Config
	log    logger.Logger
}

func NewRedisCache(
	redis *redis.Client,
	config *config.Config,
	log logger.Logger,
) RedisCache {
	return &redisCache{
		redis:  redis,
		config: config,
		log:    log,
	}
}

type Entry struct {
	Val    string
	Exists bool
}

type EntryNumber struct {
	Val    int64
	Exists bool
}

type HEntry struct {
	Exists bool
}

// setJson converts the given value as JSON into the cache
// with the given key.
func (c *redisCache) setJson(ctx context.Context, key string, value interface{}, t time.Duration) error {
	encodedData, err := json.Marshal(value)

	if err != nil {
		return err
	}

	return c.redis.Set(ctx, key, encodedData, t).Err()
}

// getJson gets the entry from the given key
// as JSON format and Unmarshal it to the given destination
// interface structure type
func (c *redisCache) getJson(ctx context.Context, key string, dest interface{}) error {
	entry := c.redis.Get(ctx, key)

	switch {
	case entry.Err() == redis.Nil:
		return entry.Err()
	case entry.Err() != nil:
		return entry.Err()
	default:
		err := json.Unmarshal([]byte(entry.Val()), dest)
		return err
	}
}

// get the entry from the given key
// as Entry. If the entry does not exist, the Entry
// has the flag Exists set false and no error will be returned.
// Or if an error happens, an error will be returned.
// Otherwise, the Value and the Exists flag will be returned, without
// any error.
func (c *redisCache) get(ctx context.Context, key string) (Entry, error) {
	entry := c.redis.Get(ctx, key)

	sReturn := Entry{Exists: false}

	switch {
	case entry.Err() == redis.Nil:
		return sReturn, nil
	case entry.Err() != nil:
		return sReturn, entry.Err()
	default:
		sReturn.Exists = true
		sReturn.Val = entry.Val()
		return sReturn, nil
	}
}

func (c *redisCache) set(ctx context.Context, key string, value interface{}, t time.Duration) error {
	err := c.redis.Set(ctx, key, value, t).Err()
	return err
}

func (c *redisCache) getIndexInList(
	ctx context.Context,
	key string,
	value string,
) (EntryNumber, error) {
	entry := c.redis.LPos(ctx, key, value, redis.LPosArgs{})

	result := EntryNumber{Exists: false}

	switch {
	case entry.Err() == redis.Nil:
		return result, nil
	case entry.Err() != nil:
		return result, entry.Err()
	default:
		result.Exists = true
		result.Val = entry.Val()
		return result, nil
	}
}

func (c *redisCache) pushToListEnd(
	ctx context.Context,
	key string,
	values ...interface{},
) (int64, error) {
	entry := c.redis.RPush(ctx, key, values)

	switch {
	case entry.Err() != nil:
		return 0, entry.Err()
	default:
		return entry.Val(), nil
	}
}

// FlushAll the cache and flush the db
func (c *redisCache) FlushAll() error {
	ctx := context.Background()
	return c.redis.FlushDB(ctx).Err()
}